package aiflow

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"paas-ai/toolkit/aiflow/node"
	"paas-ai/toolkit/aiflow/properties"
	"sync"
	"time"
)

// LogicFlowEngine 是主流程执行器，负责控制流程图中各节点的执行。
type LogicFlowEngine struct {
	nodes               map[string]node.Node     // 所有注册的节点
	adjacencyList       map[string][]*node.Edge  // 邻接表表示的边（节点间连接）
	pendingPredecessors sync.Map                 // 使用 sync.Map 来优化并发控制
	input               map[string]interface{}   // 输入上下文
	output              map[string]interface{}   // 输出上下文
	loopIndexes         map[string]int           // 每个 loopVariable 当前索引
	loopValuesMap       map[string][]interface{} // 每个 loopVariable 对应的迭代集合
	startNodeID         string                   // 起始节点 ID
	mu                  sync.Mutex               // 保护队列的并发访问
}

// NewLogicFlowEngine 创建并初始化流程引擎实例
func NewLogicFlowEngine() *LogicFlowEngine {
	return &LogicFlowEngine{
		nodes:               make(map[string]node.Node),
		adjacencyList:       make(map[string][]*node.Edge),
		pendingPredecessors: sync.Map{},
		input:               make(map[string]interface{}),
		output:              make(map[string]interface{}),
		loopIndexes:         make(map[string]int),
		loopValuesMap:       make(map[string][]interface{}),
	}
}

// SetInput 设置输入参数
func (e *LogicFlowEngine) SetInput(name string, value interface{}) {
	e.input[name] = value
}

// SetOutput 设置输出参数
func (e *LogicFlowEngine) SetOutput(name string, value interface{}) {
	e.output[name] = value
}

// RegisterNode 注册节点到流程图
func (e *LogicFlowEngine) RegisterNode(nd node.Node) {
	e.nodes[nd.ID()] = nd
	e.adjacencyList[nd.ID()] = []*node.Edge{}
	if e.startNodeID == "" {
		e.startNodeID = nd.ID()
	}
}

// RegisterEdge 注册边（连接关系）
func (e *LogicFlowEngine) RegisterEdge(edge *node.Edge) {
	e.adjacencyList[edge.SourceNodeId] = append(e.adjacencyList[edge.SourceNodeId], edge)
	if edge.LoopType != "2" { // 非 LoopEnd 节点才计入前置数量
		count, _ := e.pendingPredecessors.LoadOrStore(edge.TargetNodeId, 0)
		e.pendingPredecessors.Store(edge.TargetNodeId, count.(int)+1)
	}
}

// NodeResult 表示一个节点的执行结果
type NodeResult struct {
	NodeID string
	Error  error
}

// Execute 启动流程执行（并发控制和调度器集中）
func (e *LogicFlowEngine) Execute(emitter chan string) {
	// 执行队列（线程安全，利用 channel 或 slice + sync.Mutex 也可）
	executionQueue := make([]string, 0)
	executionQueue = append(executionQueue, e.startNodeID)
	for len(executionQueue) > 0 {
		fmt.Println("执行队列", executionQueue)

		// 当前批次节点
		currentBatch := make([]string, len(executionQueue))
		copy(currentBatch, executionQueue)
		executionQueue = []string{} // 清空队列，准备下一轮

		var wg sync.WaitGroup
		resultChan := make(chan *NodeResult, len(currentBatch))

		// 并发执行当前批次所有节点
		for _, nodeID := range currentBatch {
			nd := e.nodes[nodeID]
			if nd == nil {
				log.Printf("节点 %s 不存在，跳过", nodeID)
				continue
			}
			wg.Add(1)
			go func(n node.Node, nid string) {
				defer wg.Done()
				err := n.Execute(e.input, e.output, emitter)
				resultChan <- &NodeResult{
					NodeID: nid,
					Error:  err,
				}
			}(nd, nodeID)
		}

		// 等待所有节点完成
		wg.Wait()
		close(resultChan)

		// 检查结果，处理下一批或中断流程
		hasError := false
		for result := range resultChan {
			if result.Error != nil {
				// 出错则中断流程
				log.Printf("流程执行失败，节点 %s 异常：%v", result.NodeID, result.Error)
				emitter <- base64.StdEncoding.EncodeToString([]byte(result.Error.Error()))
				hasError = true
				break
			}
			// 成功执行，推进下一批节点
			e.processNextNodes(result.NodeID, &executionQueue)
		}

		if hasError {
			return
		}
	}
	time.Sleep(1 * time.Second) // 确保所有输出都能被消费
}

// processNextNodes 决定当前节点完成后，下一个执行哪些节点
func (e *LogicFlowEngine) processNextNodes(currentNodeID string, enqueue *[]string) {
	for _, edge := range e.adjacencyList[currentNodeID] {
		if !edge.EvaluateCondition(e.input) {
			continue
		}

		switch edge.LoopType {
		case "1": // LoopStart
			var values []interface{}
			if raw, ok := e.input[edge.LoopVariable].(string); ok {
				_ = json.Unmarshal([]byte(raw), &values)
			}
			if len(values) == 0 {
				continue
			}
			e.loopValuesMap[edge.LoopVariable] = values
			e.loopIndexes[edge.LoopVariable+"Index"] = 0
			e.input[edge.LoopVariable+"Value"] = values[0]
			*enqueue = append(*enqueue, edge.TargetNodeId)

		case "2": // LoopEnd
			values := e.loopValuesMap[edge.LoopVariable]
			idx := e.loopIndexes[edge.LoopVariable+"Index"]
			if idx+1 < len(values) {
				e.loopIndexes[edge.LoopVariable+"Index"] = idx + 1
				e.input[edge.LoopVariable+"Value"] = values[idx+1]
				loopStartID := e.findLoopStartNodeID(edge.LoopVariable)
				if loopStartID != "" {
					*enqueue = append(*enqueue, loopStartID)
				}
			} else if e.decrementAndCheck(edge.TargetNodeId) {
				*enqueue = append(*enqueue, edge.TargetNodeId)
			}
		default: // 普通边
			if e.decrementAndCheck(edge.TargetNodeId) {
				*enqueue = append(*enqueue, edge.TargetNodeId)
			}
		}
	}
}

// decrementAndCheck 前置计数递减，并判断是否可以执行
func (e *LogicFlowEngine) decrementAndCheck(nodeID string) bool {
	count, _ := e.pendingPredecessors.Load(nodeID)
	if count == nil {
		return true
	}

	// 原子操作递减
	newCount := count.(int) - 1
	e.pendingPredecessors.Store(nodeID, newCount)
	return newCount == 0
}

// findLoopStartNodeID 找到给定 loopVariable 的 LoopStart 节点
func (e *LogicFlowEngine) findLoopStartNodeID(loopVar string) string {
	for _, edges := range e.adjacencyList {
		for _, edge := range edges {
			if edge.LoopType == "1" && edge.LoopVariable == loopVar {
				return edge.TargetNodeId
			}
		}
	}
	return ""
}

// LoadFromJson 从 JSON 中加载节点和边
func (e *LogicFlowEngine) LoadFromJson(jsonString string) error {
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &jsonData); err != nil {
		return err
	}

	// 加载节点
	nodes, ok := jsonData["nodes"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid node structure")
	}
	for _, nodeItem := range nodes {
		nodeData := nodeItem.(map[string]interface{})
		id := nodeData["id"].(string)
		tType := nodeData["type"].(string)
		props := nodeData["properties"]
		bs, _ := json.Marshal(props)
		var nd node.Node
		switch tType {
		case "StartNode":
			e.startNodeID = id
			var p properties.ScriptNodeProperties
			_ = json.Unmarshal(bs, &p)
			nd = node.NewStartNode(id, p)
		case "DatabaseNode":
			var p properties.DatabaseNodeProperties
			_ = json.Unmarshal(bs, &p)
			nd = node.NewDatabaseNode(id, p)
		case "FunctionNode":
			var p properties.FunctionNodeProperties
			_ = json.Unmarshal(bs, &p)
			nd = node.NewFunctionNode(id, p)
		case "HttpNode":
			var p properties.HttpNodeProperties
			_ = json.Unmarshal(bs, &p)
			nd = node.NewHttpNode(id, p)
		case "LLMNode":
			var p properties.LlmNodeProperties
			_ = json.Unmarshal(bs, &p)
			nd = node.NewLlmNode(id, p)
		case "McpNode":
			var p properties.McpNodeProperties
			_ = json.Unmarshal(bs, &p)
			nd = node.NewMcpNode(id, p)
		case "ScriptNode":
			var p properties.ScriptNodeProperties
			_ = json.Unmarshal(bs, &p)
			nd = node.NewScriptNode(id, p)
		case "VectorNode":
			var p properties.VectorNodeProperties
			_ = json.Unmarshal(bs, &p)
			nd = node.NewVectorNode(id, p)
		case "EndNode":
			var p properties.EndNodeProperties
			_ = json.Unmarshal(bs, &p)
			nd = node.NewEndNode(id, p)
		}
		e.RegisterNode(nd)
	}

	// 加载边
	edges, ok := jsonData["edges"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid edges structure")
	}
	for _, edgeItem := range edges {
		edgeData := edgeItem.(map[string]interface{})
		sourceNodeID, _ := edgeData["sourceNodeId"].(string)
		targetNodeID, _ := edgeData["targetNodeId"].(string)
		condition, _ := edgeData["condition"].(string)
		loopType, _ := edgeData["loopType"].(string)
		loopVariable, _ := edgeData["loopVariable"].(string)
		fmt.Println(sourceNodeID, targetNodeID, condition, loopType, loopVariable)
		edge := node.NewEdge(sourceNodeID, targetNodeID, condition, loopType, loopVariable)
		e.RegisterEdge(edge)
	}

	return nil
}

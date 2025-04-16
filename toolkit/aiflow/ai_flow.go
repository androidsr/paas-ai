package aiflow

import (
	"encoding/json"
	"fmt"
	"paas-ai/toolkit/aiflow/node"
	"paas-ai/toolkit/aiflow/properties"
	"sync"
)

type Callback func(line string)

type LogicFlowEngine struct {
	nodes              map[string]node.Node
	adjacencyList      map[string][]Edge
	input              map[string]any
	output             map[string]any
	pendingPredecessor map[string]int
	startNodeId        string
	mu                 sync.Mutex
	queue              *[]string
}

func NewEngine() *LogicFlowEngine {
	return &LogicFlowEngine{
		nodes:              make(map[string]node.Node),
		adjacencyList:      make(map[string][]Edge),
		input:              make(map[string]any),
		output:             make(map[string]any),
		pendingPredecessor: make(map[string]int),
	}
}

func (e *LogicFlowEngine) SetInput(name string, value any) {
	e.input[name] = value
}

func (e *LogicFlowEngine) GetOutput(name string) any {
	return e.output[name]
}

// registerNode 注册节点
func (engine *LogicFlowEngine) RegisterNode(nd node.Node) {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	engine.nodes[nd.ID()] = nd
	if _, ok := nd.(*node.StartNode); ok {
		engine.startNodeId = nd.ID()
	}
}

// registerEdge 注册边
func (engine *LogicFlowEngine) RegisterEdge(edge Edge) {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	engine.adjacencyList[edge.SourceNodeId] = append(engine.adjacencyList[edge.SourceNodeId], edge)
	engine.pendingPredecessor[edge.TargetNodeId]++
}

func (e *LogicFlowEngine) Execute(emitter chan string) {
	fmt.Println("开始执行流程")

	// 每次执行时重置 pendingPredecessor
	e.pendingPredecessor = map[string]int{}
	for _, edges := range e.adjacencyList {
		for _, edge := range edges {
			e.pendingPredecessor[edge.TargetNodeId]++
		}
	}

	var wg sync.WaitGroup
	queue := make(chan string, 100)
	active := make(chan struct{}, 100) // 用于追踪活跃的任务数（非阻塞控制）

	queue <- e.startNodeId

	for {
		select {
		case nodeId := <-queue:
			active <- struct{}{} // 一个节点开始执行
			wg.Add(1)
			go func(nodeId string) {
				defer func() {
					<-active // 一个节点执行完成
					wg.Done()
				}()
				nd := e.nodes[nodeId]
				ok := nd.Execute(e.input, e.output, emitter)
				if !ok {
					emitter <- fmt.Sprintf("节点 %s 执行失败，中断流程", nd.ID())
					return
				}
				e.processNext(nodeId, queue)
			}(nodeId)

		default:
			if len(queue) == 0 && len(active) == 0 {
				wg.Wait()
				close(emitter)
				fmt.Println("流程执行完成")
				return
			}
		}
	}
}

func (e *LogicFlowEngine) processNext(currentNodeId string, queue chan string) {
	for _, edge := range e.adjacencyList[currentNodeId] {
		if edge.EvaluateCondition(e.input) {
			targetId := edge.TargetNodeId
			e.mu.Lock()
			e.pendingPredecessor[targetId]--
			if e.pendingPredecessor[targetId] <= 0 {
				queue <- targetId // 不限制是否执行过
			}
			e.mu.Unlock()
		}
	}
}

// loadFromJson 从 JSON 字符串加载配置
func (engine *LogicFlowEngine) LoadFromJson(jsonString string) error {
	var jsonNode map[string]any
	err := json.Unmarshal([]byte(jsonString), &jsonNode)
	if err != nil {
		return err
	}
	nodes, ok := jsonNode["nodes"].([]any)
	if !ok {
		return fmt.Errorf("Invalid nodes structure")
	}
	for _, nodeItem := range nodes {
		nodeData, ok := nodeItem.(map[string]any)
		if !ok {
			return fmt.Errorf("Invalid node data")
		}
		id := nodeData["id"].(string)
		tType := nodeData["type"].(string)
		propertiesStr := nodeData["properties"].(interface{})
		if propertiesStr == nil {
			continue
		}
		bs, err := json.Marshal(propertiesStr)
		if err != nil {
			fmt.Println(err)
			return err
		}
		var nd node.Node
		switch tType {
		case "StartNode":
			var nodeProperties properties.ScriptNodeProperties
			json.Unmarshal(bs, &nodeProperties)
			nd = node.NewStartNode(id, nodeProperties)
		case "DatabaseNode":
			var nodeProperties properties.DatabaseNodeProperties
			json.Unmarshal(bs, &nodeProperties)
			nd = node.NewDatabaseNode(id, nodeProperties)
		case "FunctionNode":
			var nodeProperties properties.FunctionNodeProperties
			json.Unmarshal(bs, &nodeProperties)
			nd = node.NewFunctionNode(id, nodeProperties)
		case "HttpNode":
			var nodeProperties properties.HttpNodeProperties
			json.Unmarshal(bs, &nodeProperties)
			nd = node.NewHttpNode(id, nodeProperties)
		case "LLMNode":
			var nodeProperties properties.LlmNodeProperties
			json.Unmarshal(bs, &nodeProperties)
			nd = node.NewLlmNode(id, nodeProperties)
		case "McpNode":
			var nodeProperties properties.McpNodeProperties
			json.Unmarshal(bs, &nodeProperties)
			nd = node.NewMcpNode(id, nodeProperties)
		case "ScriptNode":
			var nodeProperties properties.ScriptNodeProperties
			json.Unmarshal(bs, &nodeProperties)
			nd = node.NewScriptNode(id, nodeProperties)
		case "VectorNode":
			var nodeProperties properties.VectorNodeProperties
			json.Unmarshal(bs, &nodeProperties)
			nd = node.NewVectorNode(id, nodeProperties)
		case "EndNode":
			var nodeProperties properties.EndNodeProperties
			json.Unmarshal(bs, &nodeProperties)
			nd = node.NewEndNode(id, nodeProperties)
		}
		engine.RegisterNode(nd)
	}
	// 加载边
	edges, ok := jsonNode["edges"].([]any)
	if !ok {
		return fmt.Errorf("Invalid edges structure")
	}
	for _, edgeItem := range edges {
		edgeData, ok := edgeItem.(map[string]any)
		if !ok {
			return fmt.Errorf("Invalid edge data")
		}
		condition, ok := edgeData["condition"].(string)
		if !ok {
			condition = ""
		}
		edge := Edge{
			SourceNodeId: edgeData["sourceNodeId"].(string),
			TargetNodeId: edgeData["targetNodeId"].(string),
			Condition:    condition,
		}
		engine.RegisterEdge(edge)
	}
	return nil
}

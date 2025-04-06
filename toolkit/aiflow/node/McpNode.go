package node

import (
	"fmt"
	"paas-ai/toolkit/aiflow/properties"
)

type McpNode struct {
	NodeId     string
	properties properties.McpNodeProperties
}

func NewMcpNode(nodeId string, properties properties.McpNodeProperties) *McpNode {
	return &McpNode{NodeId: nodeId, properties: properties}
}

func (s *McpNode) ID() string {
	return s.NodeId
}

func (m *McpNode) Execute(input map[string]any, output map[string]any, emitter chan string) bool {
	emitter <- "执行mcp节点" + m.NodeId + "开始执行"
	fmt.Println("执行mcp节点", m.NodeId)
	return true
}

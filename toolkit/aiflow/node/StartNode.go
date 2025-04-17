package node

import (
	"fmt"
	"paas-ai/toolkit/aiflow/properties"
)

type StartNode struct {
	NodeId     string
	properties properties.ScriptNodeProperties
}

func NewStartNode(nodeId string, properties properties.ScriptNodeProperties) *StartNode {
	return &StartNode{NodeId: nodeId, properties: properties}
}

func (s *StartNode) ID() string {
	return s.NodeId
}

func (m *StartNode) Execute(input map[string]any, output map[string]any, emitter chan string) error {
	fmt.Println("执行开始节点", m.NodeId)
	return nil
}

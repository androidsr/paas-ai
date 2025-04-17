package node

import (
	"fmt"
	"paas-ai/toolkit/aiflow/properties"
)

type ScriptNode struct {
	NodeId     string
	properties properties.ScriptNodeProperties
}

func NewScriptNode(nodeId string, properties properties.ScriptNodeProperties) *ScriptNode {
	return &ScriptNode{NodeId: nodeId, properties: properties}
}

func (s *ScriptNode) ID() string {
	return s.NodeId
}

func (m *ScriptNode) Execute(input map[string]any, output map[string]any, emitter chan string) error {
	//emitter <- "执行脚本节点" + m.NodeId + "开始执行"
	fmt.Println("执行脚本节点", m.NodeId)
	return nil
}

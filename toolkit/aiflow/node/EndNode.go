package node

import (
	"encoding/json"
	"fmt"
	"paas-ai/toolkit/aiflow/properties"
	"paas-ai/toolkit/aiflow/utils"
)

type EndNode struct {
	NodeId     string
	properties properties.EndNodeProperties
}

func NewEndNode(nodeId string, properties properties.EndNodeProperties) *EndNode {
	return &EndNode{NodeId: nodeId, properties: properties}
}

func (s *EndNode) ID() string {
	return s.NodeId
}

func (m *EndNode) Execute(input map[string]any, output map[string]any, emitter chan string) bool {
	messages := m.properties.Messages
	result := make(map[string]any, 0)
	fmt.Println(output)
	for _, message := range messages {
		fmt.Println(message.VarValue)
		data := utils.ReplaceExpression(message.VarValue, input, output)
		if message.VarName != "" {
			result[message.VarName] = data
		} else {
			utils.OutputPrint(emitter, data)
		}
	}
	if len(result) > 0 {
		jsonData, err := json.Marshal(result)
		if err != nil {
			utils.OutputError(emitter, "结束", "输出结果转换为JSON失败")
			return false
		}
		utils.OutputPrint(emitter, string(jsonData))
	}
	return true
}

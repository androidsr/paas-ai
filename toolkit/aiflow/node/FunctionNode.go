package node

import (
	"context"
	"encoding/json"
	"paas-ai/entity"
	"paas-ai/toolkit"
	"paas-ai/toolkit/aiflow/properties"
	"paas-ai/toolkit/aiflow/utils"

	"github.com/androidsr/sc-go/mapper"
	"github.com/tmc/langchaingo/llms"
)

type FunctionNode struct {
	NodeId     string
	properties properties.FunctionNodeProperties
}

func NewFunctionNode(nodeId string, properties properties.FunctionNodeProperties) *FunctionNode {
	return &FunctionNode{NodeId: nodeId, properties: properties}
}

func (s *FunctionNode) ID() string {
	return s.NodeId
}

func (m *FunctionNode) Execute(input map[string]any, output map[string]any, emitter chan string) bool {
	funcDb := mapper.NewHelper[entity.Function]()
	funcInfo := new(entity.Function)
	funcInfo.Id = m.properties.FuncCall
	err := funcDb.SelectOne(funcInfo)
	if err != nil {
		utils.OutputError(emitter, m.properties.Name, "未找到对应的函数内容")
		return false
	}

	db := mapper.NewHelper[entity.AiChannel]()
	llmChannel := new(entity.AiChannel)
	llmChannel.Id = m.properties.ChannelId
	db.SelectOne(llmChannel)
	llm, err := toolkit.NewOpenAI(llmChannel.Url, llmChannel.Token, m.properties.Model)
	if err != nil {
		utils.OutputError(emitter, m.properties.Name, err.Error())
		return false
	}

	for _, v := range m.properties.Messages {
		if v.Message == "" {
			continue
		}
		msg := utils.ReplaceExpression(v.Message, input, output)
		if v.MessageType == "human" {
			llm.SetContent(llms.TextParts(llms.ChatMessageTypeHuman, msg))
		} else if v.MessageType == "system" {
			llm.SetContent(llms.TextParts(llms.ChatMessageTypeSystem, msg))
		} else if v.MessageType == "ai" {
			llm.SetContent(llms.TextParts(llms.ChatMessageTypeAI, msg))
		}
	}

	if m.properties.PrintInput {
		bs, _ := json.Marshal(llm.GetContent())
		utils.InputPrint(emitter, m.properties.Name, string(bs))
	}
	data, err := llm.GenerateFunction(context.Background(), "["+funcInfo.FuncContent+"]", func(fun *llms.FunctionCall) (bool, string, error) {
		return false, "", nil //ExecFuncCall(fun)
	}, nil)
	if err != nil {
		utils.OutputError(emitter, m.properties.Name, err.Error())
		return false
	}
	if m.properties.PrintOutput {
		utils.OutputPrint(emitter, data)
	}
	if m.properties.InputHistory {
		input[m.properties.ParameterName] = data
	}
	if m.properties.ResultHistory {
		output[m.properties.ParameterName] = data
	}
	return true
}

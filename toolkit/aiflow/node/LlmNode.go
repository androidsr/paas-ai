package node

import (
	"context"
	"encoding/json"
	"errors"
	"paas-ai/entity"
	"paas-ai/toolkit"
	"paas-ai/toolkit/aiflow/properties"
	"paas-ai/toolkit/aiflow/utils"

	"github.com/androidsr/sc-go/mapper"
	"github.com/tmc/langchaingo/llms"
)

type LlmNode struct {
	NodeId     string
	properties properties.LlmNodeProperties
}

func NewLlmNode(nodeId string, properties properties.LlmNodeProperties) *LlmNode {
	return &LlmNode{NodeId: nodeId, properties: properties}
}

func (s *LlmNode) ID() string {
	return s.NodeId
}

func (m *LlmNode) Execute(input map[string]any, output map[string]any, emitter chan string) error {
	db := mapper.NewHelper[entity.AiChannel]()
	llmChannel := new(entity.AiChannel)
	llmChannel.Id = m.properties.ChannelId
	db.SelectOne(llmChannel)
	llm, err := toolkit.NewOpenAI(llmChannel.Url, llmChannel.Token, m.properties.Model)
	if err != nil {
		utils.OutputError(emitter, m.properties.Name, err.Error())
		return errors.New("调用大模型失败")
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
	bs, _ := json.Marshal(llm.GetContent())
	inputValue := string(bs)
	if m.properties.PrintInput {
		utils.InputPrint(emitter, m.properties.Name, inputValue)
	}
	resp, err := llm.GenerateCallback(context.Background(), nil)
	if err != nil {
		utils.OutputError(emitter, m.properties.Name, err.Error())
		return errors.New("大模型响应失败")
	}
	data := resp.Choices[0].Content
	if m.properties.PrintOutput {
		utils.OutputPrint(emitter, data)
	}
	if m.properties.InputHistory {
		input[m.properties.ParameterName] = llm.GetContent()
	}
	if m.properties.ResultHistory {
		output[m.properties.ParameterName] = data
	}
	return nil
}

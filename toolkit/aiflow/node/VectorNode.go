package node

import (
	"bytes"
	"context"
	"encoding/json"
	"paas-ai/entity"
	"paas-ai/toolkit"
	"paas-ai/toolkit/aiflow/properties"
	"paas-ai/toolkit/aiflow/utils"

	"github.com/androidsr/sc-go/mapper"
	"github.com/tmc/langchaingo/llms"
)

type VectorNode struct {
	NodeId     string
	properties properties.VectorNodeProperties
}

func NewVectorNode(nodeId string, properties properties.VectorNodeProperties) *VectorNode {
	return &VectorNode{NodeId: nodeId, properties: properties}
}

func (s *VectorNode) ID() string {
	return s.NodeId
}

func (m *VectorNode) Execute(input map[string]any, output map[string]any, emitter chan string) bool {
	coll := toolkit.NewCollection()
	c := coll.Get(m.properties.CollectionId)
	ctx := context.Background()
	doc := toolkit.NewDocument(c)
	filter := make(map[string]string, 0)
	if m.properties.Filename != "" {
		filter["filename"] = m.properties.Filename
	}
	message := utils.ReplaceExpression(m.properties.Message, input, output)

	docs, err := doc.Query(message, m.properties.Limit, filter)
	if err != nil {
		utils.OutputError(emitter, m.properties.Name, err.Error())
		return false
	}
	content := bytes.Buffer{}
	for _, doc := range docs {
		if doc.Similarity >= m.properties.SimilarityScore {
			continue
		}
		content.WriteString(doc.Content)
		content.WriteString("\n")
	}
	if m.properties.Analyse && m.properties.ResultHistory {
		text := content.String()
		output[m.properties.ParameterName] = text
		if m.properties.PrintOutput {
			utils.OutputPrint(emitter, text)
		}
		return true
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
	if m.properties.System != "" {
		llm.SetContent(llms.TextParts(llms.ChatMessageTypeSystem, m.properties.System))
	}
	inputMsg := content.String() + message
	llm.SetContent(llms.TextParts(llms.ChatMessageTypeHuman, inputMsg))
	if m.properties.PrintInput {
		bs, _ := json.Marshal(llm.GetContent())
		utils.InputPrint(emitter, m.properties.Name, string(bs))
	}
	resp, err := llm.GenerateCallback(ctx, nil)
	if err != nil {
		utils.OutputError(emitter, m.properties.Name, err.Error())
		return false
	}
	data := resp.Choices[0].Content
	if m.properties.PrintOutput {
		utils.OutputPrint(emitter, data)
	}
	if m.properties.InputHistory {
		output[m.properties.ParameterName] = inputMsg
	}
	if m.properties.ResultHistory {
		output[m.properties.ParameterName] = data
	}
	return true
}

package toolkit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

var openAI *OpenAI

type OpenAI struct {
	model   string
	llm     *openai.LLM
	content []llms.MessageContent
}

/** 创建模型连接**/
func NewOpenAI(url string, token string, model string) (*OpenAI, error) {
	openAI = &OpenAI{model: model}
	var err error
	openAI.llm, err = openai.New(openai.WithBaseURL(url), openai.WithModel(model), openai.WithToken(token))
	if err != nil {
		return nil, err
	}
	openAI.content = []llms.MessageContent{}
	return openAI, nil
}

/** 设置调用模型提示词**/
func (m *OpenAI) SetContent(messages ...llms.MessageContent) {
	m.content = append(m.content, messages...)
}
func (m *OpenAI) GetContent() []llms.MessageContent {
	return m.content
}

func (m *OpenAI) GetContentSize() int {
	return len(m.content)
}
func (m *OpenAI) RemoveContent(index int) {
	if index != -1 {
		m.content = append(m.content[:index], m.content[index+1:]...)
	}
}
func (m *OpenAI) GenerateFunction(ctx context.Context, toolsContent string, funcCallback func(fun *llms.FunctionCall) (bool, string, error), callback func(ctx context.Context, chunk []byte) error, options ...llms.CallOption) (string, error) {
	if toolsContent != "" {
		tools := make([]llms.Tool, 0)
		err := json.Unmarshal([]byte(toolsContent), &tools)
		if err != nil {
			log.Printf("函数定义错误：%v\n", err)
			return "", err
		}
		options = append(options, llms.WithTools(tools))
	}
	log.Printf("提示词：%s\n", m.content)
	var funcResult string
	var isAi bool
	for {
		resp, err := m.llm.GenerateContent(ctx, m.content, options...)
		if err != nil {
			return "", err
		}
		choice1 := resp.Choices[0]
		if choice1.FuncCall == nil {
			if callback == nil {
				return choice1.Content, nil
			}
			return "", nil
		} else {
			toolCall := choice1.ToolCalls[0]
			log.Printf("执行函数：%s; 输入参数：%s\n\n", choice1.FuncCall.Name, choice1.FuncCall.Arguments)
			isAi, funcResult, err = funcCallback(choice1.FuncCall)
			if err != nil {
				if !isAi {
					return "", err
				}
				log.Printf("执行函数失败：%s\n", err.Error())
				contentPart := []llms.ContentPart{llms.ToolCallResponse{ToolCallID: toolCall.ID, Name: choice1.FuncCall.Name, Content: err.Error()}}
				m.content = append(m.content, llms.MessageContent{Role: llms.ChatMessageTypeTool, Parts: contentPart})
			} else {
				if callback != nil {
					options = append(options, llms.WithStreamingFunc(callback))
				}
				if funcResult == "" || !isAi {
					break
				}
				contentPart := []llms.ContentPart{llms.ToolCallResponse{ToolCallID: toolCall.ID, Name: choice1.FuncCall.Name, Content: funcResult}}
				m.content = append(m.content, llms.MessageContent{Role: llms.ChatMessageTypeTool, Parts: contentPart})
			}
		}
	}
	return funcResult, nil
}

/** 调用模型生成结果；异步回调 **/
func (m *OpenAI) GenerateCallback(ctx context.Context, callback func(ctx context.Context, chunk []byte) error, options ...llms.CallOption) (*llms.ContentResponse, error) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("接收消息错误：%v", err)
		}
	}()
	if callback != nil {
		options = append(options, llms.WithStreamingFunc(callback))
	}
	resp, err := m.llm.GenerateContent(ctx, m.content, options...)
	m.content = []llms.MessageContent{}
	return resp, err
}

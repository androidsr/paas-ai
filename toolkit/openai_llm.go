package toolkit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

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

func (m *OpenAI) GenerateFunction(ctx context.Context, toolsContent string,
	funcCallback func(fun *llms.FunctionCall) (bool, string, error), callback func(ctx context.Context, chunk []byte) error, options ...llms.CallOption) (string, error) {
	// 加载 tools
	if toolsContent != "" {
		var tools []llms.Tool
		if err := json.Unmarshal([]byte(toolsContent), &tools); err != nil {
			log.Printf("函数定义错误：%v\n", err)
			return "", err
		}
		options = append(options, llms.WithTools(tools))
	}

	log.Printf("提示词：%s\n", m.content)
	resultList := make([]string, 0)

	for {
		resp, err := m.llm.GenerateContent(ctx, m.content, options...)
		if err != nil {
			return "", err
		}

		choice := resp.Choices[0]
		if len(choice.ToolCalls) == 0 {
			// 无函数调用
			if callback == nil {
				resultList = append(resultList, choice.Content)
				return strings.Join(resultList, " "), nil
			}
			return "", nil
		}

		var (
			successParts []llms.ContentPart
			errorParts   []llms.ContentPart
		)

		for _, toolCall := range choice.ToolCalls {
			log.Printf("执行函数：%s; 输入参数：%s\n", toolCall.FunctionCall.Name, toolCall.FunctionCall.Arguments)

			isAi, funcResult, err := funcCallback(&llms.FunctionCall{
				Name:      toolCall.FunctionCall.Name,
				Arguments: toolCall.FunctionCall.Arguments,
			})

			if err != nil {
				if !isAi {
					return "", err
				}
				log.Printf("执行函数失败：%s\n", err.Error())
				errorParts = append(errorParts, llms.ToolCallResponse{
					ToolCallID: toolCall.ID,
					Name:       toolCall.FunctionCall.Name,
					Content:    err.Error(),
				})
			} else {
				// 保存每个函数调用的返回内容
				resultList = append(resultList, funcResult)
				successParts = append(successParts, llms.ToolCallResponse{
					ToolCallID: toolCall.ID,
					Name:       toolCall.FunctionCall.Name,
					Content:    funcResult,
				})
			}
		}

		// 返回所有成功/失败内容
		if len(errorParts) > 0 {
			m.content = append(m.content, llms.MessageContent{
				Role:  llms.ChatMessageTypeTool,
				Parts: errorParts,
			})
			continue
		}

		// 避免重复添加 callback
		if callback != nil {
			cbOption := llms.WithStreamingFunc(callback)
			if !containsOption(options, cbOption) {
				options = append(options, cbOption)
			}
		}

		m.content = append(m.content, llms.MessageContent{
			Role:  llms.ChatMessageTypeTool,
			Parts: successParts,
		})

		break
	}

	return strings.Join(resultList, " "), nil
}

func containsOption(options []llms.CallOption, target llms.CallOption) bool {
	for _, opt := range options {
		if fmt.Sprintf("%v", opt) == fmt.Sprintf("%v", target) {
			return true
		}
	}
	return false
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

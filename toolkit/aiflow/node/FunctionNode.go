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
		return ExecFuncCall(fun)
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

// 函数工具调用逻辑
func ExecFuncCall(function *llms.FunctionCall) (bool, string, error) {
	input := make(map[string]interface{}, 0)
	err := json.Unmarshal([]byte(function.Arguments), &input)
	if err != nil {
		return true, "", errors.New("解析AI参数出错：" + err.Error())
	}

	switch function.Name {
	case "getTables":
		dbname, ok := input["dbname"].(string)
		if dbname == "" || !ok {
			return true, "", errors.New("无效函数，未找到dbname参数")
		}
		db := toolkit.NewDbName(dbname)
		data, err := db.GetTables(dbname)
		if err != nil {
			return true, "", errors.New("查询数据库表信息出错：" + err.Error())
		}
		bs, err := json.Marshal(data)
		if err != nil {
			return true, "", errors.New("表信息转换JSON出错：" + err.Error())
		}
		return false, string(bs), nil
	case "getTableDetail":
		dbname, ok := input["dbname"].(string)
		if dbname == "" || !ok {
			return true, "", errors.New("无效函数，未找到dbname参数")
		}
		tableNames, ok := input["tableNames"].([]interface{})
		if !ok {
			return true, "", errors.New("无效函数，未找到tableNames参数")
		}
		db := toolkit.NewDbName(dbname)
		data, err := db.GetTableDetail(dbname, tableNames)
		if err != nil {
			return true, "", errors.New("查询数据库表信息出错：" + err.Error())
		}
		bs, err := json.Marshal(data)
		if err != nil {
			return true, "", errors.New("表信息转换JSON出错：" + err.Error())
		}
		return false, string(bs), nil
	case "queryExecSql":
		dbname, ok := input["dbname"].(string)
		if dbname == "" || !ok {
			return true, "", errors.New("无效函数，未找到dbname参数")
		}
		sql, ok := input["sql"].(string)
		if sql == "" || !ok {
			return true, "", errors.New("无效函数，未找到sql参数")
		}
		db := toolkit.NewDbName(dbname)
		data, err := db.QuerySql(sql)
		if err != nil {
			return true, "", errors.New("查询数据库表信息出错：" + err.Error())
		}
		bs, err := json.Marshal(data)
		if err != nil {
			return true, "", errors.New("表信息转换JSON出错：" + err.Error())
		}
		return false, string(bs), nil
	}
	return false, "", errors.New("未找到对应的函数内容")
}

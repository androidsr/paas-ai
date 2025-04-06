package node

import (
	"encoding/json"
	"fmt"
	"paas-ai/toolkit/aiflow/properties"
	"paas-ai/toolkit/aiflow/utils"
	"strings"

	"github.com/androidsr/sc-go/shttp"
)

type HttpNode struct {
	NodeId     string
	properties properties.HttpNodeProperties
}

func NewHttpNode(nodeId string, properties properties.HttpNodeProperties) *HttpNode {
	return &HttpNode{NodeId: nodeId, properties: properties}
}

func (s *HttpNode) ID() string {
	return s.NodeId
}

func (m *HttpNode) Execute(input map[string]any, output map[string]any, emitter chan string) bool {
	body := utils.ReplaceExpression(m.properties.Body, input, output)
	url := strings.TrimSpace(m.properties.Url)
	if m.properties.PrintInput && body != "" {
		bs, err := json.Marshal(body)
		var p any
		if err != nil {
			p = body
		} else {
			p = string(bs)
		}
		utils.InputPrint(emitter, m.properties.Name, fmt.Sprint(p))
	}
	params := make(map[string]interface{}, 0)
	if m.properties.Method != "GET" && body != "" {
		err := json.Unmarshal([]byte(body), &params)
		if err != nil {
			utils.OutputError(emitter, m.properties.Name, err.Error())
			return false
		}
	}
	headers := make(map[string]string, 0)
	if m.properties.Headers != "" {
		kvs := strings.Split(m.properties.Headers, ";")
		for _, v := range kvs {
			if strings.Contains(v, "=") {
				kv := strings.Split(v, "=")
				headers[kv[0]] = kv[1]
			}
		}
	}
	if m.properties.Cookies != "" {
		headers["Cookie"] = m.properties.Cookies
	}
	var data []byte
	var err error
	switch m.properties.Method {
	case "GET":
		data, err = shttp.Get(url+body, m.properties.ContentType, headers)
	case "POST":
		data, err = shttp.Post(url, strings.TrimSpace(m.properties.ContentType), headers, params)
	case "PUT":
		data, err = shttp.Put(url, m.properties.ContentType, headers, params)
	case "DELETE":
		data, err = shttp.Delete(url, m.properties.ContentType, headers, params)
	}
	if err != nil {
		utils.OutputError(emitter, m.properties.Name, err.Error())
		return false
	}

	var result interface{}
	resultMap := make(map[string]interface{}, 0)
	err = json.Unmarshal(data, &resultMap)
	if err != nil {
		result = string(data)
	}
	if m.properties.DataKey != "" && m.properties.Expression != "" && len(resultMap) != 0 {
		/* if !m.IsExpression(expression, fmt.Sprint(resultMap[validateKey]), validateValue) {
			return errors.New("响应数据失败：" + string(data))
		}
		if dataKey != "" {
			result = resultMap[dataKey]
		} else {
			result = resultMap
		} */
		result = resultMap
	} else if resultMap != nil {
		result = resultMap
	}

	if m.properties.PrintOutput && result != nil {
		bs, err := json.Marshal(result)
		var p any
		if err != nil {
			p = body
		} else {
			p = string(bs)
		}
		utils.OutputPrint(emitter, fmt.Sprint(p))
	}
	if m.properties.InputHistory {
		input[m.properties.ParameterName] = body
	}
	if m.properties.ResultHistory {
		output[m.properties.ParameterName] = result
	}
	return true
}

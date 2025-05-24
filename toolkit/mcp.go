package toolkit

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/tmc/langchaingo/llms"
)

func GetToolsContent(mType int, command string) (client.MCPClient, string) {
	var cli client.MCPClient
	var err error
	if mType == 1 {
		commands := strings.Split(command, " ")
		cli, err = client.NewStdioMCPClient(commands[0], []string{}, commands[1:]...)
		if err != nil {
			fmt.Println(err)
			return nil, ""
		}
	} else {
		cli, err = client.NewSSEMCPClient(command)
	}
	_, err = cli.Initialize(context.Background(), mcp.InitializeRequest{})
	if err != nil {
		fmt.Println(err)
		return nil, ""
	}
	r, err := cli.ListTools(context.Background(), mcp.ListToolsRequest{})
	if err != nil {
		fmt.Println(err)
		return nil, ""
	}
	bs, _ := json.Marshal(r.Tools)
	text := strings.ReplaceAll(string(bs), "inputSchema", "parameters")
	funcList := make([]*llms.FunctionDefinition, 0)
	err = json.Unmarshal([]byte(text), &funcList)
	if err != nil {
		return nil, ""
	}
	tools := make([]llms.Tool, 0)
	fmt.Println(funcList)
	for _, f := range funcList {
		tools = append(tools, llms.Tool{Type: "function", Function: f})
	}
	jsonData, err := json.Marshal(tools)
	if err != nil {
		return nil, ""
	}
	return cli, string(jsonData)
}

func ExecTools(cli client.MCPClient, name string, arguments string) (bool, string, error) {
	fmt.Println("接收到的参数信息：", name, arguments)
	args := make(map[string]any, 0)
	err := json.Unmarshal([]byte(arguments), &args)
	if err != nil {
		return true, "", errors.New("解析工具参数格式不正确：" + err.Error())
	}
	req := mcp.CallToolRequest{
		Params: struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Name:      name,
			Arguments: args,
			Meta:      nil,
		},
	}
	resp, err := cli.CallTool(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return true, "", errors.New("调用工具失败：" + err.Error())
	}
	buf := bytes.Buffer{}
	for _, content := range resp.Content {
		text, ok := content.(mcp.TextContent)
		if ok {
			buf.WriteString(text.Text)
			buf.WriteString("\n")
		}
	}
	fmt.Println(buf.String())
	return false, buf.String(), nil
}

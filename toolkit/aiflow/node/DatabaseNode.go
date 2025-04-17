package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"paas-ai/entity"
	"paas-ai/toolkit"
	"paas-ai/toolkit/aiflow/properties"
	"paas-ai/toolkit/aiflow/utils"
	"regexp"
	"strings"

	"github.com/androidsr/sc-go/mapper"
)

type DatabaseNode struct {
	NodeId     string
	properties properties.DatabaseNodeProperties
}

func NewDatabaseNode(nodeId string, properties properties.DatabaseNodeProperties) *DatabaseNode {
	return &DatabaseNode{NodeId: nodeId, properties: properties}
}

func (s *DatabaseNode) ID() string {
	return s.NodeId
}

func (m *DatabaseNode) Execute(input map[string]any, output map[string]any, emitter chan string) error {
	db := mapper.NewHelper[entity.DbConfig]()
	data := new(entity.DbConfig)
	data.Id = m.properties.SourceId
	err := db.SelectOne(data)
	if err != nil {
		utils.OutputError(emitter, m.properties.Name, "查询链接信息失败")
		return errors.New("查询链接信息失败")
	}
	var whereString string

	for _, w := range m.properties.Wheres {
		sqlWhere := fmt.Sprint(w)
		re := regexp.MustCompile(`#\{([^}]*)\}`)
		matches := re.FindStringSubmatch(sqlWhere)
		if len(matches) <= 1 {
			continue
		}
		key := matches[1]
		value := utils.ReplaceExpression(key, input, output)
		if value == "" {
			continue
		}
		whereString += strings.ReplaceAll(sqlWhere, "#{"+key+"}", " '"+fmt.Sprint(value)+"' ")
	}
	sql := strings.ReplaceAll(m.properties.Sql, "$where", whereString)

	result, err := toolkit.NewDbSearch(data).QuerySql(sql)
	if err != nil {
		utils.OutputError(emitter, m.properties.Name, "执行SQL语句失败"+err.Error())
		return errors.New("执行SQL语句失败")
	}
	if m.properties.PrintOutput {
		bs, err := json.Marshal(result)
		if err != nil {
			utils.OutputPrint(emitter, fmt.Sprint(result))
		} else {
			utils.OutputPrint(emitter, string(bs))
		}
	}
	if m.properties.InputHistory {
		input[m.properties.ParameterName] = whereString
	}
	if m.properties.ResultHistory {
		output[m.properties.ParameterName] = result
	}
	return nil
}

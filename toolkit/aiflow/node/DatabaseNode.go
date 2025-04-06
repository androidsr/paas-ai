package node

import (
	"encoding/json"
	"fmt"
	"paas-ai/entity"
	"paas-ai/toolkit/aiflow/properties"
	"paas-ai/toolkit/aiflow/utils"
	"regexp"
	"strings"

	"github.com/androidsr/sc-go/mapper"
	"gorm.io/gorm"
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

func (m *DatabaseNode) Execute(input map[string]any, output map[string]any, emitter chan string) bool {
	db := mapper.NewHelper[entity.DbConfig]()
	data := new(entity.DbConfig)
	data.Id = m.properties.SourceId
	err := db.SelectOne(data)
	if err != nil {
		utils.OutputError(emitter, m.properties.Name, "查询链接信息失败")
		return false
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

	result, err := m.querySql(db.DB, sql)
	if err != nil {
		utils.OutputError(emitter, m.properties.Name, "执行SQL语句失败"+err.Error())
		return false
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
	return true
}

func (m *DatabaseNode) querySql(db *gorm.DB, sql string) ([]map[string]string, error) {
	rows, err := db.Exec(sql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, _ := rows.Columns()
	values := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range values {
		scans[i] = &values[i]
	}
	results := make([]map[string]string, 0)
	for rows.Next() {
		err := rows.Scan(scans...)
		if err != nil {
			return nil, err
		}
		row := make(map[string]string, 0)
		for k, v := range values {
			key := cols[k]
			row[key] = string(v)
		}
		results = append(results, row)
	}
	return results, nil
}

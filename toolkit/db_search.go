package toolkit

import (
	"encoding/json"
	"fmt"
	"log"
	"paas-ai/dv"
	"paas-ai/entity"

	"github.com/androidsr/sc-go/mapper"
	"github.com/androidsr/sc-go/syaml"
	"gorm.io/gorm"
)

type DbSearch struct {
	db *gorm.DB
}

func NewDbSearch(config *entity.DbConfig) *DbSearch {
	db := mapper.Initdb(&syaml.GormInfo{
		Driver:  config.DbType,
		Url:     config.Url,
		ShowSql: true,
	})
	return &DbSearch{db: db}
}

func NewDbName(dbname string) *DbSearch {
	config := &entity.DbConfig{Dbname: dbname}
	err := mapper.NewHelper[entity.DbConfig]().SelectOne(config)
	if err != nil {
		log.Printf("获取数据库配置信息失败：%v", err)
		return nil
	}
	db := mapper.Initdb(&syaml.GormInfo{
		Driver:  config.DbType,
		Url:     config.Url,
		ShowSql: true,
	})
	return &DbSearch{db: db}
}

func (m DbSearch) GetTables(dbname string) ([]dv.TableInfo, error) {
	sql := fmt.Sprintf(`SELECT table_name as name,table_comment as table_desc  FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = '%s'`, dbname)
	rows, err := m.db.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]dv.TableInfo, 0)
	item := dv.TableInfo{}
	for rows.Next() {
		if err := rows.Scan(&item.TableName, &item.TableDesc); err != nil {
			log.Printf("获取表结构返回数据错误：%v", err)
			continue
		}
		result = append(result, item)
	}
	return result, nil
}

/** 获取表结构定义 **/
func (m DbSearch) GetTableDetail(dbname string, tableNames []interface{}) (string, error) {
	result := make([]string, 0)
	sql := `SHOW CREATE TABLE %s`
	for _, tableName := range tableNames {
		rows, err := m.db.Raw(fmt.Sprintf(sql, tableName)).Rows()
		if err != nil {
			log.Printf("获取表结构返回数据错误：%v", err)
			return "", err
		}
		var detail string
		for rows.Next() {
			err := rows.Scan(&tableName, &detail)
			if err != nil {
				log.Printf("获取表结构返回数据错误：%v", err)
				break
			}
			result = append(result, detail)
		}
		rows.Close()
	}
	bs, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func (m *DbSearch) QuerySql(sql string) ([]map[string]string, error) {
	rows, err := m.db.Raw(sql).Rows()
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

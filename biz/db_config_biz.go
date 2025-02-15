package biz

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"paas-ai/dv"
	"paas-ai/entity"

	"github.com/androidsr/sc-go/mapper"
	"github.com/androidsr/sc-go/model"
	"github.com/androidsr/sc-go/sbuilder"
	"github.com/androidsr/sc-go/sno"
)

type DbConfigBiz struct {
	db *mapper.Mapper[entity.DbConfig]
}

func NewDbConfigBiz() *DbConfigBiz {
	db := mapper.NewHelper[entity.DbConfig]()
	return &DbConfigBiz{db: db}
}

func (m *DbConfigBiz) Get(id string) model.HttpResult {
	data := new(entity.DbConfig)
	data.Id = id
	err := m.db.SelectOne(data)
	if err != nil {
		log.Printf("查询失败: %v\n", err)
		return model.NewFailDefault("查询失败")
	}
	return model.NewOK(data)
}

func (m *DbConfigBiz) Add(task *entity.DbConfig) model.HttpResult {
	task.Id = sno.GetString()
	err := m.db.Insert(task)
	if err != nil {
		log.Printf("添加失败: %v\n", err)
		return model.NewFailDefault("添加失败")
	}
	return model.NewOK(task)
}

func (m *DbConfigBiz) Edit(task *entity.DbConfig) model.HttpResult {
	err := m.db.Update(task, "id", task.Id).Error
	if err != nil {
		log.Printf("修改失败: %v\n", err)
		return model.NewFailDefault("修改失败")
	}
	return model.NewOK(task)
}

func (t *DbConfigBiz) Delete(id string) model.HttpResult {
	err := t.db.Delete("id", id).Error
	if err != nil {
		log.Printf("删除失败: %v\n", err)
		return model.NewFailDefault("删除失败")
	}
	return model.NewOK(true)
}

type DbDataSourceQuery struct {
	Page   *model.PageInfo `json:"page" column:"-"`
	DbType string          `json:"name" column:"a.db_type" keyword:"like"`
	Dbname string          `json:"dbname" column:"a.dbname" keyword:"eq"`
}

func (m *DbConfigBiz) Page(query *DbDataSourceQuery) model.HttpResult {
	sql := `select * from db_config a where 1=1 `
	data := make([]entity.DbConfig, 0)
	b := sbuilder.StructToBuilder(query, sql)
	sql, values := b.Build()
	sql += " order by dbname,id desc "
	result := m.db.SelectPage(&data, query.Page, sql, values...)
	return model.NewOK(result)
}

func (m *DbConfigBiz) GetList() model.HttpResult {
	sql := `select id as value,dbname as label from db_config a where 1=1 `
	data := make([]model.SelectVO, 0)
	b := sbuilder.Builder(sql)
	sql, values := b.Build()
	sql += " order by dbname,id desc "
	err := m.db.SelectSQL(&data, sql, values...)
	if err != nil {
		log.Printf("查询列表失败: %v\n", err)
		return model.NewFailDefault("查询列表失败")
	}
	return model.NewOK(data)
}

/** 获取所有表信息 **/
func (m DbConfigBiz) GetTables(dbname string) ([]dv.TableInfo, error) {
	config := m.FindByConfig(dbname)
	db, err := m.GetConnect(config.DbType, config.Url)
	if err != nil {
		return nil, err
	}
	sql := fmt.Sprintf(`SELECT table_name as name,table_comment as table_desc  FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = '%s'`, dbname)
	rows, err := db.Query(sql)
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
func (m DbConfigBiz) GetTableDetail(dbname string, tableNames []interface{}) (string, error) {
	config := m.FindByConfig(dbname)
	db, err := m.GetConnect(config.DbType, config.Url)
	if err != nil {
		return "", err
	}
	result := make([]string, 0)
	sql := `SHOW CREATE TABLE %s`
	for _, tableName := range tableNames {
		rows, err := db.Query(fmt.Sprintf(sql, tableName))
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

/** 获取表结构定义 **/
func (m DbConfigBiz) QueryExecSql(dbname, sql string) (string, error) {
	log.Printf("执行数据库：%s -> %s", dbname, sql)
	config := m.FindByConfig(dbname)
	db, err := m.GetConnect(config.DbType, config.Url)
	if err != nil {
		return "", err
	}
	results, err := m.QuerySql(db, sql)
	if err != nil {
		return "", err
	}
	bs, err := json.Marshal(results)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func (m DbConfigBiz) QuerySql(db *sql.DB, sql string) ([]map[string]string, error) {
	rows, err := db.Query(sql)
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

func (m DbConfigBiz) FindByConfig(dbname string) *entity.DbConfig {
	data := &entity.DbConfig{}
	data.Dbname = dbname
	err := mapper.NewHelper[entity.DbConfig]().SelectOne(data)
	if err != nil {
		log.Printf("获取数据源失败：%v", err)
		return nil
	}
	return data
}

func (m DbConfigBiz) GetConnect(dbtype string, url string) (*sql.DB, error) {
	db, err := sql.Open(dbtype, url)
	if err != nil {
		return nil, err
	}
	return db, nil
}

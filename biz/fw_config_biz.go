package biz

import (
	"log"
	"paas-ai/entity"

	"github.com/androidsr/sc-go/mapper"
	"github.com/androidsr/sc-go/model"
	"github.com/androidsr/sc-go/sbuilder"
	"github.com/androidsr/sc-go/sno"
)

type FwConfigBiz struct {
	db *mapper.Mapper[entity.FwConfig]
}

func NewFwConfigBiz() *FwConfigBiz {
	db := mapper.NewHelper[entity.FwConfig]()
	return &FwConfigBiz{db: db}
}

func (m *FwConfigBiz) Get(id string) model.HttpResult {
	data := new(entity.FwConfig)
	data.Id = id
	err := m.db.SelectOne(data)
	if err != nil {
		log.Printf("查询失败: %v\n", err)
		return model.NewFailDefault("查询失败")
	}
	return model.NewOK(data)
}

func (m *FwConfigBiz) Add(task *entity.FwConfig) model.HttpResult {
	task.Id = sno.GetString()
	err := m.db.Insert(task)
	if err != nil {
		log.Printf("添加失败: %v\n", err)
		return model.NewFailDefault("添加失败")
	}
	return model.NewOK(task)
}

func (m *FwConfigBiz) Edit(task *entity.FwConfig) model.HttpResult {
	err := m.db.Update(task, "id", task.Id).Error
	if err != nil {
		log.Printf("修改失败: %v\n", err)
		return model.NewFailDefault("修改失败")
	}
	return model.NewOK(task)
}

func (t *FwConfigBiz) Delete(id string) model.HttpResult {
	err := t.db.Delete("id", id).Error
	if err != nil {
		log.Printf("删除失败: %v\n", err)
		return model.NewFailDefault("删除失败")
	}
	return model.NewOK(true)
}

type FwConfigQuery struct {
	Page     *model.PageInfo `json:"page" column:"-"`
	Name     string          `json:"name" column:"a.name" keyword:"like"`
	FlowType string          `json:"flow_type" column:"a.flow_type" keyword:"eq"`
	TenantId string          `json:"tenant_id" column:"a.tenant_id" keyword:"eq"`
}

func (m *FwConfigBiz) Page(query *FwConfigQuery) model.HttpResult {
	sql := `select * from fw_config a where 1=1 `
	data := make([]entity.FwConfig, 0)
	b := sbuilder.StructToBuilder(query, sql)
	sql, values := b.Build()
	sql += " order by name,id desc "
	result := m.db.SelectPage(&data, query.Page, sql, values...)
	return model.NewOK(result)
}

func (m *FwConfigBiz) GetList() model.HttpResult {
	sql := `select id as value, name as label from fw_config a where 1=1 `
	data := make([]model.SelectVO, 0)
	b := sbuilder.Builder(sql)
	sql, values := b.Build()
	sql += " order by name,id desc "
	err := m.db.SelectSQL(&data, sql, values...)
	if err != nil {
		log.Printf("查询列表失败: %v\n", err)
		return model.NewFailDefault("查询列表失败")
	}
	return model.NewOK(data)
}

package biz

import (
	"log"
	"paas-ai/entity"

	"github.com/androidsr/sc-go/mapper"
	"github.com/androidsr/sc-go/model"
	"github.com/androidsr/sc-go/sbuilder"
	"github.com/androidsr/sc-go/sno"
)

type PromptBiz struct {
	db *mapper.Mapper[entity.Prompt]
}

func NewPromptBiz() *PromptBiz {
	db := mapper.NewHelper[entity.Prompt]()
	return &PromptBiz{db: db}
}

func (m *PromptBiz) Get(id string) model.HttpResult {
	data := new(entity.Prompt)
	data.Id = id
	err := m.db.SelectOne(data)
	if err != nil {
		log.Printf("查询失败: %v\n", err)
		return model.NewFailDefault("查询失败")
	}
	return model.NewOK(data)
}

func (m *PromptBiz) Add(task *entity.Prompt) model.HttpResult {
	task.Id = sno.GetString()
	err := m.db.Insert(task)
	if err != nil {
		log.Printf("添加失败: %v\n", err)
		return model.NewFailDefault("添加失败")
	}
	return model.NewOK(task)
}

func (m *PromptBiz) Edit(task *entity.Prompt) model.HttpResult {
	err := m.db.Update(task, "id", task.Id).Error
	if err != nil {
		log.Printf("修改失败: %v\n", err)
		return model.NewFailDefault("修改失败")
	}
	return model.NewOK(task)
}

func (t *PromptBiz) Delete(id string) model.HttpResult {
	err := t.db.Delete("id", id).Error
	if err != nil {
		log.Printf("删除失败: %v\n", err)
		return model.NewFailDefault("删除失败")
	}
	return model.NewOK(true)
}

type PromptQuery struct {
	Page     *model.PageInfo `json:"page" column:"-"`
	RoleName string          `json:"roleName" column:"a.role_name" keyword:"like"`
}

func (m *PromptBiz) Page(query *PromptQuery) model.HttpResult {
	sql := `select * from prompt a where 1=1 `
	data := make([]entity.Prompt, 0)
	b := sbuilder.StructToBuilder(query, sql)
	sql, values := b.Build()
	sql += " order by role_name,id desc "
	result := m.db.SelectPage(&data, query.Page, sql, values...)
	return model.NewOK(result)
}

func (m *PromptBiz) GetList() model.HttpResult {
	sql := `select id as value, role_name as label from prompt a where 1=1 `
	data := make([]model.SelectVO, 0)
	b := sbuilder.Builder(sql)
	sql, values := b.Build()
	sql += " order by role_name,id desc "
	err := m.db.SelectSQL(&data, sql, values...)
	if err != nil {
		log.Printf("查询列表失败: %v\n", err)
		return model.NewFailDefault("查询列表失败")
	}
	return model.NewOK(data)
}

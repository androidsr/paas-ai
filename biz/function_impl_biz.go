package biz

import (
	"log"
	"paas-ai/entity"

	"github.com/androidsr/sc-go/mapper"
	"github.com/androidsr/sc-go/model"
	"github.com/androidsr/sc-go/sbuilder"
	"github.com/androidsr/sc-go/sno"
)

type FunctionImplBiz struct {
	db *mapper.Mapper[entity.FunctionImpl]
}

func NewFunctionImplBiz() *FunctionImplBiz {
	db := mapper.NewHelper[entity.FunctionImpl]()
	return &FunctionImplBiz{db: db}
}

func (m *FunctionImplBiz) Get(id string) model.HttpResult {
	data := new(entity.FunctionImpl)
	data.Id = id
	err := m.db.SelectOne(data)
	if err != nil {
		log.Printf("查询失败: %v\n", err)
		return model.NewFailDefault("查询失败")
	}
	return model.NewOK(data)
}

func (m *FunctionImplBiz) FindByName(name string) *entity.FunctionImpl {
	data := new(entity.FunctionImpl)
	data.Name = name
	err := m.db.SelectOne(data)
	if err != nil {
		return nil
	}
	return data
}

func (m *FunctionImplBiz) Add(task *entity.FunctionImpl) model.HttpResult {
	task.Id = sno.GetString()
	err := m.db.Insert(task)
	if err != nil {
		log.Printf("添加失败: %v\n", err)
		return model.NewFailDefault("添加失败")
	}
	return model.NewOK(task)
}

func (m *FunctionImplBiz) Edit(task *entity.FunctionImpl) model.HttpResult {
	err := m.db.Update(task, "id", task.Id).Error
	if err != nil {
		log.Printf("修改失败: %v\n", err)
		return model.NewFailDefault("修改失败")
	}
	return model.NewOK(task)
}

func (t *FunctionImplBiz) Delete(id string) model.HttpResult {
	err := t.db.Delete("id", id).Error
	if err != nil {
		log.Printf("删除失败: %v\n", err)
		return model.NewFailDefault("删除失败")
	}
	return model.NewOK(true)
}

type FunctionImplQuery struct {
	Page  *model.PageInfo `json:"page" column:"-"`
	Name  string          `json:"name" column:"a.name" keyword:"like"`
	Title string          `json:"title" column:"a.title" keyword:"like"`
}

func (m *FunctionImplBiz) Page(query *FunctionImplQuery) model.HttpResult {
	sql := `select * from function_impl a where 1=1 `
	data := make([]entity.FunctionImpl, 0)
	b := sbuilder.StructToBuilder(query, sql)
	sql, values := b.Build()
	sql += " order by title,id desc "
	result := m.db.SelectPage(&data, query.Page, sql, values...)
	return model.NewOK(result)
}

func (m *FunctionImplBiz) GetList() model.HttpResult {
	sql := `select id as value, title as label from function_impl a where 1=1 `
	data := make([]model.SelectVO, 0)
	b := sbuilder.Builder(sql)
	sql, values := b.Build()
	sql += " order by title,id desc "
	err := m.db.SelectSQL(&data, sql, values...)
	if err != nil {
		log.Printf("查询列表失败: %v\n", err)
		return model.NewFailDefault("查询列表失败")
	}
	return model.NewOK(data)
}

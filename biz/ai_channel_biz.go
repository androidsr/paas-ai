package biz

import (
	"fmt"
	"log"
	"paas-ai/entity"

	"github.com/androidsr/sc-go/mapper"
	"github.com/androidsr/sc-go/model"
	"github.com/androidsr/sc-go/sbuilder"
	"github.com/androidsr/sc-go/sno"
)

type AiChannelBiz struct {
	db *mapper.Mapper[entity.AiChannel]
}

func NewAiChannelBiz() *AiChannelBiz {
	db := mapper.NewHelper[entity.AiChannel]()
	return &AiChannelBiz{db: db}
}

func (m *AiChannelBiz) Get(id string) model.HttpResult {
	data := new(entity.AiChannel)
	data.Id = id
	err := m.db.SelectOne(data)
	if err != nil {
		log.Printf("查询失败: %v\n", err)
		return model.NewFailDefault("查询失败")
	}
	return model.NewOK(data)
}

func (m *AiChannelBiz) GetById(id string) *entity.AiChannel {
	data := new(entity.AiChannel)
	data.Id = id
	err := m.db.SelectOne(data)
	if err != nil {
		return nil
	}
	return data
}

func (m *AiChannelBiz) FindByUrl(url string) model.HttpResult {
	data := new(entity.AiChannel)
	data.Url = url
	err := m.db.SelectOne(data)
	if err != nil {
		log.Printf("查询失败: %v\n", err)
		return model.NewFailDefault("查询失败")
	}
	return model.NewOK(data)
}

func (m *AiChannelBiz) Add(task *entity.AiChannel) model.HttpResult {
	task.Id = sno.GetString()
	err := m.db.Insert(task)
	if err != nil {
		log.Printf("添加失败: %v\n", err)
		return model.NewFailDefault("添加失败")
	}
	return model.NewOK(task)
}

func (m *AiChannelBiz) Edit(task *entity.AiChannel) model.HttpResult {
	err := m.db.Update(task, "id", task.Id).Error
	if err != nil {
		log.Printf("修改失败: %v\n", err)
		return model.NewFailDefault("修改失败")
	}
	return model.NewOK(task)
}

func (t *AiChannelBiz) Delete(id string) model.HttpResult {
	err := t.db.Delete("id", id).Error
	if err != nil {
		log.Printf("删除失败: %v\n", err)
		return model.NewFailDefault("删除失败")
	}
	return model.NewOK(true)
}

type AiChannelQuery struct {
	Page *model.PageInfo `json:"page" column:"-"`
	Name string          `json:"name" column:"a.name" keyword:"like"`
}

func (m *AiChannelBiz) Page(query *AiChannelQuery) model.HttpResult {
	sql := `select * from ai_channel a where 1=1 `
	data := make([]entity.AiChannel, 0)
	b := sbuilder.StructToBuilder(query, sql)
	sql, values := b.Build()
	sql += " order by priority,id desc "
	fmt.Println(query.Page)
	result := m.db.SelectPage(&data, query.Page, sql, values...)
	return model.NewOK(result)
}

func (m *AiChannelBiz) GetList() model.HttpResult {
	sql := `select id as value, name as label,models as supper_id from ai_channel a where 1=1 `
	data := make([]model.SelectVO, 0)
	b := sbuilder.Builder(sql)
	sql, values := b.Build()
	sql += " order by priority,id desc "
	err := m.db.SelectSQL(&data, sql, values...)
	if err != nil {
		log.Printf("查询列表失败: %v\n", err)
		return model.NewFailDefault("查询列表失败")
	}
	return model.NewOK(data)
}

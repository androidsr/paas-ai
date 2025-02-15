package biz

import (
	"encoding/json"
	"fmt"
	"log"
	"paas-ai/entity"
	"paas-ai/toolkit"

	"github.com/androidsr/sc-go/mapper"
	"github.com/androidsr/sc-go/model"
	"github.com/androidsr/sc-go/sbuilder"
	"github.com/androidsr/sc-go/sno"
	"github.com/philippgille/chromem-go"
)

type DocumentBiz struct {
	db *mapper.Mapper[entity.Document]
}

func NewDocumentBiz() *DocumentBiz {
	db := mapper.NewHelper[entity.Document]()
	return &DocumentBiz{db: db}
}

func (m *DocumentBiz) Get(id string) model.HttpResult {
	data := new(entity.Document)
	data.Id = id
	err := m.db.SelectOne(data)
	if err != nil {
		log.Printf("查询失败: %v\n", err)
		return model.NewFailDefault("查询失败")
	}
	return model.NewOK(data)
}

func (m *DocumentBiz) Add(dto *entity.Document) model.HttpResult {
	dto.Id = sno.GetString()
	err := m.db.Insert(dto)
	if err != nil {
		log.Printf("添加失败: %v\n", err)
		return model.NewFailDefault("添加失败")
	}
	matedata := make(map[string]string, 0)
	if dto.Metadata != "" {
		if err := json.Unmarshal([]byte(dto.Metadata), &matedata); err != nil {
			log.Printf("解析metadata失败: %v\n", err)
			return model.NewFailDefault("附加信息JSON格式不正确")
		}
	}
	matedata["filename"] = dto.Filename
	c := toolkit.NewCollection().Get(dto.CollectionName)
	doc := toolkit.NewDocument(c)
	documents := make([]chromem.Document, 0)
	documents = append(documents, chromem.Document{
		ID:       dto.Id,
		Metadata: matedata,
		Content:  dto.Content,
	})
	doc.Add(documents)
	return model.NewOK(dto)
}

func (m *DocumentBiz) Edit(dto *entity.Document) model.HttpResult {
	err := m.db.Update(dto, "id", dto.Id).Error
	if err != nil {
		log.Printf("修改失败: %v\n", err)
		return model.NewFailDefault("修改失败")
	}
	matedata := make(map[string]string, 0)
	if dto.Metadata != "" {
		if err := json.Unmarshal([]byte(dto.Metadata), &matedata); err != nil {
			log.Printf("解析metadata失败: %v\n", err)
			return model.NewFailDefault("附加信息JSON格式不正确")
		}
	}
	matedata["filename"] = dto.Filename
	c := toolkit.NewCollection().Get(dto.CollectionName)
	doc := toolkit.NewDocument(c)
	doc.Delete(nil, dto.Id)
	documents := make([]chromem.Document, 0)
	documents = append(documents, chromem.Document{
		ID:       dto.Id,
		Metadata: matedata,
		Content:  dto.Content,
	})
	doc.Add(documents)
	return model.NewOK(dto)
}

func (t *DocumentBiz) Delete(id string) model.HttpResult {
	info := entity.Document{}
	info.Id = id
	t.db.SelectOne(&info)
	c := toolkit.NewCollection().Get(info.CollectionName)
	doc := toolkit.NewDocument(c)
	doc.Delete(nil, id)
	err := t.db.Delete("id", id).Error
	if err != nil {
		log.Printf("删除失败: %v\n", err)
		return model.NewFailDefault("删除失败")
	}
	return model.NewOK(true)
}

func (t *DocumentBiz) DeleteByCollection(collectionName string) model.HttpResult {
	info := entity.Document{}
	info.CollectionName = collectionName
	docs, err := t.db.SelectList(&info)
	if err != nil {
		log.Printf("查询失败: %v\n", err)
		return model.NewFailDefault("查询失败")
	}
	c := toolkit.NewCollection().Get(info.CollectionName)
	for _, item := range docs {
		doc := toolkit.NewDocument(c)
		doc.Delete(nil, item.Id)
		err := t.db.Delete("id", item.Id).Error
		if err != nil {
			log.Printf("删除失败: %v\n", err)
			return model.NewFailDefault("删除失败")
		}
	}
	return model.NewOK(true)
}

type DocumentQuery struct {
	Page           *model.PageInfo `json:"page" column:"-"`
	CollectionName string          `json:"collectionName" column:"a.collection_name"  keyword:"like"`
	Filename       string          `json:"filename" column:"a.filename"  keyword:"like"`
	Content        string          `json:"content" column:"a.content" keyword:"like"`
	Metadata       string          `json:"metadata" column:"a.metadata" keyword:"like"`
}

func (m *DocumentBiz) Page(query *DocumentQuery) model.HttpResult {
	fmt.Println(query.Filename, query.CollectionName)
	sql := `select * from document a where 1=1 `
	data := make([]entity.Document, 0)
	b := sbuilder.StructToBuilder(query, sql)
	sql, values := b.Build()
	sql += " order by filename,id desc "
	result := m.db.SelectPage(&data, query.Page, sql, values...)
	return model.NewOK(result)
}

func (m *DocumentBiz) GetList(collectionName string) model.HttpResult {
	fmt.Println("向量名称", collectionName)
	sql := `select DISTINCT filename as value, filename as label from document a where 1=1 `
	data := make([]model.SelectVO, 0)
	b := sbuilder.Builder(sql)
	b.Eq("a.collection_name", collectionName)
	sql, values := b.Build()
	sql += " order by filename desc "
	err := m.db.SelectSQL(&data, sql, values...)
	if err != nil {
		log.Printf("查询列表失败: %v\n", err)
		return model.NewFailDefault("查询列表失败")
	}
	return model.NewOK(data)
}

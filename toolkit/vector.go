package toolkit

import (
	"context"
	"paas-ai/dv"
	"runtime"

	"github.com/philippgille/chromem-go"
)

var (
	db            *chromem.DB
	embeddingFunc chromem.EmbeddingFunc
)

func Init(config *dv.Config) {
	db, _ = chromem.NewPersistentDB("./db", false)
	var normalized = new(bool)
	*normalized = true
	embeddingFunc = chromem.NewEmbeddingFuncOpenAICompat(config.EmbeddingUrl, config.EmbeddingToken, config.EmbeddingModel, normalized)
}

type Collection struct {
}

func NewCollection() *Collection {
	return &Collection{}
}

// 创建一个集合
func (m *Collection) Add(collectionName string, metadata map[string]string) (*chromem.Collection, error) {
	c, err := db.GetOrCreateCollection(collectionName, metadata, embeddingFunc)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// 获取一个集合
func (m *Collection) Get(collectionName string) *chromem.Collection {
	c := db.GetCollection(collectionName, embeddingFunc)
	return c
}

// 删除集合
func (m *Collection) Delete(collectionName string, metadata map[string]string) error {
	return db.DeleteCollection(collectionName)
}

// 获取集合列表
func (m *Collection) GetList() []string {
	result := make([]string, 0)
	for k := range db.ListCollections() {
		result = append(result, k)
	}
	return result
}

type Document struct {
	c *chromem.Collection
}

func NewDocument(c *chromem.Collection) *Document {
	return &Document{c: c}
}

// 创建文档
func (m *Document) Add(documents []chromem.Document) error {
	return m.c.AddDocuments(context.Background(), documents, runtime.NumCPU())
}

// 按ID删除文档
func (m *Document) Delete(metadata map[string]string, ids ...string) error {
	return m.c.Delete(context.Background(), metadata, nil, ids...)
}

// 查询单个文档
func (m *Document) Get(id string) (chromem.Document, error) {
	return m.c.GetByID(context.Background(), id)
}

// 向量查询
func (m *Document) Query(queryText string, mResult int, metadata map[string]string) ([]chromem.Result, error) {
	count := m.c.Count()
	if mResult > count {
		mResult = count
	}
	return m.c.Query(context.Background(), queryText, mResult, metadata, nil)
}

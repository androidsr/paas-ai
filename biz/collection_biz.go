package biz

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"paas-ai/access"
	"paas-ai/entity"
	"paas-ai/toolkit"
	"strings"

	"github.com/androidsr/sc-go/model"
)

type CollectionBiz struct {
}

func NewCollectionBiz() *CollectionBiz {
	return &CollectionBiz{}
}

func (m *CollectionBiz) Add(collectionName string) model.HttpResult {
	toolkit.NewCollection().Add(collectionName, nil)
	return model.NewOK(true)
}

func (t *CollectionBiz) Delete(collectionName string) model.HttpResult {
	NewDocumentBiz().DeleteByCollection(collectionName)
	toolkit.NewCollection().Delete(collectionName, nil)
	return model.NewOK(true)
}

func (m *CollectionBiz) GetList() model.HttpResult {
	data := toolkit.NewCollection().GetList()
	fmt.Println(data)
	return model.NewOK(data)
}

func (m *CollectionBiz) Upload(collectionName string, files map[string]string) any {
	err := os.MkdirAll("temp/", os.ModePerm)
	if err != nil {
		return model.NewFailDefault(err.Error())
	}

	for filename, dataStr := range files {
		dst, err := os.Create("temp/" + filename)
		if err != nil {
			return model.NewFailDefault(err.Error())
		}
		defer dst.Close()
		idx := strings.Index(dataStr, ",")
		data, err := base64.StdEncoding.DecodeString(dataStr[idx+1:])
		if err != nil {
			return model.NewFailDefault(err.Error())
		}
		dst.Write(data)
		docs, err := access.AutoLoad("temp/" + filename)
		if err != nil {
			continue
		}
		for _, doc := range docs {
			metadata := make(map[string]string, 0)
			for k, v := range doc.Metadata {
				metadata[k] = fmt.Sprint(v)
			}
			metadata["filename"] = filename
			v, _ := json.Marshal(metadata)
			NewDocumentBiz().Add(&entity.Document{
				CollectionName: collectionName,
				Filename:       filename,
				Content:        doc.PageContent,
				Metadata:       string(v),
			})
		}
	}
	return model.NewOK(true)
}

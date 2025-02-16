package main

import (
	"context"
	"os"
	"paas-ai/biz"
	"paas-ai/entity"

	"github.com/androidsr/sc-go/mapper"
	"github.com/androidsr/sc-go/model"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Close() {
	os.Exit(0)

}

func (a *App) SetConfig(config *entity.Config) model.HttpResult {
	db := mapper.NewHelper[entity.Config]()
	err := db.SaveOrUpdate(config).Error
	if err != nil {
		return model.NewFailDefault("设置参数失败")
	}
	return model.NewOK(nil)
}

func (a *App) GetConfig() model.HttpResult {
	db := mapper.NewHelper[entity.Config]()
	data := new(entity.Config)
	data.Id = "1"
	err := db.SelectOne(data)
	if err != nil {
		return model.NewFailDefault("设置参数失败")
	}
	if data.LegalStatement == "1" {
		biz.LoadConfig()
		return model.NewOK(data)
	} else {
		return model.NewFailDefault("不同意法律声明")
	}
}

// 前端路由页面跳转
func (a *App) GoToPage(page string) {
	runtime.EventsEmit(a.ctx, "go-to-page", page)
}

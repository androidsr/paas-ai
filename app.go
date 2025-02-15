package main

import (
	"context"
	"fmt"
	"os"
	"paas-ai/biz"
	"paas-ai/entity"

	"github.com/androidsr/sc-go/mapper"
	"github.com/androidsr/sc-go/model"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
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

func (a *App) AddMenus() *menu.Menu {
	appMenu := menu.NewMenu()

	aiChat := appMenu.AddSubmenu("对话管理")
	i := 1
	aiChat.AddText("AI对话", keys.CmdOrCtrl(fmt.Sprintf("%d", i)), func(cd *menu.CallbackData) {
		a.GoToPage("/aichat")
	})
	i++
	aiChat.AddSeparator()
	aiChat.AddText("提示词管理", keys.CmdOrCtrl(fmt.Sprintf("%d", i)), func(cd *menu.CallbackData) {
		a.GoToPage("/prompt")
	})

	knowledge := appMenu.AddSubmenu("知识库管理")
	i++
	knowledge.AddText("集合管理", keys.CmdOrCtrl(fmt.Sprintf("%d", i)), func(cd *menu.CallbackData) {
		a.GoToPage("/collection")
	})
	i++
	knowledge.AddSeparator()
	knowledge.AddText("文档管理", keys.CmdOrCtrl(fmt.Sprintf("%d", i)), func(cd *menu.CallbackData) {
		a.GoToPage("/document")
	})

	tools := appMenu.AddSubmenu("工具管理")
	i++
	tools.AddSeparator()
	tools.AddText("函数定义", keys.CmdOrCtrl(fmt.Sprintf("%d", i)), func(cd *menu.CallbackData) {
		a.GoToPage("/function")
	})
	i++
	tools.AddSeparator()
	tools.AddText("函数实现", keys.CmdOrCtrl(fmt.Sprintf("%d", i)), func(cd *menu.CallbackData) {
		a.GoToPage("/func-impl")
	})
	i++
	tools.AddSeparator()
	tools.AddText("数据库管理", keys.CmdOrCtrl(fmt.Sprintf("%d", i)), func(cd *menu.CallbackData) {
		a.GoToPage("/dbconfig")
	})
	i++
	tools.AddSeparator()
	tools.AddText("流程管理", keys.CmdOrCtrl(fmt.Sprintf("%d", i)), func(cd *menu.CallbackData) {
		a.GoToPage("/flow")
	})

	other := appMenu.AddSubmenu("设置")
	i++
	other.AddSeparator()
	other.AddText("渠道管理", keys.CmdOrCtrl(fmt.Sprintf("%d", i)), func(cd *menu.CallbackData) {
		a.GoToPage("/aichannel")
	})
	i++
	other.AddSeparator()
	other.AddText("参数配置", keys.CmdOrCtrl(fmt.Sprintf("%d", i)), func(cd *menu.CallbackData) {
		a.GoToPage("/config")
	})
	i++
	other.AddSeparator()
	other.AddText("关于", keys.CmdOrCtrl(fmt.Sprintf("%d", i)), func(cd *menu.CallbackData) {
		a.GoToPage("/about")
	})

	//appMenu.AddSeparator()

	/* if runtime.GOOS == "darwin" {
		AppMenu.Append(menu.EditMenu()) // on macos platform, we should append EditMenu to enable Cmd+C,Cmd+V,Cmd+Z... shortcut
	} */

	return appMenu
}

// 前端路由页面跳转
func (a *App) GoToPage(page string) {
	runtime.EventsEmit(a.ctx, "go-to-page", page)
}

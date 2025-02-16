package main

import (
	"embed"
	"paas-ai/biz"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	/* if !toolkit.Signature("10-F6-0A-AD-A3-81,04-42-1A-D7-FE-8B", "2025-03-01") {
		return
	} */
	dbBiz := biz.NewDbConfigBiz()
	fwBiz := biz.NewFwConfigBiz()
	functionBiz := biz.NewFunctionBiz()
	functionImplBiz := biz.NewFunctionImplBiz()
	roleBiz := biz.NewPromptBiz()
	collectionBiz := biz.NewCollectionBiz()
	documentBiz := biz.NewDocumentBiz()
	chatBiz := biz.NewChatBiz()
	channelBiz := biz.NewAiChannelBiz()

	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "那个谁AI桌面端",
		Width:  1100,
		Height: 700,
		//HideWindowOnClose: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.startup,
		Menu:             app.AddMenus(),
		Bind: []interface{}{
			app,
			dbBiz,
			fwBiz,
			functionBiz,
			functionImplBiz,
			roleBiz,
			collectionBiz,
			documentBiz,
			chatBiz,
			channelBiz,
		},
		Windows: &windows.Options{
			WebviewIsTransparent:              false,
			WindowIsTranslucent:               false,
			BackdropType:                      windows.Mica,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
			WebviewUserDataPath:               "",
			WebviewBrowserPath:                "",
			Theme:                             windows.Light,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
	//systray.Run(onReady, onExit)

}

/*
func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("那个谁AI桌面端")
	systray.SetTooltip("那个谁AI桌面端")
	mQuit := systray.AddMenuItem("退出", "点击退出")
	mQuit.SetIcon(icon.Data)
}

func onExit() {
	os.Exit(0)
}
*/

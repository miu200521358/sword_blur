//go:build windows
// +build windows

package main

import (
	"embed"
	"fmt"
	"log"
	"runtime"

	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mwidget"
	"github.com/miu200521358/sword_blur/pkg/model"
	"github.com/miu200521358/sword_blur/pkg/ui"
)

func init() {
	runtime.LockOSThread()

	walk.AppendToWalkInit(func() {
		walk.MustRegisterWindowClass(mwidget.FilePickerClass)
		walk.MustRegisterWindowClass(mwidget.MotionPlayerClass)
		walk.MustRegisterWindowClass(mwidget.ConsoleViewClass)
	})
}

var env string

//go:embed app/*
var appFiles embed.FS

//go:embed i18n/*
var appI18nFiles embed.FS

func main() {
	var mWindow *mwidget.MWindow
	var err error

	appConfig := mconfig.LoadAppConfig(appFiles)
	appConfig.Env = env
	mi18n.Initialize(appI18nFiles)

	if appConfig.IsEnvProd() || appConfig.IsEnvDev() {
		defer mwidget.RecoverFromPanic(mWindow)
	}

	iconImg, err := mconfig.LoadIconFile(appFiles)
	mwidget.CheckError(err, nil, mi18n.T("アイコン生成エラー"))

	glWindow, err := mwidget.NewGlWindow(512, 768, 0, iconImg, appConfig, nil, nil)
	mwidget.CheckError(err, mWindow, mi18n.T("ビューワーウィンドウ生成エラー"))

	go func() {
		mWindow, err = mwidget.NewMWindow(512, 768, ui.GetMenuItems, iconImg, appConfig, true)
		mwidget.CheckError(err, nil, mi18n.T("メインウィンドウ生成エラー"))

		blurModel := model.NewBlurModel()

		step1Page, err := ui.NewStep1TabPage(mWindow)
		mwidget.CheckError(err, nil, mi18n.T("タブページ生成エラー"))

		step2Page, err := ui.NewStep2TabPage(mWindow, step1Page, blurModel)
		mwidget.CheckError(err, nil, mi18n.T("タブページ生成エラー"))

		step3Page, err := ui.NewStep3TabPage(mWindow, step2Page, blurModel)
		mwidget.CheckError(err, nil, mi18n.T("タブページ生成エラー"))

		step4Page, err := ui.NewStep4TabPage(mWindow, step3Page, blurModel)
		mwidget.CheckError(err, nil, mi18n.T("タブページ生成エラー"))

		step5Page, err := ui.NewStep5TabPage(mWindow, step4Page, blurModel)
		mwidget.CheckError(err, nil, mi18n.T("タブページ生成エラー"))

		// 関数紐付け切り替え
		mWindow.TabWidget.CurrentIndexChanged().Attach(func() {
			if mWindow.TabWidget.CurrentIndex() == 2 {
				mWindow.GetMainGlWindow().SetFuncWorldPos(step3Page.FuncWorldPos(blurModel))
			} else if mWindow.TabWidget.CurrentIndex() == 3 {
				mWindow.GetMainGlWindow().SetFuncWorldPos(step4Page.FuncWorldPos(blurModel))
			} else if mWindow.TabWidget.CurrentIndex() == 4 {
				mWindow.GetMainGlWindow().SetFuncWorldPos(step5Page.FuncWorldPos(blurModel))
			} else {
				mWindow.GetMainGlWindow().SetFuncWorldPos(nil)
			}
		})

		// コンソールはタブ外に表示
		mWindow.ConsoleView, err = mwidget.NewConsoleView(mWindow, 256, 30)
		mwidget.CheckError(err, mWindow, mi18n.T("コンソール生成エラー"))
		log.SetOutput(mWindow.ConsoleView)

		glWindow.SetMotionPlayer(step1Page.Items.MotionPlayer)
		glWindow.SetTitle(fmt.Sprintf("%s %s", mWindow.Title(), mi18n.T("ビューワー")))
		mWindow.AddGlWindow(glWindow)

		mWindow.AsFormBase().Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
			go func() {
				mWindow.GetMainGlWindow().IsClosedChannel <- true
			}()
			mWindow.Close()
		})

		step2Page.SetEnabled(false)

		mWindow.Center()
		mWindow.Run()
	}()

	glWindow.Run()
}

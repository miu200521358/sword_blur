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

//go:embed resources/*
var resourceFiles embed.FS

func main() {
	var mWindow *mwidget.MWindow
	var err error

	appConfig := mconfig.LoadAppConfig(resourceFiles)
	appConfig.Env = env
	mi18n.Initialize(resourceFiles)

	if appConfig.IsEnvProd() || appConfig.IsEnvDev() {
		defer mwidget.RecoverFromPanic(mWindow)
	}

	glWindow, err := mwidget.NewGlWindow(mi18n.T("ビューワー"), 512, 768, 0, resourceFiles, nil, nil)
	mwidget.CheckError(err, mWindow, mi18n.T("ビューワーウィンドウ生成エラー"))

	go func() {
		mWindow, err = mwidget.NewMWindow(resourceFiles, appConfig, true, 512, 768, ui.GetMenuItems)
		mwidget.CheckError(err, nil, mi18n.T("メインウィンドウ生成エラー"))

		step1Page, err := ui.NewStep1TabPage(mWindow, resourceFiles)
		mwidget.CheckError(err, nil, mi18n.T("タブページ生成エラー"))

		step2Page, err := ui.NewStep2TabPage(mWindow, step1Page)
		mwidget.CheckError(err, nil, mi18n.T("タブページ生成エラー"))

		step3Page, err := ui.NewStep3TabPage(mWindow, step2Page)
		mwidget.CheckError(err, nil, mi18n.T("タブページ生成エラー"))

		step4Page, err := ui.NewStep4TabPage(mWindow, step3Page)
		mwidget.CheckError(err, nil, mi18n.T("タブページ生成エラー"))

		// 関数紐付け切り替え
		mWindow.TabWidget.CurrentIndexChanged().Attach(func() {
			if mWindow.TabWidget.CurrentIndex() == 2 {
				mWindow.GetMainGlWindow().SetFuncWorldPos(step3Page.Items.FuncWorldPos)
			} else if mWindow.TabWidget.CurrentIndex() == 3 {
				mWindow.GetMainGlWindow().SetFuncWorldPos(step4Page.Items.FuncWorldPos)
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

		mWindow.Center()
		mWindow.Run()
	}()

	glWindow.Run()
}

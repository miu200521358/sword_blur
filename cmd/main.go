//go:build windows
// +build windows

package main

import (
	"embed"
	"log"
	"runtime"

	"github.com/miu200521358/sword_blur/pkg/ui"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/interface/app"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/interface/viewer"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
)

var env string

func init() {
	runtime.LockOSThread()

	// システム上のすべての論理プロセッサを使用させる
	runtime.GOMAXPROCS(runtime.NumCPU())

	walk.AppendToWalkInit(func() {
		walk.MustRegisterWindowClass(widget.FilePickerClass)
		walk.MustRegisterWindowClass(widget.MotionPlayerClass)
		walk.MustRegisterWindowClass(widget.ConsoleViewClass)
	})
}

//go:embed app/*
var appFiles embed.FS

//go:embed i18n/*
var appI18nFiles embed.FS

func main() {
	appConfig := mconfig.LoadAppConfig(appFiles)
	appConfig.Env = env
	mi18n.Initialize(appI18nFiles)

	mApp := app.NewMApp(appConfig)

	controlState := controller.NewControlState(mApp)
	controlState.Run()

	go func() {
		// 操作ウィンドウは別スレッドで起動
		controlWindow := controller.NewControlWindow(appConfig, controlState, ui.GetMenuItems, 2)
		mApp.SetControlWindow(controlWindow)

		controlWindow.InitTabWidget()
		ui.NewToolState(mApp, controlWindow)

		consoleView := widget.NewConsoleView(controlWindow.MainWindow, 256, 50)
		log.SetOutput(consoleView)

		mApp.ControllerRun()
	}()

	mApp.AddViewWindow(viewer.NewViewWindow(mApp.ViewerCount(), appConfig, mApp, "メインビューワー", nil))
	mApp.AddViewWindow(viewer.NewViewWindow(mApp.ViewerCount(), appConfig, mApp, "プレビュービューワー", mApp.MainViewWindow().GetWindow()))

	mApp.ExtendAnimationState(0, 0)
	mApp.ExtendAnimationState(1, 0)

	mApp.Center()
	mApp.ViewerRun()
}

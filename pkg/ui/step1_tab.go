package ui

import (
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/sword_blur/pkg/usecase"
	"github.com/miu200521358/walk/pkg/walk"
)

func newStep1Tab(controlWindow *controller.ControlWindow, toolState *ToolState) {
	toolState.Step1 = widget.NewMTabPage("Step.1")
	controlWindow.AddTabPage(toolState.Step1.TabPage)

	toolState.Step1.SetLayout(walk.NewVBoxLayout())

	var err error

	{
		// Step1. ファイル選択文言
		label, err := walk.NewTextLabel(toolState.Step1)
		if err != nil {
			widget.RaiseError(err)
		}
		label.SetText(mi18n.T("Step1Label"))
	}

	walk.NewVSeparator(toolState.Step1)

	{
		toolState.OriginalPmxPicker = widget.NewPmxReadFilePicker(
			controlWindow,
			toolState.Step1,
			"OriginalPmx",
			mi18n.T("ブレ生成対象モデル(Pmx)"),
			mi18n.T("ブレ生成対象モデルPmxファイルを選択してください"),
			mi18n.T("ブレ生成対象モデルの使い方"))

		toolState.OriginalPmxPicker.SetOnPathChanged(func(path string) {
			if data, err := toolState.OriginalPmxPicker.Load(); err == nil {
				// 出力パス設定
				outputPath := mutils.CreateOutputPath(path, "blur")
				toolState.OutputPmxPicker.SetPath(outputPath)

				toolState.BlurModel.Model = data.(*pmx.PmxModel)

				if !toolState.OriginalVmdPicker.Exists() {
					toolState.BlurModel.Motion = vmd.NewVmdMotion("")
				}
			} else {
				mlog.E(mi18n.T("読み込み失敗"), err)
			}
		})
	}

	{
		toolState.OriginalVmdPicker = widget.NewVmdVpdReadFilePicker(
			controlWindow,
			toolState.Step1,
			"OriginalVmd",
			mi18n.T("確認用モーション(Vmd/Vpd)"),
			mi18n.T("確認用モーション(Vmd/Vpd)ファイルを選択してください"),
			mi18n.T("確認用モーションの使い方"))

		toolState.OriginalVmdPicker.SetOnPathChanged(func(path string) {
			if data, err := toolState.OriginalVmdPicker.Load(); err == nil {
				motion := data.(*vmd.VmdMotion)
				controlWindow.UpdateMaxFrame(motion.MaxFrame())
				toolState.BlurModel.Motion = data.(*vmd.VmdMotion)
			} else {
				mlog.E(mi18n.T("読み込み失敗"), err)
			}
		})
	}

	{
		toolState.OutputPmxPicker = widget.NewPmxSaveFilePicker(
			controlWindow,
			toolState.Step1,
			mi18n.T("出力モデル(Pmx)"),
			mi18n.T("出力モデル(Pmx)ファイルパスを指定してください"),
			mi18n.T("出力モデルの使い方"))
	}

	walk.NewVSpacer(toolState.Step1)

	// OKボタン
	{
		toolState.Step1OkButton, err = walk.NewPushButton(toolState.Step1)
		if err != nil {
			widget.RaiseError(err)
		}
		toolState.Step1OkButton.SetText(mi18n.T("次へ進む"))
		toolState.Step1OkButton.Clicked().Attach(toolState.onClickStep1Ok)
	}

	toolState.App.SetFuncGetModels(
		func() [][]*pmx.PmxModel {
			return [][]*pmx.PmxModel{
				{toolState.BlurModel.Model},
				{toolState.BlurModel.OutputModel},
			}
		},
	)

	toolState.App.SetFuncGetMotions(
		func() [][]*vmd.VmdMotion {
			return [][]*vmd.VmdMotion{
				{toolState.BlurModel.Motion},
				{toolState.BlurModel.OutputMotion},
			}
		},
	)
}

func (toolState *ToolState) onClickStep1Ok() {
	if !toolState.OriginalPmxPicker.Exists() {
		mlog.ILT(mi18n.T("設定失敗"), mi18n.T("Step1失敗"))
		return
	}

	// 追加セットアップ
	usecase.SetupModel(toolState.BlurModel)

	toolState.BlurModel.OutputModelPath = toolState.OutputPmxPicker.GetPath()
	toolState.BlurModel.OutputModel = nil
	toolState.BlurModel.OutputMotion = nil

	// Step2.材質リストボックス設定
	toolState.MaterialListBox.SetMaterials(
		toolState.BlurModel.Model.Materials,
		toolState.onChangeMaterialListBox())

	toolState.ControlWindow.SetTabIndex(1) // Step2へ移動
	toolState.SetEnabled(2)
	mlog.IL(mi18n.T("Step1成功"))
}

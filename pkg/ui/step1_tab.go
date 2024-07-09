package ui

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/mlib_go/pkg/mwidget"
	"github.com/miu200521358/mlib_go/pkg/pmx"
	"github.com/miu200521358/mlib_go/pkg/vmd"
	"github.com/miu200521358/sword_blur/pkg/model"
	"github.com/miu200521358/sword_blur/pkg/usecase"
	"github.com/miu200521358/walk/pkg/walk"
)

func NewStep1TabPage(mWindow *mwidget.MWindow) (*Step1TabPage, error) {
	page, err := mwidget.NewMTabPage(mWindow, mWindow.TabWidget, "Step. 1")
	if err != nil {
		return nil, err
	}
	page.SetLayout(walk.NewVBoxLayout())

	stp := &Step1TabPage{
		MTabPage: page,
		mWindow:  mWindow,
		Items:    &Step1Items{},
	}

	stp.Items.composite, err = walk.NewComposite(stp)
	if err != nil {
		return nil, err
	}
	layout := walk.NewVBoxLayout()
	stp.Items.composite.SetLayout(layout)

	// Step1. ファイル選択文言
	stp.Items.label, err = walk.NewTextLabel(stp.Items.composite)
	if err != nil {
		return nil, err
	}
	stp.Items.label.SetText(mi18n.T("Step1Label"))

	walk.NewVSeparator(stp.Items.composite)

	stp.Items.OriginalPmxPicker, err = (mwidget.NewPmxReadFilePicker(
		mWindow,
		stp.Items.composite,
		"org_pmx",
		mi18n.T("ブレ生成対象モデル(Pmx)"),
		mi18n.T("ブレ生成対象モデルPmxファイルを選択してください"),
		mi18n.T("ブレ生成対象モデルの使い方"),
		nil))
	if err != nil {
		return nil, err
	}

	stp.Items.OriginalVmdPicker, err = (mwidget.NewVmdVpdReadFilePicker(
		mWindow,
		stp.Items.composite,
		"vmd",
		mi18n.T("確認用モーション(Vmd/Vpd)"),
		mi18n.T("確認用モーション(Vmd/Vpd)ファイルを選択してください"),
		mi18n.T("確認用モーションの使い方"),
		func(path string) {}))
	if err != nil {
		return nil, err
	}

	stp.Items.OutputPmxPicker, err = (mwidget.NewPmxSaveFilePicker(
		mWindow,
		stp.Items.composite,
		mi18n.T("出力モデル(Pmx)"),
		mi18n.T("出力モデル(Pmx)ファイルパスを指定してください"),
		mi18n.T("出力モデルの使い方"),
		func(path string) {}))
	if err != nil {
		return nil, err
	}

	// プレイヤーBox
	stp.Items.MotionPlayer, err = mwidget.NewMotionPlayer(stp.Items.composite, mWindow)
	if err != nil {
		return nil, err
	}

	stp.Items.MotionPlayer.OnPlay = func(isPlaying bool) error {
		// 入力欄は全部再生中は無効化
		stp.Items.OriginalPmxPicker.SetEnabled(!isPlaying)
		stp.Items.OriginalVmdPicker.SetEnabled(!isPlaying)
		stp.Items.OutputPmxPicker.SetEnabled(!isPlaying)
		stp.Items.MotionPlayer.SetEnabled(!isPlaying)
		stp.Items.MotionPlayer.PlayButton.SetEnabled(true)
		for _, glWindow := range mWindow.GlWindows {
			glWindow.TriggerPlay(isPlaying)
		}

		return nil
	}

	stp.Items.OriginalPmxPicker.PathLineEdit.SetFocus()
	stp.Items.OutputPmxPicker.SetEnabled(false)
	stp.Items.MotionPlayer.SetEnabled(false)

	// 下を埋める
	stretchWidget, err := walk.NewTextLabel(stp.Items.composite)
	if err != nil {
		return nil, err
	}
	layout.SetStretchFactor(stretchWidget, 100)

	// OKボタン
	stp.Items.okButton, err = walk.NewPushButton(stp.Items.composite)
	if err != nil {
		return nil, err
	}

	stp.Items.okButton.SetText(mi18n.T("次へ進む"))

	// 元モデル設定時
	stp.Items.OriginalPmxPicker.OnPathChanged = stp.funcOriginalPmxModelChanged()

	// モーション設定時
	stp.Items.OriginalVmdPicker.OnPathChanged = stp.funcOriginalVmdModelChanged()

	return stp, nil
}

// ------------------------------

type Step1TabPage struct {
	*mwidget.MTabPage
	mWindow  *mwidget.MWindow
	Items    *Step1Items
	nextStep *Step2TabPage
}

// Step1. OKボタンクリック時
func (step1Page Step1TabPage) funcOkButton(blurModel *model.BlurModel) {
	if step1Page.Items.OriginalPmxPicker.Exists() {
		model := step1Page.Items.OriginalPmxPicker.GetCache().(*pmx.PmxModel)
		var motion *vmd.VmdMotion
		if step1Page.Items.OriginalVmdPicker.IsCached() {
			motion = step1Page.Items.OriginalVmdPicker.GetCache().(*vmd.VmdMotion)
		} else {
			motion = vmd.NewVmdMotion("")
		}

		blurModel.Model = model
		blurModel.Motion = motion
		blurModel.OutputModelPath = step1Page.Items.OutputPmxPicker.GetPath()
		blurModel.OutputModel = nil
		blurModel.OutputMotion = nil

		// 追加セットアップ
		usecase.SetupModel(blurModel)

		step1Page.Items.MotionPlayer.SetEnabled(true)
		step1Page.Items.MotionPlayer.SetValue(0)
		step1Page.nextStep.SetEnabled(true)

		// 材質リストボックス設定
		step1Page.nextStep.Items.MaterialListBox.SetMaterials(
			blurModel.Model.Materials,
			step1Page.nextStep.funcMaterialListBoxChanged(blurModel))

		step1Page.mWindow.TabWidget.SetCurrentIndex(1) // Step2へ移動
		mlog.IL(mi18n.T("Step1成功"))
	} else {
		step1Page.nextStep.SetEnabled(false)
		mlog.ILT(mi18n.T("設定失敗"), mi18n.T("Step1失敗"))
		return
	}
}

// Step1. 元モデル設定時
func (stp Step1TabPage) funcOriginalPmxModelChanged() func(path string) {
	return func(path string) {
		func(path string) {
			isExist, err := mutils.ExistsFile(path)
			if !isExist || err != nil {
				stp.Items.OutputPmxPicker.PathLineEdit.SetText("")
				mlog.IL(mi18n.T("Step1失敗"))
				return
			}

			stp.Items.OutputPmxPicker.SetEnabled(true)

			// 出力パス設定
			dir, file := filepath.Split(path)
			ext := filepath.Ext(file)
			fileName := fmt.Sprintf("%s_blur_%s%s", file[:len(file)-len(ext)], time.Now().Format("20060102_150405"), ext)
			outputPath := filepath.Join(dir, fileName)
			stp.Items.OutputPmxPicker.PathLineEdit.SetText(outputPath)

			if stp.Items.OriginalPmxPicker.Exists() {
				data, err := stp.Items.OriginalPmxPicker.GetData()
				if err != nil {
					mlog.E(mi18n.T("Pmxファイル読み込みエラー"), err.Error())
					mlog.IL(mi18n.T("Step1失敗"))
					return
				}
				model := data.(*pmx.PmxModel)

				var motion *vmd.VmdMotion
				if stp.Items.OriginalVmdPicker.IsCached() {
					motion = stp.Items.OriginalVmdPicker.GetCache().(*vmd.VmdMotion)
				} else {
					motion = vmd.NewVmdMotion("")
				}

				go func() {
					stp.mWindow.GetMainGlWindow().FrameChannel <- 0
					stp.mWindow.GetMainGlWindow().IsPlayingChannel <- false
					stp.mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]*mwidget.ModelSet{0: {NextModel: model, NextMotion: motion}}
				}()

			} else {
				go func() {
					stp.mWindow.GetMainGlWindow().RemoveModelSetIndexChannel <- 0
				}()
				stp.Items.MotionPlayer.SetEnabled(false)
			}
		}(path)
	}
}

// Step1. モーション設定時
func (stp *Step1TabPage) funcOriginalVmdModelChanged() func(path string) {
	return func(path string) {
		if stp.Items.OriginalVmdPicker.Exists() {
			motionData, err := stp.Items.OriginalVmdPicker.GetData()
			if err != nil {
				mlog.E(mi18n.T("Vmdファイル読み込みエラー"), err.Error())
				return
			}
			motion := motionData.(*vmd.VmdMotion)

			stp.Items.MotionPlayer.SetEnabled(true)
			stp.Items.MotionPlayer.SetRange(0, motion.GetMaxFrame()+1)
			stp.Items.MotionPlayer.SetValue(0)

			if stp.Items.OriginalPmxPicker.Exists() {
				go func() {
					stp.mWindow.GetMainGlWindow().FrameChannel <- 0
					stp.mWindow.GetMainGlWindow().IsPlayingChannel <- false
					stp.mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]*mwidget.ModelSet{0: {NextMotion: motion}}
				}()
			} else {
				go func() {
					stp.mWindow.GetMainGlWindow().RemoveModelSetIndexChannel <- 0
				}()
			}
		}
	}
}

// ------------------------------

type stepItems struct {
	label     *walk.TextLabel
	composite *walk.Composite
}

func (si *stepItems) SetEnabled(enabled bool) {
	si.label.SetEnabled(enabled)
}

func (si *stepItems) Dispose() {
	si.label.Dispose()
	si.composite.Dispose()
}

type Step1Items struct {
	stepItems
	MotionPlayer      *mwidget.MotionPlayer
	OriginalPmxPicker *mwidget.FilePicker
	OriginalVmdPicker *mwidget.FilePicker
	OutputPmxPicker   *mwidget.FilePicker
	okButton          *walk.PushButton
}

func (si *Step1Items) SetEnabled(enabled bool) {
	si.stepItems.SetEnabled(enabled)
	si.MotionPlayer.SetEnabled(enabled)
	si.OriginalPmxPicker.SetEnabled(enabled)
	si.OriginalVmdPicker.SetEnabled(enabled)
	si.OutputPmxPicker.SetEnabled(enabled)
}

func (si *Step1Items) Dispose() {
	si.stepItems.Dispose()
	si.MotionPlayer.Dispose()
	si.OriginalPmxPicker.Dispose()
	si.OriginalVmdPicker.Dispose()
	si.OutputPmxPicker.Dispose()
}

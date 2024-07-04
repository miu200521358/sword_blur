package ui

import (
	"fmt"
	"path/filepath"

	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/mlib_go/pkg/mwidget"
	"github.com/miu200521358/mlib_go/pkg/pmx"
	"github.com/miu200521358/mlib_go/pkg/vmd"
	"github.com/miu200521358/walk/pkg/walk"
)

func NewStep2TabPage(mWindow *mwidget.MWindow, step1Page *Step1TabPage) (*Step2TabPage, error) {
	page, err := mwidget.NewMTabPage(mWindow, mWindow.TabWidget, "Step. 2")
	if err != nil {
		return nil, err
	}
	page.SetLayout(walk.NewVBoxLayout())

	stp := &Step2TabPage{
		MTabPage: page,
		mWindow:  mWindow,
		Items:    &Step2Items{},
	}

	// Step2. 材質選択

	stp.Items.composite, err = walk.NewComposite(stp)
	if err != nil {
		return nil, err
	}
	stp.Items.composite.SetLayout(walk.NewVBoxLayout())

	stp.Items.label, err = walk.NewTextLabel(stp.Items.composite)
	if err != nil {
		return nil, err
	}
	stp.Items.label.SetText(mi18n.T("Step2Label"))

	walk.NewVSeparator(stp.Items.composite)

	// 材質選択リストボックス
	stp.Items.MaterialListBox, err = NewMaterialListBox(stp.Items.composite)
	if err != nil {
		return nil, err
	}

	// 元モデル設定時
	step1Page.Items.OriginalPmxPicker.OnPathChanged = func(path string) {
		func(path string) {
			isExist, err := mutils.ExistsFile(path)
			if !isExist || err != nil {
				step1Page.Items.OutputPmxPicker.PathLineEdit.SetText("")
				return
			}

			step1Page.Items.OutputPmxPicker.SetEnabled(true)

			// 出力パス設定
			dir, file := filepath.Split(path)
			ext := filepath.Ext(file)
			outputPath := filepath.Join(dir, file[:len(file)-len(ext)]+"_blur"+ext)
			step1Page.Items.OutputPmxPicker.PathLineEdit.SetText(outputPath)

			if step1Page.Items.OriginalPmxPicker.Exists() {
				data, err := step1Page.Items.OriginalPmxPicker.GetData()
				if err != nil {
					mlog.E(mi18n.T("Pmxファイル読み込みエラー"), err.Error())
					return
				}
				model := data.(*pmx.PmxModel)
				var motion *vmd.VmdMotion
				if step1Page.Items.OriginalVmdPicker.IsCached() {
					motion = step1Page.Items.OriginalVmdPicker.GetCache().(*vmd.VmdMotion)
				} else {
					motion = vmd.NewVmdMotion("")
				}

				step1Page.Items.MotionPlayer.SetEnabled(true)
				step1Page.Items.MotionPlayer.SetValue(0)
				stp.SetEnabled(true)
				stp.Items.MaterialListBox.SetMaterials(model.Materials, func(indexes []int) {
					mlog.I(fmt.Sprintf("材質Indexes: %v", indexes))
				})
				mlog.IL(mi18n.T("Step1モデル設定完了"))

				go func() {
					mWindow.GetMainGlWindow().FrameChannel <- 0
					mWindow.GetMainGlWindow().IsPlayingChannel <- false
					mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]mwidget.ModelSet{0: {Model: model, Motion: motion}}
				}()
			} else {
				go func() {
					mWindow.GetMainGlWindow().RemoveModelSetIndexChannel <- 0
				}()
				step1Page.Items.MotionPlayer.SetEnabled(false)
				stp.SetEnabled(false)
			}
		}(path)
	}

	// モーション設定時
	step1Page.Items.OriginalVmdPicker.OnPathChanged = func(path string) {
		if step1Page.Items.OriginalVmdPicker.Exists() {
			motionData, err := step1Page.Items.OriginalVmdPicker.GetData()
			if err != nil {
				mlog.E(mi18n.T("Vmdファイル読み込みエラー"), err.Error())
				return
			}
			motion := motionData.(*vmd.VmdMotion)

			step1Page.Items.MotionPlayer.SetRange(0, motion.GetMaxFrame()+1)
			step1Page.Items.MotionPlayer.SetValue(0)

			if step1Page.Items.OriginalPmxPicker.Exists() {
				model := step1Page.Items.OriginalPmxPicker.GetCache().(*pmx.PmxModel)

				go func() {
					mWindow.GetMainGlWindow().FrameChannel <- 0
					mWindow.GetMainGlWindow().IsPlayingChannel <- false
					mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]mwidget.ModelSet{0: {Model: model, Motion: motion}}
				}()
			} else {
				go func() {
					mWindow.GetMainGlWindow().RemoveModelSetIndexChannel <- 0
				}()
			}
		}
	}

	stp.SetEnabled(false)

	return stp, nil
}

// ------------------------------

type Step2TabPage struct {
	*mwidget.MTabPage
	mWindow *mwidget.MWindow
	Items   *Step2Items
}

// ------------------------------

type Step2Items struct {
	stepItems
	MaterialListBox *MaterialListBox
}

func (si *Step2Items) SetEnabled(enabled bool) {
	si.stepItems.SetEnabled(enabled)
	si.MaterialListBox.SetEnabled(enabled)
}

func (si *Step2Items) Dispose() {
	si.stepItems.Dispose()
	si.MaterialListBox.Dispose()
}

package ui

import (
	"fmt"
	"path/filepath"
	"slices"
	"time"

	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/mlib_go/pkg/mwidget"
	"github.com/miu200521358/mlib_go/pkg/pmx"
	"github.com/miu200521358/mlib_go/pkg/vmd"
	"github.com/miu200521358/sword_blur/pkg/usecase"
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
		prevStep: step1Page,
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

	// OKボタン
	stp.Items.okButton, err = walk.NewPushButton(stp.Items.composite)
	if err != nil {
		return nil, err
	}
	stp.Items.okButton.SetText(mi18n.T("次へ進む"))

	// 元モデル設定時
	step1Page.Items.OriginalPmxPicker.OnPathChanged = func(path string) {
		func(path string) {
			isExist, err := mutils.ExistsFile(path)
			if !isExist || err != nil {
				step1Page.Items.OutputPmxPicker.PathLineEdit.SetText("")
				mlog.IL(mi18n.T("Step1失敗"))
				return
			}

			step1Page.Items.OutputPmxPicker.SetEnabled(true)

			// 出力パス設定
			dir, file := filepath.Split(path)
			ext := filepath.Ext(file)
			fileName := fmt.Sprintf("%s_blur_%s%s", file[:len(file)-len(ext)], time.Now().Format("20060102_150405"), ext)
			outputPath := filepath.Join(dir, fileName)
			step1Page.Items.OutputPmxPicker.PathLineEdit.SetText(outputPath)

			if step1Page.Items.OriginalPmxPicker.Exists() {
				data, err := step1Page.Items.OriginalPmxPicker.GetData()
				if err != nil {
					mlog.E(mi18n.T("Pmxファイル読み込みエラー"), err.Error())
					mlog.IL(mi18n.T("Step1失敗"))
					return
				}
				model := data.(*pmx.PmxModel)

				var motion *vmd.VmdMotion
				if step1Page.Items.OriginalVmdPicker.IsCached() {
					motion = step1Page.Items.OriginalVmdPicker.GetCache().(*vmd.VmdMotion)
				} else {
					motion = vmd.NewVmdMotion("")
				}

				go func() {
					mWindow.GetMainGlWindow().FrameChannel <- 0
					mWindow.GetMainGlWindow().IsPlayingChannel <- false
					mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]*mwidget.ModelSet{0: {NextModel: model, NextMotion: motion}}
				}()

			} else {
				go func() {
					mWindow.GetMainGlWindow().RemoveModelSetIndexChannel <- 0
				}()
				step1Page.Items.MotionPlayer.SetEnabled(false)
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

			step1Page.Items.MotionPlayer.SetEnabled(true)
			step1Page.Items.MotionPlayer.SetRange(0, motion.GetMaxFrame()+1)
			step1Page.Items.MotionPlayer.SetValue(0)

			if step1Page.Items.OriginalPmxPicker.Exists() {
				go func() {
					mWindow.GetMainGlWindow().FrameChannel <- 0
					mWindow.GetMainGlWindow().IsPlayingChannel <- false
					mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]*mwidget.ModelSet{0: {NextMotion: motion}}
				}()
			} else {
				go func() {
					mWindow.GetMainGlWindow().RemoveModelSetIndexChannel <- 0
				}()
			}
		}
	}

	stp.SetEnabled(false)

	// Step1. OKボタンクリック時
	step1Page.Items.okButton.Clicked().Attach(func() {
		if step1Page.Items.OriginalPmxPicker.Exists() {
			model := step1Page.Items.OriginalPmxPicker.GetCache().(*pmx.PmxModel)
			var motion *vmd.VmdMotion
			if step1Page.Items.OriginalVmdPicker.IsCached() {
				motion = step1Page.Items.OriginalVmdPicker.GetCache().(*vmd.VmdMotion)
			} else {
				motion = vmd.NewVmdMotion("")
			}

			// 追加セットアップ
			usecase.SetupModel(model)

			step1Page.Items.MotionPlayer.SetEnabled(true)
			step1Page.Items.MotionPlayer.SetValue(0)
			stp.SetEnabled(true)
			stp.Items.MaterialListBox.SetMaterials(model.Materials, func(indexes []int) {
				invisibleMaterialIndexes := make([]int, 0)
				for i := range model.Materials.Len() {
					material := model.Materials.Get(i)
					mf := vmd.NewMorphFrame(0)
					if slices.Contains(indexes, i) {
						mf.Ratio = 0.0
					} else {
						mf.Ratio = 1.0
						invisibleMaterialIndexes = append(invisibleMaterialIndexes, i)
					}
					motion.AppendMorphFrame(usecase.GetVisibleMorphName(material), mf)
				}

				go func() {
					mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]*mwidget.ModelSet{0: {NextMotion: motion, NextInvisibleMaterialIndexes: invisibleMaterialIndexes}}
				}()
			})

			stp.mWindow.TabWidget.SetCurrentIndex(1) // Step2へ移動
			mlog.IL(mi18n.T("Step1成功"))
		} else {
			stp.SetEnabled(false)
			mlog.IL(mi18n.T("Step1失敗"))
			return
		}
	})

	return stp, nil
}

// ------------------------------

type Step2TabPage struct {
	*mwidget.MTabPage
	mWindow  *mwidget.MWindow
	prevStep *Step1TabPage
	Items    *Step2Items
}

// ------------------------------

type Step2Items struct {
	stepItems
	MaterialListBox *MaterialListBox
	okButton        *walk.PushButton
}

func (si *Step2Items) SetEnabled(enabled bool) {
	si.stepItems.SetEnabled(enabled)
	si.MaterialListBox.SetEnabled(enabled)
	si.okButton.SetEnabled(enabled)
}

func (si *Step2Items) Dispose() {
	si.stepItems.Dispose()
	si.MaterialListBox.Dispose()
	si.okButton.Dispose()
}

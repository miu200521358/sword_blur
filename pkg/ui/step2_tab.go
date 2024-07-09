package ui

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/mlib_go/pkg/mwidget"
	"github.com/miu200521358/mlib_go/pkg/vmd"
	"github.com/miu200521358/sword_blur/pkg/model"
	"github.com/miu200521358/sword_blur/pkg/usecase"
	"github.com/miu200521358/walk/pkg/walk"
)

func NewStep2TabPage(
	mWindow *mwidget.MWindow, step1Page *Step1TabPage, blurModel *model.BlurModel,
) (*Step2TabPage, error) {
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
	step1Page.nextStep = stp

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

	// Step1. OKボタンクリック時
	step1Page.Items.okButton.Clicked().Attach(func() {
		step1Page.funcOkButton(blurModel)
	})

	return stp, nil
}

// ------------------------------

type Step2TabPage struct {
	*mwidget.MTabPage
	mWindow  *mwidget.MWindow
	prevStep *Step1TabPage
	nextStep *Step3TabPage
	Items    *Step2Items
}

// Step2. OKボタンクリック時
func (step2Page *Step2TabPage) funcOkButton(blurModel *model.BlurModel) {
	if len(step2Page.Items.MaterialListBox.SelectedIndexes()) == 0 {
		step2Page.nextStep.SetEnabled(false)
		mlog.ILT(mi18n.T("設定失敗"), mi18n.T("Step2失敗"))
		return
	} else {
		blurModel.BlurMaterialIndexes = step2Page.Items.MaterialListBox.SelectedIndexes()
		blurModel.OutputModel = nil
		blurModel.OutputMotion = nil

		step2Page.nextStep.SetEnabled(true)
		step2Page.mWindow.SetCheckWireDebugView(true)
		step2Page.mWindow.SetCheckSelectedVertexDebugView(true)
		step2Page.mWindow.TabWidget.SetCurrentIndex(2)                                                  // Step3へ移動
		step2Page.mWindow.GetMainGlWindow().SetFuncWorldPos(step2Page.nextStep.FuncWorldPos(blurModel)) // 頂点選択時のターゲットfunction変更

		go func() {
			step2Page.mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]*mwidget.ModelSet{0: {NextSelectedVertexIndexes: []int{}}}
		}()

		mlog.IL(mi18n.T("Step2成功"))
	}
}

// Step2. 材質選択
func (stp *Step2TabPage) funcMaterialListBoxChanged(blurModel *model.BlurModel) func(indexes []int) {
	return func(indexes []int) {
		if !stp.Items.MaterialListBox.Enabled() {
			return
		}

		invisibleMaterialIndexes := make([]int, 0)
		for i := range blurModel.Model.Materials.Len() {
			material := blurModel.Model.Materials.Get(i)
			mf := vmd.NewMorphFrame(0)
			if slices.Contains(indexes, i) {
				mf.Ratio = 0.0
			} else {
				mf.Ratio = 1.0
				invisibleMaterialIndexes = append(invisibleMaterialIndexes, i)
			}
			blurModel.Motion.AppendMorphFrame(usecase.GetVisibleMorphName(material), mf)
		}

		go func() {
			stp.mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]*mwidget.ModelSet{0: {NextMotion: blurModel.Motion, NextInvisibleMaterialIndexes: invisibleMaterialIndexes}}
		}()
	}
}

func (stp *Step2TabPage) SetEnabled(enabled bool) {
	stp.MTabPage.SetEnabled(enabled)
	stp.Items.SetEnabled(enabled)
	stp.nextStep.SetEnabled(enabled)
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

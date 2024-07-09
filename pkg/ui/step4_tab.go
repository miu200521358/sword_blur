package ui

import (
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/mlib_go/pkg/mwidget"
	"github.com/miu200521358/mlib_go/pkg/vmd"
	"github.com/miu200521358/sword_blur/pkg/model"
	"github.com/miu200521358/walk/pkg/walk"
)

func NewStep4TabPage(
	mWindow *mwidget.MWindow, step3Page *Step3TabPage, blurModel *model.BlurModel,
) (*Step4TabPage, error) {
	page, err := mwidget.NewMTabPage(mWindow, mWindow.TabWidget, "Step. 4")
	if err != nil {
		return nil, err
	}
	page.SetLayout(walk.NewVBoxLayout())

	stp := &Step4TabPage{
		MTabPage: page,
		mWindow:  mWindow,
		prevStep: step3Page,
		Items:    &Step4Items{},
	}
	step3Page.nextStep = stp

	// Step4. 切っ先選択

	stp.Items.composite, err = walk.NewComposite(stp)
	if err != nil {
		return nil, err
	}
	stp.Items.composite.SetLayout(walk.NewVBoxLayout())

	stp.Items.label, err = walk.NewTextLabel(stp.Items.composite)
	if err != nil {
		return nil, err
	}
	stp.Items.label.SetText(mi18n.T("Step4Label"))

	walk.NewVSeparator(stp.Items.composite)

	// 頂点選択リストボックス
	stp.Items.VertexListBox, err = NewVertexListBox(stp.Items.composite)
	if err != nil {
		return nil, err
	}

	// OKボタン
	stp.Items.okButton, err = walk.NewPushButton(stp.Items.composite)
	if err != nil {
		return nil, err
	}
	stp.Items.okButton.SetText(mi18n.T("次へ進む"))

	// Step3. OKボタンクリック時
	step3Page.Items.okButton.Clicked().Attach(func() {
		step3Page.funcOkButton(blurModel)
	})

	return stp, nil
}

// ------------------------------

type Step4TabPage struct {
	*mwidget.MTabPage
	mWindow  *mwidget.MWindow
	prevStep *Step3TabPage
	nextStep *Step5TabPage
	Items    *Step4Items
}

// Step4. OKボタンクリック時
func (step4Page *Step4TabPage) funcOkButton(blurModel *model.BlurModel) {
	if len(step4Page.Items.VertexListBox.GetItemValues()) == 0 {
		step4Page.nextStep.SetEnabled(false)
		mlog.ILT(mi18n.T("設定失敗"), mi18n.T("Step4失敗"))
		return
	} else {
		blurModel.EdgeTailVertexIndexes = step4Page.Items.VertexListBox.GetItemValues()
		blurModel.OutputModel = nil
		blurModel.OutputMotion = nil

		step4Page.nextStep.SetEnabled(true)
		step4Page.nextStep.Items.retryButton.SetEnabled(false)
		step4Page.nextStep.Items.saveButton.SetEnabled(false)

		step4Page.mWindow.SetCheckWireDebugView(true)
		step4Page.mWindow.SetCheckSelectedVertexDebugView(true)
		step4Page.mWindow.TabWidget.SetCurrentIndex(4) // Step5へ移動

		mlog.IL(mi18n.T("Step4成功"))
	}
}

// Step3. マウスカーソル位置の頂点選択
func (stp Step4TabPage) FuncWorldPos(
	blurModel *model.BlurModel,
) func(worldPos *mmath.MVec3, vmdDeltas []*vmd.VmdDeltas, viewMat *mmath.MMat4) {
	return func(worldPos *mmath.MVec3, vmdDeltas []*vmd.VmdDeltas, viewMat *mmath.MMat4) {
		if !stp.Enabled() {
			return
		}
		if ok, _ := mutils.ExistsFile(blurModel.Model.GetPath()); ok {
			// 表示されている材質からのみ直近頂点を選ぶ
			nearestVertexIndexes := vmdDeltas[0].Vertices.GetNearestVertexIndexes(
				worldPos, blurModel.BlurMaterialIndexes)
			stp.Items.VertexListBox.SetItem(nearestVertexIndexes)

			go func() {
				stp.mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]*mwidget.ModelSet{0: {NextSelectedVertexIndexes: stp.Items.VertexListBox.GetItemValues()}}
			}()
		}
	}
}

func (stp *Step4TabPage) SetEnabled(enabled bool) {
	stp.MTabPage.SetEnabled(enabled)
	stp.Items.SetEnabled(enabled)
	stp.nextStep.SetEnabled(enabled)
}

// ------------------------------

type Step4Items struct {
	stepItems
	VertexListBox *VertexListBox
	okButton      *walk.PushButton
}

func (si *Step4Items) SetEnabled(enabled bool) {
	si.stepItems.SetEnabled(enabled)
	si.VertexListBox.SetEnabled(enabled)
	si.okButton.SetEnabled(enabled)
}

func (si *Step4Items) Dispose() {
	si.stepItems.Dispose()
	si.VertexListBox.Dispose()
	si.okButton.Dispose()
}

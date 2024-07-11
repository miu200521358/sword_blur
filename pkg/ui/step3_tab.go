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

func NewStep3TabPage(
	mWindow *mwidget.MWindow, step2Page *Step2TabPage, blurModel *model.BlurModel,
) (*Step3TabPage, error) {
	page, err := mwidget.NewMTabPage(mWindow, mWindow.TabWidget, "Step. 3")
	if err != nil {
		return nil, err
	}
	page.SetLayout(walk.NewVBoxLayout())

	stp := &Step3TabPage{
		MTabPage: page,
		mWindow:  mWindow,
		prevStep: step2Page,
		Items:    &Step3Items{},
	}
	step2Page.nextStep = stp

	// Step3. 峰手元選択

	stp.Items.composite, err = walk.NewComposite(stp)
	if err != nil {
		return nil, err
	}
	stp.Items.composite.SetLayout(walk.NewVBoxLayout())

	stp.Items.label, err = walk.NewTextLabel(stp.Items.composite)
	if err != nil {
		return nil, err
	}
	stp.Items.label.SetText(mi18n.T("Step3Label"))

	walk.NewVSeparator(stp.Items.composite)

	// クリアボタン
	stp.Items.clearButton, err = walk.NewPushButton(stp.Items.composite)
	if err != nil {
		return nil, err
	}
	stp.Items.clearButton.SetText(mi18n.T("クリア"))

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

	// Step2. OKボタンクリック時
	step2Page.Items.okButton.Clicked().Attach(func() {
		step2Page.funcOkButton(blurModel)
	})

	// Step3. clearボタンクリック時
	stp.Items.clearButton.Clicked().Attach(func() {
		stp.funcClearButton()
	})

	return stp, nil
}

// ------------------------------

type Step3TabPage struct {
	*mwidget.MTabPage
	mWindow  *mwidget.MWindow
	prevStep *Step2TabPage
	nextStep *Step4TabPage
	Items    *Step3Items
}

// Step3. OKボタンクリック時
func (step3Page *Step3TabPage) funcOkButton(blurModel *model.BlurModel) {
	if len(step3Page.Items.VertexListBox.GetItemValues()) == 0 {
		step3Page.nextStep.SetEnabled(false)
		mlog.ILT(mi18n.T("設定失敗"), mi18n.T("Step3失敗"))
		return
	} else {
		blurModel.BackRootVertexIndexes = step3Page.Items.VertexListBox.GetItemValues()
		blurModel.OutputModel = nil
		blurModel.OutputMotion = nil

		step3Page.nextStep.SetEnabled(true)
		step3Page.mWindow.SetCheckWireDebugView(true)
		step3Page.mWindow.SetCheckSelectedVertexDebugView(true)
		step3Page.mWindow.TabWidget.SetCurrentIndex(3)                                                  // Step4へ移動
		step3Page.mWindow.GetMainGlWindow().SetFuncWorldPos(step3Page.nextStep.FuncWorldPos(blurModel)) // 頂点選択時のターゲットfunction変更

		go func() {
			step3Page.mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]*mwidget.ModelSet{0: {NextSelectedVertexIndexes: []int{}}}
		}()

		mlog.IL(mi18n.T("Step3成功"))
	}
}

func (stp *Step3TabPage) funcClearButton() {
	stp.Items.VertexListBox.Clear()

	go func() {
		stp.mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]*mwidget.ModelSet{0: {NextSelectedVertexIndexes: []int{}}}
	}()
}

// Step3. マウスカーソル位置の頂点選択
func (stp Step3TabPage) FuncWorldPos(
	blurModel *model.BlurModel,
) func(prevXprevYFrontPos, prevXprevYBackPos, prevXnowYFrontPos, prevXnowYBackPos,
	nowXprevYFrontPos, nowXprevYBackPos, nowXnowYFrontPos, nowXnowYBackPos *mmath.MVec3, vmdDeltas []*vmd.VmdDeltas) {
	return func(prevXprevYFrontPos, prevXprevYBackPos, prevXnowYFrontPos, prevXnowYBackPos,
		nowXprevYFrontPos, nowXprevYBackPos, nowXnowYFrontPos, nowXnowYBackPos *mmath.MVec3, vmdDeltas []*vmd.VmdDeltas) {
		if !stp.Enabled() {
			return
		}
		if ok, _ := mutils.ExistsFile(blurModel.Model.GetPath()); ok {
			var nearestVertexIndexes [][]int
			// 直近頂点を取得
			if prevXnowYFrontPos == nil {
				nearestVertexIndexes = vmdDeltas[0].Vertices.FindNearestVertexIndexes(
					prevXprevYFrontPos, blurModel.BlurMaterialIndexes)
			} else {
				nearestVertexIndexes = vmdDeltas[0].Vertices.FindVerticesInBox(prevXprevYFrontPos, prevXprevYBackPos,
					prevXnowYFrontPos, prevXnowYBackPos, nowXprevYFrontPos, nowXprevYBackPos, nowXnowYFrontPos,
					nowXnowYBackPos, blurModel.BlurMaterialIndexes)
			}

			for _, vertexIndexes := range nearestVertexIndexes {
				// 表示されている材質からのみ直近頂点を選ぶ
				stp.Items.VertexListBox.SetItem(vertexIndexes)
			}

			go func() {
				stp.mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]*mwidget.ModelSet{0: {NextSelectedVertexIndexes: stp.Items.VertexListBox.GetItemValues()}}
			}()
		}
	}
}

func (stp *Step3TabPage) SetEnabled(enabled bool) {
	stp.MTabPage.SetEnabled(enabled)
	stp.Items.SetEnabled(enabled)
	stp.nextStep.SetEnabled(enabled)
}

// ------------------------------

type Step3Items struct {
	stepItems
	VertexListBox *VertexListBox
	okButton      *walk.PushButton
	clearButton   *walk.PushButton
}

func (si *Step3Items) SetEnabled(enabled bool) {
	si.stepItems.SetEnabled(enabled)
	si.VertexListBox.SetEnabled(enabled)
	si.okButton.SetEnabled(enabled)
	si.clearButton.SetEnabled(enabled)
}

func (si *Step3Items) Dispose() {
	si.stepItems.Dispose()
	si.VertexListBox.Dispose()
	si.okButton.Dispose()
	si.clearButton.Dispose()
}

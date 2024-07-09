package ui

import (
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/mlib_go/pkg/mwidget"
	"github.com/miu200521358/mlib_go/pkg/vmd"
	"github.com/miu200521358/walk/pkg/walk"
)

func NewStep4TabPage(mWindow *mwidget.MWindow, step3Page *Step3TabPage) (*Step4TabPage, error) {
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

	stp.SetEnabled(false)

	// Step3. OKボタンクリック時
	step3Page.Items.okButton.Clicked().Attach(func() {
		if len(step3Page.Items.VertexListBox.GetItemValues()) == 0 {
			stp.SetEnabled(false)
			mlog.IL(mi18n.T("Step3失敗"))
			return
		} else {
			stp.SetEnabled(true)
			stp.mWindow.SetCheckWireDebugView(true)
			stp.mWindow.SetCheckSelectedVertexDebugView(true)
			stp.mWindow.TabWidget.SetCurrentIndex(3)                              // Step4へ移動
			stp.mWindow.GetMainGlWindow().SetFuncWorldPos(stp.Items.FuncWorldPos) // 頂点選択時のターゲットfunction変更

			go func() {
				mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]*mwidget.ModelSet{0: {NextSelectedVertexIndexes: []int{}}}
			}()

			mlog.IL(mi18n.T("Step3成功"))
		}
	})

	stp.Items.FuncWorldPos = func(worldPos *mmath.MVec3, vmdDeltas []*vmd.VmdDeltas, viewMat *mmath.MMat4) {
		if step3Page.prevStep.prevStep.Items.OriginalPmxPicker.Exists() && stp.Enabled() {
			// 表示されている材質からのみ直近頂点を選ぶ
			nearestVertexIndexes := vmdDeltas[0].Vertices.GetNearestVertexIndexes(
				worldPos, step3Page.prevStep.Items.MaterialListBox.SelectedIndexes())
			stp.Items.VertexListBox.SetItem(nearestVertexIndexes)

			go func() {
				mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]*mwidget.ModelSet{0: {NextSelectedVertexIndexes: stp.Items.VertexListBox.GetItemValues()}}
			}()
		}
	}

	return stp, nil
}

// ------------------------------

type Step4TabPage struct {
	*mwidget.MTabPage
	mWindow  *mwidget.MWindow
	prevStep *Step3TabPage
	Items    *Step4Items
}

// ------------------------------

type Step4Items struct {
	stepItems
	VertexListBox *VertexListBox
	okButton      *walk.PushButton
	FuncWorldPos  func(worldPos *mmath.MVec3, vmdDeltas []*vmd.VmdDeltas, viewMat *mmath.MMat4)
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

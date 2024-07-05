package ui

import (
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/mlib_go/pkg/mwidget"
	"github.com/miu200521358/walk/pkg/walk"
)

func NewStep3TabPage(mWindow *mwidget.MWindow, step2Page *Step2TabPage) (*Step3TabPage, error) {
	page, err := mwidget.NewMTabPage(mWindow, mWindow.TabWidget, "Step. 3")
	if err != nil {
		return nil, err
	}
	page.SetLayout(walk.NewVBoxLayout())

	stp := &Step3TabPage{
		MTabPage: page,
		mWindow:  mWindow,
		Items:    &Step3Items{},
	}

	// Step3. 峰選択

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

	// 材質選択リストボックス
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

	// Step2. OKボタンクリック時
	step2Page.Items.okButton.Clicked().Attach(func() {
		if len(step2Page.Items.MaterialListBox.SelectedIndexes()) == 0 {
			mlog.IL(mi18n.T("Step2材質設定失敗"))
			return
		} else {
			stp.SetEnabled(true)
			mlog.IL(mi18n.T("Step2材質設定完了"))
		}
	})

	return stp, nil
}

// ------------------------------

type Step3TabPage struct {
	*mwidget.MTabPage
	mWindow *mwidget.MWindow
	Items   *Step3Items
}

// ------------------------------

type Step3Items struct {
	stepItems
	VertexListBox *VertexListBox
	okButton      *walk.PushButton
}

func (si *Step3Items) SetEnabled(enabled bool) {
	si.stepItems.SetEnabled(enabled)
	si.VertexListBox.SetEnabled(enabled)
	si.okButton.SetEnabled(enabled)
}

func (si *Step3Items) Dispose() {
	si.stepItems.Dispose()
	si.VertexListBox.Dispose()
	si.okButton.Dispose()
}

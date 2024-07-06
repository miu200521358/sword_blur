package ui

import (
	"embed"

	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mwidget"
	"github.com/miu200521358/walk/pkg/walk"
)

func NewStep1TabPage(mWindow *mwidget.MWindow, resourceFiles embed.FS) (*Step1TabPage, error) {
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
	stp.Items.MotionPlayer, err = mwidget.NewMotionPlayer(stp.Items.composite, mWindow, resourceFiles)
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

	return stp, nil
}

// ------------------------------

type Step1TabPage struct {
	*mwidget.MTabPage
	mWindow *mwidget.MWindow
	Items   *Step1Items
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

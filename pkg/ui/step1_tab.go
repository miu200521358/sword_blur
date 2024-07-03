package ui

import (
	"embed"

	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mwidget"
	"github.com/miu200521358/walk/pkg/walk"
)

type Step1Page struct {
	*mwidget.MTabPage
	mWindow           *mwidget.MWindow
	MotionPlayer      *mwidget.MotionPlayer
	OriginalPmxPicker *mwidget.FilePicker
	OriginalVmdPicker *mwidget.FilePicker
	OutputPmxPicker   *mwidget.FilePicker
}

func NewStep1TabPage(mWindow *mwidget.MWindow, resourceFiles embed.FS) (*Step1Page, error) {
	page, err := mwidget.NewMTabPage(mWindow, mWindow.TabWidget, "Step1")
	if err != nil {
		return nil, err
	}
	page.SetLayout(walk.NewVBoxLayout())

	sp := &Step1Page{
		MTabPage: page,
		mWindow:  mWindow,
	}

	headerComposite, err := walk.NewComposite(sp)
	if err != nil {
		return nil, err
	}
	headerComposite.SetLayout(walk.NewVBoxLayout())

	sp.OriginalPmxPicker, err = (mwidget.NewPmxReadFilePicker(
		mWindow,
		headerComposite,
		"org_pmx",
		mi18n.T("ブレ生成対象モデル(Pmx)"),
		mi18n.T("ブレ生成対象モデルPmxファイルを選択してください"),
		mi18n.T("ブレ生成対象モデルの使い方"),
		func(path string) {}))
	if err != nil {
		return nil, err
	}

	sp.OriginalVmdPicker, err = (mwidget.NewVmdVpdReadFilePicker(
		mWindow,
		headerComposite,
		"vmd",
		mi18n.T("確認用モーション(Vmd/Vpd)"),
		mi18n.T("確認用モーション(Vmd/Vpd)ファイルを選択してください"),
		mi18n.T("確認用モーションの使い方"),
		func(path string) {}))
	if err != nil {
		return nil, err
	}

	sp.OutputPmxPicker, err = (mwidget.NewPmxSaveFilePicker(
		mWindow,
		headerComposite,
		mi18n.T("出力モデル(Pmx)"),
		mi18n.T("出力モデル(Pmx)ファイルパスを指定してください"),
		mi18n.T("出力モデルの使い方"),
		func(path string) {}))
	if err != nil {
		return nil, err
	}

	// プレイヤーBox
	sp.MotionPlayer, err = mwidget.NewMotionPlayer(headerComposite, mWindow, resourceFiles)
	if err != nil {
		return nil, err
	}

	sp.MotionPlayer.OnPlay = func(isPlaying bool) error {
		// 入力欄は全部再生中は無効化
		sp.OriginalPmxPicker.SetEnabled(!isPlaying)
		sp.OriginalVmdPicker.SetEnabled(!isPlaying)
		sp.OutputPmxPicker.SetEnabled(!isPlaying)
		sp.MotionPlayer.SetEnabled(!isPlaying)
		sp.MotionPlayer.PlayButton.SetEnabled(true)
		for _, glWindow := range mWindow.GlWindows {
			glWindow.Play(isPlaying)
		}

		return nil
	}

	sp.OriginalPmxPicker.PathLineEdit.SetFocus()

	return sp, nil
}

func (sp *Step1Page) Dispose() {
	sp.MotionPlayer.Dispose()
	sp.OriginalPmxPicker.Dispose()
	sp.OriginalVmdPicker.Dispose()
	sp.OutputPmxPicker.Dispose()
}

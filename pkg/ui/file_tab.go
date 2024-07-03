package ui

import (
	"embed"
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

func NewFileTabPage(mWindow *mwidget.MWindow, resourceFiles embed.FS) (*FileTabPage, error) {
	page, err := mwidget.NewMTabPage(mWindow, mWindow.TabWidget, mi18n.T("ファイル"))
	if err != nil {
		return nil, err
	}
	page.SetLayout(walk.NewVBoxLayout())

	ftp := &FileTabPage{
		MTabPage:     page,
		mWindow:      mWindow,
		currentIndex: -1,
		Step1:        &Step1Items{},
		Step2:        &Step2Items{},
	}

	ftp.Step1.step = 0
	ftp.Step2.step = 1

	// ナビゲーション用ツールバー
	ftp.navToolBar, err = walk.NewToolBarWithOrientationAndButtonStyle(
		page, walk.Horizontal, walk.ToolBarButtonTextOnly)
	if err != nil {
		return nil, err
	}

	ftp.Step1.composite, err = walk.NewComposite(ftp)
	if err != nil {
		return nil, err
	}
	ftp.Step1.composite.SetLayout(walk.NewVBoxLayout())

	// Step1. ファイル選択文言
	ftp.Step1.label, err = walk.NewTextLabel(ftp.Step1.composite)
	if err != nil {
		return nil, err
	}
	ftp.Step1.label.SetText(mi18n.T("Step1Label"))

	walk.NewVSeparator(ftp.Step1.composite)

	ftp.Step1.OriginalPmxPicker, err = (mwidget.NewPmxReadFilePicker(
		mWindow,
		ftp.Step1.composite,
		"org_pmx",
		mi18n.T("ブレ生成対象モデル(Pmx)"),
		mi18n.T("ブレ生成対象モデルPmxファイルを選択してください"),
		mi18n.T("ブレ生成対象モデルの使い方"),
		nil))
	if err != nil {
		return nil, err
	}

	ftp.Step1.OriginalVmdPicker, err = (mwidget.NewVmdVpdReadFilePicker(
		mWindow,
		ftp.Step1.composite,
		"vmd",
		mi18n.T("確認用モーション(Vmd/Vpd)"),
		mi18n.T("確認用モーション(Vmd/Vpd)ファイルを選択してください"),
		mi18n.T("確認用モーションの使い方"),
		func(path string) {}))
	if err != nil {
		return nil, err
	}

	ftp.Step1.OutputPmxPicker, err = (mwidget.NewPmxSaveFilePicker(
		mWindow,
		ftp.Step1.composite,
		mi18n.T("出力モデル(Pmx)"),
		mi18n.T("出力モデル(Pmx)ファイルパスを指定してください"),
		mi18n.T("出力モデルの使い方"),
		func(path string) {}))
	if err != nil {
		return nil, err
	}

	// プレイヤーBox
	ftp.Step1.MotionPlayer, err = mwidget.NewMotionPlayer(ftp.Step1.composite, mWindow, resourceFiles)
	if err != nil {
		return nil, err
	}

	ftp.Step1.MotionPlayer.OnPlay = func(isPlaying bool) error {
		// 入力欄は全部再生中は無効化
		ftp.Step1.OriginalPmxPicker.SetEnabled(!isPlaying)
		ftp.Step1.OriginalVmdPicker.SetEnabled(!isPlaying)
		ftp.Step1.OutputPmxPicker.SetEnabled(!isPlaying)
		ftp.Step1.MotionPlayer.SetEnabled(!isPlaying)
		ftp.Step1.MotionPlayer.PlayButton.SetEnabled(true)
		for _, glWindow := range mWindow.GlWindows {
			glWindow.Play(isPlaying)
		}

		return nil
	}

	ftp.Step1.OriginalPmxPicker.OnPathChanged = func(path string) {
		func(path string) {
			isExist, err := mutils.ExistsFile(path)
			if !isExist || err != nil {
				ftp.Step1.OutputPmxPicker.PathLineEdit.SetText("")
				return
			}

			// とりあえず次を有効に
			ftp.Step1.OutputPmxPicker.SetEnabled(true)
			ftp.Step2.SetEnabled(true)

			// 出力パス設定
			dir, file := filepath.Split(path)
			ext := filepath.Ext(file)
			outputPath := filepath.Join(dir, file[:len(file)-len(ext)]+"_blur"+ext)
			ftp.Step1.OutputPmxPicker.PathLineEdit.SetText(outputPath)

			if ftp.Step1.OriginalPmxPicker.Exists() {
				data, err := ftp.Step1.OriginalPmxPicker.GetData()
				if err != nil {
					mlog.E(mi18n.T("Pmxファイル読み込みエラー"), err.Error())
					return
				}
				model := data.(*pmx.PmxModel)
				var motion *vmd.VmdMotion
				if ftp.Step1.OriginalVmdPicker.IsCached() {
					motion = ftp.Step1.OriginalVmdPicker.GetCache().(*vmd.VmdMotion)
				} else {
					motion = vmd.NewVmdMotion("")
				}

				ftp.Step1.MotionPlayer.SetEnabled(true)
				mWindow.GetMainGlWindow().SetFrame(0)
				ftp.Step1.MotionPlayer.SetValue(0)
				mWindow.GetMainGlWindow().Play(false)
				mWindow.GetMainGlWindow().ClearData()
				mWindow.GetMainGlWindow().AddData(model, motion)
				mWindow.GetMainGlWindow().Run()
			}
		}(path)
	}

	ftp.Step1.OriginalPmxPicker.PathLineEdit.SetFocus()
	ftp.Step1.OutputPmxPicker.SetEnabled(false)
	ftp.Step1.MotionPlayer.SetEnabled(false)

	ftp.addStep(ftp.Step1)

	// ------------------------------
	// Step2. 材質選択

	ftp.Step2.composite, err = walk.NewComposite(ftp)
	if err != nil {
		return nil, err
	}
	ftp.Step2.composite.SetLayout(walk.NewVBoxLayout())

	// Step2. ファイル選択文言
	ftp.Step2.label, err = walk.NewTextLabel(ftp.Step2.composite)
	if err != nil {
		return nil, err
	}
	ftp.Step2.label.SetText(mi18n.T("Step2Label"))
	walk.NewVSeparator(ftp.Step2.composite)

	// 材質選択リストボックス
	ftp.Step2.MaterialListBox, err = NewMaterialListBox(ftp.Step2.composite)
	if err != nil {
		return nil, err
	}

	// 最初は操作不可
	ftp.addStep(ftp.Step2)
	ftp.Step2.SetEnabled(false)

	// 最初は Step1 から
	ftp.setCurrentAction(ftp.Step1.Step())

	return ftp, nil
}

// ------------------------------

type FileTabPage struct {
	*mwidget.MTabPage
	mWindow                     *mwidget.MWindow
	navToolBar                  *walk.ToolBar
	currentIndex                int
	currentPageChangedPublisher walk.EventPublisher
	Step1                       *Step1Items
	Step2                       *Step2Items
}

func (ftp *FileTabPage) setCurrentAction(index int) error {
	ftp.SetFocus()

	for i := range 2 {
		ftp.navToolBar.Actions().At(i).SetChecked(false)
	}
	ftp.currentIndex = index
	ftp.navToolBar.Actions().At(index).SetChecked(true)
	ftp.currentPageChangedPublisher.Publish()

	// // 選択したステップのみ表示
	// ftp.Step1.SetVisible(index == 0)
	// ftp.Step2.SetVisible(index == 1)

	return nil
}

func (ftp *FileTabPage) newPageAction(step int) (*walk.Action, error) {
	action := walk.NewAction()
	action.SetCheckable(true)
	action.SetExclusive(true)
	action.SetText(fmt.Sprintf("Step. %d", step+1))

	action.Triggered().Attach(func() {
		ftp.setCurrentAction(step)
	})

	return action, nil
}

func (ftp *FileTabPage) addStep(item iStepItems) error {
	action, err := ftp.newPageAction(item.Step())
	if err != nil {
		return err
	}
	ftp.navToolBar.Actions().Add(action)

	return nil
}

func (sp *FileTabPage) Dispose() {
	sp.Step1.Dispose()
	sp.Step2.Dispose()
}

// ------------------------------

type iStepItems interface {
	SetVisible(bool)
	SetEnabled(bool)
	Dispose()
	Step() int
}

type stepItems struct {
	step      int
	label     *walk.TextLabel
	composite *walk.Composite
}

func (si *stepItems) Step() int {
	return si.step
}

func (si *stepItems) SetEnabled(enabled bool) {
	si.label.SetEnabled(enabled)
}

func (si *stepItems) SetVisible(visible bool) {
	si.label.SetVisible(visible)
	si.composite.SetVisible(visible)
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
}

func (si *Step1Items) SetVisible(visible bool) {
	si.stepItems.SetVisible(visible)
	si.MotionPlayer.SetVisible(visible)
	si.OriginalPmxPicker.SetVisible(visible)
	si.OriginalVmdPicker.SetVisible(visible)
	si.OutputPmxPicker.SetVisible(visible)
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

type Step2Items struct {
	stepItems
	MaterialListBox *MaterialListBox
}

func (si *Step2Items) SetEnabled(enabled bool) {
	si.stepItems.SetEnabled(enabled)
	si.MaterialListBox.SetEnabled(enabled)
}

func (si *Step2Items) SetVisible(visible bool) {
	si.stepItems.SetVisible(visible)
	si.MaterialListBox.SetVisible(visible)
}

func (si *Step2Items) Dispose() {
	si.stepItems.Dispose()
	si.MaterialListBox.Dispose()
}

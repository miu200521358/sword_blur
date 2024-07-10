package ui

import (
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/mlib_go/pkg/mwidget"
	"github.com/miu200521358/mlib_go/pkg/pmx"
	"github.com/miu200521358/mlib_go/pkg/vmd"
	"github.com/miu200521358/sword_blur/pkg/model"
	"github.com/miu200521358/sword_blur/pkg/usecase"
	"github.com/miu200521358/walk/pkg/walk"
)

func NewStep5TabPage(
	mWindow *mwidget.MWindow, step4Page *Step4TabPage, blurModel *model.BlurModel,
) (*Step5TabPage, error) {
	page, err := mwidget.NewMTabPage(mWindow, mWindow.TabWidget, "Step. 5")
	if err != nil {
		return nil, err
	}
	page.SetLayout(walk.NewVBoxLayout())

	stp := &Step5TabPage{
		MTabPage: page,
		mWindow:  mWindow,
		prevStep: step4Page,
		Items:    &Step5Items{},
	}
	step4Page.nextStep = stp

	// Step5. 刃選択

	stp.Items.composite, err = walk.NewComposite(stp)
	if err != nil {
		return nil, err
	}
	stp.Items.composite.SetLayout(walk.NewVBoxLayout())

	stp.Items.label, err = walk.NewTextLabel(stp.Items.composite)
	if err != nil {
		return nil, err
	}
	stp.Items.label.SetText(mi18n.T("Step5Label"))

	walk.NewVSeparator(stp.Items.composite)

	// 材質選択リストボックス
	stp.Items.VertexListBox, err = NewVertexListBox(stp.Items.composite)
	if err != nil {
		return nil, err
	}

	// プレビューボタン
	stp.Items.previewButton, err = walk.NewPushButton(stp.Items.composite)
	if err != nil {
		return nil, err
	}
	stp.Items.previewButton.SetText(mi18n.T("プレビュー"))

	// リトライボタン
	stp.Items.retryButton, err = walk.NewPushButton(stp.Items.composite)
	if err != nil {
		return nil, err
	}
	stp.Items.retryButton.SetText(mi18n.T("リトライ"))

	// 保存ボタン
	stp.Items.saveButton, err = walk.NewPushButton(stp.Items.composite)
	if err != nil {
		return nil, err
	}
	stp.Items.saveButton.SetText(mi18n.T("保存"))

	// Step5. プレビューボタンクリック時
	stp.Items.previewButton.Clicked().Attach(func() {
		stp.funcPreviewButton(blurModel)
	})

	// Step5. リトライボタンクリック時
	stp.Items.retryButton.Clicked().Attach(func() {
		stp.funcRetryButton(blurModel)
	})

	// Step5. 保存ボタンクリック時
	stp.Items.saveButton.Clicked().Attach(func() {
		stp.funcSaveButton(blurModel)
	})

	// Step4. OKボタンクリック時
	step4Page.Items.okButton.Clicked().Attach(func() {
		step4Page.funcOkButton(blurModel)
	})

	return stp, nil
}

// ------------------------------

type Step5TabPage struct {
	*mwidget.MTabPage
	mWindow  *mwidget.MWindow
	prevStep *Step4TabPage
	Items    *Step5Items
}

// Step3. マウスカーソル位置の頂点選択
func (stp Step5TabPage) FuncWorldPos(
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

// Step5. プレビューボタンクリック時
func (stp *Step5TabPage) funcPreviewButton(blurModel *model.BlurModel) {
	if len(stp.Items.VertexListBox.GetItemValues()) == 0 {
		mlog.ILT(mi18n.T("設定失敗"), mi18n.T("Step5失敗"))
		return
	} else {
		blurModel.EdgeVertexIndexes = stp.Items.VertexListBox.GetItemValues()
		blurModel.OutputModel = stp.prevStep.prevStep.prevStep.prevStep.Items.OriginalPmxPicker.GetDataForce().(*pmx.PmxModel)
		blurModel.OutputMotion = nil

		err := usecase.Preview(blurModel)

		if err != nil {
			mlog.ET(mi18n.T("生成失敗"), mi18n.T("生成失敗メッセージ", map[string]interface{}{"Error": err.Error()}))
		} else {
			mlog.IT(mi18n.T("生成成功"), mi18n.T("生成成功メッセージ"))

			nowMaxFrame := stp.prevStep.prevStep.prevStep.prevStep.Items.MotionPlayer.FrameEdit.MaxValue()
			if blurModel.OutputMotion.BoneFrames.GetMaxFrame() > int(nowMaxFrame) {
				stp.prevStep.prevStep.prevStep.prevStep.Items.MotionPlayer.SetRange(0, blurModel.OutputMotion.BoneFrames.GetMaxFrame()+1)
			}
			stp.Items.saveButton.SetEnabled(true)
			stp.Items.retryButton.SetEnabled(true)
			stp.mWindow.SetCheckWireDebugView(false)
			stp.mWindow.SetCheckSelectedVertexDebugView(false)
			stp.prevStep.prevStep.prevStep.prevStep.Items.MotionPlayer.SetValue(0)
			stp.prevStep.prevStep.prevStep.prevStep.Items.MotionPlayer.Play(true)

			go func() {
				stp.mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]*mwidget.ModelSet{1: {NextModel: blurModel.OutputModel, NextMotion: blurModel.OutputMotion}}
				stp.mWindow.GetMainGlWindow().IsPlayingChannel <- true
			}()
		}
	}
}

// Step5. リトライボタンクリック時
func (stp *Step5TabPage) funcRetryButton(blurModel *model.BlurModel) {
	stp.mWindow.SetCheckWireDebugView(true)
	stp.mWindow.SetCheckSelectedVertexDebugView(true)
	stp.Items.saveButton.SetEnabled(false)
	stp.Items.retryButton.SetEnabled(false)
	stp.prevStep.prevStep.prevStep.prevStep.Items.MotionPlayer.SetValue(0)
	stp.prevStep.prevStep.prevStep.prevStep.Items.MotionPlayer.SetRange(0, 1)
	stp.prevStep.prevStep.prevStep.prevStep.Items.MotionPlayer.Play(true)

	blurModel.OutputModel = nil
	blurModel.OutputMotion = nil

	go func() {
		stp.mWindow.GetMainGlWindow().RemoveModelSetIndexChannel <- 1
		stp.mWindow.GetMainGlWindow().IsPlayingChannel <- false
	}()
}

// Step5. 保存ボタンクリック時
func (stp *Step5TabPage) funcSaveButton(blurModel *model.BlurModel) {
	if blurModel.OutputModel == nil {
		mlog.IT(mi18n.T("出力モデルなし"), mi18n.T("出力モデルなしメッセージ"))
		return
	} else {
		if err := usecase.Save(blurModel); err != nil {
			mlog.ET(mi18n.T("出力失敗"), mi18n.T("出力失敗メッセージ", map[string]interface{}{"Error": err.Error()}))
		} else {
			mlog.IT(mi18n.T("出力成功"), mi18n.T("出力成功メッセージ", map[string]interface{}{"Path": blurModel.OutputModelPath}))
		}

		stp.mWindow.Beep()
	}
}

func (stp *Step5TabPage) SetEnabled(enabled bool) {
	stp.MTabPage.SetEnabled(enabled)
	stp.Items.SetEnabled(enabled)
}

// ------------------------------

type Step5Items struct {
	stepItems
	VertexListBox *VertexListBox
	previewButton *walk.PushButton
	retryButton   *walk.PushButton
	saveButton    *walk.PushButton
}

func (si *Step5Items) SetEnabled(enabled bool) {
	si.stepItems.SetEnabled(enabled)
	si.VertexListBox.SetEnabled(enabled)
	si.previewButton.SetEnabled(enabled)
	si.retryButton.SetEnabled(enabled)
	si.saveButton.SetEnabled(enabled)
}

func (si *Step5Items) Dispose() {
	si.stepItems.Dispose()
	si.VertexListBox.Dispose()
	si.previewButton.Dispose()
	si.retryButton.Dispose()
	si.saveButton.Dispose()
}

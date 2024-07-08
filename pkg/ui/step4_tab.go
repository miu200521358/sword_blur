package ui

import (
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/mlib_go/pkg/mwidget"
	"github.com/miu200521358/mlib_go/pkg/pmx"
	"github.com/miu200521358/mlib_go/pkg/vmd"
	"github.com/miu200521358/sword_blur/pkg/usecase"
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

	// Step4. 刃選択

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

	// 保存ボタン
	stp.Items.saveButton, err = walk.NewPushButton(stp.Items.composite)
	if err != nil {
		return nil, err
	}
	stp.Items.saveButton.SetText(mi18n.T("保存"))

	stp.SetEnabled(false)

	// Step3. OKボタンクリック時
	step3Page.Items.okButton.Clicked().Attach(func() {
		if len(step3Page.Items.VertexListBox.GetItemValues()) == 0 {
			stp.SetEnabled(false)
			mlog.IL(mi18n.T("Step3峰頂点設定失敗"))
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

			mlog.IL(mi18n.T("Step3峰頂点設定完了"))
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

	// Step4. プレビューボタンクリック時
	stp.Items.previewButton.Clicked().Attach(func() {
		if len(stp.Items.VertexListBox.GetItemValues()) == 0 {
			stp.SetEnabled(false)
			mlog.IL(mi18n.T("Step4刃頂点設定失敗"))
			return
		} else {
			outputModel, previewVmd, err := usecase.Preview(
				stp.prevStep.prevStep.prevStep.Items.OriginalPmxPicker.GetDataForce().(*pmx.PmxModel),
				stp.prevStep.prevStep.Items.MaterialListBox.SelectedIndexes(),
				stp.prevStep.Items.VertexListBox.GetItemValues(),
				stp.Items.VertexListBox.GetItemValues(),
			)

			if err != nil {
				mlog.ET(mi18n.T("生成失敗"), mi18n.T("生成失敗メッセージ", map[string]interface{}{"Error": err.Error()}))
			} else {
				stp.Items.saveButton.SetEnabled(true)
				mlog.IT(mi18n.T("生成成功"), mi18n.T("生成成功メッセージ"))

				nowMaxFrame := stp.prevStep.prevStep.prevStep.Items.MotionPlayer.FrameEdit.MaxValue()
				if previewVmd.BoneFrames.GetMaxFrame() > int(nowMaxFrame) {
					stp.prevStep.prevStep.prevStep.Items.MotionPlayer.SetRange(0, previewVmd.BoneFrames.GetMaxFrame()+1)
				}
				stp.mWindow.SetCheckWireDebugView(false)
				stp.mWindow.SetCheckSelectedVertexDebugView(false)
				stp.prevStep.prevStep.prevStep.Items.MotionPlayer.SetValue(0)
				stp.prevStep.prevStep.prevStep.Items.MotionPlayer.Play(true)
				stp.outputModel = outputModel

				go func() {
					mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]*mwidget.ModelSet{1: {NextModel: outputModel, NextMotion: previewVmd}}
					mWindow.GetMainGlWindow().IsPlayingChannel <- true
				}()
			}
		}
	})

	stp.Items.saveButton.Clicked().Attach(func() {
		if stp.outputModel == nil {
			mlog.IT(mi18n.T("出力モデルなし"), mi18n.T("出力モデルなしメッセージ"))
			return
		} else {
			outputPath := stp.prevStep.prevStep.prevStep.Items.OutputPmxPicker.GetPath()

			if err := usecase.Save(stp.outputModel, outputPath); err != nil {
				mlog.ET(mi18n.T("出力失敗"), mi18n.T("出力失敗メッセージ", map[string]interface{}{"Error": err.Error()}))
			} else {
				mlog.IT(mi18n.T("出力成功"), mi18n.T("出力成功メッセージ", map[string]interface{}{"Path": outputPath}))
			}
		}
	})

	return stp, nil
}

// ------------------------------

type Step4TabPage struct {
	*mwidget.MTabPage
	mWindow     *mwidget.MWindow
	prevStep    *Step3TabPage
	Items       *Step4Items
	outputModel *pmx.PmxModel
}

// ------------------------------

type Step4Items struct {
	stepItems
	VertexListBox *VertexListBox
	previewButton *walk.PushButton
	saveButton    *walk.PushButton
	FuncWorldPos  func(worldPos *mmath.MVec3, vmdDeltas []*vmd.VmdDeltas, viewMat *mmath.MMat4)
}

func (si *Step4Items) SetEnabled(enabled bool) {
	si.stepItems.SetEnabled(enabled)
	si.VertexListBox.SetEnabled(enabled)
	si.previewButton.SetEnabled(enabled)
	si.saveButton.SetEnabled(enabled)
}

func (si *Step4Items) Dispose() {
	si.stepItems.Dispose()
	si.VertexListBox.Dispose()
	si.previewButton.Dispose()
	si.saveButton.Dispose()
}

package ui

import (
	"github.com/miu200521358/mlib_go/pkg/interface/app"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/sword_blur/pkg/model"
	"github.com/miu200521358/walk/pkg/walk"
)

type ToolState struct {
	App                    *app.MApp
	ControlWindow          *controller.ControlWindow
	Player                 *widget.MotionPlayer
	BlurModel              *model.BlurModel
	Step1                  *widget.MTabPage
	OriginalPmxPicker      *widget.FilePicker
	OriginalVmdPicker      *widget.FilePicker
	OutputPmxPicker        *widget.FilePicker
	Step1OkButton          *walk.PushButton
	Step2                  *widget.MTabPage
	MaterialListBox        *MaterialListBox
	Step2OkButton          *walk.PushButton
	Step2ClearButton       *walk.PushButton
	Step3                  *widget.MTabPage
	RootVertexListBox      *VertexListBox
	Step3OkButton          *walk.PushButton
	Step3ClearButton       *walk.PushButton
	RootVertexSelectedFunc func([][][]int)
	Step4                  *widget.MTabPage
	TipVertexListBox       *VertexListBox
	Step4OkButton          *walk.PushButton
	Step4ClearButton       *walk.PushButton
	TipVertexSelectedFunc  func([][][]int)
	Step5                  *widget.MTabPage
	EdgeVertexListBox      *VertexListBox
	Step5PreviewButton     *walk.PushButton
	Step5RetryButton       *walk.PushButton
	Step5SaveButton        *walk.PushButton
	Step5ClearButton       *walk.PushButton
	EdgeVertexSelectedFunc func([][][]int)
}

func NewToolState(app *app.MApp, controlWindow *controller.ControlWindow) *ToolState {
	toolState := &ToolState{
		App:           app,
		ControlWindow: controlWindow,
		BlurModel:     model.NewBlurModel(),
	}

	newStep1Tab(controlWindow, toolState)
	newStep2Tab(controlWindow, toolState)
	newStep3Tab(controlWindow, toolState)
	newStep4Tab(controlWindow, toolState)
	newStep5Tab(controlWindow, toolState)

	player := widget.NewMotionPlayer(controlWindow.MainWindow, controlWindow)
	player.SetOnTriggerPlay(func(play bool) {
		toolState.SetEnabled(6)
	})
	controlWindow.SetPlayer(player)
	toolState.Player = player

	toolState.SetEnabled(1)

	// タブ切り替え時に選択頂点リストボックスの更新メソッドを切り替える
	toolState.ControlWindow.TabWidget.CurrentIndexChanged().Attach(func() {
		if toolState.ControlWindow.TabWidget.CurrentIndex() == 2 {
			// 根元選択頂点に切り替え
			toolState.ResetSelectedVertexes(true, false, false, nil)
			toolState.ControlWindow.SetFuncSetSelectedVertexes(toolState.RootVertexSelectedFunc)
		} else if toolState.ControlWindow.TabWidget.CurrentIndex() == 3 {
			// 切っ先選択頂点に切り替え
			toolState.ResetSelectedVertexes(false, true, false, nil)
			toolState.ControlWindow.SetFuncSetSelectedVertexes(toolState.TipVertexSelectedFunc)
		} else if toolState.ControlWindow.TabWidget.CurrentIndex() == 4 {
			// 刃選択頂点に切り替え
			toolState.ResetSelectedVertexes(false, false, true, nil)
			toolState.ControlWindow.SetFuncSetSelectedVertexes(toolState.EdgeVertexSelectedFunc)
		}
	})

	return toolState
}

func (toolState *ToolState) ResetSelectedVertexes(
	isSelectRoot, isSelectTip, isSelectEdge bool, addNoSelectedIndexes []int,
) {
	selectedIndexes := make([]int, 0)
	noSelectedIndexes := make([]int, 0)
	if isSelectRoot {
		selectedIndexes = append(selectedIndexes, toolState.BlurModel.RootVertexIndexes...)
	} else {
		noSelectedIndexes = append(noSelectedIndexes, toolState.BlurModel.RootVertexIndexes...)
	}
	if isSelectTip {
		selectedIndexes = append(selectedIndexes, toolState.BlurModel.TipVertexIndexes...)
	} else {
		noSelectedIndexes = append(noSelectedIndexes, toolState.BlurModel.TipVertexIndexes...)
	}
	if isSelectEdge {
		selectedIndexes = append(selectedIndexes, toolState.BlurModel.EdgeVertexIndexes...)
	} else {
		noSelectedIndexes = append(noSelectedIndexes, toolState.BlurModel.EdgeVertexIndexes...)
	}
	if addNoSelectedIndexes != nil {
		noSelectedIndexes = append(noSelectedIndexes, addNoSelectedIndexes...)
	}

	toolState.App.ControlToViewerChannel().SetSelectedVertexesChannel([][][]int{{selectedIndexes}})
	toolState.App.ControlToViewerChannel().SetNoSelectedVertexesChannel([][][]int{{noSelectedIndexes}})
}

func (toolState *ToolState) SetEnabled(step int) {
	// Step.1
	toolState.Step1.SetEnabled(step >= 1)
	toolState.OriginalPmxPicker.SetEnabled(step >= 1)
	toolState.OriginalVmdPicker.SetEnabled(step >= 1)
	toolState.OutputPmxPicker.SetEnabled(step >= 1)
	toolState.Step1OkButton.SetEnabled(step >= 1)
	// Step.2
	toolState.Step2.SetEnabled(step >= 2)
	toolState.MaterialListBox.SetEnabled(step >= 2)
	toolState.Step2OkButton.SetEnabled(step >= 2)
	toolState.Step2ClearButton.SetEnabled(step >= 2)
	// Step.3
	toolState.Step3.SetEnabled(step >= 3)
	toolState.RootVertexListBox.SetEnabled(step >= 3)
	toolState.Step3OkButton.SetEnabled(step >= 3)
	toolState.Step3ClearButton.SetEnabled(step >= 3)
	// Step.4
	toolState.Step4.SetEnabled(step >= 4)
	toolState.TipVertexListBox.SetEnabled(step >= 4)
	toolState.Step4OkButton.SetEnabled(step >= 4)
	toolState.Step4ClearButton.SetEnabled(step >= 4)
	// Step.5
	toolState.Step5.SetEnabled(step >= 5)
	toolState.EdgeVertexListBox.SetEnabled(step >= 5)
	toolState.Step5PreviewButton.SetEnabled(step >= 5)
	// Step.5 (プレビュー後)
	toolState.Step5RetryButton.SetEnabled(step >= 6)
	toolState.Step5SaveButton.SetEnabled(step >= 6)
}

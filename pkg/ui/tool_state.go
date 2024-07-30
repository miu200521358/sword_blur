package ui

import (
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/sword_blur/pkg/model"
	"github.com/miu200521358/walk/pkg/walk"
)

type ToolState struct {
	AppState          state.IAppState
	ControlWindow     *controller.ControlWindow
	BlurModel         *model.BlurModel
	Step1             *widget.MTabPage
	OriginalPmxPicker *widget.FilePicker
	OriginalVmdPicker *widget.FilePicker
	OutputPmxPicker   *widget.FilePicker
	Step1OkButton     *walk.PushButton
	Step2             *widget.MTabPage
	MaterialListBox   *MaterialListBox
	Step2OkButton     *walk.PushButton
	Step2ClearButton  *walk.PushButton
	Step3             *widget.MTabPage
	RootVertexListBox *VertexListBox
	Step3OkButton     *walk.PushButton
	Step3ClearButton  *walk.PushButton
}

func NewToolState(appState state.IAppState, controlWindow *controller.ControlWindow) *ToolState {
	toolState := &ToolState{
		AppState:      appState,
		ControlWindow: controlWindow,
		BlurModel:     model.NewBlurModel(),
	}

	newStep1Tab(controlWindow, toolState)
	newStep2Tab(controlWindow, toolState)
	newStep3Tab(controlWindow, toolState)

	toolState.SetEnabled(1)

	return toolState
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
}

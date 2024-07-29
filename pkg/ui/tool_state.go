package ui

import (
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/sword_blur/pkg/model"
	"github.com/miu200521358/walk/pkg/walk"
)

// ------------------------------

type ToolState struct {
	AppState          state.IAppState
	BlurModel         *model.BlurModel
	Step1             *widget.MTabPage
	OriginalPmxPicker *widget.FilePicker
	OriginalVmdPicker *widget.FilePicker
	OutputPmxPicker   *widget.FilePicker
	Step1OkButton     *walk.PushButton
}

func NewToolState(appState state.IAppState, controlWindow *controller.ControlWindow) *ToolState {
	toolState := &ToolState{
		AppState:  appState,
		BlurModel: model.NewBlurModel(),
	}

	newStep1Tab(controlWindow, toolState)

	return toolState
}

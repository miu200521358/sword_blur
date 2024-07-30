package ui

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/animation"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/sword_blur/pkg/usecase"
	"github.com/miu200521358/walk/pkg/walk"
)

func newStep3Tab(controlWindow *controller.ControlWindow, toolState *ToolState) {
	{
		toolState.Step3 = widget.NewMTabPage("Step.3")
		controlWindow.AddTabPage(toolState.Step3.TabPage)

		toolState.Step3.SetLayout(walk.NewVBoxLayout())

		{
			// Step3.文言
			label, err := walk.NewTextLabel(toolState.Step3)
			if err != nil {
				widget.RaiseError(err)
			}
			label.SetText(mi18n.T("Step3Label"))
		}

		walk.NewVSeparator(toolState.Step3)

		var err error
		{
			// クリアボタン
			toolState.Step3ClearButton, err = walk.NewPushButton(toolState.Step3)
			if err != nil {
				widget.RaiseError(err)
			}
			toolState.Step3ClearButton.SetText(mi18n.T("クリア"))
			toolState.Step3ClearButton.Clicked().Attach(toolState.onClickStep3Clear)
		}

		{
			// 頂点選択リストボックス
			toolState.RootVertexListBox, err = NewVertexListBox(toolState.Step3)
			if err != nil {
				widget.RaiseError(err)
			}
		}

		{
			// OKボタン
			toolState.Step3OkButton, err = walk.NewPushButton(toolState.Step3)
			if err != nil {
				widget.RaiseError(err)
			}
			toolState.Step3OkButton.SetText(mi18n.T("次へ進む"))
			toolState.Step3OkButton.Clicked().Attach(toolState.onClickStep3Ok)
		}
	}
}

func (toolState *ToolState) onClickStep3Clear() {
}

// Step3. 頂点選択時
func (toolState *ToolState) onChangeRootVertexListBox() func(indexes []int) {
	return func(indexes []int) {
		if !toolState.MaterialListBox.Enabled() {
			return
		}

		invisibleMaterialIndexes := make([]int, 0)
		for i := range toolState.BlurModel.Model.Materials.Len() {
			material := toolState.BlurModel.Model.Materials.Get(i)
			mf := vmd.NewMorphFrame(0)
			if slices.Contains(indexes, i) {
				mf.Ratio = 0.0
			} else {
				mf.Ratio = 1.0
				invisibleMaterialIndexes = append(invisibleMaterialIndexes, i)
			}
			toolState.BlurModel.Motion.AppendRegisteredMorphFrame(usecase.GetVisibleMorphName(material), mf)
		}

		// outputPath := mutils.CreateOutputPath(
		// 	strings.ReplaceAll(toolState.BlurModel.Model.Path(), ".pmx", ".vmd"), "mat_off")
		// repository.NewVmdRepository().Save(outputPath, toolState.BlurModel.Motion, true)

		animationState := animation.NewAnimationState(0, 0)
		animationState.SetMotion(toolState.BlurModel.Motion)
		animationState.SetInvisibleMaterialIndexes(invisibleMaterialIndexes)
		toolState.ControlWindow.SetAnimationState(animationState)
	}
}

func (toolState *ToolState) onClickStep3Ok() {
	if len(toolState.MaterialListBox.SelectedIndexes()) == 0 {
		mlog.ILT(mi18n.T("設定失敗"), mi18n.T("Step3失敗"))
		return
	}

	// 材質選択設定
	toolState.BlurModel.BlurMaterialIndexes = toolState.MaterialListBox.SelectedIndexes()
	toolState.BlurModel.OutputModel = nil
	toolState.BlurModel.OutputMotion = nil

}

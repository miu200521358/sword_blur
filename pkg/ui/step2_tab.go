package ui

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/sword_blur/pkg/usecase"
	"github.com/miu200521358/walk/pkg/walk"
)

func newStep2Tab(controlWindow *controller.ControlWindow, toolState *ToolState) {
	{
		toolState.Step2 = widget.NewMTabPage("Step.2")
		controlWindow.AddTabPage(toolState.Step2.TabPage)

		toolState.Step2.SetLayout(walk.NewVBoxLayout())

		{
			// Step2.文言
			label, err := walk.NewTextLabel(toolState.Step2)
			if err != nil {
				widget.RaiseError(err)
			}
			label.SetText(mi18n.T("Step2Label"))
		}

		walk.NewVSeparator(toolState.Step2)

		var err error
		{
			// クリアボタン
			toolState.Step2ClearButton, err = walk.NewPushButton(toolState.Step2)
			if err != nil {
				widget.RaiseError(err)
			}
			toolState.Step2ClearButton.SetText(mi18n.T("クリア"))
			toolState.Step2ClearButton.Clicked().Attach(toolState.onClickStep2Clear)
		}

		{
			// 材質選択リストボックス
			toolState.MaterialListBox, err = NewMaterialListBox(toolState.Step2)
			if err != nil {
				widget.RaiseError(err)
			}
		}

		{
			// OKボタン
			toolState.Step2OkButton, err = walk.NewPushButton(toolState.Step2)
			if err != nil {
				widget.RaiseError(err)
			}
			toolState.Step2OkButton.SetText(mi18n.T("次へ進む"))
			toolState.Step2OkButton.Clicked().Attach(toolState.onClickStep2Ok)
		}
	}
}

func (toolState *ToolState) onClickStep2Clear() {
	// 材質リストボックス設定
	toolState.MaterialListBox.SetMaterials(
		toolState.BlurModel.Model.Materials,
		toolState.onChangeMaterialListBox())
}

// Step2. 材質選択時
func (toolState *ToolState) onChangeMaterialListBox() func(indexes []int) {
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
			// 強制的にモーション更新するようハッシュ更新
			toolState.BlurModel.Motion.SetRandHash()
		}

		// outputPath := mutils.CreateOutputPath(
		// 	strings.ReplaceAll(toolState.BlurModel.Model.Path(), ".pmx", ".vmd"), "mat_off")
		// repository.NewVmdRepository().Save(outputPath, toolState.BlurModel.Motion, true)

		// 材質選択し直したら後続クリア
		toolState.SetEnabled(2)
		// 非表示材質設定
		toolState.ControlWindow.ChannelState().SetInvisibleMaterialsChannel([][][]int{{invisibleMaterialIndexes}})
	}
}

func (toolState *ToolState) onClickStep2Ok() {
	materialIndexes := toolState.MaterialListBox.SelectedIndexes()
	if len(materialIndexes) == 0 {
		mlog.ILT(mi18n.T("設定失敗"), mi18n.T("Step2失敗"))
		return
	}

	// 選択材質設定
	toolState.BlurModel.BlurMaterialIndexes = materialIndexes
	toolState.BlurModel.OutputModel = nil
	toolState.BlurModel.OutputMotion = nil

	// ワイヤーフレーム表示
	toolState.ControlWindow.SetShowWire(true)
	// 頂点選択ON
	toolState.ControlWindow.SetShowSelectedVertex(true)

	// 根元選択頂点に切り替え
	toolState.ResetSelectedVertexes(true, false, false, nil)
	// 選択更新メソッド設定
	toolState.ControlWindow.SetFuncSetSelectedVertexes(toolState.RootVertexSelectedFunc)

	toolState.ControlWindow.SetTabIndex(2) // Step3へ移動
	toolState.SetEnabled(3)
	mlog.IL(mi18n.T("Step2成功"))
}

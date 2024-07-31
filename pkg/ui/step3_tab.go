package ui

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
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

		{
			// 頂点選択時メソッド
			toolState.RootVertexSelectedFunc = func(indexes [][][]int) {
				// 重複頂点を同じINDEX位置で扱う
				indexMap := make(map[mmath.MVec3][]int)
				for _, vertexIndex := range indexes[0][0] {
					vertex := toolState.BlurModel.Model.Vertices.Get(vertexIndex)
					if _, ok := indexMap[*vertex.Position]; !ok {
						indexMap[*vertex.Position] = make([]int, 0)
					}
					indexMap[*vertex.Position] = append(indexMap[*vertex.Position], vertexIndex)
				}
				// 頂点リストボックス入替
				toolState.RootVertexListBox.ReplaceItems(indexMap)
			}
		}
	}
}

func (toolState *ToolState) onClickStep3Clear() {
	// 一旦選択機能OFFにして解除
	toolState.ControlWindow.SetShowSelectedVertex(false)
	// 再度有効
	toolState.ControlWindow.SetShowSelectedVertex(true)
}

func (toolState *ToolState) onClickStep3Ok() {
	if len(toolState.RootVertexListBox.SelectedIndexes()) == 0 {
		mlog.ILT(mi18n.T("設定失敗"), mi18n.T("Step3失敗"))
		return
	}

	// 根元選択設定
	toolState.BlurModel.BackRootVertexIndexes = toolState.RootVertexListBox.SelectedIndexes()
	toolState.BlurModel.OutputModel = nil
	toolState.BlurModel.OutputMotion = nil

}

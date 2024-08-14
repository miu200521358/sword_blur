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
			if !toolState.ControlWindow.IsShowSelectedVertex() {
				return
			}

			// 頂点選択し直したら後続クリア
			toolState.SetEnabled(3)

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

func (toolState *ToolState) onClickStep3Clear() {
	toolState.BlurModel.RootVertexIndexes = make([]int, 0)
	toolState.ResetSelectedVertexes(true, false, false, toolState.RootVertexListBox.GetItemValues())
}

func (toolState *ToolState) onClickStep3Ok() {
	rootVertexIndexes := toolState.RootVertexListBox.GetItemValues()
	if len(rootVertexIndexes) == 0 {
		mlog.ILT(mi18n.T("設定失敗"), mi18n.T("Step3失敗"))
		return
	}

	// 根元選択設定
	toolState.BlurModel.RootVertexIndexes = rootVertexIndexes
	toolState.BlurModel.OutputModel = nil
	toolState.BlurModel.OutputMotion = nil

	// ワイヤーフレーム表示
	toolState.ControlWindow.SetShowWire(true)
	// 頂点選択ON
	toolState.ControlWindow.SetShowSelectedVertex(true)

	// 切っ先頂点を選択
	toolState.ResetSelectedVertexes(false, true, false, nil)
	// 選択更新メソッド設定
	toolState.ControlWindow.SetUpdateSelectedVertexesFunc(toolState.TipVertexSelectedFunc)

	toolState.ControlWindow.SetTabIndex(3) // Step4へ移動
	toolState.SetEnabled(4)
	mlog.IL(mi18n.T("Step3成功"))
}

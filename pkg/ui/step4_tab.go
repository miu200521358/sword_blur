package ui

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/walk/pkg/walk"
)

func newStep4Tab(controlWindow *controller.ControlWindow, toolState *ToolState) {
	{
		toolState.Step4 = widget.NewMTabPage("Step.4")
		controlWindow.AddTabPage(toolState.Step4.TabPage)

		toolState.Step4.SetLayout(walk.NewVBoxLayout())

		{
			// Step4.文言
			label, err := walk.NewTextLabel(toolState.Step4)
			if err != nil {
				widget.RaiseError(err)
			}
			label.SetText(mi18n.T("Step4Label"))
		}

		walk.NewVSeparator(toolState.Step4)

		var err error
		{
			// クリアボタン
			toolState.Step4ClearButton, err = walk.NewPushButton(toolState.Step4)
			if err != nil {
				widget.RaiseError(err)
			}
			toolState.Step4ClearButton.SetText(mi18n.T("クリア"))
			toolState.Step4ClearButton.Clicked().Attach(toolState.onClickStep4Clear)
		}

		{
			// 頂点選択リストボックス
			toolState.TipVertexListBox, err = NewVertexListBox(toolState.Step4)
			if err != nil {
				widget.RaiseError(err)
			}
		}

		{
			// OKボタン
			toolState.Step4OkButton, err = walk.NewPushButton(toolState.Step4)
			if err != nil {
				widget.RaiseError(err)
			}
			toolState.Step4OkButton.SetText(mi18n.T("次へ進む"))
			toolState.Step4OkButton.Clicked().Attach(toolState.onClickStep4Ok)
		}

		{
			// 頂点選択時メソッド
			toolState.TipVertexSelectedFunc = func(indexes [][][]int) {
				// 頂点選択し直したら後続クリア
				toolState.SetEnabled(4)

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
				toolState.TipVertexListBox.ReplaceItems(indexMap)
			}
		}
	}
}

func (toolState *ToolState) onClickStep4Clear() {
	toolState.BlurModel.TipVertexIndexes = make([]int, 0)
	toolState.ResetSelectedVertexIndexes(false, true, false, toolState.TipVertexListBox.GetItemValues())
}

func (toolState *ToolState) onClickStep4Ok() {
	tipVertexIndexes := toolState.TipVertexListBox.GetItemValues()
	if len(tipVertexIndexes) == 0 {
		mlog.ILT(mi18n.T("設定失敗"), mi18n.T("Step4失敗"))
		return
	}

	// 根元選択設定
	toolState.BlurModel.TipVertexIndexes = tipVertexIndexes
	toolState.BlurModel.OutputModel = nil
	toolState.BlurModel.OutputMotion = nil

	// ワイヤーフレーム表示
	toolState.ControlWindow.SetShowWire(true)
	// 頂点選択ON
	toolState.ControlWindow.SetShowSelectedVertex(true)

	// 切っ先頂点を選択
	toolState.ResetSelectedVertexIndexes(false, false, true, nil)
	// 選択更新メソッド設定
	toolState.ControlWindow.SetUpdateSelectedVertexIndexesFunc(toolState.EdgeVertexSelectedFunc)

	toolState.ControlWindow.TriggerPlay(false)

	toolState.ControlWindow.SetTabIndex(4) // Step5へ移動
	toolState.SetEnabled(5)
	mlog.IL(mi18n.T("Step4成功"))
}

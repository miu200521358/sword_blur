package ui

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/sword_blur/pkg/usecase"
	"github.com/miu200521358/walk/pkg/walk"
)

func newStep5Tab(controlWindow *controller.ControlWindow, toolState *ToolState) {

	toolState.Step5 = widget.NewMTabPage("Step.5")
	controlWindow.AddTabPage(toolState.Step5.TabPage)

	toolState.Step5.SetLayout(walk.NewVBoxLayout())

	{
		// Step5.文言
		label, err := walk.NewTextLabel(toolState.Step5)
		if err != nil {
			widget.RaiseError(err)
		}
		label.SetText(mi18n.T("Step5Label"))
	}

	walk.NewVSeparator(toolState.Step5)

	var err error
	{
		// クリアボタン
		toolState.Step5ClearButton, err = walk.NewPushButton(toolState.Step5)
		if err != nil {
			widget.RaiseError(err)
		}
		toolState.Step5ClearButton.SetText(mi18n.T("クリア"))
		toolState.Step5ClearButton.Clicked().Attach(toolState.onClickStep5Clear)
	}

	{
		// 頂点選択リストボックス
		toolState.EdgeVertexListBox, err = NewVertexListBox(toolState.Step5)
		if err != nil {
			widget.RaiseError(err)
		}
	}

	{
		// プレビューボタン
		toolState.Step5PreviewButton, err = walk.NewPushButton(toolState.Step5)
		if err != nil {
			widget.RaiseError(err)
		}
		toolState.Step5PreviewButton.SetText(mi18n.T("プレビュー"))
		toolState.Step5PreviewButton.Clicked().Attach(toolState.onClickStep5Preview)
	}

	{
		// リトライボタン
		toolState.Step5RetryButton, err = walk.NewPushButton(toolState.Step5)
		if err != nil {
			widget.RaiseError(err)
		}
		toolState.Step5RetryButton.SetText(mi18n.T("リトライ"))
		toolState.Step5RetryButton.Clicked().Attach(toolState.onClickStep5Retry)
	}

	{
		// 保存ボタン
		toolState.Step5SaveButton, err = walk.NewPushButton(toolState.Step5)
		if err != nil {
			widget.RaiseError(err)
		}
		toolState.Step5SaveButton.SetText(mi18n.T("保存"))
		toolState.Step5SaveButton.Clicked().Attach(toolState.onClickStep5Save)
	}

	{
		// 頂点選択時メソッド
		toolState.EdgeVertexSelectedFunc = func(indexes [][][]int) {
			if !toolState.ControlWindow.IsShowSelectedVertex() {
				return
			}

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
			if toolState.EdgeVertexListBox.ReplaceItems(indexMap) {
				// 選択頂点INDEX更新
				toolState.BlurModel.EdgeVertexIndexes = make([]int, 0)
				// 頂点選択し直したら後続クリア
				toolState.SetEnabled(5)
			}
		}
	}
}

func (toolState *ToolState) onClickStep5Clear() {
	toolState.BlurModel.EdgeVertexIndexes = make([]int, 0)
	toolState.ResetSelectedVertexes(false, false, true, toolState.EdgeVertexListBox.GetItemValues())
}

func (toolState *ToolState) onClickStep5Preview() {
	edgeVertexIndexes := toolState.EdgeVertexListBox.GetItemValues()
	if len(edgeVertexIndexes) == 0 {
		mlog.ILT(mi18n.T("設定失敗"), mi18n.T("Step5失敗"))
		return
	}

	// 根元選択設定
	toolState.BlurModel.EdgeVertexIndexes = edgeVertexIndexes
	// 出力用モデルを別で読み込み
	outputModel := toolState.OriginalPmxPicker.LoadForce(toolState.OriginalPmxPicker.GetPath()).(*pmx.PmxModel)

	var err error
	toolState.BlurModel.OutputModel, toolState.BlurModel.OutputMotion, err = usecase.Preview(
		toolState.BlurModel, outputModel)
	if err != nil {
		mlog.ET(mi18n.T("生成失敗"), mi18n.T("生成失敗メッセージ", map[string]interface{}{"Error": err.Error()}))
		return
	} else {
		mlog.IT(mi18n.T("生成成功"), mi18n.T("生成成功メッセージ"))
	}

	toolState.ControlWindow.UpdateMaxFrame(toolState.BlurModel.OutputMotion.MaxFrame())

	// ワイヤーフレーム非表示
	toolState.ControlWindow.SetShowWire(false)
	// 頂点選択OFF
	toolState.ControlWindow.SetShowSelectedVertex(false)

	// 再生ON
	toolState.Player.SetPlaying(true)
}

func (toolState *ToolState) onClickStep5Retry() {
	// ワイヤーフレーム切り替え
	toolState.ControlWindow.SetShowWire(true)
	// 頂点選択切り替え
	toolState.ControlWindow.SetShowSelectedVertex(true)
	toolState.ResetSelectedVertexes(false, false, true, nil)
	// 再生停止
	toolState.Player.SetPlaying(false)
}

func (toolState *ToolState) onClickStep5Save() {
	toolState.Player.SetPlaying(false)

	if toolState.BlurModel.OutputModel == nil {
		mlog.IT(mi18n.T("出力モデルなし"), mi18n.T("出力モデルなしメッセージ"))
		return
	}

	if err := usecase.Save(toolState.BlurModel); err != nil {
		mlog.ET(mi18n.T("出力失敗"), mi18n.T("出力失敗メッセージ", map[string]interface{}{"Error": err.Error()}))
	} else {
		mlog.IT(mi18n.T("出力成功"), mi18n.T("出力成功メッセージ",
			map[string]interface{}{"Path": toolState.BlurModel.OutputModelPath}))
	}

	widget.Beep()
}

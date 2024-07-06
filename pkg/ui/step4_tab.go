package ui

import (
	"slices"

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

	// OKボタン
	stp.Items.okButton, err = walk.NewPushButton(stp.Items.composite)
	if err != nil {
		return nil, err
	}
	stp.Items.okButton.SetText(mi18n.T("出力"))

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
			model := step3Page.prevStep.prevStep.Items.OriginalPmxPicker.GetCache().(*pmx.PmxModel)
			// 直近頂点を取得
			tempVertex := pmx.NewVertex()
			tempVertex.Position = worldPos
			vertexIndexes, vertexPositions := model.Vertices.GetMapValues(tempVertex)
			if len(vertexIndexes) > 0 {
				visibleVertexIndexes := make([]int, 0)
				visibleVertexPositions := make([]*mmath.MVec3, 0)
				for i, vertexIndex := range vertexIndexes {
					for _, materialIndex := range step3Page.prevStep.Items.MaterialListBox.SelectedIndexes() {
						// 表示されている材質からのみ直近頂点を選ぶ
						if slices.Contains(model.Vertices.Get(vertexIndex).MaterialIndexes, materialIndex) &&
							!slices.Contains(visibleVertexIndexes, vertexIndex) {
							visibleVertexIndexes = append(visibleVertexIndexes, vertexIndex)
							visibleVertexPositions = append(visibleVertexPositions, vertexPositions[i])
							break
						}
					}
				}
				distances := mmath.Float64Slice(mmath.Distances(worldPos, visibleVertexPositions))
				nearVertexIndexes := mmath.ArgSort(distances)
				targetVertexIndexes := make([]int, 0)
				for i, nearVertexIndex := range nearVertexIndexes {
					if i == 0 || visibleVertexPositions[nearVertexIndexes[0]].NearEquals(visibleVertexPositions[nearVertexIndexes[i]], 1e-2) {
						// 直近とほぼ同じ位置の頂点を選択
						targetVertexIndexes = append(targetVertexIndexes, visibleVertexIndexes[nearVertexIndex])
					}
				}
				stp.Items.VertexListBox.SetItem(targetVertexIndexes)

				go func() {
					mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]*mwidget.ModelSet{0: {NextSelectedVertexIndexes: stp.Items.VertexListBox.GetItemValues()}}
				}()
			}
		}
	}

	// Step4. OKボタンクリック時
	stp.Items.okButton.Clicked().Attach(func() {
		if len(stp.Items.VertexListBox.GetItemValues()) == 0 {
			stp.SetEnabled(false)
			mlog.IL(mi18n.T("Step4刃頂点設定失敗"))
			return
		} else {
			err := usecase.Save(
				stp.prevStep.prevStep.prevStep.Items.OriginalPmxPicker.GetDataForce().(*pmx.PmxModel),
				stp.prevStep.prevStep.prevStep.Items.OutputPmxPicker.GetPath(),
				stp.prevStep.prevStep.Items.MaterialListBox.SelectedIndexes(),
				stp.prevStep.Items.VertexListBox.GetItemValues(),
				stp.Items.VertexListBox.GetItemValues(),
			)

			if err != nil {
				mlog.ET(mi18n.T("出力失敗"), mi18n.T("出力失敗メッセージ", map[string]interface{}{"Error": err.Error()}))
			} else {
				mlog.IT(mi18n.T("出力成功"), mi18n.T("出力成功メッセージ", map[string]interface{}{"Path": stp.prevStep.prevStep.prevStep.Items.OutputPmxPicker.GetPath()}))
			}
		}
	})

	return stp, nil
}

// ------------------------------

type Step4TabPage struct {
	*mwidget.MTabPage
	mWindow  *mwidget.MWindow
	prevStep *Step3TabPage
	Items    *Step4Items
}

// ------------------------------

type Step4Items struct {
	stepItems
	VertexListBox *VertexListBox
	okButton      *walk.PushButton
	FuncWorldPos  func(worldPos *mmath.MVec3, vmdDeltas []*vmd.VmdDeltas, viewMat *mmath.MMat4)
}

func (si *Step4Items) SetEnabled(enabled bool) {
	si.stepItems.SetEnabled(enabled)
	si.VertexListBox.SetEnabled(enabled)
	si.okButton.SetEnabled(enabled)
}

func (si *Step4Items) Dispose() {
	si.stepItems.Dispose()
	si.VertexListBox.Dispose()
	si.okButton.Dispose()
}

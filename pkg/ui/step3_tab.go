package ui

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/mlib_go/pkg/mwidget"
	"github.com/miu200521358/mlib_go/pkg/pmx"
	"github.com/miu200521358/mlib_go/pkg/vmd"
	"github.com/miu200521358/walk/pkg/walk"
)

func NewStep3TabPage(mWindow *mwidget.MWindow, step2Page *Step2TabPage) (*Step3TabPage, error) {
	page, err := mwidget.NewMTabPage(mWindow, mWindow.TabWidget, "Step. 3")
	if err != nil {
		return nil, err
	}
	page.SetLayout(walk.NewVBoxLayout())

	stp := &Step3TabPage{
		MTabPage: page,
		mWindow:  mWindow,
		prevStep: step2Page,
		Items:    &Step3Items{},
	}

	// Step3. 峰選択

	stp.Items.composite, err = walk.NewComposite(stp)
	if err != nil {
		return nil, err
	}
	stp.Items.composite.SetLayout(walk.NewVBoxLayout())

	stp.Items.label, err = walk.NewTextLabel(stp.Items.composite)
	if err != nil {
		return nil, err
	}
	stp.Items.label.SetText(mi18n.T("Step3Label"))

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
	stp.Items.okButton.SetText(mi18n.T("次へ進む"))

	stp.SetEnabled(false)

	// Step2. OKボタンクリック時
	step2Page.Items.okButton.Clicked().Attach(func() {
		if len(step2Page.Items.MaterialListBox.SelectedIndexes()) == 0 {
			stp.SetEnabled(false)
			mlog.IL(mi18n.T("Step2材質設定失敗"))
			return
		} else {
			stp.SetEnabled(true)
			stp.mWindow.SetCheckWireDebugView(true)
			stp.mWindow.SetCheckSelectedVertexDebugView(true)
			stp.mWindow.TabWidget.SetCurrentIndex(2)                              // Step3へ移動
			stp.mWindow.GetMainGlWindow().SetFuncWorldPos(stp.Items.FuncWorldPos) // 頂点選択時のターゲットfunction変更

			go func() {
				mWindow.GetMainGlWindow().ReplaceModelSetChannel <- map[int]*mwidget.ModelSet{0: {NextSelectedVertexIndexes: []int{}}}
			}()

			mlog.IL(mi18n.T("Step2材質設定完了"))
		}
	})

	stp.Items.FuncWorldPos = func(worldPos *mmath.MVec3, vmdDeltas []*vmd.VmdDeltas, viewMat *mmath.MMat4) {
		if step2Page.prevStep.Items.OriginalPmxPicker.Exists() && stp.Enabled() {
			model := step2Page.prevStep.Items.OriginalPmxPicker.GetCache().(*pmx.PmxModel)
			// 直近頂点を取得
			tempVertex := pmx.NewVertex()
			tempVertex.Position = worldPos
			vertexIndexes, vertexPositions := model.Vertices.GetMapValues(tempVertex)
			if len(vertexIndexes) > 0 {
				visibleVertexIndexes := make([]int, 0)
				visibleVertexPositions := make([]*mmath.MVec3, 0)
				for i, vertexIndex := range vertexIndexes {
					for _, materialIndex := range step2Page.Items.MaterialListBox.SelectedIndexes() {
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

	return stp, nil
}

// ------------------------------

type Step3TabPage struct {
	*mwidget.MTabPage
	mWindow  *mwidget.MWindow
	prevStep *Step2TabPage
	Items    *Step3Items
}

// ------------------------------

type Step3Items struct {
	stepItems
	VertexListBox *VertexListBox
	okButton      *walk.PushButton
	FuncWorldPos  func(worldPos *mmath.MVec3, vmdDeltas []*vmd.VmdDeltas, viewMat *mmath.MMat4)
}

func (si *Step3Items) SetEnabled(enabled bool) {
	si.stepItems.SetEnabled(enabled)
	si.VertexListBox.SetEnabled(enabled)
	si.okButton.SetEnabled(enabled)
}

func (si *Step3Items) Dispose() {
	si.stepItems.Dispose()
	si.VertexListBox.Dispose()
	si.okButton.Dispose()
}

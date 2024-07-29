package usecase

import (
	"fmt"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/sword_blur/pkg/model"
)

// SetupModel ツール専用セットアップ
func SetupModel(blurModel *model.BlurModel) {
	// 材質モーフ追加
	addMaterialMorph(blurModel.Model)

	outputPath := mutils.CreateOutputPath(blurModel.Model.Path(), "debug")
	repository.NewPmxRepository().Save(outputPath, blurModel.Model, true)
}

func addMaterialMorph(model *pmx.PmxModel) {
	// 材質モーフ追加
	for i := range model.Materials.Len() {
		material := model.Materials.Get(i)
		offset := pmx.NewMaterialMorphOffset(
			material.Index(),
			pmx.CALC_MODE_MULTIPLICATION,
			// 透明度だけを操作する（最初から透明なものはそのまま）
			&mmath.MVec4{X: 1.0, Y: 1.0, Z: 1.0, W: 0.0},
			&mmath.MVec4{X: 1.0, Y: 1.0, Z: 1.0, W: 1.0},
			&mmath.MVec3{X: 1.0, Y: 1.0, Z: 1.0},
			&mmath.MVec4{X: 1.0, Y: 1.0, Z: 1.0, W: 0.0},
			0.0,
			&mmath.MVec4{X: 1.0, Y: 1.0, Z: 1.0, W: 1.0},
			&mmath.MVec4{X: 1.0, Y: 1.0, Z: 1.0, W: 1.0},
			&mmath.MVec4{X: 1.0, Y: 1.0, Z: 1.0, W: 1.0},
		)
		morph := pmx.NewMorph()
		morph.SetIndex(model.Morphs.Len())
		morph.SetName(getVisibleMorphName(material))
		morph.Offsets = append(morph.Offsets, offset)
		morph.MorphType = pmx.MORPH_TYPE_MATERIAL
		morph.Panel = pmx.MORPH_PANEL_OTHER_LOWER_RIGHT
		morph.IsSystem = true
		model.Morphs.Append(morph)
	}
}

func getVisibleMorphName(material *pmx.Material) string {
	return fmt.Sprintf("%s_%d_%s", pmx.MLIB_PREFIX, material.Index(), material.Name())
}

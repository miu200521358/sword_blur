package usecase

import (
	"fmt"

	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

// SetupModel ツール専用セットアップ
func SetupModel(model *pmx.PmxModel) {
	// 材質モーフ追加
	addMaterialMorph(model)

	// dir, file := filepath.Split(model.GetPath())
	// ext := filepath.Ext(file)
	// outputPath := filepath.Join(dir, file[:len(file)-len(ext)]+"_debug"+ext)

	// model.Save(true, outputPath)
}

func addMaterialMorph(model *pmx.PmxModel) {
	// 材質モーフ追加
	for i := range model.Materials.Len() {
		material := model.Materials.Get(i)
		offset := pmx.NewMaterialMorphOffset(
			material.Index,
			pmx.CALC_MODE_MULTIPLICATION,
			// 透明度だけを操作する（最初から透明なものはそのまま）
			&mmath.MVec4{1.0, 1.0, 1.0, 0.0},
			&mmath.MVec4{1.0, 1.0, 1.0, 1.0},
			&mmath.MVec3{1.0, 1.0, 1.0},
			&mmath.MVec4{1.0, 1.0, 1.0, 0.0},
			0.0,
			&mmath.MVec4{1.0, 1.0, 1.0, 1.0},
			&mmath.MVec4{1.0, 1.0, 1.0, 1.0},
			&mmath.MVec4{1.0, 1.0, 1.0, 1.0},
		)
		morph := pmx.NewMorph()
		morph.Name = GetVisibleMorphName(material)
		morph.Offsets = append(morph.Offsets, offset)
		morph.MorphType = pmx.MORPH_TYPE_MATERIAL
		morph.Panel = pmx.MORPH_PANEL_OTHER_LOWER_RIGHT
		morph.IsSystem = true
		model.Morphs.Append(morph)
	}
}

func GetVisibleMorphName(material *pmx.Material) string {
	return fmt.Sprintf("%s_%d_%s", pmx.MLIB_PREFIX, material.Index, material.Name)
}

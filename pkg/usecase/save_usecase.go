package usecase

import (
	"fmt"

	"github.com/miu200521358/mlib_go/pkg/pmx"
)

// ブレ設定入りモデル出力処理
func Save(
	model *pmx.PmxModel, outputPath string, blurMaterialIndexes []int, backVertexIndexes []int, edgeVertexIndexes []int,
) error {
	// ブレ材質の複製
	for _, materialIndex := range blurMaterialIndexes {
		m := model.Materials.Get(materialIndex)
		blurName := fmt.Sprintf("%sブレ", m.Name)
		blurEnglishName := fmt.Sprintf("%s_blur", m.EnglishName)
		model.DuplicateMaterial(materialIndex, blurName, blurEnglishName)
	}

	// 保存
	model.Save(false, outputPath)

	return nil
}

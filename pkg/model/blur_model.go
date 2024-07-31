package model

import (
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
)

type BlurModel struct {
	Model               *pmx.PmxModel  // 処理対象モデル
	Motion              *vmd.VmdMotion // 処理対象モーション
	OutputModelPath     string         // 出力パス
	BlurMaterialIndexes []int          // ブレ対象材質インデックス
	RootVertexIndexes   []int          // 棟区頂点INDEX
	TipVertexIndexes    []int          // 切っ先頂点INDEX
	EdgeVertexIndexes   []int          // 刃頂点INDEX
	OutputModel         *pmx.PmxModel  // 出力モデル
	OutputMotion        *vmd.VmdMotion // 出力モーション
}

func NewBlurModel() *BlurModel {
	return &BlurModel{}
}

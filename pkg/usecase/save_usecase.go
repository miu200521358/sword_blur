package usecase

import (
	"fmt"

	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

// ブレ設定入りモデル出力処理
func Save(
	model *pmx.PmxModel, outputPath string, blurMaterialIndexes []int, backVertexIndexes []int, edgeVertexIndexes []int,
) error {
	// 表示枠の追加
	blurDisplaySlot := createBlurDisplaySlot(model)

	backRootVertex := model.Vertices.Get(backVertexIndexes[0])
	edgeRootVertex := model.Vertices.Get(edgeVertexIndexes[0])

	backTailVertex := model.Vertices.Get(backVertexIndexes[len(backVertexIndexes)-1])
	edgeTailVertex := model.Vertices.Get(edgeVertexIndexes[len(edgeVertexIndexes)-1])

	rootVector := edgeRootVertex.Position.Subed(backRootVertex.Position).Normalize()
	backVector := backTailVertex.Position.Subed(backRootVertex.Position).Normalize()
	rootPosition := backRootVertex.Position.Added(edgeRootVertex.Position).MuledScalar(0.5)
	tailPosition := backTailVertex.Position.Added(edgeTailVertex.Position).MuledScalar(0.5)

	// ボーンの追加
	blurRootBone := createBlurRootBone(model, blurDisplaySlot, backRootVertex, rootPosition)
	blurBone := createBlurBone(model, blurDisplaySlot, blurRootBone, rootPosition, tailPosition, rootVector)
	blurTailBone := createBlurTailBone(model, blurDisplaySlot, blurRootBone, tailPosition)
	blurIkBone := createBlurIkBone(model, blurDisplaySlot, blurRootBone, blurBone, blurTailBone)
	blurWeightBone := createBlurWeightBone(model, blurDisplaySlot, blurRootBone, blurBone)
	blurWeightTailBone := createBlurWeightTailBone(model, blurDisplaySlot, blurRootBone, blurTailBone)

	mlog.V("blurIkBone: %v", blurIkBone)
	mlog.V("blurWeightTailBone: %v", blurWeightTailBone)

	blurMaterials := make([]*pmx.Material, 0)

	// ブレ材質の複製
	for _, materialIndex := range blurMaterialIndexes {
		m := model.Materials.Get(materialIndex)
		blurName := fmt.Sprintf("%s_ブレ", m.Name)
		blurEnglishName := fmt.Sprintf("%s_blur", m.EnglishName)
		duplicateMaterial(model, materialIndex, blurName, blurEnglishName,
			[]int{blurWeightBone.Index}, blurWeightTailBone.Index, edgeVertexIndexes)

		// ブレ材質の設定
		blurMaterial := model.Materials.GetByName(blurName)
		blurMaterial.DrawFlag = blurMaterial.DrawFlag.SetDoubleSidedDrawing(true)       // 両面描画ON
		blurMaterial.DrawFlag = blurMaterial.DrawFlag.SetDrawingOnSelfShadowMaps(false) // セルフシャドウマップ描画OFF
		blurMaterial.DrawFlag = blurMaterial.DrawFlag.SetDrawingSelfShadows(false)      // セルフシャドウ描画OFF
		blurMaterial.DrawFlag = blurMaterial.DrawFlag.SetGroundShadow(false)            // 地面影描画OFF
		blurMaterial.DrawFlag = blurMaterial.DrawFlag.SetDrawingEdge(false)             // エッジ描画OFF

		// エッジの設定
		blurMaterial.EdgeSize = 0.0
		// エッジの色
		blurMaterial.Edge = mmath.NewMVec4()

		// 非透過度(デフォルト透明)
		blurMaterial.Diffuse.SetW(0.0)

		// 反射強度
		blurMaterial.Specular.SetW(1.0)

		blurMaterials = append(blurMaterials, blurMaterial)
	}

	// 表示モーフ
	createDiffuseMorph(model, blurMaterials, &mmath.MVec4{0.0, 0.0, 0.0, 1.0}, "表示", "visible")
	createTextureMorph(model, blurMaterials, &mmath.MVec4{1.0, 0.0, 0.0, 0.0}, "赤", "red")
	createTextureMorph(model, blurMaterials, &mmath.MVec4{0.0, 1.0, 0.0, 0.0}, "緑", "green")
	createTextureMorph(model, blurMaterials, &mmath.MVec4{0.0, 0.0, 1.0, 0.0}, "青", "blue")
	createSpecularMorph(model, blurMaterials, &mmath.MVec4{0.0, 0.0, 0.0, 150.0}, "AL", "AL")

	// 剛体
	blurRootRigidBody := createBlurRootRigidBody(model, blurRootBone)
	blurRigidBody := createBlurRigidBody(model, blurBone, backVertexIndexes, edgeVertexIndexes)

	// ジョイント
	createJoint(model, blurBone, blurRootRigidBody, blurRigidBody, rootVector, backVector)

	// 保存
	model.Save(false, outputPath)

	return nil
}

// DuplicateMaterial は指定された材質を複製します。
func duplicateMaterial(
	model *pmx.PmxModel, materialIndex int, materialName, materialEnglishName string,
	weightBoneIndexes []int, weightTailBoneIndex int, edgeVertexIndexes []int,
) error {
	material := model.Materials.Get(materialIndex)
	duplicatedMaterial := material.Copy().(*pmx.Material)
	duplicatedMaterial.Index = model.Materials.Len()
	duplicatedMaterial.Name = materialName
	duplicatedMaterial.EnglishName = materialEnglishName

	// 該当材質内の頂点と面を取得
	prevVerticesCount := 0
	for i := range model.Materials.Len() {
		if i < materialIndex {
			prevVerticesCount += int(model.Materials.Get(i).VerticesCount / 3)
			continue
		}

		// 該当材質の場合
		m := model.Materials.Get(i)
		for j := prevVerticesCount; j < prevVerticesCount+int(m.VerticesCount/3); j++ {
			originalFace := model.Faces.Get(j)
			duplicatedFace := originalFace.Copy().(*pmx.Face)
			duplicatedFace.Index = model.Faces.Len()

			for k, vertexIndex := range originalFace.VertexIndexes {
				originalVertex := model.Vertices.Get(vertexIndex)
				duplicatedVertex := originalVertex.Copy().(*pmx.Vertex)
				duplicatedVertex.Index = model.Vertices.Len()

				if mmath.Contains(edgeVertexIndexes, vertexIndex) {
					// 刃側の頂点ウェイトはW先
					duplicatedVertex.Deform = pmx.NewBdef1(weightTailBoneIndex)
					duplicatedVertex.DeformType = pmx.BDEF1
				} else if len(weightBoneIndexes) > 0 {
					// ウェイトボーンが指定されている場合、ウェイトボーンのインデックスを更新
					duplicatedVertex.Deform = pmx.NewBdef1(weightBoneIndexes[0])
					duplicatedVertex.DeformType = pmx.BDEF1
				}

				model.Vertices.Append(duplicatedVertex)
				duplicatedFace.VertexIndexes[2-k] = duplicatedVertex.Index
			}

			model.Faces.Append(duplicatedFace)
		}

		// 該当材質の複製が終わったら終了
		break
	}

	model.Materials.Append(duplicatedMaterial)

	return nil
}

func createBlurDisplaySlot(model *pmx.PmxModel) *pmx.DisplaySlot {
	blurDisplaySlot := pmx.NewDisplaySlot()
	blurDisplaySlot.Index = model.DisplaySlots.Len()
	blurDisplaySlot.Name = "ブレ"
	blurDisplaySlot.EnglishName = "Blur"
	blurDisplaySlot.SpecialFlag = pmx.SPECIAL_FLAG_OFF

	model.DisplaySlots.Append(blurDisplaySlot)

	return blurDisplaySlot

}

func createBlurRootBone(
	model *pmx.PmxModel, blurDisplaySlot *pmx.DisplaySlot, backRootVertex *pmx.Vertex, rootPosition *mmath.MVec3,
) *pmx.Bone {
	blurRootBone := pmx.NewBone()
	blurRootBone.Index = model.Bones.Len()
	blurRootBone.Name = "ブレ根元"
	blurRootBone.EnglishName = "Blur Root"
	blurRootBone.IsSystem = false // システム用ではなく出力用

	// 元々刀身のウェイトが乗っていたボーンを親ボーンとする
	blurRootBone.ParentIndex = backRootVertex.Deform.GetAllIndexes()[0]

	blurRootBone.Position = rootPosition.Copy()
	// 移動可能, 回転可能, 操作可能, 表示枠追加
	blurRootBone.BoneFlag = pmx.BONE_FLAG_CAN_TRANSLATE | pmx.BONE_FLAG_CAN_ROTATE | pmx.BONE_FLAG_CAN_MANIPULATE | pmx.BONE_FLAG_IS_VISIBLE
	// 表示枠追加
	blurDisplaySlot.References = append(blurDisplaySlot.References,
		pmx.Reference{DisplayType: pmx.DISPLAY_TYPE_BONE, DisplayIndex: blurRootBone.Index})

	model.Bones.Append(blurRootBone)

	return blurRootBone
}

func createBlurBone(
	model *pmx.PmxModel, blurDisplaySlot *pmx.DisplaySlot, blurRootBone *pmx.Bone,
	rootPosition, tailPosition, rootVector *mmath.MVec3,
) *pmx.Bone {
	blurBone := pmx.NewBone()
	blurBone.Index = model.Bones.Len()
	blurBone.Name = "ブレ"
	blurBone.EnglishName = "Blur"
	blurBone.ParentIndex = blurRootBone.Index
	blurBone.IsSystem = false // システム用ではなく出力用

	blurBone.Position = blurRootBone.Position.Copy()
	// 移動可能, 回転可能, 操作可能, 表示枠追加, 軸制限, 物理後
	blurBone.BoneFlag = pmx.BONE_FLAG_CAN_TRANSLATE | pmx.BONE_FLAG_CAN_ROTATE | pmx.BONE_FLAG_CAN_MANIPULATE | pmx.BONE_FLAG_IS_VISIBLE | pmx.BONE_FLAG_HAS_FIXED_AXIS | pmx.BONE_FLAG_IS_AFTER_PHYSICS_DEFORM
	// 表示枠追加
	blurDisplaySlot.References = append(blurDisplaySlot.References,
		pmx.Reference{DisplayType: pmx.DISPLAY_TYPE_BONE, DisplayIndex: blurBone.Index})

	// 軸制限(先端から根元までのベクトルの外積を軸にする)
	blurBone.FixedAxis = tailPosition.Sub(rootPosition).Normalized().Cross(rootVector)

	model.Bones.Append(blurBone)

	return blurBone
}

func createBlurTailBone(
	model *pmx.PmxModel, blurDisplaySlot *pmx.DisplaySlot, blurRootBone *pmx.Bone, tailPosition *mmath.MVec3,
) *pmx.Bone {
	blurTailBone := pmx.NewBone()
	blurTailBone.Index = model.Bones.Len()
	blurTailBone.Name = "ブレ先"
	blurTailBone.EnglishName = "Blur Tail"
	blurTailBone.ParentIndex = blurRootBone.Index
	blurTailBone.IsSystem = false // システム用ではなく出力用
	blurTailBone.IsSystem = false // システム用ではなく出力用

	// 移動可能, 回転可能, 操作可能, 表示枠追加, 物理後
	blurTailBone.BoneFlag = pmx.BONE_FLAG_CAN_TRANSLATE | pmx.BONE_FLAG_CAN_ROTATE | pmx.BONE_FLAG_CAN_MANIPULATE | pmx.BONE_FLAG_IS_VISIBLE | pmx.BONE_FLAG_IS_AFTER_PHYSICS_DEFORM
	// 表示枠追加
	blurDisplaySlot.References = append(blurDisplaySlot.References,
		pmx.Reference{DisplayType: pmx.DISPLAY_TYPE_BONE, DisplayIndex: blurTailBone.Index})

	// 選択の頂点位置の中間をボーン位置とする
	blurTailBone.Position = tailPosition.Copy()

	model.Bones.Append(blurTailBone)

	return blurTailBone
}

func createBlurIkBone(
	model *pmx.PmxModel, blurDisplaySlot *pmx.DisplaySlot, blurRootBone *pmx.Bone,
	blurBone *pmx.Bone, blurTailBone *pmx.Bone,
) *pmx.Bone {
	blurIkBone := pmx.NewBone()
	blurIkBone.Index = model.Bones.Len()
	blurIkBone.Name = "ブレIK"
	blurIkBone.EnglishName = "Blur IK"
	blurIkBone.ParentIndex = blurRootBone.Index
	blurIkBone.IsSystem = false // システム用ではなく出力用

	blurIkBone.Position = blurTailBone.Position.Copy()
	// 移動可能, 回転可能, 操作可能, 表示枠追加, 物理後, IK
	blurIkBone.BoneFlag = pmx.BONE_FLAG_CAN_TRANSLATE | pmx.BONE_FLAG_CAN_ROTATE | pmx.BONE_FLAG_CAN_MANIPULATE | pmx.BONE_FLAG_IS_VISIBLE | pmx.BONE_FLAG_IS_AFTER_PHYSICS_DEFORM | pmx.BONE_FLAG_IS_IK
	// 表示枠追加
	blurDisplaySlot.References = append(blurDisplaySlot.References,
		pmx.Reference{DisplayType: pmx.DISPLAY_TYPE_BONE, DisplayIndex: blurIkBone.Index})

	// IK設定
	blurIk := pmx.NewIk()
	blurIk.BoneIndex = blurTailBone.Index
	blurIk.LoopCount = 1
	blurIk.UnitRotation.GetRadians().SetX(2.0)

	blurIkLink := pmx.NewIkLink()
	blurIkLink.BoneIndex = blurBone.Index

	// blurIkLink.AngleLimit = true
	// // 軸制限の向きに固定角度
	// limitDegrees := blurBone.FixedAxis.ToMat4().MulVec3(&mmath.MVec3{45, 0, 0})
	// minLimitDegrees := mmath.NewMVec3()
	// maxLimitDegrees := mmath.NewMVec3()
	// if limitDegrees.GetX() < 0 {
	// 	minLimitDegrees.SetX(limitDegrees.GetX())
	// } else {
	// 	maxLimitDegrees.SetX(limitDegrees.GetX())
	// }
	// if limitDegrees.GetY() < 0 {
	// 	minLimitDegrees.SetY(limitDegrees.GetY())
	// } else {
	// 	maxLimitDegrees.SetY(limitDegrees.GetY())
	// }
	// if limitDegrees.GetZ() < 0 {
	// 	minLimitDegrees.SetZ(limitDegrees.GetZ())
	// } else {
	// 	maxLimitDegrees.SetZ(limitDegrees.GetZ())
	// }
	// blurIkLink.MinAngleLimit.SetDegrees(minLimitDegrees)
	// blurIkLink.MaxAngleLimit.SetDegrees(maxLimitDegrees)

	blurIk.Links = append(blurIk.Links, blurIkLink)
	blurIkBone.Ik = blurIk

	model.Bones.Append(blurIkBone)

	return blurIkBone
}

func createBlurWeightBone(
	model *pmx.PmxModel, blurDisplaySlot *pmx.DisplaySlot, blurRootBone *pmx.Bone, blurBone *pmx.Bone,
) *pmx.Bone {
	blurWeightBone := pmx.NewBone()
	blurWeightBone.Index = model.Bones.Len()
	blurWeightBone.Name = "ブレW"
	blurWeightBone.EnglishName = "Blur W"
	blurWeightBone.ParentIndex = blurRootBone.Index
	blurWeightBone.IsSystem = false // システム用ではなく出力用

	blurWeightBone.Position = blurRootBone.Position.Copy()
	// 移動可能, 回転可能, 操作可能, 表示枠追加, 物理後, 回転付与
	blurWeightBone.BoneFlag = pmx.BONE_FLAG_CAN_TRANSLATE | pmx.BONE_FLAG_CAN_ROTATE | pmx.BONE_FLAG_CAN_MANIPULATE | pmx.BONE_FLAG_IS_VISIBLE | pmx.BONE_FLAG_IS_AFTER_PHYSICS_DEFORM | pmx.BONE_FLAG_IS_EXTERNAL_ROTATION
	// 表示枠追加
	blurDisplaySlot.References = append(blurDisplaySlot.References,
		pmx.Reference{DisplayType: pmx.DISPLAY_TYPE_BONE, DisplayIndex: blurWeightBone.Index})

	// 付与親
	blurWeightBone.EffectIndex = blurBone.Index
	blurWeightBone.EffectFactor = 1.0

	// 変形階層（親のひとつ後）
	blurWeightBone.Layer = 1

	model.Bones.Append(blurWeightBone)

	return blurWeightBone
}

func createBlurWeightTailBone(
	model *pmx.PmxModel, blurDisplaySlot *pmx.DisplaySlot, blurRootBone *pmx.Bone, blurTailBone *pmx.Bone,
) *pmx.Bone {
	blurWeightTailBone := pmx.NewBone()
	blurWeightTailBone.Index = model.Bones.Len()
	blurWeightTailBone.Name = "ブレW先"
	blurWeightTailBone.EnglishName = "Blur W Tail"
	blurWeightTailBone.ParentIndex = blurRootBone.Index
	blurWeightTailBone.IsSystem = false // システム用ではなく出力用

	blurWeightTailBone.Position = blurTailBone.Position.Copy()
	// 移動可能, 回転可能, 操作可能, 表示枠追加
	blurWeightTailBone.BoneFlag = pmx.BONE_FLAG_CAN_TRANSLATE | pmx.BONE_FLAG_CAN_ROTATE | pmx.BONE_FLAG_CAN_MANIPULATE | pmx.BONE_FLAG_IS_VISIBLE
	// 表示枠追加
	blurDisplaySlot.References = append(blurDisplaySlot.References,
		pmx.Reference{DisplayType: pmx.DISPLAY_TYPE_BONE, DisplayIndex: blurWeightTailBone.Index})

	model.Bones.Append(blurWeightTailBone)

	return blurWeightTailBone
}

func createDiffuseMorph(
	model *pmx.PmxModel, blurMaterials []*pmx.Material, diffuse *mmath.MVec4, diffuseName, diffuseEnglishName string,
) {
	// 非透過モーフ追加
	morph := pmx.NewMorph()
	morph.Index = model.Morphs.Len()
	morph.Name = fmt.Sprintf("ブレ_%s", diffuseName)
	morph.EnglishName = fmt.Sprintf("blur_%s", diffuseEnglishName)

	for _, blurMaterial := range blurMaterials {
		offset := pmx.NewMaterialMorphOffset(
			blurMaterial.Index,
			pmx.CALC_MODE_ADDITION,
			diffuse,
			&mmath.MVec4{0.0, 0.0, 0.0, 0.0},
			&mmath.MVec3{0.0, 0.0, 0.0},
			&mmath.MVec4{0.0, 0.0, 0.0, 0.0},
			0.0,
			&mmath.MVec4{0.0, 0.0, 0.0, 0.0},
			&mmath.MVec4{0.0, 0.0, 0.0, 0.0},
			&mmath.MVec4{0.0, 0.0, 0.0, 0.0},
		)
		morph.Offsets = append(morph.Offsets, offset)
	}

	morph.MorphType = pmx.MORPH_TYPE_MATERIAL
	morph.Panel = pmx.MORPH_PANEL_OTHER_LOWER_RIGHT
	morph.IsSystem = false

	model.DisplaySlots.Get(1).References = append(model.DisplaySlots.Get(1).References,
		pmx.Reference{DisplayType: pmx.DISPLAY_TYPE_MORPH, DisplayIndex: morph.Index})

	model.Morphs.Append(morph)
}

func createTextureMorph(
	model *pmx.PmxModel, blurMaterials []*pmx.Material, textureFactor *mmath.MVec4, name, englishName string,
) {
	// 色モーフ追加
	morph := pmx.NewMorph()
	morph.Index = model.Morphs.Len()
	morph.Name = fmt.Sprintf("ブレ_%s", name)
	morph.EnglishName = fmt.Sprintf("blur_%s", englishName)

	for _, blurMaterial := range blurMaterials {
		offset := pmx.NewMaterialMorphOffset(
			blurMaterial.Index,
			pmx.CALC_MODE_ADDITION,
			&mmath.MVec4{0.0, 0.0, 0.0, 0.0},
			&mmath.MVec4{0.0, 0.0, 0.0, 0.0},
			&mmath.MVec3{0.0, 0.0, 0.0},
			&mmath.MVec4{0.0, 0.0, 0.0, 0.0},
			0.0,
			textureFactor,
			&mmath.MVec4{0.0, 0.0, 0.0, 0.0},
			&mmath.MVec4{0.0, 0.0, 0.0, 0.0},
		)
		morph.Offsets = append(morph.Offsets, offset)
	}

	morph.MorphType = pmx.MORPH_TYPE_MATERIAL
	morph.Panel = pmx.MORPH_PANEL_OTHER_LOWER_RIGHT
	morph.IsSystem = false

	model.DisplaySlots.Get(1).References = append(model.DisplaySlots.Get(1).References,
		pmx.Reference{DisplayType: pmx.DISPLAY_TYPE_MORPH, DisplayIndex: morph.Index})

	model.Morphs.Append(morph)
}

func createSpecularMorph(
	model *pmx.PmxModel, blurMaterials []*pmx.Material, specular *mmath.MVec4, specularName, specularEnglishName string,
) {
	// ALモーフ追加
	morph := pmx.NewMorph()
	morph.Index = model.Morphs.Len()
	morph.Name = fmt.Sprintf("ブレ_%s", specularName)
	morph.EnglishName = fmt.Sprintf("blur_%s", specularEnglishName)

	for _, blurMaterial := range blurMaterials {
		offset := pmx.NewMaterialMorphOffset(
			blurMaterial.Index,
			pmx.CALC_MODE_ADDITION,
			&mmath.MVec4{0.0, 0.0, 0.0, 0.0},
			specular,
			&mmath.MVec3{0.0, 0.0, 0.0},
			&mmath.MVec4{0.0, 0.0, 0.0, 0.0},
			0.0,
			&mmath.MVec4{0.0, 0.0, 0.0, 0.0},
			&mmath.MVec4{0.0, 0.0, 0.0, 0.0},
			&mmath.MVec4{0.0, 0.0, 0.0, 0.0},
		)
		morph.Offsets = append(morph.Offsets, offset)
	}

	morph.MorphType = pmx.MORPH_TYPE_MATERIAL
	morph.Panel = pmx.MORPH_PANEL_OTHER_LOWER_RIGHT
	morph.IsSystem = false

	model.DisplaySlots.Get(1).References = append(model.DisplaySlots.Get(1).References,
		pmx.Reference{DisplayType: pmx.DISPLAY_TYPE_MORPH, DisplayIndex: morph.Index})

	model.Morphs.Append(morph)
}

func createBlurRootRigidBody(model *pmx.PmxModel, blurRootBone *pmx.Bone) *pmx.RigidBody {
	blurRootRigidBody := pmx.NewRigidBody()
	blurRootRigidBody.IsSystem = false
	blurRootRigidBody.Index = model.RigidBodies.Len()
	blurRootRigidBody.Name = blurRootBone.Name
	blurRootRigidBody.EnglishName = blurRootBone.EnglishName
	blurRootRigidBody.BoneIndex = blurRootBone.Index
	blurRootRigidBody.CollisionGroup = 0
	blurRootRigidBody.CollisionGroupMaskValue = 0
	blurRootRigidBody.ShapeType = pmx.SHAPE_SPHERE          // 球剛体
	blurRootRigidBody.PhysicsType = pmx.PHYSICS_TYPE_STATIC // ボーン追従
	blurRootRigidBody.Size = &mmath.MVec3{0.5, 0, 0}
	blurRootRigidBody.Position = blurRootBone.Position.Copy() // ボーン位置

	blurRootRigidBody.RigidBodyParam = pmx.NewRigidBodyParam()
	blurRootRigidBody.RigidBodyParam.Mass = 50.0
	blurRootRigidBody.RigidBodyParam.LinearDamping = 0.5
	blurRootRigidBody.RigidBodyParam.AngularDamping = 0.5
	blurRootRigidBody.RigidBodyParam.Restitution = 0.5
	blurRootRigidBody.RigidBodyParam.Friction = 0.5

	model.RigidBodies.Append(blurRootRigidBody)

	return blurRootRigidBody
}

func createBlurRigidBody(
	model *pmx.PmxModel, blurBone *pmx.Bone, backVertexIndexes []int, edgeVertexIndexes []int,
) *pmx.RigidBody {
	blurRigidBody := pmx.NewRigidBody()
	blurRigidBody.IsSystem = false
	blurRigidBody.Index = model.RigidBodies.Len()
	blurRigidBody.Name = blurBone.Name
	blurRigidBody.EnglishName = blurBone.EnglishName
	blurRigidBody.BoneIndex = blurBone.Index
	blurRigidBody.CollisionGroup = 0
	blurRigidBody.CollisionGroupMaskValue = 0
	blurRigidBody.ShapeType = pmx.SHAPE_BOX              // 箱剛体
	blurRigidBody.PhysicsType = pmx.PHYSICS_TYPE_DYNAMIC // 物理

	backRootVertex := model.Vertices.Get(backVertexIndexes[0])
	edgeRootVertex := model.Vertices.Get(edgeVertexIndexes[0])

	backTailVertex := model.Vertices.Get(backVertexIndexes[len(backVertexIndexes)-1])
	edgeTailVertex := model.Vertices.Get(edgeVertexIndexes[len(edgeVertexIndexes)-1])

	rootPosition := backRootVertex.Position.Added(edgeRootVertex.Position).MuledScalar(0.5)
	tailPosition := backTailVertex.Position.Added(edgeTailVertex.Position).MuledScalar(0.5)

	blurRigidBody.Size = &mmath.MVec3{0.5, 0.5, rootPosition.Distance(tailPosition) * 0.5}
	// 選択の頂点位置の中間
	blurRigidBody.Position = rootPosition.Added(tailPosition).MuledScalar(0.5)
	// 選択の頂点位置の方向
	blurRigidBody.Rotation.SetQuaternion(
		mmath.NewMQuaternionFromDirection(tailPosition.Sub(rootPosition), mmath.MVec3UnitY).Shorten())

	blurRigidBody.RigidBodyParam = pmx.NewRigidBodyParam()
	blurRigidBody.RigidBodyParam.Mass = 5.0
	blurRigidBody.RigidBodyParam.LinearDamping = 0
	blurRigidBody.RigidBodyParam.AngularDamping = 0.1
	blurRigidBody.RigidBodyParam.Restitution = 0
	blurRigidBody.RigidBodyParam.Friction = 0

	model.RigidBodies.Append(blurRigidBody)

	return blurRigidBody
}

func createJoint(
	model *pmx.PmxModel, blurBone *pmx.Bone, blurRootRigidBody, blurRigidBody *pmx.RigidBody,
	rootVector, backVector *mmath.MVec3,
) {
	joint := pmx.NewJoint()
	joint.Index = model.Joints.Len()
	joint.Name = fmt.Sprintf("%s_%s", blurRootRigidBody.Name, blurRigidBody.Name)
	joint.EnglishName = fmt.Sprintf("%s_%s", blurRootRigidBody.EnglishName, blurRigidBody.EnglishName)
	joint.RigidbodyIndexA = blurRootRigidBody.Index
	joint.RigidbodyIndexB = blurRigidBody.Index
	joint.Position = blurRootRigidBody.Position.Copy()

	joint.JointParam = pmx.NewJointParam()

	limitFlag := backVector.Cross(rootVector).Dot(mmath.MVec3UnitY.Cross(rootVector))
	jointRotationLimit := blurBone.FixedAxis.ToMat4().MulVec3(&mmath.MVec3{-20 * mmath.BoolToFlag(limitFlag > 0), 0, 0})
	minRotationLimit := mmath.NewMVec3()
	maxRotationLimit := mmath.NewMVec3()

	if jointRotationLimit.GetX() < 0 {
		minRotationLimit.SetX(jointRotationLimit.GetX())
	} else {
		maxRotationLimit.SetX(jointRotationLimit.GetX())
	}

	if jointRotationLimit.GetY() < 0 {
		minRotationLimit.SetY(jointRotationLimit.GetY())
	} else {
		maxRotationLimit.SetY(jointRotationLimit.GetY())
	}

	if jointRotationLimit.GetZ() < 0 {
		minRotationLimit.SetZ(jointRotationLimit.GetZ())
	} else {
		maxRotationLimit.SetZ(jointRotationLimit.GetZ())
	}

	joint.JointParam.RotationLimitMin.SetDegrees(minRotationLimit)
	joint.JointParam.RotationLimitMax.SetDegrees(maxRotationLimit)

	model.Joints.Append(joint)
}

package usecase_test

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/pmx"
	"github.com/miu200521358/sword_blur/pkg/usecase"
)

func TestSave01(t *testing.T) {
	r := &pmx.PmxReader{}

	data, err := r.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/017_数珠丸恒次/数珠丸恒次 hzeo式/数珠丸恒次（本体）.pmx")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	model := data.(*pmx.PmxModel)

	outputPath := "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/017_数珠丸恒次/数珠丸恒次 hzeo式/数珠丸恒次（本体）_test.pmx"
	blurMaterialIndexes := []int{2}
	backVertexIndexes := []int{143, 1770, 144, 1769, 128, 1766, 120, 1764, 104, 1760, 96, 1758, 88, 1756, 80, 1754,
		72, 1751, 71, 1752, 54, 1748, 46, 1746, 38, 1744, 28, 1741, 19, 1739, 12, 1737, 1, 1736, 0, 1734}
	edgeVertexIndexes := []int{130, 1768, 129, 1767, 113, 1763, 105, 1761, 97, 1759, 89, 1757, 81, 1755, 73, 1753,
		60, 1750, 57, 1749, 51, 1747, 43, 1745, 35, 1743, 29, 1742, 22, 1740, 13, 1738, 5, 1735, 0, 1734}

	err = usecase.Save(model, outputPath, blurMaterialIndexes, backVertexIndexes, edgeVertexIndexes)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	data2, err := r.ReadByFilepath(outputPath)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	outputModel := data2.(*pmx.PmxModel)

	blurMaterial := outputModel.Materials.GetByName("頭身_ブレ")
	if blurMaterial == nil {
		t.Errorf("Expected blurMaterial to be not nil, got nil")
		return
	}

	if blurMaterial.DrawFlag.IsDrawingEdge() {
		t.Errorf("Expected blurMaterial.DrawFlag.IsDrawingEdge() to be false, got true")
	}

	if !blurMaterial.DrawFlag.IsDoubleSidedDrawing() {
		t.Errorf("Expected blurMaterial.DrawFlag.IsDrawingBack() to be true, got false")
	}
}

func TestSave02(t *testing.T) {
	r := &pmx.PmxReader{}

	data, err := r.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/093_陸奥守吉行/陸奥守吉行 むつ式 ver.2.0/むつ式むっちゃん刀.pmx")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	model := data.(*pmx.PmxModel)

	outputPath := "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/093_陸奥守吉行/陸奥守吉行 むつ式 ver.2.0/むつ式むっちゃん刀_test.pmx"
	blurMaterialIndexes := []int{3}
	backVertexIndexes := []int{164, 163, 162, 161, 160, 159, 158, 157, 196, 197}
	edgeVertexIndexes := []int{146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 197}

	err = usecase.Save(model, outputPath, blurMaterialIndexes, backVertexIndexes, edgeVertexIndexes)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	data2, err := r.ReadByFilepath(outputPath)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	outputModel := data2.(*pmx.PmxModel)

	blurMaterial := outputModel.Materials.GetByName("刀身_ブレ")
	if blurMaterial == nil {
		t.Errorf("Expected blurMaterial to be not nil, got nil")
		return
	}

	if blurMaterial.DrawFlag.IsDrawingEdge() {
		t.Errorf("Expected blurMaterial.DrawFlag.IsDrawingEdge() to be false, got true")
	}

	if !blurMaterial.DrawFlag.IsDoubleSidedDrawing() {
		t.Errorf("Expected blurMaterial.DrawFlag.IsDrawingBack() to be true, got false")
	}
}

func TestSave03(t *testing.T) {
	r := &pmx.PmxReader{}

	data, err := r.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/079_江雪左文字/江雪左文字 AKI式 ver.1.51/江雪左文字本体/江雪左文字本体.pmx")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	model := data.(*pmx.PmxModel)

	outputPath := "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/079_江雪左文字/江雪左文字 AKI式 ver.1.51/江雪左文字本体/江雪左文字本体_test.pmx"
	blurMaterialIndexes := []int{3}
	backVertexIndexes := []int{876, 993, 994, 1417, 877, 878, 879, 1003, 1175, 1174, 1172, 1170, 1168, 1162, 1142, 922, 1135, 1120, 1121, 1461, 1124, 1462, 1123, 1126, 1127, 1186, 1163, 1466}
	edgeVertexIndexes := []int{947, 984, 1027, 1419, 1421, 1424, 946, 1420, 1026, 1450, 1076, 1453, 1077, 1454, 941, 1455, 940, 1456, 832, 1390, 831, 1389, 937, 1459, 936, 1460, 935, 134, 933, 1111, 1112, 1116, 1132, 1128, 1122, 1126, 1127, 1186, 1463, 1466}

	err = usecase.Save(model, outputPath, blurMaterialIndexes, backVertexIndexes, edgeVertexIndexes)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	data2, err := r.ReadByFilepath(outputPath)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	outputModel := data2.(*pmx.PmxModel)

	blurMaterial := outputModel.Materials.GetByName("刀身_ブレ")
	if blurMaterial == nil {
		t.Errorf("Expected blurMaterial to be not nil, got nil")
		return
	}

	if blurMaterial.DrawFlag.IsDrawingEdge() {
		t.Errorf("Expected blurMaterial.DrawFlag.IsDrawingEdge() to be false, got true")
	}

	if !blurMaterial.DrawFlag.IsDoubleSidedDrawing() {
		t.Errorf("Expected blurMaterial.DrawFlag.IsDrawingBack() to be true, got false")
	}
}

func TestSave04(t *testing.T) {
	r := &pmx.PmxReader{}

	data, err := r.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/118_へし切長谷部/へし切長谷部 ｻｸﾗｺ式 ver1.20/長谷部本体(鞘無し).pmx")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	model := data.(*pmx.PmxModel)

	outputPath := "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/118_へし切長谷部/へし切長谷部 ｻｸﾗｺ式 ver1.20/長谷部本体(鞘無し)_test.pmx"
	blurMaterialIndexes := []int{0}
	backVertexIndexes := []int{1137, 1128, 1129, 1130, 1131, 1132, 1133, 1134, 1135, 1136, 1140, 1143, 1144, 1145, 980, 1112, 1146, 944, 1104, 1209}
	edgeVertexIndexes := []int{962, 1222, 959, 1219, 958, 1218, 957, 1217, 956, 1216, 955, 1215, 954, 1214, 953, 1213, 952, 1212, 951, 1211, 978, 1225, 987, 1229, 983, 1226, 984, 1227, 965, 1224, 945, 1210, 944, 1104, 1209}

	err = usecase.Save(model, outputPath, blurMaterialIndexes, backVertexIndexes, edgeVertexIndexes)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	data2, err := r.ReadByFilepath(outputPath)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	outputModel := data2.(*pmx.PmxModel)

	blurMaterial := outputModel.Materials.GetByName("刃_ブレ")
	if blurMaterial == nil {
		t.Errorf("Expected blurMaterial to be not nil, got nil")
		return
	}

	if blurMaterial.DrawFlag.IsDrawingEdge() {
		t.Errorf("Expected blurMaterial.DrawFlag.IsDrawingEdge() to be false, got true")
	}

	if !blurMaterial.DrawFlag.IsDoubleSidedDrawing() {
		t.Errorf("Expected blurMaterial.DrawFlag.IsDrawingBack() to be true, got false")
	}
}

func TestSave05(t *testing.T) {
	r := &pmx.PmxReader{}

	data, err := r.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/132_太郎太刀/太郎太刀 AKI式 ver.1.00/太郎太刀本体/太郎太刀本体.pmx")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	model := data.(*pmx.PmxModel)

	outputPath := "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/132_太郎太刀/太郎太刀 AKI式 ver.1.00/太郎太刀本体/太郎太刀本体_test.pmx"
	blurMaterialIndexes := []int{1}
	backVertexIndexes := []int{127, 125, 124, 122, 120, 118, 116, 114, 112, 110, 108, 126, 104, 102, 194}
	edgeVertexIndexes := []int{158, 134, 135, 141, 144, 164, 131, 98, 99, 101, 102, 194, 61, 37, 39, 138, 41, 44, 46, 47, 67, 36, 34, 2, 1, 4, 5, 195}

	err = usecase.Save(model, outputPath, blurMaterialIndexes, backVertexIndexes, edgeVertexIndexes)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	data2, err := r.ReadByFilepath(outputPath)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	outputModel := data2.(*pmx.PmxModel)

	blurMaterial := outputModel.Materials.GetByName("鞘_ブレ")
	if blurMaterial == nil {
		t.Errorf("Expected blurMaterial to be not nil, got nil")
		return
	}

	if blurMaterial.DrawFlag.IsDrawingEdge() {
		t.Errorf("Expected blurMaterial.DrawFlag.IsDrawingEdge() to be false, got true")
	}

	if !blurMaterial.DrawFlag.IsDoubleSidedDrawing() {
		t.Errorf("Expected blurMaterial.DrawFlag.IsDrawingBack() to be true, got false")
	}
}

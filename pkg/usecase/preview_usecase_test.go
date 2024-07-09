package usecase_test

import (
	"fmt"
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

	outputPath := "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/017_数珠丸恒次/数珠丸恒次 hzeo式/数珠丸恒次（本体）_test"
	blurMaterialIndexes := []int{2}
	backRootVertexIndexes := []int{143, 1770}
	edgeTailVertexIndexes := []int{0, 1734}
	edgeVertexIndexes := []int{130, 1768, 129, 1767, 113, 1763, 105, 1761, 97, 1759, 89, 1757, 81, 1755, 73, 1753,
		60, 1750, 57, 1749, 51, 1747, 43, 1745, 35, 1743, 29, 1742, 22, 1740, 13, 1738, 5, 1735, 0, 1734}

	outputModel, previewVmd, err := usecase.Preview(model, blurMaterialIndexes, backRootVertexIndexes, edgeTailVertexIndexes, edgeVertexIndexes)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	usecase.Save(outputModel, outputPath+".pmx")
	previewVmd.Save("Preview", fmt.Sprintf("%s_preview.vmd", outputPath))

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

	outputPath := "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/093_陸奥守吉行/陸奥守吉行 むつ式 ver.2.0/むつ式むっちゃん刀_test"
	blurMaterialIndexes := []int{3}
	backRootVertexIndexes := []int{164}
	edgeTailVertexIndexes := []int{197}
	edgeVertexIndexes := []int{146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 197}

	outputModel, previewVmd, err := usecase.Preview(model, blurMaterialIndexes, backRootVertexIndexes, edgeTailVertexIndexes, edgeVertexIndexes)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	usecase.Save(outputModel, outputPath+".pmx")
	previewVmd.Save("Preview", fmt.Sprintf("%s_preview.vmd", outputPath))

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

	outputPath := "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/079_江雪左文字/江雪左文字 AKI式 ver.1.51/江雪左文字本体/江雪左文字本体_test"
	blurMaterialIndexes := []int{3}
	backRootVertexIndexes := []int{876, 993, 994, 1417}
	edgeTailVertexIndexes := []int{1126, 1127, 1186, 1163, 1466}
	edgeVertexIndexes := []int{947, 984, 1027, 1419, 1421, 1424, 946, 1420, 1026, 1450, 1076, 1453, 1077, 1454, 941, 1455, 940, 1456, 832, 1390, 831, 1389, 937, 1459, 936, 1460, 935, 134, 933, 1111, 1112, 1116, 1132, 1128, 1122, 1126, 1127, 1186, 1463, 1466}

	outputModel, previewVmd, err := usecase.Preview(model, blurMaterialIndexes, backRootVertexIndexes, edgeTailVertexIndexes, edgeVertexIndexes)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	usecase.Save(outputModel, outputPath+".pmx")
	previewVmd.Save("Preview", fmt.Sprintf("%s_preview.vmd", outputPath))

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

	outputPath := "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/118_へし切長谷部/へし切長谷部 ｻｸﾗｺ式 ver1.20/長谷部本体(鞘無し)_test"
	blurMaterialIndexes := []int{0}
	backRootVertexIndexes := []int{1137, 1128, 1129}
	edgeTailVertexIndexes := []int{944, 1104, 1209}
	edgeVertexIndexes := []int{962, 1222, 959, 1219, 958, 1218, 957, 1217, 956, 1216, 955, 1215, 954, 1214, 953, 1213, 952, 1212, 951, 1211, 978, 1225, 987, 1229, 983, 1226, 984, 1227, 965, 1224, 945, 1210, 944, 1104, 1209}

	outputModel, previewVmd, err := usecase.Preview(model, blurMaterialIndexes, backRootVertexIndexes, edgeTailVertexIndexes, edgeVertexIndexes)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	usecase.Save(outputModel, outputPath+".pmx")
	previewVmd.Save("Preview", fmt.Sprintf("%s_preview.vmd", outputPath))

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

	outputPath := "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/132_太郎太刀/太郎太刀 AKI式 ver.1.00/太郎太刀本体/太郎太刀本体_test"
	blurMaterialIndexes := []int{1}
	backRootVertexIndexes := []int{127}
	edgeTailVertexIndexes := []int{194}
	edgeVertexIndexes := []int{158, 134, 135, 141, 144, 164, 131, 98, 99, 101, 102, 194, 61, 37, 39, 138, 41, 44, 46, 47, 67, 36, 34, 2, 1, 4, 5, 195}

	outputModel, previewVmd, err := usecase.Preview(model, blurMaterialIndexes, backRootVertexIndexes, edgeTailVertexIndexes, edgeVertexIndexes)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	usecase.Save(outputModel, outputPath+".pmx")
	previewVmd.Save("Preview", fmt.Sprintf("%s_preview.vmd", outputPath))

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

func TestSave06(t *testing.T) {
	r := &pmx.PmxReader{}

	data, err := r.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/118_へし切長谷部/へし切長谷部 026式 ver.1.57/(3)圧し切長谷部_本体/(3)抜刀ver.1.10.pmx")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	model := data.(*pmx.PmxModel)

	outputPath := "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/118_へし切長谷部/へし切長谷部 026式 ver.1.57/(3)圧し切長谷部_本体/(3)抜刀ver.1.10_test"
	blurMaterialIndexes := []int{0, 1, 2}
	backRootVertexIndexes := []int{1376, 1596}
	edgeTailVertexIndexes := []int{1108, 1138}
	edgeVertexIndexes := []int{1127, 1116, 1126, 1128, 1114, 1118, 1142, 1115, 1117, 1313, 1533, 1294, 1512, 1292, 1512, 1290, 1510, 1288, 1508, 1286, 1506, 1284, 1504, 1282, 1502, 1280, 1500, 1278, 1499, 1276, 1496, 1274, 1495, 1267, 1489, 1266, 1486, 1395, 1611, 1400, 1618, 1264, 1485, 1262, 1483, 1259, 1481, 1258, 1478, 1270, 1490, 1272, 1492, 1297, 1519, 1296, 1516, 1311, 1531, 1412, 1629, 1410, 1627, 1408, 1625, 1406, 1623, 1404, 1621, 1109, 1137, 1108, 1138}

	outputModel, previewVmd, err := usecase.Preview(model, blurMaterialIndexes, backRootVertexIndexes, edgeTailVertexIndexes, edgeVertexIndexes)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	usecase.Save(outputModel, outputPath+".pmx")
	previewVmd.Save("Preview", fmt.Sprintf("%s_preview.vmd", outputPath))

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

func TestSave07(t *testing.T) {
	r := &pmx.PmxReader{}

	data, err := r.ReadByFilepath("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/097_山伏国広/山伏国広ver0_57 ぴえ式/本体.pmx")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	model := data.(*pmx.PmxModel)

	outputPath := "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/097_山伏国広/山伏国広ver0_57 ぴえ式/本体_test"
	blurMaterialIndexes := []int{1}
	backRootVertexIndexes := []int{6, 5}
	edgeTailVertexIndexes := []int{45, 82, 105, 107}
	edgeVertexIndexes := []int{35, 95, 33, 94, 34, 93, 36, 96, 37, 97, 38, 98, 39, 99, 40, 100, 41, 101, 42, 102, 43, 103, 44, 104, 45, 82, 105, 107}

	outputModel, previewVmd, err := usecase.Preview(model, blurMaterialIndexes, backRootVertexIndexes, edgeTailVertexIndexes, edgeVertexIndexes)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}
	usecase.Save(outputModel, outputPath+".pmx")
	previewVmd.Save("Preview", fmt.Sprintf("%s_preview.vmd", outputPath))

	blurMaterial := outputModel.Materials.GetByName("_刀身_ブレ")
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

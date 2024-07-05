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
}

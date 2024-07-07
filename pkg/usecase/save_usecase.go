package usecase

import "github.com/miu200521358/mlib_go/pkg/pmx"

func Save(model *pmx.PmxModel, outputPath string) error {
	return model.Save(false, outputPath)
}

package usecase

import (
	"github.com/miu200521358/sword_blur/pkg/model"
)

func Save(blurModel *model.BlurModel) error {
	err := blurModel.OutputModel.Save(false, "")
	if err != nil {
		return err
	}

	err = blurModel.OutputMotion.Save("Blur Preview", "")
	if err != nil {
		return err
	}

	return nil
}

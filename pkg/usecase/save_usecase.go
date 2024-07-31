package usecase

import (
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/sword_blur/pkg/model"
)

func Save(blurModel *model.BlurModel) error {
	pmxRep := repository.NewPmxRepository()
	if err := pmxRep.Save(blurModel.OutputModel.Path(), blurModel.OutputModel, false); err != nil {
		return err
	}

	vmdRep := repository.NewVmdRepository()
	if err := vmdRep.Save(blurModel.OutputMotion.Path(), blurModel.OutputMotion, false); err != nil {
		return err
	}

	return nil
}

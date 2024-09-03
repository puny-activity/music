package fileserviceuc

import (
	"context"
	"github.com/puny-activity/music/pkg/werr"
)

func (u *UseCase) ReloadClients(ctx context.Context) error {
	err := u.fileServiceController.Reset()
	if err != nil {
		return werr.WrapSE("failed to reload clients", err)
	}

	fileServices, err := u.fileServiceRepository.GetAll(ctx)
	if err != nil {
		return werr.WrapSE("failed to reload clients", err)
	}

	for i := range fileServices {
		err := u.fileServiceController.Add(fileServices[i])
		if err != nil {
			u.log.Warn().Err(err).Msg("failed to reload clients")
		}
	}

	return nil
}

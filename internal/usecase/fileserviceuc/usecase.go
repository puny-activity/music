package fileserviceuc

import (
	"context"
	"github.com/puny-activity/music/internal/entity/remotefile/fileservice"
	"github.com/puny-activity/music/internal/infrastructure/fileserviceclient"
	"github.com/puny-activity/music/pkg/txmanager"
	"github.com/rs/zerolog"
)

type UseCase struct {
	fileServiceRepository fileServiceRepository
	fileServiceController *fileserviceclient.Controller
	txManager             txmanager.Transactor
	log                   *zerolog.Logger
}

func New(fileServiceRepository fileServiceRepository, fileServiceController *fileserviceclient.Controller,
	txManager txmanager.Transactor, log *zerolog.Logger) *UseCase {
	return &UseCase{
		fileServiceRepository: fileServiceRepository,
		fileServiceController: fileServiceController,
		txManager:             txManager,
		log:                   log,
	}
}

type fileServiceRepository interface {
	GetAll(ctx context.Context) ([]fileservice.FileService, error)
}

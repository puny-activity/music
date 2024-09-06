package updatesonguc

import (
	"context"
	"github.com/golang-module/carbon"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/remotefile"
	"github.com/puny-activity/music/internal/entity/remotefile/fileservice"
	"github.com/puny-activity/music/internal/entity/song/album"
	"github.com/puny-activity/music/internal/entity/song/artist"
	"github.com/puny-activity/music/internal/entity/song/genre"
	"github.com/puny-activity/music/pkg/util"
	"github.com/puny-activity/music/pkg/werr"
)

func (u *UseCase) Update(ctx context.Context) error {
	fileServicesInfo, err := u.fileServiceRepository.GetAll(ctx)
	if err != nil {
		return werr.WrapSE("failed to get file services", err)
	}

	for i := range fileServicesInfo {
		err = u.updateForOneFileService(ctx, fileServicesInfo[i])
		if err != nil {
			u.log.Warn().Err(err).Msg("failed to update file service")
		}
	}

	err = u.deleteOrphanedGenres(ctx)
	if err != nil {
		u.log.Warn().Err(err).Msg("failed to delete orphaned genres")
	}

	err = u.deleteOrphanedAlbums(ctx)
	if err != nil {
		u.log.Warn().Err(err).Msg("failed to delete orphaned albums")
	}

	err = u.deleteOrphanedArtists(ctx)
	if err != nil {
		u.log.Warn().Err(err).Msg("failed to delete orphaned artists")
	}

	err = u.setNewCovers(ctx)
	if err != nil {
		u.log.Warn().Err(err).Msg("failed to set new covers")
	}

	return nil
}

func (u *UseCase) updateForOneFileService(ctx context.Context, serviceInfo fileservice.FileService) error {
	fileServiceClient, err := u.fileServiceController.Get(*serviceInfo.ID)
	if err != nil {
		return werr.WrapSE("failed to get file service client", err)
	}

	since := carbon.Parse("0001-01-01")
	if serviceInfo.ScannedAt != nil {
		since = *serviceInfo.ScannedAt
	}
	changedFiles, err := fileServiceClient.GetChangedFiles(ctx, since)
	if err != nil {
		return werr.WrapSE("failed to get file from service", err)
	}

	err = u.txManager.Transaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		err := u.createNewParts(ctx, tx, changedFiles.Updated)
		if err != nil {
			return werr.WrapSE("failed to create new parts", err)
		}

		for _, deletedFile := range changedFiles.Deleted {
			err = u.fileRepository.DeleteTx(ctx, tx, deletedFile.ID)
			if err != nil {
				u.log.Warn().Err(err).Msg("failed to delete file")
				continue
			}
		}

		allFiles, err := u.fileRepository.GetAllTx(ctx, tx)
		if err != nil {
			return werr.WrapSE("failed to get all files", err)
		}
		savedFileIDs := make(map[remotefile.ID]struct{}, len(allFiles))
		for i := range allFiles {
			savedFileIDs[allFiles[i].ID] = struct{}{}
		}

		// TODO: Сделать так, чтобы старые песни могли переиспользовать новые файлы
		for _, updatedFile := range changedFiles.Updated {
			if !updatedFile.ContentType.IsAudio() {
				continue
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *UseCase) createNewParts(ctx context.Context, tx *sqlx.Tx, updatedFiles []remotefile.Updated) error {
	savedGenres, err := u.genreRepository.GetAllTx(ctx, tx)
	if err != nil {
		return werr.WrapSE("failed to get saved genres", err)
	}
	genres := make(map[string]genre.ID)
	genresToCreate := make([]genre.Base, 0)
	for _, genreItem := range savedGenres {
		genres[genreItem.Name] = *genreItem.ID
	}

	savedAlbums, err := u.albumRepository.GetAllTx(ctx, tx)
	if err != nil {
		return werr.WrapSE("failed to get saved albums", err)
	}
	albums := make(map[string]album.ID)
	albumsToCreate := make([]album.Base, 0)
	for _, albumItem := range savedAlbums {
		albums[albumItem.Title] = *albumItem.ID
	}

	savedArtists, err := u.artistRepository.GetAllTx(ctx, tx)
	if err != nil {
		return werr.WrapSE("failed to get saved artists", err)
	}
	artists := make(map[string]artist.ID)
	artistsToCreate := make([]artist.Base, 0)
	for _, artistItem := range savedArtists {
		artists[artistItem.Name] = *artistItem.ID
	}

	for _, updatedFile := range updatedFiles {
		if updatedFile.ContentType.IsAudio() {
			audioMetadata := updatedFile.GetAudioMetadata()
			if audioMetadata.Genre != nil {
				_, ok := genres[*audioMetadata.Genre]
				if !ok {
					newGenre := genre.Base{
						Name: *audioMetadata.Genre,
					}
					newGenre.ID = util.ToPointer(genre.GenerateID())
					genres[newGenre.Name] = *newGenre.ID
					genresToCreate = append(genresToCreate, newGenre)
				}
			}
			if audioMetadata.Album != nil {
				_, ok := albums[*audioMetadata.Album]
				if !ok {
					newAlbum := album.Base{
						Title: *audioMetadata.Album,
					}
					newAlbum.ID = util.ToPointer(album.GenerateID())
					albums[newAlbum.Title] = *newAlbum.ID
					albumsToCreate = append(albumsToCreate, newAlbum)
				}
			}
			if audioMetadata.Artist != nil {
				_, ok := artists[*audioMetadata.Artist]
				if !ok {
					newArtist := artist.Base{
						Name: *audioMetadata.Artist,
					}
					newArtist.ID = util.ToPointer(artist.GenerateID())
					artists[newArtist.Name] = *newArtist.ID
					artistsToCreate = append(artistsToCreate, newArtist)
				}
			}
		} else if updatedFile.ContentType.IsImage() {

		}
	}

	err = u.genreRepository.CreateAllTx(ctx, tx, genresToCreate)
	if err != nil {
		return werr.WrapSE("failed to save genres", err)
	}

	err = u.albumRepository.CreateAllTx(ctx, tx, albumsToCreate)
	if err != nil {
		return werr.WrapSE("failed to save albums", err)
	}

	err = u.artistRepository.CreateAllTx(ctx, tx, artistsToCreate)
	if err != nil {
		return werr.WrapSE("failed to save artists", err)
	}

	return nil
}

func (u *UseCase) deleteOrphanedGenres(ctx context.Context) error {
	// TODO Реализовать удаление жанров
	return nil
}

func (u *UseCase) deleteOrphanedAlbums(ctx context.Context) error {
	// TODO Реализовать удаление альбомов
	return nil
}

func (u *UseCase) deleteOrphanedArtists(ctx context.Context) error {
	// TODO Реализовать удаление исполнителей
	return nil
}

func (u *UseCase) setNewCovers(ctx context.Context) error {
	// TODO
	return nil
}

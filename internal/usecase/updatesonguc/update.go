package updatesonguc

import (
	"context"
	"github.com/jmoiron/sqlx"
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

	return nil
}

func (u *UseCase) updateForOneFileService(ctx context.Context, serviceInfo fileservice.FileService) error {
	fileServiceClient, err := u.fileServiceController.Get(*serviceInfo.ID)
	if err != nil {
		return werr.WrapSE("failed to get file service client", err)
	}

	remoteFiles, err := fileServiceClient.GetAllFiles(ctx)
	if err != nil {
		return werr.WrapSE("failed to get file from service", err)
	}

	err = u.txManager.Transaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		savedGenres, err := u.genreRepository.GetAllTx(ctx, tx)
		if err != nil {
			return werr.WrapSE("failed to get saved genres", err)
		}
		genres := make(map[string]genre.ID)
		genresToCreate := make([]genre.Genre, 0)
		for _, genreItem := range savedGenres {
			genres[genreItem.Name] = *genreItem.ID
		}

		savedAlbums, err := u.albumRepository.GetAllTx(ctx, tx)
		if err != nil {
			return werr.WrapSE("failed to get saved albums", err)
		}
		albums := make(map[string]album.ID)
		albumsToCreate := make([]album.Album, 0)
		for _, albumItem := range savedAlbums {
			albums[albumItem.Title] = *albumItem.ID
		}

		savedArtists, err := u.artistRepository.GetAllTx(ctx, tx)
		if err != nil {
			return werr.WrapSE("failed to get saved artists", err)
		}
		artists := make(map[string]artist.ID)
		artistsToCreate := make([]artist.Artist, 0)
		for _, artistItem := range savedArtists {
			artists[artistItem.Name] = *artistItem.ID
		}

		for _, remoteFile := range remoteFiles {
			if remoteFile.ContentType.IsAudio() {
				audioMetadata := remoteFile.GetAudioMetadata()
				if audioMetadata.Genre != nil {
					_, ok := genres[*audioMetadata.Genre]
					if !ok {
						newGenre := genre.Genre{
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
						newAlbum := album.Album{
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
						newArtist := artist.Artist{
							Name: *audioMetadata.Artist,
						}
						newArtist.ID = util.ToPointer(artist.GenerateID())
						artists[newArtist.Name] = *newArtist.ID
						artistsToCreate = append(artistsToCreate, newArtist)
					}
				}
			} else if remoteFile.ContentType.IsImage() {

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
	})
	if err != nil {
		return err
	}

	return nil
}

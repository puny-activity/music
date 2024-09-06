package updatesonguc

import (
	"context"
	"github.com/golang-module/carbon"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/music/internal/entity/remotefile"
	"github.com/puny-activity/music/internal/entity/remotefile/fileservice"
	"github.com/puny-activity/music/internal/entity/song"
	"github.com/puny-activity/music/internal/entity/song/album"
	"github.com/puny-activity/music/internal/entity/song/artist"
	"github.com/puny-activity/music/internal/entity/song/genre"
	"github.com/puny-activity/music/pkg/util"
	"github.com/puny-activity/music/pkg/werr"
	"time"
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
		err := u.createNewParts(ctx, tx, append(changedFiles.Created, changedFiles.Updated...))
		if err != nil {
			return werr.WrapSE("failed to create new parts for created files", err)
		}

		err = u.fileRepository.DeleteAllTx(ctx, tx, changedFiles.Deleted)
		if err != nil {
			u.log.Warn().Err(err).Msg("failed to delete files")
		}

		err = u.createNewFiles(ctx, tx, serviceInfo, changedFiles.Created)
		if err != nil {
			u.log.Warn().Err(err).Msg("failed to delete files")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *UseCase) createNewParts(ctx context.Context, tx *sqlx.Tx, updatedFiles []remotefile.FileInfo) error {
	savedGenres, err := u.genreRepository.GetAllTx(ctx, tx)
	if err != nil {
		return werr.WrapSE("failed to get saved genres", err)
	}
	genres := make(map[string]genre.Base)
	genresToCreate := make([]genre.Base, 0)
	for _, genreItem := range savedGenres {
		genres[genreItem.Name] = genre.Base{
			ID:   genreItem.ID,
			Name: genreItem.Name,
		}
	}

	savedAlbums, err := u.albumRepository.GetAllTx(ctx, tx)
	if err != nil {
		return werr.WrapSE("failed to get saved albums", err)
	}
	albums := make(map[string]album.Base)
	albumsToCreate := make([]album.Base, 0)
	for _, albumItem := range savedAlbums {
		albums[albumItem.Title] = album.Base{
			ID:    albumItem.ID,
			Title: albumItem.Title,
		}
	}

	savedArtists, err := u.artistRepository.GetAllTx(ctx, tx)
	if err != nil {
		return werr.WrapSE("failed to get saved artists", err)
	}
	artists := make(map[string]artist.Base)
	artistsToCreate := make([]artist.Base, 0)
	for _, artistItem := range savedArtists {
		artists[artistItem.Name] = artist.Base{
			ID:   artistItem.ID,
			Name: artistItem.Name,
		}
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
					genres[newGenre.Name] = genre.Base{
						ID:   newGenre.ID,
						Name: newGenre.Name,
					}
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
					albums[newAlbum.Title] = album.Base{
						ID:    newAlbum.ID,
						Title: newAlbum.Title,
					}
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
					artists[newArtist.Name] = artist.Base{
						ID:   newArtist.ID,
						Name: newArtist.Name,
					}
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

func (u *UseCase) createNewFiles(ctx context.Context, tx *sqlx.Tx, serviceInfo fileservice.FileService, created []remotefile.FileInfo) error {
	savedGenres, err := u.genreRepository.GetAllTx(ctx, tx)
	if err != nil {
		return werr.WrapSE("failed to get saved genres", err)
	}
	genres := make(map[string]genre.Base)
	for _, genreItem := range savedGenres {
		genres[genreItem.Name] = genre.Base{
			ID:   genreItem.ID,
			Name: genreItem.Name,
		}
	}

	savedAlbums, err := u.albumRepository.GetAllTx(ctx, tx)
	if err != nil {
		return werr.WrapSE("failed to get saved albums", err)
	}
	albums := make(map[string]album.Base)
	for _, albumItem := range savedAlbums {
		albums[albumItem.Title] = album.Base{
			ID:    albumItem.ID,
			Title: albumItem.Title,
		}
	}

	savedArtists, err := u.artistRepository.GetAllTx(ctx, tx)
	if err != nil {
		return werr.WrapSE("failed to get saved artists", err)
	}
	artists := make(map[string]artist.Base)
	for _, artistItem := range savedArtists {
		artists[artistItem.Name] = artist.Base{
			ID:   artistItem.ID,
			Name: artistItem.Name,
		}
	}

	filesToCreate := make([]remotefile.File, len(created))
	for i, createdFile := range created {
		filesToCreate[i] = remotefile.File{
			ID:   createdFile.ID,
			Name: createdFile.Name,
			Path: createdFile.Path,
		}
	}
	err = u.fileRepository.CreateAllTx(ctx, tx, *serviceInfo.ID, filesToCreate)
	if err != nil {
		u.log.Warn().Err(err).Msg("failed to update files")
	}

	// TODO: Обрабатывать возможность перемещения файлов отслеживанием через md5

	songsToCreate := make([]song.Song, 0)
	for _, createdFile := range created {
		if createdFile.ContentType.IsImage() {
			// TODO: Создать обложки
		} else if createdFile.ContentType.IsAudio() {
			metadata := createdFile.GetAudioMetadata()
			title := createdFile.Name
			if metadata.Title != nil {
				title = *metadata.Title
			}
			songGenre := genre.DefaultGenre
			if metadata.Genre != nil {
				genreFromMap, ok := genres[*metadata.Genre]
				if ok {
					songGenre = genreFromMap
				}
			}
			songAlbum := album.DefaultAlbum
			if metadata.Album != nil {
				albumFromMap, ok := albums[*metadata.Album]
				if ok {
					songAlbum = albumFromMap
				}
			}
			songArtist := artist.DefaultArtist
			if metadata.Artist != nil {
				artistFromMap, ok := artists[*metadata.Artist]
				if ok {
					songArtist = artistFromMap
				}
			}
			songsToCreate = append(songsToCreate, song.Song{
				ID:           util.ToPointer(song.GenerateID()),
				FileID:       &createdFile.ID,
				Title:        title,
				Duration:     time.Duration(metadata.DurationNs),
				Cover:        nil,
				Genre:        songGenre,
				Album:        songAlbum,
				Artist:       songArtist,
				Year:         metadata.Year,
				Number:       metadata.TrackNumber,
				Comment:      metadata.Comment,
				Channels:     metadata.Channels,
				BitrateKbps:  metadata.BitrateKbps,
				SampleRateHz: metadata.SampleRateHz,
				MD5:          createdFile.MD5,
			})
		}
	}

	err = u.songRepository.CreateAllTx(ctx, tx, songsToCreate)
	if err != nil {
		u.log.Warn().Err(err).Msg("failed to save songs")
	}

	return nil
}

func (u *UseCase) setNewCovers(ctx context.Context) error {
	// TODO
	return nil
}

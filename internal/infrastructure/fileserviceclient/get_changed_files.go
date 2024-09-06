package fileserviceclient

import (
	"context"
	"github.com/golang-module/carbon"
	"github.com/puny-activity/music/internal/entity/remotefile"
	"github.com/puny-activity/music/internal/entity/remotefile/contenttype"
	"github.com/puny-activity/music/pkg/proto/gen/fileserviceproto"
	"github.com/puny-activity/music/pkg/werr"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *Client) GetChangedFiles(ctx context.Context, since carbon.Carbon) (remotefile.Changes, error) {
	clientResponse, err := c.client.GetChanges(ctx, &fileserviceproto.GetChangesRequest{
		Since: timestamppb.New(since.ToStdTime()),
	})
	if err != nil {
		return remotefile.Changes{}, werr.WrapSE("failed to fetch all files", err)
	}

	createdFiles := make([]remotefile.FileInfo, 0)
	for _, file := range clientResponse.CreatedFiles {
		id, err := remotefile.ParseID(file.Id.Id)
		if err != nil {
			c.log.Debug().Msg("failed to parse file id")
			continue
		}
		contentType, err := contenttype.New(file.ContentType)
		if err != nil {
			c.log.Debug().Msg("failed to parse content type")
			continue
		}
		if !contentType.IsAudio() && !contentType.IsImage() {
			continue
		}
		createdFiles = append(createdFiles, remotefile.FileInfo{
			ID:          id,
			Name:        file.Name,
			ContentType: contentType,
			Path:        file.Path,
			Size:        file.Size,
			Metadata:    file.Metadata,
			MD5:         file.Md5,
		})
	}

	updatedFiles := make([]remotefile.FileInfo, 0)
	for _, file := range clientResponse.UpdatedFiles {
		id, err := remotefile.ParseID(file.Id.Id)
		if err != nil {
			c.log.Debug().Msg("failed to parse file id")
			continue
		}
		contentType, err := contenttype.New(file.ContentType)
		if err != nil {
			c.log.Debug().Msg("failed to parse content type")
			continue
		}
		if !contentType.IsAudio() && !contentType.IsImage() {
			continue
		}
		updatedFiles = append(updatedFiles, remotefile.FileInfo{
			ID:          id,
			Name:        file.Name,
			ContentType: contentType,
			Path:        file.Path,
			Size:        file.Size,
			Metadata:    file.Metadata,
			MD5:         file.Md5,
		})
	}

	deletedFiles := make([]remotefile.ID, 0)
	for _, fileIDStr := range clientResponse.DeletedFileIds {
		fileID, err := remotefile.ParseID(fileIDStr.Id)
		if err != nil {
			continue
		}
		deletedFiles = append(deletedFiles, fileID)
	}

	return remotefile.Changes{
		Created: createdFiles,
		Updated: updatedFiles,
		Deleted: deletedFiles,
	}, nil
}

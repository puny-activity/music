package fileserviceclient

import (
	"context"
	"github.com/golang-module/carbon"
	"github.com/google/uuid"
	"github.com/puny-activity/music/internal/entity/remotefile"
	"github.com/puny-activity/music/internal/entity/remotefile/contenttype"
	"github.com/puny-activity/music/pkg/proto/gen/fileserviceproto"
	"github.com/puny-activity/music/pkg/werr"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *Client) GetChangedFiles(ctx context.Context, since carbon.Carbon) (remotefile.Changed, error) {
	clientResponse, err := c.client.GetChanged(ctx, &fileserviceproto.GetChangedRequest{
		Since: timestamppb.New(since.ToStdTime()),
	})
	if err != nil {
		return remotefile.Changed{}, werr.WrapSE("failed to fetch all files", err)
	}

	updatedFiles := make([]remotefile.Updated, 0)
	for _, file := range clientResponse.UpdatedFiles {
		id, err := uuid.Parse(file.Id)
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

		updatedFiles = append(updatedFiles, remotefile.Updated{
			ID:          remotefile.NewID(id),
			Name:        file.Name,
			ContentType: contentType,
			Path:        file.Path,
			Size:        file.Size,
			Metadata:    file.Metadata,
			MD5:         file.Md5,
		})
	}

	deletedFiles := make([]remotefile.Deleted, 0)
	for _, fileIDStr := range clientResponse.DeletedFileIds {
		fileID, err := uuid.Parse(fileIDStr)
		if err != nil {
			continue
		}
		deletedFiles = append(deletedFiles, remotefile.Deleted{
			ID: remotefile.NewID(fileID),
		})
	}

	return remotefile.Changed{
		Updated: updatedFiles,
		Deleted: deletedFiles,
	}, nil
}

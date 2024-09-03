package fileserviceclient

import (
	"context"
	"github.com/google/uuid"
	"github.com/puny-activity/music/internal/entity/remotefile"
	"github.com/puny-activity/music/internal/entity/remotefile/contenttype"
	"github.com/puny-activity/music/pkg/proto/gen/fileserviceproto"
	"github.com/puny-activity/music/pkg/werr"
)

func (c *Client) GetAllFiles(ctx context.Context) ([]remotefile.Full, error) {
	clientResponse, err := c.client.GetAllFiles(ctx, &fileserviceproto.GetAllFilesRequest{})
	if err != nil {
		return nil, werr.WrapSE("failed to get all files", err)
	}

	fullFiles := make([]remotefile.Full, 0)
	for _, file := range clientResponse.Files {
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

		fullFiles = append(fullFiles, remotefile.Full{
			ID:          remotefile.NewID(id),
			Name:        file.Name,
			ContentType: contentType,
			Path:        file.Path,
			Size:        file.Size,
			Metadata:    file.Metadata,
			MD5:         file.Md5,
		})
	}

	return fullFiles, nil
}

package fileserviceclient

import (
	"github.com/puny-activity/music/pkg/proto/gen/fileserviceproto"
	"github.com/puny-activity/music/pkg/werr"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	address string
	conn    *grpc.ClientConn
	client  fileserviceproto.FileServiceClient
	log     *zerolog.Logger
}

func New(address string, log *zerolog.Logger) *Client {
	return &Client{
		address: address,
		log:     log,
	}
}

func (c *Client) Start() error {
	var err error

	creds := insecure.NewCredentials()
	c.conn, err = grpc.NewClient(c.address, grpc.WithTransportCredentials(creds))
	if err != nil {
		return werr.WrapSE("failed to create grpc client", err)
	}

	client := fileserviceproto.NewFileServiceClient(c.conn)
	c.client = client

	return nil
}

func (c *Client) Stop() error {
	return c.conn.Close()
}

package fileserviceclient

import (
	"errors"
	"github.com/puny-activity/music/internal/entity/remotefile/fileservice"
	"github.com/puny-activity/music/pkg/werr"
	"github.com/rs/zerolog"
)

type Controller struct {
	clients map[fileservice.ID]*Client
	log     *zerolog.Logger
}

func NewController(log *zerolog.Logger) *Controller {
	return &Controller{
		clients: make(map[fileservice.ID]*Client),
		log:     log,
	}
}

func (c *Controller) Add(fileService fileservice.FileService) error {
	newClient := New(fileService.GRPCAddress, c.log)
	err := newClient.Start()
	if err != nil {
		return werr.WrapSE("failed to start client", err)
	}
	c.clients[*fileService.ID] = newClient
	return nil
}

func (c *Controller) Get(id fileservice.ID) (*Client, error) {
	client, ok := c.clients[id]
	if !ok {
		return nil, errors.New("client not found")
	}
	return client, nil
}

func (c *Controller) Remove(fileServiceID fileservice.ID) {
	delete(c.clients, fileServiceID)
}

func (c *Controller) Reset() error {
	for i := range c.clients {
		err := c.clients[i].Stop()
		if err != nil {
			c.log.Error().Err(err).Msg("failed to stop client")
		}
	}
	c.clients = make(map[fileservice.ID]*Client)
	return nil
}

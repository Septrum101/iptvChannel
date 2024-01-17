package zteg

import (
	"github.com/thank243/iptvChannel/api"
	"github.com/thank243/iptvChannel/config"
)

func New(conf *config.Config) *Client {
	return new(Client)
}

func (c *Client) GetEPGs(id int) ([]api.Epg, error) {
	return []api.Epg{}, nil
}

func (c *Client) GetChannels() ([]api.Channel, error) {
	return []api.Channel{}, nil
}

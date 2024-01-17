package controller

import (
	"github.com/robfig/cron/v3"

	"github.com/thank243/iptvChannel/api"
	"github.com/thank243/iptvChannel/app/server"
	"github.com/thank243/iptvChannel/config"
)

type Controller struct {
	conf          *config.Config
	cli           api.Client
	server        *server.Server
	cron          *cron.Cron
	maxConcurrent int
}

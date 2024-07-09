package controller

import (
	"github.com/robfig/cron/v3"

	"github.com/Septrum101/iptvChannel/api"
	"github.com/Septrum101/iptvChannel/app/server"
	"github.com/Septrum101/iptvChannel/config"
)

type Controller struct {
	conf          *config.Config
	cli           api.Client
	server        *server.Server
	cron          *cron.Cron
	maxConcurrent int
}

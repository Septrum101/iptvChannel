package server

import (
	"sync/atomic"

	"github.com/labstack/echo/v4"

	"github.com/thank243/iptvChannel/api"
)

type Server struct {
	Echo     *echo.Echo
	EPGs     *atomic.Pointer[[]api.Epg]
	Channels *atomic.Pointer[[]api.Channel]

	mode      string
	udpxyHost string
}

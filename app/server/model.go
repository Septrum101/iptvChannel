package server

import (
	"sync/atomic"

	"github.com/labstack/echo/v4"

	"github.com/thank243/iptvChannel/common/channel"
	"github.com/thank243/iptvChannel/common/epg"
)

type Server struct {
	Echo     *echo.Echo
	EPGs     *atomic.Pointer[[]epg.Epg]
	Channels *atomic.Pointer[[]channel.Channel]

	udpxyHost string
}

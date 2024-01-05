package handler

import (
	"sync/atomic"

	"github.com/labstack/echo/v4"

	"github.com/thank243/iptvChannel/common/channel"
	"github.com/thank243/iptvChannel/common/epg"
	"github.com/thank243/iptvChannel/common/req"
)

type Handler struct {
	Echo     *echo.Echo
	EPGs     *atomic.Pointer[[]epg.Epg]
	Channels *atomic.Pointer[[]channel.Channel]

	req       *req.Req
	udpxyHost string
}

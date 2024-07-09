package server

import (
	"sync"
	"sync/atomic"

	"github.com/labstack/echo/v4"

	"github.com/Septrum101/iptvChannel/api"
)

type Server struct {
	Echo     *echo.Echo
	EPGs     *atomic.Pointer[[]api.Epg]
	Channels *atomic.Pointer[[]api.Channel]
	DiypEPGs *sync.Map

	mode      string
	udpxyHost string
}

type DiypEPG struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Title string `json:"title"`
}

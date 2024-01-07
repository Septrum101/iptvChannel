package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync/atomic"

	"github.com/beevik/etree"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/thank243/iptvChannel/common/channel"
	"github.com/thank243/iptvChannel/common/epg"
	"github.com/thank243/iptvChannel/common/req"
	"github.com/thank243/iptvChannel/config"
)

// New creates a new Handler instance and configures it based on the provided config.
// It sets up the Echo instance, RequestLogger middleware, and the route handlers.
// It returns the created Handler instance.
func New(c *config.Config) *Handler {
	handler := &Handler{
		Echo:      echo.New(),
		req:       req.New(c),
		udpxyHost: c.UdpxyHost,
		Channels:  new(atomic.Pointer[[]channel.Channel]),
		EPGs:      new(atomic.Pointer[[]epg.Epg]),
	}

	handler.Echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogError:     true,
		LogRemoteIP:  true,
		LogMethod:    true,
		LogUserAgent: true,
		HandleError:  true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				log.WithFields(log.Fields{
					"RemoteIP":  v.RemoteIP,
					"Method":    v.Method,
					"URI":       v.URI,
					"UserAgent": v.UserAgent,
					"status":    v.Status,
				}).Info("request")
			} else {
				log.WithFields(log.Fields{
					"RemoteIP":  v.RemoteIP,
					"Method":    v.Method,
					"URI":       v.URI,
					"UserAgent": v.UserAgent,
					"status":    v.Status,
					"error":     v.Error,
				}).Error("request error")
			}
			return nil
		},
	}))

	level, err := log.ParseLevel(c.LogLevel)
	if err != nil {
		log.Panic(err)
	}
	log.SetLevel(level)
	if level == log.DebugLevel || level == log.TraceLevel {
		log.SetReportCaller(true)
	}

	g := handler.Echo.Group("/api/v1")
	g.GET("/getChannels", handler.getChannels)
	g.GET("/getEpgs", handler.getEPGs)

	return handler
}

func (h *Handler) getChannels(c echo.Context) error {
	if h.Channels.Load() == nil {
		return c.String(http.StatusServiceUnavailable, "no valid channels")
	}

	channels := *h.Channels.Load()

	b := bytes.Buffer{}
	b.WriteString("#EXTM3U\n")

	for _, ch := range channels {
		name := ch.ChannelName
		addr, err := url.Parse(ch.ChannelURL)
		if err != nil {
			continue
		}

		b.WriteString(fmt.Sprintf("#EXTINF:-1, tvg-id=\"%d\" tvg-name=\"%s\", %s\n", ch.ChannelID, name, name))
		b.WriteString(fmt.Sprintf("%s/rtp/%s\n", h.udpxyHost, addr.Host))
	}

	return c.Blob(http.StatusOK, "text/plain;charset=UTF-8", b.Bytes())
}

func (h *Handler) getEPGs(c echo.Context) error {
	if h.Channels.Load() == nil {
		return c.String(http.StatusServiceUnavailable, "no valid channels")
	}

	channels := *h.Channels.Load()

	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	tv := doc.CreateElement("tv")

	// add channel to doc
	for i := range channels {
		ch := channels[i]
		// create channel, format: <channel id="1"><display-name lang="zh">CCTV1</display-name></channel>
		channelXml := tv.CreateElement("channel")
		channelXml.CreateAttr("id", strconv.Itoa(ch.ChannelID))
		name := channelXml.CreateElement("display-name")
		name.CreateAttr("lang", "zh")
		name.CreateText(ch.ChannelName)
	}

	// add EPGs to doc
	// create programme, format:
	// <programme start="20240103215500 +0800" stop="20240103232500 +0800" channel="7249">
	// <title lang="zh">实况录像</title><desc lang="zh"></desc></programme>

	validEpg := h.EPGs.Load()
	if validEpg != nil {
		for i := range *validEpg {
			e := (*validEpg)[i]
			programmeXml := tv.CreateElement("programme")
			programmeXml.CreateAttr("start", fmt.Sprintf("%s +0800", e.BeginTimeFormat))
			programmeXml.CreateAttr("stop", fmt.Sprintf("%s +0800", e.EndTimeFormat))
			programmeXml.CreateAttr("channel", e.ChannelId)
			title := programmeXml.CreateElement("title")
			title.CreateAttr("lang", "zh")
			title.CreateText(e.ProgramName)
		}
	}

	b, _ := doc.WriteToBytes()
	return c.Blob(http.StatusOK, "text/xml", b)
}

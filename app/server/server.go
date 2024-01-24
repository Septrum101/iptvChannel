package server

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"sync/atomic"

	"github.com/beevik/etree"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/thank243/iptvChannel/api"
	"github.com/thank243/iptvChannel/config"
)

// New creates a new Server instance and configures it based on the provided config.
// It sets up the Echo instance, RequestLogger middleware, and the route handlers.
// It returns the created Server instance.
func New(c *config.Config) *Server {
	s := &Server{
		Echo:     echo.New(),
		Channels: new(atomic.Pointer[[]api.Channel]),
		EPGs:     new(atomic.Pointer[[]api.Epg]),

		mode:      c.Mode,
		udpxyHost: c.UdpxyHost,
	}

	s.Echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
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
					"remote_ip":  v.RemoteIP,
					"method":     v.Method,
					"URI":        v.URI,
					"user_agent": v.UserAgent,
					"status":     v.Status,
				}).Info("request")
			} else {
				log.WithFields(log.Fields{
					"remote_ip":  v.RemoteIP,
					"method":     v.Method,
					"URI":        v.URI,
					"user_agent": v.UserAgent,
					"status":     v.Status,
					"error":      v.Error,
				}).Error("request error")
			}
			return nil
		},
	}),
		middleware.Recover(),
		middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}),
	)

	g := s.Echo.Group("/api/v1")
	g.GET("/getChannels", s.getChannels)
	g.GET("/getEpgs", s.getEPGs)

	return s
}

func (s *Server) getChannels(c echo.Context) error {
	if s.Channels.Load() == nil {
		return c.String(http.StatusServiceUnavailable, "no valid channels")
	}

	channels := *s.Channels.Load()

	b := bytes.Buffer{}
	b.WriteString("#EXTM3U\n")

	for _, ch := range channels {
		name := ch.ChannelName
		addr, err := s.buildChannelUrl(&ch)
		if err != nil {
			logger := log.WithFields(log.Fields{
				"ChannelName": ch.ChannelName,
			})
			logger.Debug(err)
			continue
		}

		b.WriteString(fmt.Sprintf("#EXTINF:-1, tvg-id=\"%s\" tvg-name=\"%s\", %s\n", ch.ChannelID, name, name))
		b.WriteString(fmt.Sprintf("%s\n", addr))
	}

	return c.Blob(http.StatusOK, "text/plain;charset=UTF-8", b.Bytes())
}

func (s *Server) getEPGs(c echo.Context) error {
	if s.Channels.Load() == nil {
		return c.String(http.StatusServiceUnavailable, "no valid channels")
	}

	channels := *s.Channels.Load()

	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	tv := doc.CreateElement("tv")
	tv.CreateAttr("ext-info", config.GetVersion())

	// add channel to doc
	for i := range channels {
		ch := channels[i]
		// create channel, format: <channel id="1"><display-name lang="zh">CCTV1</display-name></channel>
		channelXml := tv.CreateElement("channel")
		channelXml.CreateAttr("id", ch.ChannelID)
		name := channelXml.CreateElement("display-name")
		name.CreateAttr("lang", "zh")
		name.CreateText(ch.ChannelName)
	}

	// add EPGs to doc
	// create programme, format:
	// <programme start="20240103215500 +0800" stop="20240103232500 +0800" channel="7249">
	// <title lang="zh">实况录像</title><desc lang="zh"></desc></programme>

	validEpg := s.EPGs.Load()
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

func (s *Server) buildChannelUrl(ch *api.Channel) (string, error) {
	switch s.mode {
	case "UDPXY":
		addr, err := url.Parse(ch.ChannelURL)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s/rtp/%s", s.udpxyHost, addr.Host), nil
	case "IGMP":
		addr, err := url.Parse(ch.ChannelURL)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("rtp://%s", addr.Host), nil
	case "URL":
		return ch.TimeShiftURL, nil
	default:
		return "", fmt.Errorf("unsupported mode: %s", s.mode)
	}
}

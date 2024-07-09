package server

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/beevik/etree"
	"github.com/bitly/go-simplejson"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/Septrum101/iptvChannel/api"
	"github.com/Septrum101/iptvChannel/config"
)

// New creates a new Server instance and configures it based on the provided config.
// It sets up the Echo instance, RequestLogger middleware, and the route handlers.
// It returns the created Server instance.
func New(c *config.Config) *Server {
	s := &Server{
		Echo:     echo.New(),
		Channels: new(atomic.Pointer[[]api.Channel]),
		EPGs:     new(atomic.Pointer[[]api.Epg]),
		DiypEPGs: new(sync.Map),

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
	var b []byte

	switch c.QueryParam("type") {
	case strings.ToLower("diyp"):
		b = s.buildDiyp(channels)
	default:
		b = s.buildM3U(channels)
	}

	return c.Blob(http.StatusOK, "text/plain;charset=UTF-8", b)
}

func (s *Server) buildM3U(channels []api.Channel) []byte {
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
	return b.Bytes()
}

func (s *Server) buildDiyp(channels []api.Channel) []byte {
	var (
		b         bytes.Buffer
		cctv      []api.Channel
		satellite []api.Channel
		local     []api.Channel
	)

	for _, ch := range channels {
		name := ch.ChannelName
		switch {
		case strings.Contains(name, "CCTV"):
			cctv = append(cctv, ch)
		case strings.Contains(name, "卫视"):
			satellite = append(satellite, ch)
		default:
			local = append(local, ch)
		}
	}

	// cctv
	b.Write(s.buildDiypData(cctv, 0))
	// satellite
	b.Write(s.buildDiypData(satellite, 1))
	// local
	b.Write(s.buildDiypData(local, 2))

	return b.Bytes()
}

// buildDiypData builds a string representation of DIYP data for the given channels and channel type.
// It starts by writing the appropriate header based on the channel type (0 for "央视", 1 for "卫视", and default for "本地").
func (s *Server) buildDiypData(channels []api.Channel, channelType int) []byte {
	var b bytes.Buffer

	switch channelType {
	case 0:
		b.WriteString("央视,#genre#\n")
	case 1:
		b.WriteString("卫视,#genre#\n")
	default:
		b.WriteString("本地,#genre#\n")
	}

	for _, ch := range channels {
		addr, err := s.buildChannelUrl(&ch)
		if err != nil {
			logger := log.WithFields(log.Fields{
				"ChannelName": ch.ChannelName,
			})
			logger.Debug(err)
			continue
		}
		b.WriteString(fmt.Sprintf("%s,%s\n", ch.ChannelName, addr))
	}
	b.WriteString("\n")

	return b.Bytes()
}

func (s *Server) getEPGs(c echo.Context) error {
	if s.Channels.Load() == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "no valid channels")
	}

	channels := *s.Channels.Load()

	if c.QueryParam("ch") != "" {
		qDate := c.QueryParam("date")
		if qDate == "" {
			qDate = time.Now().Format("2006-01-02")
		}

		val, ok := s.DiypEPGs.Load(c.QueryParam("ch"))
		if !ok {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid channel name")
		}

		diypMaps := val.(map[string][]DiypEPG)
		diypJs := simplejson.New()
		diypJs.Set("date", qDate)
		diypJs.Set("channel_name", c.QueryParam("ch"))
		diypJs.Set("epg_data", diypMaps[qDate])

		return c.JSON(http.StatusOK, diypJs)
	}

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
	// fix channel url
	if strings.Contains(ch.ChannelURL, "|") {
		ch.ChannelURL = strings.SplitN(ch.ChannelURL, "|", 2)[0]
	}

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

func (s *Server) GetChannelNameFromID(id string) string {
	channels := *s.Channels.Load()
	for _, ch := range channels {
		if ch.ChannelID == id {
			return ch.ChannelName
		}
	}

	return ""
}

package controller

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"

	"github.com/Septrum101/iptvChannel/api"
	"github.com/Septrum101/iptvChannel/api/hwtc"
	"github.com/Septrum101/iptvChannel/api/zteg"
	"github.com/Septrum101/iptvChannel/app/server"
	"github.com/Septrum101/iptvChannel/config"
)

func New(c *config.Config) (*Controller, error) {
	// set log level
	level, err := log.ParseLevel(c.LogLevel)
	if err != nil {
		log.Panic(err)
	}
	log.SetLevel(level)
	if level == log.DebugLevel || level == log.TraceLevel {
		log.SetReportCaller(true)
	}

	// set provide mode
	c.Mode = strings.ToUpper(c.Mode)
	switch c.Mode {
	case "UDPXY":
		if c.UdpxyHost == "" {
			return nil, errors.New("udpxy host is null")
		}
	case "IGMP", "URL":
	default:
		return nil, fmt.Errorf("unsupported mode: %s", c.Mode)
	}

	ctrl := &Controller{
		conf:          c,
		server:        server.New(c),
		cron:          cron.New(),
		maxConcurrent: c.MaxConcurrent,
	}

	// set api provider
	switch strings.ToLower(c.Api.Provider) {
	case "hwtc":
		ctrl.cli = hwtc.New(c)
	case "zteg":
		ctrl.cli = zteg.New(c)
	default:
		return nil, fmt.Errorf("unsupported mode: %s", c.Api.Provider)
	}

	// check max concurrent
	if c.MaxConcurrent > 16 {
		ctrl.maxConcurrent = 16
	}

	// set cron job skip if still running
	if _, err := ctrl.cron.AddJob(c.Cron, cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger)).Then(ctrl)); err != nil {
		return nil, err
	}

	return ctrl, nil
}

func (c *Controller) Start() error {
	fmt.Printf("%s\nLogLevel: %s, MaxConcurrent: %d, Mode: %s, Provider: %s\n",
		config.GetVersion(), c.conf.LogLevel, c.maxConcurrent, strings.ToUpper(c.conf.Mode), c.conf.Api.Provider)

	log.Info("Starting service..")
	log.Info("Fetch Channels and EPGs data on initial startup")
	c.Run()
	time.Sleep(time.Second)

	// start cron job
	c.cron.Start()

	// start http server
	if err := c.server.Echo.Start(c.conf.Address); err != nil {
		return err
	}
	return nil
}

func (c *Controller) Stop() error {
	log.Info("Closing service..")
	if err := c.server.Echo.Shutdown(c.cron.Stop()); err != nil {
		return err
	}

	return nil
}

// Run fetch EPGs and Channels
func (c *Controller) Run() {
	if err := c.fetchChannels(); err != nil {
		log.Error(err)
	}

	if err := c.fetchEPGs(); err != nil {
		log.Error(err)
	}
}

func (c *Controller) fetchChannels() error {
	log.Info("Fetch Channels")

	channels, err := c.cli.GetChannels()
	if err != nil {
		return err
	}

	// filet valid channels
	if c.conf.Exclude != "" {
		re := regexp.MustCompile(c.conf.Exclude)
		var filteredChannels []api.Channel
		for _, ch := range channels {
			if !re.MatchString(ch.ChannelName) {
				filteredChannels = append(filteredChannels, ch)
			}
		}
		channels = filteredChannels
	}

	c.server.Channels.Store(&channels)
	log.Infof("Get channels: %d", len(channels))

	return nil
}

func (c *Controller) fetchEPGs() error {
	log.Info("Fetch EPGs")

	if c.server.Channels.Load() == nil {
		log.Info("Channels is null, fetch channels first")
		if err := c.fetchChannels(); err != nil {
			return err
		}
	}

	channels := *c.server.Channels.Load()

	var epgChan = make(chan api.Epg)
	var wg sync.WaitGroup

	sem := make(chan bool, c.maxConcurrent) // This is used to limit the number of goroutines to maxConcurrent
	for i := range channels {
		wg.Add(1)

		go func(i int) {
			defer func() {
				<-sem // leave semaphore
				wg.Done()
			}()
			sem <- true // enter semaphore, will block if there are maxConcurrent tasks running already

			ch := channels[i]
			logger := log.WithFields(log.Fields{
				"ChannelId":   ch.ChannelID,
				"ChannelName": ch.ChannelName,
			})
			logger.Debug("start get EPGs")

			epgs, err := c.cli.GetEPGs(ch.ChannelID)
			if err != nil {
				logger.Debug(err)
				return
			}
			for i := range epgs {
				epgChan <- epgs[i]
			}
		}(i)
	}

	// Close the channel after all work has been done
	go func() {
		wg.Wait()
		close(sem)
		close(epgChan)
	}()

	// Consume results from the channel and append to slice
	var epgList []api.Epg
	diypEPGMaps := make(map[string]map[string][]server.DiypEPG)
	for e := range epgChan {
		chName := c.server.GetChannelNameFromID(e.ChannelId)
		if chName != "" {
			start := toDiypTimeFormat(e.BeginTimeFormat)
			end := toDiypTimeFormat(e.EndTimeFormat)
			programDay := start.Format("2006-01-02")

			if diyp, ok := diypEPGMaps[chName]; !ok {
				diypEPGMaps[chName] = make(map[string][]server.DiypEPG)
			} else {
				diyp[programDay] = append(diyp[programDay], server.DiypEPG{
					Start: start.Format("15:04"),
					End:   end.Format("15:04"),
					Title: e.ProgramName,
				})
			}
		}

		epgList = append(epgList, e)
	}

	c.server.EPGs.Store(&epgList)

	for chName := range diypEPGMaps {
		c.server.DiypEPGs.Store(chName, diypEPGMaps[chName])
	}

	log.Infof("Get EPGs: %d", len(epgList))

	return nil
}

// convert time	20240103215500 to 21:55
func toDiypTimeFormat(t string) *time.Time {
	parseTime, _ := time.Parse("20060102150405", t)
	return &parseTime
}

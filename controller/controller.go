package controller

import (
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"

	"github.com/thank243/iptvChannel/app/server"
	"github.com/thank243/iptvChannel/common/channel"
	"github.com/thank243/iptvChannel/common/epg"
	"github.com/thank243/iptvChannel/common/req"
	"github.com/thank243/iptvChannel/config"
)

func New(c *config.Config) (*Controller, error) {
	ctrl := &Controller{
		conf:          c,
		req:           req.New(c),
		server:        server.New(c),
		cron:          cron.New(),
		maxConcurrent: c.MaxConcurrent,
	}
	if c.MaxConcurrent > 16 {
		ctrl.maxConcurrent = 16
	}

	_, err := ctrl.cron.AddJob(c.Cron, cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger)).Then(ctrl))
	if err != nil {
		return nil, err
	}

	return ctrl, nil
}

func (c *Controller) Start() error {
	fmt.Printf("%s\nLogLevel: %s, MaxConcurrent: %d\n", config.GetVersion(), c.conf.LogLevel, c.maxConcurrent)
	log.Info("Starting service..")

	log.Info("Fetch EPGs and Channels data on initial startup")
	c.Run()

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
		return
	}
}

func (c *Controller) fetchChannels() error {
	log.Info("Fetch Channels")

	buf, err := c.req.GetChannelBytes()
	if err != nil {
		return err
	}

	channels, err := channel.BytesToChannels(buf)
	if err != nil {
		return err
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

	var epgChan = make(chan epg.Epg)
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
			logger := log.WithField("channelId", ch.ChannelID)
			logger.Debug("start get EPGs")
			resp, err := c.req.GetEPGBytes(ch.ChannelID)
			if err != nil {
				logger.Error(err)
				return
			}
			epgs, err := epg.BytesToValidEPGs(resp)
			if err != nil {
				logger.Error(err)
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
		close(epgChan)
	}()

	// Consume results from the channel and append to slice
	var esSlice []epg.Epg
	for e := range epgChan {
		esSlice = append(esSlice, e)
	}

	c.server.EPGs.Store(&esSlice)
	log.Infof("Get EPGs: %d", len(esSlice))

	return nil
}

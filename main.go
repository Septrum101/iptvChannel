package main

import (
	"errors"
	"net/http"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"github.com/thank243/iptvChannel/config"
	"github.com/thank243/iptvChannel/controller"
)

func main() {
	c, err := controller.New(config.ReadConfig())
	if err != nil {
		log.Panic(err)
	}

	go func() {
		if err := c.Start(); !errors.Is(err, http.ErrServerClosed) {
			log.Panic(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	if err := c.Stop(); err != nil {
		log.Panic(err)
	}
}

package controller

import (
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/Septrum101/iptvChannel/config"
)

func TestFetch(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	ctrl, err := New(config.ReadConfig())
	if err != nil {
		t.Fatal(err)
	}

	if err := ctrl.fetchEPGs(); err != nil {
		t.Fatal(err)
	}

	es := ctrl.server.EPGs.Load()
	for i := range *es {
		t.Log((*es)[i])
	}
}

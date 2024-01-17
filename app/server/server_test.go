package server

import (
	"testing"

	"github.com/thank243/iptvChannel/api"
	"github.com/thank243/iptvChannel/config"
)

func TestHandlerGetEPGs(t *testing.T) {
	s := New(config.ReadConfig())
	s.Channels.Store(&[]api.Channel{{ChannelID: 3954}})

	if err := s.getEPGs(nil); err != nil {
		t.Error(err)
	}
}

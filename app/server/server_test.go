package server

import (
	"testing"

	"github.com/thank243/iptvChannel/common/channel"
	"github.com/thank243/iptvChannel/config"
)

func TestHandlerGetEPGs(t *testing.T) {
	s := New(config.ReadConfig())
	s.Channels.Store(&[]channel.Channel{{ChannelID: 3954}})

	if err := s.getEPGs(nil); err != nil {
		t.Error(err)
	}
}

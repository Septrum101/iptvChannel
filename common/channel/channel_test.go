package channel

import (
	"testing"

	"github.com/thank243/iptvChannel/common/req"
	"github.com/thank243/iptvChannel/config"
)

func TestGetChannels(t *testing.T) {
	c := config.ReadConfig()
	r := req.New(c)

	resp, err := r.GetChannelBytes()
	if err != nil {
		t.Error(err)
	}

	channels, err := BytesToChannels(resp)
	if err != nil {
		t.Log(err)
	}

	for _, channel := range channels {
		t.Log(channel)
	}
}

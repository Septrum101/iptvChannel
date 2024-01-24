package hwtc

import (
	"testing"

	"github.com/thank243/iptvChannel/config"
)

func TestGetChannels(t *testing.T) {
	c := config.ReadConfig()
	r := New(c)

	resp, err := r.getChannelBytes()
	if err != nil {
		t.Error(err)
	}

	// resp, err := os.ReadFile("channel.bin")
	// if err != nil {
	// 	t.Error(err)
	// }

	channels, err := bytesToChannels(resp)
	if err != nil {
		t.Log(err)
	}

	for _, channel := range channels {
		t.Log(channel)
	}
}

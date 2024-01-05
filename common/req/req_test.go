package req

import (
	"testing"

	"github.com/thank243/iptvChannel/config"
)

func TestGetCookie(t *testing.T) {
	c := config.ReadConfig()
	r := New(c)

	t.Log(r.updateCookie())
}

func TestGetEPG(t *testing.T) {
	c := config.ReadConfig()
	r := New(c)

	_, err := r.GetEPGBytes(655980640)
	if err != nil {
		t.Error(err)
	}

}

func TestGetChannel(t *testing.T) {
	c := config.ReadConfig()
	r := New(c)

	_, err := r.GetChannelBytes()
	if err != nil {
		t.Error(err)
	}
}

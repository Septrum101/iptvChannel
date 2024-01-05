package handler

import (
	"testing"

	"github.com/thank243/iptvChannel/config"
)

func TestGetEpgs(t *testing.T) {
	c := config.ReadConfig()
	h := New(c)

	if err := h.getEPGs(nil); err != nil {
		t.Error(err)
	}
}

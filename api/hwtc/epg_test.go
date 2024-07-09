package hwtc

import (
	"testing"

	"github.com/Septrum101/iptvChannel/config"
)

func TestGetEpg(t *testing.T) {
	c := config.ReadConfig()
	r := New(c)

	resp, err := r.getEPGBytes("3954")
	if err != nil {
		t.Error(err)
	}

	validEs, err := bytesToValidEPGs(resp)
	if err != nil {
		t.Error(err)
	}
	for i := range validEs {
		t.Log(validEs[i])
	}
	t.Logf("Valid EPGs: %d", len(validEs))
}

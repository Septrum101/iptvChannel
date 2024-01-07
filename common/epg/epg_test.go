package epg

import (
	"testing"

	"github.com/thank243/iptvChannel/common/req"
	"github.com/thank243/iptvChannel/config"
)

func TestGetEpg(t *testing.T) {
	c := config.ReadConfig()
	r := req.New(c)

	resp, err := r.GetEPGBytes(3954)
	if err != nil {
		t.Error(err)
	}

	validEs, err := BytesToValidEPGs(resp)
	if err != nil {
		t.Error(err)
	}
	for i := range validEs {
		t.Log(validEs[i])
	}
	t.Logf("Valid EPGs: %d", len(validEs))
}

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
	es, err := BytesToAllEPGs(resp)
	if err != nil {
		t.Error(err)
	}
	for i := range es {
		t.Log(es[i])
	}
	t.Logf("Total EPGs: %d", len(es))

	validEs, err := BytesToValidEPGs(resp)
	if err != nil {
		t.Error(err)
	}
	for i := range validEs {
		t.Log(validEs[i])
	}
	t.Logf("Valid EPGs: %d", len(validEs))
}

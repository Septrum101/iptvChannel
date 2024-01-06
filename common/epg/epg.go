package epg

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/thank243/iptvChannel/infra"
)

func BytesToAllEPGs(resp []byte) ([]Epg, error) {
	re := regexp.MustCompile(`(?s)\[.*]`)
	b := re.FindSubmatch(resp)
	if b == nil {
		return nil, errors.New("not found valid data")
	}

	var es []any
	if err := json.Unmarshal(b[0], &es); err != nil {
		return nil, err
	}

	if len(es) != 2 {
		return nil, errors.New("the length of data must equal 2")
	}

	var data [][]Epg
	if err := mapstructure.Decode(es[1], &data); err != nil {
		return nil, err
	}

	var epgs []Epg
	for i := range data {
		e := data[i]
		for ii := range e {
			epgs = append(epgs, e[ii])
		}
	}

	return epgs, nil
}

func BytesToValidEPGs(resp []byte) ([]Epg, error) {
	allEPGs, err := BytesToAllEPGs(resp)
	if err != nil {
		return nil, err
	}

	var epgs []Epg
	tz := time.FixedZone("CST", 8*60*60)
	for i := range allEPGs {
		if err := allEPGs[i].filterValidEPG(tz); err == nil {
			epgs = append(epgs, allEPGs[i])
		}
	}

	return epgs, nil
}

func (e *Epg) filterValidEPG(tz *time.Location) error {
	// time format: 20231228001700
	endTime, err := infra.StrToTime(e.EndTimeFormat, tz)
	if err != nil {
		return err
	}

	beginTime, err := infra.StrToTime(e.BeginTimeFormat, tz)
	if err != nil {
		return err
	}

	if beginTime.Sub(*endTime) > 0 {
		*endTime = endTime.AddDate(0, 0, 1)
		e.EndTimeFormat = endTime.Format("20060102150405")
	}

	if time.Since(*endTime) > time.Hour {
		return fmt.Errorf("not a valid EPG: %s [%s] -> %s", e.ChannelId, e.ProgramName, e.EndTimeFormat)

	}

	return nil
}

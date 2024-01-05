package epg

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/mitchellh/mapstructure"
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
	tz := time.FixedZone("CST", 3600*8)
	for i := range allEPGs {
		epg, err := filterEPGs(&allEPGs[i], tz)
		if err != nil {
			continue
		}

		epgs = append(epgs, *epg)
	}

	return epgs, nil
}

func filterEPGs(e *Epg, tz *time.Location) (*Epg, error) {
	// time format: 20231228001700
	endTime, err := time.ParseInLocation("20060102150405", e.EndTimeFormat, tz)
	if err != nil {
		return nil, err
	}

	beginTime, err := time.ParseInLocation("20060102150405", e.BeginTimeFormat, tz)
	if err != nil {
		return nil, err
	}

	if beginTime.Sub(endTime) > 0 {
		endTime = beginTime
	}

	if time.Since(endTime) > time.Hour {
		return nil, fmt.Errorf("not a valid EPG: %s [%s] -> %s", e.ChannelId, e.ProgramName, e.EndTimeFormat)

	}

	return e, nil
}

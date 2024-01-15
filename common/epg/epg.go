package epg

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

func BytesToValidEPGs(resp []byte) ([]Epg, error) {
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

	var validEPGs []Epg
	tz := time.FixedZone("CST", 8*60*60)
	for i := range data {
		e := data[i]
		for ii := range e {
			if err := e[ii].filterValidEPG(tz); err != nil {
				continue
			}
			validEPGs = append(validEPGs, e[ii])
		}
	}

	return validEPGs, nil
}

func (e *Epg) filterValidEPG(tz *time.Location) error {
	endTime, err := e.fixEndTime(tz)
	if err != nil {
		return err
	}

	if time.Since(endTime) > time.Hour {
		return fmt.Errorf("not a valid EPG: %s [%s] -> %s", e.ChannelId, e.ProgramName, e.EndTimeFormat)

	}

	// fix char 65533 (Replacement Character)
	if strings.Contains(e.ProgramName, string(rune(65533))) {
		e.ProgramName = strings.ReplaceAll(e.ProgramName, string(rune(65533)), "")
	}

	return nil
}

func (e *Epg) fixEndTime(tz *time.Location) (time.Time, error) {
	// time format: 20231228001700
	endTime, err := strToTime(e.EndTimeFormat, tz)
	if err != nil {
		return time.Time{}, err
	}

	beginTime, err := strToTime(e.BeginTimeFormat, tz)
	if err != nil {
		return time.Time{}, err
	}

	if beginTime.Sub(endTime) > 0 {
		endTime = endTime.AddDate(0, 0, 1)
		e.EndTimeFormat = endTime.Format("20060102150405")
	}

	return endTime, nil
}

func strToTime(t string, tz *time.Location) (time.Time, error) {
	toTime, err := time.ParseInLocation("20060102150405", t, tz)
	if err != nil {
		return time.Time{}, err
	}

	return toTime, nil
}

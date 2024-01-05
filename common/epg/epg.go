package epg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
)

func byteToEpg(resp []byte) ([]Epg, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return nil, err
	}

	vm := otto.New()
	if err := vm.Set("parent", map[string]any{}); err != nil {
		return nil, err
	}

	if _, err := vm.Run(doc.Find("script").Text()); err != nil {
		return nil, err
	}

	objParent, _ := vm.Object("parent")
	obj, _ := objParent.Get("jsonBackLookStr")

	ex, _ := obj.Export()
	val, ok := ex.([]any)
	if !ok {
		return nil, fmt.Errorf("invalid type of object")
	}
	if len(val) < 2 {
		return nil, errors.New("invalid length of object")
	}
	es := val[1]

	bb, _ := json.Marshal(es)
	var data [][]Epg
	if err := json.Unmarshal(bb, &data); err != nil {
		return nil, err
	}

	var epgs []Epg
	for i := range data {
		epgsDay := data[i]
		for ii := range epgsDay {
			epgs = append(epgs, epgsDay[ii])
		}
	}

	return epgs, nil
}

func GetEPGs(resp []byte) ([]Epg, error) {
	epgs, err := byteToEpg(resp)
	if err != nil {
		return nil, err
	}

	var validEpgs []Epg
	tz := time.FixedZone("CST", 3600*8)
	for i := range epgs {
		endTime, err := time.ParseInLocation("20060102150405", epgs[i].EndTimeFormat, tz) // time format: 20231228001700
		if err != nil {
			continue
		}

		beginTime, err := time.ParseInLocation("20060102150405", epgs[i].BeginTimeFormat, tz) // time format: 20231228001700
		if err != nil {
			continue
		}

		if beginTime.Sub(endTime) > 0 {
			endTime = beginTime
		}

		if time.Since(endTime) < time.Hour {
			validEpgs = append(validEpgs, epgs[i])
		}
	}
	return validEpgs, nil
}

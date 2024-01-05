package channel

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mitchellh/mapstructure"
)

func BytesToChannels(resp []byte) ([]Channel, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(resp))
	if err != nil {
		return nil, err
	}

	var data []string
	re := regexp.MustCompile(`\((.*?)\)`)  // filter strings in ()
	re1 := regexp.MustCompile(`'([^']*)'`) // filter strings in ''

	// use the goquery document...
	_ = doc.Find("script:contains(ChannelName)").Each(func(i int, selection *goquery.Selection) {
		matches := re.FindStringSubmatch(selection.Text())
		if len(matches) > 1 {
			matches2 := re1.FindAllStringSubmatch(matches[1], -1)
			data = append(data, matches2[1][1]) // 用1而非0作为索引
		}
	})

	var channelMaps []map[string]any
	re2 := regexp.MustCompile(`画中画|单音轨`)
	for i := range data {
		if re2.MatchString(data[i]) {
			continue
		}

		res := strings.Split(data[i], ",")
		kvMap := make(map[string]any)
		for ii := range res {
			kvs := strings.SplitN(res[ii], "=", 2)
			val := strings.Trim(kvs[1], "\"")

			if kvs[0] == "ChannelID" {
				channelID, _ := strconv.Atoi(val)
				kvMap[kvs[0]] = channelID
			} else {
				kvMap[kvs[0]] = val
			}

		}
		channelMaps = append(channelMaps, kvMap)
	}

	var channels []Channel
	if err := mapstructure.Decode(&channelMaps, &channels); err != nil {
		return nil, err
	}

	return channels, nil
}

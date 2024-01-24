package hwtc

import (
	"regexp"
	"strings"

	"github.com/mitchellh/mapstructure"
)

func bytesToChannels(resp []byte) ([]Channel, error) {
	re := regexp.MustCompile(`'ChannelID=.*?'`)
	data := re.FindAll(resp, -1)

	var channelMaps []map[string]string
	re2 := regexp.MustCompile(`画中画|单音轨`)
	for i := range data {
		if re2.Match(data[i]) {
			continue
		}

		res := strings.Split(strings.ReplaceAll(strings.Trim(string(data[i]), "'"), "\"", ""), ",")
		kvMap := make(map[string]string)
		for ii := range res {
			kvs := strings.SplitN(res[ii], "=", 2)
			kvMap[kvs[0]] = kvs[1]
		}
		channelMaps = append(channelMaps, kvMap)
	}

	var channels []Channel
	if err := mapstructure.Decode(&channelMaps, &channels); err != nil {
		return nil, err
	}

	return channels, nil
}

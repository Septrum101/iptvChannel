package zteg

import (
	"regexp"
	"strings"

	"github.com/mitchellh/mapstructure"
)

func bytesToChannels(resp []byte) ([]Channel, error) {
	re := regexp.MustCompile(`'ChannelID=.*?'`)
	data := re.FindAllString(string(resp), -1)

	var channelMaps []map[string]string
	for i := range data {
		res := strings.Split(strings.Trim(strings.ReplaceAll(data[i], "\"", ""), "'"), ",")
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

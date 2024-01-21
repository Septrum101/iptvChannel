package zteg

import (
	"regexp"
	"strings"

	"github.com/mitchellh/mapstructure"
)

func bytesToChannels(resp []byte) ([]Channel, error) {
	re := regexp.MustCompile(`ChannelID="\w+".*?ChannelFCCPort="\d+"`)
	data := re.FindAll(resp, -1)

	var channelMaps []map[string]string
	re2 := regexp.MustCompile(`PIP`)
	for i := range data {
		if re2.Match(data[i]) {
			continue
		}

		d := data[i]
		res := strings.Split(string(d), ",")
		kvMap := make(map[string]string)
		for ii := range res {
			kvs := strings.SplitN(res[ii], "=", 2)
			val := strings.Trim(kvs[1], "\"")
			kvMap[kvs[0]] = val
		}
		channelMaps = append(channelMaps, kvMap)
	}

	var channels []Channel
	if err := mapstructure.Decode(&channelMaps, &channels); err != nil {
		return nil, err
	}

	return channels, nil
}

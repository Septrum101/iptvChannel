package channel

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
)

func BytesToChannels(resp []byte) ([]Channel, error) {
	re := regexp.MustCompile(`(?s)ChannelID="\d*".*?ChannelFECPort="\d*"`)
	data := re.FindAllString(string(resp), -1)

	var channelMaps []map[string]any
	re2 := regexp.MustCompile(`画中画|单音轨`)
	for i := range data {
		if re2.MatchString(data[i]) {
			continue
		}

		d := data[i]
		res := strings.Split(d, ",")
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

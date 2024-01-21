package zteg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/thank243/iptvChannel/api"
	"github.com/thank243/iptvChannel/config"
)

func New(conf *config.Config) *Client {
	r := &Client{
		cli:           resty.New().SetRetryCount(3).SetBaseURL(fmt.Sprintf("%s/iptvepg", conf.Api.ApiHost)),
		userId:        conf.Api.Auth["userid"],
		authenticator: conf.Api.Auth["authenticator"],
	}

	return r
}

// GetEPGs todo
func (c *Client) GetEPGs(id string) ([]api.Epg, error) {
	return []api.Epg{}, nil
}

func (c *Client) getChannelBytes() ([]byte, error) {
	for i := 0; i < 3; i++ {
		resp, err := c.cli.R().SetFormData(map[string]string{
			"MAIN_WIN_SRC":    "/iptvepg/empty.jsp",
			"NEED_UPDATE_STB": "1",
			"BUILD_ACTION":    "FRAMESET_BUILDER",
		}).Post("function/frameset_builder.jsp")
		if err != nil {
			return nil, err
		}

		if strings.Contains(resp.String(), "resignon") {
			time.Sleep(time.Second * 3)
			if err := c.updateCookie(); err != nil {
				return nil, err
			}
			continue
		}

		// convert gbk to utf-8
		buf, err := io.ReadAll(transform.NewReader(bytes.NewReader(resp.Body()), simplifiedchinese.GBK.NewDecoder()))
		if err != nil {
			return nil, err
		}
		return buf, nil
	}
	return nil, errors.New("retry after 3 times")
}

func (c *Client) GetChannels() ([]api.Channel, error) {
	buf, err := c.getChannelBytes()
	if err != nil {
		return nil, err
	}

	chs, err := bytesToChannels(buf)
	if err != nil {
		return nil, err
	}

	var channels []api.Channel
	for i := range chs {
		ch := chs[i]
		channels = append(channels, api.Channel{
			ChannelID:    ch.ChannelID,
			ChannelName:  ch.ChannelName,
			ChannelURL:   ch.ChannelURL,
			TimeShiftURL: ch.TimeShiftURL,
		})
	}

	return channels, nil
}

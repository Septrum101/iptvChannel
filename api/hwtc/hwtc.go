package hwtc

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/thank243/iptvChannel/api"
	"github.com/thank243/iptvChannel/config"
)

func New(conf *config.Config) *Client {
	r := &Client{
		cli:           resty.New().SetRetryCount(3).SetBaseURL(fmt.Sprintf("%s/EPG/jsp", conf.Api.ApiHost)),
		userId:        conf.Api.Auth["userid"],
		authenticator: conf.Api.Auth["authenticator"],
		epgPath:       conf.Api.EPGPath,
	}

	return r
}

func (c *Client) getEPGBytes(channelId string) ([]byte, error) {
	var buf []byte
	for i := 0; i < 3; i++ {
		resp, err := c.cli.R().ForceContentType("text/html;charset=UTF-8").SetQueryParam("channelId", channelId).
			Get(c.epgPath)
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
		buf = resp.Body()
		break
	}
	return buf, nil
}

func (c *Client) GetEPGs(id string) ([]api.Epg, error) {
	buf, err := c.getEPGBytes(id)
	if err != nil {
		return nil, err
	}

	epgs, err := bytesToValidEPGs(buf)
	if err != nil {
		return nil, err
	}

	var es []api.Epg
	for i := range epgs {
		e := epgs[i]
		es = append(es, api.Epg{
			ChannelId:       e.ChannelId,
			BeginTimeFormat: e.BeginTimeFormat,
			EndTimeFormat:   e.EndTimeFormat,
			ProgramName:     e.ProgramName,
		})
	}

	return es, nil
}

func (c *Client) getChannelBytes() ([]byte, error) {
	for i := 0; i < 3; i++ {
		resp, err := c.cli.R().SetFormData(map[string]string{
			"UserToken": c.userToken,
			"UserID":    c.userId,
		}).Post("getchannellistHWCTC.jsp")
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
		return resp.Body(), nil
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

		// fix channel url
		if strings.Contains(ch.ChannelURL, "|") {
			ch.ChannelURL = strings.SplitN(ch.ChannelURL, "|", 2)[0]
		}

		channels = append(channels, api.Channel{
			ChannelID:    ch.ChannelID,
			ChannelName:  ch.ChannelName,
			ChannelURL:   ch.ChannelURL,
			TimeShiftURL: ch.TimeShiftURL,
		})
	}

	return channels, nil
}

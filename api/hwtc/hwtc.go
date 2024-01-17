package hwtc

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/thank243/iptvChannel/api"
	"github.com/thank243/iptvChannel/config"
)

func New(conf *config.Config) *Client {
	r := &Client{
		cli: resty.New().SetRetryCount(3).SetBaseURL(fmt.Sprintf("%s/EPG/jsp", conf.Api.ApiHost)),

		userId:        conf.Api.Auth["userid"],
		authenticator: conf.Api.Auth["authenticator"],
	}

	return r
}

func (c *Client) getEPGBytes(channelId int) ([]byte, error) {
	var buf []byte
	for i := 0; i < 3; i++ {
		resp, err := c.cli.R().ForceContentType("text/html;charset=UTF-8").SetQueryParam("channelId", strconv.Itoa(channelId)).
			Get("stliveplay_30/en/getTvodData.jsp")
		if err != nil {
			return nil, err
		}

		if strings.Contains(resp.String(), "(\"resignon\",\"1\")") {
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

func (c *Client) GetEPGs(id int) ([]api.Epg, error) {
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
	var buf []byte

	for i := 0; i < 3; i++ {
		resp, err := c.cli.R().
			Get("getchannellistHWCTC.jsp")
		if err != nil {
			return nil, err
		}

		if strings.Contains(resp.String(), "(\"resignon\",\"1\")") {
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

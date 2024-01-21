package zteg

import (
	"github.com/go-resty/resty/v2"
)

type Client struct {
	cli *resty.Client

	userId        string
	authenticator string
}

type Channel struct {
	BeginTime            string `mapstructure:"BeginTime"`
	ChannelPurchased     string `mapstructure:"ChannelPurchased"`
	LocalTimeShift       string `mapstructure:"LocalTimeShift"`
	UserTeamChannelID    string `mapstructure:"UserTeamChannelID"`
	ChannelFCCServerAddr string `mapstructure:"ChannelFCCServerAddr"`
	ChannelFCCIP         string `mapstructure:"ChannelFCCIP"`
	ChannelFCCPort       string `mapstructure:"ChannelFCCPort"`
	ChannelID            string `mapstructure:"ChannelID"`
	ChannelLogURL        string `mapstructure:"ChannelLogURL"`
	ChannelName          string `mapstructure:"ChannelName"`
	UserChannelID        string `mapstructure:"UserChannelID"`
	ChannelSDP           string `mapstructure:"ChannelSDP"`
	ChannelType          string `mapstructure:"ChannelType"`
	ChannelURL           string `mapstructure:"ChannelURL"`
	Interval             string `mapstructure:"Interval"`
	Lasting              string `mapstructure:"Lasting"`
	PositionX            string `mapstructure:"PositionX"`
	PositionY            string `mapstructure:"PositionY"`
	TimeShift            string `mapstructure:"TimeShift"`
	TimeShiftURL         string `mapstructure:"TimeShiftURL"`
}

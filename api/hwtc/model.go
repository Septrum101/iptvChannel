package hwtc

import (
	"github.com/go-resty/resty/v2"
)

type Client struct {
	cli *resty.Client

	userId        string
	authenticator string
}

type Epg struct {
	BeginTimeFormat string `json:"beginTimeFormat" mapstructure:"beginTimeFormat"`
	ChannelId       string `json:"channelId" mapstructure:"channelId"`
	ContentId       string `json:"contentId" mapstructure:"contentId"`
	EndTime         string `json:"endTime" mapstructure:"endTime"`
	EndTimeFormat   string `json:"endTimeFormat" mapstructure:"endTimeFormat"`
	Index           string `json:"index" mapstructure:"index"`
	IsPlayable      string `json:"isPlayable" mapstructure:"isPlayable"`
	ProgramName     string `json:"programName" mapstructure:"programName"`
	StartTime       string `json:"startTime" mapstructure:"startTime"`
}

type Channel struct {
	ActionType       string `mapstructure:"ActionType"`
	BeginTime        string `mapstructure:"BeginTime"`
	ChannelFCCIP     string `mapstructure:"ChannelFCCIP"`
	ChannelFCCPort   string `mapstructure:"ChannelFCCPort"`
	ChannelFECPort   string `mapstructure:"ChannelFECPort"`
	ChannelID        int    `mapstructure:"ChannelID"`
	ChannelLocked    string `mapstructure:"ChannelLocked"`
	ChannelLogURL    string `mapstructure:"ChannelLogURL"`
	ChannelName      string `mapstructure:"ChannelName"`
	ChannelPurchased string `mapstructure:"ChannelPurchased"`
	ChannelSDP       string `mapstructure:"ChannelSDP"`
	ChannelType      string `mapstructure:"ChannelType"`
	ChannelURL       string `mapstructure:"ChannelURL"`
	FCCEnable        string `mapstructure:"FCCEnable"`
	Interval         string `mapstructure:"Interval"`
	IsHDChannel      string `mapstructure:"IsHDChannel"`
	Lasting          string `mapstructure:"Lasting"`
	PositionX        string `mapstructure:"PositionX"`
	PositionY        string `mapstructure:"PositionY"`
	PreviewEnable    string `mapstructure:"PreviewEnable"`
	TimeShift        string `mapstructure:"TimeShift"`
	TimeShiftLength  string `mapstructure:"TimeShiftLength"`
	TimeShiftURL     string `mapstructure:"TimeShiftURL"`
	UserChannelID    string `mapstructure:"UserChannelID"`
}

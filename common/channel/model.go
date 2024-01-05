package channel

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

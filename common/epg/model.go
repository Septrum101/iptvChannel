package epg

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

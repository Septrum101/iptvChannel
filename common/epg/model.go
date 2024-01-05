package epg

type Epg struct {
	BeginTimeFormat string `json:"beginTimeFormat"`
	ChannelId       string `json:"channelId"`
	ContentId       string `json:"contentId"`
	EndTime         string `json:"endTime"`
	EndTimeFormat   string `json:"endTimeFormat"`
	Index           string `json:"index"`
	IsPlayable      string `json:"isPlayable"`
	ProgramName     string `json:"programName"`
	StartTime       string `json:"startTime"`
}

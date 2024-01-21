package api

type Client interface {
	GetChannels() ([]Channel, error)
	GetEPGs(id string) ([]Epg, error)
}

type Channel struct {
	ChannelID    string
	ChannelName  string
	ChannelURL   string
	TimeShiftURL string
}

type Epg struct {
	ChannelId       string
	BeginTimeFormat string
	EndTimeFormat   string
	ProgramName     string
}

package config

type Config struct {
	LogLevel      string `yaml:"LogLevel"`
	Cron          string `yaml:"Cron"`
	MaxConcurrent int    `yaml:"MaxConcurrent"`

	// Req
	ApiHost       string `yaml:"ApiHost"`
	UserID        string `yaml:"UserID"`
	Authenticator string `yaml:"Authenticator"`

	// Controller
	Mode      string `yaml:"Mode"`
	Address   string `yaml:"Address"`
	UdpxyHost string `yaml:"UdpxyHost"`
}

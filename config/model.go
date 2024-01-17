package config

type Config struct {
	LogLevel      string `yaml:"LogLevel"`
	Cron          string `yaml:"Cron"`
	MaxConcurrent int    `yaml:"MaxConcurrent"`

	// Client
	Api api `yaml:"Api"`

	// Controller
	Mode      string `yaml:"Mode"`
	Address   string `yaml:"Address"`
	UdpxyHost string `yaml:"UdpxyHost"`
}

type api struct {
	Provider string            `yaml:"Provider"`
	ApiHost  string            `yaml:"ApiHost"`
	Auth     map[string]string `yaml:"Auth"`
}

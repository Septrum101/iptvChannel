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
	Address   string `yaml:"Address"`
	UdpxyHost string `yaml:"UdpxyHost"`
}

var LogLevel = map[string]uint8{
	"DEBUG": 1,
	"INFO":  2,
	"WARN":  3,
	"ERROR": 4,
	"OFF":   5,
}

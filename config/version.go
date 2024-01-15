package config

import (
	"fmt"
)

const (
	appName = "iptvChannel"
	desc    = "Sources: https://www.github.com/thank243/iptvChannel"
)

var (
	version = "dev"
	date    = "unknown"
)

func GetVersion() string {
	return fmt.Sprintf("%s %s, built at %s\n%s", appName, version, date, desc)
}

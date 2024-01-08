package config

import (
	"fmt"
)

const (
	appName = "iptvChannel"
	version = "0.0.3"
	desc    = "(github.com/thank243/iptvChannel)"
)

func GetVersion() string {
	return fmt.Sprintf("%s %s %s", appName, version, desc)
}

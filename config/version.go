package config

import (
	"fmt"
)

const (
	appName = "iptvChannel"
	version = "0.0.3"
	desc    = "Sources: https://www.github.com/thank243/iptvChannel"
)

var commit = "dev"

func GetVersion() string {
	return fmt.Sprintf("%s %s-%s \n%s", appName, version, commit, desc)
}

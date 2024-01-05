package config

import (
	"fmt"
)

const (
	appName = "iptvChannel"
	version = "0.0.2"
	desc    = "(github.com/thank243/iptvChannel)"
)

func ShowVersion() {
	fmt.Println(appName, version, desc)
}

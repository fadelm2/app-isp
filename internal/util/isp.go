package util

import (
	"os"
	"strings"
)

// GetISPName returns the ISP Name, defaulting to "Greenet" if not set in environment.
func GetISPName() string {
	name := os.Getenv("ISP_NAME")
	if name == "" {
		name = "Greenet"
	}
	return name
}

// GetISPNameUpper returns the ISP Name in uppercase, defaulting to "GREENET" if not set in environment.
func GetISPNameUpper() string {
	name := os.Getenv("ISP_NAME")
	if name == "" {
		name = "GREENET"
	}
	return strings.ToUpper(name)
}

package cli

import (
	
)

const (
)

func ParseBoolFromString(s string) bool {

	switch s {
	case "1":
		return true
	default:
		return false
	}
}

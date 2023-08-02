package cli

import (
	flag "github.com/spf13/pflag"
	"strings"
)

const (
	FlagWhitelistLiquidation = "add-liquidation-whitelisting"
)

func ParseBoolFromString(s string) bool {
	switch s {
	case "1":
		return true
	default:
		return false
	}
}

//func ParseUint64SliceFromString(s string, separator string) ([]uint64, error) {
//	var parsedInts []uint64
//	for _, s := range strings.Split(s, separator) {
//		s = strings.TrimSpace(s)
//
//		parsed, err := strconv.ParseUint(s, 10, 64)
//		if err != nil {
//			return []uint64{}, err
//		}
//		parsedInts = append(parsedInts, parsed)
//	}
//	return parsedInts, nil
//}

func ParseStringFromString(s string, separator string) ([]string, error) {
	var parsedStrings []string
	for _, s := range strings.Split(s, separator) {
		s = strings.TrimSpace(s)

		parsedStrings = append(parsedStrings, s)
	}
	return parsedStrings, nil
}

func FlagSetWhitelistLiquidation() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagWhitelistLiquidation, "", "add liquidation whitelisting json file path")
	return fs
}

type createWhitelistLiquidationInputs struct {
	AppId                  string `json:"app_id"`
	Initiator              string `json:"initiator"`
	IsDutchActivated       string `json:"is_dutch_activated"`
	Premium                string `json:"premium"`
	Discount               string `json:"discount"`
	DecrementFactor        string `json:"decrement_factor"`
	IsEnglishActivated     string `json:"is_english_activated"`
	DecrementFactorEnglish string `json:"decrement_factor_english"`
	Title                  string
	Description            string
	Deposit                string
}

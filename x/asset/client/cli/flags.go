package cli

import (
	flag "github.com/spf13/pflag"
	"strconv"
	"strings"
)

const (
	flagName                = "name"
	flagDenom               = "denom"
	flagDecimals            = "decimals"
	FlagAddAssetMappingFile = "add-asset-mapping-file"
)

func ParseStringFromString(s string, separator string) ([]string, error) {
	var parsedStrings []string
	for _, s := range strings.Split(s, separator) {
		s = strings.TrimSpace(s)

		parsedStrings = append(parsedStrings, s)
	}
	return parsedStrings, nil
}

func ParseBoolFromString(s string) bool {
	switch s {
	case "1":
		return true
	default:
		return false
	}
}

func ParseInt64SliceFromString(s string, separator string) ([]int64, error) {
	var parsedInts []int64
	for _, s := range strings.Split(s, separator) {
		s = strings.TrimSpace(s)

		parsed, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return []int64{}, err
		}
		parsedInts = append(parsedInts, parsed)
	}
	return parsedInts, nil
}

func ParseUint64SliceFromString(s string, separator string) ([]uint64, error) {
	var parsedInts []uint64
	for _, s := range strings.Split(s, separator) {
		s = strings.TrimSpace(s)

		parsed, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return []uint64{}, err
		}
		parsedInts = append(parsedInts, parsed)
	}
	return parsedInts, nil
}

func FlagSetCreateAssetMapping() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagAddAssetMappingFile, "", "add asset mapping json file path")
	return fs
}


type createAddAssetMappingInputs struct {
	AppID         string `json:"app_id"`
	AssetID       string `json:"asset_id"`
	GenesisSupply string `json:"genesis_supply"`
	IsGovToken    string `json:"is_gov_token"`
	Recipient     string `json:"recipient"`
	Title         string
	Description   string
	Deposit       string
}

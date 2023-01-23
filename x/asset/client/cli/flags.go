package cli

import (
	"strings"
	"strconv"
	flag "github.com/spf13/pflag"
)

const (
	FlagAddAssetMappingFile  = "add-asset-mapping-file"
	FlagAddAssetsMappingFile = "add-assets-file"
)

func ParseBoolFromString(s string) bool {
	switch s {
	case "1":
		return true
	default:
		return false
	}
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

func ParseStringFromString(s string, separator string) ([]string, error) {
	var parsedStrings []string
	for _, s := range strings.Split(s, separator) {
		s = strings.TrimSpace(s)

		parsedStrings = append(parsedStrings, s)
	}
	return parsedStrings, nil
}

func FlagSetCreateAssetMapping() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagAddAssetMappingFile, "", "add asset mapping json file path")
	return fs
}

func FlagSetCreateAssetsMapping() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagAddAssetsMappingFile, "", "add assets json file path")
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

type createAddAssetsMappingInputs struct {
	Name             string `json:"name"`
	Denom            string `json:"denom"`
	Decimals         string `json:"decimals"`
	IsOnChain        string `json:"is_on_chain"`
	AssetOraclePrice string `json:"asset_oracle_price"`
	IsCdpMintable    string `json:"is_cdp_mintable"`
	Title            string
	Description      string
	Deposit          string
}

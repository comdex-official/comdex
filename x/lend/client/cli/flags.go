package cli

import (
	flag "github.com/spf13/pflag"
	"strconv"
	"strings"
)

const (
	FlagNewLendPairFile        = "add-lend-pair-file"
	FlagAddLendPoolFile        = "add-lend-pool-file"
	FlagAddAssetRatesStatsFile = "add-asset-rates-stats-file"
)

func ParseStringFromString(s string, separator string) ([]string, error) {
	var parsedStrings []string
	for _, s := range strings.Split(s, separator) {
		s = strings.TrimSpace(s)

		parsedStrings = append(parsedStrings, s)
	}
	return parsedStrings, nil
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

func ParseBoolFromString(s uint64) bool {
	switch s {
	case 1:
		return true
	default:
		return false
	}
}

func FlagSetNewLendPairsMapping() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagNewLendPairFile, "", "add new lend pairs json file path")
	return fs
}

func FlagSetAddLendPoolMapping() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagAddLendPoolFile, "", "add new lend pool json file path")
	return fs
}

func FlagSetAddAssetRatesStatsMapping() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagAddAssetRatesStatsFile, "", "add asset rates stats json file path")
	return fs
}

type addNewLendPairsInputs struct {
	AssetIn        string `json:"asset_in"`
	AssetOut       string `json:"asset_out"`
	IsInterPool    string `json:"is_inter_pool"`
	AssetOutPoolID string `json:"asset_out_pool_id"`
	Title          string
	Description    string
	Deposit        string
}

type addLendPoolInputs struct {
	ModuleName           string `json:"module_name"`
	FirstBridgedAssetID  string `json:"first_bridged_asset_id"`
	SecondBridgedAssetID string `json:"second_bridged_asset_id"`
	AssetID              string `json:"asset_id"`
	IsBridgedAsset       string `json:"is_bridged_asset"`
	Title                string
	Description          string
	Deposit              string
}

type addAssetRatesStatsInputs struct {
	AssetID              string `json:"asset_id"`
	UOptimal             string `json:"u_optimal"`
	Base                 string `json:"base"`
	Slope1               string `json:"slope_1"`
	Slope2               string `json:"slope_2"`
	StableBase           string `json:"stable_base"`
	StableSlope1         string `json:"stable_slope_1"`
	StableSlope2         string `json:"stable_slope_2"`
	LTV                  string `json:"ltv"`
	LiquidationThreshold string `json:"liquidation_threshold"`
	LiquidationPenalty   string `json:"liquidation_penalty"`
	ReserveFactor        string `json:"reserve_factor"`
	CAssetId             string `json:"c_asset_id"`
	Title                string
	Description          string
	Deposit              string
}

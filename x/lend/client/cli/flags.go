package cli

import (
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

const (
	FlagNewLendPairFile         = "add-lend-pair-file"
	FlagAddLendPoolFile         = "add-lend-pool-file"
	FlagAddAssetRatesParamsFile = "add-asset-rates-params-file"
	FlagSetAuctionParamsFile    = "add-auction-params-file"
)

func ParseUint64SliceFromString(s string, separator string) ([]uint64, error) {
	var parsedInt []uint64
	for _, s := range strings.Split(s, separator) {
		s = strings.TrimSpace(s)

		parsed, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return []uint64{}, err
		}
		parsedInt = append(parsedInt, parsed)
	}
	return parsedInt, nil
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

func FlagSetAddAssetRatesParamsMapping() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagAddAssetRatesParamsFile, "", "add asset rates stats json file path")
	return fs
}

func FlagSetAuctionParams() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagSetAuctionParamsFile, "", "add auction params json file path")
	return fs
}

type addNewLendPairsInputs struct {
	AssetIn         string `json:"asset_in"`
	AssetOut        string `json:"asset_out"`
	IsInterPool     string `json:"is_inter_pool"`
	AssetOutPoolID  string `json:"asset_out_pool_id"`
	MinUSDValueLeft string `json:"min_usd_value_left"`
	Title           string
	Description     string
	Deposit         string
}

type addLendPoolInputs struct {
	ModuleName       string `json:"module_name"`
	AssetID          string `json:"asset_id"`
	AssetTransitType string `json:"asset_transit_type"`
	SupplyCap        string `json:"supply_cap"`
	CPoolName        string `json:"c_pool_name"`
	ReserveFunds     string `json:"reserve_funds"`
	Title            string
	Description      string
	Deposit          string
}

type addAssetRatesParamsInputs struct {
	AssetID              string `json:"asset_id"`
	UOptimal             string `json:"u_optimal"`
	Base                 string `json:"base"`
	Slope1               string `json:"slope_1"`
	Slope2               string `json:"slope_2"`
	EnableStableBorrow   string `json:"enable_stable_borrow"`
	StableBase           string `json:"stable_base"`
	StableSlope1         string `json:"stable_slope_1"`
	StableSlope2         string `json:"stable_slope_2"`
	LTV                  string `json:"ltv"`
	LiquidationThreshold string `json:"liquidation_threshold"`
	LiquidationPenalty   string `json:"liquidation_penalty"`
	LiquidationBonus     string `json:"liquidation_bonus"`
	ReserveFactor        string `json:"reserve_factor"`
	CAssetID             string `json:"c_asset_id"`
	Title                string
	Description          string
	Deposit              string
}

type addNewAuctionParamsInputs struct {
	AppID                  string `json:"app_id"`
	AuctionDurationSeconds string `json:"auction_duration_seconds"`
	Buffer                 string `json:"buffer"`
	Cusp                   string `json:"cusp"`
	Step                   string `json:"step"`
	PriceFunctionType      string `json:"price_function_type"`
	DutchID                string `json:"dutch_id"`
	BidDurationSeconds     string `json:"bid_duration_seconds"`
	Title                  string
	Description            string
	Deposit                string
}

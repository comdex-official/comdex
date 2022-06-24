package cli

import (
	flag "github.com/spf13/pflag"
	"strconv"
	"strings"
)

const (
	FlagAddLookupParamsTable = "collector-lookup-table-file"
	FlagAuctionControlParams = "auction-control-params-file"
)

func ParseStringFromString(s string, separator string) ([]string, error) {
	var parsedStrings []string
	for _, s := range strings.Split(s, separator) {
		s = strings.TrimSpace(s)

		parsedStrings = append(parsedStrings, s)
	}
	return parsedStrings, nil
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

func FlagSetCreateLookupTableParamsMapping() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagAddLookupParamsTable, "", "create lookup table params json file path")
	return fs
}

func FlagSetAuctionControlParamsMapping() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagAuctionControlParams, "", "auction control params json file path")
	return fs
}

type createLookupTableParamsInputs struct {
	AppID            string `json:"app_id"`
	CollectorAssetID string `json:"collector_asset_id"`
	SecondaryAssetID string `json:"secondary_asset_id"`
	SurplusThreshold string `json:"surplus_threshold"`
	DebtThreshold    string `json:"debt_threshold"`
	LockerSavingRate string `json:"locker_saving_rate"`
	LotSize          string `json:"lot_size"`
	BidFactor        string `json:"bid_factor"`
	DebtLotSize      string `json:"debt_lot_size"`
	Title            string
	Description      string
	Deposit          string
}

type auctionControlParamsInputs struct {
	AppID               string `json:"app_id"`
	AssetID             string `json:"asset_id"`
	SurplusAuction      string `json:"surplus_auction"`
	DebtAuction         string `json:"debt_auction"`
	AssetOutOraclePrice string `json:"asset_out_oracle_price"`
	AssetOutPrice       string `json:"asset_out_price"`
	Title               string
	Description         string
	Deposit             string
}

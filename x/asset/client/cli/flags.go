package cli

import (
	flag "github.com/spf13/pflag"
	"strconv"
	"strings"
)

const (
	flagName                    = "name"
	flagDenom                   = "denom"
	flagDecimals                = "decimals"
	flagCollateralWeight        = "collateralWeight"
	flagLiquidationThreshold    = "liquidationThreshold"
	FlagExtendedPairVaultFile   = "extended-pair-vault-file"
	FlagAddAssetMappingFile     = "add-asset-mapping-file"
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

func FlagSetCreateExtendedPairVault() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagExtendedPairVaultFile, "", "extended json file path")
	return fs
}

func FlagSetCreateAssetMapping() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagAddAssetMappingFile, "", "add asset mapping json file path")
	return fs
}


type createExtPairVaultInputs struct {
	AppMappingID        string `json:"app_mapping_id"`
	PairID              string `json:"pair_id"`
	StabilityFee        string `json:"stability_fee"`
	ClosingFee          string `json:"closing_fee"`
	LiquidationPenalty  string `json:"liquidation_penalty"`
	DrawDownFee         string `json:"draw_down_fee"`
	IsVaultActive       string `json:"is_vault_active"`
	DebtCeiling         string `json:"debt_ceiling"`
	DebtFloor           string `json:"debt_floor"`
	IsPsmPair           string `json:"is_psm_pair"`
	MinCr               string `json:"min_cr"`
	PairName            string `json:"pair_name"`
	AssetOutOraclePrice string `json:"asset_out_oracle_price"`
	AssetOutPrice       string `json:"asset_out_price"`
	MinUsdValueLeft     string `json:"min_usd_value_left"`
	Title               string
	Description         string
	Deposit             string
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


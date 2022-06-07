package cli

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"strconv"
	"strings"
)

const (
	flagLiquidationRatio      = "liquidation-ratio"
	flagName                  = "name"
	flagDenom                 = "denom"
	flagDecimals              = "decimals"
	flagCollateralWeight      = "collateralWeight"
	flagLiquidationThreshold  = "liquidationThreshold"
	flagIsBridgedAsset        = "isBridgedAsset"
	flagbaseborrowrateasset1  = "baseBorrowRate1"
	flagbaseborrowrateasset2  = "baseBorrowRate2"
	flagbaselendrateasset1    = "baseLendRate1"
	flagbaselendrateasset2    = "baseLendRate2"
	flagModuleAcc             = "moduleAcc"
	FlagExtendedPairVaultFile = "extended-pair-vault-file"
)

func GetLiquidationRatio(cmd *cobra.Command) (sdk.Dec, error) {
	s, err := cmd.Flags().GetString(flagLiquidationRatio)
	if err != nil {
		return sdk.Dec{}, err
	}

	return sdk.NewDecFromStr(s)
}

func ParseStringFromString(s string, seperator string) ([]string, error) {
	var parsedStrings []string
	for _, s := range strings.Split(s, seperator) {
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

func ParseInt64SliceFromString(s string, seperator string) ([]int64, error) {
	var parsedInts []int64
	for _, s := range strings.Split(s, seperator) {
		s = strings.TrimSpace(s)

		parsed, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return []int64{}, err
		}
		parsedInts = append(parsedInts, parsed)
	}
	return parsedInts, nil
}

func ParseUint64SliceFromString(s string, seperator string) ([]uint64, error) {
	var parsedInts []uint64
	for _, s := range strings.Split(s, seperator) {
		s = strings.TrimSpace(s)

		parsed, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return []uint64{}, err
		}
		parsedInts = append(parsedInts, parsed)
	}
	return parsedInts, nil
}

func FlagSetCreateExtendedPaiVault() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagExtendedPairVaultFile, "", "extended json file path")
	return fs
}

type createExtPairVaultInputs struct {
	AppMappingId        string `json:"app_mapping_id"`
	PairId              string `json:"pair_id"`
	LiquidationRatio    string `json:"liquidation_ratio"`
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
}

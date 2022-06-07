package cli

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"strconv"
	"strings"
)

const (
	flagLiquidationRatio        = "liquidation-ratio"
	flagName                    = "name"
	flagDenom                   = "denom"
	flagDecimals                = "decimals"
	flagCollateralWeight        = "collateralWeight"
	flagLiquidationThreshold    = "liquidationThreshold"
	flagIsBridgedAsset          = "isBridgedAsset"
	flagBaseBorrowRateAsset1    = "baseBorrowRate1"
	flagBaseBorrowRateAsset2    = "baseBorrowRate2"
	flagBaseLendRateAsset1      = "baseLendRate1"
	flagBaseLendRateAsset2      = "baseLendRate2"
	flagModuleAcc               = "moduleAcc"
	FlagExtendedPairVaultFile   = "extended-pair-vault-file"
	FlagAddAssetMappingFile     = "add-asset-mapping-file"
	FlagAddWhiteListedPairsFile = "add-white-whitelisted-pairs-file"
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

func FlagSetCreateWhiteListedPairsMapping() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagAddWhiteListedPairsFile, "", "add white listed asset pairs json file path")
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
	Title               string
	Description         string
	Deposit             string
}

type createAddAssetMappingInputs struct {
	AppId         string `json:"app_id"`
	AssetId       string `json:"asset_id"`
	GenesisSupply string `json:"genesis_supply"`
	IsGovToken    string `json:"is_gov_token"`
	Recipient     string `json:"recipient"`
	Title         string
	Description   string
	Deposit       string
}

type createAddWhiteListedPairsInputs struct {
	PairId               string `json:"pair_id"`
	ModuleAccount        string `json:"module-account"`
	BaseBorrowRateAsset1 string `json:"base_borrow_rate_asset_1"`
	BaseBorrowRateAsset2 string `json:"base_borrow_rate_asset_2"`
	BaseLendRateAsset1   string `json:"base_lend_rate_asset_1"`
	BaseLendRateAsset2   string `json:"base_lend_rate_asset_2"`
	Title                string
	Description          string
	Deposit              string
}

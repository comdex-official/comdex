package cli

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

const (
	flagLiquidationRatio     = "liquidation-ratio"
	flagName                 = "name"
	flagDenom                = "denom"
	flagDecimals             = "decimals"
	flagCollateralWeight     = "collateralWeight"
	flagLiquidationThreshold = "liquidationThreshold"
	flagIsBridgedAsset       = "isBridgedAsset"
	flagbaseborrowrateasset1 = "baseBorrowRate1"
	flagbaseborrowrateasset2 = "baseBorrowRate2"
	flagbaselendrateasset1   = "baseLendRate1"
	flagbaselendrateasset2   = "baseLendRate2"
	flagModuleAcc            = "moduleAcc"
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
	if s == "1" {
		return true

	} else {
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

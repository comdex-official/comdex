package cli

import (
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
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

func ParseStringFromString(s string, separator string) ([]string, error) {
	stringsSlice := strings.Split(s, separator)
	parsedStrings := make([]string, 0, len(stringsSlice))
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
	stringsSlice := strings.Split(s, separator)
	parsedInts := make([]int64, 0, len(stringsSlice))
	for _, s := range stringsSlice {
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
	stringsSlice := strings.Split(s, separator)
	parsedInts := make([]uint64, 0, len(stringsSlice))
	for _, s := range stringsSlice {
		s = strings.TrimSpace(s)

		parsed, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return []uint64{}, err
		}
		parsedInts = append(parsedInts, parsed)
	}
	return parsedInts, nil
}

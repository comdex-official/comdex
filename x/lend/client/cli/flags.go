package cli

import (
	"strconv"
	"strings"
)

const (
	flagName                 = "name"
	flagDenom                = "denom"
	flagDecimal              = "decimal"
	flagCollateralWeight     = "collateralWeight"
	flagLiquidationThreshold = "liquidationThreshold"
	flagBaseBorrowRate       = "baseBorrowRate"
	flagBaseLendRate         = "baseLendRate"
	flagAssetOne             = "assetOne"
	flagAssetTwo             = "assetTwo"
	flagModuleAcc            = "moduleAcc"
)

func ParseStringFromString(s string, seperator string) ([]string, error) {
	var parsedStrings []string
	for _, s := range strings.Split(s, seperator) {
		s = strings.TrimSpace(s)

		parsedStrings = append(parsedStrings, s)
	}
	return parsedStrings, nil
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

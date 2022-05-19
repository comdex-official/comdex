package keeper

import (
	"fmt"
	"strings"

	"github.com/comdex-official/comdex/x/incentives/types"
)

func IntegerArrayToString(intArray []uint64) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(types.ValidGaugeTypeIds)), ","), "[]")
}

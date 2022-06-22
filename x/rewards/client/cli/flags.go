package cli

import (
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

// flags for incentives module tx commands.
const (
	FlagStartTime = "start-time"

	FlagPoolID       = "pool-id"
	FlagAppID        = "app-id"
	FlagIsMasterPool = "is-master-pool"
	FlagChildPoolIds = "child-pool-ids"
)

// FlagSetCreateGauge returns flags for creating gauge.
func FlagSetCreateGauge() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	// Global Flags
	fs.String(FlagStartTime, "", "Timestamp to begin distribution")

	// Msg Specific Flags - Liquidity GaugeType Flags.
	fs.Uint64(FlagPoolID, 0, "Pool Id")
	fs.Uint64(FlagAppID, 0, "App Id")
	fs.Bool(FlagIsMasterPool, false, "If gauge is for master pool or not, default false")
	fs.String(FlagChildPoolIds, "", "List of child pool ids, default [] i.e all pools")

	return fs
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

package cli

import (
	"strconv"
	"strings"

	"time"

	flag "github.com/spf13/pflag"
)

// flags for incentives module tx commands.
const (
	// Global Flags.
	FlagStartTime = "start-time"

	// Msg Specific Flags - Liquidity GaugeType Flags.
	FlagPoolID       = "pool-id"
	FlagIsMasterPool = "is-master-pool"
	FlagChildPoolIds = "child-pool-ids"
	FlagLockDuration = "lock-duration"
)

// FlagSetCreateGauge returns flags for creating gauge.
func FlagSetCreateGauge() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	// Global Flags
	fs.String(FlagStartTime, "", "Timestamp to begin distribution")

	// Msg Specific Flags - Liquidity GaugeType Flags.
	dur, _ := time.ParseDuration("24h")
	fs.Uint64(FlagPoolID, 0, "Pool Id")
	fs.Bool(FlagIsMasterPool, false, "If gauge is for master pool or not, default false")
	fs.String(FlagChildPoolIds, "", "List of child pool ids, default [] i.e all pools")
	fs.Duration(FlagLockDuration, dur, "Bonding duration, default 24h")

	return fs
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

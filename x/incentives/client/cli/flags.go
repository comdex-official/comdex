package cli

import (
	"time"

	flag "github.com/spf13/pflag"
)

// flags for incentives module tx commands
const (
	// Global Flags
	FlagStartTime = "start-time"

	// Msg Specific Flags - Liquidity GaugeType Flags
	FlagPoolId       = "pool-id"
	FlagIsMasterPool = "is-master-pool"
	FlagChildPoolIds = "child-pool-ids"
	FlagLockDuration = "lock-duration"
)

// FlagSetCreateGauge returns flags for creating gauge
func FlagSetCreateGauge() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	// Global Flags
	fs.String(FlagStartTime, "", "Timestamp to begin distribution")

	// Msg Specific Flags - Liquidity GaugeType Flags
	dur, _ := time.ParseDuration("24h")
	fs.Uint64(FlagPoolId, 0, "Pool Id")
	fs.Bool(FlagIsMasterPool, false, "If gauge is for master pool or not, default false")
	fs.String(FlagChildPoolIds, "", "List of child pool ids, default [] i.e all pools")
	fs.Duration(FlagLockDuration, dur, "Bonding duration, default 24h")

	return fs
}

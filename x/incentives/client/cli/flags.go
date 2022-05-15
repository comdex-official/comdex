package cli

import (
	"time"

	flag "github.com/spf13/pflag"
)

// flags for incentives module tx commands
const (

	// Liquidity GaugeType Flags
	FlagPoolId       = "pool-id"
	FlagIsMasterPool = "is-master-pool"
	FlagChildPoolIds = "child-pool-ids"
	FlagStartTime    = "start-time"
	FlagLockDuration = "lock-duration"
)

// FlagSetCreateGauge returns flags for creating gauge
func FlagSetCreateGauge() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	// Liquidity Metadata GaugeType Flags
	dur, _ := time.ParseDuration("24h")
	fs.Uint64(FlagPoolId, 0, "Pool Id")
	fs.Bool(FlagIsMasterPool, false, "If gauge is for master pool or not, default false")
	fs.String(FlagChildPoolIds, "", "List of child pool ids, default [] i.e all pools")
	fs.String(FlagStartTime, "", "Timestamp to begin distribution")
	fs.Duration(FlagLockDuration, dur, "Bonding duration, default 24h")

	return fs
}

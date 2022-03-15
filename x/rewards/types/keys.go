package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	ModuleName   = "rewards"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
	MemStoreKey  = "mem_rewards"
)

var (
	MintingRewardsIdKey     = []byte{0x01}
	MintingRewardsKeyPrefix = []byte{0x11}
)

func MintingRewardsKey(id uint64) []byte {
	return append(MintingRewardsKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

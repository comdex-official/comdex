package types

const (
	ModuleName   = "rewards"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
	MemStoreKey  = "mem_rewards"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

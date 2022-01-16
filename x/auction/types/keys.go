package types

const (
	ModuleName     = "auction"
	ParamsSubspace = ModuleName
	QuerierRoute   = ModuleName
	RouterKey      = ModuleName
	StoreKey       = ModuleName
	MemStoreKey    = ModuleName
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

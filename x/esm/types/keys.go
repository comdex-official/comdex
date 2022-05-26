package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "esm"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

    // QuerierRoute defines the module's query routing key
    QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_esm"

    
)
var (
	TriggerEsmNewKey = [] byte{0x01}
)


func KeyPrefix(p string) []byte {
    return []byte(p)
}


func TriggerEsmKey(app_id uint64) []byte {
	return append(TriggerEsmNewKey, sdk.Uint64ToBigEndian(app_id)...)
}

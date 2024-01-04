package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName defines the module name.
	ModuleName   = "rwa"
	QuerierRoute = ModuleName

	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName
)

var (
	RwaUSerKeyPrefix      = []byte{0x10}
	CounterPartyKeyPrefix = []byte{0x11}
	InvoiceKeyPrefix      = []byte{0x12}
)

func RwaUserKey(address string) []byte {
	return append(RwaUSerKeyPrefix, address...)
}

func CounterPartyKey(ID uint64) []byte {
	return append(CounterPartyKeyPrefix, sdk.Uint64ToBigEndian(ID)...)
}

func InvoiceKey(ID uint64) []byte {
	return append(InvoiceKeyPrefix, sdk.Uint64ToBigEndian(ID)...)
}

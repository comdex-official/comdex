package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName   = "tokenmint"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
	MemStoreKey  = "mem_tokenmint"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

var TypeMsgMintNewTokensRequest = ModuleName + ":mintnewtokens"
var TypeMsgBurnHarborTokensRequest = ModuleName + ":burntokens"

var TokenMintKeyPrefix = []byte{0x10}

func TokenMintKey(appMappingID uint64) []byte {
	return append(TokenMintKeyPrefix, sdk.Uint64ToBigEndian(appMappingID)...)
}

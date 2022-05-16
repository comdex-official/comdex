package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName   = "vault"
	QuerierRoute = ModuleName
	RouterKey    = ModuleName
	StoreKey     = ModuleName
)

var (
	TypeMsgCreateRequest    = ModuleName + ":create"
	TypeMsgDepositRequest   = ModuleName + ":deposit"
	TypeMsgWithdrawRequest  = ModuleName + ":withdraw"
	TypeMsgDrawRequest      = ModuleName + ":draw"
	TypeMsgRepayRequest     = ModuleName + ":repay"
	TypeMsgLiquidateRequest = ModuleName + ":liquidate"
)

var (
	IDKey                          = []byte{0x00}
	VaultKeyPrefix                 = []byte{0x10}
	VaultForAddressByPairKeyPrefix = []byte{0x20}
	LookupTableVaultPrefix         = []byte{0x30}
	CounterVaultKeyPrefix = []byte{0x40}
	VaultIDPrefix = []byte{0x50}
	PairIDPrefix = []byte{0x60}
	ExtendedIDPrefix = []byte{0x70}
	ExtendedpairIDPrefix = []byte{0x80}
	TokenMintedIDPrefix = []byte{0x90}
)

func VaultKey(AppVaultTypeId string) []byte {
	return append(VaultKeyPrefix, []byte(AppVaultTypeId)...)
}

func VaultForAddressByPair(address sdk.AccAddress, ExtendedPairVaultID uint64) []byte {
	v := append(VaultForAddressByPairKeyPrefix, address.Bytes()...)
	if len(v) != 1+20 {
		panic(fmt.Errorf("invalid key length %d; expected %d", len(v), 1+20))
	}

	return append(v, sdk.Uint64ToBigEndian(ExtendedPairVaultID)...)
}

func LookupTableVaultKey(AppId uint64) []byte {
	return append(LookupTableVaultPrefix, sdk.Uint64ToBigEndian(AppId)...)
}

func CounterVaultKey(AppId uint64) []byte {
	return append(CounterVaultKeyPrefix, sdk.Uint64ToBigEndian(AppId)...)
}

func UserVaultIdMappingKey(owner string) []byte {
	return append(VaultIDPrefix, []byte(owner)...)

}

func PairIDKey(PairId uint64) []byte {
	return append(PairIDPrefix, sdk.Uint64ToBigEndian(PairId)...)
}

func ExtendedIDKey(Appid uint64) []byte {
	return append(ExtendedIDPrefix, sdk.Uint64ToBigEndian(Appid)...)
}

func ExtendedpairIDKey(Pairid uint64) []byte {
	return append(ExtendedpairIDPrefix, sdk.Uint64ToBigEndian(Pairid)...)
}

func TokenMintedIDKey(AppId uint64) []byte {
	return append(TokenMintedIDPrefix, sdk.Uint64ToBigEndian(AppId)...)
}
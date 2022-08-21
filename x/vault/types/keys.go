package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName   = "vaultv1"
	QuerierRoute = ModuleName
	RouterKey    = ModuleName
	StoreKey     = ModuleName
)

var (
	TypeMsgCreateRequest             = ModuleName + ":create"
	TypeMsgDepositRequest            = ModuleName + ":deposit"
	TypeMsgWithdrawRequest           = ModuleName + ":withdraw"
	TypeMsgDrawRequest               = ModuleName + ":draw"
	TypeMsgRepayRequest              = ModuleName + ":repay"
	TypeMsgLiquidateRequest          = ModuleName + ":liquidate"
	TypeMsgCreateStableMintRequest   = ModuleName + ":create_stablemint"
	TypeMsgDepositStableMintRequest  = ModuleName + ":deposit_stablemint"
	TypeMsgWithdrawStableMintRequest = ModuleName + ":withdraw_stablemint"
)

var (
	VaultKeyPrefix                        = []byte{0x10}
	UserVaultExtendedPairMappingKeyPrefix = []byte{0x12}
	AppExtendedPairVaultMappingKeyPrefix  = []byte{0x13}
	StableMintVaultKeyPrefix              = []byte{0x14}
	VaultIDPrefix                         = []byte{0x15}
	StableVaultIDPrefix                   = []byte{0x16}
)

func VaultKey(vaultID uint64) []byte {
	return append(VaultKeyPrefix, sdk.Uint64ToBigEndian(vaultID)...)
}
func StableMintVaultKey(stableVaultID uint64) []byte {
	return append(StableMintVaultKeyPrefix, sdk.Uint64ToBigEndian(stableVaultID)...)
}

// func UserVaultExtendedPairMappingKey(address string) []byte {
// 	return append(UserVaultExtendedPairMappingKeyPrefix, address...)
// }

func AppExtendedPairVaultMappingKey(appMappingID uint64, pairVaultID uint64) []byte {
	return append(append(AppExtendedPairVaultMappingKeyPrefix, sdk.Uint64ToBigEndian(appMappingID)...), sdk.Uint64ToBigEndian(pairVaultID)...)
}

func AppMappingKey(appMappingID uint64) []byte {
	return append(AppExtendedPairVaultMappingKeyPrefix, sdk.Uint64ToBigEndian(appMappingID)...)
}

func UserAppExtendedPairMappingKey(address string, appID uint64, pairVaultID uint64) []byte {
	return append(append(append(UserVaultExtendedPairMappingKeyPrefix, address...), sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(pairVaultID)...)
}

func UserAppMappingKey(address string, appID uint64) []byte {
	return append(append(UserVaultExtendedPairMappingKeyPrefix, address...), sdk.Uint64ToBigEndian(appID)...)
}

func UserKey(address string) []byte {
	return append(UserVaultExtendedPairMappingKeyPrefix, address...)
}

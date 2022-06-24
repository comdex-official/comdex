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
	TypeMsgCreateRequest    = ModuleName + ":create"
	TypeMsgDepositRequest   = ModuleName + ":deposit"
	TypeMsgWithdrawRequest  = ModuleName + ":withdraw"
	TypeMsgDrawRequest      = ModuleName + ":draw"
	TypeMsgRepayRequest     = ModuleName + ":repay"
	TypeMsgLiquidateRequest = ModuleName + ":liquidate"
)

var (
	VaultKeyPrefix                        = []byte{0x10}
	UserVaultExtendedPairMappingKeyPrefix = []byte{0x12}
	AppExtendedPairVaultMappingKeyPrefix  = []byte{0x13}
	StableMintVaultKeyPrefix              = []byte{0x14}
)

func VaultKey(vaultID string) []byte {
	return append(VaultKeyPrefix, vaultID...)
}
func StableMintVaultKey(stableVaultID string) []byte {
	return append(StableMintVaultKeyPrefix, stableVaultID...)
}

func UserVaultExtendedPairMappingKey(address string) []byte {
	return append(UserVaultExtendedPairMappingKeyPrefix, address...)
}

func AppExtendedPairVaultMappingKey(appMappingID uint64) []byte {
	return append(AppExtendedPairVaultMappingKeyPrefix, sdk.Uint64ToBigEndian(appMappingID)...)
}

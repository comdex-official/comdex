package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName   = "vaultV1"
	QuerierRoute = ModuleName
	RouterKey    = ModuleName
	StoreKey     = ModuleName
)

var (
	TypeMsgCreateRequest                    = ModuleName + ":create"
	TypeMsgDepositRequest                   = ModuleName + ":deposit"
	TypeMsgWithdrawRequest                  = ModuleName + ":withdraw"
	TypeMsgDrawRequest                      = ModuleName + ":draw"
	TypeMsgRepayRequest                     = ModuleName + ":repay"
	TypeMsgLiquidateRequest                 = ModuleName + ":liquidate"
	TypeMsgDepositDrawRequest               = ModuleName + ":deposit_draw"
	TypeMsgCreateStableMintRequest          = ModuleName + ":create_stablemint"
	TypeMsgDepositStableMintRequest         = ModuleName + ":deposit_stablemint"
	TypeMsgWithdrawStableMintRequest        = ModuleName + ":withdraw_stablemint"
	TypeMsgVaultInterestCalcRequest         = ModuleName + ":calculate_interest"
	TypeMsgWithdrawStableMintControlRequest = ModuleName + ":withdraw_stablemint_control"
)

var (
	VaultKeyPrefix                        = []byte{0x10}
	UserVaultExtendedPairMappingKeyPrefix = []byte{0x12}
	AppExtendedPairVaultMappingKeyPrefix  = []byte{0x13}
	StableMintVaultKeyPrefix              = []byte{0x14}
	VaultIDPrefix                         = []byte{0x15}
	StableVaultIDPrefix                   = []byte{0x16}
	VaultLengthPrefix                     = []byte{0x17}
	StableVaultRewardsKeyPrefix           = []byte{0x18}
	StableVaultControlKeyPrefix           = []byte{0x19}
)

func VaultKey(vaultID uint64) []byte {
	return append(VaultKeyPrefix, sdk.Uint64ToBigEndian(vaultID)...)
}

func StableMintVaultKey(stableVaultID uint64) []byte {
	return append(StableMintVaultKeyPrefix, sdk.Uint64ToBigEndian(stableVaultID)...)
}

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

func StableMintVaultRewardsKey(appID uint64, address string, blockHeight uint64) []byte {
	return append(append(append(StableVaultRewardsKeyPrefix, sdk.Uint64ToBigEndian(appID)...), address...), sdk.Uint64ToBigEndian(blockHeight)...)
}

func StableMintRewardsKey(appID uint64, address string) []byte {
	return append(append(StableVaultRewardsKeyPrefix, sdk.Uint64ToBigEndian(appID)...), address...)
}

func StableMintRewardsAppKey(appID uint64) []byte {
	return append(StableVaultRewardsKeyPrefix, sdk.Uint64ToBigEndian(appID)...)
}

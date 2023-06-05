package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName defines the module name
	ModuleName = "liquidationsV2"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key.
	MemStoreKey = "mem_newliq"
)

var (
	TypeMsgLiquidateRequest          = ModuleName + ":liquidate"
	TypeMsgLiquidateExternalRequest  = ModuleName + ":liquidate_external"
	TypeAppReserveFundsRequest       = ModuleName + ":app_reserve_funds"
	AppIdsKeyPrefix                  = []byte{0x01}
	LiquidationOffsetHolderKeyPrefix = []byte{0x02}
	LockedVaultIDKey                 = []byte{0x03}
	LockedVaultKeyPrefix             = []byte{0x04}
	LiquidationWhiteListingKeyPrefix = []byte{0x05}
	AppReserveFundsKeyPrefix         = []byte{0x06}
	AppReserveFundsTxDataKeyPrefix   = []byte{0x07}
)

// LengthPrefixString returns length-prefixed bytes representation
// of a string.
func LengthPrefixString(s string) []byte {
	bz := []byte(s)
	bzLen := len(bz)
	if bzLen == 0 {
		return bz
	}
	return append([]byte{byte(bzLen)}, bz...)
}

// GetLiquidationOffsetHolderKey returns the index key to look offset value for liquidation.
func GetLiquidationOffsetHolderKey(appID uint64, liquidationForPrefix string) []byte {
	return append(append(LiquidationOffsetHolderKeyPrefix, sdk.Uint64ToBigEndian(appID)...), LengthPrefixString(liquidationForPrefix)...)
}

// WhitelistAppKeyByApp whitelisting kv
func WhitelistAppKeyByApp(appID uint64) []byte {
	return append(AppIdsKeyPrefix, sdk.Uint64ToBigEndian(appID)...)
}

func LockedVaultKey(appID, lockedVaultID uint64) []byte {
	return append(append(LockedVaultKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(lockedVaultID)...)
}

func LockedVaultKeyByApp(appID uint64) []byte {
	return append(LockedVaultKeyPrefix, sdk.Uint64ToBigEndian(appID)...)
}

func LiquidationWhiteListingKey(appId uint64) []byte {
	return append(LiquidationWhiteListingKeyPrefix, sdk.Uint64ToBigEndian(appId)...)
}

func AppReserveFundsKey(appID, assetID uint64) []byte {
	return append(append(AppReserveFundsKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(assetID)...)
}

func AppReserveFundsTxDataKey(appId uint64) []byte {
	return append(AppReserveFundsTxDataKeyPrefix, sdk.Uint64ToBigEndian(appId)...)
}

package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name.
	ModuleName = "liquidityV1"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// RouterKey is the message route for slashing.
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key.
	QuerierRoute = ModuleName
)

var (
	LastPairIDKey = []byte{0xa0} // key for the latest pair id
	LastPoolIDKey = []byte{0xa1} // key for the latest pool id

	PairKeyPrefix               = []byte{0xa2}
	PairIndexKeyPrefix          = []byte{0xa3}
	PairsByDenomsIndexKeyPrefix = []byte{0xa4}

	PoolKeyPrefix                      = []byte{0xa5}
	PoolByReserveAddressIndexKeyPrefix = []byte{0xa6}
	PoolsByPairIndexKeyPrefix          = []byte{0xa7}

	DepositRequestKeyPrefix       = []byte{0xa8}
	DepositRequestIndexKeyPrefix  = []byte{0xa9}
	WithdrawRequestKeyPrefix      = []byte{0xb0}
	WithdrawRequestIndexKeyPrefix = []byte{0xb1}
	OrderKeyPrefix                = []byte{0xb2}
	OrderIndexKeyPrefix           = []byte{0xb3}

	ActiveFarmerKeyPrefix = []byte{0xb4}
	QueuedFarmerKeyPrefix = []byte{0xb5}

	GenericParamsKey = []byte{0xb6}
)

// GetLastPairIDKey returns the store key to retrieve the last pair id.
func GetLastPairIDKey(appID uint64) []byte {
	return append(LastPairIDKey, sdk.Uint64ToBigEndian(appID)...)
}

// GetPairKey returns the store key to retrieve pair object from the pair id.
func GetPairKey(appID, pairID uint64) []byte {
	return append(append(PairKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(pairID)...)
}

// GetAllPairsKey returns the store key to retrieve all pairs.
func GetAllPairsKey(appID uint64) []byte {
	return append(PairKeyPrefix, sdk.Uint64ToBigEndian(appID)...)
}

// GetPairIndexKey returns the index key to get a pair by denoms.
func GetPairIndexKey(appID uint64, baseCoinDenom, quoteCoinDenom string) []byte {
	return append(append(append(PairIndexKeyPrefix, sdk.Uint64ToBigEndian(appID)...), LengthPrefixString(baseCoinDenom)...), LengthPrefixString(quoteCoinDenom)...)
}

// GetPairsByDenomsIndexKey returns the index key to lookup pairs with given denoms.
func GetPairsByDenomsIndexKey(appID uint64, denomA, denomB string, pairID uint64) []byte {
	return append(append(append(append(PairsByDenomsIndexKeyPrefix, sdk.Uint64ToBigEndian(appID)...), LengthPrefixString(denomA)...), LengthPrefixString(denomB)...), sdk.Uint64ToBigEndian(pairID)...)
}

// GetPairsByDenomIndexKeyPrefix returns the index key prefix to lookup pairs with given denom.
func GetPairsByDenomIndexKeyPrefix(appID uint64, denomA string) []byte {
	return append(append(PairsByDenomsIndexKeyPrefix, sdk.Uint64ToBigEndian(appID)...), LengthPrefixString(denomA)...)
}

// GetPairsByDenomsIndexKeyPrefix returns the index key prefix to lookup pairs with given denoms.
func GetPairsByDenomsIndexKeyPrefix(appID uint64, denomA, denomB string) []byte {
	return append(append(append(PairsByDenomsIndexKeyPrefix, sdk.Uint64ToBigEndian(appID)...), LengthPrefixString(denomA)...), LengthPrefixString(denomB)...)
}

// GetLastPoolIDKey returns the store key to retrieve the last pool id.
func GetLastPoolIDKey(appID uint64) []byte {
	return append(LastPoolIDKey, sdk.Uint64ToBigEndian(appID)...)
}

// GetPoolKey returns the store key to retrieve pool object from the pool id.
func GetPoolKey(appID, poolID uint64) []byte {
	return append(append(PoolKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(poolID)...)
}

// GetAllPoolsKey returns the store key to retrieve all pools.
func GetAllPoolsKey(appID uint64) []byte {
	return append(PoolKeyPrefix, sdk.Uint64ToBigEndian(appID)...)
}

// GetPoolByReserveAddressIndexKey returns the index key to retrieve the particular pool.
func GetPoolByReserveAddressIndexKey(appID uint64, reserveAddr sdk.AccAddress) []byte {
	return append(append(PoolByReserveAddressIndexKeyPrefix, sdk.Uint64ToBigEndian(appID)...), address.MustLengthPrefix(reserveAddr)...)
}

// GetPoolsByPairIndexKey returns the index key to retrieve pool id that is used to iterate pools.
func GetPoolsByPairIndexKey(appID, pairID, poolID uint64) []byte {
	return append(append(append(PoolsByPairIndexKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(pairID)...), sdk.Uint64ToBigEndian(poolID)...)
}

// GetPoolsByPairIndexKeyPrefix returns the store key to retrieve pool id to iterate pools.
func GetPoolsByPairIndexKeyPrefix(appID, pairID uint64) []byte {
	return append(append(PoolsByPairIndexKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(pairID)...)
}

// GetAllDepositRequestKey returns the store key to retrieve all deposit requests.
func GetAllDepositRequestKey(appID uint64) []byte {
	return append(DepositRequestKeyPrefix, sdk.Uint64ToBigEndian(appID)...)
}

// GetDepositRequestKey returns the store key to retrieve deposit request object from the pool id and request id.
func GetDepositRequestKey(appID, poolID, id uint64) []byte {
	return append(append(append(DepositRequestKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(poolID)...), sdk.Uint64ToBigEndian(id)...)
}

// GetDepositRequestIndexKey returns the index key to map deposit requests
// with a depositor.
func GetDepositRequestIndexKey(
	appID uint64,
	depositor sdk.AccAddress,
	poolID, reqID uint64,
) []byte {
	return append(append(append(append(DepositRequestIndexKeyPrefix, sdk.Uint64ToBigEndian(appID)...), address.MustLengthPrefix(depositor)...),
		sdk.Uint64ToBigEndian(poolID)...), sdk.Uint64ToBigEndian(reqID)...)
}

// GetDepositRequestIndexKeyPrefix returns the index key prefix to iterate
// deposit requests by a depositor.
func GetDepositRequestIndexKeyPrefix(appID uint64, depositor sdk.AccAddress) []byte {
	return append(append(DepositRequestIndexKeyPrefix, sdk.Uint64ToBigEndian(appID)...), address.MustLengthPrefix(depositor)...)
}

// GetAllWithdrawRequestKey returns the store key to retrieve all withdraw requests.
func GetAllWithdrawRequestKey(appID uint64) []byte {
	return append(WithdrawRequestKeyPrefix, sdk.Uint64ToBigEndian(appID)...)
}

// GetWithdrawRequestKey returns the store key to retrieve withdraw request object from the pool id and request id.
func GetWithdrawRequestKey(appID, poolID, id uint64) []byte {
	return append(append(append(WithdrawRequestKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(poolID)...), sdk.Uint64ToBigEndian(id)...)
}

// GetWithdrawRequestIndexKey returns the index key to map withdraw requests
// with a withdrawer.
func GetWithdrawRequestIndexKey(
	appID uint64,
	withdrawer sdk.AccAddress,
	poolID, reqID uint64,
) []byte {
	return append(append(append(append(WithdrawRequestIndexKeyPrefix, sdk.Uint64ToBigEndian(appID)...), address.MustLengthPrefix(withdrawer)...),
		sdk.Uint64ToBigEndian(poolID)...), sdk.Uint64ToBigEndian(reqID)...)
}

// GetWithdrawRequestIndexKeyPrefix returns the index key prefix to iterate
// withdraw requests by a withdrawer.
func GetWithdrawRequestIndexKeyPrefix(appID uint64, depositor sdk.AccAddress) []byte {
	return append(append(WithdrawRequestIndexKeyPrefix, sdk.Uint64ToBigEndian(appID)...), address.MustLengthPrefix(depositor)...)
}

// GetOrderKey returns the store key to retrieve order object from the pair id and request id.
func GetOrderKey(appID, pairID, id uint64) []byte {
	return append(append(append(OrderKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(pairID)...), sdk.Uint64ToBigEndian(id)...)
}

// GetAllOrdersKey returns the store key to retrieve all orders.
func GetAllOrdersKey(appID uint64) []byte {
	return append(OrderKeyPrefix, sdk.Uint64ToBigEndian(appID)...)
}

// GetOrdersByPairKeyPrefix returns the store key to iterate orders by pair.
func GetOrdersByPairKeyPrefix(appID uint64, pairID uint64) []byte {
	return append(append(OrderKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(pairID)...)
}

// GetOrderIndexKey returns the index key to map orders with an orderer.
func GetOrderIndexKey(
	appID uint64,
	orderer sdk.AccAddress,
	pairID, orderID uint64,
) []byte {
	return append(append(append(append(OrderIndexKeyPrefix, sdk.Uint64ToBigEndian(appID)...), address.MustLengthPrefix(orderer)...),
		sdk.Uint64ToBigEndian(pairID)...), sdk.Uint64ToBigEndian(orderID)...)
}

// GetOrderIndexKeyPrefix returns the index key prefix to iterate orders
// by an orderer.
func GetOrderIndexKeyPrefix(appID uint64, orderer sdk.AccAddress) []byte {
	return append(append(OrderIndexKeyPrefix, sdk.Uint64ToBigEndian(appID)...), address.MustLengthPrefix(orderer)...)
}

// GetGenericParamsKey returns the store key to retrieve params object.
func GetGenericParamsKey(appID uint64) []byte {
	return append(GenericParamsKey, sdk.Uint64ToBigEndian(appID)...)
}

// ParsePairsByDenomsIndexKey parses a pair by denom index key.
func ParsePairsByDenomsIndexKey(key []byte) (denomA, denomB string, pairID uint64) {
	if !bytes.HasPrefix(key, PairsByDenomsIndexKeyPrefix) {
		panic("key does not have proper prefix")
	}
	denomALen := key[9]
	denomA = string(key[10 : 10+denomALen])
	denomBLen := key[10+denomALen]
	denomB = string(key[11+denomALen : 11+denomALen+denomBLen])
	pairID = sdk.BigEndianToUint64(key[11+denomALen+denomBLen:])
	return
}

// ParsePoolsByPairIndexKey parses a pool id from the index key.
func ParsePoolsByPairIndexKey(key []byte) (poolID uint64) {
	if !bytes.HasPrefix(key, PoolsByPairIndexKeyPrefix) {
		panic("key does not have proper prefix")
	}

	bytesLen := 8
	poolID = sdk.BigEndianToUint64(key[1+bytesLen+bytesLen:])
	return
}

// ParseDepositRequestIndexKey parses a deposit request index key.
func ParseDepositRequestIndexKey(key []byte) (depositor sdk.AccAddress, poolID, reqID uint64) {
	if !bytes.HasPrefix(key, DepositRequestIndexKeyPrefix) {
		panic("key does not have proper prefix")
	}

	addrLen := key[9]
	depositor = key[10 : 10+addrLen]
	poolID = sdk.BigEndianToUint64(key[10+addrLen : 10+addrLen+8])
	reqID = sdk.BigEndianToUint64(key[10+addrLen+8:])
	return
}

// ParseWithdrawRequestIndexKey parses a withdraw request index key.
func ParseWithdrawRequestIndexKey(key []byte) (withdrawer sdk.AccAddress, poolID, reqID uint64) {
	if !bytes.HasPrefix(key, WithdrawRequestIndexKeyPrefix) {
		panic("key does not have proper prefix")
	}

	addrLen := key[9]
	withdrawer = key[10 : 10+addrLen]
	poolID = sdk.BigEndianToUint64(key[10+addrLen : 10+addrLen+8])
	reqID = sdk.BigEndianToUint64(key[10+addrLen+8:])
	return
}

// ParseOrderIndexKey parses an order index key.
func ParseOrderIndexKey(key []byte) (orderer sdk.AccAddress, pairID, orderID uint64) {
	if !bytes.HasPrefix(key, OrderIndexKeyPrefix) {
		panic("key does not have proper prefix")
	}

	addrLen := key[9]
	orderer = key[10 : 10+addrLen]
	pairID = sdk.BigEndianToUint64(key[10+addrLen : 10+addrLen+8])
	orderID = sdk.BigEndianToUint64(key[10+addrLen+8:])
	return
}

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

// GetActiveFarmerKey returns the store key to retrieve active farmer object from the app id, pool id and farmer address.
func GetActiveFarmerKey(appID, poolID uint64, farmer sdk.AccAddress) []byte {
	return append(append(append(ActiveFarmerKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(poolID)...), address.MustLengthPrefix(farmer)...)
}

// GetAllActiveFarmerKey returns the store key to retrieve all active farmers.
func GetAllActiveFarmersKey(appID, poolID uint64) []byte {
	return append(append(ActiveFarmerKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(poolID)...)
}

// GetQueuedFarmerKey returns the store key to retrieve queued farmer object from the app id, pool id and farmer address.
func GetQueuedFarmerKey(appID, poolID uint64, farmer sdk.AccAddress) []byte {
	return append(append(append(QueuedFarmerKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(poolID)...), address.MustLengthPrefix(farmer)...)
}

// GetAllQueuedFarmerKey returns the store key to retrieve all queued farmers.
func GetAllQueuedFarmersKey(appID, poolID uint64) []byte {
	return append(append(QueuedFarmerKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(poolID)...)
}

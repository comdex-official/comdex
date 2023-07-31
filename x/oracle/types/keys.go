package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"

	"github.com/comdex-official/comdex/types"
)

const (
	// ModuleName is the name of the oracle module
	ModuleName = "oracle"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// RouterKey is the message route
	RouterKey = ModuleName

	// QuerierRoute is the query router key for the oracle module
	QuerierRoute = ModuleName
)

// KVStore key prefixes
var (
	KeyPrefixExchangeRate                 = []byte{0x01} // prefix for each key to a rate
	KeyPrefixFeederDelegation             = []byte{0x02} // prefix for each key to a feeder delegation
	KeyPrefixMissCounter                  = []byte{0x03} // prefix for each key to a miss counter
	KeyPrefixAggregateExchangeRatePrevote = []byte{0x04} // prefix for each key to a aggregate prevote
	KeyPrefixAggregateExchangeRateVote    = []byte{0x05} // prefix for each key to a aggregate vote
	KeyPrefixMedian                       = []byte{0x06} // prefix for each key to a price median
	KeyPrefixMedianDeviation              = []byte{0x07} // prefix for each key to a price median standard deviation
	KeyPrefixHistoricPrice                = []byte{0x08} // prefix for each key to a historic price
	KeyPrefixValidatorRewardSet           = []byte{0x09} // prefix for each key to a validator reward set
)

// GetExchangeRateKey - stored by *denom*
func GetExchangeRateKey(denom string) (key []byte) {
	key = append(key, KeyPrefixExchangeRate...)
	key = append(key, []byte(denom)...)
	return append(key, 0) // append 0 for null-termination
}

// GetFeederDelegationKey - stored by *Validator* address
func GetFeederDelegationKey(v sdk.ValAddress) (key []byte) {
	key = append(key, KeyPrefixFeederDelegation...)
	return append(key, address.MustLengthPrefix(v)...)
}

// GetMissCounterKey - stored by *Validator* address
func GetMissCounterKey(v sdk.ValAddress) (key []byte) {
	key = append(key, KeyPrefixMissCounter...)
	return append(key, address.MustLengthPrefix(v)...)
}

// GetAggregateExchangeRatePrevoteKey - stored by *Validator* address
func GetAggregateExchangeRatePrevoteKey(v sdk.ValAddress) (key []byte) {
	key = append(key, KeyPrefixAggregateExchangeRatePrevote...)
	return append(key, address.MustLengthPrefix(v)...)
}

// GetAggregateExchangeRateVoteKey - stored by *Validator* address
func GetAggregateExchangeRateVoteKey(v sdk.ValAddress) (key []byte) {
	key = append(key, KeyPrefixAggregateExchangeRateVote...)
	return append(key, address.MustLengthPrefix(v)...)
}

// KeyMedian - stored by *denom*
func KeyMedian(denom string, blockNum uint64) (key []byte) {
	return types.ConcatBytes(0, KeyPrefixMedian, []byte(denom), types.UintWithNullPrefix(blockNum))
}

// KeyMedianDeviation - stored by *denom*
func KeyMedianDeviation(denom string, blockNum uint64) (key []byte) {
	return types.ConcatBytes(0, KeyPrefixMedianDeviation, []byte(denom), types.UintWithNullPrefix(blockNum))
}

// KeyHistoricPrice - stored by *denom* and *block*
func KeyHistoricPrice(denom string, blockNum uint64) (key []byte) {
	return types.ConcatBytes(0, KeyPrefixHistoricPrice, []byte(denom), types.UintWithNullPrefix(blockNum))
}

// KeyValidatorRewardSet - stored by *block*
func KeyValidatorRewardSet() (key []byte) {
	return types.ConcatBytes(0, KeyPrefixValidatorRewardSet)
}

// ParseDenomAndBlockFromKey returns the denom and block contained in the *key*
// that has a uint64 at the end with a null prefix (length 9).
func ParseDenomAndBlockFromKey(key []byte, prefix []byte) (string, uint64) {
	return string(key[len(prefix) : len(key)-9]), binary.LittleEndian.Uint64(key[len(key)-8:])
}

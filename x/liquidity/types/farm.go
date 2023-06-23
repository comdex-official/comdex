package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewActivefarmer returns a new active farmer object.
func NewActivefarmer(appID, poolID uint64, farmer sdk.AccAddress, poolCoin sdk.Coin) ActiveFarmer {
	return ActiveFarmer{
		AppId:          appID,
		PoolId:         poolID,
		Farmer:         farmer.String(),
		FarmedPoolCoin: poolCoin,
	}
}

// Validate validates ActiveFarmer.
func (activeFarmer ActiveFarmer) Validate() error {
	if activeFarmer.AppId == 0 {
		return fmt.Errorf("app id must not be 0")
	}
	if activeFarmer.PoolId == 0 {
		return fmt.Errorf("pool id must not be 0")
	}
	if _, err := sdk.AccAddressFromBech32(activeFarmer.Farmer); err != nil {
		return fmt.Errorf("invalid farmer address %s: %w", activeFarmer.Farmer, err)
	}
	if err := activeFarmer.FarmedPoolCoin.Validate(); err != nil {
		return fmt.Errorf("invalid farmed-pool-coin %s: %w", activeFarmer.FarmedPoolCoin, err)
	}
	if activeFarmer.FarmedPoolCoin.IsZero() {
		return fmt.Errorf("farmed-pool-coin must not be 0")
	}
	return nil
}

// MustMarshalActiveFarmer returns the active farmer bytes.
// It throws panic if it fails.
func MustMarshalActiveFarmer(cdc codec.BinaryCodec, activeFarmer ActiveFarmer) []byte {
	return cdc.MustMarshal(&activeFarmer)
}

// MustUnmarshalActiveFarmer return the unmarshalled active farmer from bytes.
// It throws panic if it fails.
func MustUnmarshalActiveFarmer(cdc codec.BinaryCodec, value []byte) ActiveFarmer {
	activeFarmer, err := UnmarshalActiveFarmer(cdc, value)
	if err != nil {
		panic(err)
	}

	return activeFarmer
}

// UnmarshalActiveFarmer returns the active farmer from bytes.
func UnmarshalActiveFarmer(cdc codec.BinaryCodec, value []byte) (activeFarmer ActiveFarmer, err error) {
	err = cdc.Unmarshal(value, &activeFarmer)
	return activeFarmer, err
}

// NewQueuedfarmer returns a new queued farmer object.
func NewQueuedfarmer(appID, poolID uint64, farmer sdk.AccAddress) QueuedFarmer {
	return QueuedFarmer{
		AppId:      appID,
		PoolId:     poolID,
		Farmer:     farmer.String(),
		QueudCoins: []*QueuedCoin{},
	}
}

// Validate validates QueuedFarmer.
func (queuedFarmer QueuedFarmer) Validate() error {
	if queuedFarmer.AppId == 0 {
		return fmt.Errorf("app id must not be 0")
	}
	if queuedFarmer.PoolId == 0 {
		return fmt.Errorf("pool id must not be 0")
	}
	if _, err := sdk.AccAddressFromBech32(queuedFarmer.Farmer); err != nil {
		return fmt.Errorf("invalid farmer address %s: %w", queuedFarmer.Farmer, err)
	}
	return nil
}

// MustMarshalQueuedFarmer returns the queued farmer bytes.
// It throws panic if it fails.
func MustMarshalQueuedFarmer(cdc codec.BinaryCodec, queuedFarmer QueuedFarmer) []byte {
	return cdc.MustMarshal(&queuedFarmer)
}

// MustUnmarshalQueuedFarmer return the unmarshalled queued farmer from bytes.
// It throws panic if it fails.
func MustUnmarshalQueuedFarmer(cdc codec.BinaryCodec, value []byte) QueuedFarmer {
	queuedFarmer, err := UnmarshalQueuedFarmer(cdc, value)
	if err != nil {
		panic(err)
	}

	return queuedFarmer
}

// UnmarshalQueuedFarmer returns the queued farmer from bytes.
func UnmarshalQueuedFarmer(cdc codec.BinaryCodec, value []byte) (queuedFarmer QueuedFarmer, err error) {
	err = cdc.Unmarshal(value, &queuedFarmer)
	return queuedFarmer, err
}

package types

import (

	"github.com/cosmos/cosmos-sdk/codec"
)

const (
	VaultLiquidationsOffsetPrefix = "vault-liquidations"
)

// NewLiquidationOffsetHolder returns a new LiquidationOffsetHolder object.
func NewLiquidationOffsetHolder( currentOffset uint64) LiquidationOffsetHolder {
	return LiquidationOffsetHolder{
		CurrentOffset: currentOffset,
	}
}

func GetSliceStartEndForLiquidations(sliceLen, offset, batchSize int) (int, int) {
	if offset >= sliceLen || offset < 0 || batchSize < 0 {
		return sliceLen, sliceLen
	}
	start := offset
	end := offset + batchSize
	if end >= sliceLen {
		return start, sliceLen
	}
	return start, end
}

// Validate validates ActiveFarmer.
// func (liquidationOffsetHolder LiquidationOffsetHolder) Validate() error {
// 	if liquidationOffsetHolder.AppId == 0 {
// 		return fmt.Errorf("app id must not be 0")
// 	}
// 	return nil
// }

// MustMarshalLiquidationOffsetHolder returns the LiquidationOffsetHolder bytes.
// It throws panic if it fails.
func MustMarshalLiquidationOffsetHolder(cdc codec.BinaryCodec, liquidationOffsetHolder LiquidationOffsetHolder) []byte {
	return cdc.MustMarshal(&liquidationOffsetHolder)
}

// MustUnmarshalLiquidationOffsetHolder return the unmarshalledLiquidationOffsetHolder from bytes.
// It throws panic if it fails.
func MustUnmarshalLiquidationOffsetHolder(cdc codec.BinaryCodec, value []byte) LiquidationOffsetHolder {
	liquidationOffsetHolder, err := UnmarshalLiquidationOffsetHolder(cdc, value)
	if err != nil {
		panic(err)
	}

	return liquidationOffsetHolder
}

// UnmarshalLiquidationOffsetHolder returns the current offeset from bytes.
func UnmarshalLiquidationOffsetHolder(cdc codec.BinaryCodec, value []byte) (liquidationOffsetHolder LiquidationOffsetHolder, err error) {
	err = cdc.Unmarshal(value, &liquidationOffsetHolder)
	return liquidationOffsetHolder, err
}

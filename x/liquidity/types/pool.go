package types

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	time "time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/liquidity/amm"
)

const (
	DefaultSwapFeeDistributionDuration time.Duration = time.Hour * 24
)

var (
	_ amm.OrderSource = (*BasicPoolOrderSource)(nil)

	poolCoinDenomRegexp = regexp.MustCompile(`^pool([1-9]-[1-9]\d*)$`)
)

// PoolReserveAddress returns a unique pool reserve account address for each pool.
func PoolReserveAddress(appID, poolID uint64) sdk.AccAddress {
	return DeriveAddress(
		AddressType32Bytes,
		ModuleName,
		strings.Join([]string{PoolReserveAddressPrefix, strconv.FormatUint(appID, 10), strconv.FormatUint(poolID, 10)}, ModuleAddressNameSplitter),
	)
}

// NewPool returns a new pool object.
func NewPool(appID, id, pairID uint64) Pool {
	return Pool{
		Id:                    id,
		PairId:                pairID,
		ReserveAddress:        PoolReserveAddress(appID, id).String(),
		PoolCoinDenom:         PoolCoinDenom(appID, id),
		LastDepositRequestId:  0,
		LastWithdrawRequestId: 0,
		Disabled:              false,
		AppId:                 appID,
	}
}

// PoolCoinDenom returns a unique pool coin denom for a pool.
func PoolCoinDenom(appID, poolID uint64) string {
	return fmt.Sprintf("pool%d-%d", appID, poolID)
}

// ParsePoolCoinDenom parses a pool coin denom and returns the pool id.
func ParsePoolCoinDenom(denom string) (appID, poolID uint64, err error) {
	chunks := poolCoinDenomRegexp.FindStringSubmatch(denom)
	if len(chunks) == 0 {
		return 0, 0, fmt.Errorf("%s is not a pool coin denom", denom)
	}
	appID, err = strconv.ParseUint(strings.Split(chunks[1], "-")[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parse app id: %w", err)
	}
	poolID, err = strconv.ParseUint(strings.Split(chunks[1], "-")[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parse pool id: %w", err)
	}
	return appID, poolID, nil
}

func (pool Pool) GetReserveAddress() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(pool.ReserveAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// Validate validates Pool for genesis.
func (pool Pool) Validate() error {
	if pool.Id == 0 {
		return fmt.Errorf("pool id must not be 0")
	}
	if pool.PairId == 0 {
		return fmt.Errorf("pair id must not be 0")
	}
	if _, err := sdk.AccAddressFromBech32(pool.ReserveAddress); err != nil {
		return fmt.Errorf("invalid reserve address %s: %w", pool.ReserveAddress, err)
	}
	if err := sdk.ValidateDenom(pool.PoolCoinDenom); err != nil {
		return fmt.Errorf("invalid pool coin denom: %w", err)
	}
	return nil
}

// BasicPoolOrderSource is the order source for a pool which implements
// amm.OrderSource.
type BasicPoolOrderSource struct {
	amm.Pool
	PoolID                        uint64
	PoolReserveAddress            sdk.AccAddress
	BaseCoinDenom, QuoteCoinDenom string
}

// NewBasicPoolOrderSource returns a new BasicPoolOrderSource.
func NewBasicPoolOrderSource(
	pool amm.Pool,
	poolID uint64,
	//nolint
	reserveAddr sdk.AccAddress,
	baseCoinDenom, quoteCoinDenom string,
) *BasicPoolOrderSource {
	return &BasicPoolOrderSource{
		Pool:               pool,
		PoolID:             poolID,
		PoolReserveAddress: reserveAddr,
		BaseCoinDenom:      baseCoinDenom,
		QuoteCoinDenom:     quoteCoinDenom,
	}
}

func (os *BasicPoolOrderSource) BuyOrdersOver(price sdk.Dec) []amm.Order {
	amt := os.BuyAmountOver(price)
	if IsTooSmallOrderAmount(amt, price) {
		return nil
	}
	quoteCoin := sdk.NewCoin(os.QuoteCoinDenom, amm.OfferCoinAmount(amm.Buy, price, amt))
	return []amm.Order{NewPoolOrder(os.PoolID, os.PoolReserveAddress, amm.Buy, price, amt, quoteCoin, os.BaseCoinDenom)}
}

func (os *BasicPoolOrderSource) SellOrdersUnder(price sdk.Dec) []amm.Order {
	amt := os.SellAmountUnder(price)
	if IsTooSmallOrderAmount(amt, price) {
		return nil
	}
	baseCoin := sdk.NewCoin(os.BaseCoinDenom, amt)
	return []amm.Order{NewPoolOrder(os.PoolID, os.PoolReserveAddress, amm.Sell, price, amt, baseCoin, os.QuoteCoinDenom)}
}

// MustMarshalPool returns the pool bytes.
// It throws panic if it fails.
func MustMarshalPool(cdc codec.BinaryCodec, pool Pool) []byte {
	return cdc.MustMarshal(&pool)
}

// MustUnmarshalPool return the unmarshalled pool from bytes.
// It throws panic if it fails.
func MustUnmarshalPool(cdc codec.BinaryCodec, value []byte) Pool {
	pool, err := UnmarshalPool(cdc, value)
	if err != nil {
		panic(err)
	}

	return pool
}

// UnmarshalPool returns the pool from bytes.
func UnmarshalPool(cdc codec.BinaryCodec, value []byte) (pool Pool, err error) {
	err = cdc.Unmarshal(value, &pool)
	return pool, err
}

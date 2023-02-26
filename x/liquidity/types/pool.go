package types

import (
	"fmt"
	math "math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/liquidity/amm"
)

const (
	DefaultSwapFeeDistributionDuration = time.Hour * 24
)

var (
	_ amm.Orderer = (*PoolOrderer)(nil)

	poolCoinDenomRegexp = regexp.MustCompile(`^pool([1-9]-[1-9]\d*)$`)
	farmCoinDenomRegexp = regexp.MustCompile(`^farm([1-9]-[1-9]\d*)$`)
)

type PoolTokenDeserializerKit struct {
	Pair                 Pair
	Pool                 Pool
	PoolCoinSupply       sdk.Int
	QuoteCoinPoolBalance sdk.Coin
	BaseCoinPoolBalance  sdk.Coin
	AmmPoolObject        amm.Pool
}

// PoolReserveAddress returns a unique pool reserve account address for each pool.
func PoolReserveAddress(appID, poolID uint64) sdk.AccAddress {
	return DeriveAddress(
		AddressType32Bytes,
		ModuleName,
		strings.Join([]string{PoolReserveAddressPrefix, strconv.FormatUint(appID, 10), strconv.FormatUint(poolID, 10)}, ModuleAddressNameSplitter),
	)
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

// FarmCoinDenom returns a unique farm coin denom for a pool.
func FarmCoinDenom(appID, poolID uint64) string {
	return fmt.Sprintf("farm%d-%d", appID, poolID)
}

// ParseFarmCoinDenom parses a farm coin denom and returns the app and pool id.
func ParseFarmCoinDenom(denom string) (appID, poolID uint64, err error) {
	chunks := farmCoinDenomRegexp.FindStringSubmatch(denom)
	if len(chunks) == 0 {
		return 0, 0, fmt.Errorf("%s is not a farm coin denom", denom)
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

func IsFarmCoinDenom(denom string) bool {
	return len(farmCoinDenomRegexp.FindStringSubmatch(denom)) != 0
}

// FarmCoinDenom returns a unique farm coin denom for a pool.
func FarmCoinDecimals(baseDecimals, quoteDecimals uint64) uint64 {
	baseExponent := math.Log10(float64(baseDecimals))
	quoteExponent := math.Log10(float64(quoteDecimals))
	if baseExponent+quoteExponent >= 18 {
		return uint64(math.Pow10(18))
	}
	return uint64(math.Pow10(int(baseExponent + quoteExponent)))
}

// NewBasicPool returns a new basic pool object.
func NewBasicPool(appID, id, pairID uint64, creator sdk.AccAddress, baseAsset, quoteAsset assettypes.Asset) Pool {
	return Pool{
		Type:                  PoolTypeBasic,
		Id:                    id,
		PairId:                pairID,
		Creator:               creator.String(),
		ReserveAddress:        PoolReserveAddress(appID, id).String(),
		PoolCoinDenom:         PoolCoinDenom(appID, id),
		LastDepositRequestId:  0,
		LastWithdrawRequestId: 0,
		Disabled:              false,
		AppId:                 appID,
		FarmCoin: &FarmCoin{
			Denom:    FarmCoinDenom(appID, id),
			Decimals: FarmCoinDecimals(baseAsset.Decimals.Uint64(), quoteAsset.Decimals.Uint64()),
		},
	}
}

// NewRangedPool returns a new ranged pool object.
func NewRangedPool(appID, id, pairId uint64, creator sdk.AccAddress, minPrice, maxPrice sdk.Dec, baseAsset, quoteAsset assettypes.Asset) Pool {
	return Pool{
		Type:                  PoolTypeRanged,
		Id:                    id,
		PairId:                pairId,
		Creator:               creator.String(),
		ReserveAddress:        PoolReserveAddress(appID, id).String(),
		PoolCoinDenom:         PoolCoinDenom(appID, id),
		MinPrice:              &minPrice,
		MaxPrice:              &maxPrice,
		LastDepositRequestId:  0,
		LastWithdrawRequestId: 0,
		Disabled:              false,
		AppId:                 appID,
		FarmCoin: &FarmCoin{
			Denom:    FarmCoinDenom(appID, id),
			Decimals: FarmCoinDecimals(baseAsset.Decimals.Uint64(), quoteAsset.Decimals.Uint64()),
		},
	}
}

func (pool Pool) GetCreator() sdk.AccAddress {
	if pool.Creator == "" {
		return nil
	}
	addr, err := sdk.AccAddressFromBech32(pool.Creator)
	if err != nil {
		panic(err)
	}
	return addr
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

// AMMPool constructs amm.Pool interface from Pool.
func (pool Pool) AMMPool(rx, ry, ps sdk.Int) amm.Pool {
	switch pool.Type {
	case PoolTypeBasic:
		return amm.NewBasicPool(rx, ry, ps)
	case PoolTypeRanged:
		return amm.NewRangedPool(rx, ry, ps, *pool.MinPrice, *pool.MaxPrice)
	default:
		panic(fmt.Errorf("invalid pool type: %s", pool.Type))
	}
}

type PoolOrderer struct {
	amm.Pool
	Id                            uint64
	ReserveAddress                sdk.AccAddress
	BaseCoinDenom, QuoteCoinDenom string
}

// NewBasicPoolOrderSource returns a new BasicPoolOrderSource.
func NewPoolOrderer(
	pool amm.Pool,
	poolID uint64,
	reserveAddr sdk.AccAddress,
	baseCoinDenom, quoteCoinDenom string,
) *PoolOrderer {
	return &PoolOrderer{
		Pool:           pool,
		Id:             poolID,
		ReserveAddress: reserveAddr,
		BaseCoinDenom:  baseCoinDenom,
		QuoteCoinDenom: quoteCoinDenom,
	}
}

func (orderer *PoolOrderer) Order(dir amm.OrderDirection, price sdk.Dec, amt sdk.Int) amm.Order {
	var offerCoinDenom, demandCoinDenom string
	switch dir {
	case amm.Buy:
		offerCoinDenom, demandCoinDenom = orderer.QuoteCoinDenom, orderer.BaseCoinDenom
	case amm.Sell:
		offerCoinDenom, demandCoinDenom = orderer.BaseCoinDenom, orderer.QuoteCoinDenom
	}
	return NewPoolOrder(orderer.Id, orderer.ReserveAddress, dir, price, amt, offerCoinDenom, demandCoinDenom)
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

package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/liquidity/amm"
)

// OrderDirectionFromAMM converts amm.OrderDirection to liquidity module's
// OrderDirection.
func OrderDirectionFromAMM(dir amm.OrderDirection) OrderDirection {
	switch dir {
	case amm.Buy:
		return OrderDirectionBuy
	case amm.Sell:
		return OrderDirectionSell
	default:
		panic(fmt.Errorf("invalid order direction: %s", dir))
	}
}

// UserOrder is the user order type.
type UserOrder struct {
	*amm.BaseOrder
	Orderer                         sdk.AccAddress
	OrderID                         uint64
	BatchID                         uint64
	OfferCoinDenom, DemandCoinDenom string
}

// NewUserOrder returns a new user order.
func NewUserOrder(order Order) *UserOrder {
	var dir amm.OrderDirection
	var amt sdk.Int
	switch order.Direction {
	case OrderDirectionBuy:
		dir = amm.Buy
		utils.SafeMath(func() {
			amt = sdk.MinInt(
				order.OpenAmount,
				order.RemainingOfferCoin.Amount.ToDec().QuoTruncate(order.Price).TruncateInt(),
			)
		}, func() {
			amt = order.OpenAmount
		})
	case OrderDirectionSell:
		dir = amm.Sell
		amt = order.OpenAmount
	}
	return &UserOrder{
		BaseOrder:       amm.NewBaseOrder(dir, order.Price, amt, order.RemainingOfferCoin.Amount),
		Orderer:         order.GetOrderer(),
		OrderID:         order.Id,
		BatchID:         order.BatchId,
		OfferCoinDenom:  order.OfferCoin.Denom,
		DemandCoinDenom: order.ReceivedCoin.Denom,
	}
}

func (order *UserOrder) GetBatchID() uint64 {
	return order.BatchID
}

func (order *UserOrder) HasPriority(other amm.Order) bool {
	if !order.Amount.Equal(other.GetAmount()) {
		return order.BaseOrder.HasPriority(other)
	}
	switch other := other.(type) {
	case *UserOrder:
		return order.OrderID < other.OrderID
	case *PoolOrder:
		return true
	default:
		panic(fmt.Errorf("invalid order type: %T", other))
	}
}

func (order *UserOrder) String() string {
	return fmt.Sprintf("UserOrder(%d,%d,%s,%s,%s)",
		order.OrderID, order.BatchID, order.Direction, order.Price, order.Amount)
}

// PoolOrder is the pool order type.
type PoolOrder struct {
	*amm.BaseOrder
	PoolID                          uint64
	ReserveAddress                  sdk.AccAddress
	OfferCoinDenom, DemandCoinDenom string
}

// NewPoolOrder returns a new pool order.
func NewPoolOrder(
	poolID uint64,

	reserveAddr sdk.AccAddress,
	dir amm.OrderDirection,
	price sdk.Dec,
	amt sdk.Int,
	offerCoinDenom, demandCoinDenom string,
) *PoolOrder {
	return &PoolOrder{
		BaseOrder:       amm.NewBaseOrder(dir, price, amt, amm.OfferCoinAmount(dir, price, amt)),
		PoolID:          poolID,
		ReserveAddress:  reserveAddr,
		OfferCoinDenom:  offerCoinDenom,
		DemandCoinDenom: demandCoinDenom,
	}
}

func (order *PoolOrder) HasPriority(other amm.Order) bool {
	if !order.Amount.Equal(other.GetAmount()) {
		return order.BaseOrder.HasPriority(other)
	}
	switch other := other.(type) {
	case *UserOrder:
		return false
	case *PoolOrder:
		return order.PoolID < other.PoolID
	default:
		panic(fmt.Errorf("invalid order type: %T", other))
	}
}

func (order *PoolOrder) String() string {
	return fmt.Sprintf("PoolOrder(%d,%s,%s,%s)",
		order.PoolID, order.Direction, order.Price, order.Amount)
}

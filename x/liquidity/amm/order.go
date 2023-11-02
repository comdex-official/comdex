package amm

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
)

var (
	_ Order   = (*BaseOrder)(nil)
	_ Orderer = (*BaseOrderer)(nil)

	DefaultOrderer = BaseOrderer{}
)

// OrderDirection specifies an order direction, either buy or sell.
type OrderDirection int

// OrderDirection enumerations.
const (
	Buy OrderDirection = iota + 1
	Sell
)

func (dir OrderDirection) String() string {
	switch dir {
	case Buy:
		return "Buy"
	case Sell:
		return "Sell"
	default:
		return fmt.Sprintf("OrderDirection(%d)", dir)
	}
}

type Orderer interface {
	Order(dir OrderDirection, price sdkmath.LegacyDec, amt sdkmath.Int) Order
}

// BaseOrderer creates new BaseOrder with sufficient offer coin amount
// considering price and amount.
type BaseOrderer struct{}

func (orderer BaseOrderer) Order(dir OrderDirection, price sdkmath.LegacyDec, amt sdkmath.Int) Order {
	return NewBaseOrder(dir, price, amt, OfferCoinAmount(dir, price, amt))
}

// Order is the universal interface of an order.
type Order interface {
	GetDirection() OrderDirection
	// GetBatchId returns the batch id where the order was created.
	// Batch id of 0 means the current batch.
	GetBatchID() uint64
	GetPrice() sdkmath.LegacyDec
	GetAmount() sdkmath.Int // The original order amount
	GetOfferCoinAmount() sdkmath.Int
	GetPaidOfferCoinAmount() sdkmath.Int
	SetPaidOfferCoinAmount(amt sdkmath.Int)
	GetReceivedDemandCoinAmount() sdkmath.Int
	SetReceivedDemandCoinAmount(amt sdkmath.Int)
	GetOpenAmount() sdkmath.Int
	SetOpenAmount(amt sdkmath.Int)
	IsMatched() bool
	// HasPriority returns true if the order has higher priority
	// than the other order.
	HasPriority(other Order) bool
	String() string
}

// BaseOrder is the base struct for an Order.
type BaseOrder struct {
	Direction       OrderDirection
	Price           sdkmath.LegacyDec
	Amount          sdkmath.Int
	OfferCoinAmount sdkmath.Int

	// Match info
	OpenAmount               sdkmath.Int
	PaidOfferCoinAmount      sdkmath.Int
	ReceivedDemandCoinAmount sdkmath.Int
}

// NewBaseOrder returns a new BaseOrder.
func NewBaseOrder(dir OrderDirection, price sdkmath.LegacyDec, amt, offerCoinAmt sdkmath.Int) *BaseOrder {
	return &BaseOrder{
		Direction:                dir,
		Price:                    price,
		Amount:                   amt,
		OfferCoinAmount:          offerCoinAmt,
		OpenAmount:               amt,
		PaidOfferCoinAmount:      sdkmath.ZeroInt(),
		ReceivedDemandCoinAmount: sdkmath.ZeroInt(),
	}
}

// GetDirection returns the order direction.
func (order *BaseOrder) GetDirection() OrderDirection {
	return order.Direction
}

func (order *BaseOrder) GetBatchID() uint64 {
	return 0
}

// GetPrice returns the order price.
func (order *BaseOrder) GetPrice() sdkmath.LegacyDec {
	return order.Price
}

// GetAmount returns the order amount.
func (order *BaseOrder) GetAmount() sdkmath.Int {
	return order.Amount
}

func (order *BaseOrder) GetOfferCoinAmount() sdkmath.Int {
	return order.OfferCoinAmount
}

func (order *BaseOrder) GetPaidOfferCoinAmount() sdkmath.Int {
	return order.PaidOfferCoinAmount
}

func (order *BaseOrder) SetPaidOfferCoinAmount(amt sdkmath.Int) {
	order.PaidOfferCoinAmount = amt
}

func (order *BaseOrder) GetReceivedDemandCoinAmount() sdkmath.Int {
	return order.ReceivedDemandCoinAmount
}

func (order *BaseOrder) SetReceivedDemandCoinAmount(amt sdkmath.Int) {
	order.ReceivedDemandCoinAmount = amt
}

func (order *BaseOrder) GetOpenAmount() sdkmath.Int {
	return order.OpenAmount
}

func (order *BaseOrder) SetOpenAmount(amt sdkmath.Int) {
	order.OpenAmount = amt
}

func (order *BaseOrder) IsMatched() bool {
	return order.OpenAmount.LT(order.Amount)
}

// HasPriority returns whether the order has higher priority than
// the other order.
func (order *BaseOrder) HasPriority(other Order) bool {
	return order.Amount.GT(other.GetAmount())
}

func (order *BaseOrder) String() string {
	return fmt.Sprintf("BaseOrder(%s,%s,%s)", order.Direction, order.Price, order.Amount)
}

package types

import (
	"time"

	errorsmod "cosmossdk.io/errors"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/comdex-official/comdex/x/liquidity/amm"
)

var (
	_ sdk.Msg = (*MsgCreatePair)(nil)
	_ sdk.Msg = (*MsgCreatePool)(nil)
	_ sdk.Msg = (*MsgCreateRangedPool)(nil)
	_ sdk.Msg = (*MsgDeposit)(nil)
	_ sdk.Msg = (*MsgWithdraw)(nil)
	_ sdk.Msg = (*MsgLimitOrder)(nil)
	_ sdk.Msg = (*MsgMarketOrder)(nil)
	_ sdk.Msg = (*MsgMMOrder)(nil)
	_ sdk.Msg = (*MsgCancelOrder)(nil)
	_ sdk.Msg = (*MsgCancelAllOrders)(nil)
	_ sdk.Msg = (*MsgCancelMMOrder)(nil)
	_ sdk.Msg = (*MsgFarm)(nil)
	_ sdk.Msg = (*MsgUnfarm)(nil)
	_ sdk.Msg = (*MsgDepositAndFarm)(nil)
	_ sdk.Msg = (*MsgUnfarmAndWithdraw)(nil)
)

// Message types for the liquidity module.
const (
	TypeMsgCreatePair        = "create_pair"
	TypeMsgCreatePool        = "create_pool"
	TypeMsgCreateRangedPool  = "create_ranged_pool"
	TypeMsgDeposit           = "deposit"
	TypeMsgWithdraw          = "withdraw"
	TypeMsgLimitOrder        = "limit_order"
	TypeMsgMarketOrder       = "market_order"
	TypeMsgMMOrder           = "mm_order"
	TypeMsgCancelOrder       = "cancel_order"
	TypeMsgCancelAllOrders   = "cancel_all_orders"
	TypeMsgCancelMMOrder     = "cancel_mm_order"
	TypeMsgFarm              = "farm"
	TypeMsgUnfarm            = "unfarm"
	TypeMsgDepositAndFarm    = "deposit_and_farm"
	TypeMsgUnfarmAndWithdraw = "unfarm_and_withdraw"
)

// NewMsgCreatePair returns a new MsgCreatePair.
func NewMsgCreatePair(
	appID uint64,

	creator sdk.AccAddress,
	baseCoinDenom, quoteCoinDenom string,
) *MsgCreatePair {
	return &MsgCreatePair{
		AppId:          appID,
		Creator:        creator.String(),
		BaseCoinDenom:  baseCoinDenom,
		QuoteCoinDenom: quoteCoinDenom,
	}
}

func (msg MsgCreatePair) Route() string { return RouterKey }

func (msg MsgCreatePair) Type() string { return TypeMsgCreatePair }

func (msg MsgCreatePair) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address: %v", err)
	}
	if err := sdk.ValidateDenom(msg.BaseCoinDenom); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	if err := sdk.ValidateDenom(msg.QuoteCoinDenom); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	if msg.BaseCoinDenom == msg.QuoteCoinDenom {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "cannot use same denom for both base coin and quote coin")
	}
	return nil
}

func (msg MsgCreatePair) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgCreatePair) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgCreatePair) GetCreator() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgCreatePool creates a new MsgCreatePool.
func NewMsgCreatePool(
	appID uint64,

	creator sdk.AccAddress,
	pairID uint64,
	depositCoins sdk.Coins,
) *MsgCreatePool {
	return &MsgCreatePool{
		AppId:        appID,
		Creator:      creator.String(),
		PairId:       pairID,
		DepositCoins: depositCoins,
	}
}

func (msg MsgCreatePool) Route() string { return RouterKey }

func (msg MsgCreatePool) Type() string { return TypeMsgCreatePool }

func (msg MsgCreatePool) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address: %v", err)
	}
	if msg.PairId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "pair id must not be 0")
	}
	if err := msg.DepositCoins.Validate(); err != nil {
		return err
	}
	if len(msg.DepositCoins) != 2 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "wrong number of deposit coins: %d", len(msg.DepositCoins))
	}
	for _, coin := range msg.DepositCoins {
		if coin.Amount.GT(amm.MaxCoinAmount) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "deposit coin %s is bigger than the max amount %s", coin, amm.MaxCoinAmount)
		}
	}
	return nil
}

func (msg MsgCreatePool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgCreatePool) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgCreatePool) GetCreator() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgCreateRangedPool creates a new MsgCreateRangedPool.
func NewMsgCreateRangedPool(
	appID uint64,
	creator sdk.AccAddress,
	pairID uint64,
	depositCoins sdk.Coins,
	minPrice sdkmath.LegacyDec,
	maxPrice sdkmath.LegacyDec,
	initialPrice sdkmath.LegacyDec,
) *MsgCreateRangedPool {
	return &MsgCreateRangedPool{
		AppId:        appID,
		Creator:      creator.String(),
		PairId:       pairID,
		DepositCoins: depositCoins,
		MinPrice:     minPrice,
		MaxPrice:     maxPrice,
		InitialPrice: initialPrice,
	}
}

func (msg MsgCreateRangedPool) Route() string { return RouterKey }

func (msg MsgCreateRangedPool) Type() string { return TypeMsgCreateRangedPool }

func (msg MsgCreateRangedPool) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address: %v", err)
	}
	if msg.PairId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "pair id must not be 0")
	}
	if err := msg.DepositCoins.Validate(); err != nil {
		return err
	}
	if len(msg.DepositCoins) == 0 || len(msg.DepositCoins) > 2 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "wrong number of deposit coins: %d", len(msg.DepositCoins))
	}
	for _, coin := range msg.DepositCoins {
		if coin.Amount.GT(amm.MaxCoinAmount) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "deposit coin %s is bigger than the max amount %s", coin, amm.MaxCoinAmount)
		}
	}
	if err := amm.ValidateRangedPoolParams(msg.MinPrice, msg.MaxPrice, msg.InitialPrice); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	return nil
}

func (msg MsgCreateRangedPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgCreateRangedPool) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgCreateRangedPool) GetCreator() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgDeposit creates a new MsgDeposit.
func NewMsgDeposit(
	appID uint64,

	depositor sdk.AccAddress,
	poolID uint64,
	depositCoins sdk.Coins,
) *MsgDeposit {
	return &MsgDeposit{
		AppId:        appID,
		Depositor:    depositor.String(),
		PoolId:       poolID,
		DepositCoins: depositCoins,
	}
}

func (msg MsgDeposit) Route() string { return RouterKey }

func (msg MsgDeposit) Type() string { return TypeMsgDeposit }

func (msg MsgDeposit) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Depositor); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid depositor address: %v", err)
	}
	if msg.PoolId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "pool id must not be 0")
	}
	if err := msg.DepositCoins.Validate(); err != nil {
		return err
	}
	if len(msg.DepositCoins) == 0 || len(msg.DepositCoins) > 2 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "wrong number of deposit coins: %d", len(msg.DepositCoins))
	}
	return nil
}

func (msg MsgDeposit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgDeposit) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgDeposit) GetDepositor() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgWithdraw creates a new MsgWithdraw.
func NewMsgWithdraw(
	appID uint64,

	withdrawer sdk.AccAddress,
	poolID uint64,
	poolCoin sdk.Coin,
) *MsgWithdraw {
	return &MsgWithdraw{
		AppId:      appID,
		Withdrawer: withdrawer.String(),
		PoolId:     poolID,
		PoolCoin:   poolCoin,
	}
}

func (msg MsgWithdraw) Route() string { return RouterKey }

func (msg MsgWithdraw) Type() string { return TypeMsgWithdraw }

func (msg MsgWithdraw) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Withdrawer); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid withdrawer address: %v", err)
	}
	if msg.PoolId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "pool id must not be 0")
	}
	if err := msg.PoolCoin.Validate(); err != nil {
		return err
	}
	if !msg.PoolCoin.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "pool coin must be positive")
	}
	return nil
}

func (msg MsgWithdraw) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgWithdraw) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Withdrawer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgWithdraw) GetWithdrawer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Withdrawer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgLimitOrder creates a new MsgLimitOrder.
func NewMsgLimitOrder(
	appID uint64,

	orderer sdk.AccAddress,
	pairID uint64,
	dir OrderDirection,
	offerCoin sdk.Coin,
	demandCoinDenom string,
	price sdkmath.LegacyDec,
	amt sdkmath.Int,
	orderLifespan time.Duration,
) *MsgLimitOrder {
	return &MsgLimitOrder{
		AppId:           appID,
		Orderer:         orderer.String(),
		PairId:          pairID,
		Direction:       dir,
		OfferCoin:       offerCoin,
		DemandCoinDenom: demandCoinDenom,
		Price:           price,
		Amount:          amt,
		OrderLifespan:   orderLifespan,
	}
}

func (msg MsgLimitOrder) Route() string { return RouterKey }

func (msg MsgLimitOrder) Type() string { return TypeMsgLimitOrder }

func (msg MsgLimitOrder) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Orderer); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid orderer address: %v", err)
	}
	if msg.PairId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "pair id must not be 0")
	}
	if msg.Direction != OrderDirectionBuy && msg.Direction != OrderDirectionSell {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid order direction: %s", msg.Direction)
	}
	if err := sdk.ValidateDenom(msg.DemandCoinDenom); err != nil {
		return errorsmod.Wrap(err, "invalid demand coin denom")
	}
	if !msg.Price.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "price must be positive")
	}
	if err := msg.OfferCoin.Validate(); err != nil {
		return errorsmod.Wrap(err, "invalid offer coin")
	}
	if msg.OfferCoin.Amount.LT(amm.MinCoinAmount) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "offer coin %s is smaller than the min amount %s", msg.OfferCoin, amm.MinCoinAmount)
	}
	if msg.OfferCoin.Amount.GT(amm.MaxCoinAmount) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "offer coin %s is bigger than the max amount %s", msg.OfferCoin, amm.MaxCoinAmount)
	}
	if msg.Amount.LT(amm.MinCoinAmount) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "order amount %s is smaller than the min amount %s", msg.Amount, amm.MinCoinAmount)
	}
	if msg.Amount.GT(amm.MaxCoinAmount) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "order amount %s is bigger than the max amount %s", msg.Amount, amm.MaxCoinAmount)
	}
	var minOfferCoin sdk.Coin
	switch msg.Direction {
	case OrderDirectionBuy:
		minOfferCoin = sdk.NewCoin(msg.OfferCoin.Denom, amm.OfferCoinAmount(amm.Buy, msg.Price, msg.Amount))
	case OrderDirectionSell:
		minOfferCoin = sdk.NewCoin(msg.OfferCoin.Denom, msg.Amount)
	}
	if msg.OfferCoin.IsLT(minOfferCoin) {
		return errorsmod.Wrapf(ErrInsufficientOfferCoin, "%s is less than %s", msg.OfferCoin, minOfferCoin)
	}
	if msg.OfferCoin.Denom == msg.DemandCoinDenom {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "offer coin denom and demand coin denom must not be same")
	}
	if msg.OrderLifespan < 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "order lifespan must not be negative: %s", msg.OrderLifespan)
	}
	return nil
}

func (msg MsgLimitOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgLimitOrder) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Orderer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgLimitOrder) GetOrderer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Orderer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgMarketOrder creates a new MsgMarketOrder.
func NewMsgMarketOrder(
	appID uint64,

	orderer sdk.AccAddress,
	pairID uint64,
	dir OrderDirection,
	offerCoin sdk.Coin,
	demandCoinDenom string,
	amt sdkmath.Int,
	orderLifespan time.Duration,
) *MsgMarketOrder {
	return &MsgMarketOrder{
		AppId:           appID,
		Orderer:         orderer.String(),
		PairId:          pairID,
		Direction:       dir,
		OfferCoin:       offerCoin,
		DemandCoinDenom: demandCoinDenom,
		Amount:          amt,
		OrderLifespan:   orderLifespan,
	}
}

func (msg MsgMarketOrder) Route() string { return RouterKey }

func (msg MsgMarketOrder) Type() string { return TypeMsgMarketOrder }

func (msg MsgMarketOrder) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Orderer); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid orderer address: %v", err)
	}
	if msg.PairId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "pair id must not be 0")
	}
	if msg.Direction != OrderDirectionBuy && msg.Direction != OrderDirectionSell {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid order direction: %s", msg.Direction)
	}
	if err := sdk.ValidateDenom(msg.DemandCoinDenom); err != nil {
		return errorsmod.Wrap(err, "invalid demand coin denom")
	}
	if err := msg.OfferCoin.Validate(); err != nil {
		return errorsmod.Wrap(err, "invalid offer coin")
	}
	if msg.OfferCoin.Amount.LT(amm.MinCoinAmount) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "offer coin %s is smaller than the min amount %s", msg.OfferCoin, amm.MinCoinAmount)
	}
	if msg.OfferCoin.Amount.GT(amm.MaxCoinAmount) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "offer coin %s is bigger than the max amount %s", msg.OfferCoin, amm.MaxCoinAmount)
	}
	if msg.Amount.LT(amm.MinCoinAmount) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "order amount %s is smaller than the min amount %s", msg.Amount, amm.MinCoinAmount)
	}
	if msg.Amount.GT(amm.MaxCoinAmount) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "order amount %s is bigger than the max amount %s", msg.Amount, amm.MaxCoinAmount)
	}
	if msg.OfferCoin.Denom == msg.DemandCoinDenom {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "offer coin denom and demand coin denom must not be same")
	}
	if msg.OrderLifespan < 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "order lifespan must not be negative: %s", msg.OrderLifespan)
	}
	return nil
}

func (msg MsgMarketOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgMarketOrder) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Orderer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgMarketOrder) GetOrderer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Orderer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgMMOrder creates a new MsgMMOrder.
func NewMsgMMOrder(
	appID uint64,
	orderer sdk.AccAddress,
	pairID uint64,
	maxSellPrice, minSellPrice sdkmath.LegacyDec, sellAmt sdkmath.Int,
	maxBuyPrice, minBuyPrice sdkmath.LegacyDec, buyAmt sdkmath.Int,
	orderLifespan time.Duration,
) *MsgMMOrder {
	return &MsgMMOrder{
		AppId:         appID,
		Orderer:       orderer.String(),
		PairId:        pairID,
		MaxSellPrice:  maxSellPrice,
		MinSellPrice:  minSellPrice,
		SellAmount:    sellAmt,
		MaxBuyPrice:   maxBuyPrice,
		MinBuyPrice:   minBuyPrice,
		BuyAmount:     buyAmt,
		OrderLifespan: orderLifespan,
	}
}

func (msg MsgMMOrder) Route() string { return RouterKey }

func (msg MsgMMOrder) Type() string { return TypeMsgMMOrder }

func (msg MsgMMOrder) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Orderer); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid orderer address: %v", err)
	}
	if msg.PairId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "pair id must not be 0")
	}
	if msg.SellAmount.IsZero() && msg.BuyAmount.IsZero() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "sell amount and buy amount must not be zero at the same time")
	}
	if !msg.SellAmount.IsZero() {
		if msg.SellAmount.LT(amm.MinCoinAmount) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "sell amount %s is smaller than the min amount %s", msg.SellAmount, amm.MinCoinAmount)
		}
		if !msg.MaxSellPrice.IsPositive() {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "max sell price must be positive: %s", msg.MaxSellPrice)
		}
		if !msg.MinSellPrice.IsPositive() {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "min sell price must be positive: %s", msg.MinSellPrice)
		}
		if msg.MinSellPrice.GT(msg.MaxSellPrice) {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "max sell price must not be lower than min sell price")
		}
	}
	if !msg.BuyAmount.IsZero() {
		if msg.BuyAmount.LT(amm.MinCoinAmount) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "buy amount %s is smaller than the min amount %s", msg.BuyAmount, amm.MinCoinAmount)
		}
		if !msg.MinBuyPrice.IsPositive() {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "min buy price must be positive: %s", msg.MinBuyPrice)
		}
		if !msg.MaxBuyPrice.IsPositive() {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "max buy price must be positive: %s", msg.MaxBuyPrice)
		}
		if msg.MinBuyPrice.GT(msg.MaxBuyPrice) {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "max buy price must not be lower than min buy price")
		}
	}
	if msg.OrderLifespan < 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "order lifespan must not be negative: %s", msg.OrderLifespan)
	}
	return nil
}

func (msg MsgMMOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgMMOrder) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Orderer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgMMOrder) GetOrderer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Orderer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgCancelOrder creates a new MsgCancelOrder.
func NewMsgCancelOrder(
	appID uint64,

	orderer sdk.AccAddress,
	pairID uint64,
	orderID uint64,
) *MsgCancelOrder {
	return &MsgCancelOrder{
		AppId:   appID,
		OrderId: orderID,
		PairId:  pairID,
		Orderer: orderer.String(),
	}
}

func (msg MsgCancelOrder) Route() string { return RouterKey }

func (msg MsgCancelOrder) Type() string { return TypeMsgCancelOrder }

func (msg MsgCancelOrder) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Orderer); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid orderer address: %v", err)
	}
	if msg.PairId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "pair id must not be 0")
	}
	if msg.OrderId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "order id must not be 0")
	}
	return nil
}

func (msg MsgCancelOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgCancelOrder) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Orderer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgCancelOrder) GetOrderer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Orderer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgCancelAllOrders creates a new MsgCancelAllOrders.
func NewMsgCancelAllOrders(
	appID uint64,

	orderer sdk.AccAddress,
	pairIDs []uint64,
) *MsgCancelAllOrders {
	return &MsgCancelAllOrders{
		AppId:   appID,
		Orderer: orderer.String(),
		PairIds: pairIDs,
	}
}

func (msg MsgCancelAllOrders) Route() string { return RouterKey }

func (msg MsgCancelAllOrders) Type() string { return TypeMsgCancelAllOrders }

func (msg MsgCancelAllOrders) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Orderer); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid orderer address: %v", err)
	}
	pairIDSet := map[uint64]struct{}{}
	for _, pairID := range msg.PairIds {
		if pairID == 0 {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "pair id must not be 0")
		}
		if _, ok := pairIDSet[pairID]; ok {
			return ErrDuplicatePairID
		}
		pairIDSet[pairID] = struct{}{}
	}
	return nil
}

func (msg MsgCancelAllOrders) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgCancelAllOrders) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Orderer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgCancelAllOrders) GetOrderer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Orderer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgCancelMMOrder creates a new MsgCancelMMOrder.
func NewMsgCancelMMOrder(
	appID uint64,
	orderer sdk.AccAddress,
	pairID uint64,
) *MsgCancelMMOrder {
	return &MsgCancelMMOrder{
		AppId:   appID,
		Orderer: orderer.String(),
		PairId:  pairID,
	}
}

func (msg MsgCancelMMOrder) Route() string { return RouterKey }

func (msg MsgCancelMMOrder) Type() string { return TypeMsgCancelMMOrder }

func (msg MsgCancelMMOrder) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Orderer); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid orderer address: %v", err)
	}
	if msg.PairId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "pair id must not be 0")
	}
	return nil
}

func (msg MsgCancelMMOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgCancelMMOrder) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Orderer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgCancelMMOrder) GetOrderer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Orderer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgFarm creates a new MsgFarm.
func NewMsgFarm(
	appID uint64,
	poolID uint64,

	farmer sdk.AccAddress,
	poolCoin sdk.Coin,
) *MsgFarm {
	return &MsgFarm{
		AppId:           appID,
		PoolId:          poolID,
		Farmer:          farmer.String(),
		FarmingPoolCoin: poolCoin,
	}
}

func (msg MsgFarm) Route() string { return RouterKey }

func (msg MsgFarm) Type() string { return TypeMsgFarm }

func (msg MsgFarm) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Farmer); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid withdrawer address: %v", err)
	}
	if msg.PoolId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "pool id must not be 0")
	}
	if msg.AppId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "app id must not be 0")
	}
	if err := msg.FarmingPoolCoin.Validate(); err != nil {
		return err
	}
	if !msg.FarmingPoolCoin.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "coin must be positive")
	}
	return nil
}

func (msg MsgFarm) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgFarm) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgFarm) GetFarmer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgUnfarm creates a new MsgUnfarm.
func NewMsgUnfarm(
	appID uint64,
	poolID uint64,

	farmer sdk.AccAddress,
	poolCoin sdk.Coin,
) *MsgUnfarm {
	return &MsgUnfarm{
		AppId:             appID,
		PoolId:            poolID,
		Farmer:            farmer.String(),
		UnfarmingPoolCoin: poolCoin,
	}
}

func (msg MsgUnfarm) Route() string { return RouterKey }

func (msg MsgUnfarm) Type() string { return TypeMsgUnfarm }

func (msg MsgUnfarm) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Farmer); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid withdrawer address: %v", err)
	}
	if msg.PoolId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "pool id must not be 0")
	}
	if msg.AppId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "app id must not be 0")
	}
	if err := msg.UnfarmingPoolCoin.Validate(); err != nil {
		return err
	}
	if !msg.UnfarmingPoolCoin.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "coin must be positive")
	}
	return nil
}

func (msg MsgUnfarm) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgUnfarm) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgUnfarm) GetFarmer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgDeposit creates a new MsgDepositAndFarm.
func NewMsgDepositAndFarm(
	appID uint64,
	depositor sdk.AccAddress,
	poolID uint64,
	depositCoins sdk.Coins,
) *MsgDepositAndFarm {
	return &MsgDepositAndFarm{
		AppId:        appID,
		Depositor:    depositor.String(),
		PoolId:       poolID,
		DepositCoins: depositCoins,
	}
}

func (msg MsgDepositAndFarm) Route() string { return RouterKey }

func (msg MsgDepositAndFarm) Type() string { return TypeMsgDeposit }

func (msg MsgDepositAndFarm) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Depositor); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid depositor address: %v", err)
	}
	if msg.PoolId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "pool id must not be 0")
	}
	if err := msg.DepositCoins.Validate(); err != nil {
		return err
	}
	if len(msg.DepositCoins) == 0 || len(msg.DepositCoins) > 2 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "wrong number of deposit coins: %d", len(msg.DepositCoins))
	}
	return nil
}

func (msg MsgDepositAndFarm) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgDepositAndFarm) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgDepositAndFarm) GetDepositor() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgUnfarmAndWithdraw creates a new MsgUnfarmAndWithdraw.
func NewMsgUnfarmAndWithdraw(
	appID uint64,
	poolID uint64,
	farmer sdk.AccAddress,
	poolCoin sdk.Coin,
) *MsgUnfarmAndWithdraw {
	return &MsgUnfarmAndWithdraw{
		AppId:             appID,
		PoolId:            poolID,
		Farmer:            farmer.String(),
		UnfarmingPoolCoin: poolCoin,
	}
}

func (msg MsgUnfarmAndWithdraw) Route() string { return RouterKey }

func (msg MsgUnfarmAndWithdraw) Type() string { return TypeMsgUnfarm }

func (msg MsgUnfarmAndWithdraw) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Farmer); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid withdrawer address: %v", err)
	}
	if msg.PoolId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "pool id must not be 0")
	}
	if msg.AppId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "app id must not be 0")
	}
	if err := msg.UnfarmingPoolCoin.Validate(); err != nil {
		return err
	}
	if !msg.UnfarmingPoolCoin.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "coin must be positive")
	}
	return nil
}

func (msg MsgUnfarmAndWithdraw) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgUnfarmAndWithdraw) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgUnfarmAndWithdraw) GetFarmer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		panic(err)
	}
	return addr
}

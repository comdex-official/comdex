package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/comdex-official/comdex/x/liquidity/amm"
)

var (
	_ sdk.Msg = (*MsgCreatePair)(nil)
	_ sdk.Msg = (*MsgCreatePool)(nil)
	_ sdk.Msg = (*MsgDeposit)(nil)
	_ sdk.Msg = (*MsgWithdraw)(nil)
	_ sdk.Msg = (*MsgLimitOrder)(nil)
	_ sdk.Msg = (*MsgMarketOrder)(nil)
	_ sdk.Msg = (*MsgCancelOrder)(nil)
	_ sdk.Msg = (*MsgCancelAllOrders)(nil)
)

// Message types for the liquidity module.
const (
	TypeMsgCreatePair      = "create_pair"
	TypeMsgCreatePool      = "create_pool"
	TypeMsgDeposit         = "deposit"
	TypeMsgWithdraw        = "withdraw"
	TypeMsgLimitOrder      = "limit_order"
	TypeMsgMarketOrder     = "market_order"
	TypeMsgCancelOrder     = "cancel_order"
	TypeMsgCancelAllOrders = "cancel_all_orders"
)

// NewMsgCreatePair returns a new MsgCreatePair.
func NewMsgCreatePair(
	appID uint64,
	//nolint
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address: %v", err)
	}
	if err := sdk.ValidateDenom(msg.BaseCoinDenom); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	if err := sdk.ValidateDenom(msg.QuoteCoinDenom); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
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
	//nolint
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address: %v", err)
	}
	if msg.PairId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "pair id must not be 0")
	}
	if err := msg.DepositCoins.Validate(); err != nil {
		return err
	}
	if len(msg.DepositCoins) != 2 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "wrong number of deposit coins: %d", len(msg.DepositCoins))
	}
	for _, coin := range msg.DepositCoins {
		if coin.Amount.GT(amm.MaxCoinAmount) {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "deposit coin %s is bigger than the max amount %s", coin, amm.MaxCoinAmount)
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

// NewMsgDeposit creates a new MsgDeposit.
func NewMsgDeposit(
	appID uint64,
	//nolint
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid depositor address: %v", err)
	}
	if msg.PoolId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "pool id must not be 0")
	}
	if err := msg.DepositCoins.Validate(); err != nil {
		return err
	}
	if len(msg.DepositCoins) != 2 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "wrong number of deposit coins: %d", len(msg.DepositCoins))
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
	//nolint
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid withdrawer address: %v", err)
	}
	if msg.PoolId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "pool id must not be 0")
	}
	if err := msg.PoolCoin.Validate(); err != nil {
		return err
	}
	if !msg.PoolCoin.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "pool coin must be positive")
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
	//nolint
	orderer sdk.AccAddress,
	pairID uint64,
	dir OrderDirection,
	offerCoin sdk.Coin,
	demandCoinDenom string,
	price sdk.Dec,
	amt sdk.Int,
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid orderer address: %v", err)
	}
	if msg.PairId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "pair id must not be 0")
	}
	if msg.Direction != OrderDirectionBuy && msg.Direction != OrderDirectionSell {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid order direction: %s", msg.Direction)
	}
	if err := sdk.ValidateDenom(msg.DemandCoinDenom); err != nil {
		return sdkerrors.Wrap(err, "invalid demand coin denom")
	}
	if !msg.Price.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "price must be positive")
	}
	if err := msg.OfferCoin.Validate(); err != nil {
		return sdkerrors.Wrap(err, "invalid offer coin")
	}
	if msg.OfferCoin.Amount.LT(amm.MinCoinAmount) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "offer coin %s is smaller than the min amount %s", msg.OfferCoin, amm.MinCoinAmount)
	}
	if msg.OfferCoin.Amount.GT(amm.MaxCoinAmount) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "offer coin %s is bigger than the max amount %s", msg.OfferCoin, amm.MaxCoinAmount)
	}
	if msg.Amount.LT(amm.MinCoinAmount) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "order amount %s is smaller than the min amount %s", msg.Amount, amm.MinCoinAmount)
	}
	if msg.Amount.GT(amm.MaxCoinAmount) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "order amount %s is bigger than the max amount %s", msg.Amount, amm.MaxCoinAmount)
	}
	var minOfferCoin sdk.Coin
	switch msg.Direction {
	case OrderDirectionBuy:
		minOfferCoin = sdk.NewCoin(msg.OfferCoin.Denom, amm.OfferCoinAmount(amm.Buy, msg.Price, msg.Amount))
	case OrderDirectionSell:
		minOfferCoin = sdk.NewCoin(msg.OfferCoin.Denom, msg.Amount)
	}
	if msg.OfferCoin.IsLT(minOfferCoin) {
		return sdkerrors.Wrapf(ErrInsufficientOfferCoin, "%s is less than %s", msg.OfferCoin, minOfferCoin)
	}
	if msg.OfferCoin.Denom == msg.DemandCoinDenom {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "offer coin denom and demand coin denom must not be same")
	}
	if msg.OrderLifespan < 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "order lifespan must not be negative: %s", msg.OrderLifespan)
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
	//nolint
	orderer sdk.AccAddress,
	pairID uint64,
	dir OrderDirection,
	offerCoin sdk.Coin,
	demandCoinDenom string,
	amt sdk.Int,
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid orderer address: %v", err)
	}
	if msg.PairId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "pair id must not be 0")
	}
	if msg.Direction != OrderDirectionBuy && msg.Direction != OrderDirectionSell {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid order direction: %s", msg.Direction)
	}
	if err := sdk.ValidateDenom(msg.DemandCoinDenom); err != nil {
		return sdkerrors.Wrap(err, "invalid demand coin denom")
	}
	if err := msg.OfferCoin.Validate(); err != nil {
		return sdkerrors.Wrap(err, "invalid offer coin")
	}
	if msg.OfferCoin.Amount.LT(amm.MinCoinAmount) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "offer coin %s is smaller than the min amount %s", msg.OfferCoin, amm.MinCoinAmount)
	}
	if msg.OfferCoin.Amount.GT(amm.MaxCoinAmount) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "offer coin %s is bigger than the max amount %s", msg.OfferCoin, amm.MaxCoinAmount)
	}
	if msg.Amount.LT(amm.MinCoinAmount) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "order amount %s is smaller than the min amount %s", msg.Amount, amm.MinCoinAmount)
	}
	if msg.Amount.GT(amm.MaxCoinAmount) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "order amount %s is bigger than the max amount %s", msg.Amount, amm.MaxCoinAmount)
	}
	if msg.OfferCoin.Denom == msg.DemandCoinDenom {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "offer coin denom and demand coin denom must not be same")
	}
	if msg.OrderLifespan < 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "order lifespan must not be negative: %s", msg.OrderLifespan)
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

// NewMsgCancelOrder creates a new MsgCancelOrder.
func NewMsgCancelOrder(
	appID uint64,
	//nolint
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid orderer address: %v", err)
	}
	if msg.PairId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "pair id must not be 0")
	}
	if msg.OrderId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "order id must not be 0")
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
	//nolint
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid orderer address: %v", err)
	}
	pairIDSet := map[uint64]struct{}{}
	for _, pairID := range msg.PairIds {
		if pairID == 0 {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "pair id must not be 0")
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

// NewMsgSoftLock creates a new MsgTokensSoftLock.
func NewMsgSoftLock(
	appID uint64,
	//nolint
	depositor sdk.AccAddress,
	poolID uint64,
	softLockCoin sdk.Coin,
) *MsgTokensSoftLock {
	return &MsgTokensSoftLock{
		AppId:        appID,
		Depositor:    depositor.String(),
		PoolId:       poolID,
		SoftLockCoin: softLockCoin,
	}
}

func (msg MsgTokensSoftLock) Route() string { return RouterKey }

func (msg MsgTokensSoftLock) Type() string { return TypeMsgWithdraw }

func (msg MsgTokensSoftLock) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Depositor); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid withdrawer address: %v", err)
	}
	if msg.PoolId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "pool id must not be 0")
	}
	if err := msg.SoftLockCoin.Validate(); err != nil {
		return err
	}
	if !msg.SoftLockCoin.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "coin must be positive")
	}
	return nil
}

func (msg MsgTokensSoftLock) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgTokensSoftLock) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgTokensSoftLock) GetWithdrawer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgSoftUnlock creates a new MsgTokensSoftUnlock.
func NewMsgSoftUnlock(
	appID uint64,
	//nolint
	depositor sdk.AccAddress,
	poolID uint64,
	softUnlockCoin sdk.Coin,
) *MsgTokensSoftUnlock {
	return &MsgTokensSoftUnlock{
		AppId:          appID,
		Depositor:      depositor.String(),
		PoolId:         poolID,
		SoftUnlockCoin: softUnlockCoin,
	}
}

func (msg MsgTokensSoftUnlock) Route() string { return RouterKey }

func (msg MsgTokensSoftUnlock) Type() string { return TypeMsgWithdraw }

func (msg MsgTokensSoftUnlock) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Depositor); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid withdrawer address: %v", err)
	}
	if msg.PoolId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "pool id must not be 0")
	}
	if err := msg.SoftUnlockCoin.Validate(); err != nil {
		return err
	}
	if !msg.SoftUnlockCoin.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "coin must be positive")
	}
	return nil
}

func (msg MsgTokensSoftUnlock) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgTokensSoftUnlock) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgTokensSoftUnlock) GetWithdrawer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		panic(err)
	}
	return addr
}

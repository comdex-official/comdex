package keeper

import (
	"fmt"
	"strconv"

	"github.com/comdex-official/comdex/x/lend/expected"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/comdex-official/comdex/x/lend/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   sdk.StoreKey
		memKey     sdk.StoreKey
		paramstore paramtypes.Subspace
		bank       expected.BankKeeper
		account    expected.AccountKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	bank expected.BankKeeper,
	account expected.AccountKeeper,

) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{

		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		bank:       bank,
		account:    account,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) ModuleBalance(ctx sdk.Context, moduleName string, denom string) sdk.Int {
	return k.bank.GetBalance(ctx, authtypes.NewModuleAddress(moduleName), denom).Amount
}

func (k Keeper) LendAsset(ctx sdk.Context, lenderAddr sdk.AccAddress, PairId uint64, loan sdk.Coin) error {
	if !k.IsWhitelistedAsset(ctx, loan.Denom) {
		return sdkerrors.Wrap(types.ErrInvalidAsset, loan.String())
	}

	pair, found := k.GetPair(ctx, PairId)
	if !found {
		return sdkerrors.Wrap(types.ErrorPairDoesNotExist, strconv.Itoa(int(PairId)))
	}
	asset1, _ := k.GetAsset(ctx, pair.Asset_1)
	asset2, _ := k.GetAsset(ctx, pair.Asset_2)

	if loan.Denom != asset1.Denom && loan.Denom != asset2.Denom {
		return sdkerrors.Wrap(types.ErrBadOfferCoinAmount, loan.Denom)
	}

	userLendId := k.GetUserLendIDHistory(ctx)
	k.SetUserLendHistory(ctx, lenderAddr, loan, userLendId)
	k.SetUserLendIDHistory(ctx, userLendId+1)

	// send token balance to lend module account
	// update users lending
	//TODO:
	// update reserves
	// calculate interest rate
	loanTokens := sdk.NewCoins(loan)
	if err := k.bank.SendCoinsFromAccountToModule(ctx, lenderAddr, types.ModuleName, loanTokens); err != nil {
		return err
	}

	currentCollateral := k.GetCollateralAmount(ctx, lenderAddr, loan.Denom)
	if err := k.setCollateralAmount(ctx, lenderAddr, currentCollateral.Add(loan)); err != nil {
		return err
	}

	return nil
}

func (k Keeper) WithdrawAsset(ctx sdk.Context, lenderAddr sdk.AccAddress, withdrawal sdk.Coin) error {

	if !k.IsWhitelistedAsset(ctx, withdrawal.Denom) {
		return sdkerrors.Wrap(types.ErrInvalidAsset, withdrawal.String())
	}

	// Ensure module account has sufficient unreserved tokens to withdraw
	reservedAmount := k.GetReserveFunds(ctx, withdrawal.Denom)
	currentCollateral := k.GetCollateralAmount(ctx, lenderAddr, withdrawal.Denom)
	availableAmount := k.ModuleBalance(ctx, types.ModuleName, withdrawal.Denom)

	if withdrawal.Amount.GT(availableAmount.Sub(reservedAmount)) {
		return sdkerrors.Wrap(types.ErrLendingPoolInsufficient, withdrawal.String())
	}

	if withdrawal.Amount.GT(currentCollateral.Amount) {
		return sdkerrors.Wrap(types.ErrInsufficientBalance, withdrawal.String())
	}
	// update lenders share after withdraw
	if err := k.setCollateralAmount(ctx, lenderAddr, currentCollateral.Sub(withdrawal)); err != nil {
		return err
	}
	// send the base assets to lender
	tokens := sdk.NewCoins(withdrawal)
	if err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, lenderAddr, tokens); err != nil {
		return err
	}

	return nil
}

func (k Keeper) BorrowAsset(ctx sdk.Context, lenderAddr sdk.AccAddress, loan sdk.Coin) error {
	if !k.IsWhitelistedAsset(ctx, loan.Denom) {
		return sdkerrors.Wrap(types.ErrInvalidAsset, loan.String())
	}

	// send token balance to lend module account
	loanTokens := sdk.NewCoins(loan)
	if err := k.bank.SendCoinsFromAccountToModule(ctx, lenderAddr, types.ModuleName, loanTokens); err != nil {
		return err
	}

	return nil
}

func (k Keeper) RepayAsset(ctx sdk.Context, borrowerAddr sdk.AccAddress, payment sdk.Coin) (sdk.Int, error) {
	if !payment.IsValid() {
		return sdk.ZeroInt(), sdkerrors.Wrap(types.ErrInvalidAsset, payment.String())
	}

	if !k.IsWhitelistedAsset(ctx, payment.Denom) {
		return sdk.ZeroInt(), sdkerrors.Wrap(types.ErrInvalidAsset, payment.String())
	}

	return payment.Amount, nil
}

func (k Keeper) FundModAcc(ctx sdk.Context, moduleName string, lenderAddr sdk.AccAddress, payment sdk.Coin) error {
	if !k.IsWhitelistedAsset(ctx, payment.Denom) {
		return sdkerrors.Wrap(types.ErrInvalidAsset, payment.String())
	}

	loanTokens := sdk.NewCoins(payment)
	if err := k.bank.SendCoinsFromAccountToModule(ctx, lenderAddr, moduleName, loanTokens); err != nil {
		return err
	}

	currentCollateral := k.GetCollateralAmount(ctx, lenderAddr, payment.Denom)
	if err := k.setCollateralAmount(ctx, lenderAddr, currentCollateral.Add(payment)); err != nil {
		return err
	}

	return nil
}

func (k *Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

func (k *Keeper) SetUserLendIDHistory(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendHistoryPrefix
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}


func (k *Keeper) SetUserLendHistory(ctx sdk.Context, lenderAddr sdk.AccAddress, loan sdk.Coin, id uint64) {
	
	user_lend := types.LendHistory{
		Owner: lenderAddr.String(),
		Amount: &loan,
	}
	var (
		store = k.Store(ctx)
		key   = types.LendUserHistoryKey(id)
		value = k.cdc.MustMarshal(&user_lend)
	)
	store.Set(key, value)
}

func (k *Keeper) GetUserLendHistory(ctx sdk.Context, id uint64) (user_lend types.LendHistory, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LendUserHistoryKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return user_lend, false
	}

	k.cdc.MustUnmarshal(value, &user_lend)
	return user_lend, true
}

func (k *Keeper) GetUserLendIDHistory(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.LendHistoryPrefix
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}
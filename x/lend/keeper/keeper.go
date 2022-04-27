package keeper

import (
	"fmt"
	"github.com/comdex-official/comdex/x/lend/expected"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/comdex-official/comdex/x/lend/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
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

func (k Keeper) LendAsset(ctx sdk.Context, lenderAddr sdk.AccAddress, loan sdk.Coin) error {
	if !k.IsWhitelistedAsset(ctx, loan.Denom) {
		return sdkerrors.Wrap(types.ErrInvalidAsset, loan.String())
	}
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

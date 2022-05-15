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
		asset      expected.AssetKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	bank expected.BankKeeper,
	account expected.AccountKeeper,
	asset expected.AssetKeeper,

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
		asset:      asset,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) ModuleBalance(ctx sdk.Context, moduleName string, denom string) sdk.Int {
	return k.bank.GetBalance(ctx, authtypes.NewModuleAddress(moduleName), denom).Amount
}

func (k Keeper) LendAsset(ctx sdk.Context, lenderAddr sdk.AccAddress, pairID uint64, lent sdk.Coin) error {
	lentTokens := sdk.NewCoins(lent)

	ExtPairID, _ := k.asset.GetWhitelistPair(ctx, pairID)

	currentCollateral := k.GetCollateralAmount(ctx, lenderAddr, lent.Denom)
	if err := k.setCollateralAmount(ctx, lenderAddr, currentCollateral.Add(lent)); err != nil {
		return err
	}
	basePair, found := k.asset.GetPair(ctx, ExtPairID.PairId)
	if found != true {
		return types.ErrorPairDoesNotExist
	}

	Asset, found := k.asset.GetAsset(ctx, basePair.AssetIn)
	if found != true {
		return types.ErrorAssetDoesNotExist
	}

	if Asset.Denom != lent.Denom {
		return types.ErrInvalidAsset
	}

	cToken, err := k.ExchangeToken(ctx, lent, Asset.Name)
	if err != nil {
		return err
	}

	if err := k.bank.SendCoinsFromAccountToModule(ctx, lenderAddr, ExtPairID.ModuleAcc, lentTokens); err != nil {
		return err
	}
	// mint c/Token and set new total cToken supply

	cTokens := sdk.NewCoins(cToken)
	if err = k.bank.MintCoins(ctx, ExtPairID.ModuleAcc, cTokens); err != nil {
		return err
	}
	if err = k.setCTokenSupply(ctx, k.GetCTokenSupply(ctx, cToken.Denom).Add(cToken)); err != nil {
		return err
	}

	err = k.bank.SendCoinsFromModuleToAccount(ctx, ExtPairID.ModuleAcc, lenderAddr, cTokens)
	if err != nil {
		return err
	}

	Id := k.GetLendID(ctx)
	lendPositon := types.Lend_Asset{
		ID:          Id + 1,
		PairID:      pairID,
		Owner:       lenderAddr.String(),
		AmountIn:    lent,
		LendingTime: ctx.BlockTime(),
		Reward:      sdk.NewCoin("cCMDX", sdk.NewInt(0)),
	}

	k.SetLendID(ctx, lendPositon.ID)
	k.SetLend(ctx, lendPositon)
	return nil
}

func (k Keeper) WithdrawAsset(ctx sdk.Context, lenderAddr sdk.AccAddress, withdrawal sdk.Coin) error {

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

func (k Keeper) FundModAcc(ctx sdk.Context, moduleName string, lenderAddr sdk.AccAddress, payment sdk.Coin) error {
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

package keeper

import (
	"fmt"
	"github.com/comdex-official/comdex/x/lend/expected"
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

	/*if k.GetOracleValidationResult(ctx) == false{
		return nil, types.ErrorOraclePriceExpired
	}*/

	if k.HasLendForAddressByPair(ctx, lenderAddr, pairID) {
		return types.ErrorDuplicateLend
	}

	lentTokens := sdk.NewCoins(lent)

	LendPairID, found := k.GetLendPair(ctx, pairID)
	if found != true {
		return types.ErrorPairDoesNotExist
	}

	currentCollateral := k.GetCollateralAmount(ctx, lenderAddr, lent.Denom)
	if err := k.setCollateralAmount(ctx, lenderAddr, currentCollateral.Add(lent)); err != nil {
		return err
	}

	Asset, found := k.asset.GetAsset(ctx, LendPairID.AssetIn)
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

	if err := k.bank.SendCoinsFromAccountToModule(ctx, lenderAddr, LendPairID.ModuleAcc, lentTokens); err != nil {
		return err
	}
	// mint c/Token and set new total cToken supply

	cTokens := sdk.NewCoins(cToken)
	if err = k.bank.MintCoins(ctx, LendPairID.ModuleAcc, cTokens); err != nil {
		return err
	}
	if err = k.setCTokenSupply(ctx, k.GetCTokenSupply(ctx, cToken.Denom).Add(cToken)); err != nil {
		return err
	}

	err = k.bank.SendCoinsFromModuleToAccount(ctx, LendPairID.ModuleAcc, lenderAddr, cTokens)
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
	k.SetLendForAddressByPair(ctx, lenderAddr, pairID, lendPositon.ID)

	return nil
}

func (k Keeper) DepositAsset(ctx sdk.Context, lendID uint64, lenderAddr sdk.AccAddress, deposit sdk.Coin) error {

	//if k.GetOracleValidationResult(ctx) == false{
	//	return nil, types.ErrorOraclePriceExpired
	//}

	lend, found := k.GetLend(ctx, lendID)
	if !found {
		return types.ErrorLendDoesNotExist
	}
	if lenderAddr.String() != lend.Owner {
		return types.ErrorUnauthorized
	}

	pair, _ := k.GetLendPair(ctx, lend.PairID)
	if !found {
		return types.ErrorPairDoesNotExist
	}

	assetIn, found := k.GetAsset(ctx, pair.AssetIn)
	if !found {
		return types.ErrorAssetDoesNotExist
	}

	lend.AmountIn = lend.AmountIn.Add(deposit)
	if !lend.AmountIn.IsPositive() {
		return types.ErrorUnauthorized
	}

	if err := k.SendCoinFromAccountToModule(ctx, lenderAddr, pair.ModuleAcc, sdk.NewCoin(assetIn.Denom, deposit.Amount)); err != nil {
		return err
	}

	cToken, err := k.ExchangeToken(ctx, deposit, assetIn.Name)
	if err != nil {
		return err
	}

	cTokens := sdk.NewCoins(cToken)
	if err = k.bank.MintCoins(ctx, pair.ModuleAcc, cTokens); err != nil {
		return err
	}
	if err = k.setCTokenSupply(ctx, k.GetCTokenSupply(ctx, cToken.Denom).Add(cToken)); err != nil {
		return err
	}

	lend.AmountOut = lend.AmountOut.Add(cToken)

	err = k.bank.SendCoinsFromModuleToAccount(ctx, pair.ModuleAcc, lenderAddr, cTokens)

	k.SetLend(ctx, lend)
	return nil
}

func (k Keeper) WithdrawAsset(ctx sdk.Context, lendID uint64, lenderAddr sdk.AccAddress, withdrawal sdk.Coin) error {

	lend, found := k.GetLend(ctx, lendID)
	if !found {
		return types.ErrorLendDoesNotExist
	}
	if lenderAddr.String() != lend.Owner {
		return types.ErrorUnauthorized
	}

	pair, found := k.GetLendPair(ctx, lendID)
	if !found {
		return types.ErrorPairDoesNotExist
	}

	assetIn, found := k.GetAsset(ctx, pair.AssetIn)
	if !found {
		return types.ErrorAssetDoesNotExist
	}

	/*assetOut, found := k.GetAsset(ctx, pair.AssetOut)
	if !found {
		return types.ErrorAssetDoesNotExist
	}*/

	lend.AmountIn = lend.AmountIn.Sub(withdrawal)
	if !lend.AmountIn.IsPositive() {
		return types.ErrorInvalidAmount
	}

	/*if err := k.VerifyCollaterlizationRatio(ctx, lend.AmountIn, assetIn, lend.AmountOut, assetOut); err != nil {
		return err
	}*/

	if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, lenderAddr, sdk.NewCoin(assetIn.Denom, withdrawal.Amount)); err != nil {
		return err
	}

	k.SetLend(ctx, lend)
	return nil
}

func (k Keeper) BorrowAsset(ctx sdk.Context, lenderAddr sdk.AccAddress, borrow sdk.Coin) error {
	return nil
}

func (k Keeper) DrawAsset(ctx sdk.Context, borrowID uint64, lenderAddr sdk.AccAddress, draw sdk.Coin) error {
	return nil
}

func (k Keeper) RepayAsset(ctx sdk.Context, borrowID uint64, lenderAddr sdk.AccAddress, repay sdk.Coin) error {
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

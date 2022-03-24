package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/vault/types"
)

var (
	_ types.MsgServiceServer = (*msgServer)(nil)
)

type msgServer struct {
	Keeper
}

func NewMsgServiceServer(keeper Keeper) types.MsgServiceServer {
	return &msgServer{
		Keeper: keeper,
	}
}

func (k *msgServer) MsgCreate(c context.Context, msg *types.MsgCreateRequest) (*types.MsgCreateResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if k.GetOracleValidationResult(ctx) == false{
		return nil, types.ErrorOraclePriceExpired
	}

	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	if k.HasVaultForAddressByPair(ctx, from, msg.PairID) {
		return nil, types.ErrorDuplicateVault
	}

	pair, found := k.GetPair(ctx, msg.PairID)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}

	assetIn, found := k.GetAsset(ctx, pair.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	assetOut, found := k.GetAsset(ctx, pair.AssetOut)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	if err := k.VerifyCollaterlizationRatio(ctx, msg.AmountIn, assetIn, msg.AmountOut, assetOut, pair.LiquidationRatio); err != nil {
		return nil, err
	}

	if err := k.SendCoinFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoin(assetIn.Denom, msg.AmountIn)); err != nil {
		return nil, err
	}
	if err := k.MintCAssets(ctx, types.ModuleName, assetIn.Denom, assetOut.Denom, msg.AmountOut); err != nil {
		return nil, err
	}
	if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, from, sdk.NewCoin(assetOut.Denom, msg.AmountOut)); err != nil {
		return nil, err
	}
	k.CreteNewVault(ctx, msg.PairID, msg.From, assetIn, msg.AmountIn, assetOut, msg.AmountOut)
	return &types.MsgCreateResponse{}, nil
}

func (k *msgServer) MsgDeposit(c context.Context, msg *types.MsgDepositRequest) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if k.GetOracleValidationResult(ctx) == false{
		return nil, types.ErrorOraclePriceExpired
	}

	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	vault, found := k.GetVault(ctx, msg.ID)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	if msg.From != vault.Owner {
		return nil, types.ErrorUnauthorized
	}

	pair, found := k.GetPair(ctx, vault.PairID)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}

	assetIn, found := k.GetAsset(ctx, pair.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	vault.AmountIn = vault.AmountIn.Add(msg.Amount)
	if !vault.AmountIn.IsPositive() {
		return nil, types.ErrorInvalidAmount
	}

	if err := k.SendCoinFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoin(assetIn.Denom, msg.Amount)); err != nil {
		return nil, err
	}

	k.SetVault(ctx, vault)
	return &types.MsgDepositResponse{}, nil
}

func (k *msgServer) MsgWithdraw(c context.Context, msg *types.MsgWithdrawRequest) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if k.GetOracleValidationResult(ctx) == false{
		return nil, types.ErrorOraclePriceExpired
	}

	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	vault, found := k.GetVault(ctx, msg.ID)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	if msg.From != vault.Owner {
		return nil, types.ErrorUnauthorized
	}

	pair, found := k.GetPair(ctx, vault.PairID)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}

	assetIn, found := k.GetAsset(ctx, pair.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	assetOut, found := k.GetAsset(ctx, pair.AssetOut)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	vault.AmountIn = vault.AmountIn.Sub(msg.Amount)
	if !vault.AmountIn.IsPositive() {
		return nil, types.ErrorInvalidAmount
	}

	if err := k.VerifyCollaterlizationRatio(ctx, vault.AmountIn, assetIn, vault.AmountOut, assetOut, pair.LiquidationRatio); err != nil {
		return nil, err
	}

	if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, from, sdk.NewCoin(assetIn.Denom, msg.Amount)); err != nil {
		return nil, err
	}

	k.SetVault(ctx, vault)
	return &types.MsgWithdrawResponse{}, nil
}

func (k *msgServer) MsgDraw(c context.Context, msg *types.MsgDrawRequest) (*types.MsgDrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if k.GetOracleValidationResult(ctx) == false{
		return nil, types.ErrorOraclePriceExpired
	}

	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	vault, found := k.GetVault(ctx, msg.ID)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	if msg.From != vault.Owner {
		return nil, types.ErrorUnauthorized
	}

	pair, found := k.GetPair(ctx, vault.PairID)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}

	assetIn, found := k.GetAsset(ctx, pair.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	assetOut, found := k.GetAsset(ctx, pair.AssetOut)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	vault.AmountOut = vault.AmountOut.Add(msg.Amount)
	if !vault.AmountOut.IsPositive() {
		return nil, types.ErrorInvalidAmount
	}

	if err := k.VerifyCollaterlizationRatio(ctx, vault.AmountIn, assetIn, vault.AmountOut, assetOut, pair.LiquidationRatio); err != nil {
		return nil, err
	}

	if err := k.MintCAssets(ctx, types.ModuleName, assetIn.Denom, assetOut.Denom, msg.Amount); err != nil {
		return nil, err
	}
	if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, from, sdk.NewCoin(assetOut.Denom, msg.Amount)); err != nil {
		return nil, err
	}

	k.SetVault(ctx, vault)
	return &types.MsgDrawResponse{}, nil
}

func (k *msgServer) MsgRepay(c context.Context, msg *types.MsgRepayRequest) (*types.MsgRepayResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if k.GetOracleValidationResult(ctx) == false{
		return nil, types.ErrorOraclePriceExpired
	}

	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	vault, found := k.GetVault(ctx, msg.ID)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}

	if msg.From != vault.Owner {
		return nil, types.ErrorUnauthorized
	}

	if !msg.Amount.LT(vault.AmountOut) {
		return nil, types.ErrorInvalidAmount
	}

	pair, found := k.GetPair(ctx, vault.PairID)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}

	assetIn, found := k.GetAsset(ctx, pair.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	assetOut, found := k.GetAsset(ctx, pair.AssetOut)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	if err := k.SendCoinFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoin(assetOut.Denom, msg.Amount)); err != nil {
		return nil, err
	}
	if err := k.BurnCAssets(ctx, types.ModuleName, assetIn.Denom, assetOut.Denom, msg.Amount); err != nil {
		return nil, err
	}

	vault.AmountOut = vault.AmountOut.Sub(msg.Amount)
	if !vault.AmountOut.IsPositive() {
		return nil, types.ErrorInvalidAmount
	}

	k.SetVault(ctx, vault)
	return &types.MsgRepayResponse{}, nil
}

func (k *msgServer) MsgClose(c context.Context, msg *types.MsgCloseRequest) (*types.MsgCloseResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	vault, found := k.GetVault(ctx, msg.ID)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	if msg.From != vault.Owner {
		return nil, types.ErrorUnauthorized
	}

	pair, found := k.GetPair(ctx, vault.PairID)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}

	assetIn, found := k.GetAsset(ctx, pair.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	assetOut, found := k.GetAsset(ctx, pair.AssetOut)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	if err := k.SendCoinFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoin(assetOut.Denom, vault.AmountOut)); err != nil {
		return nil, err
	}
	if err := k.BurnCAssets(ctx, types.ModuleName, assetIn.Denom, assetOut.Denom, vault.AmountOut); err != nil {
		return nil, err
	}
	if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, from, sdk.NewCoin(assetIn.Denom, vault.AmountIn)); err != nil {
		return nil, err
	}

	k.DeleteVault(ctx, vault.ID)
	k.DeleteVaultForAddressByPair(ctx, from, vault.PairID)
	k.UpdateUserVaultIdMapping(ctx, msg.From, vault.ID, false)
	k.UpdateCollateralVaultIdMapping(ctx, assetIn.Denom, assetOut.Denom, vault.ID, false)

	return &types.MsgCloseResponse{}, nil
}

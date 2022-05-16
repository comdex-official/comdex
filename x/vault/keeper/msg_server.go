package keeper

import (
	"context"
	"strconv"

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

	//checks if extended pair exists
	pairs, found := k.GetPairsVault(ctx, msg.ExtendedPairVaultID)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}

	//getting appMappingId from ExtendedPairVaultId
	appMappingId := pairs.AppMappingId

	//checking if appMappingId for appMappingId in ExtendedPairVault
	if appMappingId != msg.AppMappingId {
		return nil, types.ErrorAppIstoExtendedAppId
	}

	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	//check for duplicate vault
	if k.HasVaultForAddressByPair(ctx, from, msg.ExtendedPairVaultID) {
		return nil, types.ErrorDuplicateVault
	}

	// check for isPsmPair
	if pairs.IsPsmPair {
		return nil, types.ErrorCannotCreateStableSwapVault
	}

	if !(msg.AmountOut.GTE(pairs.DebtFloor)) {
		return nil, types.ErrorAmtGreaterDebt

	}

	Minted, foundMinted := k.GetTokenMintedID(ctx, msg.ExtendedPairVaultID)
	if !foundMinted {
		Minted = sdk.ZeroInt()
	}
	if !((pairs.DebtCieling).GT(msg.AmountOut.Add(Minted))) {
		return nil, types.ErrorAmtGreaterDebt
	}
	if pairs.AssetOutOraclePrice { //fetch oracle price from band

	} else { //don't consider oracle price, consider 1$

	}

	//get shortName for App
	app, _ := k.GetApp(ctx, appMappingId)
	sName := app.ShortName

	value, Notfound := k.GetCounterID(ctx, appMappingId)
	if !Notfound {
		count := 0
		k.SetCounterID(ctx, appMappingId, uint64(count))
	} else {
		k.SetCounterID(ctx, appMappingId, value)
	}

	pId := pairs.PairId
	assets, _ := k.GetPair(ctx, pId)

	assetDetailIN, _ := k.GetAsset(ctx, assets.AssetIn)
	assetDetailOUT, _ := k.GetAsset(ctx, assets.AssetOut)

	if err := k.VerifyCollaterlizationRatio(ctx, msg.AmountIn, assetDetailIN, msg.AmountOut, assetDetailOUT, pairs.LiquidationRatio); err != nil {
		return nil, err
	}

	if err := k.SendCoinFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoin(assetDetailIN.Denom, msg.AmountIn)); err != nil {
		return nil, err
	}
	if err := k.MintCoin(ctx, types.ModuleName, sdk.NewCoin(assetDetailOUT.Denom, msg.AmountOut)); err != nil {
		return nil, err
	}

	if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, from, sdk.NewCoin(assetDetailOUT.Denom, (msg.AmountOut-(sdk.Int(pairs.CreationFee)/100))*msg.AmountOut)); err != nil {
		return nil, err
	}

	var (
		// id, _  = k.GetCounterID(ctx, appMappingId)
		NewAppVaultTypeId = sName + strconv.Itoa(int(value+1))
		vault             = types.Vault{
			AppVaultTypeId:      NewAppVaultTypeId,
			ExtendedPairVaultID: msg.ExtendedPairVaultID,
			Owner:               msg.From,
			AmountIn:            msg.AmountIn,
			AmountOut:           msg.AmountOut,
		}
	)

	lookupVault := types.LookupTableVault{
		AppMappingId: appMappingId,
		Counter:      value,
	}
	lookupVault.AppVaultIds = append(lookupVault.AppVaultIds, NewAppVaultTypeId)

	k.SetLookupTableVault(ctx, lookupVault, appMappingId)
	k.SetVault(ctx, vault, sName)
	k.SetVaultForAddressByPair(ctx, from, vault.ExtendedPairVaultID, vault.Id)

	UserVaultId := types.UserVaultIdMapping{
		Owner: msg.From,
	}
	vaultToApp := types.VaultToAppMapping{
		AppMappingId: appMappingId,
	}
	vaultToApp.AppVaultTypeId = append(vaultToApp.AppVaultTypeId, NewAppVaultTypeId)
	UserVaultId.UserVaultIds = append(UserVaultId.UserVaultIds, &vaultToApp)
	k.SetUserVaultIdMapping(ctx, UserVaultId)

	PairtoVaultMap := types.PairToVaultMapping{
		PairId: pairs.PairId,
	}
	ExtendedPairvault := types.ExtendedPairVaultMap{
		ExtendedPairId: msg.ExtendedPairVaultID,
	}
	ExtendedPairvault.AppVaultTypeIds = append(ExtendedPairvault.AppVaultTypeIds, NewAppVaultTypeId)
	PairtoVaultMap.ExtendedpairVault = append(PairtoVaultMap.ExtendedpairVault, &ExtendedPairvault)
	k.SetPairVaultMapping(ctx, PairtoVaultMap)

	AppExtendedpPairVault := types.AppExtendedPairVaultMapping{
		AppId: appMappingId,
	}
	AppExtendedpPairVault.ExtendedpairVault = append(AppExtendedpPairVault.ExtendedpairVault, &ExtendedPairvault)
	k.SetExtendedPairMapping(ctx, AppExtendedpPairVault)

	valueMinted, foundMinted := k.GetTokenMintedID(ctx, msg.ExtendedPairVaultID)
	if !foundMinted {
		TokenMinted := sdk.ZeroInt()
		k.SetTokenMintedID(ctx, appMappingId, TokenMinted)
	} else {
		valueMinted = valueMinted.Add(msg.AmountOut)
		k.SetTokenMintedID(ctx, appMappingId, valueMinted)
	}
	ExtendedpPairVaultMap := types.ExtendedPairVaultMapping{
		ExtendedPairId: msg.ExtendedPairVaultID,
		TokenMinted:    value,
	}
	ExtendedpPairVaultMap.AppVaultTypeIds = append(ExtendedpPairVaultMap.AppVaultTypeIds, NewAppVaultTypeId)
	k.SetExtendedVaultPairMapping(ctx, ExtendedpPairVaultMap)

	// opening fee reduction from amout

	return &types.MsgCreateResponse{}, nil
}

func (k *msgServer) MsgDeposit(c context.Context, msg *types.MsgDepositRequest) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	//checks if extended pair exists
	pairs, found := k.GetPairsVault(ctx, msg.ExtendedPairVaultID)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}
	//getting appMappingId from ExtendedPairVaultId
	appMappingId := pairs.AppMappingId

	//checking if appMappingId for appMappingId in ExtendedPairVault
	if appMappingId != msg.AppMappingId {
		return nil, types.ErrorAppIstoExtendedAppId
	}
	//get shortName for App
	app, _ := k.GetApp(ctx, appMappingId)
	sName := app.ShortName

	vault, found := k.GetVault(ctx, msg.ID)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	if msg.From != vault.Owner {
		return nil, types.ErrorUnauthorized
	}

	pair, found := k.GetPair(ctx, vault.ExtendedPairVaultID)
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

	k.SetVault(ctx, vault, sName)
	return &types.MsgDepositResponse{}, nil
}

func (k *msgServer) MsgWithdraw(c context.Context, msg *types.MsgWithdrawRequest) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	//checks if extended pair exists
	pairs, found := k.GetPairsVault(ctx, msg.ExtendedPairVaultID)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}

	//getting appMappingId from ExtendedPairVaultId
	appMappingId := pairs.AppMappingId

	//checking if appMappingId for appMappingId in ExtendedPairVault
	if appMappingId != msg.AppMappingId {
		return nil, types.ErrorAppIstoExtendedAppId
	}
	//get shortName for App
	app, _ := k.GetApp(ctx, appMappingId)
	sName := app.ShortName

	vault, found := k.GetVault(ctx, msg.ID)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	if msg.From != vault.Owner {
		return nil, types.ErrorUnauthorized
	}

	pair, found := k.GetPair(ctx, vault.ExtendedPairVaultID)
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

	if err := k.VerifyCollaterlizationRatio(ctx, vault.AmountIn, assetIn, vault.AmountOut, assetOut, pairs.LiquidationRatio); err != nil {
		return nil, err
	}

	if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, from, sdk.NewCoin(assetIn.Denom, msg.Amount)); err != nil {
		return nil, err
	}

	k.SetVault(ctx, vault, sName)
	return &types.MsgWithdrawResponse{}, nil
}

func (k *msgServer) MsgDraw(c context.Context, msg *types.MsgDrawRequest) (*types.MsgDrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}
	//checks if extended pair exists
	pairs, found := k.GetPairsVault(ctx, msg.ExtendedPairVaultID)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}

	//getting appMappingId from ExtendedPairVaultId
	appMappingId := pairs.AppMappingId

	//checking if appMappingId for appMappingId in ExtendedPairVault
	if appMappingId != msg.AppMappingId {
		return nil, types.ErrorAppIstoExtendedAppId
	}
	//get shortName for App
	app, _ := k.GetApp(ctx, appMappingId)
	sName := app.ShortName

	vault, found := k.GetVault(ctx, msg.ID)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	if msg.From != vault.Owner {
		return nil, types.ErrorUnauthorized
	}

	pair, found := k.GetPair(ctx, vault.ExtendedPairVaultID)
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

	if err := k.VerifyCollaterlizationRatio(ctx, vault.AmountIn, assetIn, vault.AmountOut, assetOut, pairs.LiquidationRatio); err != nil {
		return nil, err
	}

	if err := k.MintCoin(ctx, types.ModuleName, sdk.NewCoin(assetOut.Denom, msg.Amount)); err != nil {
		return nil, err
	}
	if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, from, sdk.NewCoin(assetOut.Denom, msg.Amount)); err != nil {
		return nil, err
	}

	k.SetVault(ctx, vault, sName)
	return &types.MsgDrawResponse{}, nil
}

func (k *msgServer) MsgRepay(c context.Context, msg *types.MsgRepayRequest) (*types.MsgRepayResponse, error) {
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
	if !msg.Amount.Equal(vault.AmountOut) {
		return nil, types.ErrorInvalidAmount
	}

	pair, found := k.GetPair(ctx, vault.ExtendedPairVaultID)
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
	if err := k.BurnCoin(ctx, types.ModuleName, sdk.NewCoin(assetOut.Denom, vault.AmountOut)); err != nil {
		return nil, err
	}
	if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, from, sdk.NewCoin(assetIn.Denom, vault.AmountIn)); err != nil {
		return nil, err
	}

	k.DeleteVault(ctx, vault.AppVaultTypeId)
	k.DeleteVaultForAddressByPair(ctx, from, vault.ExtendedPairVaultID)

	return &types.MsgRepayResponse{}, nil
}

func (k *msgServer) MsgClose(c context.Context, msg *types.MsgCloseRequest) (*types.MsgCloseResponse, error) {
	panic("implement me")
}

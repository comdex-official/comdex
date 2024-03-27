package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/tokenmint/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

func NewMsgServer(keeper Keeper) types.MsgServer {
	return &msgServer{
		Keeper: keeper,
	}
}

func (k msgServer) MsgMintNewTokens(c context.Context, msg *types.MsgMintNewTokensRequest) (*types.MsgMintNewTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	assetData, found := k.asset.GetAsset(ctx, msg.AssetId)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	appMappingData, found := k.asset.GetApp(ctx, msg.AppId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExists
	}
	// Checking if asset exists in the app

	assetDataInApp, found := k.asset.GetMintGenesisTokenData(ctx, appMappingData.Id, assetData.Id)
	if !found {
		return nil, types.ErrorAssetNotWhiteListedForGenesisMinting
	}

	mintData, found := k.GetTokenMint(ctx, msg.AppId)
	if !found {
		var newTokenMintAppData types.MintedTokens
		var appData types.TokenMint

		if err := k.bank.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetData.Denom, assetDataInApp.GenesisSupply))); err != nil {
			return nil, err
		}
		userAddress, err := sdk.AccAddressFromBech32(assetDataInApp.Recipient)
		if err != nil {
			return nil, err
		}
		if assetDataInApp.GenesisSupply.GT(sdk.ZeroInt()) {
			if err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, userAddress, sdk.NewCoins(sdk.NewCoin(assetData.Denom, assetDataInApp.GenesisSupply))); err != nil {
				return nil, err
			}
		}

		newTokenMintAppData.AssetId = msg.AssetId
		newTokenMintAppData.CreatedAt = ctx.BlockTime()
		newTokenMintAppData.GenesisSupply = assetDataInApp.GenesisSupply
		newTokenMintAppData.CurrentSupply = newTokenMintAppData.GenesisSupply

		appData.AppId = appMappingData.Id
		appData.MintedTokens = append(appData.MintedTokens, &newTokenMintAppData)

		k.SetTokenMint(ctx, appData)
	} else {
		// AppData in TokenMint exists
		_, found := k.GetAssetDataInTokenMintByApp(ctx, appMappingData.Id, assetData.Id)
		if found {
			return nil, types.ErrorGenesisMintingForTokenAlreadyDone
		}
		userAddress, err := sdk.AccAddressFromBech32(assetDataInApp.Recipient)
		if err != nil {
			return nil, err
		}
		if assetDataInApp.GenesisSupply.GT(sdk.ZeroInt()) {
			if err = k.bank.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetData.Denom, assetDataInApp.GenesisSupply))); err != nil {
				return nil, err
			}
			if err = k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, userAddress, sdk.NewCoins(sdk.NewCoin(assetData.Denom, assetDataInApp.GenesisSupply))); err != nil {
				return nil, err
			}
		}

		var newTokenMintAppData types.MintedTokens
		newTokenMintAppData.AssetId = msg.AssetId
		newTokenMintAppData.CreatedAt = ctx.BlockTime()
		newTokenMintAppData.GenesisSupply = assetDataInApp.GenesisSupply
		newTokenMintAppData.CurrentSupply = newTokenMintAppData.GenesisSupply
		mintData.MintedTokens = append(mintData.MintedTokens, &newTokenMintAppData)
		k.SetTokenMint(ctx, mintData)
	}
	ctx.GasMeter().ConsumeGas(types.TokenmintGas, "TokenmintGas")

	return &types.MsgMintNewTokensResponse{}, nil
}

func (k msgServer) MsgBurnHarborTokens(c context.Context, msg *types.MsgBurnHarborTokensRequest) (*types.MsgBurnHarborTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	burner, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}
	if msg.BurnCoins.Denom != "uharbor" {
		return &types.MsgBurnHarborTokensResponse{}, types.ErrorInvalidDenom
	}
	appData, found := k.asset.GetApp(ctx, msg.AppId)
	if !found || appData.Name != "harbor" {
		return &types.MsgBurnHarborTokensResponse{}, types.ErrorAppMappingDoesNotExists
	}
	if appData.GenesisToken[0].Recipient != msg.From {
		return &types.MsgBurnHarborTokensResponse{}, types.ErrorInvalidFrom
	} 
	err = k.BurnGovTokensForApp(ctx, msg.AppId, burner, msg.BurnCoins)
	if err != nil {
		return &types.MsgBurnHarborTokensResponse{}, err
	}

	return &types.MsgBurnHarborTokensResponse{}, nil
}

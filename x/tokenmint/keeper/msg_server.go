package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/tokenmint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ types.MsgServer = msgServer{}
)

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
	assetData, found := k.GetAsset(ctx, msg.AssetId)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	appMappingData, found := k.GetApp(ctx, msg.AppId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExists
	}
	//Checking if asset exists in the app

	assetDataInApp, found := k.GetMintGenesisTokenData(ctx, appMappingData.Id, assetData.Id)
	if !found {
		return nil, types.ErrorAssetNotWhiteListedForGenesisMinting
	}

	mintData, found := k.GetTokenMint(ctx, msg.AppId)
	if !found {
		var newTokenMintAppData types.MintedTokens
		var appData types.TokenMint

		if err := k.MintCoin(ctx, types.ModuleName, sdk.NewCoin(assetData.Denom, assetDataInApp.GenesisSupply)); err != nil {
			return nil, err
		}
		userAddress, err := sdk.AccAddressFromBech32(assetDataInApp.Recipient)

		if err != nil {
			return nil, err
		}
		if assetDataInApp.GenesisSupply.GT(sdk.ZeroInt()) {
			if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, userAddress, sdk.NewCoin(assetData.Denom, assetDataInApp.GenesisSupply)); err != nil {
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
		//AppData in TokenMint exists
		_, found := k.GetAssetDataInTokenMintByApp(ctx, appMappingData.Id, assetData.Id)
		if found {
			return nil, types.ErrorGenesisMintingForTokenAlreadyDone
		}
		userAddress, err := sdk.AccAddressFromBech32(assetDataInApp.Recipient)
		if err != nil {
			return nil, err
		}
		if assetDataInApp.GenesisSupply.GT(sdk.ZeroInt()) {
			if err = k.MintCoin(ctx, types.ModuleName, sdk.NewCoin(assetData.Denom, assetDataInApp.GenesisSupply)); err != nil {
				return nil, err
			}
			if err = k.SendCoinFromModuleToAccount(ctx, types.ModuleName, userAddress, sdk.NewCoin(assetData.Denom, assetDataInApp.GenesisSupply)); err != nil {
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

package wasm

import (
	"encoding/json"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/comdex-official/comdex/app/wasm/bindings"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func CustomQuerier(queryPlugin *QueryPlugin) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var comdexQuery bindings.ComdexQuery
		if err := json.Unmarshal(request, &comdexQuery); err != nil {
			return nil, sdkerrors.Wrap(err, "app query")
		}
		if comdexQuery.AppData != nil {
			App_Id := comdexQuery.AppData.App_Id
			MinGovDeposit, GovTimeInSeconds, AssetId, _ := queryPlugin.GetAppInfo(ctx, App_Id)
			res := bindings.AppDataResponse{
				MinGovDeposit:    MinGovDeposit,
				GovTimeInSeconds: GovTimeInSeconds,
				AssetId:          AssetId,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "App data query response")
			}
			return bz, nil
		} else if comdexQuery.AssetData != nil {
			asset_Id := comdexQuery.AssetData.Asset_Id
			denom, _ := queryPlugin.GetAssetInfo(ctx, asset_Id)
			res := bindings.AssetDataResponse{
				Denom: denom,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "App data query response")
			}
			return bz, nil
		} else if comdexQuery.MintedToken != nil {
			App_Id := comdexQuery.MintedToken.App_Id
			Asset_Id := comdexQuery.MintedToken.Asset_Id
			MintedToken, _ := queryPlugin.GetTokenMint(ctx, App_Id, Asset_Id)
			res := bindings.MintedTokenResponse{
				MintedTokens: MintedToken,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "tokenMint query response")
			}
			return bz, nil
		} else if comdexQuery.State != nil {
			address := comdexQuery.State.Address
			denom := comdexQuery.State.Denom
			height := comdexQuery.State.Height
			target := comdexQuery.State.Target
			state, _ := GetState(address, denom, height, target)
			res := bindings.StateResponse{
				Amount: state,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "locker state query response")
			}
			return bz, nil
		}

		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown App Data query variant"}
	}

}

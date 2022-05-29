package asset

import (
	"encoding/json"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func CustomQuerier(assetKeeper *QueryPlugin) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var contractQuery AppQuery
		if err := json.Unmarshal(request, &contractQuery); err != nil {
			return nil, sdkerrors.Wrap(err, "app query")
		}
		if contractQuery.AppData != nil {
			App_Id := contractQuery.AppData.App_Id
			MinGovDeposit, GovTimeInSeconds, AssetId, _ := assetKeeper.GetAppInfo(ctx, App_Id)
			res := AppDataResponse{
				MinGovDeposit:    MinGovDeposit,
				GovTimeInSeconds: GovTimeInSeconds,
				AssetId:          AssetId,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "App data query response")
			}
			return bz, nil
		} else if contractQuery.AssetData != nil {
			asset_Id := contractQuery.AssetData.Asset_Id
			denom, _ := assetKeeper.GetAssetInfo(ctx, asset_Id)
			res := AssetDataResponse{
				Denom: denom,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "App data query response")
			}
			return bz, nil
		}

		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown App Data query variant"}
	}

}

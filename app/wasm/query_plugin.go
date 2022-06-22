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
			appID := comdexQuery.AppData.App_Id
			MinGovDeposit, GovTimeInSeconds, assetID, _ := queryPlugin.GetAppInfo(ctx, appID)
			res := bindings.AppDataResponse{
				MinGovDeposit:    MinGovDeposit,
				GovTimeInSeconds: GovTimeInSeconds,
				AssetId:          assetID,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "App data query response")
			}
			return bz, nil
		} else if comdexQuery.AssetData != nil {
			assetID := comdexQuery.AssetData.Asset_Id
			denom, _ := queryPlugin.GetAssetInfo(ctx, assetID)
			res := bindings.AssetDataResponse{
				Denom: denom,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "App data query response")
			}
			return bz, nil
		} else if comdexQuery.MintedToken != nil {
			appID := comdexQuery.MintedToken.App_Id
			assetID := comdexQuery.MintedToken.Asset_Id
			MintedToken, _ := queryPlugin.GetTokenMint(ctx, appID, assetID)
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
		} else if comdexQuery.RemoveWhiteListAssetLocker != nil {
			appID := comdexQuery.RemoveWhiteListAssetLocker.App_Id
			assetID := comdexQuery.RemoveWhiteListAssetLocker.Asset_Id

			found, errormsg := queryPlugin.GetRemoveWhitelistAppIdLockerRewardsCheck(ctx, appID, assetID)
			res := bindings.RemoveWhiteListAssetResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "RemoveWhiteListAssetLocker query response")
			}
			return bz, nil
		} else if comdexQuery.WhitelistAppIdLockerRewards != nil {
			appID := comdexQuery.WhitelistAppIdLockerRewards.App_Id
			assetID := comdexQuery.WhitelistAppIdLockerRewards.Asset_Id

			found, errormsg := queryPlugin.GetWhitelistAppIdLockerRewardsCheck(ctx, appID, assetID)
			res := bindings.WhitelistAppIdLockerRewardsResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "WhitelistAppIdLockerRewards query response")
			}
			return bz, nil
		} else if comdexQuery.WhitelistAppIdVaultInterest != nil {
			appID := comdexQuery.WhitelistAppIdVaultInterest.App_Id

			found, errormsg := queryPlugin.GetWhitelistAppIdVaultInterestCheck(ctx, appID)
			res := bindings.WhitelistAppIdLockerRewardsResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "WhitelistAppIdVaultInterest query response")
			}
			return bz, nil
		} else if comdexQuery.ExternalLockerRewards != nil {
			appID := comdexQuery.ExternalLockerRewards.App_Id
			assetID := comdexQuery.ExternalLockerRewards.Asset_Id

			found, errormsg := queryPlugin.GetExternalLockerRewardsCheck(ctx, appID, assetID)
			res := bindings.WhitelistAppIdLockerRewardsResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "GetExternalLockerRewardsCheck query response")
			}
			return bz, nil
		} else if comdexQuery.ExternalVaultRewards != nil {
			appID := comdexQuery.ExternalVaultRewards.App_Id
			assetID := comdexQuery.ExternalVaultRewards.Asset_Id

			found, errormsg := queryPlugin.GetExternalVaultRewardsCheck(ctx, appID, assetID)
			res := bindings.ExternalVaultRewardsResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "ExternalVaultRewards query response")
			}
			return bz, nil
		} else if comdexQuery.CollectorLookupTableQuery != nil {
			appMappingID := comdexQuery.CollectorLookupTableQuery.AppMappingId
			collectorAssetID := comdexQuery.CollectorLookupTableQuery.CollectorAssetId
			secondaryAssetID := comdexQuery.CollectorLookupTableQuery.SecondaryAssetId
			found, errormsg := queryPlugin.CollectorLookupTableQueryCheck(ctx, appMappingID, collectorAssetID, secondaryAssetID)
			res := bindings.CollectorLookupTableQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "ExternalVaultRewards query response")
			}
			return bz, nil
		} else if comdexQuery.ExtendedPairsVaultRecordsQuery != nil {
			appMappingID := comdexQuery.ExtendedPairsVaultRecordsQuery.AppMappingId
			pairID := comdexQuery.ExtendedPairsVaultRecordsQuery.PairId
			StabilityFee := comdexQuery.ExtendedPairsVaultRecordsQuery.StabilityFee
			ClosingFee := comdexQuery.ExtendedPairsVaultRecordsQuery.ClosingFee
			DrawDownFee := comdexQuery.ExtendedPairsVaultRecordsQuery.DrawDownFee
			DebtCeiling := comdexQuery.ExtendedPairsVaultRecordsQuery.DebtCeiling
			DebtFloor := comdexQuery.ExtendedPairsVaultRecordsQuery.DebtFloor
			PairName := comdexQuery.ExtendedPairsVaultRecordsQuery.PairName

			found, errorMsg := queryPlugin.ExtendedPairsVaultRecordsQueryCheck(ctx, appMappingID, pairID, StabilityFee, ClosingFee, DrawDownFee, DebtCeiling, DebtFloor, PairName)
			res := bindings.ExtendedPairsVaultRecordsQueryResponse{
				Found: found,
				Err:   errorMsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "ExternalVaultRewards query response")
			}
			return bz, nil
		} else if comdexQuery.AuctionMappingForAppQuery != nil {
			appMappingID := comdexQuery.AuctionMappingForAppQuery.AppMappingId
			found, errormsg := queryPlugin.AuctionMappingForAppQueryCheck(ctx, appMappingID)
			res := bindings.AuctionMappingForAppQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "AuctionMappingForAppQuery query response")
			}
			return bz, nil
		} else if comdexQuery.WhiteListedAssetQuery != nil {
			appMappingID := comdexQuery.WhiteListedAssetQuery.AppMappingId
			assetID := comdexQuery.WhiteListedAssetQuery.AssetId
			found, errormsg := queryPlugin.WhiteListedAssetQueryCheck(ctx, appMappingID, assetID)
			res := bindings.WhiteListedAssetQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "WhiteListedAssetQueryCheck query response")
			}
			return bz, nil
		} else if comdexQuery.UpdatePairsVaultQuery != nil {
			appMappingID := comdexQuery.UpdatePairsVaultQuery.AppMappingId
			extPairID := comdexQuery.UpdatePairsVaultQuery.ExtPairId
			found, errormsg := queryPlugin.UpdatePairsVaultQueryCheck(ctx, appMappingID, extPairID)
			res := bindings.UpdatePairsVaultQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "UpdatePairsVaultQuery query response")
			}
			return bz, nil
		} else if comdexQuery.UpdateCollectorLookupTableQuery != nil {
			appMappingID := comdexQuery.UpdateCollectorLookupTableQuery.AppMappingId
			assetID := comdexQuery.UpdateCollectorLookupTableQuery.AssetId
			found, errormsg := queryPlugin.UpdateCollectorLookupTableQueryCheck(ctx, appMappingID, assetID)
			res := bindings.UpdateCollectorLookupTableQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "UpdatePairsVaultQuery query response")
			}
			return bz, nil
		} else if comdexQuery.RemoveWhitelistAppIdVaultInterestQuery != nil {
			appMappingID := comdexQuery.RemoveWhitelistAppIdVaultInterestQuery.AppMappingId
			found, errormsg := queryPlugin.WasmRemoveWhitelistAppIdVaultInterestQueryCheck(ctx, appMappingID)
			res := bindings.RemoveWhitelistAppIdVaultInterestQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "RemoveWhitelistAppIdVaultInterestQuery query response")
			}
			return bz, nil
		} else if comdexQuery.RemoveWhitelistAssetLockerQuery != nil {
			appMappingID := comdexQuery.RemoveWhitelistAssetLockerQuery.AppMappingId
			assetID := comdexQuery.RemoveWhitelistAssetLockerQuery.AssetId

			found, errormsg := queryPlugin.WasmRemoveWhitelistAssetLockerQueryCheck(ctx, appMappingID, assetID)
			res := bindings.RemoveWhitelistAssetLockerQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "RemoveWhitelistAssetLockerQuery query response")
			}
			return bz, nil
		} else if comdexQuery.WhitelistAppIdLiquidationQuery != nil {
			AppMappingID := comdexQuery.WhitelistAppIdLiquidationQuery.AppMappingId

			found, errormsg := queryPlugin.WasmWhitelistAppIdLiquidationQueryCheck(ctx, AppMappingID)
			res := bindings.WhitelistAppIdLiquidationQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "WhitelistAppIdLiquidationQuery query response")
			}
			return bz, nil
		} else if comdexQuery.RemoveWhitelistAppIdLiquidationQuery != nil {
			AppMappingID := comdexQuery.RemoveWhitelistAppIdLiquidationQuery.AppMappingId

			found, errormsg := queryPlugin.WasmRemoveWhitelistAppIdLiquidationQueryCheck(ctx, AppMappingID)
			res := bindings.RemoveWhitelistAppIdLiquidationQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "RemoveWhitelistAppIdLiquidationQuery query response")
			}
			return bz, nil
		}
		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown App Data query variant"}
	}
}

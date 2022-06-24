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
			appID := comdexQuery.AppData.AppMappingID
			MinGovDeposit, GovTimeInSeconds, assetID, _ := queryPlugin.GetAppInfo(ctx, appID)
			res := bindings.AppDataResponse{
				MinGovDeposit:    MinGovDeposit,
				GovTimeInSeconds: GovTimeInSeconds,
				AssetID:          assetID,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "App data query response")
			}
			return bz, nil
		} else if comdexQuery.AssetData != nil {
			assetID := comdexQuery.AssetData.AssetID
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
			appID := comdexQuery.MintedToken.AppMappingID
			assetID := comdexQuery.MintedToken.AssetID
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
			appID := comdexQuery.RemoveWhiteListAssetLocker.AppMappingID
			assetID := comdexQuery.RemoveWhiteListAssetLocker.AssetIDs

			found, errormsg := queryPlugin.GetRemoveWhitelistAppIDLockerRewardsCheck(ctx, appID, assetID)
			res := bindings.RemoveWhiteListAssetResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "RemoveWhiteListAssetLocker query response")
			}
			return bz, nil
		} else if comdexQuery.WhitelistAppIDLockerRewards != nil {
			appID := comdexQuery.WhitelistAppIDLockerRewards.AppMappingID
			assetID := comdexQuery.WhitelistAppIDLockerRewards.AssetID

			found, errormsg := queryPlugin.GetWhitelistAppIDLockerRewardsCheck(ctx, appID, assetID)
			res := bindings.WhitelistAppIDLockerRewardsResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "WhitelistAppIdLockerRewards query response")
			}
			return bz, nil
		} else if comdexQuery.WhitelistAppIDVaultInterest != nil {
			appID := comdexQuery.WhitelistAppIDVaultInterest.AppMappingID

			found, errormsg := queryPlugin.GetWhitelistAppIDVaultInterestCheck(ctx, appID)
			res := bindings.WhitelistAppIDLockerRewardsResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "WhitelistAppIdVaultInterest query response")
			}
			return bz, nil
		} else if comdexQuery.ExternalLockerRewards != nil {
			appID := comdexQuery.ExternalLockerRewards.AppMappingID
			assetID := comdexQuery.ExternalLockerRewards.AssetID

			found, errormsg := queryPlugin.GetExternalLockerRewardsCheck(ctx, appID, assetID)
			res := bindings.WhitelistAppIDLockerRewardsResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "GetExternalLockerRewardsCheck query response")
			}
			return bz, nil
		} else if comdexQuery.ExternalVaultRewards != nil {
			appID := comdexQuery.ExternalVaultRewards.AppMappingID
			assetID := comdexQuery.ExternalVaultRewards.AssetID

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
			appMappingID := comdexQuery.CollectorLookupTableQuery.AppMappingID
			collectorAssetID := comdexQuery.CollectorLookupTableQuery.CollectorAssetID
			secondaryAssetID := comdexQuery.CollectorLookupTableQuery.SecondaryAssetID
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
			appMappingID := comdexQuery.ExtendedPairsVaultRecordsQuery.AppMappingID
			pairID := comdexQuery.ExtendedPairsVaultRecordsQuery.PairID
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
			appMappingID := comdexQuery.AuctionMappingForAppQuery.AppMappingID
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
			appMappingID := comdexQuery.WhiteListedAssetQuery.AppMappingID
			assetID := comdexQuery.WhiteListedAssetQuery.AssetID
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
			appMappingID := comdexQuery.UpdatePairsVaultQuery.AppMappingID
			extPairID := comdexQuery.UpdatePairsVaultQuery.ExtPairID
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
			appMappingID := comdexQuery.UpdateCollectorLookupTableQuery.AppMappingID
			assetID := comdexQuery.UpdateCollectorLookupTableQuery.AssetID
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
		} else if comdexQuery.RemoveWhitelistAppIDVaultInterestQuery != nil {
			appMappingID := comdexQuery.RemoveWhitelistAppIDVaultInterestQuery.AppMappingID
			found, errormsg := queryPlugin.WasmRemoveWhitelistAppIDVaultInterestQueryCheck(ctx, appMappingID)
			res := bindings.RemoveWhitelistAppIDVaultInterestQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "RemoveWhitelistAppIdVaultInterestQuery query response")
			}
			return bz, nil
		} else if comdexQuery.RemoveWhitelistAssetLockerQuery != nil {
			appMappingID := comdexQuery.RemoveWhitelistAssetLockerQuery.AppMappingID
			assetID := comdexQuery.RemoveWhitelistAssetLockerQuery.AssetID

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
		} else if comdexQuery.WhitelistAppIDLiquidationQuery != nil {
			AppMappingID := comdexQuery.WhitelistAppIDLiquidationQuery.AppMappingID

			found, errormsg := queryPlugin.WasmWhitelistAppIDLiquidationQueryCheck(ctx, AppMappingID)
			res := bindings.WhitelistAppIDLiquidationQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "WhitelistAppIDLiquidationQuery query response")
			}
			return bz, nil
		} else if comdexQuery.RemoveWhitelistAppIDLiquidationQuery != nil {
			AppMappingID := comdexQuery.RemoveWhitelistAppIDLiquidationQuery.AppMappingID

			found, errormsg := queryPlugin.WasmRemoveWhitelistAppIDLiquidationQueryCheck(ctx, AppMappingID)
			res := bindings.RemoveWhitelistAppIDLiquidationQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "RemoveWhitelistAppIDLiquidationQuery query response")
			}
			return bz, nil
		}
		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown App Data query variant"}
	}
}

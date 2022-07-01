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
			appID := comdexQuery.AppData.AppID
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
			appID := comdexQuery.MintedToken.AppID
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
			appID := comdexQuery.RemoveWhiteListAssetLocker.AppID
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
			appID := comdexQuery.WhitelistAppIDLockerRewards.AppID
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
			appID := comdexQuery.WhitelistAppIDVaultInterest.AppID

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
			appID := comdexQuery.ExternalLockerRewards.AppID
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
			appID := comdexQuery.ExternalVaultRewards.AppID
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
			appID := comdexQuery.CollectorLookupTableQuery.AppID
			collectorAssetID := comdexQuery.CollectorLookupTableQuery.CollectorAssetID
			secondaryAssetID := comdexQuery.CollectorLookupTableQuery.SecondaryAssetID
			found, errormsg := queryPlugin.CollectorLookupTableQueryCheck(ctx, appID, collectorAssetID, secondaryAssetID)
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
			appID := comdexQuery.ExtendedPairsVaultRecordsQuery.AppID
			pairID := comdexQuery.ExtendedPairsVaultRecordsQuery.PairID
			StabilityFee := comdexQuery.ExtendedPairsVaultRecordsQuery.StabilityFee
			ClosingFee := comdexQuery.ExtendedPairsVaultRecordsQuery.ClosingFee
			DrawDownFee := comdexQuery.ExtendedPairsVaultRecordsQuery.DrawDownFee
			DebtCeiling := comdexQuery.ExtendedPairsVaultRecordsQuery.DebtCeiling
			DebtFloor := comdexQuery.ExtendedPairsVaultRecordsQuery.DebtFloor
			PairName := comdexQuery.ExtendedPairsVaultRecordsQuery.PairName

			found, errorMsg := queryPlugin.ExtendedPairsVaultRecordsQueryCheck(ctx, appID, pairID, StabilityFee, ClosingFee, DrawDownFee, DebtCeiling, DebtFloor, PairName)
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
			appID := comdexQuery.AuctionMappingForAppQuery.AppID
			found, errormsg := queryPlugin.AuctionMappingForAppQueryCheck(ctx, appID)
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
			appID := comdexQuery.WhiteListedAssetQuery.AppID
			assetID := comdexQuery.WhiteListedAssetQuery.AssetID
			found, errormsg := queryPlugin.WhiteListedAssetQueryCheck(ctx, appID, assetID)
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
			appID := comdexQuery.UpdatePairsVaultQuery.AppID
			extPairID := comdexQuery.UpdatePairsVaultQuery.ExtPairID
			found, errormsg := queryPlugin.UpdatePairsVaultQueryCheck(ctx, appID, extPairID)
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
			appID := comdexQuery.UpdateCollectorLookupTableQuery.AppID
			assetID := comdexQuery.UpdateCollectorLookupTableQuery.AssetID
			found, errormsg := queryPlugin.UpdateCollectorLookupTableQueryCheck(ctx, appID, assetID)
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
			appID := comdexQuery.RemoveWhitelistAppIDVaultInterestQuery.AppID
			found, errormsg := queryPlugin.WasmRemoveWhitelistAppIDVaultInterestQueryCheck(ctx, appID)
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
			appID := comdexQuery.RemoveWhitelistAssetLockerQuery.AppID
			assetID := comdexQuery.RemoveWhitelistAssetLockerQuery.AssetID

			found, errormsg := queryPlugin.WasmRemoveWhitelistAssetLockerQueryCheck(ctx, appID, assetID)
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
			AppID := comdexQuery.WhitelistAppIDLiquidationQuery.AppID

			found, errormsg := queryPlugin.WasmWhitelistAppIDLiquidationQueryCheck(ctx, AppID)
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
			AppID := comdexQuery.RemoveWhitelistAppIDLiquidationQuery.AppID

			found, errormsg := queryPlugin.WasmRemoveWhitelistAppIDLiquidationQueryCheck(ctx, AppID)
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

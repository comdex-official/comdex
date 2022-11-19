package wasm

import (
	"encoding/json"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/petrichormoney/petri/app/wasm/bindings"
)

func CustomQuerier(queryPlugin *QueryPlugin) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var petriQuery bindings.ComdexQuery
		if err := json.Unmarshal(request, &petriQuery); err != nil {
			return nil, sdkerrors.Wrap(err, "app query")
		}
		if petriQuery.AppData != nil {
			appID := petriQuery.AppData.AppID
			MinGovDeposit, GovTimeInSeconds, assetID, _ := queryPlugin.GetAppInfo(ctx, appID)
			res := bindings.AppDataResponse{
				MinGovDeposit:    MinGovDeposit.String(),
				GovTimeInSeconds: GovTimeInSeconds,
				AssetID:          assetID,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "App data query response")
			}
			return bz, nil
		} else if petriQuery.AssetData != nil {
			assetID := petriQuery.AssetData.AssetID
			denom, _ := queryPlugin.GetAssetInfo(ctx, assetID)
			res := bindings.AssetDataResponse{
				Denom: denom,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "App data query response")
			}
			return bz, nil
		} else if petriQuery.MintedToken != nil {
			appID := petriQuery.MintedToken.AppID
			assetID := petriQuery.MintedToken.AssetID
			MintedToken, _ := queryPlugin.GetTokenMint(ctx, appID, assetID)
			res := bindings.MintedTokenResponse{
				MintedTokens: MintedToken,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "tokenMint query response")
			}
			return bz, nil
		} else if petriQuery.RemoveWhiteListAssetLocker != nil {
			appID := petriQuery.RemoveWhiteListAssetLocker.AppID
			assetID := petriQuery.RemoveWhiteListAssetLocker.AssetIDs

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
		} else if petriQuery.WhitelistAppIDLockerRewards != nil {
			appID := petriQuery.WhitelistAppIDLockerRewards.AppID
			assetID := petriQuery.WhitelistAppIDLockerRewards.AssetID

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
		} else if petriQuery.WhitelistAppIDVaultInterest != nil {
			appID := petriQuery.WhitelistAppIDVaultInterest.AppID

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
		} else if petriQuery.ExternalLockerRewards != nil {
			appID := petriQuery.ExternalLockerRewards.AppID
			assetID := petriQuery.ExternalLockerRewards.AssetID

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
		} else if petriQuery.ExternalVaultRewards != nil {
			appID := petriQuery.ExternalVaultRewards.AppID
			assetID := petriQuery.ExternalVaultRewards.AssetID

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
		} else if petriQuery.CollectorLookupTableQuery != nil {
			appID := petriQuery.CollectorLookupTableQuery.AppID
			collectorAssetID := petriQuery.CollectorLookupTableQuery.CollectorAssetID
			secondaryAssetID := petriQuery.CollectorLookupTableQuery.SecondaryAssetID
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
		} else if petriQuery.ExtendedPairsVaultRecordsQuery != nil {
			appID := petriQuery.ExtendedPairsVaultRecordsQuery.AppID
			pairID := petriQuery.ExtendedPairsVaultRecordsQuery.PairID
			StabilityFee := petriQuery.ExtendedPairsVaultRecordsQuery.StabilityFee
			ClosingFee := petriQuery.ExtendedPairsVaultRecordsQuery.ClosingFee
			DrawDownFee := petriQuery.ExtendedPairsVaultRecordsQuery.DrawDownFee
			DebtCeiling := petriQuery.ExtendedPairsVaultRecordsQuery.DebtCeiling
			DebtFloor := petriQuery.ExtendedPairsVaultRecordsQuery.DebtFloor
			PairName := petriQuery.ExtendedPairsVaultRecordsQuery.PairName

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
		} else if petriQuery.AuctionMappingForAppQuery != nil {
			appID := petriQuery.AuctionMappingForAppQuery.AppID
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
		} else if petriQuery.WhiteListedAssetQuery != nil {
			appID := petriQuery.WhiteListedAssetQuery.AppID
			assetID := petriQuery.WhiteListedAssetQuery.AssetID
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
		} else if petriQuery.UpdatePairsVaultQuery != nil {
			appID := petriQuery.UpdatePairsVaultQuery.AppID
			extPairID := petriQuery.UpdatePairsVaultQuery.ExtPairID
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
		} else if petriQuery.UpdateCollectorLookupTableQuery != nil {
			appID := petriQuery.UpdateCollectorLookupTableQuery.AppID
			assetID := petriQuery.UpdateCollectorLookupTableQuery.AssetID
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
		} else if petriQuery.RemoveWhitelistAppIDVaultInterestQuery != nil {
			appID := petriQuery.RemoveWhitelistAppIDVaultInterestQuery.AppID
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
		} else if petriQuery.RemoveWhitelistAssetLockerQuery != nil {
			appID := petriQuery.RemoveWhitelistAssetLockerQuery.AppID
			assetID := petriQuery.RemoveWhitelistAssetLockerQuery.AssetID

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
		} else if petriQuery.WhitelistAppIDLiquidationQuery != nil {
			AppID := petriQuery.WhitelistAppIDLiquidationQuery.AppID

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
		} else if petriQuery.RemoveWhitelistAppIDLiquidationQuery != nil {
			AppID := petriQuery.RemoveWhitelistAppIDLiquidationQuery.AppID

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
		} else if petriQuery.AddESMTriggerParamsForAppQuery != nil {
			AppID := petriQuery.AddESMTriggerParamsForAppQuery.AppID

			found, errormsg := queryPlugin.WasmAddESMTriggerParamsQueryCheck(ctx, AppID)
			res := bindings.AddESMTriggerParamsForAppResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "AddESMTriggerParamsForAppResponse query response")
			}
			return bz, nil
		} else if petriQuery.ExtendedPairByApp != nil {
			AppID := petriQuery.ExtendedPairByApp.AppID

			extendedPair, _ := queryPlugin.WasmExtendedPairByApp(ctx, AppID)
			res := bindings.ExtendedPairByAppResponse{
				ExtendedPair: extendedPair,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "ExtendedPairByAppResponse query response")
			}
			return bz, nil
		} else if petriQuery.CheckSurplusReward != nil {
			AppID := petriQuery.CheckSurplusReward.AppID
			AssetID := petriQuery.CheckSurplusReward.AssetID
			amount := queryPlugin.WasmCheckSurplusReward(ctx, AppID, AssetID)
			res := bindings.CheckSurplusRewardResponse{
				Amount: amount,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "CheckSurplusRewardResponse query response")
			}
			return bz, nil
		} else if petriQuery.CheckWhitelistedAsset != nil {
			Denom := petriQuery.CheckWhitelistedAsset.Denom

			found := queryPlugin.WasmCheckWhitelistedAsset(ctx, Denom)
			res := bindings.CheckWhitelistedAssetResponse{
				Found: found,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "CheckWhitelistedAssetResponse query response")
			}
			return bz, nil
		} else if petriQuery.CheckVaultCreated != nil {
			Address := petriQuery.CheckVaultCreated.Address
			AppID := petriQuery.CheckVaultCreated.AppID
			found := queryPlugin.WasmCheckVaultCreated(ctx, Address, AppID)
			res := bindings.VaultCreatedResponse{
				IsCompleted: found,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "VaultCreatedResponse query response")
			}
			return bz, nil
		} else if petriQuery.CheckBorrowed != nil {
			AssetID := petriQuery.CheckBorrowed.AssetID
			Address := petriQuery.CheckBorrowed.Address
			found := queryPlugin.WasmCheckBorrowed(ctx, AssetID, Address)
			res := bindings.BorrowedResponse{
				IsCompleted: found,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "BorrowedResponse query response")
			}
			return bz, nil
		} else if petriQuery.CheckLiquidityProvided != nil {
			AppID := petriQuery.CheckLiquidityProvided.AppID
			PoolID := petriQuery.CheckLiquidityProvided.PoolID
			Address := petriQuery.CheckLiquidityProvided.Address
			found := queryPlugin.WasmCheckLiquidityProvided(ctx, AppID, PoolID, Address)
			res := bindings.LiquidityProvidedResponse{
				IsCompleted: found,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "LiquidityProvidedResponse query response")
			}
			return bz, nil
		}
		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown App Data query variant"}
	}
}

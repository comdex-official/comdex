package wasm

import (
	errorsmod "cosmossdk.io/errors"
	"encoding/json"
	"fmt"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/app/wasm/bindings"
)

func CustomQuerier(queryPlugin *QueryPlugin) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var contractQuery bindings.ContractQuery
		if err := json.Unmarshal(request, &contractQuery); err != nil {
			return nil, errorsmod.Wrap(err, "app query")
		}
		switch {
		case contractQuery.AppData != nil:
			appID := contractQuery.AppData.AppID
			MinGovDeposit, GovTimeInSeconds, assetID, _ := queryPlugin.GetAppInfo(ctx, appID)
			res := bindings.AppDataResponse{
				MinGovDeposit:    MinGovDeposit.String(),
				GovTimeInSeconds: GovTimeInSeconds,
				AssetID:          assetID,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "App data query response")
			}
			return bz, nil
		case contractQuery.AssetData != nil:
			assetID := contractQuery.AssetData.AssetID
			denom, _ := queryPlugin.GetAssetInfo(ctx, assetID)
			res := bindings.AssetDataResponse{
				Denom: denom,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "App data query response")
			}
			return bz, nil
		case contractQuery.MintedToken != nil:
			appID := contractQuery.MintedToken.AppID
			assetID := contractQuery.MintedToken.AssetID
			MintedToken, _ := queryPlugin.GetTokenMint(ctx, appID, assetID)
			res := bindings.MintedTokenResponse{
				MintedTokens: MintedToken,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "tokenMint query response")
			}
			return bz, nil
		case contractQuery.RemoveWhiteListAssetLocker != nil:
			appID := contractQuery.RemoveWhiteListAssetLocker.AppID
			assetID := contractQuery.RemoveWhiteListAssetLocker.AssetIDs

			found, errormsg := queryPlugin.GetRemoveWhitelistAppIDLockerRewardsCheck(ctx, appID, assetID)
			res := bindings.RemoveWhiteListAssetResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "RemoveWhiteListAssetLocker query response")
			}
			return bz, nil
		case contractQuery.WhitelistAppIDLockerRewards != nil:
			appID := contractQuery.WhitelistAppIDLockerRewards.AppID
			assetID := contractQuery.WhitelistAppIDLockerRewards.AssetID

			found, errormsg := queryPlugin.GetWhitelistAppIDLockerRewardsCheck(ctx, appID, assetID)
			res := bindings.WhitelistAppIDLockerRewardsResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "WhitelistAppIdLockerRewards query response")
			}
			return bz, nil
		case contractQuery.WhitelistAppIDVaultInterest != nil:
			appID := contractQuery.WhitelistAppIDVaultInterest.AppID

			found, errormsg := queryPlugin.GetWhitelistAppIDVaultInterestCheck(ctx, appID)
			res := bindings.WhitelistAppIDLockerRewardsResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "WhitelistAppIdVaultInterest query response")
			}
			return bz, nil
		case contractQuery.ExternalLockerRewards != nil:
			appID := contractQuery.ExternalLockerRewards.AppID
			assetID := contractQuery.ExternalLockerRewards.AssetID

			found, errormsg := queryPlugin.GetExternalLockerRewardsCheck(ctx, appID, assetID)
			res := bindings.WhitelistAppIDLockerRewardsResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "GetExternalLockerRewardsCheck query response")
			}
			return bz, nil
		case contractQuery.ExternalVaultRewards != nil:
			appID := contractQuery.ExternalVaultRewards.AppID
			assetID := contractQuery.ExternalVaultRewards.AssetID

			found, errormsg := queryPlugin.GetExternalVaultRewardsCheck(ctx, appID, assetID)
			res := bindings.ExternalVaultRewardsResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "ExternalVaultRewards query response")
			}
			return bz, nil
		case contractQuery.CollectorLookupTableQuery != nil:
			appID := contractQuery.CollectorLookupTableQuery.AppID
			collectorAssetID := contractQuery.CollectorLookupTableQuery.CollectorAssetID
			secondaryAssetID := contractQuery.CollectorLookupTableQuery.SecondaryAssetID
			found, errormsg := queryPlugin.CollectorLookupTableQueryCheck(ctx, appID, collectorAssetID, secondaryAssetID)
			res := bindings.CollectorLookupTableQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "ExternalVaultRewards query response")
			}
			return bz, nil
		case contractQuery.ExtendedPairsVaultRecordsQuery != nil:
			appID := contractQuery.ExtendedPairsVaultRecordsQuery.AppID
			pairID := contractQuery.ExtendedPairsVaultRecordsQuery.PairID
			StabilityFee := contractQuery.ExtendedPairsVaultRecordsQuery.StabilityFee
			ClosingFee := contractQuery.ExtendedPairsVaultRecordsQuery.ClosingFee
			DrawDownFee := contractQuery.ExtendedPairsVaultRecordsQuery.DrawDownFee
			DebtCeiling := contractQuery.ExtendedPairsVaultRecordsQuery.DebtCeiling
			DebtFloor := contractQuery.ExtendedPairsVaultRecordsQuery.DebtFloor
			PairName := contractQuery.ExtendedPairsVaultRecordsQuery.PairName

			found, errorMsg := queryPlugin.ExtendedPairsVaultRecordsQueryCheck(ctx, appID, pairID, StabilityFee, ClosingFee, DrawDownFee, DebtCeiling, DebtFloor, PairName)
			res := bindings.ExtendedPairsVaultRecordsQueryResponse{
				Found: found,
				Err:   errorMsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "ExternalVaultRewards query response")
			}
			return bz, nil
		case contractQuery.AuctionMappingForAppQuery != nil:
			appID := contractQuery.AuctionMappingForAppQuery.AppID
			found, errormsg := queryPlugin.AuctionMappingForAppQueryCheck(ctx, appID)
			res := bindings.AuctionMappingForAppQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "AuctionMappingForAppQuery query response")
			}
			return bz, nil
		case contractQuery.WhiteListedAssetQuery != nil:
			appID := contractQuery.WhiteListedAssetQuery.AppID
			assetID := contractQuery.WhiteListedAssetQuery.AssetID
			found, errormsg := queryPlugin.WhiteListedAssetQueryCheck(ctx, appID, assetID)
			res := bindings.WhiteListedAssetQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "WhiteListedAssetQueryCheck query response")
			}
			return bz, nil
		case contractQuery.UpdatePairsVaultQuery != nil:
			appID := contractQuery.UpdatePairsVaultQuery.AppID
			extPairID := contractQuery.UpdatePairsVaultQuery.ExtPairID
			found, errormsg := queryPlugin.UpdatePairsVaultQueryCheck(ctx, appID, extPairID)
			res := bindings.UpdatePairsVaultQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "UpdatePairsVaultQuery query response")
			}
			return bz, nil
		case contractQuery.UpdateCollectorLookupTableQuery != nil:
			appID := contractQuery.UpdateCollectorLookupTableQuery.AppID
			assetID := contractQuery.UpdateCollectorLookupTableQuery.AssetID
			found, errormsg := queryPlugin.UpdateCollectorLookupTableQueryCheck(ctx, appID, assetID)
			res := bindings.UpdateCollectorLookupTableQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "UpdatePairsVaultQuery query response")
			}
			return bz, nil
		case contractQuery.RemoveWhitelistAppIDVaultInterestQuery != nil:
			appID := contractQuery.RemoveWhitelistAppIDVaultInterestQuery.AppID
			found, errormsg := queryPlugin.WasmRemoveWhitelistAppIDVaultInterestQueryCheck(ctx, appID)
			res := bindings.RemoveWhitelistAppIDVaultInterestQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "RemoveWhitelistAppIdVaultInterestQuery query response")
			}
			return bz, nil
		case contractQuery.RemoveWhitelistAssetLockerQuery != nil:
			appID := contractQuery.RemoveWhitelistAssetLockerQuery.AppID
			assetID := contractQuery.RemoveWhitelistAssetLockerQuery.AssetID

			found, errormsg := queryPlugin.WasmRemoveWhitelistAssetLockerQueryCheck(ctx, appID, assetID)
			res := bindings.RemoveWhitelistAssetLockerQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "RemoveWhitelistAssetLockerQuery query response")
			}
			return bz, nil
		case contractQuery.WhitelistAppIDLiquidationQuery != nil:
			AppID := contractQuery.WhitelistAppIDLiquidationQuery.AppID

			found, errormsg := queryPlugin.WasmWhitelistAppIDLiquidationQueryCheck(ctx, AppID)
			res := bindings.WhitelistAppIDLiquidationQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "WhitelistAppIDLiquidationQuery query response")
			}
			return bz, nil
		case contractQuery.RemoveWhitelistAppIDLiquidationQuery != nil:
			AppID := contractQuery.RemoveWhitelistAppIDLiquidationQuery.AppID

			found, errormsg := queryPlugin.WasmRemoveWhitelistAppIDLiquidationQueryCheck(ctx, AppID)
			res := bindings.RemoveWhitelistAppIDLiquidationQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "RemoveWhitelistAppIDLiquidationQuery query response")
			}
			return bz, nil
		case contractQuery.AddESMTriggerParamsForAppQuery != nil:
			AppID := contractQuery.AddESMTriggerParamsForAppQuery.AppID

			found, errormsg := queryPlugin.WasmAddESMTriggerParamsQueryCheck(ctx, AppID)
			res := bindings.AddESMTriggerParamsForAppResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "AddESMTriggerParamsForAppResponse query response")
			}
			return bz, nil
		case contractQuery.ExtendedPairByApp != nil:
			AppID := contractQuery.ExtendedPairByApp.AppID

			extendedPair, _ := queryPlugin.WasmExtendedPairByApp(ctx, AppID)
			res := bindings.ExtendedPairByAppResponse{
				ExtendedPair: extendedPair,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "ExtendedPairByAppResponse query response")
			}
			return bz, nil
		case contractQuery.CheckSurplusReward != nil:
			AppID := contractQuery.CheckSurplusReward.AppID
			AssetID := contractQuery.CheckSurplusReward.AssetID
			amount := queryPlugin.WasmCheckSurplusReward(ctx, AppID, AssetID)
			res := bindings.CheckSurplusRewardResponse{
				Amount: amount,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "CheckSurplusRewardResponse query response")
			}
			return bz, nil
		case contractQuery.CheckWhitelistedAsset != nil:
			Denom := contractQuery.CheckWhitelistedAsset.Denom

			found := queryPlugin.WasmCheckWhitelistedAsset(ctx, Denom)
			res := bindings.CheckWhitelistedAssetResponse{
				Found: found,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "CheckWhitelistedAssetResponse query response")
			}
			return bz, nil
		case contractQuery.CheckVaultCreated != nil:
			Address := contractQuery.CheckVaultCreated.Address
			AppID := contractQuery.CheckVaultCreated.AppID
			found := queryPlugin.WasmCheckVaultCreated(ctx, Address, AppID)
			res := bindings.VaultCreatedResponse{
				IsCompleted: found,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "VaultCreatedResponse query response")
			}
			return bz, nil
		case contractQuery.CheckBorrowed != nil:
			AssetID := contractQuery.CheckBorrowed.AssetID
			Address := contractQuery.CheckBorrowed.Address
			found := queryPlugin.WasmCheckBorrowed(ctx, AssetID, Address)
			res := bindings.BorrowedResponse{
				IsCompleted: found,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "BorrowedResponse query response")
			}
			return bz, nil
		case contractQuery.CheckLiquidityProvided != nil:
			AppID := contractQuery.CheckLiquidityProvided.AppID
			PoolID := contractQuery.CheckLiquidityProvided.PoolID
			Address := contractQuery.CheckLiquidityProvided.Address
			found := queryPlugin.WasmCheckLiquidityProvided(ctx, AppID, PoolID, Address)
			res := bindings.LiquidityProvidedResponse{
				IsCompleted: found,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "LiquidityProvidedResponse query response")
			}
			return bz, nil
		case contractQuery.GetPoolByApp != nil:
			AppID := contractQuery.GetPoolByApp.AppID
			pools := queryPlugin.WasmGetPools(ctx, AppID)
			res := bindings.GetPoolByAppResponse{
				Pools: pools,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "GetPoolByApp query response")
			}
			return bz, nil
		case contractQuery.GetAssetPrice != nil:
			assetID := contractQuery.GetAssetPrice.AssetID
			assetPrice, _ := queryPlugin.WasmGetAssetPrice(ctx, assetID)
			res := bindings.GetAssetPriceResponse{
				Price: assetPrice,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "GetAssetPrice query response")
			}
			return bz, nil
		case contractQuery.FullDenom != nil:
			creator := contractQuery.FullDenom.CreatorAddr
			subdenom := contractQuery.FullDenom.Subdenom

			fullDenom, err := GetFullDenom(creator, subdenom)
			if err != nil {
				return nil, errorsmod.Wrap(err, "osmo full denom query")
			}

			res := bindings.FullDenomResponse{
				Denom: fullDenom,
			}

			bz, err := json.Marshal(res)
			if err != nil {
				return nil, errorsmod.Wrap(err, "failed to marshal FullDenomResponse")
			}

			return bz, nil

		case contractQuery.Admin != nil:
			res, err := queryPlugin.GetDenomAdmin(ctx, contractQuery.Admin.Denom)
			if err != nil {
				return nil, err
			}

			bz, err := json.Marshal(res)
			if err != nil {
				return nil, fmt.Errorf("failed to JSON marshal AdminResponse: %w", err)
			}

			return bz, nil

		case contractQuery.Metadata != nil:
			res, err := queryPlugin.GetMetadata(ctx, contractQuery.Metadata.Denom)
			if err != nil {
				return nil, err
			}

			bz, err := json.Marshal(res)
			if err != nil {
				return nil, fmt.Errorf("failed to JSON marshal MetadataResponse: %w", err)
			}

			return bz, nil

		case contractQuery.DenomsByCreator != nil:
			res, err := queryPlugin.GetDenomsByCreator(ctx, contractQuery.DenomsByCreator.Creator)
			if err != nil {
				return nil, err
			}

			bz, err := json.Marshal(res)
			if err != nil {
				return nil, fmt.Errorf("failed to JSON marshal DenomsByCreatorResponse: %w", err)
			}

			return bz, nil

		case contractQuery.Params != nil:
			res, err := queryPlugin.GetParams(ctx)
			if err != nil {
				return nil, err
			}

			bz, err := json.Marshal(res)
			if err != nil {
				return nil, fmt.Errorf("failed to JSON marshal ParamsResponse: %w", err)
			}

			return bz, nil

		}
		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown App Data query variant"}
	}
}

// ConvertSdkCoinsToWasmCoins converts sdk type coins to wasm vm type coins
func ConvertSdkCoinsToWasmCoins(coins []sdk.Coin) wasmvmtypes.Coins {
	var toSend wasmvmtypes.Coins
	for _, coin := range coins {
		c := ConvertSdkCoinToWasmCoin(coin)
		toSend = append(toSend, c)
	}
	return toSend
}

// ConvertSdkCoinToWasmCoin converts a sdk type coin to a wasm vm type coin
func ConvertSdkCoinToWasmCoin(coin sdk.Coin) wasmvmtypes.Coin {
	return wasmvmtypes.Coin{
		Denom: coin.Denom,
		// Note: tokenfactory tokens have 18 decimal places, so 10^22 is common, no longer in u64 range
		Amount: coin.Amount.String(),
	}
}
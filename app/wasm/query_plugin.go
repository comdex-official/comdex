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
		} else if comdexQuery.RemoveWhiteListAssetLocker != nil {
			App_Id := comdexQuery.RemoveWhiteListAssetLocker.App_Id
			Asset_Id := comdexQuery.RemoveWhiteListAssetLocker.Asset_Id

			found, errormsg := queryPlugin.GetRemoveWhitelistAppIdLockerRewardsCheck(ctx, App_Id, Asset_Id)
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
			App_Id := comdexQuery.RemoveWhiteListAssetLocker.App_Id
			Asset_Id := comdexQuery.RemoveWhiteListAssetLocker.Asset_Id

			found, errormsg := queryPlugin.GetWhitelistAppIdLockerRewardsCheck(ctx, App_Id, Asset_Id)
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
			App_Id := comdexQuery.WhitelistAppIdVaultInterest.App_Id

			found, errormsg := queryPlugin.GetWhitelistAppIdVaultInterestCheck(ctx, App_Id)
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
			App_Id := comdexQuery.ExternalLockerRewards.App_Id
			Asset_Id := comdexQuery.ExternalLockerRewards.Asset_Id

			found, errormsg := queryPlugin.GetExternalLockerRewardsCheck(ctx, App_Id, Asset_Id)
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
			App_Id := comdexQuery.ExternalVaultRewards.App_Id
			Asset_Id := comdexQuery.ExternalVaultRewards.Asset_Id

			found, errormsg := queryPlugin.GetExternalVaultRewardsCheck(ctx, App_Id, Asset_Id)
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
			AppMappingId := comdexQuery.CollectorLookupTableQuery.AppMappingId
			CollectorAssetId := comdexQuery.CollectorLookupTableQuery.CollectorAssetId
			SecondaryAssetId := comdexQuery.CollectorLookupTableQuery.SecondaryAssetId
			found, errormsg := queryPlugin.CollectorLookupTableQueryCheck(ctx, AppMappingId, CollectorAssetId, SecondaryAssetId)
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
			AppMappingId := comdexQuery.ExtendedPairsVaultRecordsQuery.AppMappingId
			PairId := comdexQuery.ExtendedPairsVaultRecordsQuery.PairId
			StabilityFee := comdexQuery.ExtendedPairsVaultRecordsQuery.StabilityFee
			ClosingFee := comdexQuery.ExtendedPairsVaultRecordsQuery.ClosingFee
			DrawDownFee := comdexQuery.ExtendedPairsVaultRecordsQuery.DrawDownFee
			DebtCeiling := comdexQuery.ExtendedPairsVaultRecordsQuery.DebtCeiling
			DebtFloor := comdexQuery.ExtendedPairsVaultRecordsQuery.DebtFloor
			PairName := comdexQuery.ExtendedPairsVaultRecordsQuery.PairName

			found, errormsg := queryPlugin.ExtendedPairsVaultRecordsQueryCheck(ctx, AppMappingId, PairId, StabilityFee, ClosingFee, DrawDownFee, DebtCeiling, DebtFloor, PairName)
			res := bindings.ExtendedPairsVaultRecordsQueryResponse{
				Found: found,
				Err:   errormsg,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "ExternalVaultRewards query response")
			}
			return bz, nil
		}

		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown App Data query variant"}
	}

}

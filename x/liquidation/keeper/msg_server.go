package keeper

import (
	"context"
	"fmt"
	"github.com/comdex-official/comdex/x/liquidation/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	vaultypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (k msgServer) MsgLiquidateVault(c context.Context, msg *types.MsgLiquidateVaultRequest) (*types.MsgLiquidateVaultResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	appID, found := k.GetAppIDByAppForLiquidation(ctx, msg.AppId)
	if !found {
		return nil, types.ErrAppIDDoesNotExists
	}
	//Call Vault Data
	vault, found := k.vault.GetVault(ctx, msg.VaultId)
	if !found {
		return nil, vaultypes.ErrorVaultDoesNotExist
	}
	if vault.AppId != appID {
		return nil, fmt.Errorf("vault and app id mismatch for vault ID %d", vault.Id)
	}
	extPair, _ := k.asset.GetPairsVault(ctx, vault.ExtendedPairVaultID)
	pair, _ := k.asset.GetPair(ctx, extPair.PairId)
	assetIn, _ := k.asset.GetAsset(ctx, pair.AssetIn)

	totalRate, err := k.market.CalcAssetPrice(ctx, assetIn.Id, vault.AmountIn)
	if err != nil {
		return nil, err
	}
	totalIn := totalRate

	liqRatio := extPair.MinCr
	totalOut := vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated)
	collateralizationRatio, err := k.vault.CalculateCollateralizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, totalOut)
	if err != nil {
		return nil, err
	}
	if collateralizationRatio.LT(liqRatio) {
		// calculate interest and update vault
		totalDebt := vault.AmountOut.Add(vault.InterestAccumulated)
		err1 := k.rewards.CalculateVaultInterest(ctx, vault.AppId, vault.ExtendedPairVaultID, vault.Id, totalDebt, vault.BlockHeight, vault.BlockTime.Unix())
		if err1 != nil {
			return nil, err
		}
		vault, _ := k.vault.GetVault(ctx, vault.Id)
		totalFees := vault.InterestAccumulated.Add(vault.ClosingFeeAccumulated)
		totalOut := vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated)
		collateralizationRatio, err := k.vault.CalculateCollateralizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, totalOut)
		if err != nil {
			return nil, err
		}
		err = k.CreateLockedVault(ctx, vault, totalIn, collateralizationRatio, msg.AppId, totalFees)
		if err != nil {
			return nil, err
		}
		k.vault.DeleteVault(ctx, vault.Id)
		var rewards rewardstypes.VaultInterestTracker
		rewards.AppMappingId = msg.AppId
		rewards.VaultId = vault.Id
		k.rewards.DeleteVaultInterestTracker(ctx, rewards)
		k.vault.DeleteAddressFromAppExtendedPairVaultMapping(ctx, vault.ExtendedPairVaultID, vault.Id, msg.AppId)
	}
	return &types.MsgLiquidateVaultResponse{}, nil
}

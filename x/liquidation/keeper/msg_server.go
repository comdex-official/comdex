package keeper

import (
	"github.com/comdex-official/comdex/x/liquidation/types"
	vaultypes "github.com/comdex-official/comdex/x/vault/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServer(keeper Keeper) types.MsgServer {
	return &msgServer{
		Keeper: keeper,
	}
}


func (k msgServer) MsgLiquidateVaults(c context.Context, msg *types.MsgLiquidateVaultRequest) (*types.MsgLiquidateVaultResponse, error) {
	appId,found := k.GetAppIDByAppForLiquidation(ctx,appId)
	if !found{
		return types.ErrAppIDDoesNotExists
	}
	//Call Vault Data
	vault,found:=k.vault.GetVault(ctx,msg.VaultId)
	if !found{
		return vaultypes.ErrorVaultDoesNotExist
	}
	//
	if vault.AppId != appId {
		return fmt.Errorf("vault and app id mismatch in Liquidation, liquidate_vaults.go for vault ID %d", vault.Id)
	}
	extPair, _ := k.asset.GetPairsVault(ctx, vault.ExtendedPairVaultID)
	pair, _ := k.asset.GetPair(ctx, extPair.PairId)
	assetIn, _ := k.asset.GetAsset(ctx, pair.AssetIn)
	
	totalRate, err := k.market.CalcAssetPrice(ctx, assetIn.Id, vault.AmountIn)
	if err != nil {
		return fmt.Errorf("error in CalcAssetPrice in Liquidation, liquidate_vaults.go for vault ID %d", vault.Id)
	}
	totalIn := totalRate

	liqRatio := extPair.MinCr
	totalOut := vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated)
	collateralizationRatio, err := k.vault.CalculateCollateralizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, totalOut)
	if err != nil {
		return fmt.Errorf("error Calculating CR in Liquidation, liquidate_vaults.go for vault ID %d", vault.Id)
	}
	if collateralizationRatio.LT(liqRatio) {
		// calculate interest and update vault
		totalDebt := vault.AmountOut.Add(vault.InterestAccumulated)
		err1 := k.rewards.CalculateVaultInterest(ctx, vault.AppId, vault.ExtendedPairVaultID, vault.Id, totalDebt, vault.BlockHeight, vault.BlockTime.Unix())
		if err1 != nil {
			return fmt.Errorf("error Calculating vault interest in Liquidation, liquidate_vaults.go for vaultID %d", vault.Id)
		}
		vault, _ := k.vault.GetVault(ctx, vault.Id)
		totalFees := vault.InterestAccumulated.Add(vault.ClosingFeeAccumulated)
		totalOut := vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated)
		collateralizationRatio, err := k.vault.CalculateCollateralizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, totalOut)
		if err != nil {
			return fmt.Errorf("error Calculating CR in Liquidation, liquidate_vaults.go for vaultID %d", vault.Id)
		}
		err = k.CreateLockedVault(ctx, vault, totalIn, collateralizationRatio, appIds[i], totalFees)
		if err != nil {
			return fmt.Errorf("error Creating Locked Vaults in Liquidation, liquidate_vaults.go for Vault %d", vault.Id)
		}
		k.vault.DeleteVault(ctx, vault.Id)
		var rewards rewardstypes.VaultInterestTracker
		rewards.AppMappingId = appIds[i]
		rewards.VaultId = vault.Id
		k.rewards.DeleteVaultInterestTracker(ctx, rewards)
		k.vault.DeleteAddressFromAppExtendedPairVaultMapping(ctx, vault.ExtendedPairVaultID, vault.Id, appIds[i])
	}
	return nil
}






	return &types.MsgLiquidateVaultResponse{}, nil


}
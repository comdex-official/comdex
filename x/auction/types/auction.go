package types

import (
	"time"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewCollateralAuction(
	vault vaulttypes.Vault,
	collateralizationRatio sdk.Dec,
	assetIn assettypes.Asset,
	assetOut assettypes.Asset,
) CollateralAuction {
	auction := CollateralAuction{
		AuctionedVault: AuctionedVault{
			PairID:               vault.PairID,
			Owner:                vault.Owner,
			AmountIn:             vault.AmountIn,
			AmountOut:            vault.AmountOut,
			CrAtLiquidation:      collateralizationRatio,
			LiquidationTimestamp: time.Now().UTC(),
		},
		Lot:        sdk.NewCoin(assetIn.Denom, vault.AmountIn),
		Bidder:     nil,
		Bid:        sdk.NewCoin(assetOut.Denom, sdk.NewInt(0)),
		EndTime:    time.Now(),
		MaxEndTime: time.Now(),
		MaxBid:     sdk.NewCoin(assetOut.Denom, vault.AmountOut),
	}
	return auction
}

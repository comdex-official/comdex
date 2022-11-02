package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/app/wasm/bindings"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
)

func (k Keeper) GetUUSDFromUSD(ctx sdk.Context, price sdk.Dec) sdk.Dec {
	usdInUUSD := sdk.MustNewDecFromStr("1000000")
	return price.Mul(usdInUUSD)
}

func (k Keeper) GetModuleAccountBalance(ctx sdk.Context, moduleName string, denom string) sdk.Int {
	address := k.account.GetModuleAddress(moduleName)
	return k.bank.GetBalance(ctx, address, denom).Amount
}

func (k Keeper) FundModule(ctx sdk.Context, moduleName string, denom string, amt uint64) error {
	err := k.bank.MintCoins(ctx, moduleName, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewIntFromUint64(amt))))
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) AddAuctionParams(ctx sdk.Context, auctionParamsBinding *bindings.MsgAddAuctionParams) error {
	newStep := sdk.NewIntFromUint64(auctionParamsBinding.Step)
	auctionParams := auctiontypes.AuctionParams{
		AppId:                  auctionParamsBinding.AppID,
		AuctionDurationSeconds: auctionParamsBinding.AuctionDurationSeconds,
		Buffer:                 auctionParamsBinding.Buffer,
		Cusp:                   auctionParamsBinding.Cusp,
		Step:                   newStep,
		PriceFunctionType:      auctionParamsBinding.PriceFunctionType,
		SurplusId:              auctionParamsBinding.SurplusID,
		DebtId:                 auctionParamsBinding.DebtID,
		DutchId:                auctionParamsBinding.DutchID,
		BidDurationSeconds:     auctionParamsBinding.BidDurationSeconds,
	}

	k.SetAuctionParams(ctx, auctionParams)

	return nil
}

func (k Keeper) makeFalseForFlags(ctx sdk.Context, appID, assetID uint64) error {
	auctionLookupTable, found := k.collector.GetAuctionMappingForApp(ctx, appID, assetID)
	if !found {
		return auctiontypes.ErrorInvalidAddress
	}

	auctionLookupTable.IsAuctionActive = false
	err := k.collector.SetAuctionMappingForApp(ctx, auctionLookupTable)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) CalcDollarValueForToken(ctx sdk.Context, id uint64, rate sdk.Dec, amt sdk.Int) (price sdk.Dec, err error) {
	asset, found := k.asset.GetAsset(ctx, id)
	if !found {
		return sdk.ZeroDec(), assettypes.ErrorAssetDoesNotExist
	}

	numerator := sdk.NewDecFromInt(amt).Mul(rate)
	denominator := sdk.NewDecFromInt(asset.Decimals)
	return numerator.Quo(denominator), nil
}

package v7

import (
	"fmt"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func ReturnAtomToVaultOwners(ctx sdk.Context,
	bankKeeper bankkeeper.Keeper,
) {
	type auctionStruct struct {
		address    string
		vaultOwner string
		amount     int64
	}
	auctionData := []auctionStruct{
		{
			address:    "comdex1vdhk6er90qckxvekwpunxer9xcm8vmty0qm8xuee8pe8sam60fc8wdpnd5ekvwpkxd4rvam6awa4v8",
			vaultOwner: "comdex1c36py3de66vmdx6ss98rxwzzpw43m3f863j6wz",
			amount:     4413023,
		},
		{
			address:    "comdex1vdhk6er90qckzcmn0fjr26n5xd4hjd340fkxzet28qekcmfnxenrvmr48968wdeswfkhxap5g4ap6c",
			vaultOwner: "comdex1acszd5jt3ky65zlaej83lm36f6lu9tw70rmst4",
			amount:     23398233,
		},
		{
			address:    "comdex1vdhk6er90qchvwthxsuhyarywp4hw7r9dfjns7rw0quhqmnx89arj6rewgex2vr5wd4h2u3cg9xqnq",
			vaultOwner: "comdex1v9w49rtdpkwxeje8xnx9pnf9z9hyr2e0tskur8",
			amount:     8739931,
		},
		{
			address:    "comdex1vdhk6er90qchs7twwda8sufhwpnxxdmrxqukudnnxp3nw6rdv4nxuwphx4nh2v3kd3582am8hxr5t6",
			vaultOwner: "comdex1xynszxq7pfc7c09n6s0c7hmefn875gu26lhuwg",
			amount:     1859402,
		},
		{
			address:    "comdex1vdhk6er90qckuepe89nhvdenxajk2wr8vc6ks6rdw56hqmn8w4jxsvrhv3krgmredpcrxar2zupxa0",
			vaultOwner: "comdex1nd99gv737ee8gf5hhmu5pngudh0wdl4lyhp3tj",
			amount:     2518788,
		},
	}
	for i := 0; i < 5; i++ {
		atomCoins := sdk.NewCoin("ibc/961FA3E54F5DCCA639F37A7C45F7BBE41815579EF1513B5AFBEFCFEB8F256352", sdk.NewInt(auctionData[i].amount))
		addr1, err := sdk.AccAddressFromBech32(auctionData[i].address)
		if err != nil {
			fmt.Println("error in address", auctionData[i].address, err)
		}
		err = bankKeeper.SendCoinsFromAccountToModule(ctx, addr1, auctiontypes.ModuleName, sdk.NewCoins(atomCoins))
		if err != nil {
			fmt.Println("error in transfer to module of amt", auctionData[i].amount, err)
		}
		vaultOwner1, err := sdk.AccAddressFromBech32(auctionData[i].vaultOwner)
		if err != nil {
			fmt.Println("error in address of owner ", auctionData[i].vaultOwner, err)
		}
		if err := bankKeeper.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, vaultOwner1, sdk.NewCoins(atomCoins)); err != nil {
			fmt.Println("error in transfer to owner ", auctionData[i].vaultOwner, err)
		}
	}
}

func InitializeLendReservesStates(
	ctx sdk.Context,
	lendKeeper lendkeeper.Keeper,
) {
	dataAsset1, _ := lendKeeper.GetReserveBuybackAssetData(ctx, 1)
	dataAsset2, _ := lendKeeper.GetReserveBuybackAssetData(ctx, 2)
	dataAsset3, _ := lendKeeper.GetReserveBuybackAssetData(ctx, 3)
	reserveStat1 := types.AllReserveStats{
		AssetID:                        1,
		AmountOutFromReserveToLenders:  sdk.ZeroInt(),
		AmountOutFromReserveForAuction: sdk.ZeroInt(),
		AmountInFromLiqPenalty:         sdk.ZeroInt(),
		AmountInFromRepayments:         dataAsset1.BuybackAmount.Add(dataAsset1.ReserveAmount),
		TotalAmountOutToLenders:        sdk.ZeroInt(),
	}

	reserveStat2 := types.AllReserveStats{
		AssetID:                        2,
		AmountOutFromReserveToLenders:  sdk.ZeroInt(),
		AmountOutFromReserveForAuction: sdk.ZeroInt(),
		AmountInFromLiqPenalty:         sdk.ZeroInt(),
		AmountInFromRepayments:         dataAsset2.BuybackAmount.Add(dataAsset2.ReserveAmount),
		TotalAmountOutToLenders:        sdk.ZeroInt(),
	}

	reserveStat3 := types.AllReserveStats{
		AssetID:                        3,
		AmountOutFromReserveToLenders:  sdk.ZeroInt(),
		AmountOutFromReserveForAuction: sdk.ZeroInt(),
		AmountInFromLiqPenalty:         sdk.ZeroInt(),
		AmountInFromRepayments:         dataAsset3.BuybackAmount.Add(dataAsset3.ReserveAmount),
		TotalAmountOutToLenders:        sdk.ZeroInt(),
	}
	lendKeeper.SetAllReserveStatsByAssetID(ctx, reserveStat1)
	lendKeeper.SetAllReserveStatsByAssetID(ctx, reserveStat2)
	lendKeeper.SetAllReserveStatsByAssetID(ctx, reserveStat3)
	borrows := lendKeeper.GetAllBorrow(ctx)
	for _, v := range borrows {
		if v.IsLiquidated {
			pair, _ := lendKeeper.GetLendPair(ctx, v.PairID)
			assetStats, _ := lendKeeper.GetAssetStatsByPoolIDAndAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut)
			assetStats.TotalBorrowed = assetStats.TotalBorrowed.Sub(v.AmountOut.Amount)
			lendKeeper.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)
		}
	}
}

func CreateUpgradeHandler700(
	mm *module.Manager,
	configurator module.Configurator,
	lendKeeper lendkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		InitializeLendReservesStates(ctx, lendKeeper)
		ReturnAtomToVaultOwners(ctx, bankKeeper)
		return newVM, err
	}
}

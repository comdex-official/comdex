package v12

import (
	auctionkeeper "github.com/comdex-official/comdex/x/auctionsV2/keeper"
	auctionsV2types "github.com/comdex-official/comdex/x/auctionsV2/types"
	liquidationkeeper "github.com/comdex-official/comdex/x/liquidationsV2/keeper"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidationsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	icqkeeper "github.com/cosmos/ibc-apps/modules/async-icq/v4/keeper"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v4/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	collectorkeeper "github.com/comdex-official/comdex/x/collector/keeper"



		"fmt"


)

// An error occurred during the creation of the CMST/STJUNO pair, as it was mistakenly created in the Harbor app (ID-2) instead of the cSwap app (ID-1).
// As a result, the transaction fee was charged to the creator of the pair, who is entitled to a refund.
// The provided code is designed to initiate the refund process.
// The transaction hash for the pair creation is EF408AD53B8BB0469C2A593E4792CB45552BD6495753CC2C810A1E4D82F3982F.
// MintScan - https://www.mintscan.io/comdex/txs/EF408AD53B8BB0469C2A593E4792CB45552BD6495753CC2C810A1E4D82F3982F

func CreateUpgradeHandlerV12(
	mm *module.Manager,
	configurator module.Configurator,
	icqkeeper *icqkeeper.Keeper,
	liquidationKeeper liquidationkeeper.Keeper,
	auctionKeeper auctionkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	collectorKeeper collectorkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Applying main net upgrade - v.12.0.0")

		icqparams := icqtypes.DefaultParams()
		icqparams.AllowQueries = append(icqparams.AllowQueries, "/cosmwasm.wasm.v1.Query/SmartContractState")
		icqkeeper.SetParams(ctx, icqparams)

		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}
		InitializeStates(ctx, liquidationKeeper, auctionKeeper,bankKeeper,collectorKeeper)
		return vm, err
	}
}

func InitializeStates(
	ctx sdk.Context,
	liquidationKeeper liquidationkeeper.Keeper,
	auctionKeeper auctionkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	collectorKeeper collectorkeeper.Keeper,

) {
	dutchAuctionParams := liquidationtypes.DutchAuctionParam{
		Premium:         newDec("1.2"),
		Discount:        newDec("0.7"),
		DecrementFactor: sdk.NewInt(1),
	}
	englishAuctionParams := liquidationtypes.EnglishAuctionParam{DecrementFactor: sdk.NewInt(1)}

	harborParams := liquidationtypes.LiquidationWhiteListing{
		AppId:               2,
		Initiator:           false,
		IsDutchActivated:    true,
		DutchAuctionParam:   &dutchAuctionParams,
		IsEnglishActivated:  false,
		EnglishAuctionParam: &englishAuctionParams,
		KeeeperIncentive:    sdk.ZeroDec(),
	}

	commodoParams := liquidationtypes.LiquidationWhiteListing{
		AppId:               3,
		Initiator:           false,
		IsDutchActivated:    true,
		DutchAuctionParam:   &dutchAuctionParams,
		IsEnglishActivated:  false,
		EnglishAuctionParam: nil,
		KeeeperIncentive:    sdk.ZeroDec(),
	}

	liquidationKeeper.SetLiquidationWhiteListing(ctx, harborParams)
	liquidationKeeper.SetLiquidationWhiteListing(ctx, commodoParams)

	appReserveFundsTxDataHbr, found := liquidationKeeper.GetAppReserveFundsTxData(ctx, 2)
	if !found {
		appReserveFundsTxDataHbr.AppId = 2
	}
	appReserveFundsTxDataHbr.AssetTxData = append(appReserveFundsTxDataHbr.AssetTxData, liquidationtypes.AssetTxData{})
	liquidationKeeper.SetAppReserveFundsTxData(ctx, appReserveFundsTxDataHbr)

	appReserveFundsTxDataCmdo, found := liquidationKeeper.GetAppReserveFundsTxData(ctx, 3)
	if !found {
		appReserveFundsTxDataCmdo.AppId = 3
	}
	appReserveFundsTxDataCmdo.AssetTxData = append(appReserveFundsTxDataCmdo.AssetTxData, liquidationtypes.AssetTxData{})
	liquidationKeeper.SetAppReserveFundsTxData(ctx, appReserveFundsTxDataCmdo)

	auctionParams := auctionsV2types.AuctionParams{
		AuctionDurationSeconds: 18000,
		Step:                   newDec("0.1"),
		WithdrawalFee:          newDec("0.0"),
		ClosingFee:             newDec("0.0"),
		MinUsdValueLeft:        100000,
		BidFactor:              newDec("0.1"),
		LiquidationPenalty:     newDec("0.1"),
		AuctionBonus:           newDec("0.0"),
	}
	auctionKeeper.SetAuctionParams(ctx, auctionParams)
	auctionKeeper.SetParams(ctx, auctionsV2types.Params{})
	auctionKeeper.SetAuctionID(ctx, 0)
	auctionKeeper.SetUserBidID(ctx, 0)


	////// refund CMST to vault owner////////

	type refundStruct struct {
		vaultOwner    string
		amount     int64
	}

	refundData := []refundStruct{
		{
			vaultOwner: "comdex1x22fak2s8a6m9gysx7y4d5794dgds0jy6jch3t",
			amount:     27380000,
		},
		{
			vaultOwner: "comdex12jhse8d8uxgkqgrvfcv5j46wqu08yru7z3ze8z",
			amount:     1142650000,
		},
		{
			vaultOwner: "comdex1w5lep3d53p5dtkg37gerq6qxdlagykyryta989",
			amount:     4363010000,
		},
		{
			vaultOwner: "comdex122esu76xehp8sq9t88kcn666ejjum5g5ynxu0k",
			amount:     32460000,
		},
		{
			vaultOwner: "comdex1reeycz4d4pu4fddzqafzsyh6vvjp3nflp84xpp",
			amount:     44960000,
		},
		{
			vaultOwner: "comdex12q0708jnrd6d5ud7ap5lz4tgu3yshppfwd9x28",
			amount:     808240000,
		},
		{
			vaultOwner: "comdex120t6ntph3za6a7trw3zegseefkyf5u8gu3q4yu",
			amount:     29310000,
		},
		{
			vaultOwner: "comdex1qmklnue6z90vlljx04ll2v0elqjnzr3fswxm2u",
			amount:     10249670000,
		},
		{
			vaultOwner: "comdex13mm0ua6c20f8jup3q2g0uuw2k5n54cgkrw3lqs",
			amount:     664440000,
		},
		{
			vaultOwner: "comdex1wk25umx7ldgnca290dlg09yssusujhfek3l38l",
			amount:     2520920000,
		},
		{
			vaultOwner: "comdex1z2cmdk7atwfefl4a3had7a2tsamxrwgucmhutx",
			amount:     24300000,
		},
		{
			vaultOwner: "comdex1snezfskvsvdav5z9rsg5pgdrwnrg77kfjrc25f",
			amount:     23090000,
		},
		{
			vaultOwner: "comdex15xvnvwffhmy5wx8y7a9rchxe4zys9pa4gv8k8r",
			amount:     23650000,
		},
		{
			vaultOwner: "comdex1dwhhjyl6luv949ekpkplwc0zhqxa2jmhv6yl2w",
			amount:     19930000,
		},
		{
			vaultOwner: "comdex1nwtwhhs3d8rjl6c3clmcxlf3qdpv8n6rc9u9uy",
			amount:     18550000,
		},
		{
			vaultOwner: "comdex15gp4hjqf79zeggxteewzu2n0qde2zzfkkgec3z",
			amount:     79060000,
		},
		{
			vaultOwner: "comdex1v3truxzuz0j7896tumz77unla4sltqlgxwzhxy",
			amount:     45560000,
		},
		{
			vaultOwner: "comdex1850jsqvx54zl0urkav9tvee20j8r5fqj98zq9p",
			amount:     21940000,
		},
		{
			vaultOwner: "comdex1qx46s5gen6c88yaauh9jfttmfgdxnxxshzhahu",
			amount:     24400000,
		},
		
	}

	for i:=0; i<len(refundData); i++ {
		cmstCoins := sdk.NewCoin("ucmst", sdk.NewInt(refundData[i].amount))

		vaultOwner1, err := sdk.AccAddressFromBech32(refundData[i].vaultOwner)
		if err != nil {
			fmt.Println("error in address of owner ", refundData[i].vaultOwner, err)
		}
		
			if err := bankKeeper.SendCoinsFromModuleToAccount(ctx, collectortypes.ModuleName, vaultOwner1, sdk.NewCoins(cmstCoins)); err != nil {
			fmt.Println("error in transfer to owner ", refundData[i].vaultOwner, err)
		}

	}



}

func newDec(i string) sdk.Dec {
	dec, _ := sdk.NewDecFromStr(i)
	return dec
}
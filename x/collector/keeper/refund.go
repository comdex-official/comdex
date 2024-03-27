package keeper

import (
	"fmt"
	"github.com/comdex-official/comdex/x/collector/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	protobuftypes "github.com/cosmos/gogoproto/types"
)

func (k Keeper) Refund(ctx sdk.Context) error {
	////// refund CMST to vault owner////////

	type refundStruct struct {
		vaultOwner string
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

	// check if collector module account has enough balance to refund
	macc := k.ModuleBalance(ctx, collectortypes.ModuleName, "ucmst")
	// Check if sufficient balance exists

	if macc.Int64() < 20163520000 {
		fmt.Println("collector module account does not have enough balance to refund")
		return types.ErrorInsufficientBalance
	} else {
		for i := 0; i < len(refundData); i++ {
			cmstCoins := sdk.NewCoin("ucmst", sdk.NewInt(refundData[i].amount))

			vaultOwner1, err := sdk.AccAddressFromBech32(refundData[i].vaultOwner)
			if err != nil {
				fmt.Println("error in address of owner ", refundData[i].vaultOwner, err)
				return err
			}

			if err := k.bank.SendCoinsFromModuleToAccount(ctx, collectortypes.ModuleName, vaultOwner1, sdk.NewCoins(cmstCoins)); err != nil {
				fmt.Println("error in transfer to owner ", refundData[i].vaultOwner, err)
				return err
			}
			fmt.Println(`refund done for `, refundData[i].vaultOwner, ` amount `, refundData[i].amount)
		}

		// decrease net fee collected
		err := k.DecreaseNetFeeCollectedData(ctx, 2, 3, sdk.NewInt(20163520000))
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) SetRefundCounterStatus(ctx sdk.Context, id uint64) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.RefundCounterStatusPrefix
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) GetRefundCounterStatus(ctx sdk.Context) uint64 {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.RefundCounterStatusPrefix
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

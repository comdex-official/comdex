package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/comdex-official/comdex/x/asset/keeper"
	"github.com/comdex-official/comdex/x/asset/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// Simulation operation weights constants
const (
	DefaultWeightMsgAddAsset   int = 10
	DefaultWeightMsgAddToGauge int = 10
	OpWeightMsgAddAsset            = "op_weight_msg_add_asset"
	OpWeightMsgAddToGauge          = "op_weight_msg_add_to_gauge"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONCodec, ak govtypes.AccountKeeper,
	bk bankkeeper.Keeper, k keeper.Keeper) simulation.WeightedOperations {
	var (
		weightMsgAddAsset int
		// weightMsgAddToGauge int
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgAddAsset, &weightMsgAddAsset, nil,
		func(_ *rand.Rand) {
			weightMsgAddAsset = simappparams.DefaultWeightMsgCreateValidator
		},
	)

	// appParams.GetOrGenerate(cdc, OpWeightMsgAddToGauge, &weightMsgAddToGauge, nil,
	// 	func(_ *rand.Rand) {
	// 		weightMsgAddToGauge = simappparams.DefaultWeightMsgCreateValidator
	// 	},
	// )

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgAddAsset,
			SimulateMsgAddAsset(ak, bk, k),
		),
		// simulation.NewWeightedOperation(
		// 	weightMsgAddToGauge,
		// 	SimulateMsgAddToGaug(ak, bk, k),
		// ),
	}
}

func SimulateMsgAddAsset(ak govtypes.AccountKeeper, bk bankkeeper.Keeper, k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account,
		chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		//comdex tx asset add-asset gold ucGOLD 1000000 --from cooluser --chain-id test-1 --keyring-backend test -y
		//AddAssetRequest creates new Asset.

		// type MsgAddAssetRequest struct {
		//     From     string
		//     Name     string
		//     Denom    string
		//     Decimals int64

		simAccount, _ := simtypes.RandomAcc(r, accounts)
		account := ak.GetAccount(ctx, simAccount.Address)
		balance := bk.SpendableCoins(ctx, simAccount.Address)

		if balance.Len() <= 0 {
			return simtypes.NoOpMsg(
				types.ModuleName, "MsgAddAssetRequest", "Account does not have any coin"), nil, nil
		}

		if !balance.IsAnyNegative() {
			return simtypes.NoOpMsg(types.ModuleName, "MsgAddAssetRequest", "balance is negative"), nil, nil
		}

		asset := sdk.Coin{
			Denom:  "ucGOLD",
			Amount: sdk.Int(sdk.NewInt(1000000)),
		}

		if balance.AmountOf(asset.Denom).LT(asset.Amount) {
			return simtypes.NoOpMsg(types.ModuleName, "MsgAddAssetRequest", "Not enough funds to create asset"), nil, nil
		}

		balance = balance.Sub(sdk.NewCoins(asset))

		feeinucdmx := sdk.Coin{
			Denom:  "ucmdx",
			Amount: sdk.Int(sdk.NewInt(40000)),
		}
		fees := sdk.Coins{
			feeinucdmx,
		}

		if balance.AmountOf(feeinucdmx.Denom).LT(feeinucdmx.Amount) {
			return simtypes.NoOpMsg(types.ModuleName, "MsgAddAssetRequest", "Not enough funds for fees"), nil, nil
		}

		msg := types.NewMsgAddAssetRequest(sdk.AccAddress(account.GetAddress().String()), "gold", "ucGOLD", asset.Amount.Int64())

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, "MsgAddAssetRequest", "unable to generate mock tx"), nil, err
		}
		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, "MsgAddAssetRequest", "unable to deliver mock tx"), nil, err
		}
		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// // SimulateMsgAddPair creates new Asset Pair
// func SimulateMsgAddPair(ak govtypes.AccountKeeper, bk bankkeeper.Keeper, k keeper.Keeper) simtypes.Operation {
// 	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account,
// 		chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

// 		//AssetIn will be the user provided asset.

// 		// type MsgAddPairRequest struct {
// 		// 	From             string
// 		// 	AssetIn          uint64
// 		// 	AssetOut         uint64
// 		// 	LiquidationRatio types.Dec
// 		// }

// 		simAccount, _ := simtypes.RandomAcc(r, accounts)
// 		account := ak.GetAccount(ctx, simAccount.Address)
// 		balance := bk.SpendableCoins(ctx, simAccount.Address)
// 		if balance.Len() <= 0 {
// 			return simtypes.NoOpMsg(
// 				types.ModuleName, "MsgAddPairRequest", "Account does not have any coin"), nil, nil
// 		}
// 		//balance exists
// 		if !balance.IsAnyNegative() {
// 			return simtypes.NoOpMsg(types.ModuleName, "MsgAddPairRequest", "balance is negative"), nil, nil
// 		}

// 		assetin := sdk.Coin{
// 			Denom:  "ucmdx",
// 			Amount: sdk.Int(sdk.NewInt(100000)),
// 		}

// 		//check whether balance is less than deposit amount
// 		if balance.AmountOf(assetin.Denom).LT(assetin.Amount) {
// 			return simtypes.NoOpMsg(types.ModuleName, "MsgAddPairRequest", "not enough funds"), nil, nil
// 		}

// 		//then subtract deposit coins from balance
// 		balance = balance.Sub(sdk.NewCoins(assetin))

// 		k.SetParams(ctx, types.Params{})
// 		//declare fees
// 		feeinucdmx := sdk.Coin{
// 			Denom:  "ucmdx",
// 			Amount: sdk.Int(sdk.NewInt(4000)),
// 		}
// 		fees := sdk.Coins{
// 			feeinucdmx,
// 		}

// 		//check whether balance is less than fees
// 		if balance.AmountOf(deposit.Denom).LT(feeinucdmx.Amount) {
// 			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateRequest, "unable to generate fees"), nil, nil
// 		}

// 		//create the msg
// 		msg := types.NewMsgCreateRequest(simAccount.Address, 1, sdk.Int(sdk.NewInt(100000)), sdk.Int(sdk.NewInt(66666)))

// 		txGen := simappparams.MakeTestEncodingConfig().TxConfig
// 		tx, err := helpers.GenTx(
// 			txGen,
// 			[]sdk.Msg{msg},
// 			fees,
// 			helpers.DefaultGenTxGas,
// 			chainID,
// 			[]uint64{account.GetAccountNumber()},
// 			[]uint64{account.GetSequence()},
// 			simAccount.PrivKey,
// 		)

// 		if err != nil {
// 			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
// 		}
// 		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
// 		if err != nil {
// 			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver mock tx"), nil, err
// 		}
// 		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
// 	}
// }

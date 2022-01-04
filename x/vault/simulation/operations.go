package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/comdex-official/comdex/x/vault/keeper"
	"github.com/comdex-official/comdex/x/vault/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// Simulation operation weights constants
const (
	DefaultWeightMsgCreateGauge int = 10
	DefaultWeightMsgAddToGauge  int = 10
	OpWeightMsgCreateGauge          = "op_weight_msg_create_gauge"
	OpWeightMsgAddToGauge           = "op_weight_msg_add_to_gauge"
)

const (
	OpWeightMsgCreateDeployment = "op_weight_msg_create_deployment"
	OpWeightMsgUpdateDeployment = "op_weight_msg_update_deployment"
	OpWeightMsgCloseDeployment  = "op_weight_msg_close_deployment"
	OpWeightMsgCloseGroup       = "op_weight_msg_close_group"
)

const (
	// DefaultWeightMsgCreateProvider int = 100
	// DefaultWeightMsgUpdateProvider int = 5

	DefaultWeightMsgCreateVault   int = 100
	DefaultWeightMsgDepositVault  int = 10
	DefaultWeightMsgWithdrawVault int = 100

	// DefaultWeightMsgCreateBid  int = 100
	// DefaultWeightMsgCloseBid   int = 100
	// DefaultWeightMsgCloseLease int = 10
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONCodec, ak govtypes.AccountKeeper,
	bk bankkeeper.Keeper, k keeper.Keeper) simulation.WeightedOperations {
	var (
		weightMsgCreateVault   int
		weightMsgDepositVault  int
		weightMsgWithdrawVault int
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgCreateDeployment, &weightMsgCreateVault, nil, func(r *rand.Rand) {
			weightMsgCreateVault = DefaultWeightMsgCreateVault
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgUpdateDeployment, &weightMsgDepositVault, nil, func(r *rand.Rand) {
			weightMsgDepositVault = DefaultWeightMsgDepositVault
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgCloseDeployment, &weightMsgWithdrawVault, nil, func(r *rand.Rand) {
			weightMsgWithdrawVault = DefaultWeightMsgWithdrawVault
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgCreateVault,
			SimulateMsgCreateVault(ak, bk, k),
		),
		// simulation.NewWeightedOperation(
		// 	weightMsgDepositVault,
		// 	SimulateMsgDepositVault(ak, bk, k),
		// ),
		// simulation.NewWeightedOperation(
		// 	weightMsgWithdrawVault,
		// 	SimulateMsgWithdrawVault(ak, bk, k),
		// ),
	}
}

// SimulateMsgCreateVault generates a NewMsgCreateRequest with random values
func SimulateMsgCreateVault(ak govtypes.AccountKeeper, bk bankkeeper.Keeper, k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account,
		chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		//A new Vault is created, sender becomes the Vault owner.
		//Collateral is taken from the sender and sent to the Vault module account.
		//Principal stable coins are minted and sent to the sender.

		simAccount, _ := simtypes.RandomAcc(r, accounts)
		account := ak.GetAccount(ctx, simAccount.Address)
		balance := bk.SpendableCoins(ctx, simAccount.Address)
		if balance.Len() <= 0 {
			return simtypes.NoOpMsg(
				types.ModuleName, types.TypeMsgCreateRequest, "Account does not have any coin"), nil, nil
		}
		//balance exists
		if balance.IsAnyNegative() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateRequest, "balance is negative"), nil, nil
		}

		deposit := sdk.Coin{
			Denom:  "ucmdx",
			Amount: sdk.Int(sdk.NewInt(100000)),
		}

		//check whether balance is less than deposit amount
		if balance.AmountOf(deposit.Denom).LT(deposit.Amount) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateRequest, "not enough funds"), nil, nil
		}

		//then subtract deposit coins from balance
		balance = balance.Sub(sdk.NewCoins(deposit))

		//declare fees
		feeinucdmx := sdk.Coin{
			Denom:  "ucmdx",
			Amount: sdk.Int(sdk.NewInt(4000)),
		}
		fees := sdk.Coins{
			feeinucdmx,
		}

		//check whether balance is less than fees
		if balance.AmountOf(deposit.Denom).LT(feeinucdmx.Amount) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateRequest, "unable to generate fees"), nil, nil
		}

		//create the msg
		msg := types.NewMsgCreateRequest(simAccount.Address, 1, sdk.Int(sdk.NewInt(100000)), sdk.Int(sdk.NewInt(66666)))

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
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}
		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver mock tx"), nil, err
		}
		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgDepositVault generates a NewMsgCreateRequest with random values
func SimulateMsgDepositVault(ak govtypes.AccountKeeper, bk bankkeeper.Keeper, k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account,
		chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		// Deposit adds collateral to a Vault in the form of a deposit. Collateral is taken from owner.

		// type MsgDepositRequest struct {
		// 	From   string
		// 	ID     uint64
		// 	Amount types.Int
		// }
		// Resulting Changes:

		// Collateral is taken from the owner and sent to the Vault module account.
		// The depositor's Deposit Collateral struct is updated, or a new one created.

		simAccount, _ := simtypes.RandomAcc(r, accounts)
		account := ak.GetAccount(ctx, simAccount.Address)
		balance := bk.SpendableCoins(ctx, simAccount.Address)
		if balance.Len() <= 0 {
			return simtypes.NoOpMsg(
				types.ModuleName, types.TypeMsgDepositRequest, "Account does not have any coin"), nil, nil
		}
		//balance exists
		if !balance.IsAnyNegative() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositRequest, "balance is negative"), nil, nil
		}

		deposit := sdk.Coin{
			Denom:  "ucmdx",
			Amount: sdk.Int(sdk.NewInt(100000)),
		}

		//check whether balance is less than deposit amount
		if balance.AmountOf(deposit.Denom).LT(deposit.Amount) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositRequest, "not enough funds"), nil, nil
		}

		//then subtract deposit coins from balance
		balance = balance.Sub(sdk.NewCoins(deposit))

		//declare fees
		feeinucdmx := sdk.Coin{
			Denom:  "ucmdx",
			Amount: sdk.Int(sdk.NewInt(4000)),
		}
		fees := sdk.Coins{
			feeinucdmx,
		}

		//check whether balance is less than fees
		if balance.AmountOf(deposit.Denom).LT(feeinucdmx.Amount) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositRequest, "unable to generate fees"), nil, nil
		}

		//create the msg
		msg := types.NewMsgDepositRequest(simAccount.Address, 1, sdk.Int(sdk.NewInt(deposit.Amount.Int64())))

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
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}
		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver mock tx"), nil, err
		}
		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// SimulateMsgWithdrawVault generates a NewMsgCreateRequest with random values
func SimulateMsgWithdrawVault(ak govtypes.AccountKeeper, bk bankkeeper.Keeper, k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account,
		chainID string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		// Withdraw removes collateral asset from the vault, provided it does not put the Vault under the liquidation ratio.

		// type MsgWithdrawRequest struct {
		// 	From   string
		// 	ID     uint64
		// 	Amount types.Int
		// }
		// State Changes:

		// Collateral assets are sent from the vault to the owner.

		simAccount, _ := simtypes.RandomAcc(r, accounts)
		account := ak.GetAccount(ctx, simAccount.Address)
		balance := bk.SpendableCoins(ctx, simAccount.Address)
		if balance.Len() <= 0 {
			return simtypes.NoOpMsg(
				types.ModuleName, types.TypeMsgDepositRequest, "Account does not have any coin"), nil, nil
		}
		//balance exists
		if !balance.IsAnyNegative() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositRequest, "balance is negative"), nil, nil
		}

		deposit := sdk.Coin{
			Denom:  "ucmdx",
			Amount: sdk.Int(sdk.NewInt(100000)),
		}

		//check whether balance is less than deposit amount
		if balance.AmountOf(deposit.Denom).LT(deposit.Amount) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositRequest, "not enough funds"), nil, nil
		}

		//then subtract deposit coins from balance
		balance = balance.Sub(sdk.NewCoins(deposit))

		//declare fees
		feeinucdmx := sdk.Coin{
			Denom:  "ucmdx",
			Amount: sdk.Int(sdk.NewInt(4000)),
		}
		fees := sdk.Coins{
			feeinucdmx,
		}

		//check whether balance is less than fees
		if balance.AmountOf(deposit.Denom).LT(feeinucdmx.Amount) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDepositRequest, "unable to generate fees"), nil, nil
		}

		//create the msg
		msg := types.NewMsgWithdrawRequest(simAccount.Address, 1, sdk.Int(sdk.NewInt(deposit.Amount.Int64())))

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
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}
		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver mock tx"), nil, err
		}
		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

package keeper

import (
	"strconv"
	"strings"

	sdkerrors "cosmossdk.io/errors"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/comdex-official/comdex/x/gasless/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) EmitFeeConsumptionEvent(
	ctx sdk.Context,
	feeSource sdk.AccAddress,
	failedGasProviderIDs []uint64,
	failedGasProviderErrors []error,
	succeededGpid uint64,
) {
	failedGasProviderIDsStr := []string{}
	for _, id := range failedGasProviderIDs {
		failedGasProviderIDsStr = append(failedGasProviderIDsStr, strconv.FormatUint(id, 10))
	}
	failedGasProviderErrorMessages := []string{}
	for _, err := range failedGasProviderErrors {
		failedGasProviderErrorMessages = append(failedGasProviderErrorMessages, err.Error())
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeFeeConsumption,
			sdk.NewAttribute(types.AttributeKeyFeeSource, feeSource.String()),
			sdk.NewAttribute(types.AttributeKeyFailedGasProviderIDs, strings.Join(failedGasProviderIDsStr, ",")),
			sdk.NewAttribute(types.AttributeKeyFailedGasProviderErrors, strings.Join(failedGasProviderErrorMessages, ",")),
			sdk.NewAttribute(types.AttributeKeySucceededGpid, strconv.FormatUint(succeededGpid, 10)),
		),
	})
}

func (k Keeper) CanGasProviderBeUsedAsSource(ctx sdk.Context, gpid uint64, consumer types.GasConsumer, fee sdk.Coin) (gasProvider types.GasProvider, isValid bool, err error) {
	gasProvider, found := k.GetGasProvider(ctx, gpid)
	// there is no gas provider with given id, likely impossible to happen
	// exists only as aditional check.
	if !found {
		return gasProvider, false, sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "gas provider not found")
	}

	// gas provider is not active and cannot be used as fee source
	if !gasProvider.IsActive {
		return gasProvider, false, sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "gas provider not active")
	}

	// fee denom does not match between gas provider and asked fee
	if fee.Denom != gasProvider.FeeDenom {
		return gasProvider, false, sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "denom mismatch between provier and asked fee")
	}

	// asked fee amount is more than the allowed fee usage for tx.
	if fee.Amount.GT(gasProvider.MaxFeeUsagePerTx) {
		return gasProvider, false, sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "fee amount more than allowed limit")
	}

	// insufficient reserve in the gas tank to fulfill the transaction fee
	gasTankReserveBalance := k.bankKeeper.GetBalance(ctx, gasProvider.GetGasTankReserveAddress(), gasProvider.FeeDenom)
	if gasTankReserveBalance.IsLT(fee) {
		return gasProvider, false, sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "funds insufficient in gas reserve tank")
	}

	// if there is no consumption for the consumer, indicates that consumer is new and 1st time visitor
	// and the consumer is considered as valid and gas provider can be used as fee source
	if consumer.Consumption == nil {
		return gasProvider, true, nil
	}

	// no need to check the consumption usage since there is no key available with given gas provider id
	// i.e the consumer has never used this gas reserve before and the first time visitor for the given gas provider
	if _, ok := consumer.Consumption[gasProvider.Id]; !ok {
		return gasProvider, true, nil
	}

	consumptionDetails := consumer.Consumption[gasProvider.Id]

	// consumer is blocked by the gas provider
	if consumptionDetails.IsBlocked {
		return gasProvider, false, sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "blocked by gas provider")
	}

	// consumer exhausted the transaction count limit, hence not eligible with given gas provider
	if consumptionDetails.TotalTxsMade >= consumptionDetails.TotalTxsAllowed {
		return gasProvider, false, sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "exhausted tx limit")
	}

	// if total fees consumed by the consumer is more than or equal to the allowed consumption
	// i.e consumer has exhausted its fee limit and hence is not eligible for the given provider
	totalFeeConsumption := consumptionDetails.TotalFeesConsumed.Add(fee)
	if totalFeeConsumption.IsGTE(consumptionDetails.TotalFeeConsumptionAllowed) {
		return gasProvider, false, sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "exhausted total fee usage or pending fee limit insufficient for tx")
	}

	return gasProvider, true, nil
}

func (k Keeper) GetFeeSource(ctx sdk.Context, sdkTx sdk.Tx, originalFeePayer sdk.AccAddress, fees sdk.Coins) sdk.AccAddress {
	if len(sdkTx.GetMsgs()) > 1 {
		k.EmitFeeConsumptionEvent(ctx, originalFeePayer, []uint64{}, []error{sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "multiple messages")}, 0)
		return originalFeePayer
	}

	// only one fee coin is supported, tx containing multiple coins as fees are not allowed.
	if len(fees) != 1 {
		k.EmitFeeConsumptionEvent(ctx, originalFeePayer, []uint64{}, []error{sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "asked fee != 1")}, 0)
		return originalFeePayer
	}

	fee := fees[0]

	msg := sdkTx.GetMsgs()[0]
	msgTypeURL := sdk.MsgTypeURL(msg)

	isContract := false
	var contractAddress string

	executeContractMessage, ok := msg.(*wasmtypes.MsgExecuteContract)
	if ok {
		isContract = true
		contractAddress = executeContractMessage.GetContract()
	}

	txIdentifier := msgTypeURL
	if isContract {
		txIdentifier = contractAddress
	}

	// check if there are any gas providers for given txIdentifier i.e msgTypeURL or Contract address
	// if there is no gas provider for the given identifier, fee source will be original feePayer
	txGpids, found := k.GetTxGPIDS(ctx, txIdentifier)
	if !found {
		k.EmitFeeConsumptionEvent(ctx, originalFeePayer, []uint64{}, []error{sdkerrors.Wrapf(types.ErrorFeeConsumptionFailure, "no gas providers found")}, 0)
		return originalFeePayer
	}

	tempConsumer, found := k.GetGasConsumer(ctx, originalFeePayer)
	if !found {
		tempConsumer = types.NewGasConsumer(originalFeePayer)
	}

	failedGpids := []uint64{}
	failedGpidErrors := []error{}
	gasProvider := types.GasProvider{}
	isValid := false
	var err error
	gasProviderIds := txGpids.GasProviderIds
	for _, gpid := range gasProviderIds {
		gasProvider, isValid, err = k.CanGasProviderBeUsedAsSource(ctx, gpid, tempConsumer, fee)
		if isValid {
			break
		}
		failedGpidErrors = append(failedGpidErrors, err)
		failedGpids = append(failedGpids, gpid)
	}

	if !isValid {
		k.EmitFeeConsumptionEvent(ctx, originalFeePayer, failedGpids, failedGpidErrors, 0)
		return originalFeePayer
	}

	// update the consumption and usage details of the consumer
	gasConsumer := k.GetOrCreateGasConsumer(ctx, originalFeePayer, gasProvider)
	gasConsumer.Consumption[gasProvider.Id].TotalTxsMade = gasConsumer.Consumption[gasProvider.Id].TotalTxsMade + 1
	gasConsumer.Consumption[gasProvider.Id].TotalFeesConsumed = gasConsumer.Consumption[gasProvider.Id].TotalFeesConsumed.Add(fee)

	usage := gasConsumer.Consumption[gasProvider.Id].Usage
	if isContract {
		if usage.Contracts == nil {
			usage.Contracts = make(map[string]*types.UsageDetails)
		}
		if _, ok := usage.Contracts[contractAddress]; !ok {
			usage.Contracts[contractAddress] = &types.UsageDetails{}
		}
		usage.Contracts[contractAddress].Details = append(usage.Contracts[contractAddress].Details, &types.UsageDetail{
			Timestamp:   ctx.BlockTime(),
			GasConsumed: fee,
		})
	} else {
		if usage.Txs == nil {
			usage.Txs = make(map[string]*types.UsageDetails)
		}
		if _, ok := usage.Txs[msgTypeURL]; !ok {
			usage.Txs[msgTypeURL] = &types.UsageDetails{}
		}
		usage.Txs[msgTypeURL].Details = append(usage.Txs[msgTypeURL].Details, &types.UsageDetail{
			Timestamp:   ctx.BlockTime(),
			GasConsumed: fee,
		})
	}
	// assign the updated usage and set it to the store
	gasConsumer.Consumption[gasProvider.Id].Usage = usage
	k.SetGasConsumer(ctx, gasConsumer)

	// shift the used gas provider at the end of all providers, so that a different gas provider can be picked
	// in next cycle if there exists any.
	txGpids.GasProviderIds = types.ShiftToEndUint64(txGpids.GasProviderIds, gasProvider.Id)
	k.SetTxGPIDS(ctx, txGpids)

	feeSource := gasProvider.GetGasTankReserveAddress()
	k.EmitFeeConsumptionEvent(ctx, feeSource, failedGpids, failedGpidErrors, gasProvider.Id)

	return feeSource
}

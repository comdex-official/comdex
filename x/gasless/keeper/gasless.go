package keeper

import (
	"strconv"

	sdkerrors "cosmossdk.io/errors"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/comdex-official/comdex/x/gasless/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAvailableMessages(ctx sdk.Context) []string {
	return k.interfaceRegistry.ListImplementations("cosmos.base.v1beta1.Msg")
}

func (k Keeper) GetAllContractInfos(ctx sdk.Context) (contractInfos []wasmtypes.ContractInfo) {
	contractInfos = []wasmtypes.ContractInfo{}
	k.wasmKeeper.IterateContractInfo(ctx, func(aa sdk.AccAddress, ci wasmtypes.ContractInfo) bool {
		contractInfos = append(contractInfos, ci)
		return false
	})
	return contractInfos
}

func (k Keeper) GetAllContractsByCode(ctx sdk.Context, codeID uint64) (contracts []string) {
	contracts = []string{}
	k.wasmKeeper.IterateContractsByCode(ctx, codeID, func(address sdk.AccAddress) bool {
		contracts = append(contracts, address.String())
		return false
	})
	return contracts
}

func (k Keeper) GetAllAvailableContracts(ctx sdk.Context) (contractsDetails []types.ContractDetails) {
	contractsDetails = []types.ContractDetails{}
	contractInfos := k.GetAllContractInfos(ctx)
	for _, ci := range contractInfos {
		contracts := k.GetAllContractsByCode(ctx, ci.CodeID)
		for _, c := range contracts {
			contractsDetails = append(contractsDetails, types.ContractDetails{
				CodeId:  ci.CodeID,
				Address: c,
				Lable:   ci.Label,
			})
		}
	}
	return contractsDetails
}

func (k Keeper) ValidateMsgCreateGasProvider(ctx sdk.Context, msg *types.MsgCreateGasProvider) error {
	allGasProviders := k.GetAllGasProviders(ctx)
	gasTanks := 0
	for _, gp := range allGasProviders {
		if gp.Creator == msg.Creator {
			gasTanks++
		}
	}
	if gasTanks >= 10 {
		return sdkerrors.Wrapf(types.ErrorMaxLimitReachedByCreator, " %d gas tanks already created by the creator", 10)
	}

	if msg.FeeDenom != msg.GasDeposit.Denom {
		return sdkerrors.Wrapf(types.ErrorInvalidrequest, " fee denom %s do not match gas depoit denom %s ", msg.FeeDenom, msg.GasDeposit.Denom)
	}

	if !msg.MaxFeeUsagePerTx.IsPositive() {
		return sdkerrors.Wrapf(types.ErrorInvalidrequest, "max_fee_usage_per_tx should be positive")
	}
	if !msg.MaxFeeUsagePerConsumer.IsPositive() {
		return sdkerrors.Wrapf(types.ErrorInvalidrequest, "max_fee_usage_per_consumer should be positive")
	}

	if len(msg.TxsAllowed) == 0 && len(msg.ContractsAllowed) == 0 {
		return sdkerrors.Wrapf(types.ErrorInvalidrequest, "request should have atleast one tx path or contract address")
	}

	if len(msg.TxsAllowed) > 0 {
		allAvailableMessages := k.GetAvailableMessages(ctx)
		for _, message := range msg.TxsAllowed {
			if !types.ItemExists(allAvailableMessages, message) {
				return sdkerrors.Wrapf(types.ErrorInvalidrequest, "invalid message - %s", message)
			}
		}
	}

	if len(msg.ContractsAllowed) > 0 {
		allAvailableContractsDetails := k.GetAllAvailableContracts(ctx)
		contracts := []string{}
		for _, cdetails := range allAvailableContractsDetails {
			contracts = append(contracts, cdetails.Address)
		}
		for _, contract := range msg.ContractsAllowed {
			if !types.ItemExists(contracts, contract) {
				return sdkerrors.Wrapf(types.ErrorInvalidrequest, "invalid contract address - %s", contract)
			}
		}
	}

	if !msg.GasDeposit.IsPositive() {
		return sdkerrors.Wrapf(types.ErrorInvalidrequest, "deposit amount should be positive")
	}

	return nil
}

func (k Keeper) CreateGasProvider(ctx sdk.Context, msg *types.MsgCreateGasProvider) (types.GasProvider, error) {
	if err := k.ValidateMsgCreateGasProvider(ctx, msg); err != nil {
		return types.GasProvider{}, err
	}
	id := k.GetNextGasProviderIDWithUpdate(ctx)
	gasProvider := types.NewGasProvider(
		id,
		sdk.MustAccAddressFromBech32(msg.GetCreator()),
		msg.MaxTxsCountPerConsumer,
		msg.MaxFeeUsagePerConsumer,
		msg.MaxFeeUsagePerTx,
		msg.TxsAllowed,
		msg.ContractsAllowed,
		msg.FeeDenom,
	)

	k.SetGasProvider(ctx, gasProvider)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateGasProvider,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyGasProviderId, strconv.FormatUint(gasProvider.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyFeeDenom, msg.FeeDenom),
		),
	})

	return gasProvider, nil
}

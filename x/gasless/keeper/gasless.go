package keeper

import (
	"strconv"

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

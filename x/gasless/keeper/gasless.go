package keeper

import (
	"strconv"
	"strings"

	sdkerrors "cosmossdk.io/errors"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/comdex-official/comdex/x/gasless/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) GetAvailableMessages(_ sdk.Context) []string {
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

	if msg.MaxTxsCountPerConsumer == 0 {
		return sdkerrors.Wrap(types.ErrorInvalidrequest, "max tx count per consumer must not be 0")
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

	// Send gas deposit coins to the gas tank's reserve account.
	creator, err := sdk.AccAddressFromBech32(msg.GetCreator())
	if err != nil {
		return types.GasProvider{}, err
	}
	if err := k.bankKeeper.SendCoins(ctx, creator, gasProvider.GetGasTankReserveAddress(), sdk.NewCoins(msg.GasDeposit)); err != nil {
		return types.GasProvider{}, err
	}

	k.SetGasProvider(ctx, gasProvider)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateGasProvider,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyGasProviderId, strconv.FormatUint(gasProvider.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyFeeDenom, msg.FeeDenom),
			sdk.NewAttribute(types.AttributeKeyMaxFeeUsagePerTx, msg.MaxFeeUsagePerTx.String()),
			sdk.NewAttribute(types.AttributeKeyMaxTxsCountPerConsumer, strconv.FormatUint(msg.MaxTxsCountPerConsumer, 10)),
			sdk.NewAttribute(types.AttributeKeyMaxFeeUsagePerConsumer, msg.MaxFeeUsagePerConsumer.String()),
			sdk.NewAttribute(types.AttributeKeyTxsAllowed, strings.Join(gasProvider.TxsAllowed, ",")),
			sdk.NewAttribute(types.AttributeKeyContractsAllowed, strings.Join(gasProvider.ContractsAllowed, ",")),
		),
	})

	return gasProvider, nil
}

func (k Keeper) ValidateMsgAuthorizeActors(ctx sdk.Context, msg *types.MsgAuthorizeActors) error {
	gasProvider, found := k.GetGasProvider(ctx, msg.GasProviderId)
	if !found {
		return sdkerrors.Wrapf(errors.ErrNotFound, "gas provider with id %d not found", msg.GasProviderId)
	}

	if !gasProvider.IsActive {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas provider inactive")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Provider); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid provider address: %v", err)
	}

	if gasProvider.Creator != msg.Provider {
		return sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized provider")
	}

	for _, actor := range msg.Actors {
		if _, err := sdk.AccAddressFromBech32(actor); err != nil {
			return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid actor address - %s : %v", actor, err)
		}
	}

	return nil
}

func (k Keeper) AuthorizeActors(ctx sdk.Context, msg *types.MsgAuthorizeActors) (types.GasProvider, error) {
	if err := k.ValidateMsgAuthorizeActors(ctx, msg); err != nil {
		return types.GasProvider{}, err
	}

	gasProvider, _ := k.GetGasProvider(ctx, msg.GasProviderId)
	gasProvider.AuthorizedActors = types.RemoveDuplicates(msg.Actors)

	k.SetGasProvider(ctx, gasProvider)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAuthorizeActors,
			sdk.NewAttribute(types.AttributeKeyProvider, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyGasProviderId, strconv.FormatUint(gasProvider.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyAuthorizedActors, strings.Join(msg.Actors, ",")),
		),
	})

	return gasProvider, nil
}

func (k Keeper) ValidatMsgUpdateGasProviderStatus(ctx sdk.Context, msg *types.MsgUpdateGasProviderStatus) error {
	gasProvider, found := k.GetGasProvider(ctx, msg.GasProviderId)
	if !found {
		return sdkerrors.Wrapf(errors.ErrNotFound, "gas provider with id %d not found", msg.GasProviderId)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Provider); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid provider address: %v", err)
	}

	if gasProvider.Creator != msg.Provider {
		return sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized provider")
	}
	return nil
}

func (k Keeper) UpdateGasProviderStatus(ctx sdk.Context, msg *types.MsgUpdateGasProviderStatus) (types.GasProvider, error) {
	if err := k.ValidatMsgUpdateGasProviderStatus(ctx, msg); err != nil {
		return types.GasProvider{}, err
	}
	gasProvider, _ := k.GetGasProvider(ctx, msg.GasProviderId)
	gasProvider.IsActive = !gasProvider.IsActive

	k.SetGasProvider(ctx, gasProvider)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateGasProviderStatus,
			sdk.NewAttribute(types.AttributeKeyProvider, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyGasProviderId, strconv.FormatUint(gasProvider.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyGasProviderStatus, strconv.FormatBool(gasProvider.IsActive)),
		),
	})

	return gasProvider, nil
}

func (k Keeper) ValidateMsgUpdateGasProviderConfig(ctx sdk.Context, msg *types.MsgUpdateGasProviderConfig) error {
	gasProvider, found := k.GetGasProvider(ctx, msg.GasProviderId)
	if !found {
		return sdkerrors.Wrapf(errors.ErrNotFound, "gas provider with id %d not found", msg.GasProviderId)
	}

	if !gasProvider.IsActive {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas provider inactive")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Provider); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid provider address: %v", err)
	}

	if gasProvider.Creator != msg.Provider {
		return sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized provider")
	}

	if msg.MaxTxsCountPerConsumer == 0 {
		return sdkerrors.Wrap(types.ErrorInvalidrequest, "max tx count per consumer must not be 0")
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

	return nil
}

func (k Keeper) UpdateGasProviderConfig(ctx sdk.Context, msg *types.MsgUpdateGasProviderConfig) (types.GasProvider, error) {
	if err := k.ValidateMsgUpdateGasProviderConfig(ctx, msg); err != nil {
		return types.GasProvider{}, err
	}

	gasProvider, _ := k.GetGasProvider(ctx, msg.GasProviderId)
	gasProvider.MaxFeeUsagePerTx = msg.MaxFeeUsagePerTx
	gasProvider.MaxTxsCountPerConsumer = msg.MaxTxsCountPerConsumer
	gasProvider.MaxFeeUsagePerConsumer = msg.MaxFeeUsagePerConsumer
	gasProvider.TxsAllowed = types.RemoveDuplicates(msg.TxsAllowed)
	gasProvider.ContractsAllowed = types.RemoveDuplicates(msg.ContractsAllowed)

	k.SetGasProvider(ctx, gasProvider)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateGasProviderConfig,
			sdk.NewAttribute(types.AttributeKeyProvider, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyGasProviderId, strconv.FormatUint(gasProvider.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyMaxFeeUsagePerTx, msg.MaxFeeUsagePerTx.String()),
			sdk.NewAttribute(types.AttributeKeyMaxTxsCountPerConsumer, strconv.FormatUint(msg.MaxTxsCountPerConsumer, 10)),
			sdk.NewAttribute(types.AttributeKeyMaxFeeUsagePerConsumer, msg.MaxFeeUsagePerConsumer.String()),
			sdk.NewAttribute(types.AttributeKeyTxsAllowed, strings.Join(gasProvider.TxsAllowed, ",")),
			sdk.NewAttribute(types.AttributeKeyContractsAllowed, strings.Join(gasProvider.ContractsAllowed, ",")),
		),
	})

	return gasProvider, nil
}

func (k Keeper) ValidateMsgBlockConsumer(ctx sdk.Context, msg *types.MsgBlockConsumer) error {
	gasProvider, found := k.GetGasProvider(ctx, msg.GasProviderId)
	if !found {
		return sdkerrors.Wrapf(errors.ErrNotFound, "gas provider with id %d not found", msg.GasProviderId)
	}

	if !gasProvider.IsActive {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas provider inactive")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Actor); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid actor address: %v", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Consumer); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid consumer address: %v", err)
	}

	authorizedActors := gasProvider.AuthorizedActors
	authorizedActors = append(authorizedActors, gasProvider.Creator)

	if !types.ItemExists(authorizedActors, msg.Actor) {
		return sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized actor")
	}
	return nil
}

func (k Keeper) BlockConsumer(ctx sdk.Context, msg *types.MsgBlockConsumer) (types.GasConsumer, error) {
	if err := k.ValidateMsgBlockConsumer(ctx, msg); err != nil {
		return types.GasConsumer{}, err
	}

	gasProvider, _ := k.GetGasProvider(ctx, msg.GasProviderId)
	gasConsumer := k.GetOrCreateGasConsumer(ctx, sdk.MustAccAddressFromBech32(msg.Consumer), gasProvider)
	gasConsumer.Consumption[msg.GasProviderId].IsBlocked = true
	k.SetGasConsumer(ctx, gasConsumer)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBlockConsumer,
			sdk.NewAttribute(types.AttributeKeyActor, msg.Actor),
			sdk.NewAttribute(types.AttributeKeyConsumer, msg.Consumer),
			sdk.NewAttribute(types.AttributeKeyGasProviderId, strconv.FormatUint(msg.GasProviderId, 10)),
		),
	})

	return gasConsumer, nil
}

func (k Keeper) ValidateMsgUnblockConsumer(ctx sdk.Context, msg *types.MsgUnblockConsumer) error {
	gasProvider, found := k.GetGasProvider(ctx, msg.GasProviderId)
	if !found {
		return sdkerrors.Wrapf(errors.ErrNotFound, "gas provider with id %d not found", msg.GasProviderId)
	}

	if !gasProvider.IsActive {
		return sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas provider inactive")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Actor); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid actor address: %v", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Consumer); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid consumer address: %v", err)
	}

	authorizedActors := gasProvider.AuthorizedActors
	authorizedActors = append(authorizedActors, gasProvider.Creator)

	if !types.ItemExists(authorizedActors, msg.Actor) {
		return sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized actor")
	}
	return nil
}

func (k Keeper) UnblockConsumer(ctx sdk.Context, msg *types.MsgUnblockConsumer) (types.GasConsumer, error) {
	if err := k.ValidateMsgUnblockConsumer(ctx, msg); err != nil {
		return types.GasConsumer{}, err
	}

	gasProvider, _ := k.GetGasProvider(ctx, msg.GasProviderId)
	gasConsumer := k.GetOrCreateGasConsumer(ctx, sdk.MustAccAddressFromBech32(msg.Consumer), gasProvider)
	gasConsumer.Consumption[msg.GasProviderId].IsBlocked = false
	k.SetGasConsumer(ctx, gasConsumer)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnblockConsumer,
			sdk.NewAttribute(types.AttributeKeyActor, msg.Actor),
			sdk.NewAttribute(types.AttributeKeyConsumer, msg.Consumer),
			sdk.NewAttribute(types.AttributeKeyGasProviderId, strconv.FormatUint(msg.GasProviderId, 10)),
		),
	})

	return gasConsumer, nil
}

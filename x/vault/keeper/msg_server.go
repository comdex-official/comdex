package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	"github.com/comdex-official/comdex/x/vault/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

func NewMsgServer(keeper Keeper) types.MsgServer {
	return &msgServer{
		Keeper: keeper,
	}
}

// MsgCreate Creating a new CDP.
func (k msgServer) MsgCreate(c context.Context, msg *types.MsgCreateRequest) (*types.MsgCreateResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	esmStatus, found := k.esm.GetESMStatus(ctx, msg.AppId)
	status := false
	if found {
		status = esmStatus.Status
	}
	if status {
		return nil, esmtypes.ErrESMAlreadyExecuted
	}
	killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, msg.AppId)
	if killSwitchParams.BreakerEnable {
		return nil, esmtypes.ErrCircuitBreakerEnabled
	}
	extendedPairVault, found := k.asset.GetPairsVault(ctx, msg.ExtendedPairVaultId)
	if !found {
		return nil, types.ErrorExtendedPairVaultDoesNotExists
	}
	pairData, found := k.asset.GetPair(ctx, extendedPairVault.PairId)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}
	assetInData, found := k.asset.GetAsset(ctx, pairData.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	assetOutData, found := k.asset.GetAsset(ctx, pairData.AssetOut)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	appMapping, found := k.asset.GetApp(ctx, msg.AppId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}

	if appMapping.Id != extendedPairVault.AppId {
		return nil, types.ErrorAppMappingIDMismatch
	}

	depositorAddress, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	// Checking if this is a stableMint pair or not  -- stableMintPair == psmPair
	if extendedPairVault.IsStableMintVault {
		return nil, types.ErrorCannotCreateStableMintVault
	}
	// Checking
	if !extendedPairVault.IsVaultActive {
		return nil, types.ErrorVaultCreationInactive
	}
	// if does then check app to extendedPair mapping has any vault key
	// if it does throw error
	_, userExists := k.GetUserAppExtendedPairMappingData(ctx, msg.From, msg.AppId, msg.ExtendedPairVaultId)

	if userExists {
		// _, alreadyExists := k.CheckUserAppToExtendedPairMapping(ctx, userVaultExtendedPairMapping, extendedPairVault.Id, appMapping.Id)
		return nil, types.ErrorUserVaultAlreadyExists
	}
	// Call CheckAppExtendedPairVaultMapping function to get counter - it also initialised the kv store if appMapping_id does not exists, or extendedPairVault_id does not exists.
	tokenMintedStatistics, _ := k.CheckAppExtendedPairVaultMapping(ctx, appMapping.Id, extendedPairVault.Id)
	// Check debt Floor
	if !msg.AmountOut.GTE(extendedPairVault.DebtFloor) {
		return nil, types.ErrorAmountOutLessThanDebtFloor
	}
	// Check Debt Ceil
	currentMintedStatistics := tokenMintedStatistics.Add(msg.AmountOut)

	if currentMintedStatistics.GT(extendedPairVault.DebtCeiling) {
		return nil, types.ErrorAmountOutGreaterThanDebtCeiling
	}

	// Calculate CR - make necessary changes to calculate collateralization function
	if err := k.VerifyCollaterlizationRatio(ctx, extendedPairVault.Id, msg.AmountIn, msg.AmountOut, extendedPairVault.MinCr, status); err != nil {
		return nil, err
	}
	// Take amount from user
	if msg.AmountIn.GT(sdk.ZeroInt()) {
		if err := k.bank.SendCoinsFromAccountToModule(ctx, depositorAddress, types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetInData.Denom, msg.AmountIn))); err != nil {
			return nil, err
		}
	}

	// Mint Tokens for user
	mintCoin := sdk.NewCoin(assetOutData.Denom, msg.AmountOut)
	if mintCoin.IsZero() {
		return nil, types.MintCoinValueInVaultIsZero
	}
	if err := k.bank.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mintCoin)); err != nil {
		return nil, err
	}

	// Send Fees to Accumulator
	// Deducting Opening Fee if 0 opening fee then act accordingly
	if extendedPairVault.DrawDownFee.IsZero() && msg.AmountOut.GT(sdk.ZeroInt()) { // Send Rest to user
		if err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositorAddress, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, msg.AmountOut))); err != nil {
			return nil, err
		}
	} else {
		// If not zero deduct send to collector//////////
		// one approach could be
		collectorShare := sdk.NewDecFromInt(msg.AmountOut).Mul(extendedPairVault.DrawDownFee).TruncateInt()

		if collectorShare.GT(sdk.ZeroInt()) {
			if err := k.bank.SendCoinsFromModuleToModule(ctx, types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, collectorShare))); err != nil {
				return nil, err
			}

			err := k.collector.UpdateCollector(ctx, appMapping.Id, pairData.AssetOut, sdk.ZeroInt(), sdk.ZeroInt(), collectorShare, sdk.ZeroInt())
			if err != nil {
				return nil, err
			}
		}

		// and send the rest to the user
		amountToUser := msg.AmountOut.Sub(collectorShare)
		if amountToUser.GT(sdk.ZeroInt()) {
			if err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositorAddress, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, amountToUser))); err != nil {
				return nil, err
			}
		}
	}
	blockHeight := ctx.BlockHeight()
	blockTime := ctx.BlockTime()
	if extendedPairVault.StabilityFee.IsZero() {
		blockHeight = 0
	}

	// If all correct  create vault
	oldID := k.GetIDForVault(ctx)
	zeroVal := sdk.ZeroInt()
	var newVault types.Vault
	updatedID := oldID + 1
	newVault.Id = updatedID
	newVault.AmountIn = msg.AmountIn

	// closingFeeVal := msg.AmountOut.Mul(sdk.Int(extendedPairVault.ClosingFee)).Quo(sdk.Int(sdk.OneDec()))
	closingFeeVal := sdk.NewDecFromInt(msg.AmountOut).Mul(extendedPairVault.ClosingFee).TruncateInt()

	newVault.ClosingFeeAccumulated = closingFeeVal
	newVault.AmountOut = msg.AmountOut
	newVault.AppId = appMapping.Id
	newVault.InterestAccumulated = zeroVal
	newVault.Owner = msg.From
	newVault.CreatedAt = ctx.BlockTime()
	newVault.BlockHeight = blockHeight
	newVault.BlockTime = blockTime
	newVault.ExtendedPairVaultID = extendedPairVault.Id

	k.SetVault(ctx, newVault)
	k.SetIDForVault(ctx, updatedID)
	length := k.GetLengthOfVault(ctx)
	k.SetLengthOfVault(ctx, length+1)

	// Update mapping data - take proper approach
	// lookup table already exists
	// only need to update counter and token statistics value
	k.UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx, newVault)

	var mappingData types.OwnerAppExtendedPairVaultMappingData
	mappingData.Owner = msg.From
	mappingData.AppId = msg.AppId
	mappingData.ExtendedPairId = msg.ExtendedPairVaultId
	mappingData.VaultId = newVault.Id

	k.SetUserAppExtendedPairMappingData(ctx, mappingData)

	ctx.GasMeter().ConsumeGas(types.CreateVaultGas, "CreateVaultGas")

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateVault,
			sdk.NewAttribute(types.AttributeKeyVaultID, strconv.FormatUint(newVault.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyAppID, strconv.FormatUint(msg.AppId, 10)),
			sdk.NewAttribute(types.AttributeKeyExtendedPairID, strconv.FormatUint(msg.ExtendedPairVaultId, 10)),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.From),
			sdk.NewAttribute(types.AttributeKeyAmountIn, newVault.AmountIn.String()),
			sdk.NewAttribute(types.AttributeKeyAmountOut, newVault.AmountOut.String()),
			sdk.NewAttribute(types.AttributeKeyCreatedAt, ctx.BlockTime().String()),
			sdk.NewAttribute(types.AttributeKeyInterestAccumulated, newVault.InterestAccumulated.String()),
			sdk.NewAttribute(types.AttributeKeyClosingFeeAccumulated, newVault.ClosingFeeAccumulated.String()),
		),
	})

	return &types.MsgCreateResponse{}, nil
}

// MsgDeposit Only for depositing new collateral.
func (k msgServer) MsgDeposit(c context.Context, msg *types.MsgDepositRequest) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	esmStatus, found := k.esm.GetESMStatus(ctx, msg.AppId)
	status := false
	if found {
		status = esmStatus.Status
	}
	if status {
		return nil, esmtypes.ErrESMAlreadyExecuted
	}
	klwsParams, _ := k.esm.GetKillSwitchData(ctx, msg.AppId)
	if klwsParams.BreakerEnable {
		return nil, esmtypes.ErrCircuitBreakerEnabled
	}
	depositor, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	// checks if extended pair exists
	extendedPairVault, found := k.asset.GetPairsVault(ctx, msg.ExtendedPairVaultId)
	if !found {
		return nil, types.ErrorExtendedPairVaultDoesNotExists
	}
	pairData, found := k.asset.GetPair(ctx, extendedPairVault.PairId)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}
	assetInData, found := k.asset.GetAsset(ctx, pairData.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	// Checking if appMapping_id exists
	appMapping, found := k.asset.GetApp(ctx, msg.AppId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	// Checking if vault access disabled
	if !extendedPairVault.IsVaultActive {
		return nil, types.ErrorVaultInactive
	}

	// Checking if the appMapping_id in the msg_create & extendedPairVault_are same or not
	if appMapping.Id != extendedPairVault.AppId {
		return nil, types.ErrorAppMappingIDMismatch
	}

	userVault, found := k.GetVault(ctx, msg.UserVaultId)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	if userVault.Owner != msg.From {
		return nil, types.ErrVaultAccessUnauthorised
	}

	if appMapping.Id != userVault.AppId {
		return nil, types.ErrorInvalidAppMappingData
	}
	if extendedPairVault.Id != userVault.ExtendedPairVaultID {
		return nil, types.ErrorInvalidExtendedPairMappingData
	}

	totalDebt := userVault.AmountOut.Add(userVault.InterestAccumulated)
	err1 := k.rewards.CalculateVaultInterest(ctx, appMapping.Id, msg.ExtendedPairVaultId, msg.UserVaultId, totalDebt, userVault.BlockHeight, userVault.BlockTime.Unix())
	if err1 != nil {
		return nil, err1
	}
	userVault, found1 := k.GetVault(ctx, msg.UserVaultId)
	if !found1 {
		return nil, types.ErrorVaultDoesNotExist
	}
	userVault.AmountIn = userVault.AmountIn.Add(msg.Amount)
	if !userVault.AmountIn.IsPositive() {
		return nil, types.ErrorInvalidAmount
	}

	if msg.Amount.GT(sdk.ZeroInt()) {
		if err := k.bank.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetInData.Denom, msg.Amount))); err != nil {
			return nil, err
		}
	}
	userVault.BlockHeight = ctx.BlockHeight()
	userVault.BlockTime = ctx.BlockTime()

	k.SetVault(ctx, userVault)
	// Updating appExtendedPairvaultMappingData data -
	appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMappingData(ctx, appMapping.Id, msg.ExtendedPairVaultId)
	k.UpdateCollateralLockedAmountLockerMapping(ctx, appExtendedPairVaultData.AppId, appExtendedPairVaultData.ExtendedPairId, msg.Amount, true)

	ctx.GasMeter().ConsumeGas(types.DepositVaultGas, "DepositVaultGas")
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDepositVault,
			sdk.NewAttribute(types.AttributeKeyVaultID, strconv.FormatUint(msg.UserVaultId, 10)),
			sdk.NewAttribute(types.AttributeKeyAppID, strconv.FormatUint(msg.AppId, 10)),
			sdk.NewAttribute(types.AttributeKeyExtendedPairID, strconv.FormatUint(msg.ExtendedPairVaultId, 10)),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.From),
			sdk.NewAttribute(types.AttributeKeyAmountIn, msg.Amount.String()),
		),
	})

	return &types.MsgDepositResponse{}, nil
}

// MsgWithdraw Withdrawing collateral.
func (k msgServer) MsgWithdraw(c context.Context, msg *types.MsgWithdrawRequest) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, msg.AppId)
	if killSwitchParams.BreakerEnable {
		return nil, esmtypes.ErrCircuitBreakerEnabled
	}
	esmStatus, found := k.esm.GetESMStatus(ctx, msg.AppId)
	status := false
	if found {
		status = esmStatus.Status
	}

	if ctx.BlockTime().After(esmStatus.EndTime) && status {
		return nil, esmtypes.ErrCoolOffPeriodPassed
	}

	depositor, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	// checks if extended pair exists
	extendedPairVault, found := k.asset.GetPairsVault(ctx, msg.ExtendedPairVaultId)
	if !found {
		return nil, types.ErrorExtendedPairVaultDoesNotExists
	}
	pairData, found := k.asset.GetPair(ctx, extendedPairVault.PairId)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}
	assetInData, found := k.asset.GetAsset(ctx, pairData.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	// Checking if appMapping_id exists
	appMapping, found := k.asset.GetApp(ctx, msg.AppId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	// Checking if vault access disabled
	if !extendedPairVault.IsVaultActive {
		return nil, types.ErrorVaultInactive
	}
	// Checking if the appMapping_id in the msg_create & extendedPairVault_are same or not
	if appMapping.Id != extendedPairVault.AppId {
		return nil, types.ErrorAppMappingIDMismatch
	}

	userVault, found := k.GetVault(ctx, msg.UserVaultId)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	if userVault.Owner != msg.From {
		return nil, types.ErrVaultAccessUnauthorised
	}

	if appMapping.Id != userVault.AppId {
		return nil, types.ErrorInvalidAppMappingData
	}
	if extendedPairVault.Id != userVault.ExtendedPairVaultID {
		return nil, types.ErrorInvalidExtendedPairMappingData
	}

	totalDebt := userVault.AmountOut.Add(userVault.InterestAccumulated)
	err1 := k.rewards.CalculateVaultInterest(ctx, appMapping.Id, msg.ExtendedPairVaultId, msg.UserVaultId, totalDebt, userVault.BlockHeight, userVault.BlockTime.Unix())
	if err1 != nil {
		return nil, err1
	}

	userVault, found1 := k.GetVault(ctx, msg.UserVaultId)
	if !found1 {
		return nil, types.ErrorVaultDoesNotExist
	}
	userVault.AmountIn = userVault.AmountIn.Sub(msg.Amount)
	if !userVault.AmountIn.IsPositive() {
		return nil, types.ErrorInvalidAmount
	}

	totalDebtCalculation := userVault.AmountOut.Add(userVault.InterestAccumulated)
	totalDebtCalculation = totalDebtCalculation.Add(userVault.ClosingFeeAccumulated)

	// Calculate CR - make necessary changes to the calculate collateralization function
	if esmStatus.Status {
		totalDebtCalculation = userVault.AmountOut
	}
	if err := k.VerifyCollaterlizationRatio(ctx, extendedPairVault.Id, userVault.AmountIn, totalDebtCalculation, extendedPairVault.MinCr, status); err != nil {
		return nil, err
	}
	if msg.Amount.GT(sdk.ZeroInt()) {
		if err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositor, sdk.NewCoins(sdk.NewCoin(assetInData.Denom, msg.Amount))); err != nil {
			return nil, err
		}
	}
	userVault.BlockHeight = ctx.BlockHeight()
	userVault.BlockTime = ctx.BlockTime()
	k.SetVault(ctx, userVault)

	// Updating appExtendedPairVaultMappingData
	appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMappingData(ctx, appMapping.Id, msg.ExtendedPairVaultId)
	k.UpdateCollateralLockedAmountLockerMapping(ctx, appExtendedPairVaultData.AppId, appExtendedPairVaultData.ExtendedPairId, msg.Amount, false)

	ctx.GasMeter().ConsumeGas(types.WithdrawVaultGas, "WithdrawVaultGas")
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawVault,
			sdk.NewAttribute(types.AttributeKeyVaultID, strconv.FormatUint(msg.UserVaultId, 10)),
			sdk.NewAttribute(types.AttributeKeyAppID, strconv.FormatUint(msg.AppId, 10)),
			sdk.NewAttribute(types.AttributeKeyExtendedPairID, strconv.FormatUint(msg.ExtendedPairVaultId, 10)),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.From),
			sdk.NewAttribute(types.AttributeKeyAmountIn, msg.Amount.String()),
		),
	})

	return &types.MsgWithdrawResponse{}, nil
}

// MsgDraw To borrow more amount.
func (k msgServer) MsgDraw(c context.Context, msg *types.MsgDrawRequest) (*types.MsgDrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	esmStatus, found := k.esm.GetESMStatus(ctx, msg.AppId)
	status := false
	if found {
		status = esmStatus.Status
	}
	if status {
		return nil, esmtypes.ErrESMAlreadyExecuted
	}
	killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, msg.AppId)
	if killSwitchParams.BreakerEnable {
		return nil, esmtypes.ErrCircuitBreakerEnabled
	}
	depositor, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	// checks if extended pair exists
	extendedPairVault, found := k.asset.GetPairsVault(ctx, msg.ExtendedPairVaultId)
	if !found {
		return nil, types.ErrorExtendedPairVaultDoesNotExists
	}
	pairData, found := k.asset.GetPair(ctx, extendedPairVault.PairId)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}
	assetOutData, found := k.asset.GetAsset(ctx, pairData.AssetOut)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	// Checking if appMapping_id exists
	appMapping, found := k.asset.GetApp(ctx, msg.AppId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	// Checking if vault access disabled
	if !extendedPairVault.IsVaultActive {
		return nil, types.ErrorVaultInactive
	}
	// Checking if the appMapping_id in the msg_create & extendedPairVault_are same or not
	if appMapping.Id != extendedPairVault.AppId {
		return nil, types.ErrorAppMappingIDMismatch
	}

	userVault, found := k.GetVault(ctx, msg.UserVaultId)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	if userVault.Owner != msg.From {
		return nil, types.ErrVaultAccessUnauthorised
	}

	if appMapping.Id != userVault.AppId {
		return nil, types.ErrorInvalidAppMappingData
	}
	if extendedPairVault.Id != userVault.ExtendedPairVaultID {
		return nil, types.ErrorInvalidExtendedPairMappingData
	}
	if msg.Amount.LTE(sdk.NewInt(0)) {
		return nil, types.ErrorInvalidAmount
	}

	totalCalDebt := userVault.AmountOut.Add(userVault.InterestAccumulated)
	err1 := k.rewards.CalculateVaultInterest(ctx, appMapping.Id, msg.ExtendedPairVaultId, msg.UserVaultId, totalCalDebt, userVault.BlockHeight, userVault.BlockTime.Unix())
	if err1 != nil {
		return nil, err1
	}

	userVault, found1 := k.GetVault(ctx, msg.UserVaultId)
	if !found1 {
		return nil, types.ErrorVaultDoesNotExist
	}

	newUpdatedAmountOut := userVault.AmountOut.Add(msg.Amount)
	totalDebt := newUpdatedAmountOut.Add(userVault.InterestAccumulated)
	totalDebt = totalDebt.Add(userVault.ClosingFeeAccumulated)

	tokenMintedStatistics, _ := k.CheckAppExtendedPairVaultMapping(ctx, appMapping.Id, extendedPairVault.Id)

	// Check Debt Ceil
	currentMintedStatistics := tokenMintedStatistics.Add(msg.Amount)

	if currentMintedStatistics.GTE(extendedPairVault.DebtCeiling) {
		return nil, types.ErrorAmountOutGreaterThanDebtCeiling
	}

	if err := k.VerifyCollaterlizationRatio(ctx, extendedPairVault.Id, userVault.AmountIn, totalDebt, extendedPairVault.MinCr, status); err != nil {
		return nil, err
	}

	mintCoin := sdk.NewCoin(assetOutData.Denom, msg.Amount)
	if mintCoin.IsZero() {
		return nil, types.MintCoinValueInVaultIsZero
	}
	if err := k.bank.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mintCoin)); err != nil {
		return nil, err
	}

	if extendedPairVault.DrawDownFee.IsZero() && msg.Amount.GT(sdk.ZeroInt()) {
		// Send Rest to user
		if err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositor, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, msg.Amount))); err != nil {
			return nil, err
		}
	} else {
		// If not zero deduct send to collector//////////
		// one approach could be
		collectorShare := sdk.NewDecFromInt(msg.Amount).Mul(extendedPairVault.DrawDownFee).TruncateInt()

		if collectorShare.GT(sdk.ZeroInt()) {
			if err := k.bank.SendCoinsFromModuleToModule(ctx, types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, collectorShare))); err != nil {
				return nil, err
			}

			err := k.collector.UpdateCollector(ctx, appMapping.Id, pairData.AssetOut, sdk.ZeroInt(), sdk.ZeroInt(), collectorShare, sdk.ZeroInt())
			if err != nil {
				return nil, err
			}
		}
		// and send the rest to the user
		amountToUser := msg.Amount.Sub(collectorShare)
		if amountToUser.GT(sdk.ZeroInt()) {
			if err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositor, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, amountToUser))); err != nil {
				return nil, err
			}
		}
	}

	userVault.AmountOut = userVault.AmountOut.Add(msg.Amount)
	userVault.BlockHeight = ctx.BlockHeight()
	userVault.BlockTime = ctx.BlockTime()
	k.SetVault(ctx, userVault)

	// Updating appExtendedPairVaultMappingData
	appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMappingData(ctx, appMapping.Id, msg.ExtendedPairVaultId)
	k.UpdateTokenMintedAmountLockerMapping(ctx, appExtendedPairVaultData.AppId, appExtendedPairVaultData.ExtendedPairId, msg.Amount, true)

	ctx.GasMeter().ConsumeGas(types.DrawVaultGas, "DrawVaultGas")

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDrawVault,
			sdk.NewAttribute(types.AttributeKeyVaultID, strconv.FormatUint(msg.UserVaultId, 10)),
			sdk.NewAttribute(types.AttributeKeyAppID, strconv.FormatUint(msg.AppId, 10)),
			sdk.NewAttribute(types.AttributeKeyExtendedPairID, strconv.FormatUint(msg.ExtendedPairVaultId, 10)),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.From),
			sdk.NewAttribute(types.AttributeKeyAmountOut, msg.Amount.String()),
		),
	})

	return &types.MsgDrawResponse{}, nil
}

func (k msgServer) MsgRepay(c context.Context, msg *types.MsgRepayRequest) (*types.MsgRepayResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	esmStatus, found := k.esm.GetESMStatus(ctx, msg.AppId)
	status := false
	if found {
		status = esmStatus.Status
	}
	if status {
		return nil, esmtypes.ErrESMAlreadyExecuted
	}
	killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, msg.AppId)
	if killSwitchParams.BreakerEnable {
		return nil, esmtypes.ErrCircuitBreakerEnabled
	}
	depositor, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	// checks if extended pair exists
	extendedPairVault, found := k.asset.GetPairsVault(ctx, msg.ExtendedPairVaultId)
	if !found {
		return nil, types.ErrorExtendedPairVaultDoesNotExists
	}
	pairData, found := k.asset.GetPair(ctx, extendedPairVault.PairId)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}
	assetOutData, found := k.asset.GetAsset(ctx, pairData.AssetOut)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	// Checking if appMapping_id exists
	appMapping, found := k.asset.GetApp(ctx, msg.AppId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	// Checking if vault acccess disabled

	// Checking if the appMapping_id in the msg_create & extendedPairVault_are same or not
	if appMapping.Id != extendedPairVault.AppId {
		return nil, types.ErrorAppMappingIDMismatch
	}

	userVault, found := k.GetVault(ctx, msg.UserVaultId)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	if userVault.Owner != msg.From {
		return nil, types.ErrVaultAccessUnauthorised
	}

	if appMapping.Id != userVault.AppId {
		return nil, types.ErrorInvalidAppMappingData
	}
	if extendedPairVault.Id != userVault.ExtendedPairVaultID {
		return nil, types.ErrorInvalidExtendedPairMappingData
	}
	if msg.Amount.LTE(sdk.NewInt(0)) {
		return nil, types.ErrorInvalidAmount
	}

	totalDebt := userVault.AmountOut.Add(userVault.InterestAccumulated)
	err1 := k.rewards.CalculateVaultInterest(ctx, appMapping.Id, msg.ExtendedPairVaultId, msg.UserVaultId, totalDebt, userVault.BlockHeight, userVault.BlockTime.Unix())
	if err1 != nil {
		return nil, err1
	}

	userVault, found1 := k.GetVault(ctx, msg.UserVaultId)
	if !found1 {
		return nil, types.ErrorVaultDoesNotExist
	}

	newAmount := userVault.AmountOut.Add(userVault.InterestAccumulated)
	newAmount = newAmount.Sub(msg.Amount)
	if newAmount.LT(sdk.NewInt(0)) {
		return nil, types.ErrorInvalidAmount
	}

	if msg.Amount.LTE(userVault.InterestAccumulated) {
		// Amount is less than equal to the interest accumulated
		// subtract that as interest
		reducedFees := userVault.InterestAccumulated.Sub(msg.Amount)
		userVault.InterestAccumulated = reducedFees
		// and send it to the collector module
		if msg.Amount.GT(sdk.ZeroInt()) {
			if err := k.bank.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, msg.Amount))); err != nil {
				return nil, err
			}
			//			SEND TO COLLECTOR- msg.Amount
			if err := k.bank.SendCoinsFromModuleToModule(ctx, types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, msg.Amount))); err != nil {
				return nil, err
			}
			err := k.collector.UpdateCollector(ctx, appMapping.Id, pairData.AssetOut, msg.Amount, sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt())
			if err != nil {
				return nil, err
			}
		}
		userVault.BlockHeight = ctx.BlockHeight()
		userVault.BlockTime = ctx.BlockTime()
		k.SetVault(ctx, userVault)
	} else {
		updatedUserSentAmountAfterFeesDeduction := msg.Amount.Sub(userVault.InterestAccumulated)

		updatedUserDebt := userVault.AmountOut.Sub(updatedUserSentAmountAfterFeesDeduction)

		// //If user's closing fees is a bigger amount than the debt floor, user will not close the debt floor

		if !updatedUserDebt.GTE(extendedPairVault.DebtFloor) {
			return nil, types.ErrorAmountOutLessThanDebtFloor
		}
		if msg.Amount.GT(sdk.ZeroInt()) {
			if err := k.bank.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, msg.Amount))); err != nil {
				return nil, err
			}
		}
		if updatedUserSentAmountAfterFeesDeduction.GT(sdk.ZeroInt()) {
			burnCoin := sdk.NewCoin(assetOutData.Denom, updatedUserSentAmountAfterFeesDeduction)
			if burnCoin.IsZero() {
				return nil, types.BurnCoinValueInVaultIsZero
			}
			if err := k.bank.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(burnCoin)); err != nil {
				return nil, err
			}
		}
		//			SEND TO COLLECTOR----userVault.InterestAccumulated
		if userVault.InterestAccumulated.GT(sdk.ZeroInt()) {
			if err := k.bank.SendCoinsFromModuleToModule(ctx, types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, userVault.InterestAccumulated))); err != nil {
				return nil, err
			}
			err := k.collector.UpdateCollector(ctx, appMapping.Id, pairData.AssetOut, userVault.InterestAccumulated, sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt())
			if err != nil {
				return nil, err
			}
		}

		userVault.AmountOut = updatedUserDebt
		zeroVal := sdk.ZeroInt()
		userVault.InterestAccumulated = zeroVal
		userVault.BlockHeight = ctx.BlockHeight()
		userVault.BlockTime = ctx.BlockTime()
		k.SetVault(ctx, userVault)
		appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMappingData(ctx, appMapping.Id, msg.ExtendedPairVaultId)
		k.UpdateTokenMintedAmountLockerMapping(ctx, appExtendedPairVaultData.AppId, appExtendedPairVaultData.ExtendedPairId, updatedUserSentAmountAfterFeesDeduction, false)
	}

	ctx.GasMeter().ConsumeGas(types.RepayVaultGas, "RepayVaultGas")

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRepayVault,
			sdk.NewAttribute(types.AttributeKeyVaultID, strconv.FormatUint(msg.UserVaultId, 10)),
			sdk.NewAttribute(types.AttributeKeyAppID, strconv.FormatUint(msg.AppId, 10)),
			sdk.NewAttribute(types.AttributeKeyExtendedPairID, strconv.FormatUint(msg.ExtendedPairVaultId, 10)),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.From),
			sdk.NewAttribute(types.AttributeKeyAmountOut, msg.Amount.String()),
		),
	})

	return &types.MsgRepayResponse{}, nil
}

func (k msgServer) MsgClose(c context.Context, msg *types.MsgCloseRequest) (*types.MsgCloseResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	esmStatus, found := k.esm.GetESMStatus(ctx, msg.AppId)
	status := false
	if found {
		status = esmStatus.Status
	}
	if status {
		return nil, esmtypes.ErrESMAlreadyExecuted
	}
	killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, msg.AppId)
	if killSwitchParams.BreakerEnable {
		return nil, esmtypes.ErrCircuitBreakerEnabled
	}
	depositor, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	// checks if extended pair exists
	extendedPairVault, found := k.asset.GetPairsVault(ctx, msg.ExtendedPairVaultId)
	if !found {
		return nil, types.ErrorExtendedPairVaultDoesNotExists
	}
	pairData, found := k.asset.GetPair(ctx, extendedPairVault.PairId)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}
	assetInData, found := k.asset.GetAsset(ctx, pairData.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	assetOutData, found := k.asset.GetAsset(ctx, pairData.AssetOut)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	// Checking if appMapping_id exists
	appMapping, found := k.asset.GetApp(ctx, msg.AppId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	// //Checking if vault acccess disabled

	// Checking if the appMapping_id in the msg_create & extendedPairVault_are same or not
	if appMapping.Id != extendedPairVault.AppId {
		return nil, types.ErrorAppMappingIDMismatch
	}

	userVault, found := k.GetVault(ctx, msg.UserVaultId)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	if userVault.Owner != msg.From {
		return nil, types.ErrVaultAccessUnauthorised
	}

	if appMapping.Id != userVault.AppId {
		return nil, types.ErrorInvalidAppMappingData
	}
	if extendedPairVault.Id != userVault.ExtendedPairVaultID {
		return nil, types.ErrorInvalidExtendedPairMappingData
	}

	totalDebt := userVault.AmountOut.Add(userVault.InterestAccumulated)
	err1 := k.rewards.CalculateVaultInterest(ctx, appMapping.Id, msg.ExtendedPairVaultId, msg.UserVaultId, totalDebt, userVault.BlockHeight, userVault.BlockTime.Unix())
	if err1 != nil {
		return nil, err1
	}

	userVault, found1 := k.GetVault(ctx, msg.UserVaultId)
	if !found1 {
		return nil, types.ErrorVaultDoesNotExist
	}

	totalUserDebt := userVault.AmountOut.Add(userVault.InterestAccumulated)
	totalUserDebt = totalUserDebt.Add(userVault.ClosingFeeAccumulated)
	if totalUserDebt.GT(sdk.ZeroInt()) {
		if err := k.bank.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, totalUserDebt))); err != nil {
			return nil, err
		}
	}

	//			SEND TO COLLECTOR----userVault.InterestAccumulated & userVault.ClosingFees

	err = k.collector.UpdateCollector(ctx, appMapping.Id, pairData.AssetOut, userVault.InterestAccumulated, userVault.ClosingFeeAccumulated, sdk.ZeroInt(), sdk.ZeroInt())
	if err != nil {
		return nil, err
	}
	if userVault.InterestAccumulated.GT(sdk.ZeroInt()) {
		if err := k.bank.SendCoinsFromModuleToModule(ctx, types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, userVault.InterestAccumulated))); err != nil {
			return nil, err
		}
	}
	if userVault.ClosingFeeAccumulated.GT(sdk.ZeroInt()) {
		if err := k.bank.SendCoinsFromModuleToModule(ctx, types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, userVault.ClosingFeeAccumulated))); err != nil {
			return nil, err
		}
	}
	if userVault.AmountOut.GT(sdk.ZeroInt()) {
		burnCoin := sdk.NewCoin(assetOutData.Denom, userVault.AmountOut)
		if burnCoin.IsZero() {
			return nil, types.BurnCoinValueInVaultIsZero
		}
		if err := k.bank.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(burnCoin)); err != nil {
			return nil, err
		}
	}
	if userVault.AmountIn.GT(sdk.ZeroInt()) {
		if err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositor, sdk.NewCoins(sdk.NewCoin(assetInData.Denom, userVault.AmountIn))); err != nil {
			return nil, err
		}
	}

	// Update LookupTable minting Status
	appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMappingData(ctx, appMapping.Id, msg.ExtendedPairVaultId)

	k.UpdateCollateralLockedAmountLockerMapping(ctx, appExtendedPairVaultData.AppId, appExtendedPairVaultData.ExtendedPairId, userVault.AmountIn, false)
	k.UpdateTokenMintedAmountLockerMapping(ctx, appExtendedPairVaultData.AppId, appExtendedPairVaultData.ExtendedPairId, userVault.AmountOut, false)

	// Remove address from lookup table
	k.DeleteAddressFromAppExtendedPairVaultMapping(ctx, extendedPairVault.Id, userVault.Id, appMapping.Id)

	// Remove user extendedPair to address field in UserLookupStruct
	k.DeleteUserVaultExtendedPairMapping(ctx, msg.From, appMapping.Id, extendedPairVault.Id)

	// Delete Vault
	k.DeleteVault(ctx, userVault.Id)

	length := k.GetLengthOfVault(ctx)
	k.SetLengthOfVault(ctx, length-1)

	var rewards rewardstypes.VaultInterestTracker
	rewards.AppMappingId = appMapping.Id
	rewards.VaultId = userVault.Id
	k.rewards.DeleteVaultInterestTracker(ctx, rewards)

	ctx.GasMeter().ConsumeGas(types.CloseVaultGas, "CloseVaultGas")

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCloseVault,
			sdk.NewAttribute(types.AttributeKeyVaultID, strconv.FormatUint(userVault.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyAppID, strconv.FormatUint(msg.AppId, 10)),
			sdk.NewAttribute(types.AttributeKeyExtendedPairID, strconv.FormatUint(msg.ExtendedPairVaultId, 10)),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.From),
			sdk.NewAttribute(types.AttributeKeyCreatedAt, userVault.CreatedAt.String()),
			sdk.NewAttribute(types.AttributeKeyInterestAccumulated, userVault.InterestAccumulated.String()),
			sdk.NewAttribute(types.AttributeKeyClosingFeeAccumulated, userVault.ClosingFeeAccumulated.String()),
		),
	})

	return &types.MsgCloseResponse{}, nil
}

func (k msgServer) MsgDepositAndDraw(c context.Context, msg *types.MsgDepositAndDrawRequest) (*types.MsgDepositAndDrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	userVault, found := k.GetVault(ctx, msg.UserVaultId)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	newAmt := k.calculateUserToken(userVault, msg.Amount)
	msgDepositReq := types.MsgDepositRequest{
		From:                msg.From,
		AppId:               msg.AppId,
		ExtendedPairVaultId: msg.ExtendedPairVaultId,
		UserVaultId:         msg.UserVaultId,
		Amount:              msg.Amount,
	}
	_, err := k.MsgDeposit(c, &msgDepositReq)
	if err != nil {
		return nil, err
	}
	msgDrawReq := types.MsgDrawRequest{
		From:                msg.From,
		AppId:               msg.AppId,
		ExtendedPairVaultId: msg.ExtendedPairVaultId,
		UserVaultId:         msg.UserVaultId,
		Amount:              newAmt,
	}
	_, err = k.MsgDraw(c, &msgDrawReq)
	if err != nil {
		return nil, err
	}
	ctx.GasMeter().ConsumeGas(types.DepositDrawVaultGas, "DepositDrawVaultGas")
	return &types.MsgDepositAndDrawResponse{}, nil
}

func (k msgServer) MsgCreateStableMint(c context.Context, msg *types.MsgCreateStableMintRequest) (*types.MsgCreateStableMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	esmStatus, found := k.esm.GetESMStatus(ctx, msg.AppId)
	status := false
	if found {
		status = esmStatus.Status
	}
	if status {
		return nil, esmtypes.ErrESMAlreadyExecuted
	}
	killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, msg.AppId)
	if killSwitchParams.BreakerEnable {
		return nil, esmtypes.ErrCircuitBreakerEnabled
	}
	// Checking if extended pair exists
	extendedPairVault, found := k.asset.GetPairsVault(ctx, msg.ExtendedPairVaultId)
	if !found {
		return nil, types.ErrorExtendedPairVaultDoesNotExists
	}
	pairData, found := k.asset.GetPair(ctx, extendedPairVault.PairId)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}
	assetInData, found := k.asset.GetAsset(ctx, pairData.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	assetOutData, found := k.asset.GetAsset(ctx, pairData.AssetOut)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	// Checking if appMapping_id exists
	appMapping, found := k.asset.GetApp(ctx, msg.AppId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}

	// Checking if the appMapping_id in the msg_create & extendedPairVault_are same or not
	if appMapping.Id != extendedPairVault.AppId {
		return nil, types.ErrorAppMappingIDMismatch
	}

	// Converting user address for bank transaction
	depositorAddress, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	// Checking if this is a stableMint pair or not  -- stableMintPair == psmPair
	if !extendedPairVault.IsStableMintVault {
		return nil, types.ErrorCannotCreateStableMintVault
	}
	// Checking
	if !extendedPairVault.IsVaultActive {
		return nil, types.ErrorVaultCreationInactive
	}
	// Call CheckAppExtendedPairVaultMapping function to get counter - it also initialised the kv store if appMapping_id does not exists, or extendedPairVault_id does not exists.

	_, tokenOutAmount, err := k.GetAmountOfOtherToken(ctx, assetInData.Id, sdk.OneDec(), msg.Amount, assetOutData.Id, sdk.OneDec())
	if err != nil {
		return nil, err
	}
	// Check debt Floor
	if !tokenOutAmount.GTE(extendedPairVault.DebtFloor) {
		return nil, types.ErrorAmountOutLessThanDebtFloor
	}

	tokenMintedStatistics, _ := k.CheckAppExtendedPairVaultMapping(ctx, appMapping.Id, extendedPairVault.Id)

	extPairData, _ := k.GetAppExtendedPairVaultMappingData(ctx, appMapping.Id, msg.ExtendedPairVaultId)
	if len(extPairData.VaultIds) >= 1 {
		return nil, types.ErrorStableMintVaultAlreadyCreated
	}

	// Check Debt Ceil
	// currentMintedStatistics := tokenMintedStatistics.Add(msg.Amount)
	currentMintedStatistics := tokenMintedStatistics.Add(tokenOutAmount)

	if currentMintedStatistics.GTE(extendedPairVault.DebtCeiling) {
		return nil, types.ErrorAmountOutGreaterThanDebtCeiling
	}
	var amountToUser sdk.Int

	if msg.Amount.GT(sdk.ZeroInt()) {
		// Take amount from user
		if err := k.bank.SendCoinsFromAccountToModule(ctx, depositorAddress, types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetInData.Denom, msg.Amount))); err != nil {
			return nil, err
		}
		// Mint Tokens for user
		// mintCoin := sdk.NewCoin(assetOutData.Denom, msg.Amount)
		mintCoin := sdk.NewCoin(assetOutData.Denom, tokenOutAmount)
		if mintCoin.IsZero() {
			return nil, types.MintCoinValueInVaultIsZero
		}
		if err := k.bank.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mintCoin)); err != nil {
			return nil, err
		}
	}

	if extendedPairVault.DrawDownFee.IsZero() && msg.Amount.GT(sdk.ZeroInt()) {
		// Send Rest to user
		if err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositorAddress, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, msg.Amount))); err != nil {
			return nil, err
		}
	} else {
		// If not zero deduct send to collector//////////
		//			COLLECTOR FUNCTION
		// collectorShare := (tokenOutAmount.Mul(sdk.Int(extendedPairVault.DrawDownFee))).Quo(sdk.Int(sdk.OneDec()))
		collectorShare := sdk.NewDecFromInt(tokenOutAmount).Mul(extendedPairVault.DrawDownFee).TruncateInt()

		if collectorShare.GT(sdk.ZeroInt()) {
			if err := k.bank.SendCoinsFromModuleToModule(ctx, types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, collectorShare))); err != nil {
				return nil, err
			}
			err := k.collector.UpdateCollector(ctx, appMapping.Id, pairData.AssetOut, sdk.ZeroInt(), sdk.ZeroInt(), collectorShare, sdk.ZeroInt())
			if err != nil {
				return nil, err
			}
		}

		// and send the rest to the user
		// amountToUser := msg.Amount.Sub(collectorShare)
		amountToUser = tokenOutAmount.Sub(collectorShare)
		if amountToUser.GT(sdk.ZeroInt()) {
			if err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositorAddress, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, amountToUser))); err != nil {
				return nil, err
			}
		}
	}
	// Create Mint Vault

	oldID := k.GetIDForStableVault(ctx)
	var stableVault types.StableMintVault
	newID := oldID + 1

	stableVault.Id = newID
	stableVault.AmountIn = msg.Amount
	// stableVault.AmountOut = msg.Amount
	stableVault.AmountOut = tokenOutAmount
	stableVault.AppId = appMapping.Id
	stableVault.CreatedAt = ctx.BlockTime()
	stableVault.ExtendedPairVaultID = extendedPairVault.Id
	k.SetStableMintVault(ctx, stableVault)
	k.SetIDForStableVault(ctx, newID)
	// update Locker Data 	//Update Amount
	k.UpdateAppExtendedPairVaultMappingDataOnMsgCreateStableMintVault(ctx, stableVault)
	found = k.rewards.VerifyAppIDInRewards(ctx, msg.AppId)
	if found {
		var stableRewards types.StableMintVaultRewards
		stableRewards.AppId = msg.AppId
		stableRewards.StableExtendedPairId = msg.ExtendedPairVaultId
		stableRewards.User = msg.From
		stableRewards.BlockHeight = uint64(ctx.BlockHeight())
		stableRewards.Amount = amountToUser
		k.SetStableMintVaultRewards(ctx, stableRewards)
	}

	ctx.GasMeter().ConsumeGas(types.CreateStableVaultGas, "CreateStableVaultGas")

	return &types.MsgCreateStableMintResponse{}, nil
}

func (k msgServer) MsgDepositStableMint(c context.Context, msg *types.MsgDepositStableMintRequest) (*types.MsgDepositStableMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	esmStatus, found := k.esm.GetESMStatus(ctx, msg.AppId)
	status := false
	if found {
		status = esmStatus.Status
	}
	if status {
		return nil, esmtypes.ErrESMAlreadyExecuted
	}
	killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, msg.AppId)
	if killSwitchParams.BreakerEnable {
		return nil, esmtypes.ErrCircuitBreakerEnabled
	}
	depositorAddress, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	// checks if extended pair exists
	extendedPairVault, found := k.asset.GetPairsVault(ctx, msg.ExtendedPairVaultId)
	if !found {
		return nil, types.ErrorExtendedPairVaultDoesNotExists
	}
	pairData, found := k.asset.GetPair(ctx, extendedPairVault.PairId)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}
	assetInData, found := k.asset.GetAsset(ctx, pairData.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	assetOutData, found := k.asset.GetAsset(ctx, pairData.AssetOut)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	// Checking if appMapping_id exists
	appMapping, found := k.asset.GetApp(ctx, msg.AppId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	// Checking if vault access disabled
	if !extendedPairVault.IsVaultActive {
		return nil, types.ErrorVaultInactive
	}
	if !extendedPairVault.IsStableMintVault {
		return nil, types.ErrorCannotCreateStableMintVault
	}
	// Checking if the appMapping_id in the msg_create & extendedPairVault_are same or not
	if appMapping.Id != extendedPairVault.AppId {
		return nil, types.ErrorAppMappingIDMismatch
	}

	stableVault, found := k.GetStableMintVault(ctx, msg.StableVaultId)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	if appMapping.Id != stableVault.AppId {
		return nil, types.ErrorInvalidAppMappingData
	}
	if extendedPairVault.Id != stableVault.ExtendedPairVaultID {
		return nil, types.ErrorInvalidExtendedPairMappingData
	}

	stableAmountIn := stableVault.AmountIn.Add(msg.Amount)
	if !stableAmountIn.IsPositive() {
		return nil, types.ErrorInvalidAmount
	}
	// Looking for a case where apart from create function , this function creates new vaults and its data.
	tokenMintedStatistics, _ := k.CheckAppExtendedPairVaultMapping(ctx, appMapping.Id, extendedPairVault.Id)

	_, tokenOutAmount, err := k.GetAmountOfOtherToken(ctx, assetInData.Id, sdk.OneDec(), msg.Amount, assetOutData.Id, sdk.OneDec())
	if err != nil {
		return nil, err
	}
	// Check debt Floor
	if !tokenOutAmount.GTE(extendedPairVault.DebtFloor) {
		return nil, types.ErrorAmountOutLessThanDebtFloor
	}
	// Check Debt Ceil
	// currentMintedStatistics := tokenMintedStatistics.Add(msg.Amount)
	currentMintedStatistics := tokenMintedStatistics.Add(tokenOutAmount)

	if currentMintedStatistics.GTE(extendedPairVault.DebtCeiling) {
		return nil, types.ErrorAmountOutGreaterThanDebtCeiling
	}
	var amountToUser sdk.Int

	if msg.Amount.GT(sdk.ZeroInt()) {
		// Take amount from user
		if err := k.bank.SendCoinsFromAccountToModule(ctx, depositorAddress, types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetInData.Denom, msg.Amount))); err != nil {
			return nil, err
		}
		// Mint Tokens for user
		// mintCoin := sdk.NewCoin(assetOutData.Denom, msg.Amount)
		mintCoin := sdk.NewCoin(assetOutData.Denom, tokenOutAmount)
		if mintCoin.IsZero() {
			return nil, types.MintCoinValueInVaultIsZero
		}
		if err := k.bank.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mintCoin)); err != nil {
			return nil, err
		}
	}
	if extendedPairVault.DrawDownFee.IsZero() && msg.Amount.GT(sdk.ZeroInt()) {
		// Send Rest to user
		if err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositorAddress, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, tokenOutAmount))); err != nil {
			return nil, err
		}
	} else {
		//If not zero deduct send to collector//////////
		//			COLLECTOR FUNCTION
		/////////////////////////////////////////////////

		collectorShare := sdk.NewDecFromInt(tokenOutAmount).Mul(extendedPairVault.DrawDownFee).TruncateInt()

		if collectorShare.GT(sdk.ZeroInt()) {
			if err := k.bank.SendCoinsFromModuleToModule(ctx, types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, collectorShare))); err != nil {
				return nil, err
			}
			err := k.collector.UpdateCollector(ctx, appMapping.Id, pairData.AssetOut, sdk.ZeroInt(), sdk.ZeroInt(), collectorShare, sdk.ZeroInt())
			if err != nil {
				return nil, err
			}
		}

		// and send the rest to the user
		// amountToUser := msg.Amount.Sub(collectorShare)
		amountToUser = tokenOutAmount.Sub(collectorShare)
		if amountToUser.GT(sdk.ZeroInt()) {
			if err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositorAddress, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, amountToUser))); err != nil {
				return nil, err
			}
		}
	}
	stableVault.AmountIn = stableVault.AmountIn.Add(msg.Amount)
	// stableVault.AmountOut = stableVault.AmountOut.Add(msg.Amount)
	stableVault.AmountOut = stableVault.AmountOut.Add(tokenOutAmount)

	k.SetStableMintVault(ctx, stableVault)
	appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMappingData(ctx, appMapping.Id, msg.ExtendedPairVaultId)
	k.UpdateCollateralLockedAmountLockerMapping(ctx, appExtendedPairVaultData.AppId, appExtendedPairVaultData.ExtendedPairId, msg.Amount, true)
	k.UpdateTokenMintedAmountLockerMapping(ctx, appExtendedPairVaultData.AppId, appExtendedPairVaultData.ExtendedPairId, tokenOutAmount, true)
	found = k.rewards.VerifyAppIDInRewards(ctx, msg.AppId)
	if found {
		var stableRewards types.StableMintVaultRewards
		stableRewards.AppId = msg.AppId
		stableRewards.StableExtendedPairId = msg.ExtendedPairVaultId
		stableRewards.User = msg.From
		stableRewards.BlockHeight = uint64(ctx.BlockHeight())
		stableRewards.Amount = amountToUser
		k.SetStableMintVaultRewards(ctx, stableRewards)
	}

	ctx.GasMeter().ConsumeGas(types.DepositStableVaultGas, "DepositStableVaultGas")
	return &types.MsgDepositStableMintResponse{}, nil
}

func (k msgServer) MsgWithdrawStableMint(c context.Context, msg *types.MsgWithdrawStableMintRequest) (*types.MsgWithdrawStableMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	getControl := k.GetWithdrawStableMintControl(ctx)
	if getControl {
		return nil, types.ErrorWithdrawStableMintVault
	}
	esmStatus, found := k.esm.GetESMStatus(ctx, msg.AppId)
	status := false
	if found {
		status = esmStatus.Status
	}
	if status {
		return nil, esmtypes.ErrESMAlreadyExecuted
	}
	killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, msg.AppId)
	if killSwitchParams.BreakerEnable {
		return nil, esmtypes.ErrCircuitBreakerEnabled
	}
	depositorAddress, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	// checks if extended pair exists
	extendedPairVault, found := k.asset.GetPairsVault(ctx, msg.ExtendedPairVaultId)
	if !found {
		return nil, types.ErrorExtendedPairVaultDoesNotExists
	}
	pairData, found := k.asset.GetPair(ctx, extendedPairVault.PairId)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}
	assetInData, found := k.asset.GetAsset(ctx, pairData.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	assetOutData, found := k.asset.GetAsset(ctx, pairData.AssetOut)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	// Checking if appMapping_id exists
	appMapping, found := k.asset.GetApp(ctx, msg.AppId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	// Checking if vault access disabled

	if !extendedPairVault.IsStableMintVault {
		return nil, types.ErrorCannotCreateStableMintVault
	}
	// Checking if the appMapping_id in the msg_create & extendedPairVault_are same or not
	if appMapping.Id != extendedPairVault.AppId {
		return nil, types.ErrorAppMappingIDMismatch
	}

	// Check debt Floor
	if !msg.Amount.GTE(extendedPairVault.DebtFloor) {
		return nil, types.ErrorAmountOutLessThanDebtFloor
	}

	stableVault, found := k.GetStableMintVault(ctx, msg.StableVaultId)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	if appMapping.Id != stableVault.AppId {
		return nil, types.ErrorInvalidAppMappingData
	}
	if extendedPairVault.Id != stableVault.ExtendedPairVaultID {
		return nil, types.ErrorInvalidExtendedPairMappingData
	}

	_, tokenOutAmount, err := k.GetAmountOfOtherToken(ctx, assetOutData.Id, sdk.OneDec(), msg.Amount, assetInData.Id, sdk.OneDec())
	if err != nil {
		return nil, nil
	}

	// stableAmountIn := stableVault.AmountIn.Sub(msg.Amount)
	stableAmountIn := stableVault.AmountIn.Sub(tokenOutAmount)
	if stableAmountIn.LT(sdk.NewInt(0)) {
		return nil, types.ErrorInvalidAmount
	}
	// updated amount is the CMST amount
	var updatedAmount sdk.Int
	// Take amount from user
	if msg.Amount.GT(sdk.ZeroInt()) {
		if err := k.bank.SendCoinsFromAccountToModule(ctx, depositorAddress, types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, msg.Amount))); err != nil {
			return nil, err
		}
	}

	if extendedPairVault.DrawDownFee.IsZero() && msg.Amount.GT(sdk.ZeroInt()) {
		// BurnTokens for user
		burnCoin := sdk.NewCoin(assetOutData.Denom, msg.Amount)
		if burnCoin.IsZero() {
			return nil, types.BurnCoinValueInVaultIsZero
		}
		if err := k.bank.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(burnCoin)); err != nil {
			return nil, err
		}

		// Send Rest to user
		if err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositorAddress, sdk.NewCoins(sdk.NewCoin(assetInData.Denom, tokenOutAmount))); err != nil {
			return nil, err
		}
		updatedAmount = msg.Amount
	} else {
		//If not zero deduct send to collector//////////
		//			COLLECTOR FUNCTION
		/////////////////////////////////////////////////
		//collectorShare := (msg.Amount.Mul(sdk.Int(extendedPairVault.DrawDownFee))).Quo(sdk.Int(sdk.OneDec()))
		collectorShare := sdk.NewDecFromInt(msg.Amount).Mul(extendedPairVault.DrawDownFee).TruncateInt()

		if collectorShare.GT(sdk.ZeroInt()) {
			if err := k.bank.SendCoinsFromModuleToModule(ctx, types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, collectorShare))); err != nil {
				return nil, err
			}
			err := k.collector.UpdateCollector(ctx, appMapping.Id, pairData.AssetOut, sdk.ZeroInt(), sdk.ZeroInt(), collectorShare, sdk.ZeroInt())
			if err != nil {
				return nil, err
			}
		}

		updatedAmount = msg.Amount.Sub(collectorShare)

		if updatedAmount.GT(sdk.ZeroInt()) {
			// BurnTokens for user
			burnCoin := sdk.NewCoin(assetOutData.Denom, updatedAmount)
			if burnCoin.IsZero() {
				return nil, types.BurnCoinValueInVaultIsZero
			}
			if err := k.bank.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(burnCoin)); err != nil {
				return nil, err
			}

			// and send the rest to the user
			_, newOutAmount, err := k.GetAmountOfOtherToken(ctx, assetOutData.Id, sdk.OneDec(), updatedAmount, assetInData.Id, sdk.OneDec())
			if err != nil {
				return nil, err
			}
			tokenOutAmount = newOutAmount

			if err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositorAddress, sdk.NewCoins(sdk.NewCoin(assetInData.Denom, newOutAmount))); err != nil {
				return nil, err
			}
		}
	}
	// stableVault.AmountIn = stableVault.AmountIn.Sub(updatedAmount)
	stableVault.AmountIn = stableVault.AmountIn.Sub(tokenOutAmount)
	stableVault.AmountOut = stableVault.AmountOut.Sub(updatedAmount)
	k.SetStableMintVault(ctx, stableVault)
	appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMappingData(ctx, appMapping.Id, msg.ExtendedPairVaultId)
	k.UpdateCollateralLockedAmountLockerMapping(ctx, appExtendedPairVaultData.AppId, appExtendedPairVaultData.ExtendedPairId, tokenOutAmount, false)
	k.UpdateTokenMintedAmountLockerMapping(ctx, appExtendedPairVaultData.AppId, appExtendedPairVaultData.ExtendedPairId, updatedAmount, false)

	// Function that deletes the entries in the stable mint rewards structure.
	k.DeleteUserStableRewardEntries(ctx, appExtendedPairVaultData.AppId, msg.From, updatedAmount)

	ctx.GasMeter().ConsumeGas(types.WithdrawStableVaultGas, "WithdrawStableVaultGas")

	return &types.MsgWithdrawStableMintResponse{}, nil
}

// take app id
// check app id
// take vault id
// check vault id
// calculate total debt
// call function
// exit function

func (k msgServer) MsgVaultInterestCalc(c context.Context, msg *types.MsgVaultInterestCalcRequest) (*types.MsgVaultInterestCalcResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	appMapping, found := k.asset.GetApp(ctx, msg.AppId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	userVault, found := k.GetVault(ctx, msg.UserVaultId)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}

	totalDebt := userVault.AmountOut.Add(userVault.InterestAccumulated)
	err1 := k.rewards.CalculateVaultInterest(ctx, appMapping.Id, userVault.ExtendedPairVaultID, msg.UserVaultId, totalDebt, userVault.BlockHeight, userVault.BlockTime.Unix())
	if err1 != nil {
		return nil, err1
	}

	return &types.MsgVaultInterestCalcResponse{}, nil
}

func (k msgServer) MsgWithdrawStableMintControl(c context.Context, msg *types.MsgWithdrawStableMintControlRequest) (*types.MsgWithdrawStableMintControlResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// check if address is admin
	getAdmin := k.esm.GetParams(ctx).Admin

	// check if address is admin in getAdmin array
	if getAdmin[0] != msg.From {
		return nil, esmtypes.ErrorUnauthorized
	}

	appMapping, found := k.asset.GetApp(ctx, msg.AppId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	// check app name is harbor
	if appMapping.Name != "harbor" {
		return nil, types.ErrorAppMappingDoesNotExist
	}

	// check GetWithdrawStableMintControl value
	control := k.GetWithdrawStableMintControl(ctx)

	if control {
		k.SetWithdrawStableMintControl(ctx, false)
	} else {
		k.SetWithdrawStableMintControl(ctx, true)
	}

	return &types.MsgWithdrawStableMintControlResponse{}, nil
}

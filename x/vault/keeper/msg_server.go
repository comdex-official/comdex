package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/types/time"

	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	"github.com/comdex-official/comdex/x/vault/types"
)

var (
	_ types.MsgServer = (*msgServer)(nil)
)

type msgServer struct {
	Keeper
}

func NewMsgServer(keeper Keeper) types.MsgServer {
	return &msgServer{
		Keeper: keeper,
	}
}

// MsgCreate Creating a new CDP.
func (k *msgServer) MsgCreate(c context.Context, msg *types.MsgCreateRequest) (*types.MsgCreateResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	extendedPairVault, found := k.GetPairsVault(ctx, msg.ExtendedPairVaultId)
	if !found {
		return nil, types.ErrorExtendedPairVaultDoesNotExists
	}
	pairData, found := k.GetPair(ctx, extendedPairVault.PairId)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}
	assetInData, found := k.GetAsset(ctx, pairData.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	assetOutData, found := k.GetAsset(ctx, pairData.AssetOut)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	appMapping, found := k.GetApp(ctx, msg.AppMappingId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}

	if appMapping.Id != extendedPairVault.AppMappingId {
		return nil, types.ErrorAppMappingIDMismatch
	}

	depositorAddress, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	// Checking if this is a stableMint pair or not  -- stableMintPair == psmPair
	if extendedPairVault.IsPsmPair {
		return nil, types.ErrorCannotCreateStableMintVault
	}
	//Checking
	if !extendedPairVault.IsVaultActive {
		return nil, types.ErrorVaultCreationInactive

	}
	//if does then check app to extendedPair mapping has any vault key
	//if it does throw error
	userVaultExtendedPairMapping, userExists := k.GetUserVaultExtendedPairMapping(ctx, msg.From)
	if userExists {
		_, alreadyExists := k.CheckUserAppToExtendedPairMapping(ctx, userVaultExtendedPairMapping, extendedPairVault.Id, appMapping.Id)
		if alreadyExists {
			return nil, types.ErrorUserVaultAlreadyExists

		}
	}
	//Call CheckAppExtendedPairVaultMapping function to get counter - it also initialised the kv store if appMapping_id does not exists, or extendedPairVault_id does not exists.

	counterVal, tokenMintedStatistics, _ := k.CheckAppExtendedPairVaultMapping(ctx, appMapping.Id, extendedPairVault.Id)

	//Check debt Floor
	if !msg.AmountOut.GTE(extendedPairVault.DebtFloor) {
		return nil, types.ErrorAmountOutLessThanDebtFloor
	}
	//Check Debt Ceil
	currentMintedStatistics := tokenMintedStatistics.Add(msg.AmountOut)

	if currentMintedStatistics.GT(extendedPairVault.DebtCeiling) {
		return nil, types.ErrorAmountOutGreaterThanDebtCeiling
	}

	//Calculate CR - make necessary changes to calculate collateralization function
	if err := k.VerifyCollaterlizationRatio(ctx, extendedPairVault.Id, msg.AmountIn, msg.AmountOut, extendedPairVault.MinCr); err != nil {
		return nil, err
	}
	//Take amount from user
	if err := k.SendCoinFromAccountToModule(ctx, depositorAddress, types.ModuleName, sdk.NewCoin(assetInData.Denom, msg.AmountIn)); err != nil {
		return nil, err
	}
	//Mint Tokens for user
	if err := k.MintCoin(ctx, types.ModuleName, sdk.NewCoin(assetOutData.Denom, msg.AmountOut)); err != nil {
		return nil, err
	}

	//Send Fees to Accumulator
	//Deducting Opening Fee if 0 opening fee then act accordingly
	if extendedPairVault.DrawDownFee.IsZero() { //Send Rest to user
		if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, depositorAddress, sdk.NewCoin(assetOutData.Denom, msg.AmountOut)); err != nil {
			return nil, err
		}
	} else {
		//If not zero deduct send to collector//////////
		//one approach could be
		collectorShare := (msg.AmountOut.Mul(sdk.Int(extendedPairVault.DrawDownFee))).Quo(sdk.Int(sdk.OneDec()))

		if err := k.SendCoinFromModuleToModule(ctx, types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, collectorShare))); err != nil {
			return nil, err
		}
		err := k.UpdateCollector(ctx, appMapping.Id, pairData.AssetOut, sdk.ZeroInt(), sdk.ZeroInt(), collectorShare, sdk.ZeroInt())
		if err != nil {
			return nil, err
		}

		// and send the rest to the user
		amountToUser := msg.AmountOut.Sub(collectorShare)
		if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, depositorAddress, sdk.NewCoin(assetOutData.Denom, amountToUser)); err != nil {
			return nil, err
		}
	}

	//If all correct  create vault
	zeroVal := sdk.ZeroInt()
	var newVault types.Vault
	updatedCounter := counterVal + 1
	newVault.Id = appMapping.ShortName + strconv.FormatUint(updatedCounter, 10)
	newVault.AmountIn = msg.AmountIn

	// closingFeeVal := (sdk.Dec(msg.AmountOut).Mul((extendedPairVault.ClosingFee)))

	closingFeeVal := msg.AmountOut.Mul(sdk.Int(extendedPairVault.ClosingFee)).Quo(sdk.Int(sdk.OneDec()))

	newVault.ClosingFeeAccumulated = closingFeeVal
	newVault.AmountOut = msg.AmountOut
	newVault.AppMappingId = appMapping.Id
	newVault.InterestAccumulated = zeroVal
	newVault.Owner = msg.From
	newVault.CreatedAt = time.Now()
	newVault.ExtendedPairVaultID = extendedPairVault.Id

	k.SetVault(ctx, newVault)

	//Update mapping data - take proper approach
	// lookup table already exists
	//only need to update counter and token statistics value
	k.UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx, updatedCounter, newVault)

	userVaultExtendedPairMappingData, userExists := k.GetUserVaultExtendedPairMapping(ctx, msg.From)
	if !userExists {
		var userMappingData types.UserVaultAssetMapping
		var userAppData types.VaultToAppMapping
		var userExtendedPairData types.ExtendedPairToVaultMapping

		userExtendedPairData.ExtendedPairId = newVault.ExtendedPairVaultID
		userExtendedPairData.VaultId = newVault.Id
		userAppData.AppMappingId = appMapping.Id
		userAppData.UserExtendedPairVault = append(userAppData.UserExtendedPairVault, &userExtendedPairData)
		userMappingData.Owner = msg.From
		userMappingData.UserVaultApp = append(userMappingData.UserVaultApp, &userAppData)

		k.SetUserVaultExtendedPairMapping(ctx, userMappingData)
	} else {
		///Check if user appMapping data exits

		appExists := k.CheckUserToAppMapping(ctx, userVaultExtendedPairMappingData, appMapping.Id)
		if appExists {
			//User has the appMapping added
			//So only need to add the locker id with asset
			var userExtendedPairData types.ExtendedPairToVaultMapping
			userExtendedPairData.VaultId = newVault.Id
			userExtendedPairData.ExtendedPairId = newVault.ExtendedPairVaultID

			for _, appData := range userVaultExtendedPairMappingData.UserVaultApp {
				if appData.AppMappingId == appMapping.Id {
					appData.UserExtendedPairVault = append(appData.UserExtendedPairVault, &userExtendedPairData)
				}
			}
			k.SetUserVaultExtendedPairMapping(ctx, userVaultExtendedPairMappingData)
		} else {
			var userAppData types.VaultToAppMapping
			var userExtendedPairData types.ExtendedPairToVaultMapping

			userExtendedPairData.ExtendedPairId = newVault.ExtendedPairVaultID
			userExtendedPairData.VaultId = newVault.Id
			userAppData.AppMappingId = appMapping.Id
			userAppData.UserExtendedPairVault = append(userAppData.UserExtendedPairVault, &userExtendedPairData)
			userVaultExtendedPairMappingData.UserVaultApp = append(userVaultExtendedPairMappingData.UserVaultApp, &userAppData)
			k.SetUserVaultExtendedPairMapping(ctx, userVaultExtendedPairMappingData)
		}
	}

	return &types.MsgCreateResponse{}, nil
}

// MsgDeposit Only for depositing new collateral.
func (k *msgServer) MsgDeposit(c context.Context, msg *types.MsgDepositRequest) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	depositor, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	//checks if extended pair exists
	extendedPairVault, found := k.GetPairsVault(ctx, msg.ExtendedPairVaultId)
	if !found {
		return nil, types.ErrorExtendedPairVaultDoesNotExists
	}
	pairData, found := k.GetPair(ctx, extendedPairVault.PairId)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}
	assetInData, found := k.GetAsset(ctx, pairData.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	//Checking if appMapping_id exists
	appMapping, found := k.GetApp(ctx, msg.AppMappingId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	//Checking if vault access disabled
	if !extendedPairVault.IsVaultActive {
		return nil, types.ErrorVaultInactive
	}

	//Checking if the appMapping_id in the msg_create & extendedPairVault_are same or not
	if appMapping.Id != extendedPairVault.AppMappingId {
		return nil, types.ErrorAppMappingIDMismatch
	}

	userVault, found := k.GetVault(ctx, msg.UserVaultId)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	if userVault.Owner != msg.From {
		return nil, types.ErrVaultAccessUnauthorised
	}

	if appMapping.Id != userVault.AppMappingId {
		return nil, types.ErrorInvalidAppMappingData
	}
	if extendedPairVault.Id != userVault.ExtendedPairVaultID {
		return nil, types.ErrorInvalidExtendedPairMappingData
	}

	userVault.AmountIn = userVault.AmountIn.Add(msg.Amount)
	if !userVault.AmountIn.IsPositive() {
		return nil, types.ErrorInvalidAmount
	}

	if err := k.SendCoinFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoin(assetInData.Denom, msg.Amount)); err != nil {
		return nil, err
	}

	k.SetVault(ctx, userVault)
	//Updating appExtendedPairvaultMappingData data -
	appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMapping(ctx, appMapping.Id)
	k.UpdateCollateralLockedAmountLockerMapping(ctx, appExtendedPairVaultData, extendedPairVault.Id, msg.Amount, true)

	return &types.MsgDepositResponse{}, nil
}

// MsgWithdraw Withdrawing collateral.
func (k *msgServer) MsgWithdraw(c context.Context, msg *types.MsgWithdrawRequest) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	depositor, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	//checks if extended pair exists
	extendedPairVault, found := k.GetPairsVault(ctx, msg.ExtendedPairVaultId)
	if !found {
		return nil, types.ErrorExtendedPairVaultDoesNotExists
	}
	pairData, found := k.GetPair(ctx, extendedPairVault.PairId)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}
	assetInData, found := k.GetAsset(ctx, pairData.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	// assetOutData, found := k.GetAsset(ctx, pairData.AssetOut)
	// if !found {
	// 	return nil, types.ErrorAssetDoesNotExist
	// }

	//Checking if appMapping_id exists
	appMapping, found := k.GetApp(ctx, msg.AppMappingId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	//Checking if vault access disabled
	if !extendedPairVault.IsVaultActive {
		return nil, types.ErrorVaultInactive
	}
	//Checking if the appMapping_id in the msg_create & extendedPairVault_are same or not
	if appMapping.Id != extendedPairVault.AppMappingId {
		return nil, types.ErrorAppMappingIDMismatch
	}

	userVault, found := k.GetVault(ctx, msg.UserVaultId)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	if userVault.Owner != msg.From {
		return nil, types.ErrVaultAccessUnauthorised
	}

	if appMapping.Id != userVault.AppMappingId {
		return nil, types.ErrorInvalidAppMappingData
	}
	if extendedPairVault.Id != userVault.ExtendedPairVaultID {
		return nil, types.ErrorInvalidExtendedPairMappingData
	}

	userVault.AmountIn = userVault.AmountIn.Sub(msg.Amount)
	if !userVault.AmountIn.IsPositive() {
		return nil, types.ErrorInvalidAmount
	}

	totalDebtCalculation := userVault.AmountOut.Add(userVault.InterestAccumulated)
	totalDebtCalculation = totalDebtCalculation.Add(userVault.ClosingFeeAccumulated)

	//Calculate CR - make necessary changes to the calculate collateralization function
	if err := k.VerifyCollaterlizationRatio(ctx, extendedPairVault.Id, userVault.AmountIn, totalDebtCalculation, extendedPairVault.MinCr); err != nil {
		return nil, err
	}

	if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, depositor, sdk.NewCoin(assetInData.Denom, msg.Amount)); err != nil {
		return nil, err
	}

	k.SetVault(ctx, userVault)

	//Updating appExtendedPairVaultMappingData
	appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMapping(ctx, appMapping.Id)
	k.UpdateCollateralLockedAmountLockerMapping(ctx, appExtendedPairVaultData, extendedPairVault.Id, msg.Amount, false)

	return &types.MsgWithdrawResponse{}, nil
}

// MsgDraw To borrow more amount.
func (k *msgServer) MsgDraw(c context.Context, msg *types.MsgDrawRequest) (*types.MsgDrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	depositor, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	//checks if extended pair exists
	extendedPairVault, found := k.GetPairsVault(ctx, msg.ExtendedPairVaultId)
	if !found {
		return nil, types.ErrorExtendedPairVaultDoesNotExists
	}
	pairData, found := k.GetPair(ctx, extendedPairVault.PairId)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}
	// assetInData, found := k.GetAsset(ctx, pairData.AssetIn)
	// if !found {
	// 	return nil, types.ErrorAssetDoesNotExist
	// }
	assetOutData, found := k.GetAsset(ctx, pairData.AssetOut)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	//Checking if appMapping_id exists
	appMapping, found := k.GetApp(ctx, msg.AppMappingId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	//Checking if vault access disabled
	if !extendedPairVault.IsVaultActive {
		return nil, types.ErrorVaultInactive
	}
	//Checking if the appMapping_id in the msg_create & extendedPairVault_are same or not
	if appMapping.Id != extendedPairVault.AppMappingId {
		return nil, types.ErrorAppMappingIDMismatch
	}

	userVault, found := k.GetVault(ctx, msg.UserVaultId)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	if userVault.Owner != msg.From {
		return nil, types.ErrVaultAccessUnauthorised
	}

	if appMapping.Id != userVault.AppMappingId {
		return nil, types.ErrorInvalidAppMappingData
	}
	if extendedPairVault.Id != userVault.ExtendedPairVaultID {
		return nil, types.ErrorInvalidExtendedPairMappingData
	}
	if msg.Amount.LTE(sdk.NewInt(0)) {
		return nil, types.ErrorInvalidAmount
	}

	newUpdatedAmountOut := userVault.AmountOut.Add(msg.Amount)
	totalDebt := newUpdatedAmountOut.Add(userVault.InterestAccumulated)
	totalDebt = totalDebt.Add(userVault.ClosingFeeAccumulated)

	_, tokenMintedStatistics, _ := k.CheckAppExtendedPairVaultMapping(ctx, appMapping.Id, extendedPairVault.Id)

	//Check Debt Ceil
	currentMintedStatistics := tokenMintedStatistics.Add(msg.Amount)

	if currentMintedStatistics.GTE(extendedPairVault.DebtCeiling) {
		return nil, types.ErrorAmountOutGreaterThanDebtCeiling
	}

	if err := k.VerifyCollaterlizationRatio(ctx, extendedPairVault.Id, userVault.AmountIn, totalDebt, extendedPairVault.MinCr); err != nil {
		return nil, err
	}

	if err := k.MintCoin(ctx, types.ModuleName, sdk.NewCoin(assetOutData.Denom, msg.Amount)); err != nil {
		return nil, err
	}

	if extendedPairVault.DrawDownFee.IsZero() {
		//Send Rest to user
		if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, depositor, sdk.NewCoin(assetOutData.Denom, msg.Amount)); err != nil {
			return nil, err
		}
	} else {
		//If not zero deduct send to collector//////////
		//one approach could be
		collectorShare := (msg.Amount.Mul(sdk.Int(extendedPairVault.DrawDownFee))).Quo(sdk.Int(sdk.OneDec()))

		if err := k.SendCoinFromModuleToModule(ctx, types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, collectorShare))); err != nil {
			return nil, err
		}
		err := k.UpdateCollector(ctx, appMapping.Id, pairData.AssetOut, sdk.ZeroInt(), sdk.ZeroInt(), collectorShare, sdk.ZeroInt())
		if err != nil {
			return nil, err
		}

		// and send the rest to the user
		amountToUser := msg.Amount.Sub(collectorShare)
		if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, depositor, sdk.NewCoin(assetOutData.Denom, amountToUser)); err != nil {
			return nil, err
		}
	}

	// if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, depositor, sdk.NewCoin(assetOutData.Denom, msg.Amount)); err != nil {
	// 	return nil, err
	// }
	userVault.AmountOut = userVault.AmountOut.Add(msg.Amount)

	k.SetVault(ctx, userVault)

	//Updating appExtendedPairVaultMappingData
	appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMapping(ctx, appMapping.Id)
	k.UpdateTokenMintedAmountLockerMapping(ctx, appExtendedPairVaultData, extendedPairVault.Id, msg.Amount, true)

	return &types.MsgDrawResponse{}, nil
}

func (k *msgServer) MsgRepay(c context.Context, msg *types.MsgRepayRequest) (*types.MsgRepayResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	depositor, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	//checks if extended pair exists
	extendedPairVault, found := k.GetPairsVault(ctx, msg.ExtendedPairVaultId)
	if !found {
		return nil, types.ErrorExtendedPairVaultDoesNotExists
	}
	pairData, found := k.GetPair(ctx, extendedPairVault.PairId)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}
	// assetInData, found := k.GetAsset(ctx, pairData.AssetIn)
	// if !found {
	// 	return nil, types.ErrorAssetDoesNotExist
	// }
	assetOutData, found := k.GetAsset(ctx, pairData.AssetOut)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	//Checking if appMapping_id exists
	appMapping, found := k.GetApp(ctx, msg.AppMappingId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	//Checking if vault acccess disabled
	if !extendedPairVault.IsVaultActive {
		return nil, types.ErrorVaultInactive
	}
	//Checking if the appMapping_id in the msg_create & extendedPairVault_are same or not
	if appMapping.Id != extendedPairVault.AppMappingId {
		return nil, types.ErrorAppMappingIDMismatch
	}

	userVault, found := k.GetVault(ctx, msg.UserVaultId)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	if userVault.Owner != msg.From {
		return nil, types.ErrVaultAccessUnauthorised
	}

	if appMapping.Id != userVault.AppMappingId {
		return nil, types.ErrorInvalidAppMappingData
	}
	if extendedPairVault.Id != userVault.ExtendedPairVaultID {
		return nil, types.ErrorInvalidExtendedPairMappingData
	}
	if msg.Amount.LTE(sdk.NewInt(0)) {
		return nil, types.ErrorInvalidAmount
	}

	newAmount := userVault.AmountOut.Add(userVault.InterestAccumulated)
	newAmount = newAmount.Sub(msg.Amount)
	if newAmount.LT(sdk.NewInt(0)) {
		return nil, types.ErrorInvalidAmount
	}

	if msg.Amount.LTE(userVault.InterestAccumulated) {
		//Amount is less than equal to the interest accumulated
		//subtract that as interest
		reducedFees := userVault.InterestAccumulated.Sub(msg.Amount)
		userVault.InterestAccumulated = reducedFees
		//and send it to the collector module
		if err := k.SendCoinFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoin(assetOutData.Denom, msg.Amount)); err != nil {
			return nil, err
		}
		//			SEND TO COLLECTOR- msg.Amount

		if err := k.SendCoinFromModuleToModule(ctx, types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, msg.Amount))); err != nil {
			return nil, err
		}
		err := k.UpdateCollector(ctx, appMapping.Id, pairData.AssetOut, msg.Amount, sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt())
		if err != nil {
			return nil, err
		}

		k.SetVault(ctx, userVault)
	} else {
		updatedUserSentAmountAfterFeesDeduction := msg.Amount.Sub(userVault.InterestAccumulated)

		updatedUserDebt := userVault.AmountOut.Sub(updatedUserSentAmountAfterFeesDeduction)

		// //If user's closing fees is a bigger amount than the debt floor, user will not close the debt floor
		// totalUpdatedDebt:=updatedUserDebt.Add(*userVault.ClosingFeeAccumulated)
		// if err := k.VerifyCollaterlizationRatio(ctx, extendedPairVault.Id, userVault.AmountIn, totalUpdatedDebt, extendedPairVault.MinCr); err != nil {
		// 	return nil, err
		// }

		if !updatedUserDebt.GTE(extendedPairVault.DebtFloor) {
			return nil, types.ErrorAmountOutLessThanDebtFloor
		}
		if err := k.SendCoinFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoin(assetOutData.Denom, msg.Amount)); err != nil {
			return nil, err
		}

		if err := k.BurnCoin(ctx, types.ModuleName, sdk.NewCoin(assetOutData.Denom, updatedUserSentAmountAfterFeesDeduction)); err != nil {
			return nil, err
		}
		//			SEND TO COLLECTOR----userVault.InterestAccumulated
		if err := k.SendCoinFromModuleToModule(ctx, types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, userVault.InterestAccumulated))); err != nil {
			return nil, err
		}
		err := k.UpdateCollector(ctx, appMapping.Id, pairData.AssetOut, userVault.InterestAccumulated, sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt())
		if err != nil {
			return nil, err
		}

		userVault.AmountOut = updatedUserDebt
		zeroval := sdk.ZeroInt()
		userVault.InterestAccumulated = zeroval
		k.SetVault(ctx, userVault)
		appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMapping(ctx, appMapping.Id)
		k.UpdateTokenMintedAmountLockerMapping(ctx, appExtendedPairVaultData, extendedPairVault.Id, updatedUserSentAmountAfterFeesDeduction, false)
	}

	return &types.MsgRepayResponse{}, nil
}

func (k *msgServer) MsgClose(c context.Context, msg *types.MsgCloseRequest) (*types.MsgCloseResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	depositor, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	//checks if extended pair exists
	extendedPairVault, found := k.GetPairsVault(ctx, msg.ExtendedPairVaultId)
	if !found {
		return nil, types.ErrorExtendedPairVaultDoesNotExists
	}
	pairData, found := k.GetPair(ctx, extendedPairVault.PairId)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}
	assetInData, found := k.GetAsset(ctx, pairData.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	assetOutData, found := k.GetAsset(ctx, pairData.AssetOut)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	//Checking if appMapping_id exists
	appMapping, found := k.GetApp(ctx, msg.AppMappingId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	// //Checking if vault acccess disabled
	// if !extendedPairVault.IsVaultActive {
	// 	return nil, types.ErrorVaultInactive

	// }

	//Checking if the appMapping_id in the msg_create & extendedPairVault_are same or not
	if appMapping.Id != extendedPairVault.AppMappingId {
		return nil, types.ErrorAppMappingIDMismatch
	}

	userVault, found := k.GetVault(ctx, msg.UserVaultId)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	if userVault.Owner != msg.From {
		return nil, types.ErrVaultAccessUnauthorised
	}

	if appMapping.Id != userVault.AppMappingId {
		return nil, types.ErrorInvalidAppMappingData
	}
	if extendedPairVault.Id != userVault.ExtendedPairVaultID {
		return nil, types.ErrorInvalidExtendedPairMappingData
	}

	totalUserDebt := userVault.AmountOut.Add(userVault.InterestAccumulated)
	totalUserDebt = totalUserDebt.Add(userVault.ClosingFeeAccumulated)
	if err := k.SendCoinFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoin(assetOutData.Denom, totalUserDebt)); err != nil {
		return nil, err
	}

	//			SEND TO COLLECTOR----userVault.InterestAccumulated & userVault.ClosingFees

	err = k.UpdateCollector(ctx, appMapping.Id, pairData.AssetOut, userVault.InterestAccumulated, userVault.ClosingFeeAccumulated, sdk.ZeroInt(), sdk.ZeroInt())
	if err != nil {
		return nil, err
	}
	if err := k.SendCoinFromModuleToModule(ctx, types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, userVault.InterestAccumulated))); err != nil {
		return nil, err
	}
	if err := k.SendCoinFromModuleToModule(ctx, types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, userVault.ClosingFeeAccumulated))); err != nil {
		return nil, err
	}
	if err := k.BurnCoin(ctx, types.ModuleName, sdk.NewCoin(assetOutData.Denom, userVault.AmountOut)); err != nil {
		return nil, err
	}

	if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, depositor, sdk.NewCoin(assetInData.Denom, userVault.AmountIn)); err != nil {
		return nil, err
	}

	//Update LookupTable minting Status
	appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMapping(ctx, appMapping.Id)

	k.UpdateCollateralLockedAmountLockerMapping(ctx, appExtendedPairVaultData, extendedPairVault.Id, userVault.AmountIn, false)
	k.UpdateTokenMintedAmountLockerMapping(ctx, appExtendedPairVaultData, extendedPairVault.Id, userVault.AmountOut, false)

	//Remove address from lookup table
	k.DeleteAddressFromAppExtendedPairVaultMapping(ctx, extendedPairVault.Id, userVault.Id, appMapping.Id)

	//Remove user extendedPair to address field in UserLookupStruct
	k.UpdateUserVaultExtendedPairMapping(ctx, extendedPairVault.Id, msg.From, appMapping.Id)

	//Delete Vault
	k.DeleteVault(ctx, userVault.Id)

	return &types.MsgCloseResponse{}, nil
}

func (k *msgServer) MsgCreateStableMint(c context.Context, msg *types.MsgCreateStableMintRequest) (*types.MsgCreateStableMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	//Checking if extended pair exists
	extendedPairVault, found := k.GetPairsVault(ctx, msg.ExtendedPairVaultId)
	if !found {
		return nil, types.ErrorExtendedPairVaultDoesNotExists
	}
	pairData, found := k.GetPair(ctx, extendedPairVault.PairId)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}
	assetInData, found := k.GetAsset(ctx, pairData.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	assetOutData, found := k.GetAsset(ctx, pairData.AssetOut)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	//Checking if appMapping_id exists
	appMapping, found := k.GetApp(ctx, msg.AppMappingId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}

	//Checking if the appMapping_id in the msg_create & extendedPairVault_are same or not
	if appMapping.Id != extendedPairVault.AppMappingId {
		return nil, types.ErrorAppMappingIDMismatch
	}

	//Converting user address for bank transaction
	depositorAddress, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	// Checking if this is a stableMint pair or not  -- stableMintPair == psmPair
	if !extendedPairVault.IsPsmPair {
		return nil, types.ErrorCannotCreateStableMintVault
	}
	//Checking
	if !extendedPairVault.IsVaultActive {
		return nil, types.ErrorVaultCreationInactive
	}
	//Call CheckAppExtendedPairVaultMapping function to get counter - it also initialised the kv store if appMapping_id does not exists, or extendedPairVault_id does not exists.

	counterVal, tokenMintedStatistics, lenOfVault := k.CheckAppExtendedPairVaultMapping(ctx, appMapping.Id, extendedPairVault.Id)

	if lenOfVault >= 1 {
		return nil, types.ErrorStableMintVaultAlreadyCreated
	}

	//Check Debt Ceil
	currentMintedStatistics := tokenMintedStatistics.Add(msg.Amount)

	if currentMintedStatistics.GTE(extendedPairVault.DebtCeiling) {
		return nil, types.ErrorAmountOutGreaterThanDebtCeiling
	}

	//Take amount from user
	if err := k.SendCoinFromAccountToModule(ctx, depositorAddress, types.ModuleName, sdk.NewCoin(assetInData.Denom, msg.Amount)); err != nil {
		return nil, err
	}
	//Mint Tokens for user

	if err := k.MintCoin(ctx, types.ModuleName, sdk.NewCoin(assetOutData.Denom, msg.Amount)); err != nil {
		return nil, err
	}
	if extendedPairVault.DrawDownFee.IsZero() {
		//Send Rest to user
		if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, depositorAddress, sdk.NewCoin(assetOutData.Denom, msg.Amount)); err != nil {
			return nil, err
		}
	} else {
		//If not zero deduct send to collector//////////
		//			COLLECTOR FUNCTION
		collectorShare := (msg.Amount.Mul(sdk.Int(extendedPairVault.DrawDownFee))).Quo(sdk.Int(sdk.OneDec()))
		if err := k.SendCoinFromModuleToModule(ctx, types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, collectorShare))); err != nil {
			return nil, err
		}
		err := k.UpdateCollector(ctx, appMapping.Id, pairData.AssetOut, sdk.ZeroInt(), sdk.ZeroInt(), collectorShare, sdk.ZeroInt())
		if err != nil {
			return nil, err
		}

		// and send the rest to the user
		amountToUser := msg.Amount.Sub(collectorShare)
		if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, depositorAddress, sdk.NewCoin(assetOutData.Denom, amountToUser)); err != nil {
			return nil, err
		}
	}
	//Create Mint Vault

	var stableVault types.StableMintVault
	updatedCounter := counterVal + 1

	stableVault.Id = appMapping.ShortName + strconv.FormatUint(updatedCounter, 10)
	stableVault.AmountIn = msg.Amount
	stableVault.AmountOut = msg.Amount
	stableVault.AppMappingId = appMapping.Id
	stableVault.CreatedAt = time.Now()
	stableVault.ExtendedPairVaultID = extendedPairVault.Id
	k.SetStableMintVault(ctx, stableVault)
	//update Locker Data 	//Update Amount
	k.UpdateAppExtendedPairVaultMappingDataOnMsgCreateStableMintVault(ctx, updatedCounter, stableVault)

	return &types.MsgCreateStableMintResponse{}, nil
}

func (k *msgServer) MsgDepositStableMint(c context.Context, msg *types.MsgDepositStableMintRequest) (*types.MsgDepositStableMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	depositorAddress, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	//checks if extended pair exists
	extendedPairVault, found := k.GetPairsVault(ctx, msg.ExtendedPairVaultId)
	if !found {
		return nil, types.ErrorExtendedPairVaultDoesNotExists
	}
	pairData, found := k.GetPair(ctx, extendedPairVault.PairId)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}
	assetInData, found := k.GetAsset(ctx, pairData.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	assetOutData, found := k.GetAsset(ctx, pairData.AssetOut)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	//Checking if appMapping_id exists
	appMapping, found := k.GetApp(ctx, msg.AppMappingId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	//Checking if vault access disabled
	if !extendedPairVault.IsVaultActive {
		return nil, types.ErrorVaultInactive
	}
	if !extendedPairVault.IsPsmPair {
		return nil, types.ErrorCannotCreateStableMintVault
	}
	//Checking if the appMapping_id in the msg_create & extendedPairVault_are same or not
	if appMapping.Id != extendedPairVault.AppMappingId {
		return nil, types.ErrorAppMappingIDMismatch
	}

	stableVault, found := k.GetStableMintVault(ctx, msg.StableVaultId)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	if appMapping.Id != stableVault.AppMappingId {
		return nil, types.ErrorInvalidAppMappingData
	}
	if extendedPairVault.Id != stableVault.ExtendedPairVaultID {
		return nil, types.ErrorInvalidExtendedPairMappingData
	}

	stableAmountIn := stableVault.AmountIn.Add(msg.Amount)
	if !stableAmountIn.IsPositive() {
		return nil, types.ErrorInvalidAmount
	}
	_, tokenMintedStatistics, _ := k.CheckAppExtendedPairVaultMapping(ctx, appMapping.Id, extendedPairVault.Id)

	//Check Debt Ceil
	currentMintedStatistics := tokenMintedStatistics.Add(msg.Amount)

	if currentMintedStatistics.GTE(extendedPairVault.DebtCeiling) {
		return nil, types.ErrorAmountOutGreaterThanDebtCeiling
	}

	//Take amount from user
	if err := k.SendCoinFromAccountToModule(ctx, depositorAddress, types.ModuleName, sdk.NewCoin(assetInData.Denom, msg.Amount)); err != nil {
		return nil, err
	}
	//Mint Tokens for user

	if err := k.MintCoin(ctx, types.ModuleName, sdk.NewCoin(assetOutData.Denom, msg.Amount)); err != nil {
		return nil, err
	}
	if extendedPairVault.DrawDownFee.IsZero() {
		//Send Rest to user
		if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, depositorAddress, sdk.NewCoin(assetOutData.Denom, msg.Amount)); err != nil {
			return nil, err
		}
	} else {
		//If not zero deduct send to collector//////////
		//
		//			COLLECTOR FUNCTION
		//
		//
		/////////////////////////////////////////////////

		collectorShare := (msg.Amount.Mul(sdk.Int(extendedPairVault.DrawDownFee))).Quo(sdk.Int(sdk.OneDec()))
		if err := k.SendCoinFromModuleToModule(ctx, types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, collectorShare))); err != nil {
			return nil, err
		}
		err := k.UpdateCollector(ctx, appMapping.Id, pairData.AssetOut, sdk.ZeroInt(), sdk.ZeroInt(), collectorShare, sdk.ZeroInt())
		if err != nil {
			return nil, err
		}

		// and send the rest to the user
		amountToUser := msg.Amount.Sub(collectorShare)
		if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, depositorAddress, sdk.NewCoin(assetOutData.Denom, amountToUser)); err != nil {
			return nil, err
		}
	}
	stableVault.AmountIn = stableVault.AmountIn.Add(msg.Amount)
	stableVault.AmountOut = stableVault.AmountOut.Add(msg.Amount)

	k.SetStableMintVault(ctx, stableVault)
	appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMapping(ctx, appMapping.Id)
	k.UpdateCollateralLockedAmountLockerMapping(ctx, appExtendedPairVaultData, extendedPairVault.Id, stableVault.AmountIn, true)
	k.UpdateTokenMintedAmountLockerMapping(ctx, appExtendedPairVaultData, extendedPairVault.Id, stableVault.AmountOut, true)

	return &types.MsgDepositStableMintResponse{}, nil
}

func (k *msgServer) MsgWithdrawStableMint(c context.Context, msg *types.MsgWithdrawStableMintRequest) (*types.MsgWithdrawStableMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	depositorAddress, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	//checks if extended pair exists
	extendedPairVault, found := k.GetPairsVault(ctx, msg.ExtendedPairVaultId)
	if !found {
		return nil, types.ErrorExtendedPairVaultDoesNotExists
	}
	pairData, found := k.GetPair(ctx, extendedPairVault.PairId)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}
	assetInData, found := k.GetAsset(ctx, pairData.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	assetOutData, found := k.GetAsset(ctx, pairData.AssetOut)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	//Checking if appMapping_id exists
	appMapping, found := k.GetApp(ctx, msg.AppMappingId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	//Checking if vault access disabled
	if !extendedPairVault.IsVaultActive {
		return nil, types.ErrorVaultInactive
	}
	if !extendedPairVault.IsPsmPair {
		return nil, types.ErrorCannotCreateStableMintVault
	}
	//Checking if the appMapping_id in the msg_create & extendedPairVault_are same or not
	if appMapping.Id != extendedPairVault.AppMappingId {
		return nil, types.ErrorAppMappingIDMismatch
	}

	stableVault, found := k.GetStableMintVault(ctx, msg.StableVaultId)
	if !found {
		return nil, types.ErrorVaultDoesNotExist
	}
	if appMapping.Id != stableVault.AppMappingId {
		return nil, types.ErrorInvalidAppMappingData
	}
	if extendedPairVault.Id != stableVault.ExtendedPairVaultID {
		return nil, types.ErrorInvalidExtendedPairMappingData
	}

	stableAmountIn := stableVault.AmountIn.Sub(msg.Amount)
	if stableAmountIn.LT(sdk.NewInt(0)) {
		return nil, types.ErrorInvalidAmount
	}
	var updatedAmount sdk.Int
	//Take amount from user
	if err := k.SendCoinFromAccountToModule(ctx, depositorAddress, types.ModuleName, sdk.NewCoin(assetOutData.Denom, msg.Amount)); err != nil {
		return nil, err
	}

	if extendedPairVault.DrawDownFee.IsZero() {
		//BurnTokens for user
		if err := k.BurnCoin(ctx, types.ModuleName, sdk.NewCoin(assetOutData.Denom, msg.Amount)); err != nil {
			return nil, err
		}

		//Send Rest to user
		if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, depositorAddress, sdk.NewCoin(assetInData.Denom, msg.Amount)); err != nil {
			return nil, err
		}
		updatedAmount = msg.Amount
	} else {
		//If not zero deduct send to collector//////////
		//
		//			COLLECTOR FUNCTION
		//
		//
		/////////////////////////////////////////////////
		collectorShare := (msg.Amount.Mul(sdk.Int(extendedPairVault.DrawDownFee))).Quo(sdk.Int(sdk.OneDec()))
		if err := k.SendCoinFromModuleToModule(ctx, types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetOutData.Denom, collectorShare))); err != nil {
			return nil, err
		}
		err := k.UpdateCollector(ctx, appMapping.Id, pairData.AssetOut, sdk.ZeroInt(), sdk.ZeroInt(), collectorShare, sdk.ZeroInt())
		if err != nil {
			return nil, err
		}

		updatedAmount = msg.Amount.Sub(collectorShare)

		//BurnTokens for user
		if err := k.BurnCoin(ctx, types.ModuleName, sdk.NewCoin(assetOutData.Denom, updatedAmount)); err != nil {
			return nil, err
		}

		// and send the rest to the user

		if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, depositorAddress, sdk.NewCoin(assetInData.Denom, updatedAmount)); err != nil {
			return nil, err
		}
	}
	stableVault.AmountIn = stableVault.AmountIn.Sub(updatedAmount)
	stableVault.AmountOut = stableVault.AmountOut.Sub(updatedAmount)
	k.SetStableMintVault(ctx, stableVault)
	appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMapping(ctx, appMapping.Id)
	k.UpdateCollateralLockedAmountLockerMapping(ctx, appExtendedPairVaultData, extendedPairVault.Id, stableVault.AmountIn, false)
	k.UpdateTokenMintedAmountLockerMapping(ctx, appExtendedPairVaultData, extendedPairVault.Id, stableVault.AmountOut, false)

	return &types.MsgWithdrawStableMintResponse{}, nil
}

package keeper

import (
	"context"
	"strconv"
	"time"

	"github.com/comdex-official/comdex/x/locker/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ types.MsgServer = (*msgServer)(nil)
)

type msgServer struct {
	Keeper
}

func NewMsgServiceServer(keeper Keeper) types.MsgServer {
	return &msgServer{
		Keeper: keeper,
	}
}

func (k *msgServer) MsgCreateLocker(c context.Context, msg *types.MsgCreateLockerRequest) (*types.MsgCreateLockerResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	asset, found := k.GetAsset(ctx, msg.AssetId)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	appMapping, found := k.GetApp(ctx, msg.AppMappingId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	//Checking if user mapping exists
	//if it does then check app to asset mapping has any locker key
	//if it does throw error
	userLockerAssetMapping, userExists := k.GetUserLockerAssetMapping(ctx, msg.Depositor)

	if userExists {
		_, alreadyExists := k.CheckUserAppToAssetMapping(ctx, userLockerAssetMapping, asset.Id, appMapping.Id)
		if alreadyExists {
			return nil, types.ErrorUserLockerAlreadyExists
		}
	}

	lockerProductAssetMapping, found := k.GetLockerProductAssetMapping(ctx, appMapping.Id)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	isFound := k.CheckLockerProductAssetMapping(ctx, asset.Id, lockerProductAssetMapping)
	if isFound {
		//This asset is accepted by the app
		//Create a new instance of locker

		//call Lookup table to get relevant data
		lookupTableData, exists := k.GetLockerLookupTable(ctx, lockerProductAssetMapping.AppMappingId)
		if !exists {
			return nil, types.ErrorAppMappingDoesNotExist
		}
		//Transferring amount from user to module
		depositor, err := sdk.AccAddressFromBech32(msg.Depositor)
		if err != nil {
			return nil, err
		}
		if err := k.SendCoinFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoin(asset.Denom, msg.Amount)); err != nil {
			return nil, err
		}
		//Creating locker instance
		var userLocker types.Locker
		counter := lookupTableData.Counter + 1
		userLocker.LockerId = appMapping.ShortName + strconv.FormatUint(counter, 10)
		userLocker.Depositor = msg.Depositor
		userLocker.AssetDepositId = asset.Id
		userLocker.CreatedAt = time.Now()
		userLocker.IsLocked = false
		userLocker.NetBalance = msg.Amount
		userLocker.ReturnsAccumulated = sdk.ZeroInt()
		userLocker.AppMappingId = appMapping.Id
		k.SetLocker(ctx, userLocker)
		//Checking if user data exits in mapping by user address
		//if not - create a new set
		userLockerAssetMappingData, userExists := k.GetUserLockerAssetMapping(ctx, msg.Depositor)
		if !userExists {
			//UserData does not exists
			//Create a new instance
			var userMappingData types.UserLockerAssetMapping
			var userAppData types.LockerToAppMapping
			var userAssetData types.AssetToLockerMapping
			var userTxData types.UserTxData

			userAssetData.AssetId = asset.Id
			userAssetData.LockerId = userLocker.LockerId
			userTxData.TxType = "Create"
			userTxData.Amount = msg.Amount
			userTxData.Balance = msg.Amount
			userTxData.TxTime = time.Now()
			userAssetData.UserData = append(userAssetData.UserData, &userTxData)

			userAppData.AppMappingId = appMapping.Id
			userAppData.UserAssetLocker = append(userAppData.UserAssetLocker, &userAssetData)
			userMappingData.Owner = msg.Depositor
			userMappingData.LockerAppMapping = append(userMappingData.LockerAppMapping, &userAppData)

			k.SetUserLockerAssetMapping(ctx, userMappingData)
		} else {
			///Check if user app_mapping data exits

			appExists := k.CheckUserToAppMapping(ctx, userLockerAssetMappingData, appMapping.Id)
			if appExists { //User has the app_mapping added
				//So only need to add the locker id with asset
				var userAssetData types.AssetToLockerMapping
				var userTxData types.UserTxData
				userAssetData.AssetId = asset.Id
				userAssetData.LockerId = userLocker.LockerId
				userTxData.TxType = "Create"
				userTxData.Amount = msg.Amount
				userTxData.Balance = msg.Amount
				userTxData.TxTime = time.Now()
				userAssetData.UserData = append(userAssetData.UserData, &userTxData)

				for _, appData := range userLockerAssetMappingData.LockerAppMapping {
					if appData.AppMappingId == appMapping.Id {
						appData.UserAssetLocker = append(appData.UserAssetLocker, &userAssetData)
					}
				}
				k.SetUserLockerAssetMapping(ctx, userLockerAssetMappingData)
			} else {
				//Will need to create new app and add it to the user
				var userAssetData types.AssetToLockerMapping
				var userAppData types.LockerToAppMapping
				var userTxData types.UserTxData

				userAssetData.AssetId = asset.Id
				userAssetData.LockerId = userLocker.LockerId
				userAppData.AppMappingId = appMapping.Id
				userTxData.TxType = "Create"
				userTxData.Amount = msg.Amount
				userTxData.Balance = msg.Amount
				userTxData.TxTime = time.Now()
				userAssetData.UserData = append(userAssetData.UserData, &userTxData)

				userAppData.UserAssetLocker = append(userAppData.UserAssetLocker, &userAssetData)
				userLockerAssetMappingData.LockerAppMapping = append(userLockerAssetMappingData.LockerAppMapping, &userAppData)
				k.SetUserLockerAssetMapping(ctx, userLockerAssetMappingData)
			}
		}
		k.UpdateTokenLockerMapping(ctx, lookupTableData, counter, userLocker)
	} else {
		//Not a whitelisted asset , return err
		return nil, types.ErrorLockerProductAssetMappingDoesNotExists
	}
	return &types.MsgCreateLockerResponse{}, nil
}

// MsgDepositAsset Remove asset id from Deposit & Withdraw redundant
func (k *msgServer) MsgDepositAsset(c context.Context, msg *types.MsgDepositAssetRequest) (*types.MsgDepositAssetResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	asset, found := k.GetAsset(ctx, msg.AssetId)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	appMapping, found := k.GetApp(ctx, msg.AppMappingId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}

	lockerData, found := k.GetLocker(ctx, msg.LockerId)

	if !found {
		return nil, types.ErrorLockerDoesNotExists
	}
	if lockerData.AssetDepositId != asset.Id {
		return nil, types.ErrorInvalidAssetID
	}
	if msg.Depositor != lockerData.Depositor {
		return nil, types.ErrorUnauthorized
	}
	if appMapping.Id != lockerData.AppMappingId {
		return nil, types.ErrorAppMappingDoesNotExist
	}

	lookupTableData, exists := k.GetLockerLookupTable(ctx, appMapping.Id)
	if !exists {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	depositor, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}
	if err := k.SendCoinFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoin(asset.Denom, msg.Amount)); err != nil {
		return nil, err
	}

	// if err := k.SendCoinFromModuleToModule(ctx, types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(asset.Denom, msg.Amount))); err != nil {
	// 	return nil, err
	// }

	lockerData.NetBalance = lockerData.NetBalance.Add(msg.Amount)
	k.SetLocker(ctx, lockerData)

	//Update  Amount in Locker Mapping
	k.UpdateAmountLockerMapping(ctx, lookupTableData, asset.Id, msg.Amount, true)

	userLockerAssetMappingData, _ := k.GetUserLockerAssetMapping(ctx, msg.Depositor)
	var userHisData types.UserTxData
	userHisData.TxType = "Deposit"
	userHisData.Amount = msg.Amount
	userHisData.Balance = lockerData.NetBalance
	userHisData.TxTime = time.Now()
	for _, userLockerAppData := range userLockerAssetMappingData.LockerAppMapping {
		if userLockerAppData.AppMappingId == msg.AppMappingId {
			for _, assetData := range userLockerAppData.UserAssetLocker {
				if assetData.AssetId == msg.AssetId {
					assetData.UserData = append(assetData.UserData, &userHisData)
				}
			}
		}
	}

	k.SetUserLockerAssetMapping(ctx, userLockerAssetMappingData)

	return &types.MsgDepositAssetResponse{}, nil
}

// MsgWithdrawAsset Remove asset id from Deposit & Withdraw-redundant
func (k *msgServer) MsgWithdrawAsset(c context.Context, msg *types.MsgWithdrawAssetRequest) (*types.MsgWithdrawAssetResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	asset, found := k.GetAsset(ctx, msg.AssetId)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	appMapping, found := k.GetApp(ctx, msg.AppMappingId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}

	lockerData, found := k.GetLocker(ctx, msg.LockerId)

	if !found {
		return nil, types.ErrorLockerDoesNotExists
	}
	if lockerData.AssetDepositId != asset.Id {
		return nil, types.ErrorInvalidAssetID
	}
	if msg.Depositor != lockerData.Depositor {
		return nil, types.ErrorUnauthorized
	}
	if appMapping.Id != lockerData.AppMappingId {
		return nil, types.ErrorAppMappingDoesNotExist
	}

	lookupTableData, exists := k.GetLockerLookupTable(ctx, appMapping.Id)
	if !exists {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	depositor, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}
	if lockerData.NetBalance.LT(msg.Amount) {
		return nil, types.ErrorRequestedAmountExceedsDepositAmount
	}

	lockerData.NetBalance = lockerData.NetBalance.Sub(msg.Amount)

	// if err := k.SendCoinFromModuleToModule(ctx, collectortypes.ModuleName, types.ModuleName, sdk.NewCoins(sdk.NewCoin(asset.Denom, msg.Amount))); err != nil {
	// 	return nil, err
	// }

	if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, depositor, sdk.NewCoin(asset.Denom, msg.Amount)); err != nil {
		return nil, err
	}

	k.SetLocker(ctx, lockerData)

	//Update  Amount in Locker Mapping
	k.UpdateAmountLockerMapping(ctx, lookupTableData, asset.Id, msg.Amount, false)

	userLockerAssetMappingData, _ := k.GetUserLockerAssetMapping(ctx, msg.Depositor)

	var userTxData types.UserTxData
	for _, userLockerAppData := range userLockerAssetMappingData.LockerAppMapping {
		if userLockerAppData.AppMappingId == msg.AppMappingId {
			for _, assetData := range userLockerAppData.UserAssetLocker {
				if assetData.AssetId == msg.AssetId {
					userTxData.TxType = "Withdraw"
					userTxData.Amount = msg.Amount
					userTxData.Balance = lockerData.NetBalance
					userTxData.TxTime = time.Now()
					assetData.UserData = append(assetData.UserData, &userTxData)
				}
			}
		}
	}
	k.SetUserLockerAssetMapping(ctx, userLockerAssetMappingData)

	return &types.MsgWithdrawAssetResponse{}, nil
}

func (k *msgServer) MsgAddWhiteListedAsset(c context.Context, msg *types.MsgAddWhiteListedAssetRequest) (*types.MsgAddWhiteListedAssetResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	appMapping, found := k.GetApp(ctx, msg.AppMappingId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	asset, found := k.GetAsset(ctx, msg.AssetId)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	lockerProductAssetMapping, found := k.GetLockerProductAssetMapping(ctx, msg.AppMappingId)

	if !found {
		//Set a new instance of Locker Product Asset  Mapping

		var locker types.LockerProductAssetMapping
		locker.AppMappingId = appMapping.Id
		locker.AssetIds = append(locker.AssetIds, asset.Id)
		k.SetLockerProductAssetMapping(ctx, locker)

		//Also Create a LockerLookup table Instance and set it with the new asset id
		var lockerLookupData types.LockerLookupTable
		var lockerAssetData types.TokenToLockerMapping

		lockerAssetData.AssetId = asset.Id
		lockerLookupData.Counter = 0
		lockerLookupData.AppMappingId = appMapping.Id
		lockerLookupData.Lockers = append(lockerLookupData.Lockers, &lockerAssetData)
		k.SetLockerLookupTable(ctx, lockerLookupData)
		return &types.MsgAddWhiteListedAssetResponse{}, nil
	}
	// Check if the asset from msg exists or not ,
	found = k.CheckLockerProductAssetMapping(ctx, msg.AssetId, lockerProductAssetMapping)
	if found {
		return nil, types.ErrorLockerProductAssetMappingExists
	}
	lockerProductAssetMapping.AssetIds = append(lockerProductAssetMapping.AssetIds, asset.Id)
	k.SetLockerProductAssetMapping(ctx, lockerProductAssetMapping)
	lockerLookupTableData, _ := k.GetLockerLookupTable(ctx, appMapping.Id)
	var lockerAssetData types.TokenToLockerMapping
	lockerAssetData.AssetId = asset.Id
	lockerLookupTableData.Lockers = append(lockerLookupTableData.Lockers, &lockerAssetData)
	k.SetLockerLookupTable(ctx, lockerLookupTableData)
	return &types.MsgAddWhiteListedAssetResponse{}, nil
}

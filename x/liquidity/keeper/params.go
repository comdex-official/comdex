package keeper

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/comdex-official/comdex/x/liquidity/types"
)

// GetParams returns the parameters for the liquidity module.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return
}

// SetParams sets the parameters for the liquidity module.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetGenericParams returns the parameters for the liquidity module.
func (k Keeper) GetGenericParams(ctx sdk.Context, appID uint64) (types.GenericParams, error) {
	genericParams, found := k.GetGenericLiquidityParams(ctx, appID)
	if !found {
		if ctx.IsCheckTx() {
			return types.GenericParams{}, status.Errorf(codes.NotFound, "params for app-id %d doesn't exist", appID)
		}
		_, found := k.assetKeeper.GetApp(ctx, appID)
		if !found {
			return types.GenericParams{}, sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", appID)
		}
		genericParams = types.DefaultGenericParams(appID)
		k.SetGenericParams(ctx, genericParams)
	}
	return genericParams, nil
}

// SetGenericParams sets the parameters for the liquidity module.
func (k Keeper) SetGenericParams(ctx sdk.Context, genericParams types.GenericParams) {
	k.SetGenericLiquidityParams(ctx, genericParams)
}

func (k Keeper) UpdateGenericParams(ctx sdk.Context, appID uint64, keys, values []string) error {
	_, found := k.assetKeeper.GetApp(ctx, appID)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", appID)
	}

	if len(keys) != len(values) {
		return fmt.Errorf("keys and values list length mismatch")
	}

	for _, key := range keys {
		keyFound := false
		for _, uKey := range types.UpdatableKeys {
			if uKey == key {
				keyFound = true
			}
		}
		if !keyFound {
			return fmt.Errorf("invalid key for update: %s", key)
		}
	}

	genericParams, err := k.GetGenericParams(ctx, appID)
	if err != nil {
		return err
	}

	parseValidateFunctionMap := types.KeyParseValidateFuncMap()

	for i, k := range keys {
		if parseValidateFunctionMap[k] == nil {
			return fmt.Errorf("invalid key for update: %s", k)
		}
		value := values[i]

		parsedValueInterface, err := parseValidateFunctionMap[k][0].(func(string) (interface{}, error))(value)
		if err != nil {
			return err
		}
		validationErr := parseValidateFunctionMap[k][1].(func(interface{}) error)(parsedValueInterface) //nolint:forcetypeassert //this is just for validation of error
		if validationErr != nil {
			return validationErr
		}
		reflect.ValueOf(&genericParams).Elem().FieldByName(k).Set(reflect.ValueOf(parsedValueInterface))
	}
	k.SetGenericParams(ctx, genericParams)
	return nil
}

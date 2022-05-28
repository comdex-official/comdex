package locker

import (
	"encoding/json"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func CustomQuerier(lockerKeeper *QueryPlugin) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var contractQuery LockerQuery
		if err := json.Unmarshal(request, &contractQuery); err != nil {
			return nil, sdkerrors.Wrap(err, "locker query")
			if contractQuery.State != nil {
				address := contractQuery.State.Address
				denom := contractQuery.State.Denom
				height := contractQuery.State.Height
				target := contractQuery.State.Target
				res := State{
					address,
					denom,
					height,
					target,
				}
				bz, err := json.Marshal(res)
				if err != nil {
					return nil, sdkerrors.Wrap(err, "locker state query response")
				}
				return bz, nil
			}
		}
		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown osmosis query variant"}
	}

}

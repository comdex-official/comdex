package expected

import (
	"context"
	wasmvmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ContractOpsKeeper contains mutable operations on a contract.
type ContractOpsKeeper interface {
	// Sudo allows to call privileged entry point of a contract.
	Sudo(ctx context.Context, contractAddress sdk.AccAddress, msg []byte) ([]byte, error)
	GetContractInfo(ctx context.Context, contractAddress sdk.AccAddress) *wasmvmtypes.ContractInfo
}
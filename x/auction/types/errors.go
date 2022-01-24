package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidInitialAuctionID = errors.Register(ModuleName, 101, "initial auction ID hasn't been set")
)

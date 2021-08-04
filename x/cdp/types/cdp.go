package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewCDP(id uint64, owner sdk.AccAddress, collateral sdk.Coin, collateralType string, debt sdk.Coin) CDP {
	return CDP{
		Id:         id,
		Owner:      owner.String(),
		Type:       collateralType,
		Collateral: collateral,
		Debt:       debt,
	}
}

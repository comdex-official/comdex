package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewCDP(id uint64, owner sdk.AccAddress, collateral sdk.Coin, collateralType string, principal sdk.Coin) CDP {
	return CDP{
		Id:         id,
		Owner:      owner.String(),
		Type:       collateralType,
		Collateral: collateral,
		Principal:  principal,
	}
}

func NewDeposit(cdpID uint64, depositor sdk.AccAddress, amount sdk.Coin) Deposit {
	return Deposit{cdpID, depositor.String(), amount}
}

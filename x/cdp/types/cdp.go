package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

func NewCDP(id uint64, owner sdk.AccAddress, collateral sdk.Coin, collateralType string, principal sdk.Coin, time time.Time ) CDP {
	fees := sdk.NewCoin(principal.Denom, sdk.ZeroInt())
	return CDP{
		Id: id,
		Owner: owner.String(),
		Type: collateralType,
		Collateral: collateral,
		Principal: principal,
		AccumulatedFees: fees,
		FeesUpdated: time,
		InterestFactor: sdk.ZeroDec(),
	}
}

func NewDeposit(cdpID uint64, depositor sdk.AccAddress, amount sdk.Coin) Deposit {
	return Deposit{cdpID, depositor.String(), amount}
}
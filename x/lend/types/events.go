package types

// Event types for the Lend module.
const (
	EventTypeLend            = "lend"
	EventTypeWithdraw        = "withdraw"
	EventTypeDeposit         = "deposit"
	EventTypeClose           = "close"
	EventTypeBorrow          = "borrow"
	EventTypeRepay           = "repay"
	EventTypeDepositBorrow   = "depositBorrow"
	EventTypeDraw            = "draw"
	EventTypeCloseBorrow     = "closeBorrow"
	EventTypeBorrowAlternate = "borrowAlternate"
	EventTypeFundModuleAccn  = "fundModuleAccn"
	EventTypeBorrowInterest  = "borrowInterest"
	EventTypeLendRewards     = "lendRewards"

	AttributeKeyCreator   = "creator"
	AttributeKeyAppID     = "appId"
	AttributeKeyPoolID    = "PoolId"
	AttributeKeyAssetID   = "AssetId"
	AttributeKeyAmountIn  = "amountIn"
	AttributeKeyAmountOut = "amountOut"
	AttributeKeyTimestamp = "timestamp"
	AttributeKeyLendID    = "lendId"
	AttributeKeyPairID    = "pairId"
	AttributeKeyIsStable  = "isStableBorrow"
	AttributeKeyBorrowID  = "borrowId"
)

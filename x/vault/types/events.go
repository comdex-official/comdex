package types

// Event types for the Vault module.
const (
	EventTypeCreateVault   = "create_vault"
	EventTypeDepositVault  = "deposit_vault"
	EventTypeWithdrawVault = "withdraw_vault"
	EventTypeDrawVault     = "draw_vault"
	EventTypeRepayVault    = "repay_vault"
	EventTypeCloseVault    = "close_vault"

	AttributeKeyVaultID               = "vaultId"
	AttributeKeyCreator               = "creator"
	AttributeKeyAppID                 = "appId"
	AttributeKeyExtendedPairID        = "extendedPairId"
	AttributeKeyAmountIn              = "amountIn"
	AttributeKeyAmountOut             = "amountOut"
	AttributeKeyCreatedAt             = "createdAt"
	AttributeKeyInterestAccumulated   = "interestAccumulated"
	AttributeKeyClosingFeeAccumulated = "closingFeeAccumulated"
)

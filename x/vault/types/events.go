package types

// Event types for the Vault module.
const (
	EventTypeCreateVault = "create_vault"
	EventTypeCloseVault  = "close_vault"

	AttributeKeyVaultID               = "vaultId"
	AttributeKeyCreator               = "creator"
	AttributeKeyAppID                 = "appId"
	AttributeKeyAmountIn              = "amountIn"
	AttributeKeyAmountOut             = "amountOut"
	AttributeKeyCreatedAt             = "createdAt"
	AttributeKeyInterestAccumulated   = "interestAccumulated"
	AttributeKeyClosingFeeAccumulated = "closingFeeAccumulated"
)

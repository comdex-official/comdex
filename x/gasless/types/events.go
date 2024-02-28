package types

// Event types for the gasless module.
const (
	EventTypeCreateGasProvider = "create_gas_provider"

	AttributeKeyCreator       = "creator"
	AttributeKeyGasProviderId = "gas_provider_id"
	AttributeKeyFeeDenom      = "fee_denom"
)

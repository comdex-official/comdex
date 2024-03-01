package types

// Event types for the gasless module.
const (
	EventTypeCreateGasProvider       = "create_gas_provider"
	EventTypeAuthorizeActors         = "authorize_actors"
	EventTypeUpdateGasProviderStatus = "update_gas_provider_status"
	EventTypeUpdateGasProviderConfig = "update_gas_provider_config"
	EventTypeBlockConsumer           = "block_consumer"
	EventTypeUnblockConsumer         = "unblock_consumer"
	EventTypeFeeConsumption          = "fee_consumption"

	AttributeKeyCreator                 = "creator"
	AttributeKeyProvider                = "provider"
	AttributeKeyActor                   = "actor"
	AttributeKeyConsumer                = "consumer"
	AttributeKeyGasProviderID           = "gas_provider_id"
	AttributeKeyFeeDenom                = "fee_denom"
	AttributeKeyAuthorizedActors        = "authorized_actors"
	AttributeKeyGasProviderStatus       = "gas_provider_status"
	AttributeKeyMaxFeeUsagePerTx        = "max_fee_usage_per_tx"
	AttributeKeyMaxTxsCountPerConsumer  = "max_txs_count_per_consumer"
	AttributeKeyMaxFeeUsagePerConsumer  = "max_fee_usage_per_consumer"
	AttributeKeyTxsAllowed              = "txs_allowed"
	AttributeKeyContractsAllowed        = "contracts_allowed"
	AttributeKeyFeeConsumptionMessage   = "message"
	AttributeKeyFeeSource               = "fee_source"
	AttributeKeyFailedGasProviderIDs    = "failed_gas_provider_ids"
	AttributeKeyFailedGasProviderErrors = "failed_gas_provider_errors"
	AttributeKeySucceededGpid           = "succeeded_gas_provider_id"
)

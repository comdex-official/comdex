package types

const (
	EventTypeSurplusActivatorErr   = "surplus_activator_err"
	EventTypeDebtActivatorErr      = "debt_activator_err"
	EventTypeRestartDutchErr       = "restart_dutch_err"
	EventTypeRestartLendDutchErr   = "restart_lend_dutch_err"
	EventTypeDutchNewAuction       = "dutch_new_auction"
	AttributeKeyOwner              = "vault_owner"
	AttributeKeyCollateral         = "collateral_token"
	AttributeKeyDebt               = "debt_token"
	AttributeKeyStartTime          = "start_time"
	AttributeKeyEndTime            = "end_time"
	DataAppID                      = "data_app_id"
	DataAssetID                    = "data_asset_id"
	DataAssetOutOraclePrice        = "data_asset_out_oracle_price"
	DataAssetOutPrice              = "data_asset_out_price"
	DatIsAuctionActive             = "data_is_auction_active"
	DataIsDebtAuction              = "data_is_debt_auction"
	DataIsDistributor              = "data_is_distributor"
	DataIsSurplusAuction           = "data_is_surplus_auction"
	KillSwitchParamsBreakerEnabled = "data_is_surplus_auction"
	Status                         = "status"
)

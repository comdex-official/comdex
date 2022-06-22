package bindings

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ComdexMessages struct {
	MsgWhiteListAssetLocker              *MsgWhiteListAssetLocker              `json:"msg_white_list_asset_locker,omitempty"`
	MsgWhitelistAppIDVaultInterest       *MsgWhitelistAppIDVaultInterest       `json:"msg_whitelist_app_id_vault_interest,omitempty"`
	MsgWhitelistAppIDLockerRewards       *MsgWhitelistAppIDLockerRewards       `json:"msg_whitelist_app_id_locker_rewards,omitempty"`
	MsgAddExtendedPairsVault             *MsgAddExtendedPairsVault             `json:"msg_add_extended_pairs_vault,omitempty"`
	MsgSetCollectorLookupTable           *MsgSetCollectorLookupTable           `json:"msg_set_collector_lookup_table,omitempty"`
	MsgSetAuctionMappingForApp           *MsgSetAuctionMappingForApp           `json:"msg_set_auction_mapping_for_app,omitempty"`
	MsgUpdatePairsVault                  *MsgUpdatePairsVault                  `json:"msg_update_pairs_vault,omitempty"`
	MsgUpdateCollectorLookupTable        *MsgUpdateCollectorLookupTable        `json:"msg_update_collector_lookup_table,omitempty"`
	MsgRemoveWhitelistAssetLocker        *MsgRemoveWhitelistAssetLocker        `json:"msg_remove_whitelist_asset_locker,omitempty"`
	MsgRemoveWhitelistAppIDVaultInterest *MsgRemoveWhitelistAppIDVaultInterest `json:"msg_remove_whitelist_app_id_vault_interest,omitempty"`
	MsgWhitelistAppIDLiquidation         *MsgWhitelistAppIDLiquidation         `json:"msg_whitelist_app_id_liquidation,omitempty"`
	MsgRemoveWhitelistAppIDLiquidation   *MsgRemoveWhitelistAppIDLiquidation   `json:"msg_remove_whitelist_app_id_liquidation,omitempty"`
	MsgAddAuctionParams                  *MsgAddAuctionParams                  `json:"msg_add_auction_params,omitempty"`
	MsgBurnGovTokensForApp               *MsgBurnGovTokensForApp               `json:"msg_burn_gov_tokens_for_app,omitempty"`
}

type MsgWhiteListAssetLocker struct {
	AppMappingID uint64 `json:"app_mapping_id"`
	AssetID      uint64 `json:"asset_id"`
}

type MsgWhitelistAppIDVaultInterest struct {
	AppMappingID uint64 `json:"app_mapping_id"`
}

type MsgWhitelistAppIDLockerRewards struct {
	AppMappingID uint64   `json:"app_mapping_id"`
	AssetIDs     []uint64 `json:"asset_ids"`
}

type MsgAddExtendedPairsVault struct {
	AppMappingID        uint64  `json:"app_mapping_id"`
	PairID              uint64  `json:"pair_id"`
	LiquidationRatio    sdk.Dec `json:"liquidation_ratio"`
	StabilityFee        sdk.Dec `json:"stability_fee"`
	ClosingFee          sdk.Dec `json:"closing_fee"`
	LiquidationPenalty  sdk.Dec `json:"liquidation_penalty"`
	DrawDownFee         sdk.Dec `json:"draw_down_fee"`
	IsVaultActive       bool    `json:"is_vault_active"`
	DebtCeiling         uint64  `json:"debt_ceiling"`
	DebtFloor           uint64  `json:"debt_floor"`
	IsPsmPair           bool    `json:"is_psm_pair"`
	MinCr               sdk.Dec `json:"min_cr"`
	PairName            string  `json:"pair_name"`
	AssetOutOraclePrice bool    `json:"asset_out_oracle_price"`
	AssetOutPrice       uint64  `json:"asset_out_price"`
	MinUsdValueLeft     uint64  `json:"min_usd_value_left"`
}

type MsgSetCollectorLookupTable struct {
	AppMappingID     uint64  `json:"app_mapping_id"`
	CollectorAssetID uint64  `json:"collector_asset_id"`
	SecondaryAssetID uint64  `json:"secondary_asset_id"`
	SurplusThreshold uint64  `json:"surplus_threshold"`
	DebtThreshold    uint64  `json:"debt_threshold"`
	LockerSavingRate sdk.Dec `json:"locker_saving_rate"`
	LotSize          uint64  `json:"lot_size"`
	BidFactor        sdk.Dec `json:"bid_factor"`
	DebtLotSize      uint64  `json:"debt_lot_size"`
}

type MsgSetAuctionMappingForApp struct {
	AppMappingID        uint64   `json:"app_mapping_id"`
	AssetID             []uint64 `json:"asset_id"`
	IsSurplusAuction    []bool   `json:"is_surplus_auction"`
	IsDebtAuction       []bool   `json:"is_debt_auction"`
	AssetOutOraclePrice []bool   `json:"asset_out_oracle_price"`
	AssetOutPrice       []uint64 `json:"asset_out_price"`
}

type MsgUpdatePairsVault struct {
	AppMappingID       uint64  `json:"app_mapping_id"`
	ExtPairID          uint64  `json:"ext_pair_id"`
	LiquidationRatio   sdk.Dec `json:"liquidation_ratio"`
	StabilityFee       sdk.Dec `json:"stability_fee"`
	ClosingFee         sdk.Dec `json:"closing_fee"`
	LiquidationPenalty sdk.Dec `json:"liquidation_penalty"`
	DrawDownFee        sdk.Dec `json:"draw_down_fee"`
	MinCr              sdk.Dec `json:"min_cr"`
	DebtCeiling        uint64  `json:"debt_ceiling"`
	DebtFloor          uint64  `json:"debt_floor"`
	MinUsdValueLeft    uint64  `json:"min_usd_value_left"`
}

type MsgUpdateCollectorLookupTable struct {
	AppMappingID     uint64  `json:"app_mapping_id"`
	AssetID          uint64  `json:"asset_id"`
	DebtThreshold    uint64  `json:"debt_threshold"`
	SurplusThreshold uint64  `json:"surplus_threshold"`
	LotSize          uint64  `json:"lot_size"`
	DebtLotSize      uint64  `json:"debt_lot_size"`
	BidFactor        sdk.Dec `json:"bid_factor"`
	LSR              sdk.Dec `json:"lsr"`
}

type MsgRemoveWhitelistAssetLocker struct {
	AppMappingID uint64 `json:"app_mapping_id"`
	AssetID      uint64 `json:"asset_id"`
}

type MsgRemoveWhitelistAppIDVaultInterest struct {
	AppMappingID uint64 `json:"app_mapping_id"`
}

type MsgWhitelistAppIDLiquidation struct {
	AppMappingID uint64 `json:"app_mapping_id"`
}

type MsgRemoveWhitelistAppIDLiquidation struct {
	AppMappingID uint64 `json:"app_mapping_id"`
}

type MsgAddAuctionParams struct {
	AppMappingID           uint64  `json:"app_mapping_id"`
	AuctionDurationSeconds uint64  `json:"auction_duration_seconds"`
	Buffer                 sdk.Dec `json:"buffer"`
	Cusp                   sdk.Dec `json:"cusp"`
	Step                   uint64  `json:"step"`
	PriceFunctionType      uint64  `json:"price_function_type"`
	SurplusID              uint64  `json:"surplus_id"`
	DebtID                 uint64  `json:"debt_id"`
	DutchID                uint64  `json:"dutch_id"`
	BidDurationSeconds     uint64  `json:"bid_duration_seconds"`
}

type MsgBurnGovTokensForApp struct {
	AppMappingID uint64         `json:"app_mapping_id"`
	From         sdk.AccAddress `json:"from"`
	Amount       sdk.Coin       `json:"amount"`
}

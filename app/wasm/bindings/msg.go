package bindings

import (
	sdkmath "cosmossdk.io/math"
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
	MsgAddESMTriggerParams               *MsgAddESMTriggerParams               `json:"msg_add_e_s_m_trigger_params,omitempty"`
	MsgEmissionRewards                   *MsgEmissionRewards                   `json:"msg_emission_rewards,omitempty"`
	MsgFoundationEmission                *MsgFoundationEmission                `json:"msg_foundation_emission,omitempty"`
	MsgRebaseMint                        *MsgRebaseMint                        `json:"msg_rebase_mint,omitempty"`
	MsgGetSurplusFund                    *MsgGetSurplusFund                    `json:"msg_get_surplus_fund,omitempty"`
	MsgEmissionPoolRewards               *MsgEmissionPoolRewards               `json:"msg_emission_pool_rewards,omitempty"`
}

type MsgWhiteListAssetLocker struct {
	AppID   uint64 `json:"app_id"`
	AssetID uint64 `json:"asset_id"`
}

type MsgWhitelistAppIDVaultInterest struct {
	AppID uint64 `json:"app_id"`
}

type MsgWhitelistAppIDLockerRewards struct {
	AppID   uint64 `json:"app_id"`
	AssetID uint64 `json:"asset_id"`
}

type MsgAddExtendedPairsVault struct {
	AppID               uint64            `json:"app_id"`
	PairID              uint64            `json:"pair_id"`
	StabilityFee        sdkmath.LegacyDec `json:"stability_fee"`
	ClosingFee          sdkmath.LegacyDec `json:"closing_fee"`
	LiquidationPenalty  sdkmath.LegacyDec `json:"liquidation_penalty"`
	DrawDownFee         sdkmath.LegacyDec `json:"draw_down_fee"`
	IsVaultActive       bool              `json:"is_vault_active"`
	DebtCeiling         sdkmath.Int       `json:"debt_ceiling"`
	DebtFloor           sdkmath.Int       `json:"debt_floor"`
	IsStableMintVault   bool              `json:"is_stable_mint_vault"`
	MinCr               sdkmath.LegacyDec `json:"min_cr"`
	PairName            string            `json:"pair_name"`
	AssetOutOraclePrice bool              `json:"asset_out_oracle_price"`
	AssetOutPrice       uint64            `json:"asset_out_price"`
	MinUsdValueLeft     uint64            `json:"min_usd_value_left"`
}

type MsgSetCollectorLookupTable struct {
	AppID            uint64            `json:"app_id"`
	CollectorAssetID uint64            `json:"collector_asset_id"`
	SecondaryAssetID uint64            `json:"secondary_asset_id"`
	SurplusThreshold sdkmath.Int       `json:"surplus_threshold"`
	DebtThreshold    sdkmath.Int       `json:"debt_threshold"`
	LockerSavingRate sdkmath.LegacyDec `json:"locker_saving_rate"`
	LotSize          sdkmath.Int       `json:"lot_size"`
	BidFactor        sdkmath.LegacyDec `json:"bid_factor"`
	DebtLotSize      sdkmath.Int       `json:"debt_lot_size"`
}

type MsgSetAuctionMappingForApp struct {
	AppID                uint64 `json:"app_id"`
	AssetIDs             uint64 `json:"asset_id"`
	IsSurplusAuctions    bool   `json:"is_surplus_auction"`
	IsDebtAuctions       bool   `json:"is_debt_auction"`
	IsDistributor        bool   `json:"is_distributor"`
	AssetOutOraclePrices bool   `json:"asset_out_oracle_price"`
	AssetOutPrices       uint64 `json:"asset_out_price"`
}

type MsgUpdatePairsVault struct {
	AppID              uint64            `json:"app_id"`
	ExtPairID          uint64            `json:"ext_pair_id"`
	StabilityFee       sdkmath.LegacyDec `json:"stability_fee"`
	ClosingFee         sdkmath.LegacyDec `json:"closing_fee"`
	LiquidationPenalty sdkmath.LegacyDec `json:"liquidation_penalty"`
	DrawDownFee        sdkmath.LegacyDec `json:"draw_down_fee"`
	IsVaultActive      bool              `json:"is_vault_active"`
	MinCr              sdkmath.LegacyDec `json:"min_cr"`
	DebtCeiling        sdkmath.Int       `json:"debt_ceiling"`
	DebtFloor          sdkmath.Int       `json:"debt_floor"`
	MinUsdValueLeft    uint64            `json:"min_usd_value_left"`
}

type MsgUpdateCollectorLookupTable struct {
	AppID            uint64            `json:"app_id"`
	AssetID          uint64            `json:"asset_id"`
	DebtThreshold    sdkmath.Int       `json:"debt_threshold"`
	SurplusThreshold sdkmath.Int       `json:"surplus_threshold"`
	LotSize          sdkmath.Int       `json:"lot_size"`
	DebtLotSize      sdkmath.Int       `json:"debt_lot_size"`
	BidFactor        sdkmath.LegacyDec `json:"bid_factor"`
	LSR              sdkmath.LegacyDec `json:"lsr"`
}

type MsgRemoveWhitelistAssetLocker struct {
	AppID   uint64 `json:"app_id"`
	AssetID uint64 `json:"asset_id"`
}

type MsgRemoveWhitelistAppIDVaultInterest struct {
	AppMappingID uint64 `json:"app_id"`
}

type MsgWhitelistAppIDLiquidation struct {
	AppID uint64 `json:"app_id"`
}

type MsgRemoveWhitelistAppIDLiquidation struct {
	AppID uint64 `json:"app_id"`
}

type MsgAddAuctionParams struct {
	AppID                  uint64            `json:"app_id"`
	AuctionDurationSeconds uint64            `json:"auction_duration_seconds"`
	Buffer                 sdkmath.LegacyDec `json:"buffer"`
	Cusp                   sdkmath.LegacyDec `json:"cusp"`
	Step                   uint64            `json:"step"`
	PriceFunctionType      uint64            `json:"price_function_type"`
	SurplusID              uint64            `json:"surplus_id"`
	DebtID                 uint64            `json:"debt_id"`
	DutchID                uint64            `json:"dutch_id"`
	BidDurationSeconds     uint64            `json:"bid_duration_seconds"`
}

type MsgBurnGovTokensForApp struct {
	AppID  uint64         `json:"app_id"`
	From   sdk.AccAddress `json:"from"`
	Amount sdk.Coin       `json:"amount"`
}

type MsgAddESMTriggerParams struct {
	AppID         uint64   `json:"app_id"`
	TargetValue   sdk.Coin `json:"target_value"`
	CoolOffPeriod uint64   `json:"cool_off_period"`
	AssetID       []uint64 `json:"asset_id"`
	Rates         []uint64 `json:"rates"`
}

type MsgEmissionRewards struct {
	AppID          uint64        `json:"app_id"`
	Amount         sdkmath.Int   `json:"amount"`
	EmissionAmount uint64        `json:"emission_amount"`
	ExtendedPair   []uint64      `json:"extended_pair"`
	VotingRatio    []sdkmath.Int `json:"voting_ratio"`
}

type MsgFoundationEmission struct {
	AppID             uint64      `json:"app_id"`
	Amount            sdkmath.Int `json:"amount"`
	FoundationAddress []string    `json:"foundation_address"`
}

type MsgRebaseMint struct {
	AppID        uint64         `json:"app_id"`
	Amount       sdkmath.Int    `json:"amount"`
	ContractAddr sdk.AccAddress `json:"contract_addr"`
}

type MsgGetSurplusFund struct {
	AppID        uint64         `json:"app_id"`
	AssetID      uint64         `json:"asset_id"`
	ContractAddr sdk.AccAddress `json:"contract_addr"`
	Amount       sdk.Coin       `json:"amount"`
}

type MsgEmissionPoolRewards struct {
	AppID       uint64        `json:"app_id"`
	CswapAppID  uint64        `json:"cswap_app_id"`
	Amount      sdkmath.Int   `json:"amount"`
	Pools       []uint64      `json:"pools"`
	VotingRatio []sdkmath.Int `json:"voting_ratio"`
}

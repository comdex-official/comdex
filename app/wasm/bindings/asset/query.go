package asset

type AppQuery struct {
	AppData *AppData `json:"get_app,omitempty"`
}

type AppData struct {
	App_Id uint64 `json:"app_mapping_id"`
}

type AppDataResponse struct {
	MinGovDeposit    int64  `json:"min_gov_deposit"`
	GovTimeInSeconds int64  `json:"gov_time_in_seconds"`
	AssetId          uint64 `json:"gov_token_id"` //only when isGovToken true
}

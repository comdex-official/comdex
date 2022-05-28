package asset

import "github.com/comdex-official/comdex/x/asset/types"

type AppQuery struct {
	AppData *AppData `json:"app_data,omitempty"`
}

type AppData struct {
	App_Id uint64 `json:"app_id"`
}

type AppDataResponse struct {
	AppDatas types.AppMapping `json:"app_datas"`
}

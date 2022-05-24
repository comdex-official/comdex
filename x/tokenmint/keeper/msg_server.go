package keeper

import (
	"context"
	"time"

	"github.com/comdex-official/comdex/x/tokenmint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ types.MsgServiceServer = (*msgServer)(nil)
)

type msgServer struct {
	Keeper
}

func NewMsgServiceServer(keeper Keeper) types.MsgServiceServer {
	return &msgServer{
		Keeper: keeper,
	}
}

//protostructure
//	message TokenMint
//	app_id:=
//	repeated token data
//			--- asset id
//          ----genesis supply
//			--current stats
//
//
//

//MsgMintGenesisToken
//Take app_mapping_id from user , take asset id
//check if app mapping exists ,
//    if exists
//check for if asset exists
//if asset exissts and genesis supplymint is true
// send error already minted
//if asset does not exists
//check data in app_mapping table - if data recieved from proposal and data for the asset is in the app mapping kv
//if does
// mint asset send to user save data to kv store
//else
//error
//if does not exists
//go to app mapping and check if asset data for that app exists in app_mapping_table or not
//if it does
//mint and update
//if not
//error

//Mint request function --- to mint asset
// if app & asset id exits 	--mint and update current data - have user address to send
//Do Same for burn as well
//only difference you dont need to have user address, will get tokens from the user

//List Of function required

//1. GetTokenMintData
//2. SetTOkenMintData
//3. GetAssetDataInTokenMintByApp
//4. UpdateAssetDataInTOkenMintByApp--- +- of current stats only
//5. CheckAppMappingData- write here check in asseetmodule

func (k *msgServer) MsgMintNewTokens(c context.Context, msg *types.MsgMintNewTokensRequest) (*types.MsgMintNewTokensResponse, error) {

	ctx := sdk.UnwrapSDKContext(c)

	
	mintData,found:=k.GetTokenMint(ctx,msg.AppMappingId)
	if !found {
		var newTokenMintappData types.MintedTokens
		var appData types.TokenMint

		newTokenMintappData.AssetId=msg.AssetId
		newTokenMintappData.CreatedAt=time.Now()
		newTokenMintappData.
		
	}else{

	}



	return &types.MsgMintNewTokensResponse{}, nil

}

package tokenmint

import (
	tokenMintKeeper "github.com/comdex-official/comdex/x/tokenmint/keeper"
	tokenMinttypes "github.com/comdex-official/comdex/x/tokenmint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type QueryPlugin struct {
	tokenMintKeeper *tokenMintKeeper.Keeper
}

func NewQueryPlugin(
	tokenMintKeeper *tokenMintKeeper.Keeper,
) *QueryPlugin {
	return &QueryPlugin{
		tokenMintKeeper: tokenMintKeeper,
	}
}

func (qp QueryPlugin) GetTokenMint(ctx sdk.Context, appMappingId, assetId uint64) (tokenMinttypes.MintedTokens, error) {
	tokenData, err := qp.tokenMintKeeper.GetAssetDataInTokenMintByApp(ctx, appMappingId, assetId)
	if err != true {
		return tokenData, nil
	}
	return tokenData, nil
}

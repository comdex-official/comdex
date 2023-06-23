package v2_test

// import (
// 	"fmt"
// 	"strings"
// 	"testing"

// 	chain "github.com/comdex-official/comdex/app"
// 	assettypes "github.com/comdex-official/comdex/x/asset/types"
// 	v1liquidity "github.com/comdex-official/comdex/x/liquidity/legacy/v1"
// 	v2liquidity "github.com/comdex-official/comdex/x/liquidity/legacy/v2"
// 	"github.com/comdex-official/comdex/x/liquidity/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/stretchr/testify/require"
// 	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
// )

// func CreateNewTestApp(t *testing.T, app *chain.App, ctx sdk.Context, appName string) uint64 {
// 	if app.AssetKeeper.HasAppForName(ctx, appName) {
// 		panic(fmt.Sprintf("app %s alredy exists", appName))
// 	}
// 	err := app.AssetKeeper.AddAppRecords(ctx, assettypes.AppData{
// 		Name:             strings.ToLower(appName),
// 		ShortName:        strings.ToLower(appName),
// 		MinGovDeposit:    sdk.NewInt(0),
// 		GovTimeInSeconds: 0,
// 		GenesisToken:     []assettypes.MintGenesisToken{},
// 	})
// 	require.NoError(t, err)
// 	found := app.AssetKeeper.HasAppForName(ctx, appName)
// 	require.True(t, found)

// 	apps, found := app.AssetKeeper.GetApps(ctx)
// 	require.True(t, found)
// 	var appID uint64
// 	for _, app := range apps {
// 		if app.Name == appName {
// 			appID = app.Id
// 			break
// 		}
// 	}
// 	require.NotZero(t, appID)
// 	return appID

// }

// func GetOldPool(appID, poolId uint64, disabled bool) v1liquidity.Pool {
// 	return v1liquidity.Pool{
// 		Id:                    poolId,
// 		PairId:                poolId,
// 		ReserveAddress:        string(types.PoolReserveAddress(appID, poolId)),
// 		PoolCoinDenom:         fmt.Sprintf("pool%d-%d", appID, poolId),
// 		LastDepositRequestId:  2,
// 		LastWithdrawRequestId: 3,
// 		Disabled:              disabled,
// 		AppId:                 appID,
// 	}
// }

// func TestMigratePools(t *testing.T) {
// 	app := chain.Setup(false)
// 	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

// 	appID := CreateNewTestApp(t, app, ctx, "appone")
// 	fmt.Println("......", appID)

// 	cdc := chain.MakeTestEncodingConfig().Marshaler
// 	store := ctx.KVStore(sdk.NewKVStoreKey(types.StoreKey))

// 	oldPool1 := GetOldPool(appID, 1, false)
// 	store.Set(types.GetPoolKey(oldPool1.AppId, oldPool1.Id), cdc.MustMarshal(&oldPool1))

// 	oldPool2 := GetOldPool(appID, 2, true)
// 	store.Set(types.GetPoolKey(oldPool2.AppId, oldPool2.Id), cdc.MustMarshal(&oldPool2))

// 	oldPool3 := GetOldPool(appID, 3, false)
// 	store.Set(types.GetPoolKey(oldPool3.AppId, oldPool3.Id), cdc.MustMarshal(&oldPool3))

// 	require.NoError(t, v2liquidity.MigratePools(appID, store, cdc))
// }

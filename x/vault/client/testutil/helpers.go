package testutil

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/app/wasm/bindings"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"
	"github.com/comdex-official/comdex/x/vault/client/cli"
	"github.com/comdex-official/comdex/x/vault/types"
)

var commonArgs = []string{
	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)).String()),
}

// via cli
func MsgCreate(
	clientCtx client.Context,
	appMappingID, extendedPairVaultID uint64,
	amountIn, amountOut sdk.Int,
	from string,
	extraArgs ...string,
) (testutil.BufferWriter, error) {
	args := append(append([]string{
		strconv.Itoa(int(appMappingID)),
		strconv.Itoa(int(extendedPairVaultID)),
		amountIn.String(),
		amountOut.String(),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...), extraArgs...)

	resp, err := clitestutil.ExecTestCLICmd(clientCtx, cli.Create(), args)
	if err != nil {
		return resp, err
	}
	var respJSON map[string]interface{}
	err = json.Unmarshal([]byte(resp.String()), &respJSON)
	if err != nil {
		return nil, err
	}
	if respJSON["code"] != 0 {
		errLog, _ := respJSON["raw_log"].(string)
		err = fmt.Errorf(errLog)
	}
	return resp, err
}

func (s *VaultIntegrationTestSuite) fundAddr(addr sdk.AccAddress, amt sdk.Coins) { //nolint:unused
	s.T().Helper()
	err := s.app.BankKeeper.MintCoins(s.ctx, types.ModuleName, amt)
	s.Require().NoError(err)
	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, types.ModuleName, addr, amt)
	s.Require().NoError(err)
}

func (s *VaultIntegrationTestSuite) CreateNewApp(appName string) uint64 {
	err := s.app.AssetKeeper.AddAppRecords(s.ctx, assettypes.AppData{
		Name:             appName,
		ShortName:        appName,
		MinGovDeposit:    sdk.NewInt(0),
		GovTimeInSeconds: 0,
		GenesisToken:     []assettypes.MintGenesisToken{},
	})
	s.Require().NoError(err)
	found := s.app.AssetKeeper.HasAppForName(s.ctx, appName)
	s.Require().True(found)

	apps, found := s.app.AssetKeeper.GetApps(s.ctx)
	s.Require().True(found)
	var appID uint64
	for _, app := range apps {
		if app.Name == appName {
			appID = app.Id
			break
		}
	}
	s.Require().NotZero(appID)
	return appID
}

func (s *VaultIntegrationTestSuite) SetOraclePrice(symbol string, price uint64) {
	var (
		store = s.app.MarketKeeper.Store(s.ctx)
		key   = markettypes.PriceForMarketKey(symbol)
	)
	value := s.cfg.Codec.MustMarshal(
		&protobuftypes.UInt64Value{
			Value: price,
		},
	)
	store.Set(key, value)
}

func (s *VaultIntegrationTestSuite) CreateNewAsset(name, denom string, price uint64) uint64 {
	err := s.app.AssetKeeper.AddAssetRecords(s.ctx, assettypes.Asset{
		Name:                  name,
		Denom:                 denom,
		Decimals:              1000000,
		IsOnChain:             true,
		IsOraclePriceRequired: true,
	})
	s.Require().NoError(err)
	assets := s.app.AssetKeeper.GetAssets(s.ctx)
	var assetID uint64
	for _, asset := range assets {
		if asset.Denom == denom {
			assetID = asset.Id
			break
		}
	}
	s.Require().NotZero(assetID)

	market := markettypes.Market{
		Symbol:   name,
		ScriptID: 12,
		Rates:    price,
	}
	s.app.MarketKeeper.SetMarket(s.ctx, market)

	exists := s.app.MarketKeeper.HasMarketForAsset(s.ctx, assetID)
	s.Suite.Require().False(exists)
	s.app.MarketKeeper.SetMarketForAsset(s.ctx, assetID, name)
	exists = s.app.MarketKeeper.HasMarketForAsset(s.ctx, assetID)
	s.Suite.Require().True(exists)

	s.SetOraclePrice(name, price)

	return assetID
}

func (s *VaultIntegrationTestSuite) CreateNewPair(assetIn, assetOut uint64) uint64 {
	_, err := s.app.AssetKeeper.NewAddPair(s.ctx, &assettypes.MsgAddPairRequest{
		From:     s.val.Address.String(),
		AssetIn:  assetIn,
		AssetOut: assetOut,
	})
	s.Suite.NoError(err)
	pairs := s.app.AssetKeeper.GetPairs(s.ctx)
	var pairID uint64
	for _, pair := range pairs {
		if pair.AssetIn == assetIn && pair.AssetOut == assetOut {
			pairID = pair.Id
			break
		}
	}
	s.Require().NotZero(pairID)
	return pairID
}

func (s *VaultIntegrationTestSuite) CreateNewExtendedVaultPair(pairName string, appMappingID, pairID uint64) uint64 {
	err := s.app.AssetKeeper.WasmAddExtendedPairsVaultRecords(s.ctx, &bindings.MsgAddExtendedPairsVault{
		AppID:               appMappingID,
		PairID:              pairID,
		StabilityFee:        sdk.NewDecWithPrec(2, 2), // 0.02
		ClosingFee:          sdk.NewDec(0),
		LiquidationPenalty:  sdk.NewDecWithPrec(15, 2), // 0.15
		DrawDownFee:         sdk.NewDecWithPrec(1, 2),  // 0.01
		IsVaultActive:       true,
		DebtCeiling:         1000000000000000000,
		DebtFloor:           100000000,
		IsStableMintVault:   false,
		MinCr:               sdk.NewDecWithPrec(23, 1), // 2.3
		PairName:            pairName,
		AssetOutOraclePrice: true,
		AssetOutPrice:       1000000,
		MinUsdValueLeft:     1000000,
	})
	s.Suite.Require().NoError(err)

	extendedVaultPairs, found := s.app.AssetKeeper.GetPairsVaults(s.ctx)
	s.Suite.Require().True(found)

	var extendedVaultPairID uint64
	for _, extendedVaultPair := range extendedVaultPairs {
		if extendedVaultPair.PairName == pairName && extendedVaultPair.AppId == appMappingID {
			extendedVaultPairID = extendedVaultPair.Id
			break
		}
	}
	s.Require().NotZero(extendedVaultPairID)
	return extendedVaultPairID
}

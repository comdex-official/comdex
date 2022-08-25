package keeper_test

import "github.com/comdex-official/comdex/x/lend/types"

func (s *KeeperTestSuite) Test() {

	//addr1 := s.addr(1)
	//addr2 := s.addr(2)

	assetOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	assetTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	assetThreeID := s.CreateNewAsset("ASSETHREE", "uasset3", 2000000)
	assetFourID := s.CreateNewAsset("ASSETFOUR", "uasset4", 2000000)
	cAssetOneID := s.CreateNewAsset("CASSET1", "ucasset1", 1000000)
	cAssetTwoID := s.CreateNewAsset("CASSET2", "ucasset2", 2000000)
	cAssetThreeID := s.CreateNewAsset("CASSET3", "ucasset3", 2000000)
	cAssetFourID := s.CreateNewAsset("CASSET4", "ucasset4", 2000000)

	var (
		assetDataPoolOne []types.AssetDataPoolMapping
		assetDataPoolTwo []types.AssetDataPoolMapping
	)
	assetDataPoolOneAssetOne := types.AssetDataPoolMapping{
		AssetID:   assetOneID,
		IsBridged: false,
	}
	assetDataPoolOneAssetTwo := types.AssetDataPoolMapping{
		AssetID:   assetTwoID,
		IsBridged: true,
	}
	assetDataPoolOneAssetThree := types.AssetDataPoolMapping{
		AssetID:   assetThreeID,
		IsBridged: true,
	}
	assetDataPoolTwoAssetFour := types.AssetDataPoolMapping{
		AssetID:   assetFourID,
		IsBridged: true,
	}

	assetDataPoolOne = append(assetDataPoolOne, assetDataPoolOneAssetOne, assetDataPoolOneAssetTwo, assetDataPoolOneAssetThree)
	assetDataPoolTwo = append(assetDataPoolOne, assetDataPoolTwoAssetFour, assetDataPoolOneAssetTwo, assetDataPoolOneAssetThree)

	poolOneID := s.CreateNewPool("cmdx", "CMDX-ATOM-CMST", assetOneID, assetTwoID, assetThreeID, assetDataPoolOne)
	poolTwoID := s.CreateNewPool("osmo", "OSMO-ATOM-CMST", assetFourID, assetTwoID, assetThreeID, assetDataPoolTwo)

	s.AddAssetRatesStats(assetThreeID, newDec("0.8"), newDec("0.002"), newDec("0.06"), newDec("0.6"), true, newDec("0.04"), newDec("0.04"), newDec("0.06"), newDec("0.8"), newDec("0.85"), newDec("0.025"), newDec("0.025"), newDec("0.1"), cAssetThreeID)
	s.AddAssetRatesStats(assetOneID, newDec("0.75"), newDec("0.002"), newDec("0.07"), newDec("1.25"), false, newDec("0.0"), newDec("0.0"), newDec("0.0"), newDec("0.7"), newDec("0.75"), newDec("0.05"), newDec("0.05"), newDec("0.2"), cAssetOneID)
	s.AddAssetRatesStats(assetFourID, newDec("0.65"), newDec("0.002"), newDec("0.08"), newDec("1.5"), false, newDec("0.0"), newDec("0.0"), newDec("0.0"), newDec("0.6"), newDec("0.65"), newDec("0.05"), newDec("0.05"), newDec("0.2"), cAssetFourID)
	s.AddAssetRatesStats(assetTwoID, newDec("0.5"), newDec("0.002"), newDec("0.08"), newDec("2.0"), false, newDec("0.0"), newDec("0.0"), newDec("0.0"), newDec("0.5"), newDec("0.55"), newDec("0.05"), newDec("0.05"), newDec("0.2"), cAssetTwoID)

	pairOneID := s.AddExtendedLendPair(assetTwoID, assetThreeID, false, poolOneID, 1000000)
	pairTwoID := s.AddExtendedLendPair(assetTwoID, assetOneID, false, poolOneID, 1000000)
	pairThreeID := s.AddExtendedLendPair(assetOneID, assetTwoID, false, poolOneID, 1000000)
	pairFourID := s.AddExtendedLendPair(assetOneID, assetThreeID, false, poolOneID, 1000000)
	pairFiveID := s.AddExtendedLendPair(assetThreeID, assetTwoID, false, poolOneID, 1000000)
	pairSixID := s.AddExtendedLendPair(assetThreeID, assetOneID, false, poolOneID, 1000000)
	pairSevenID := s.AddExtendedLendPair(assetFourID, assetThreeID, false, poolTwoID, 1000000)
	pairEightID := s.AddExtendedLendPair(assetFourID, assetOneID, false, poolTwoID, 1000000)
	pairNineID := s.AddExtendedLendPair(assetOneID, assetFourID, false, poolTwoID, 1000000)
	pairTenID := s.AddExtendedLendPair(assetOneID, assetThreeID, false, poolTwoID, 1000000)
	pairElevenID := s.AddExtendedLendPair(assetThreeID, assetFourID, false, poolTwoID, 1000000)
	pairTwelveID := s.AddExtendedLendPair(assetThreeID, assetOneID, false, poolTwoID, 1000000)
	pairThirteenID := s.AddExtendedLendPair(assetTwoID, assetFourID, true, poolTwoID, 1000000)
	pairFourteenID := s.AddExtendedLendPair(assetThreeID, assetFourID, true, poolTwoID, 1000000)
	pairFifteenID := s.AddExtendedLendPair(assetOneID, assetFourID, true, poolTwoID, 1000000)
	pairSixteenID := s.AddExtendedLendPair(assetFourID, assetTwoID, true, poolOneID, 1000000)
	pairSeventeenID := s.AddExtendedLendPair(assetThreeID, assetTwoID, true, poolOneID, 1000000)
	pairEighteenID := s.AddExtendedLendPair(assetOneID, assetTwoID, true, poolOneID, 1000000)

	s.AddAssetToPair(assetOneID, poolOneID, []uint64{pairThreeID, pairFourID, pairFifteenID})
	s.AddAssetToPair(assetTwoID, poolOneID, []uint64{pairOneID, pairTwoID, pairThirteenID})
	s.AddAssetToPair(assetThreeID, poolOneID, []uint64{pairFiveID, pairSixID, pairFourteenID})
	s.AddAssetToPair(assetFourID, poolTwoID, []uint64{pairSevenID, pairEightID, pairSixteenID})
	s.AddAssetToPair(assetOneID, poolTwoID, []uint64{pairNineID, pairTenID, pairEighteenID})
	s.AddAssetToPair(assetThreeID, poolTwoID, []uint64{pairElevenID, pairTwelveID, pairSeventeenID})

}

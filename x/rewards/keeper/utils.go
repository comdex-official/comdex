package keeper

// IntegerArrayToString converts integer slice to "," separated string.

// SplitTotalAmountPerEpoch splits amount into totalEpochs
// e.g. SplitTotalAmountPerEpoch(150, 11) => [13 13 13 13 14 14 14 14 14 14 14].
func SplitTotalAmountPerEpoch(totalAmount uint64, totalEpochs uint64) []uint64 {
	var splits []uint64
	if totalAmount < totalEpochs {
		return splits
	} else if totalAmount%totalEpochs == 0 {
		for i := uint64(0); i < totalEpochs; i++ {
			splits = append(splits, totalAmount/totalEpochs)
		}
		return splits
	} else {
		zp := totalEpochs - (totalAmount % totalEpochs)
		pp := totalAmount / totalEpochs
		for i := uint64(0); i < totalEpochs; i++ {
			if i >= zp {
				splits = append(splits, pp+1)
			} else {
				splits = append(splits, pp)
			}
		}
		return splits
	}
}

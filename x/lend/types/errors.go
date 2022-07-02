package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidAsset                    = sdkerrors.Register(ModuleName, 1100, "invalid asset")
	ErrLendingPoolInsufficient         = sdkerrors.Register(ModuleName, 1103, "lending pool insufficient")
	ErrInvalidRepayment                = sdkerrors.Register(ModuleName, 1104, "invalid repayment")
	ErrNegativeTimeElapsed             = sdkerrors.Register(ModuleName, 1111, "negative time elapsed since last interest time")
	ErrorUnknownProposalType           = sdkerrors.Register(ModuleName, 1119, "unknown proposal type")
	ErrorEmptyProposalAssets           = sdkerrors.Register(ModuleName, 1120, "empty assets in proposal")
	ErrorAssetDoesNotExist             = sdkerrors.Register(ModuleName, 1121, "asset does not exist")
	ErrorPairDoesNotExist              = sdkerrors.Register(ModuleName, 1123, "pair does not exist")
	ErrBadOfferCoinAmount              = sdkerrors.Register(ModuleName, 1125, "invalid offer coin amount")
	ErrorDuplicateLendPair             = sdkerrors.Register(ModuleName, 1126, "Duplicate lend Pair")
	ErrorDuplicateLend                 = sdkerrors.Register(ModuleName, 1127, "Duplicate lend Position")
	ErrorLendOwnerNotFound             = sdkerrors.Register(ModuleName, 1128, "Lend Owner not found")
	ErrLendNotFound                    = sdkerrors.Register(ModuleName, 1129, "Lend Position not found")
	ErrWithdrawalAmountExceeds         = sdkerrors.Register(ModuleName, 1130, "Withdrawal Amount Exceeds")
	ErrLendAccessUnauthorised          = sdkerrors.Register(ModuleName, 1131, "Unauthorized user for the tx")
	ErrorPairNotFound                  = sdkerrors.Register(ModuleName, 1132, "Pair Not Found")
	ErrorInvalidCollateralizationRatio = sdkerrors.Register(ModuleName, 1133, "Error Invalid Collaterallization Ratio")
	ErrorPriceInDoesNotExist           = sdkerrors.Register(ModuleName, 1134, "Error Price In Does Not Exist")
	ErrorPriceOutDoesNotExist          = sdkerrors.Register(ModuleName, 1135, "Error Price Out Does Not Exist")
	ErrorInvalidAmountIn               = sdkerrors.Register(ModuleName, 1136, "Error Invalid Amount In")
	ErrorInvalidAmountOut              = sdkerrors.Register(ModuleName, 1137, "Error Invalid Amount Out")
	ErrorDuplicateBorrow               = sdkerrors.Register(ModuleName, 1138, "Duplicate borrow Position")
	ErrBorrowingPoolInsufficient       = sdkerrors.Register(ModuleName, 1139, "borrowing pool insufficient")
	ErrBorrowNotFound                  = sdkerrors.Register(ModuleName, 1140, "Borrow Position not found")
	ErrBorrowingPositionOpen           = sdkerrors.Register(ModuleName, 1141, "borrowing position open")
	ErrAssetStatsNotFound              = sdkerrors.Register(ModuleName, 1142, "Asset Stats Not Found")
	ErrorDuplicateAssetRatesStats      = sdkerrors.Register(ModuleName, 1143, "Duplicate Asset Rates Stats")
	ErrorAssetStatsNotFound            = sdkerrors.Register(ModuleName, 1144, "Asset Stats Not Found")
	ErrInvalidAssetIDForPool           = sdkerrors.Register(ModuleName, 1145, "Asset Id not defined in the pool")
	ErrorAssetRatesStatsNotFound       = sdkerrors.Register(ModuleName, 1146, "Asset Rates Stats not found")
	ErrPoolNotFound                    = sdkerrors.Register(ModuleName, 1147, "Pool Not Found")
	ErrAvailableToBorrowInsufficient   = sdkerrors.Register(ModuleName, 1148, "Available To Borrow Insufficient")
	ErrStableBorrowDisabled            = sdkerrors.Register(ModuleName, 1149, "Stable Borrow Rate Not Enabled for This Asset")
)

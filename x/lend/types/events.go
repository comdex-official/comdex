package types

const (
	EventTypeLoanAsset           = "loan_asset"
	EventTypeWithdrawLoanedAsset = "withdraw_loaned_asset"
	EventTypeBorrowAsset         = "borrow_asset"
	EventTypeRepayBorrowedAsset  = "repay_borrowed_asset"

	EventAttrModule    = ModuleName
	EventAttrLender    = "lender"
	EventAttrBorrower  = "borrower"
	EventAttrAttempted = "attempted"
)

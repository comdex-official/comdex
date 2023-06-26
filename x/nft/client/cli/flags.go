package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagName            = "name"
	FlagDescription     = "description"
	FlagMediaURI        = "media-uri"
	FlagPreviewURI      = "preview-uri"
	FlagData            = "data"
	FlagNonTransferable = "non-transferable"
	FlagInExtensible    = "inextensible"
	FlagRecipient       = "recipient"
	FlagOwner           = "owner"
	FlagDenomID         = "denom-id"
	FlagSchema          = "schema"
	FlagNsfw            = "nsfw"
	FlagRoyaltyShare    = "royalty-share"
)

var (
	FsCreateDenom   = flag.NewFlagSet("", flag.ContinueOnError)
	FsUpdateDenom   = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferDenom = flag.NewFlagSet("", flag.ContinueOnError)
	FsMintNFT       = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferNFT   = flag.NewFlagSet("", flag.ContinueOnError)
	FsQuerySupply   = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryOwner    = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsCreateDenom.String(FlagSchema, "", "Denom schema")
	FsCreateDenom.String(FlagName, "", "Name of the denom")
	FsCreateDenom.String(FlagDescription, "", "Description for denom")
	FsCreateDenom.String(FlagPreviewURI, "", "Preview image uri for denom")

	FsUpdateDenom.String(FlagName, "[do-not-modify]", "Name of the denom")
	FsUpdateDenom.String(FlagDescription, "[do-not-modify]", "Description for denom")
	FsUpdateDenom.String(FlagPreviewURI, "[do-not-modify]", "Preview image uri for denom")

	FsTransferDenom.String(FlagRecipient, "", "recipient of the denom")

	FsMintNFT.String(FlagMediaURI, "", "Media uri of nft")
	FsMintNFT.String(FlagRecipient, "", "Receiver of the nft. default value is sender address of transaction")
	FsMintNFT.String(FlagPreviewURI, "", "Preview uri of nft")
	FsMintNFT.String(FlagName, "", "Name of nft")
	FsMintNFT.String(FlagDescription, "", "Description of nft")
	FsMintNFT.String(FlagData, "", "custom data of nft")

	FsMintNFT.Bool(FlagNonTransferable, false, "To mint non-transferable nft")
	FsMintNFT.Bool(FlagInExtensible, false, "To mint non-extensible nft")
	FsMintNFT.Bool(FlagNsfw, false, "not safe for work flag for nft")
	FsMintNFT.String(FlagRoyaltyShare, "", "Royalty share value decimal value between 0 and 1")

	FsTransferNFT.String(FlagRecipient, "", "Receiver of the nft. default value is sender address of transaction")
	FsQuerySupply.String(FlagOwner, "", "The owner of a nft")
	FsQueryOwner.String(FlagDenomID, "", "id of the denom")
}

package app

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	AccountPubKeyPrefix        = AccountAddressPrefix + sdk.PrefixPublic
	ValidatorAddressPrefix     = AccountAddressPrefix + sdk.PrefixValidator + sdk.PrefixOperator
	ValidatorPubKeyPrefix      = ValidatorAddressPrefix + sdk.PrefixPublic
	ConsensusNodeAddressPrefix = AccountAddressPrefix + sdk.PrefixValidator + sdk.PrefixConsensus
	ConsensusNodePubKeyPrefix  = ConsensusNodeAddressPrefix + sdk.PrefixPublic
)

// SetAccountAddressPrefixes sets the global prefix to be used when serializing addresses to bech32 strings.
func SetAccountAddressPrefixes() {
	config := sdk.GetConfig()

	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsensusNodeAddressPrefix, ConsensusNodePubKeyPrefix)

	config.Seal()
}

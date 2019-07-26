package types

import (
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSubmitProposal{}, "irishub/gov/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(MsgSubmitCommunityTaxUsageProposal{}, "irishub/gov/MsgSubmitCommunityTaxUsageProposal", nil)
	cdc.RegisterConcrete(MsgSubmitSoftwareUpgradeProposal{}, "irishub/gov/MsgSubmitSoftwareUpgradeProposal", nil)
	cdc.RegisterConcrete(MsgSubmitTokenAdditionProposal{}, "irishub/gov/MsgSubmitTokenAdditionProposal", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "irishub/gov/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgVote{}, "irishub/gov/MsgVote", nil)
	cdc.RegisterConcrete(&Vote{}, "irishub/gov/Vote", nil)
	cdc.RegisterConcrete(&GovParams{}, "irishub/gov/Params", nil)
}

var msgCdc = codec.New()

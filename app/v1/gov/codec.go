package gov

import (
	"github.com/irisnet/irishub/app/v1/gov/internal/keeper"
	"github.com/irisnet/irishub/app/v1/gov/internal/types"
	"github.com/irisnet/irishub/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
	cdc.RegisterInterface((*keeper.Proposal)(nil), nil)
	cdc.RegisterConcrete(&keeper.BasicProposal{}, "irishub/gov/BasicProposal", nil)
	cdc.RegisterConcrete(&keeper.ParameterProposal{}, "irishub/gov/ParameterProposal", nil)
	cdc.RegisterConcrete(&keeper.PlainTextProposal{}, "irishub/gov/PlainTextProposal", nil)
	cdc.RegisterConcrete(&keeper.TokenAdditionProposal{}, "irishub/gov/TokenAdditionProposal", nil)
	cdc.RegisterConcrete(&keeper.SoftwareUpgradeProposal{}, "irishub/gov/SoftwareUpgradeProposal", nil)
	cdc.RegisterConcrete(&keeper.SystemHaltProposal{}, "irishub/gov/SystemHaltProposal", nil)
	cdc.RegisterConcrete(&keeper.CommunityTaxUsageProposal{}, "irishub/gov/CommunityTaxUsageProposal", nil)
}

var msgCdc = codec.New()

package keeper

import (
	"fmt"
	"github.com/irisnet/irishub/app/v1/asset/exported"
	"github.com/irisnet/irishub/app/v1/gov/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

type handler func(content types.Content) Proposal

type ProposalFactory struct {
	router map[types.ProposalKind]handler
}

var proposalFactory ProposalFactory

func init() {
	proposalFactory = ProposalFactory{
		router: map[types.ProposalKind]handler{
			types.ProposalTypePlainText:         createPlainTextInfo(),
			types.ProposalTypeParameter:         createParameterInfo(),
			types.ProposalTypeSoftwareUpgrade:   createSoftwareUpgradeInfo(),
			types.ProposalTypeSystemHalt:        createSystemHaltInfo(),
			types.ProposalTypeCommunityTaxUsage: createCommunityTaxUsageInfo(),
			types.ProposalTypeTokenAddition:     createPlainTextInfo(),
		},
	}
}

func (factory ProposalFactory) create(content types.Content) Proposal {
	handler, ok := factory.router[content.GetProposalType()]
	if !ok {
		panic(fmt.Sprintf("not found proposal type:%d", content.GetProposalType()))
	}
	return handler(content)
}

func createPlainTextInfo() handler {
	return func(content types.Content) Proposal {
		return buildProposal(content, func(p BasicProposal, content types.Content) Proposal {
			return &PlainTextProposal{
				p,
			}
		})
	}
}
func createParameterInfo() handler {
	return func(content types.Content) Proposal {
		return buildProposal(content, func(p BasicProposal, content types.Content) Proposal {
			return &ParameterProposal{
				p,
				content.GetParams(),
			}
		})
	}
}
func createSoftwareUpgradeInfo() handler {
	return func(content types.Content) Proposal {
		return buildProposal(content, func(p BasicProposal, content types.Content) Proposal {
			upgradeMsg := content.(types.MsgSubmitSoftwareUpgradeProposal)
			proposal := &SoftwareUpgradeProposal{
				p,
				sdk.ProtocolDefinition{
					Version:   upgradeMsg.Version,
					Software:  upgradeMsg.Software,
					Height:    upgradeMsg.SwitchHeight,
					Threshold: upgradeMsg.Threshold},
			}
			return proposal
		})
	}
}

func createSystemHaltInfo() handler {
	return func(content types.Content) Proposal {
		return buildProposal(content, func(p BasicProposal, content types.Content) Proposal {
			return &SystemHaltProposal{
				p,
			}
		})
	}
}

func createCommunityTaxUsageInfo() handler {
	return func(content types.Content) Proposal {
		return buildProposal(content, func(p BasicProposal, content types.Content) Proposal {
			taxMsg := content.(types.MsgSubmitCommunityTaxUsageProposal)
			proposal := &CommunityTaxUsageProposal{
				p,
				TaxUsage{
					taxMsg.Usage,
					taxMsg.DestAddress,
					taxMsg.Percent},
			}
			return proposal
		})
	}
}

func createTokenAdditionInfo() handler {
	return func(content types.Content) Proposal {
		return buildProposal(content, func(p BasicProposal, content types.Content) Proposal {
			addTokenMsg := content.(types.MsgSubmitTokenAdditionProposal)
			decimal := int(addTokenMsg.Decimal)
			initialSupply := sdk.NewIntWithDecimal(int64(addTokenMsg.InitialSupply), decimal)
			maxSupply := sdk.NewIntWithDecimal(int64(exported.MaximumAssetMaxSupply), decimal)

			fToken := exported.NewFungibleToken(exported.EXTERNAL, "", addTokenMsg.Symbol, addTokenMsg.Name, addTokenMsg.Decimal, addTokenMsg.CanonicalSymbol, addTokenMsg.MinUnitAlias, initialSupply, maxSupply, false, nil)
			proposal := &TokenAdditionProposal{
				p,
				fToken,
			}
			return proposal
		})
	}
}

func buildProposal(content types.Content, callback func(p BasicProposal, content types.Content) Proposal) Proposal {
	var p = BasicProposal{
		Title:        content.GetTitle(),
		Description:  content.GetDescription(),
		ProposalType: content.GetProposalType(),
		Status:       types.StatusDepositPeriod,
		TallyResult:  types.EmptyTallyResult(),
		TotalDeposit: sdk.Coins{},
		Proposer:     content.GetProposer(),
	}
	return callback(p, content)
}

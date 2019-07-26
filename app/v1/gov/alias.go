package gov

import (
	"github.com/irisnet/irishub/app/v1/gov/internal/keeper"
	"github.com/irisnet/irishub/app/v1/gov/internal/types"
)

type (
	Keeper               = keeper.Keeper
	Proposal             = keeper.Proposal
	Proposals            = keeper.Proposals
	GovParams            = types.GovParams
	Params               = types.Params
	Param                = types.Param
	Vote                 = types.Vote
	Votes                = types.Votes
	Deposit              = types.Deposit
	Deposits             = types.Deposits
	TallyResult          = types.TallyResult
	QueryProposalParams  = keeper.QueryProposalParams
	QueryProposalsParams = keeper.QueryProposalsParams
	QueryVoteParams      = keeper.QueryVoteParams
	QueryVotesParams     = keeper.QueryVotesParams
	QueryDepositParams   = keeper.QueryDepositParams
	QueryDepositsParams  = keeper.QueryDepositsParams
	QueryTallyParams     = keeper.QueryTallyParams
	UsageType            = types.UsageType
)

const (
	PASS       = types.PASS
	REJECT     = types.REJECT
	REJECTVETO = types.REJECTVETO

	StatusPassed   = types.StatusPassed
	StatusRejected = types.StatusRejected

	DefaultParamSpace             = types.DefaultParamSpace
	StatusDepositPeriod           = types.StatusDepositPeriod
	StatusVotingPeriod            = types.StatusVotingPeriod
	ProposalTypeParameter         = types.ProposalTypeParameter
	ProposalTypeCommunityTaxUsage = types.ProposalTypeCommunityTaxUsage
	ProposalTypeSoftwareUpgrade   = types.ProposalTypeSoftwareUpgrade
	ProposalTypeTokenAddition     = types.ProposalTypeTokenAddition
)

var (
	NewKeeper                             = keeper.NewKeeper
	NewQuerier                            = keeper.NewQuerier
	DefaultCodespace                      = types.DefaultCodespace
	PrometheusMetrics                     = types.PrometheusMetrics
	ValidateParams                        = types.ValidateParams
	ProposalStatusFromString              = types.ProposalStatusFromString
	ProposalTypeFromString                = types.ProposalTypeFromString
	NewMsgSubmitProposal                  = types.NewMsgSubmitProposal
	UsageTypeFromString                   = types.UsageTypeFromString
	NewMsgSubmitCommunityTaxUsageProposal = types.NewMsgSubmitCommunityTaxUsageProposal
	NewMsgSubmitSoftwareUpgradeProposal   = types.NewMsgSubmitSoftwareUpgradeProposal
	NewMsgSubmitTokenAdditionProposal     = types.NewMsgSubmitTokenAdditionProposal
	NewMsgDeposit                         = types.NewMsgDeposit
	NewMsgVote                            = types.NewMsgVote
	VoteOptionFromString                  = types.VoteOptionFromString
	ErrInvalidParam                       = types.ErrInvalidParam
)

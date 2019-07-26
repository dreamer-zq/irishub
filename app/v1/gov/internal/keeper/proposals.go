package keeper

import (
	"fmt"
	"github.com/irisnet/irishub/app/v1/gov/internal/types"
	"strings"
	"time"

	sdk "github.com/irisnet/irishub/types"
	"github.com/pkg/errors"
)

//-----------------------------------------------------------
// Proposal interface
type Proposal interface {
	GetProposalID() uint64
	SetProposalID(uint64)

	GetTitle() string
	SetTitle(string)

	GetDescription() string
	SetDescription(string)

	GetProposalType() types.ProposalKind
	SetProposalType(types.ProposalKind)

	GetStatus() types.ProposalStatus
	SetStatus(types.ProposalStatus)

	GetTallyResult() types.TallyResult
	SetTallyResult(types.TallyResult)

	GetSubmitTime() time.Time
	SetSubmitTime(time.Time)

	GetDepositEndTime() time.Time
	SetDepositEndTime(time.Time)

	GetTotalDeposit() sdk.Coins
	SetTotalDeposit(sdk.Coins)

	GetVotingStartTime() time.Time
	SetVotingStartTime(time.Time)

	GetVotingEndTime() time.Time
	SetVotingEndTime(time.Time)

	GetProposalLevel() types.ProposalLevel
	GetProposer() sdk.AccAddress

	String() string
	Validate(ctx sdk.Context, gk Keeper, verifyPropNum bool) sdk.Error
	Execute(ctx sdk.Context, gk Keeper) sdk.Error
}

//-----------------------------------------------------------
// Basic Proposals
type BasicProposal struct {
	ProposalID   uint64             `json:"proposal_id"`   //  ID of the proposal
	Title        string             `json:"title"`         //  Title of the proposal
	Description  string             `json:"description"`   //  Description of the proposal
	ProposalType types.ProposalKind `json:"proposal_type"` //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}

	Status      types.ProposalStatus `json:"proposal_status"` //  Status of the Proposal {Pending, Active, Passed, Rejected}
	TallyResult types.TallyResult    `json:"tally_result"`    //  Result of Tallys

	SubmitTime     time.Time `json:"submit_time"`      //  Time of the block where TxGovSubmitProposal was included
	DepositEndTime time.Time `json:"deposit_end_time"` // Time that the Proposal would expire if deposit amount isn't met
	TotalDeposit   sdk.Coins `json:"total_deposit"`    //  Current deposit on this proposal. Initial value is set at InitialDeposit

	VotingStartTime time.Time      `json:"voting_start_time"` //  Time of the block where MinDeposit was reached. -1 if MinDeposit is not reached
	VotingEndTime   time.Time      `json:"voting_end_time"`   // Time that the VotingPeriod for this proposal will end and votes will be tallied
	Proposer        sdk.AccAddress `json:"proposer"`
}

func (bp BasicProposal) String() string {
	return fmt.Sprintf(`Proposal %d:
  Title:              %s
  Type:               %s
  TagProposer:           %s
  Status:             %s
  Submit Time:        %s
  Deposit End Time:   %s
  Total Deposit:      %s
  Voting Start Time:  %s
  Voting End Time:    %s
  Description:        %s`,
		bp.ProposalID, bp.Title, bp.ProposalType, bp.Proposer.String(),
		bp.Status, bp.SubmitTime, bp.DepositEndTime,
		bp.TotalDeposit.MainUnitString(), bp.VotingStartTime, bp.VotingEndTime, bp.GetDescription(),
	)
}

func (bp BasicProposal) HumanString(converter sdk.CoinsConverter) string {
	return fmt.Sprintf(`Proposal %d:
  Title:              %s
  Type:               %s
  Status:             %s
  Submit Time:        %s
  Deposit End Time:   %s
  Total Deposit:      %s
  Voting Start Time:  %s
  Voting End Time:    %s
  Description:        %s`,
		bp.ProposalID, bp.Title, bp.ProposalType,
		bp.Status, bp.SubmitTime, bp.DepositEndTime,
		converter.ToMainUnit(bp.TotalDeposit), bp.VotingStartTime, bp.VotingEndTime, bp.GetDescription(),
	)
}

// Proposals is an array of proposal
type Proposals []Proposal

// nolint
func (p Proposals) String() string {
	if len(p) == 0 {
		return "[]"
	}
	out := "ID - (Status) [Type] [TotalDeposit] Title\n"
	for _, prop := range p {
		out += fmt.Sprintf("%d - (%s) [%s] [%s] %s\n",
			prop.GetProposalID(), prop.GetStatus(),
			prop.GetProposalType(), prop.GetTotalDeposit().MainUnitString(), prop.GetTitle())
	}
	return strings.TrimSpace(out)
}

func (p Proposals) HumanString(converter sdk.CoinsConverter) string {
	if len(p) == 0 {
		return "[]"
	}
	out := "ID - (Status) [Type] [TotalDeposit] Title\n"
	for _, prop := range p {
		out += fmt.Sprintf("%d - (%s) [%s] [%s] %s\n",
			prop.GetProposalID(), prop.GetStatus(),
			prop.GetProposalType(), converter.ToMainUnit(prop.GetTotalDeposit()), prop.GetTitle())
	}
	return strings.TrimSpace(out)
}

// Implements Proposal Interface
var _ Proposal = (*BasicProposal)(nil)

// nolint
func (bp BasicProposal) GetProposalID() uint64               { return bp.ProposalID }
func (bp *BasicProposal) SetProposalID(proposalID uint64)    { bp.ProposalID = proposalID }
func (bp BasicProposal) GetTitle() string                    { return bp.Title }
func (bp *BasicProposal) SetTitle(title string)              { bp.Title = title }
func (bp BasicProposal) GetDescription() string              { return bp.Description }
func (bp *BasicProposal) SetDescription(description string)  { bp.Description = description }
func (bp BasicProposal) GetProposalType() types.ProposalKind { return bp.ProposalType }
func (bp *BasicProposal) SetProposalType(proposalType types.ProposalKind) {
	bp.ProposalType = proposalType
}
func (bp BasicProposal) GetStatus() types.ProposalStatus               { return bp.Status }
func (bp *BasicProposal) SetStatus(status types.ProposalStatus)        { bp.Status = status }
func (bp BasicProposal) GetTallyResult() types.TallyResult             { return bp.TallyResult }
func (bp *BasicProposal) SetTallyResult(tallyResult types.TallyResult) { bp.TallyResult = tallyResult }
func (bp BasicProposal) GetSubmitTime() time.Time                      { return bp.SubmitTime }
func (bp *BasicProposal) SetSubmitTime(submitTime time.Time)           { bp.SubmitTime = submitTime }
func (bp BasicProposal) GetDepositEndTime() time.Time                  { return bp.DepositEndTime }
func (bp *BasicProposal) SetDepositEndTime(depositEndTime time.Time) {
	bp.DepositEndTime = depositEndTime
}
func (bp BasicProposal) GetTotalDeposit() sdk.Coins              { return bp.TotalDeposit }
func (bp *BasicProposal) SetTotalDeposit(totalDeposit sdk.Coins) { bp.TotalDeposit = totalDeposit }
func (bp BasicProposal) GetVotingStartTime() time.Time           { return bp.VotingStartTime }
func (bp *BasicProposal) SetVotingStartTime(votingStartTime time.Time) {
	bp.VotingStartTime = votingStartTime
}
func (bp BasicProposal) GetVotingEndTime() time.Time { return bp.VotingEndTime }
func (bp *BasicProposal) SetVotingEndTime(votingEndTime time.Time) {
	bp.VotingEndTime = votingEndTime
}
func (bp BasicProposal) GetProtocolDefinition() sdk.ProtocolDefinition {
	return sdk.ProtocolDefinition{}
}
func (bp *BasicProposal) SetProtocolDefinition(sdk.ProtocolDefinition) {}
func (bp BasicProposal) GetTaxUsage() TaxUsage                         { return TaxUsage{} }
func (bp *BasicProposal) SetTaxUsage(taxUsage TaxUsage)                {}
func (bp *BasicProposal) Validate(ctx sdk.Context, k Keeper, verify bool) sdk.Error {
	if !verify {
		return nil
	}
	pLevel := bp.ProposalType.GetProposalLevel()
	if num, ok := k.HasReachedTheMaxProposalNum(ctx, pLevel); ok {
		return types.ErrMoreThanMaxProposal(k.codespace, num, pLevel.String())
	}
	return nil
}
func (bp *BasicProposal) GetProposalLevel() types.ProposalLevel {
	return bp.ProposalType.GetProposalLevel()
}

func (bp *BasicProposal) GetProposer() sdk.AccAddress {
	return bp.Proposer
}
func (bp *BasicProposal) Execute(ctx sdk.Context, gk Keeper) sdk.Error {
	return sdk.MarshalResultErr(errors.New("BasicProposal can not execute 'Execute' method"))
}

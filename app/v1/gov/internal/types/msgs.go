package types

import (
	"fmt"
	"github.com/irisnet/irishub/app/v1/asset/exported"
	sdk "github.com/irisnet/irishub/types"
)

// name to idetify transaction types
const MsgRoute = "gov"

var _, _, _, _ sdk.Msg = MsgSubmitProposal{}, MsgSubmitCommunityTaxUsageProposal{}, MsgDeposit{}, MsgVote{}

type Content interface {
	sdk.Msg
	GetTitle() string
	GetDescription() string
	GetProposalType() ProposalKind
	GetProposer() sdk.AccAddress
	GetInitialDeposit() sdk.Coins
	GetParams() Params
}

//-----------------------------------------------------------
// MsgSubmitProposal
type MsgSubmitProposal struct {
	Title          string         `json:"title"`           //  Title of the proposal
	Description    string         `json:"description"`     //  Description of the proposal
	ProposalType   ProposalKind   `json:"proposal_type"`   //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	Proposer       sdk.AccAddress `json:"proposer"`        //  Address of the proposer
	InitialDeposit sdk.Coins      `json:"initial_deposit"` //  Initial deposit paid by sender. Must be strictly positive.
	Params         Params         `json:"params"`
}

func NewMsgSubmitProposal(title string, description string, proposalType ProposalKind, proposer sdk.AccAddress, initialDeposit sdk.Coins, params Params) MsgSubmitProposal {
	return MsgSubmitProposal{
		Title:          title,
		Description:    description,
		ProposalType:   proposalType,
		Proposer:       proposer,
		InitialDeposit: initialDeposit,
		Params:         params,
	}
}

func (msg MsgSubmitProposal) GetTitle() string {
	return msg.Title
}
func (msg MsgSubmitProposal) GetDescription() string {
	return msg.Description
}
func (msg MsgSubmitProposal) GetProposalType() ProposalKind {
	return msg.ProposalType
}
func (msg MsgSubmitProposal) GetProposer() sdk.AccAddress {
	return msg.Proposer
}
func (msg MsgSubmitProposal) GetInitialDeposit() sdk.Coins {
	return msg.InitialDeposit
}
func (msg MsgSubmitProposal) GetParams() Params {
	return msg.Params
}

//nolint
func (msg MsgSubmitProposal) Route() string { return MsgRoute }
func (msg MsgSubmitProposal) Type() string  { return "submit_proposal" }

// Implements Msg.
func (msg MsgSubmitProposal) ValidateBasic() sdk.Error {
	if len(msg.Title) == 0 {
		return ErrInvalidTitle(DefaultCodespace, msg.Title) // TODO: Proper Error
	}
	if len(msg.Description) == 0 {
		return ErrInvalidDescription(DefaultCodespace, msg.Description) // TODO: Proper Error
	}
	if !ValidProposalType(msg.ProposalType) {
		return ErrInvalidProposalType(DefaultCodespace, msg.ProposalType)
	}
	if len(msg.Proposer) == 0 {
		return sdk.ErrInvalidAddress(msg.Proposer.String())
	}
	if !msg.InitialDeposit.IsValid() {
		return sdk.ErrInvalidCoins(msg.InitialDeposit.String())
	}
	if !msg.InitialDeposit.IsNotNegative() {
		return sdk.ErrInvalidCoins(msg.InitialDeposit.String())
	}
	if err := msg.EnsureLength(); err != nil {
		return err
	}
	if msg.ProposalType == ProposalTypeParameter {
		if len(msg.Params) == 0 {
			return ErrEmptyParam(DefaultCodespace)
		}
		if len(msg.Params) > 1 {
			return ErrInvalidParamNum(DefaultCodespace)
		}
	}
	return nil
}

func (msg MsgSubmitProposal) String() string {
	return fmt.Sprintf("MsgSubmitProposal{%s, %s, %s, %v}", msg.Title, msg.Description, msg.ProposalType, msg.InitialDeposit)
}

// Implements Msg.
func (msg MsgSubmitProposal) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg MsgSubmitProposal) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgSubmitProposal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proposer}
}

type MsgSubmitSoftwareUpgradeProposal struct {
	MsgSubmitProposal
	Version      uint64  `json:"version"`
	Software     string  `json:"software"`
	SwitchHeight uint64  `json:"switch_height"`
	Threshold    sdk.Dec `json:"threshold"`
}

func NewMsgSubmitSoftwareUpgradeProposal(msgSubmitProposal MsgSubmitProposal, version uint64, software string, switchHeight uint64, threshold sdk.Dec) MsgSubmitSoftwareUpgradeProposal {
	return MsgSubmitSoftwareUpgradeProposal{
		MsgSubmitProposal: msgSubmitProposal,
		Version:           version,
		Software:          software,
		SwitchHeight:      switchHeight,
		Threshold:         threshold,
	}
}

func (msg MsgSubmitSoftwareUpgradeProposal) ValidateBasic() sdk.Error {
	err := msg.MsgSubmitProposal.ValidateBasic()
	if err != nil {
		return err
	}

	if len(msg.Software) > 70 {
		return sdk.ErrInvalidLength(DefaultCodespace, CodeInvalidProposal, "software", len(msg.Software), 70)
	}

	// if threshold not in [0.85,1), then print error
	if msg.Threshold.LT(sdk.NewDecWithPrec(80, 2)) || msg.Threshold.GTE(sdk.NewDec(1)) {
		return ErrInvalidUpgradeThreshold(DefaultCodespace, msg.Threshold)
	}

	return nil
}

func (msg MsgSubmitSoftwareUpgradeProposal) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

type MsgSubmitCommunityTaxUsageProposal struct {
	MsgSubmitProposal
	Usage       UsageType      `json:"usage"`
	DestAddress sdk.AccAddress `json:"dest_address"`
	Percent     sdk.Dec        `json:"percent"`
}

func NewMsgSubmitCommunityTaxUsageProposal(msgSubmitProposal MsgSubmitProposal, usage UsageType, destAddress sdk.AccAddress, percent sdk.Dec) MsgSubmitCommunityTaxUsageProposal {
	return MsgSubmitCommunityTaxUsageProposal{
		MsgSubmitProposal: msgSubmitProposal,
		Usage:             usage,
		DestAddress:       destAddress,
		Percent:           percent,
	}
}

func (msg MsgSubmitCommunityTaxUsageProposal) ValidateBasic() sdk.Error {
	err := msg.MsgSubmitProposal.ValidateBasic()
	if err != nil {
		return err
	}
	if !ValidUsageType(msg.Usage) {
		return ErrInvalidUsageType(DefaultCodespace, msg.Usage)
	}
	if msg.Usage != UsageTypeBurn && len(msg.DestAddress) == 0 {
		return sdk.ErrInvalidAddress(msg.DestAddress.String())
	}
	if msg.Percent.LTE(sdk.NewDec(0)) || msg.Percent.GT(sdk.NewDec(1)) {
		return ErrInvalidPercent(DefaultCodespace, msg.Percent)
	}
	return nil
}

func (msg MsgSubmitCommunityTaxUsageProposal) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

//-----------------------------------------------------------
// MsgDeposit
type MsgDeposit struct {
	ProposalID uint64         `json:"proposal_id"` // ID of the proposal
	Depositor  sdk.AccAddress `json:"depositor"`   // Address of the depositor
	Amount     sdk.Coins      `json:"amount"`      // Coins to add to the proposal's deposit
}

func NewMsgDeposit(depositor sdk.AccAddress, proposalID uint64, amount sdk.Coins) MsgDeposit {
	return MsgDeposit{
		ProposalID: proposalID,
		Depositor:  depositor,
		Amount:     amount,
	}
}

// Implements Msg.
// nolint
func (msg MsgDeposit) Route() string { return MsgRoute }
func (msg MsgDeposit) Type() string  { return "deposit" }

// Implements Msg.
func (msg MsgDeposit) ValidateBasic() sdk.Error {
	if len(msg.Depositor) == 0 {
		return sdk.ErrInvalidAddress(msg.Depositor.String())
	}
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}
	if !msg.Amount.IsNotNegative() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}
	if msg.ProposalID < 0 {
		return ErrUnknownProposal(DefaultCodespace, msg.ProposalID)
	}
	return nil
}

func (msg MsgDeposit) String() string {
	return fmt.Sprintf("MsgDeposit{%s=>%v: %v}", msg.Depositor, msg.ProposalID, msg.Amount)
}

// Implements Msg.
func (msg MsgDeposit) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg MsgDeposit) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Depositor}
}

//-----------------------------------------------------------
// MsgVote
type MsgVote struct {
	ProposalID uint64         `json:"proposal_id"` // ID of the proposal
	Voter      sdk.AccAddress `json:"voter"`       //  address of the voter
	Option     VoteOption     `json:"option"`      //  option from OptionSet chosen by the voter
}

func NewMsgVote(voter sdk.AccAddress, proposalID uint64, option VoteOption) MsgVote {
	return MsgVote{
		ProposalID: proposalID,
		Voter:      voter,
		Option:     option,
	}
}

// Implements Msg.
// nolint
func (msg MsgVote) Route() string { return MsgRoute }
func (msg MsgVote) Type() string  { return "vote" }

// Implements Msg.
func (msg MsgVote) ValidateBasic() sdk.Error {
	if len(msg.Voter.Bytes()) == 0 {
		return sdk.ErrInvalidAddress(msg.Voter.String())
	}
	if msg.ProposalID < 0 {
		return ErrUnknownProposal(DefaultCodespace, msg.ProposalID)
	}
	if !ValidVoteOption(msg.Option) {
		return ErrInvalidVote(DefaultCodespace, msg.Option)
	}
	return nil
}

func (msg MsgVote) String() string {
	return fmt.Sprintf("MsgVote{%v - %s}", msg.ProposalID, msg.Option)
}

// Implements Msg.
func (msg MsgVote) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg MsgVote) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgVote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Voter}
}

func (msg MsgSubmitProposal) EnsureLength() sdk.Error {
	if len(msg.Title) > 70 {
		return sdk.ErrInvalidLength(DefaultCodespace, CodeInvalidProposal, "title", len(msg.Title), 70)
	}
	if len(msg.Description) > 280 {
		return sdk.ErrInvalidLength(DefaultCodespace, CodeInvalidProposal, "description", len(msg.Description), 280)
	}

	return nil
}

type MsgSubmitTokenAdditionProposal struct {
	MsgSubmitProposal
	Symbol          string `json:"symbol"`
	CanonicalSymbol string `json:"canonical_symbol"`
	Name            string `json:"name"`
	Decimal         uint8  `json:"decimal"`
	MinUnitAlias    string `json:"min_unit_alias"`
	InitialSupply   uint64 `json:"initial_supply"`
}

func NewMsgSubmitTokenAdditionProposal(msgSubmitProposal MsgSubmitProposal, symbol, canonicalSymbol, name, minUnitAlias string, decimal uint8, initialSupply uint64) MsgSubmitTokenAdditionProposal {
	return MsgSubmitTokenAdditionProposal{
		MsgSubmitProposal: msgSubmitProposal,
		Symbol:            symbol,
		CanonicalSymbol:   canonicalSymbol,
		Name:              name,
		Decimal:           decimal,
		MinUnitAlias:      minUnitAlias,
		InitialSupply:     initialSupply,
	}
}

func (msg MsgSubmitTokenAdditionProposal) ValidateBasic() sdk.Error {
	err := msg.MsgSubmitProposal.ValidateBasic()
	if err != nil {
		return err
	}

	issueToken := exported.NewMsgIssueToken(exported.FUNGIBLE, exported.EXTERNAL, "", msg.Symbol, msg.CanonicalSymbol, msg.Name, msg.Decimal, msg.MinUnitAlias, msg.InitialSupply, exported.MaximumAssetMaxSupply, false, nil)
	return exported.ValidateMsgIssueToken(&issueToken)
}
func (msg MsgSubmitTokenAdditionProposal) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Vote
type Vote struct {
	Voter      sdk.AccAddress `json:"voter"`       //  address of the voter
	ProposalID uint64         `json:"proposal_id"` //  proposalID of the proposal
	Option     VoteOption     `json:"option"`      //  option from OptionSet chosen by the voter
}

func (v Vote) String() string {
	return fmt.Sprintf("TagVoter %s voted with option %s on proposal %d", v.Voter, v.Option, v.ProposalID)
}

// Votes is a collection of Vote
type Votes []Vote

func (v Votes) String() string {
	out := fmt.Sprintf("Votes for Proposal %d:", v[0].ProposalID)
	for _, vot := range v {
		out += fmt.Sprintf("\n  %s: %s", vot.Voter, vot.Option)
	}
	return out
}

// Returns whether 2 votes are equal
func (v Vote) Equals(voteB Vote) bool {
	return v.Voter.Equals(voteB.Voter) && v.ProposalID == voteB.ProposalID && v.Option == voteB.Option
}

// Returns whether a vote is empty
func (v Vote) Empty() bool {
	voteB := Vote{}
	return v.Equals(voteB)
}

// Deposit
type Deposit struct {
	Depositor  sdk.AccAddress `json:"depositor"`   //  Address of the depositor
	ProposalID uint64         `json:"proposal_id"` //  proposalID of the proposal
	Amount     sdk.Coins      `json:"amount"`      //  Deposit amount
}

func (d Deposit) String() string {
	return fmt.Sprintf("Deposit by %s on Proposal %d is for the amount %s",
		d.Depositor, d.ProposalID, d.Amount.MainUnitString())
}

func (d Deposit) HumanString(converter sdk.CoinsConverter) string {
	return fmt.Sprintf("Deposit by %s on Proposal %d is for the amount %s",
		d.Depositor, d.ProposalID, converter.ToMainUnit(d.Amount))
}

// Deposits is a collection of depoist
type Deposits []Deposit

func (d Deposits) String() string {
	if len(d) == 0 {
		return "[]"
	}
	out := fmt.Sprintf("Deposits for Proposal %d:", d[0].ProposalID)
	for _, dep := range d {
		out += fmt.Sprintf("\n  %s: %s", dep.Depositor, dep.Amount.MainUnitString())
	}
	return out
}

func (d Deposits) HumanString(converter sdk.CoinsConverter) string {
	if len(d) == 0 {
		return "[]"
	}
	out := fmt.Sprintf("Deposits for Proposal %d:", d[0].ProposalID)
	for _, dep := range d {
		out += fmt.Sprintf("\n  %s: %s", dep.Depositor, converter.ToMainUnit(dep.Amount))
	}
	return out
}

// Returns whether 2 deposits are equal
func (d Deposit) Equals(depositB Deposit) bool {
	return d.Depositor.Equals(depositB.Depositor) && d.ProposalID == depositB.ProposalID && d.Amount.IsEqual(depositB.Amount)
}

// Returns whether a deposit is empty
func (d Deposit) Empty() bool {
	depositB := Deposit{}
	return d.Equals(depositB)
}

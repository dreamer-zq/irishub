// nolint
package types

import (
	sdk "github.com/irisnet/irishub/types"
)

var (
	ActionSubmitProposal   = []byte("submit-proposal")
	ActionDeposit          = []byte("deposit")
	ActionVote             = []byte("vote")
	ActionProposalDropped  = []byte("proposal-dropped")
	ActionProposalPassed   = []byte("proposal-passed")
	ActionProposalRejected = []byte("proposal-rejected")

	TagAction            = sdk.TagAction
	TagProposer          = "proposer"
	TagProposalID        = "proposal-id"
	TagVotingPeriodStart = "voting-period-start"
	TagDepositor         = "depositor"
	TagVoter             = "voter"
	TagParam             = "param"
	TagUsage             = "usage"
	TagPercent           = "percent"
	TagDestAddress       = "dest-address"
	TagTokenId           = "token-id"
)

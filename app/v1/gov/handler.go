package gov

import (
	"github.com/irisnet/irishub/app/v1/gov/internal/keeper"
	"github.com/irisnet/irishub/app/v1/gov/internal/types"
	"strconv"

	sdk "github.com/irisnet/irishub/types"
)

// Handle all "gov" type messages.
func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgSubmitProposal,
			types.MsgSubmitSoftwareUpgradeProposal,
			types.MsgSubmitTokenAdditionProposal,
			types.MsgSubmitCommunityTaxUsageProposal:
			return handleMsgSubmitProposal(ctx, keeper, msg)
		case types.MsgDeposit:
			return handleMsgDeposit(ctx, keeper, msg)
		case types.MsgVote:
			return handleMsgVote(ctx, keeper, msg)
		default:
			errMsg := "Unrecognized gov msg type"
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSubmitProposal(ctx sdk.Context, keeper keeper.Keeper, msg sdk.Msg) sdk.Result {
	resTags, err := keeper.SubmitProposal(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: resTags,
	}
}

func handleMsgDeposit(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgDeposit) sdk.Result {

	err, votingStarted := keeper.AddDeposit(ctx, msg.ProposalID, msg.Depositor, msg.Amount)
	if err != nil {
		return err.Result()
	}

	proposalIDBytes := []byte(strconv.FormatUint(msg.ProposalID, 10))

	// TODO: Add tag for if voting period started
	resTags := sdk.NewTags(
		types.TagDepositor, []byte(msg.Depositor.String()),
		types.TagProposalID, proposalIDBytes,
	)

	if votingStarted {
		resTags = resTags.AppendTag(types.TagVotingPeriodStart, proposalIDBytes)
	}

	return sdk.Result{
		Tags: resTags,
	}
}

func handleMsgVote(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgVote) sdk.Result {

	err := keeper.AddVote(ctx, msg.ProposalID, msg.Voter, msg.Option)
	if err != nil {
		return err.Result()
	}

	proposalIDBytes := []byte(strconv.FormatUint(msg.ProposalID, 10))

	resTags := sdk.NewTags(
		types.TagVoter, []byte(msg.Voter.String()),
		types.TagProposalID, proposalIDBytes,
	)
	return sdk.Result{
		Tags: resTags,
	}
}

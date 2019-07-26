package gov

import (
	"github.com/irisnet/irishub/app/v1/gov/internal/types"
	"strconv"

	sdk "github.com/irisnet/irishub/types"
	tmstate "github.com/tendermint/tendermint/state"
)

// Called every block, process inflation, update validator set
func EndBlocker(ctx sdk.Context, keeper Keeper) (resTags sdk.Tags) {
	ctx = ctx.WithCoinFlowTrigger(sdk.GovEndBlocker)
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "endBlock").With("module", "iris/gov"))
	resTags = sdk.NewTags()

	if ctx.BlockHeight() == keeper.GetSystemHaltHeight(ctx) {
		resTags = resTags.AppendTag(tmstate.HaltTagKey, []byte(tmstate.HaltTagValue))
		ctx.Logger().Info("SystemHalt Start!!!")
	}

	//pop inactive proposal from inactive queue
	keeper.PopInactiveProposal(ctx, ctx.BlockHeader().Time, func(p Proposal) {
		resTags = resTags.AppendTag(types.TagAction, types.ActionProposalDropped)
		resTags = resTags.AppendTag(types.TagProposalID, []byte(string(p.GetProposalID())))
		ctx.Logger().Info("Proposal didn't meet minimum deposit; deleted", "TagProposalID",
			p.GetProposalID(), "MinDeposit", keeper.GetDepositProcedure(ctx, p.GetProposalLevel()).MinDeposit,
			"ActualDeposit", p.GetTotalDeposit(),
		)
		keeper.DeleteDeposits(ctx, p.GetProposalID())
		keeper.DeleteProposal(ctx, p.GetProposalID())
		keeper.SubProposalNum(ctx, p.GetProposalLevel())
	})

	//pop active proposal from active queue
	keeper.PopActiveProposal(ctx, ctx.BlockHeader().Time, func(p Proposal) {
		result, tallyResults, votingVals := keeper.Tally(ctx, p)
		var action []byte
		switch result {
		case PASS:
			keeper.Metrics.ProposalStatus.With(types.ProposalIDLabel, strconv.FormatUint(p.GetProposalID(), 10)).Set(2)
			keeper.RefundDeposits(ctx, p.GetProposalID())
			p.SetStatus(StatusPassed)
			action = types.ActionProposalPassed
			p.Execute(ctx, keeper)
			break
		case REJECT:
			keeper.Metrics.ProposalStatus.With(types.ProposalIDLabel, strconv.FormatUint(p.GetProposalID(), 10)).Set(3)
			keeper.RefundDeposits(ctx, p.GetProposalID())
			p.SetStatus(StatusRejected)
			action = types.ActionProposalRejected
			break
		case REJECTVETO:
			keeper.Metrics.ProposalStatus.With(types.ProposalIDLabel, strconv.FormatUint(p.GetProposalID(), 10)).Set(3)
			keeper.DeleteDeposits(ctx, p.GetProposalID())
			p.SetStatus(StatusRejected)
			action = types.ActionProposalRejected
		}

		p.SetTallyResult(tallyResults)
		keeper.SetProposal(ctx, p)

		keeper.Slash(ctx, p, votingVals)
		keeper.SubProposalNum(ctx, p.GetProposalLevel())
		keeper.DeleteValidatorSet(ctx, p.GetProposalID())

		resTags = resTags.AppendTag(types.TagAction, action)
		resTags = resTags.AppendTag(types.TagProposalID, []byte(string(p.GetProposalID())))

		ctx.Logger().Info("Proposal tallied", "TagProposalID", p.GetProposalID(), "result", result, "tallyResults", tallyResults)
	})
	return resTags
}

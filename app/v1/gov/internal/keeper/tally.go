package keeper

import (
	"github.com/irisnet/irishub/app/v1/gov/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

// validatorGovInfo used for tallying
type validatorGovInfo struct {
	Address             sdk.ValAddress   // address of the validator operator
	Vote                types.VoteOption // Vote of the validator
	TokenPerShare       sdk.Dec
	DelegatorShares     sdk.Dec // Total outstanding delegator shares
	DelegatorDeductions sdk.Dec // Delegator deductions from validator's delegators voting independently
}

func (keeper Keeper) Tally(ctx sdk.Context, proposal Proposal) (result types.ProposalResult, tallyResults types.TallyResult, votingVals map[string]bool) {
	results := make(map[types.VoteOption]sdk.Dec)
	results[types.OptionYes] = sdk.ZeroDec()
	results[types.OptionAbstain] = sdk.ZeroDec()
	results[types.OptionNo] = sdk.ZeroDec()
	results[types.OptionNoWithVeto] = sdk.ZeroDec()

	//voted votingPower
	totalVotingPower := sdk.ZeroDec()
	//all votingPower
	systemVotingPower := sdk.ZeroDec()
	currValidators := make(map[string]validatorGovInfo)
	votingVals = make(map[string]bool)

	keeper.vs.IterateBondedValidatorsByPower(ctx, func(index int64, validator sdk.Validator) (stop bool) {
		currValidators[validator.GetOperator().String()] = validatorGovInfo{
			Address:             validator.GetOperator(),
			TokenPerShare:       validator.GetTokens().Quo(validator.GetDelegatorShares()),
			Vote:                types.OptionEmpty,
			DelegatorShares:     validator.GetDelegatorShares(),
			DelegatorDeductions: sdk.ZeroDec(),
		}
		systemVotingPower = systemVotingPower.Add(validator.GetTokens())
		return false
	})
	// iterate over all the votes
	keeper.IterateVotes(ctx, proposal.GetProposalID(), func(vote types.Vote) {
		// if validator, just record it in the map
		valAddrStr := sdk.ValAddress(vote.Voter).String()
		if val, ok := currValidators[valAddrStr]; ok {
			val.Vote = vote.Option
			votingVals[valAddrStr] = true
			currValidators[valAddrStr] = val
		}
		// if validator is also delegator
		keeper.ds.IterateDelegations(ctx, vote.Voter, func(index int64, delegation sdk.Delegation) (stop bool) {
			valAddr := delegation.GetValidatorAddr().String()
			if valAddr == valAddrStr {
				return false
			}
			//only tally the delegator voting power under the validator
			if val, ok := currValidators[valAddr]; ok {
				val.DelegatorDeductions = val.DelegatorDeductions.Add(delegation.GetShares())
				currValidators[valAddr] = val

				votingPower := delegation.GetShares().Mul(val.TokenPerShare)
				results[vote.Option] = results[vote.Option].Add(votingPower)
				totalVotingPower = totalVotingPower.Add(votingPower)
			}
			return false
		})
	})

	// iterate over the validators again to tally their voting power
	for _, val := range currValidators {
		if val.Vote == types.OptionEmpty {
			continue
		}

		sharesAfterDeductions := val.DelegatorShares.Sub(val.DelegatorDeductions)
		votingPower := sharesAfterDeductions.Mul(val.TokenPerShare)

		results[val.Vote] = results[val.Vote].Add(votingPower)
		totalVotingPower = totalVotingPower.Add(votingPower)
	}

	tallyingProcedure := keeper.GetTallyingProcedure(ctx, proposal.GetProposalLevel())

	tallyResults = types.TallyResult{
		Yes:               results[types.OptionYes].QuoInt(sdk.AttoScaleFactor),
		Abstain:           results[types.OptionAbstain].QuoInt(sdk.AttoScaleFactor),
		No:                results[types.OptionNo].QuoInt(sdk.AttoScaleFactor),
		NoWithVeto:        results[types.OptionNoWithVeto].QuoInt(sdk.AttoScaleFactor),
		SystemVotingPower: systemVotingPower.QuoInt(sdk.AttoScaleFactor),
	}

	// If no one votes, proposal fails
	if totalVotingPower.Sub(results[types.OptionAbstain]).Equal(sdk.ZeroDec()) {
		return types.REJECT, tallyResults, votingVals
	}

	//if more than 1/3 of voters abstain, proposal fails
	if tallyingProcedure.Participation.GT(totalVotingPower.Quo(systemVotingPower)) {
		return types.REJECT, tallyResults, votingVals
	}

	// If more than 1/3 of voters veto, proposal fails
	if results[types.OptionNoWithVeto].Quo(totalVotingPower).GT(tallyingProcedure.Veto) {
		return types.REJECTVETO, tallyResults, votingVals
	}

	// If more than 1/2 of non-abstaining voters vote Yes, proposal passes
	if results[types.OptionYes].Quo(totalVotingPower).GT(tallyingProcedure.Threshold) {
		return types.PASS, tallyResults, votingVals
	}
	// If more than 1/2 of non-abstaining voters vote No, proposal fails

	return types.REJECT, tallyResults, votingVals
}

func (keeper Keeper) Slash(ctx sdk.Context, p Proposal, voter map[string]bool) {
	for _, valAddr := range keeper.GetValidatorSet(ctx, p.GetProposalID()) {
		if _, ok := voter[valAddr.String()]; !ok {
			val := keeper.ds.GetValidatorSet().Validator(ctx, valAddr)
			if val != nil && val.GetStatus() == sdk.Bonded {
				keeper.ds.GetValidatorSet().Slash(ctx,
					val.GetConsAddr(),
					ctx.BlockHeight(),
					val.GetPower().RoundInt64(),
					keeper.GetTallyingProcedure(ctx, p.GetProposalLevel()).Penalty)
			}
		}
	}
}

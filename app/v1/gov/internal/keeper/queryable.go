package keeper

import (
	"github.com/irisnet/irishub/app/v1/gov/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the governance Querier
const (
	QueryProposals = "proposals"
	QueryProposal  = "proposal"
	QueryDeposits  = "deposits"
	QueryDeposit   = "deposit"
	QueryVotes     = "votes"
	QueryVote      = "vote"
	QueryTally     = "tally"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryProposals:
			return queryProposals(ctx, path[1:], req, keeper)
		case QueryProposal:
			return queryProposal(ctx, path[1:], req, keeper)
		case QueryDeposits:
			return queryDeposits(ctx, path[1:], req, keeper)
		case QueryDeposit:
			return queryDeposit(ctx, path[1:], req, keeper)
		case QueryVotes:
			return queryVotes(ctx, path[1:], req, keeper)
		case QueryVote:
			return queryVote(ctx, path[1:], req, keeper)
		case QueryTally:
			return queryTally(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown gov query endpoint")
		}
	}
}

// Params for query 'custom/gov/proposal'
type QueryProposalParams struct {
	ProposalID uint64
}

// nolint: unparam
func queryProposal(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params QueryProposalParams
	err2 := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err2 != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	proposal := keeper.GetProposal(ctx, params.ProposalID)
	if proposal == nil {
		return nil, types.ErrUnknownProposal(types.DefaultCodespace, params.ProposalID)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, proposal)
	if err2 != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
}

// Params for query 'custom/gov/deposit'
type QueryDepositParams struct {
	ProposalID uint64
	Depositor  sdk.AccAddress
}

// nolint: unparam
func queryDeposit(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params QueryDepositParams
	err2 := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err2 != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	proposal := keeper.GetProposal(ctx, params.ProposalID)
	if proposal == nil {
		return nil, types.ErrUnknownProposal(types.DefaultCodespace, params.ProposalID)
	}

	if proposal.GetStatus() == types.StatusPassed || proposal.GetStatus() == types.StatusRejected {
		return nil, types.ErrCodeDepositDeleted(types.DefaultCodespace, params.ProposalID)
	}

	deposit, bool := keeper.GetDeposit(ctx, params.ProposalID, params.Depositor)
	if !bool {
		return nil, types.ErrCodeDepositNotExisted(types.DefaultCodespace, params.Depositor, params.ProposalID)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, deposit)
	if err2 != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
}

// Params for query 'custom/gov/vote'
type QueryVoteParams struct {
	ProposalID uint64
	Voter      sdk.AccAddress
}

// nolint: unparam
func queryVote(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params QueryVoteParams
	err2 := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err2 != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	proposal := keeper.GetProposal(ctx, params.ProposalID)
	if proposal == nil {
		return nil, types.ErrUnknownProposal(types.DefaultCodespace, params.ProposalID)
	}

	vote, bool := keeper.GetVote(ctx, params.ProposalID, params.Voter)
	if !bool {
		return nil, types.ErrCodeVoteNotExisted(types.DefaultCodespace, params.Voter, params.ProposalID)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, vote)
	if err2 != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
}

// Params for query 'custom/gov/deposits'
type QueryDepositsParams struct {
	ProposalID uint64
}

// nolint: unparam
func queryDeposits(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params QueryDepositsParams
	err2 := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err2 != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	proposal := keeper.GetProposal(ctx, params.ProposalID)
	if proposal == nil {
		return nil, types.ErrUnknownProposal(types.DefaultCodespace, params.ProposalID)
	}

	if proposal.GetStatus() == types.StatusPassed || proposal.GetStatus() == types.StatusRejected {
		return nil, types.ErrCodeDepositDeleted(types.DefaultCodespace, params.ProposalID)
	}

	var deposits []types.Deposit
	depositsIterator := keeper.GetDeposits(ctx, params.ProposalID)
	defer depositsIterator.Close()
	for ; depositsIterator.Valid(); depositsIterator.Next() {
		deposit := types.Deposit{}
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(depositsIterator.Value(), &deposit)
		deposits = append(deposits, deposit)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, deposits)
	if err2 != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
}

// Params for query 'custom/gov/votes'
type QueryVotesParams struct {
	ProposalID uint64
}

// nolint: unparam
func queryVotes(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params QueryVotesParams
	err2 := keeper.cdc.UnmarshalJSON(req.Data, &params)

	if err2 != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	proposal := keeper.GetProposal(ctx, params.ProposalID)
	if proposal == nil {
		return nil, types.ErrUnknownProposal(types.DefaultCodespace, params.ProposalID)
	}

	var votes []types.Vote
	votesIterator := keeper.GetVotes(ctx, params.ProposalID)
	defer votesIterator.Close()
	for ; votesIterator.Valid(); votesIterator.Next() {
		vote := types.Vote{}
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(votesIterator.Value(), &vote)
		votes = append(votes, vote)
	}

	if len(votes) == 0 {
		return nil, nil
	}
	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, votes)
	if err2 != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
}

// Params for query 'custom/gov/proposals'
type QueryProposalsParams struct {
	Voter          sdk.AccAddress
	Depositor      sdk.AccAddress
	ProposalStatus string
	Limit          uint64
}

// nolint: unparam
func queryProposals(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params QueryProposalsParams
	err2 := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err2 != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	var status = types.StatusNil
	if s, err := types.ProposalStatusFromString(params.ProposalStatus); err == nil {
		status = s
	}

	proposals := keeper.GetProposalsFiltered(ctx, params.Voter, params.Depositor, status, params.Limit)
	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, proposals)
	if err2 != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
}

// Params for query 'custom/gov/tally'
type QueryTallyParams struct {
	ProposalID uint64
}

// nolint: unparam
func queryTally(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	// TODO: Dependant on #1914

	var param QueryTallyParams
	err2 := keeper.cdc.UnmarshalJSON(req.Data, &param)
	if err2 != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	proposal := keeper.GetProposal(ctx, param.ProposalID)
	if proposal == nil {
		return nil, types.ErrUnknownProposal(types.DefaultCodespace, param.ProposalID)
	}

	var tallyResult types.TallyResult

	if proposal.GetStatus() == types.StatusDepositPeriod {
		tallyResult = types.EmptyTallyResult()
	} else if proposal.GetStatus() == types.StatusPassed || proposal.GetStatus() == types.StatusRejected {
		tallyResult = proposal.GetTallyResult()
	} else {
		_, tallyResult, _ = keeper.Tally(ctx, proposal)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, tallyResult)
	if err2 != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
}

package gov

import (
	"github.com/irisnet/irishub/app/v1/gov/internal/keeper"
	"github.com/irisnet/irishub/app/v1/gov/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

const StartingProposalID = 1

// GenesisState - all gov state that must be provided at genesis
type GenesisState struct {
	Params types.GovParams `json:"params"` // inflation params
}

func NewGenesisState(systemHaltPeriod int64, params types.GovParams) GenesisState {
	return GenesisState{
		Params: params,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: types.DefaultParams(),
	}
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data GenesisState) {
	err := ValidateGenesis(data)
	if err != nil {
		// TODO: Handle this with #870
		panic(err)
	}

	err = k.SetInitialProposalID(ctx, StartingProposalID)
	if err != nil {
		// TODO: Handle this with #870
		panic(err)
	}

	k.SetSystemHaltHeight(ctx, -1)
	k.SetParamSet(ctx, data.Params)
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) GenesisState {

	return GenesisState{
		Params: k.GetParamSet(ctx),
	}
}

func ValidateGenesis(data GenesisState) error {
	err := ValidateParams(data.Params)
	if err != nil {
		return err
	}
	return nil
}

// get raw genesis raw message for testing
func DefaultGenesisStateForCliTest() GenesisState {

	return GenesisState{
		Params: types.DefaultParamsForTest(),
	}
}

func PrepForZeroHeightGenesis(ctx sdk.Context, k keeper.Keeper) {
	proposals := k.GetProposalsFiltered(ctx, nil, nil, StatusDepositPeriod, 0)
	for _, proposal := range proposals {
		proposalID := proposal.GetProposalID()
		k.RefundDeposits(ctx, proposalID)
	}

	proposals = k.GetProposalsFiltered(ctx, nil, nil, StatusVotingPeriod, 0)
	for _, proposal := range proposals {
		proposalID := proposal.GetProposalID()
		k.RefundDeposits(ctx, proposalID)
	}
}

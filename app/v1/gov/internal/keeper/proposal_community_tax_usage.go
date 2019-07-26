package keeper

import (
	"github.com/irisnet/irishub/app/v1/gov/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

// Implements Proposal Interface
var _ Proposal = (*CommunityTaxUsageProposal)(nil)

type TaxUsage struct {
	Usage       types.UsageType `json:"usage"`
	DestAddress sdk.AccAddress  `json:"dest_address"`
	Percent     sdk.Dec         `json:"percent"`
}

type CommunityTaxUsageProposal struct {
	BasicProposal
	TaxUsage TaxUsage `json:"tax_usage"`
}

func (tp CommunityTaxUsageProposal) GetTaxUsage() TaxUsage { return tp.TaxUsage }
func (tp *CommunityTaxUsageProposal) SetTaxUsage(taxUsage TaxUsage) {
	tp.TaxUsage = taxUsage
}

func (tp *CommunityTaxUsageProposal) Validate(ctx sdk.Context, k Keeper, verify bool) sdk.Error {
	if err := tp.BasicProposal.Validate(ctx, k, verify); err != nil {
		return err
	}

	if tp.TaxUsage.Usage != types.UsageTypeBurn {
		_, found := k.gk.GetTrustee(ctx, tp.TaxUsage.DestAddress)
		if !found {
			return types.ErrNotTrustee(k.codespace, tp.TaxUsage.DestAddress)
		}
	}
	return nil
}

func (tp *CommunityTaxUsageProposal) Execute(ctx sdk.Context, gk Keeper) sdk.Error {
	logger := ctx.Logger()
	if err := tp.Validate(ctx, gk, false); err != nil {
		logger.Error("Execute CommunityTaxUsageProposal Failure", "info",
			"the destination address is not a trustee now", "destinationAddress", tp.TaxUsage.DestAddress)
		return err
	}
	burn := false
	if tp.TaxUsage.Usage == types.UsageTypeBurn {
		burn = true
	}
	gk.dk.AllocateFeeTax(ctx, tp.TaxUsage.DestAddress, tp.TaxUsage.Percent, burn)
	return nil
}

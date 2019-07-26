package types

import (
	"github.com/irisnet/irishub/app/v1/asset/exported"
	"github.com/irisnet/irishub/modules/guardian"
	sdk "github.com/irisnet/irishub/types"
)

// AssetKeeper expected asset keeper
type AssetKeeper interface {
	IssueToken(ctx sdk.Context, token exported.FungibleToken) (sdk.Tags, sdk.Error)
	HasToken(ctx sdk.Context, tokenId string) bool
}

type ProtocolKeeper interface {
	GetUpgradeConfig(ctx sdk.Context) (sdk.UpgradeConfig, bool)
	SetUpgradeConfig(ctx sdk.Context, upgradeConfig sdk.UpgradeConfig)
	IsValidVersion(ctx sdk.Context, version uint64) bool
}

type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)
	BurnCoins(ctx sdk.Context, fromAddr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)
}

type DistrKeeper interface {
	AllocateFeeTax(ctx sdk.Context, destAddr sdk.AccAddress, percent sdk.Dec, burn bool)
}

type GuardianKeeper interface {
	GetTrustee(ctx sdk.Context, addr sdk.AccAddress) (guardian.Guardian, bool)
	GetProfiler(ctx sdk.Context, addr sdk.AccAddress) (guardian.Guardian, bool)
}

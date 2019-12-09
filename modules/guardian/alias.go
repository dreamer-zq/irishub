// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/irisnet/irishub/modules/guardian/internal/keeper
// ALIASGEN: github.com/irisnet/irishub/modules/guardian/internal/types
package guardian

import (
	"github.com/irisnet/irishub/modules/guardian/internal/keeper"
	"github.com/irisnet/irishub/modules/guardian/internal/types"
)

const (
	ModuleName       = types.ModuleName
	DefaultCodespace = types.DefaultCodespace
	StoreKey         = types.StoreKey
	RouterKey        = types.RouterKey
	QuerierRoute     = types.QuerierRoute
	Ordinary         = types.Ordinary
	Genesis          = types.Genesis
	QueryProfilers   = types.QueryProfilers
	QueryTrustees    = types.QueryTrustees

	EventTypeAddProfiler    = types.EventTypeAddProfiler
	EventTypeAddTrustee     = types.EventTypeAddTrustee
	EventTypeDeleteProfiler = types.EventTypeDeleteProfiler
	EventTypeDeleteTrustee  = types.EventTypeDeleteTrustee

	AttributeValueCategory      = types.AttributeValueCategory
	AttributeKeyProfilerAddress = types.AttributeKeyProfilerAddress
	AttributeKeyTrusteeAddress  = types.AttributeKeyTrusteeAddress
	AttributeKeyAddedBy         = types.AttributeKeyAddedBy
	AttributeKeyDeletedBy       = types.AttributeKeyDeletedBy
)

var (
	// functions aliases
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	NewGuardian         = types.NewGuardian
	RegisterCodec       = types.RegisterCodec

	// variable aliases
	ModuleCdc = types.ModuleCdc

	// errors aliases
	ErrInvalidOperator       = types.ErrInvalidOperator
	ErrProfilerExists        = types.ErrProfilerExists
	ErrProfilerNotExists     = types.ErrProfilerNotExists
	ErrTrusteeExists         = types.ErrTrusteeExists
	ErrTrusteeNotExists      = types.ErrTrusteeNotExists
	ErrDeleteGenesisProfiler = types.ErrDeleteGenesisProfiler
	ErrDeleteGenesisTrustee  = types.ErrDeleteGenesisTrustee
)

type (
	Keeper            = keeper.Keeper
	GenesisState      = types.GenesisState
	Guardian          = types.Guardian
	Profilers         = types.Profilers
	Trustees          = types.Trustees
	MsgAddProfiler    = types.MsgAddProfiler
	MsgAddTrustee     = types.MsgAddTrustee
	MsgDeleteProfiler = types.MsgDeleteProfiler
	MsgDeleteTrustee  = types.MsgDeleteTrustee
)

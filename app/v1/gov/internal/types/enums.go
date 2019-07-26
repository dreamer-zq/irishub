package types

import (
	"encoding/json"
	"fmt"
	sdk "github.com/irisnet/irishub/types"
	"github.com/pkg/errors"
)

//-----------------------------------------------------------
// ProposalQueue
type ProposalQueue []uint64

//-----------------------------------------------------------
// ProposalKind

// Type that represents Proposal Type as a byte
type ProposalKind byte

//nolint
const (
	ProposalTypeNil               ProposalKind = 0x00
	ProposalTypeParameter         ProposalKind = 0x01
	ProposalTypeSoftwareUpgrade   ProposalKind = 0x02
	ProposalTypeSystemHalt        ProposalKind = 0x03
	ProposalTypeCommunityTaxUsage ProposalKind = 0x04
	ProposalTypePlainText         ProposalKind = 0x05
	ProposalTypeTokenAddition     ProposalKind = 0x06
)

type pTypeInfo struct {
	Type  ProposalKind
	Level ProposalLevel
}

var pTypeMap = map[string]pTypeInfo{
	"PlainText":         {},
	"Parameter":         {},
	"SoftwareUpgrade":   {},
	"SystemHalt":        {},
	"CommunityTaxUsage": {},
	"TokenAddition":     {},
}

// String to proposalType byte.  Returns ff if invalid.
func ProposalTypeFromString(str string) (ProposalKind, error) {
	kind, ok := pTypeMap[str]
	if !ok {
		return ProposalKind(0xff), errors.Errorf("'%s' is not a valid proposal type", str)
	}
	return kind.Type, nil
}

// is defined ProposalType?
func ValidProposalType(pt ProposalKind) bool {
	_, ok := pTypeMap[pt.String()]
	return ok
}

// Marshal needed for protobuf compatibility
func (pt ProposalKind) Marshal() ([]byte, error) {
	return []byte{byte(pt)}, nil
}

// Unmarshal needed for protobuf compatibility
func (pt *ProposalKind) Unmarshal(data []byte) error {
	*pt = ProposalKind(data[0])
	return nil
}

// Marshals to JSON using string
func (pt ProposalKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(pt.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (pt *ProposalKind) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	bz2, err := ProposalTypeFromString(s)
	if err != nil {
		return err
	}
	*pt = bz2
	return nil
}

// Turns VoteOption byte to String
func (pt ProposalKind) String() string {
	for k, v := range pTypeMap {
		if v.Type == pt {
			return k
		}
	}
	return ""
}

// For Printf / Sprintf, returns bech32 when using %s
// nolint: errcheck
func (pt ProposalKind) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(pt.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(pt))))
	}
}

func (pt ProposalKind) GetProposalLevel() ProposalLevel {
	return pTypeMap[pt.String()].Level
}

//-----------------------------------------------------------
// ProposalStatus

// Type that represents Proposal Status as a byte
type ProposalStatus byte

//nolint
const (
	StatusNil           ProposalStatus = 0x00
	StatusDepositPeriod ProposalStatus = 0x01
	StatusVotingPeriod  ProposalStatus = 0x02
	StatusPassed        ProposalStatus = 0x03
	StatusRejected      ProposalStatus = 0x04
)

var pStatusMap = map[string]ProposalStatus{
	"DepositPeriod": StatusDepositPeriod,
	"VotingPeriod":  StatusVotingPeriod,
	"Passed":        StatusPassed,
	"Rejected":      StatusRejected,
}

// ProposalStatusToString turns a string into a ProposalStatus
func ProposalStatusFromString(str string) (ProposalStatus, error) {
	status, ok := pStatusMap[str]
	if !ok {
		return ProposalStatus(0xff), errors.Errorf("'%s' is not a valid proposal status", str)
	}
	return status, nil
}

// is defined ProposalType?
func ValidProposalStatus(status ProposalStatus) bool {
	_, ok := pStatusMap[status.String()]
	return ok
}

// Marshal needed for protobuf compatibility
func (status ProposalStatus) Marshal() ([]byte, error) {
	return []byte{byte(status)}, nil
}

// Unmarshal needed for protobuf compatibility
func (status *ProposalStatus) Unmarshal(data []byte) error {
	*status = ProposalStatus(data[0])
	return nil
}

// Marshals to JSON using string
func (status ProposalStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(status.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (status *ProposalStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	bz2, err := ProposalStatusFromString(s)
	if err != nil {
		return err
	}
	*status = bz2
	return nil
}

// Turns VoteStatus byte to String
func (status ProposalStatus) String() string {
	for k, v := range pStatusMap {
		if v == status {
			return k
		}
	}
	return ""
}

// For Printf / Sprintf, returns bech32 when using %s
// nolint: errcheck
func (status ProposalStatus) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(status.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(status))))
	}
}

//-----------------------------------------------------------
// Tally Results
type TallyResult struct {
	Yes               sdk.Dec `json:"yes"`
	Abstain           sdk.Dec `json:"abstain"`
	No                sdk.Dec `json:"no"`
	NoWithVeto        sdk.Dec `json:"no_with_veto"`
	SystemVotingPower sdk.Dec `json:"system_voting_power"`
}

// checks if two proposals are equal
func EmptyTallyResult() TallyResult {
	return TallyResult{
		Yes:               sdk.ZeroDec(),
		Abstain:           sdk.ZeroDec(),
		No:                sdk.ZeroDec(),
		NoWithVeto:        sdk.ZeroDec(),
		SystemVotingPower: sdk.ZeroDec(),
	}
}

// checks if two proposals are equal
func (tr TallyResult) Equals(resultB TallyResult) bool {
	return tr.Yes.Equal(resultB.Yes) &&
		tr.Abstain.Equal(resultB.Abstain) &&
		tr.No.Equal(resultB.No) &&
		tr.NoWithVeto.Equal(resultB.NoWithVeto) &&
		tr.SystemVotingPower.Equal(resultB.SystemVotingPower)
}

func (tr TallyResult) String() string {
	return fmt.Sprintf(`Tally Result:
  Yes:                %s
  Abstain:            %s
  No:                 %s
  NoWithVeto:         %s
  SystemVotingPower:  %s`, tr.Yes.String(), tr.Abstain.String(), tr.No.String(), tr.NoWithVeto.String(), tr.SystemVotingPower.String())
}

//-----------------------------------------------------------
// ProposalLevel

// Type that represents Proposal Level as a byte
type ProposalLevel byte

//nolint
const (
	ProposalLevelNil       ProposalLevel = 0x00
	ProposalLevelCritical  ProposalLevel = 0x01
	ProposalLevelImportant ProposalLevel = 0x02
	ProposalLevelNormal    ProposalLevel = 0x03
)

func (p ProposalLevel) String() string {
	switch p {
	case ProposalLevelCritical:
		return "critical"
	case ProposalLevelImportant:
		return "important"
	case ProposalLevelNormal:
		return "normal"
	default:
		return " "
	}
}

type UsageType byte

const (
	UsageTypeBurn       UsageType = 0x01
	UsageTypeDistribute UsageType = 0x02
	UsageTypeGrant      UsageType = 0x03
)

// String to UsageType byte.  Returns ff if invalid.
func UsageTypeFromString(str string) (UsageType, error) {
	switch str {
	case "Burn":
		return UsageTypeBurn, nil
	case "Distribute":
		return UsageTypeDistribute, nil
	case "Grant":
		return UsageTypeGrant, nil
	default:
		return UsageType(0xff), errors.Errorf("'%s' is not a valid usage type", str)
	}
}

// is defined UsageType?
func ValidUsageType(ut UsageType) bool {
	if ut == UsageTypeBurn ||
		ut == UsageTypeDistribute ||
		ut == UsageTypeGrant {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility
func (ut UsageType) Marshal() ([]byte, error) {
	return []byte{byte(ut)}, nil
}

// Unmarshal needed for protobuf compatibility
func (ut *UsageType) Unmarshal(data []byte) error {
	*ut = UsageType(data[0])
	return nil
}

// Marshals to JSON using string
func (ut UsageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(ut.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (ut *UsageType) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	bz2, err := UsageTypeFromString(s)
	if err != nil {
		return err
	}
	*ut = bz2
	return nil
}

// Turns VoteOption byte to String
func (ut UsageType) String() string {
	switch ut {
	case UsageTypeBurn:
		return "Burn"
	case UsageTypeDistribute:
		return "Distribute"
	case UsageTypeGrant:
		return "Grant"
	default:
		return ""
	}
}

// For Printf / Sprintf, returns bech32 when using %s
// nolint: errcheck
func (ut UsageType) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(ut.String()))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(ut))))
	}
}

type ProposalResult string

const (
	PASS       ProposalResult = "pass"
	REJECT     ProposalResult = "reject"
	REJECTVETO ProposalResult = "reject-veto"
)

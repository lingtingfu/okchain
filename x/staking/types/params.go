package types

import (
	"bytes"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/config"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Staking params default values
const (
	// DefaultUnbondingTime reflects three weeks in seconds as the default
	// unbonding time.
	// TODO: Justify our choice of default here.
	DefaultUnbondingTime = config.DefaultUnbondingTime

	// Default maximum number of bonded validators
	DefaultMaxValidators = config.DefaultMaxValidators

	DefaultEpoch         uint16 = config.DefaultBlocksPerEpoch
	DefaultMaxValsToVote uint16 = config.DefaultMaxValsToVote
)

var (
	// DefaultMinSelfDelegationLimit is the limit value of min self delegation
	DefaultMinSelfDelegationLimit = config.DefaultMinSelfDelegationLimit
	// DefaultMinDelegation is the limit value of delegation or undelegation
	DefaultMinDelegation = config.DefaultMinDelegation
)

// nolint - Keys for parameter access
var (
	KeyUnbondingTime     = []byte("UnbondingTime")
	KeyMaxValidators     = []byte("MaxValidators")
	KeyBondDenom         = []byte("BondDenom")
	KeyEpoch             = []byte("BlocksPerEpoch")    // how many blocks each epoch has
	KeyTheEndOfLastEpoch = []byte("TheEndOfLastEpoch") // a block height that is the end of last epoch

	KeyMaxValsToVote          = []byte("MaxValsToVote")
	KeyMinSelfDelegationLimit = []byte("MinSelfDelegationLimit")
	KeyMinDelegation          = []byte("MinDelegation")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for staking
type Params struct {
	// time duration of unbonding
	UnbondingTime time.Duration `json:"unbonding_time" yaml:"unbonding_time"`
	// note: we need to be a bit careful about potential overflow here, since this is user-determined
	// maximum number of validators (max uint16 = 65535)
	MaxValidators uint16 `json:"max_bonded_validators" yaml:"max_bonded_validators"`
	//epoch for validator update
	Epoch         uint16 `json:"epoch" yaml:"epoch"`
	MaxValsToVote uint16 `json:"max_validators_to_vote" yaml:"max_validators_to_vote"`
	// bondable coin denomination
	BondDenom string `json:"bond_denom" yaml:"bond_denom"`
	// limited amount of the msd
	MinSelfDelegationLimit sdk.Dec `json:"min_self_delegation" yaml:"min_self_delegation"`
	//limited amount of delegate
	MinDelegation sdk.Dec `json:"min_delegation" yaml:"min_delegation"`
}

// NewParams creates a new Params instance
func NewParams(unbondingTime time.Duration, maxValidators uint16, bondDenom string, epoch uint16, maxValsToVote uint16,
	minSelfDelegationLimited sdk.Dec, minDelegation sdk.Dec) Params {

	return Params{
		UnbondingTime:          unbondingTime,
		MaxValidators:          maxValidators,
		BondDenom:              bondDenom,
		Epoch:                  epoch,
		MaxValsToVote:          maxValsToVote,
		MinSelfDelegationLimit: minSelfDelegationLimited,
		MinDelegation:          minDelegation,
	}
}

// ParamSetPairs is the implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyUnbondingTime, Value: &p.UnbondingTime},
		{Key: KeyMaxValidators, Value: &p.MaxValidators},
		{Key: KeyBondDenom, Value: &p.BondDenom},
		{Key: KeyEpoch, Value: &p.Epoch},
		{Key: KeyMaxValsToVote, Value: &p.MaxValsToVote},
		{Key: KeyMinSelfDelegationLimit, Value: &p.MinSelfDelegationLimit},
		{Key: KeyMinDelegation, Value: &p.MinDelegation},
	}
}

// Equal returns a boolean determining if two Param types are identical
// TODO: This is slower than comparing struct fields directly
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(DefaultUnbondingTime, DefaultMaxValidators,
		sdk.DefaultBondDenom, DefaultEpoch, DefaultMaxValsToVote,
		DefaultMinSelfDelegationLimit, DefaultMinDelegation)
}

// String returns a human readable string representation of the Params
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Unbonding Time:    		%s
  Max Validators:   	 	%d
  Epoch: 					%d
  Bonded Coin Denom: 		%s
  MaxValsToVote:     		%d
  MinSelfDelegationLimited  %d
  MinDelegation				%d`, p.UnbondingTime,
		p.MaxValidators, p.Epoch, p.BondDenom, p.MaxValsToVote, p.MinSelfDelegationLimit, p.MinDelegation)
}

// Validate gives a quick validity check for a set of params
func (p Params) Validate() error {
	if p.BondDenom == "" {
		return fmt.Errorf("staking parameter BondDenom can't be an empty string")
	}
	if p.MaxValidators == 0 {
		return fmt.Errorf("staking parameter MaxValidators must be a positive integer")
	}
	if p.Epoch == 0 {
		return fmt.Errorf("staking parameter Epoch must be a positive integer")
	}
	if p.MaxValsToVote == 0 {
		return fmt.Errorf("staking parameter MaxValsToVote must be a positive integer")
	}
	if p.MinSelfDelegationLimit.LTE(sdk.ZeroDec()) {
		return fmt.Errorf("staking parameter MinSelfDelegationLimit cannot be a negative integer")
	}
	return nil
}

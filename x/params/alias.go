// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/pokt-network/posmint/x/params/subspace
// ALIASGEN: github.com/pokt-network/posmint/x/params/types
package params

import (
	"github.com/pokt-network/posmint/x/params/subspace"
	"github.com/pokt-network/posmint/x/params/types"
)

const (
	StoreKey             = subspace.StoreKey
	TStoreKey            = subspace.TStoreKey
	TestParamStore       = subspace.TestParamStore
	DefaultCodespace     = types.DefaultCodespace
	CodeUnknownSubspace  = types.CodeUnknownSubspace
	CodeSettingParameter = types.CodeSettingParameter
	CodeEmptyData        = types.CodeEmptyData
	ModuleName           = types.ModuleName
	RouterKey            = types.RouterKey
)

var (
	// functions aliases
	NewSubspace           = subspace.NewSubspace
	NewKeyTable           = subspace.NewKeyTable
	DefaultTestComponents = subspace.DefaultTestComponents
	RegisterCodec         = types.RegisterCodec
	ErrUnknownSubspace    = types.ErrUnknownSubspace
	ErrSettingParameter   = types.ErrSettingParameter
	ErrEmptyChanges       = types.ErrEmptyChanges
	ErrEmptySubspace      = types.ErrEmptySubspace
	ErrEmptyKey           = types.ErrEmptyKey
	ErrEmptyValue         = types.ErrEmptyValue

	// variable aliases
	ModuleCdc = types.ModuleCdc
)

type (
	ParamSetPair     = subspace.ParamSetPair
	ParamSetPairs    = subspace.ParamSetPairs
	ParamSet         = subspace.ParamSet
	Subspace         = subspace.Subspace
	ReadOnlySubspace = subspace.ReadOnlySubspace
	KeyTable         = subspace.KeyTable
)

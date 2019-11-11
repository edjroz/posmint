package params

import (
	"encoding/json"

	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/params/types"
)

var (
	_ module.AppModuleBasic = AppModuleBasic{}
)

const moduleName = "params"

// app module basics object
type AppModuleBasic struct{}

// module name
func (AppModuleBasic) Name() string {
	return moduleName
}

// register module codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

// default genesis state
func (AppModuleBasic) DefaultGenesis() json.RawMessage { return nil }

// module validate genesis
func (AppModuleBasic) ValidateGenesis(_ json.RawMessage) error { return nil }

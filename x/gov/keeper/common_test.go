package keeper

import (
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/crypto"
	"github.com/pokt-network/posmint/store"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/keeper"
	govTypes "github.com/pokt-network/posmint/x/gov/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"math/rand"
	"testing"
)

var (
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
	)
)

// nolint: deadcode unused
// create a codec used only for testing
func makeTestCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	govTypes.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

func getRandomPubKey() crypto.Ed25519PublicKey {
	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])
	return pub
}

func getRandomValidatorAddress() sdk.Address {
	return sdk.Address(getRandomPubKey().Address())
}

// nolint: deadcode unused
func createTestKeeperAndContext(t *testing.T, isCheckTx bool) (sdk.Context, Keeper) {
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(sdk.ParamsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(sdk.ParamsTKey, sdk.StoreTypeTransient, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain"}, isCheckTx, log.NewNopLogger())
	ctx = ctx.WithConsensusParams(
		&abci.ConsensusParams{
			Validator: &abci.ValidatorParams{
				PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519},
			},
		},
	)
	cdc := makeTestCodec()
	maccPerms := map[string][]string{
		auth.FeeCollectorName:   nil,
		govTypes.DAOAccountName: {"burner", "staking", "minter"},
		"FAKE":                  {"burner", "staking", "minter"},
	}
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[auth.NewModuleAddress(acc).String()] = true
	}
	akSubspace := sdk.NewSubspace(auth.DefaultParamspace)
	ak := keeper.NewKeeper(cdc, keyAcc, akSubspace, maccPerms)
	ak.GetModuleAccount(ctx, "FAKE")
	pk := NewKeeper(cdc, sdk.ParamsKey, sdk.ParamsTKey, govTypes.DefaultParamspace, ak, akSubspace)
	moduleManager := module.NewManager(
		auth.NewAppModule(ak),
	)
	genesisState := ModuleBasics.DefaultGenesis()
	moduleManager.InitGenesis(ctx, genesisState)
	params := govTypes.DefaultParams()
	pk.SetParams(ctx, params)
	gs := govTypes.DefaultGenesisState()
	acl := createTestACL()
	gs.Params.ACL = acl
	pk.InitGenesis(ctx, gs)
	return ctx, pk
}

var testACL govTypes.ACL

func createTestACL() govTypes.ACL {
	if testACL == nil {
		acl := govTypes.ACL(make([]govTypes.ACLPair, 0))
		acl.SetOwner("auth/MaxMemoCharacters", getRandomValidatorAddress())
		acl.SetOwner("auth/TxSigLimit", getRandomValidatorAddress())
		acl.SetOwner("gov/daoOwner", getRandomValidatorAddress())
		acl.SetOwner("gov/acl", getRandomValidatorAddress())
		acl.SetOwner("gov/upgrade", getRandomValidatorAddress())
		testACL = acl
	}
	return testACL
}

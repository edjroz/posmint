package keeper

import (
	kb "github.com/pokt-network/posmint/crypto/keys"
	"github.com/pokt-network/posmint/types/module"
	"github.com/pokt-network/posmint/x/supply"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	dbm "github.com/tendermint/tm-db"

	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/store"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/bank"
	"github.com/pokt-network/posmint/x/params"
	"github.com/pokt-network/posmint/x/pos/types"
	//"github.com/pokt-network/posmint/x/supply/internal/types"

	sdk "github.com/pokt-network/posmint/types"
)

// nolint: deadcode unused
var (
	multiPerm    = "multiple permissions account"
	randomPerm   = "random permission"
	holder       = "holder"
	Keybase      kb.Keybase
	TMNode       *node.Node
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		params.AppModuleBasic{},
		supply.AppModuleBasic{},
	)
)

// nolint: deadcode unused
// create a codec used only for testing
func makeTestCodec() *codec.Codec {
	var cdc = codec.New()

	bank.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

// nolint: deadcode unused
func createTestInput(t *testing.T, isCheckTx bool, initPower int64, nAccs int64) (sdk.Context, []auth.Account, Keeper) {
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
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
		auth.FeeCollectorName: nil,
		types.StakedPoolName:  {supply.Burner, supply.Staking},
	}
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}
	valTokens := sdk.TokensFromConsensusPower(initPower)

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(cdc, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, modAccAddrs)
	sk := supply.NewKeeper(cdc, keySupply, ak, bk, maccPerms)

	moduleManager := module.NewManager(
		auth.NewAppModule(ak, TMNode, Keybase),
		bank.NewAppModule(bk, ak, TMNode, Keybase),
		supply.NewAppModule(sk, ak, TMNode, Keybase),
	)

	genesisState := ModuleBasics.DefaultGenesis()
	moduleManager.InitGenesis(ctx, genesisState)

	posSubSpace := pk.Subspace(DefaultParamspace)
	initialCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, valTokens))
	accs := createTestAccs(ctx, int(nAccs), initialCoins, &ak)

	keeper := NewKeeper(cdc, keySupply, bk, sk, posSubSpace, sdk.CodespaceType("pos"))

	params := types.DefaultParams()
	keeper.SetParams(ctx, params)
	return ctx, accs, keeper
}

// nolint: unparam deadcode unused
func createTestAccs(ctx sdk.Context, numAccs int, initialCoins sdk.Coins, ak *auth.AccountKeeper) (accs []auth.Account) {
	for i := 0; i < numAccs; i++ {
		privKey := secp256k1.GenPrivKey()
		pubKey := privKey.PubKey()
		addr := sdk.AccAddress(pubKey.Address())
		acc := auth.NewBaseAccountWithAddress(addr)
		acc.Coins = initialCoins
		acc.PubKey = pubKey
		acc.AccountNumber = uint64(i)
		ak.SetAccount(ctx, &acc)
	}
	return
}

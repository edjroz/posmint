package types

import (
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	stakeDenom = "stake"
	feeDenom   = "fee"
)

func TestBaseAddressPubKey(t *testing.T) {
	_, pub1, addr1 := KeyTestPubAddr()
	_, pub2, addr2 := KeyTestPubAddr()
	acc := NewBaseAccountWithAddress(addr1)

	// check the address (set) and pubkey (not set)
	require.EqualValues(t, addr1, acc.GetAddress())
	require.EqualValues(t, nil, acc.GetPubKey())

	// can't override address
	err := acc.SetAddress(addr2)
	require.NotNil(t, err)
	require.EqualValues(t, addr1, acc.GetAddress())

	// set the pubkey
	err = acc.SetPubKey(pub1)
	require.Nil(t, err)
	require.Equal(t, pub1, acc.GetPubKey())

	// can override pubkey
	err = acc.SetPubKey(pub2)
	require.Nil(t, err)
	require.Equal(t, pub2, acc.GetPubKey())

	//------------------------------------

	// can set address on empty account
	acc2 := BaseAccount{}
	err = acc2.SetAddress(addr2)
	require.Nil(t, err)
	require.EqualValues(t, addr2, acc2.GetAddress())
}

func TestBaseAccountCoins(t *testing.T) {
	_, _, addr := KeyTestPubAddr()
	acc := NewBaseAccountWithAddress(addr)

	someCoins := sdk.Coins{sdk.NewInt64Coin("atom", 123), sdk.NewInt64Coin("eth", 246)}

	err := acc.SetCoins(someCoins)
	require.Nil(t, err)
	require.Equal(t, someCoins, acc.GetCoins())
}

func TestBaseAccountSequence(t *testing.T) {
	_, _, addr := KeyTestPubAddr()
	acc := NewBaseAccountWithAddress(addr)

	seq := uint64(7)

	err := acc.SetSequence(seq)
	require.Nil(t, err)
	require.Equal(t, seq, acc.GetSequence())
}

func TestBaseAccountMarshal(t *testing.T) {
	_, pub, addr := KeyTestPubAddr()
	acc := NewBaseAccountWithAddress(addr)

	someCoins := sdk.Coins{sdk.NewInt64Coin("atom", 123), sdk.NewInt64Coin("eth", 246)}
	seq := uint64(7)

	// set everything on the account
	err := acc.SetPubKey(pub)
	require.Nil(t, err)
	err = acc.SetSequence(seq)
	require.Nil(t, err)
	err = acc.SetCoins(someCoins)
	require.Nil(t, err)

	// need a codec for marshaling
	cdc := codec.New()
	codec.RegisterCrypto(cdc)

	b, err := cdc.MarshalBinaryLengthPrefixed(acc)
	require.Nil(t, err)

	acc2 := BaseAccount{}
	err = cdc.UnmarshalBinaryLengthPrefixed(b, &acc2)
	require.Nil(t, err)
	require.Equal(t, acc, acc2)

	// error on bad bytes
	acc2 = BaseAccount{}
	err = cdc.UnmarshalBinaryLengthPrefixed(b[:len(b)/2], &acc2)
	require.NotNil(t, err)
}

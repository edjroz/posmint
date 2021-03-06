package keys

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/pokt-network/posmint/crypto/keys/mintkey"
	"github.com/pokt-network/posmint/types"
)

func init() {
	mintkey.BcryptSecurityParameter = 1
}

// TestKeyManagement makes sure we can manipulate these keys well
func TestKeyManagement(t *testing.T) {
	// make the storage with reasonable defaults
	cstore := NewInMemory()

	//n1, n2, n3 := "personal", "business", "other"
	p1, p2 := "1234", "really-secure!@#$"

	// Check empty state
	l, err := cstore.List()
	require.Nil(t, err)
	assert.Empty(t, l)

	// Fetching a non existent address should throw an error
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])
	blankAddress := types.AccAddress(pub.Address())
	_, err = cstore.Get(blankAddress)
	require.Error(t, err)

	// create some keys
	kp1, err := cstore.Create(p1)
	require.NoError(t, err)

	kp2, err := cstore.Create(p2)
	require.NoError(t, err)

	// we can get these keys
	keyPairList, err := cstore.List()
	require.NotEmpty(t, keyPairList)
	retrievedKp1, err := cstore.Get(kp1.GetAddress())
	require.NoError(t, err)
	require.NotEqual(t, KeyPair{}, retrievedKp1)
	retrievedKp2, err := cstore.Get(kp2.GetAddress())
	require.NoError(t, err)
	require.NotEqual(t, KeyPair{}, retrievedKp2)

	// List retrieves all keypairs
	keyS, err := cstore.List()
	require.NoError(t, err)
	require.Equal(t, 2, len(keyS))

	// deleting a key removes it
	err = cstore.Delete(blankAddress, "foo")
	require.NotNil(t, err)
	err = cstore.Delete(kp1.GetAddress(), p1)
	require.NoError(t, err)
	keyS, err = cstore.List()
	require.NoError(t, err)
	require.Equal(t, 1, len(keyS))
	_, err = cstore.Get(kp1.GetAddress())
	require.Error(t, err)
}

// TestSignVerify does some detailed checks on how we sign and validate
// signatures
func TestSignVerify(t *testing.T) {
	cstore := NewInMemory()

	//n1, n2, n3 := "some dude", "a dudette", "dude-ish"
	//p1, p2, p3 := "1234", "foobar", "foobar"

	// create a user and get their info
	passphrase := "1234"
	kp, err := cstore.Create(passphrase)
	require.Nil(t, err)

	// let's try to sign some messages
	d1 := []byte("my first message")
	d2 := []byte("some other important info!")
	d3 := []byte("feels like I forgot something...")

	// try signing both data with both ..
	s1, pub1, err := cstore.Sign(kp.GetAddress(), passphrase, d1)
	require.Nil(t, err)
	require.Equal(t, kp.PubKey, pub1)

	s2, pub2, err := cstore.Sign(kp.GetAddress(), passphrase, d2)
	require.Nil(t, err)
	require.Equal(t, kp.PubKey, pub2)

	s3, pub3, err := cstore.Sign(kp.GetAddress(), passphrase, d3)
	require.Nil(t, err)
	require.Equal(t, kp.PubKey, pub3)

	// let's try to validate and make sure it only works when everything is proper
	cases := []struct {
		key   crypto.PubKey
		data  []byte
		sig   []byte
		valid bool
	}{
		// proper matches
		{kp.PubKey, d1, s1, true},
		{kp.PubKey, d2, s2, true},
		{kp.PubKey, d3, s3, true},
		// change data, pubkey, or signature leads to fail
		{kp.PubKey, d1, s2, false},
		{kp.PubKey, d2, s3, false},
		{kp.PubKey, d3, s1, false},
	}

	for i, tc := range cases {
		valid := tc.key.VerifyBytes(tc.data, tc.sig)
		require.Equal(t, tc.valid, valid, "%d", i)
	}
}

// TestExportImport tests exporting and importing
func TestArmoredExportImport(t *testing.T) {
	// make the storage with reasonable defaults
	cstore := NewInMemory()

	// Create an account
	passphrase := "1234"
	kp, err := cstore.Create(passphrase)
	require.NoError(t, err)

	// Export the account armored
	armoredKey, err := cstore.ExportPrivKeyEncryptedArmor(kp.GetAddress(), passphrase, passphrase)
	require.NoError(t, err)
	require.NotEmpty(t, armoredKey)

	// Import armored account, expect error because it already exists in the keybase
	_, err = cstore.ImportPrivKey(armoredKey, passphrase, passphrase)
	require.Error(t, err)

	// Delete the account, because otherwise it would error out
	err = cstore.Delete(kp.GetAddress(), passphrase)
	require.NoError(t, err)

	// Import the account armored
	importedKp, err := cstore.ImportPrivKey(armoredKey, passphrase, passphrase)
	require.NoError(t, err)
	fetchedKp, err := cstore.Get(importedKp.GetAddress())
	require.NoError(t, err)
	require.Equal(t, fetchedKp, importedKp)
}

func TestRawExportImport(t *testing.T) {
	// make the storage with reasonable defaults
	cstore := NewInMemory()

	// Create an account
	passphrase := "1234"
	kp, err := cstore.Create(passphrase)
	require.NoError(t, err)

	// Export the raw account
	rawPk, err := cstore.ExportPrivateKeyObject(kp.GetAddress(), passphrase)
	require.NoError(t, err)
	require.NotEmpty(t, rawPk)
	require.Equal(t, kp.PubKey.Address().String(), rawPk.PubKey().Address().String())
	kpList, err := cstore.List()
	require.NoError(t, err)
	require.NotEmpty(t, kpList)
	_, err = cstore.ImportPrivateKeyObject(rawPk.(ed25519.PrivKeyEd25519), passphrase)
	require.Error(t, err)

	// Delete the account, because otherwise it would error out
	err = cstore.Delete(kp.GetAddress(), passphrase)
	require.NoError(t, err)

	// Import the raw account succesfully
	importedKp, err := cstore.ImportPrivateKeyObject(rawPk.(ed25519.PrivKeyEd25519), passphrase)
	require.NoError(t, err)
	fetchedKp, err := cstore.Get(importedKp.GetAddress())
	require.NoError(t, err)
	require.Equal(t, fetchedKp, importedKp)
}

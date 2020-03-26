package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/gov/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModifyParam(t *testing.T) {
	addr := getRandomValidatorAddress()
	var aclKey = types.NewACLKey(types.ModuleName, string(types.DAOOwnerKey))
	ctx, k := createTestKeeperAndContext(t, false)
	res := k.ModifyParam(ctx, aclKey, addr, k.GetACL(ctx).GetOwner(aclKey))
	assert.Zero(t, res.Code)
	s, ok := k.GetSubspace(types.DefaultParamspace)
	assert.True(t, ok)
	var b sdk.Address
	s.Get(ctx, []byte("daoOwner"), &b)
	assert.Equal(t, addr, b)
}

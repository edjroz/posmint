package keeper

import (
	"fmt"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/pos/exported"
	"github.com/pokt-network/posmint/x/pos/types"
)

// get the staked validators from the previous state
func (k Keeper) getValsFromPrevState(ctx sdk.Context) (validators []types.Validator) {
	store := ctx.KVStore(k.storeKey)
	maxValidators := k.MaxValidators(ctx)
	validators = make([]types.Validator, maxValidators)
	iterator := sdk.KVStorePrefixIterator(store, types.PrevStateValidatorsPowerKey)
	defer iterator.Close()
	i := 0
	for ; iterator.Valid(); iterator.Next() {
		if i >= int(maxValidators) {
			panic("more validators than maxValidators found")
		}
		address := types.AddressFromPrevStateValidatorPowerKey(iterator.Key())
		validator := k.mustGetValidator(ctx, address)
		validators[i] = validator
		i++
	}
	return validators[:i] // trim
}

// Load the prevState total validator power.
func (k Keeper) PrevStateValidatorsPower(ctx sdk.Context) (power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.PrevStateTotalPowerKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &power)
	return
}

// Set the prevState total validator power (used in moving the curr to prev)
func (k Keeper) SetPrevStateValidatorsPower(ctx sdk.Context, power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(power)
	store.Set(types.PrevStateTotalPowerKey, b)
}

// returns an iterator for the consensus validators in the prevState block
func (k Keeper) prevStateValidatorsIterator(ctx sdk.Context) (iterator sdk.Iterator) {
	store := ctx.KVStore(k.storeKey)
	iterator = sdk.KVStorePrefixIterator(store, types.PrevStateValidatorsPowerKey)
	return iterator
}

// Iterate over prevState validator powers and perform a function on each validator.
func (k Keeper) IterateAndExecuteOverPrevStateValsByPower(
	ctx sdk.Context, handler func(address sdk.ValAddress, power int64) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.PrevStateValidatorsPowerKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		addr := sdk.ValAddress(iter.Key()[len(types.PrevStateValidatorsPowerKey):])
		var power int64
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &power)
		if handler(addr, power) {
			break
		}
	}
}

// iterate through the active validator set and perform the provided function
func (k Keeper) IterateAndExecuteOverPrevStateVals(
	ctx sdk.Context, fn func(index int64, validator exported.ValidatorI) (stop bool)) {
	iterator := k.prevStateValidatorsIterator(ctx)
	defer iterator.Close()
	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		address := types.AddressFromPrevStateValidatorPowerKey(iterator.Key())
		validator, found := k.GetValidator(ctx, address)
		if !found {
			panic(fmt.Sprintf("validator record not found for address: %v\n", address))
		}
		stop := fn(i, validator) // XXX is this safe will the validator unexposed fields be able to get written to?
		if stop {
			break
		}
		i++
	}
}

// get the power of a SINGLE staked validator from the previous state
func (k Keeper) PrevStateValidatorPower(ctx sdk.Context, addr sdk.ValAddress) (power int64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyForValidatorPrevStateStateByPower(addr))
	if bz == nil {
		return 0
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &power)
	return
}

// set the power of a SINGLE staked validator from the previous state
func (k Keeper) SetPrevStateValPower(ctx sdk.Context, addr sdk.ValAddress, power int64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(power)
	store.Set(types.KeyForValidatorPrevStateStateByPower(addr), bz)
}

// Delete the power of a SINGLE staked validator from the previous state
func (k Keeper) DeletePrevStateValPower(ctx sdk.Context, addr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForValidatorPrevStateStateByPower(addr))
}

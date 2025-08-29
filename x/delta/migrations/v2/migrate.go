package v2

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/delta-chain/delta/v2/x/delta/exported"
	"github.com/delta-chain/delta/v2/x/delta/types"
)

// Migrate migrates the x/delta module state from the consensus version 1 to
// version 2. Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the x/delta
// module state.
func Migrate(ctx sdk.Context, store storetypes.KVStore, legacySubspace exported.Subspace, cdc codec.BinaryCodec) error {
	var currParams types.Params
	legacySubspace.GetParamSetIfExists(ctx, &currParams)

	if err := currParams.Validate(); err != nil {
		return err
	}
	if currParams.GetMaxCallbackGas() == 0 {
		currParams.MaxCallbackGas = types.MaxCallbackGasDefaultValue
	}
	bz := cdc.MustMarshal(&currParams)
	store.Set(types.ParamsKey, bz)

	return nil
}

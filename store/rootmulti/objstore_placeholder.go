//go:build !objstore
// +build !objstore

package rootmulti

import (
	"fmt"

	"cosmossdk.io/store/types"
	"github.com/delta-chain/delta/memiavl"
)

func (rs *Store) loadExtraStore(db *memiavl.DB, key types.StoreKey, params storeParams) (types.CommitStore, error) {
	panic(fmt.Sprintf("unrecognized store type %v", params.typ))
}

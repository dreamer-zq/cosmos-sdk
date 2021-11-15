package v016

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft/example/nft/exported"
	"github.com/cosmos/cosmos-sdk/x/nft/example/nft/types"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	migrateNFTs(store, cdc)
	return nil
}

func migrateNFTs(store sdk.KVStore, cdc codec.BinaryCodec) error {
	var denoms []types.Denom
	iterator := sdk.KVStorePrefixIterator(store, types.KeyDenomID(""))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var denom types.Denom
		cdc.MustUnmarshal(iterator.Value(), &denom)
		denoms = append(denoms, denom)
	}
	return nil
}

// GetNFTs returns all NFTs by the specified denom ID
func GetNFTs(ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	denom string) (nfts []exported.NFT) {
	store := ctx.KVStore(storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.KeyNFT(denom, ""))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var baseNFT types.BaseNFT
		cdc.MustUnmarshal(iterator.Value(), &baseNFT)
		nfts = append(nfts, baseNFT)
	}
	return nfts
}

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/x/nft/example/nft/types"
)

// SetCollection saves all NFTs and returns an error if there already exists
func (k Keeper) SetCollection(ctx sdk.Context, collection types.Collection) error {
	for _, nft := range collection.NFTs {
		if err := k.MintNFT(
			ctx,
			collection.Denom.Id,
			nft.GetID(),
			nft.GetName(),
			nft.GetURI(),
			nft.GetData(),
			nft.GetOwner(),
		); err != nil {
			return err
		}
	}
	return nil
}

// GetCollection returns the collection by the specified denom ID
func (k Keeper) GetCollection(ctx sdk.Context, denomID string) (types.Collection, error) {
	denom, found := k.GetDenom(ctx, denomID)
	if !found {
		return types.Collection{}, sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s not existed ", denomID)
	}

	nfts, err := k.GetNFTs(ctx, denomID)
	if err != nil {
		return types.Collection{}, err
	}
	return types.NewCollection(denom, nfts), nil
}

// GetCollections returns all the collections
func (k Keeper) GetCollections(ctx sdk.Context) (cs []types.Collection, err error) {
	for _, denom := range k.GetDenoms(ctx) {
		nfts, err := k.GetNFTs(ctx, denom.Id)
		if err != nil {
			return cs, err
		}
		cs = append(cs, types.NewCollection(denom, nfts))
	}
	return cs, nil
}

// GetTotalSupply returns the number of NFTs by the specified denom ID
func (k Keeper) GetTotalSupply(ctx sdk.Context, denomID string) uint64 {
	return k.nk.GetTotalSupply(ctx, denomID)
}

// GetTotalSupplyOfOwner returns the amount of NFTs by the specified conditions
func (k Keeper) GetTotalSupplyOfOwner(ctx sdk.Context, id string, owner sdk.AccAddress) (supply uint64) {
	return k.nk.GetBalance(ctx, id, owner)
}

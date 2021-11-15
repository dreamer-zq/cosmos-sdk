package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gogotypes "github.com/gogo/protobuf/types"

	"github.com/cosmos/cosmos-sdk/x/nft/example/nft/exported"
	"github.com/cosmos/cosmos-sdk/x/nft/example/nft/types"
)

// GetNFT gets the the specified NFT
func (k Keeper) GetNFT(ctx sdk.Context, denomID, tokenID string) (nft exported.NFT, err error) {
	token, exist := k.nk.GetNFT(ctx, denomID, tokenID)
	if !exist {
		return nil, sdkerrors.Wrapf(types.ErrUnknownNFT, "not found NFT: %s", denomID)
	}
	owner := k.nk.GetOwner(ctx, denomID, tokenID)
	return types.BaseNFT{
		Id:    token.GetId(),
		Name:  "",
		URI:   token.GetUri(),
		Data:  "",
		Owner: owner.String(),
	}, nil
}

// GetNFTs returns all NFTs by the specified denom ID
func (k Keeper) GetNFTs(ctx sdk.Context, denom string) (nfts []exported.NFT, err error) {
	tokens := k.nk.GetNFTsOfClass(ctx, denom)
	for _, token := range tokens {
		var data = &gogotypes.StringValue{}
		if err := k.cdc.Unmarshal(token.GetData().Value, data); err != nil {
			return nil, err
		}
		nfts = append(nfts, types.BaseNFT{
			Id:    token.GetId(),
			Name:  "",
			URI:   token.GetUri(),
			Data:  data.GetValue(),
			Owner: k.nk.GetOwner(ctx, denom, token.GetClassId()).String(),
		})
	}
	return nfts, nil
}

// Authorize checks if the sender is the owner of the given NFT
// Return the NFT if true, an error otherwise
func (k Keeper) Authorize(ctx sdk.Context, denomID, tokenID string, owner sdk.AccAddress) error {
	if !owner.Equals(k.nk.GetOwner(ctx, denomID, tokenID)) {
		return sdkerrors.Wrap(types.ErrUnauthorized, owner.String())
	}
	return nil
}

// HasNFT checks if the specified NFT exists
func (k Keeper) HasNFT(ctx sdk.Context, denomID, tokenID string) bool {
	return k.nk.HasNFT(ctx, denomID, tokenID)
}

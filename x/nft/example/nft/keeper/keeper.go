package keeper

import (
	"fmt"
	"reflect"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gogotypes "github.com/gogo/protobuf/types"

	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/cosmos/cosmos-sdk/x/nft/example/nft/types"
	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey storetypes.StoreKey // Unexposed key to access store from sdk.Context
	cdc      codec.Codec
	nk       nftkeeper.Keeper
}

// NewKeeper creates a new instance of the NFT Keeper
func NewKeeper(cdc codec.Codec, storeKey storetypes.StoreKey, ak nft.AccountKeeper, bk nft.BankKeeper) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
		nk:       nftkeeper.NewKeeper(storeKey, cdc, ak, bk),
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("irismod/%s", types.ModuleName))
}

// IssueDenom issues a denom according to the given params
func (k Keeper) IssueDenom(ctx sdk.Context,
	id, name, schema, symbol string,
	creator sdk.AccAddress,
	mintRestricted, updateRestricted bool,
) error {
	if err := k.nk.SaveClass(ctx, nft.Class{
		Id:          id,
		Name:        name,
		Symbol:      schema,
		Description: "",
		Uri:         "",
		UriHash:     "",
		Data:        nil,
	}); err != nil {
		return err
	}
	return k.SetDenom(ctx, types.NewDenom(id, name, schema, symbol, creator, mintRestricted, updateRestricted))
}

// MintNFT mints an NFT and manages the NFT's existence within Collections and Owners
func (k Keeper) MintNFT(
	ctx sdk.Context, denomID, tokenID, tokenNm,
	tokenURI, tokenData string, owner sdk.AccAddress,
) error {
	data, err := codectypes.NewAnyWithValue(&gogotypes.StringValue{Value: tokenData})
	if err != nil {
		return err
	}
	return k.nk.Mint(ctx, nft.NFT{
		ClassId: denomID,
		Id:      tokenID,
		Uri:     tokenURI,
		UriHash: "",
		Data:    data,
	}, owner)
}

// EditNFT updates an already existing NFT
func (k Keeper) EditNFT(
	ctx sdk.Context, denomID, tokenID, tokenNm,
	tokenURI, tokenData string, owner sdk.AccAddress,
) error {
	denom, found := k.GetDenom(ctx, denomID)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denom ID %s not exists", denomID)
	}

	if denom.UpdateRestricted {
		// if true , nobody can update the NFT under this denom
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "nobody can update the NFT under this denom %s", denom.Id)
	}

	// just the owner of NFT can edit
	if err := k.Authorize(ctx, denomID, tokenID, owner); err != nil {
		return err
	}

	// if types.Modified(tokenNm) {
	// 	token.Name = tokenNm
	// } // TODO
	token, exist := k.nk.GetNFT(ctx, denomID, tokenID)
	if !exist {
		return sdkerrors.Wrapf(types.ErrUnknownNFT, "nft ID %s not exists", tokenID)
	}
	if types.Modified(tokenURI) {
		token.Uri = tokenURI
	}

	if types.Modified(tokenData) {
		metadata, err := codectypes.NewAnyWithValue(&gogotypes.StringValue{Value: tokenData})
		if err != nil {
			return err
		}
		token.Data = metadata
	}
	return k.nk.Update(ctx, token)
}

// TransferOwner transfers the ownership of the given NFT to the new owner
func (k Keeper) TransferOwner(
	ctx sdk.Context, denomID, tokenID, tokenNm, tokenURI,
	tokenData string, srcOwner, dstOwner sdk.AccAddress,
) error {
	if err := k.Authorize(ctx, denomID, tokenID, srcOwner); err != nil {
		return err
	}
	token, exist := k.nk.GetNFT(ctx, denomID, tokenID)
	if !exist {
		return sdkerrors.Wrapf(types.ErrInvalidTokenID, "nft ID %s not exists", tokenID)
	}

	data, err := codectypes.NewAnyWithValue(&gogotypes.StringValue{Value: tokenData})
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(token, nft.NFT{
		ClassId: denomID,
		Id:      tokenID,
		Uri:     tokenURI,
		Data:    data,
	}) {
		if err := k.EditNFT(ctx, denomID, tokenID, tokenNm, tokenURI, tokenData, srcOwner); err != nil {
			return err
		}
	}
	return k.nk.Transfer(ctx, denomID, tokenID, dstOwner)
}

// BurnNFT deletes a specified NFT
func (k Keeper) BurnNFT(ctx sdk.Context, denomID, tokenID string, owner sdk.AccAddress) error {
	if err := k.Authorize(ctx, denomID, tokenID, owner); err != nil {
		return err
	}
	return k.nk.Burn(ctx, denomID, tokenID)
}

// TransferDenomOwner transfers the ownership of the given denom to the new owner
func (k Keeper) TransferDenomOwner(
	ctx sdk.Context, denomID string, srcOwner, dstOwner sdk.AccAddress,
) error {
	denom, found := k.GetDenom(ctx, denomID)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denom ID %s not exists", denomID)
	}

	// authorize
	if srcOwner.String() != denom.Creator {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to transfer denom %s", srcOwner.String(), denomID)
	}

	denom.Creator = dstOwner.String()

	err := k.UpdateDenom(ctx, denom)
	if err != nil {
		return err
	}

	return nil
}

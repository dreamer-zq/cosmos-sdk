package v160

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	PrefixNFT        = []byte{0x01}
	PrefixOwners     = []byte{0x02} // key for a owner
	PrefixCollection = []byte{0x03} // key for balance of NFTs held by the denom

	delimiter = []byte("/")
)

type (
	IssueDenomFn = func(
		ctx sdk.Context,
		id, name, schema, symbol string,
		creator sdk.AccAddress,
		mintRestricted, updateRestricted bool,
	) error

	MintNFTFn = func(
		ctx sdk.Context, denomID, tokenID, tokenNm,
		tokenURI, tokenData string, owner sdk.AccAddress,
	) error
)

func keyCollection(denomID string) []byte {
	key := append(PrefixCollection, delimiter...)
	return append(key, []byte(denomID)...)
}

// KeyOwner gets the key of a collection owned by an account address
func keyOwner(address sdk.AccAddress, denomID, tokenID string) []byte {
	key := append(PrefixOwners, delimiter...)
	if address != nil {
		key = append(key, []byte(address.String())...)
		key = append(key, delimiter...)
	}

	if address != nil && len(denomID) > 0 {
		key = append(key, []byte(denomID)...)
		key = append(key, delimiter...)
	}

	if address != nil && len(denomID) > 0 && len(tokenID) > 0 {
		key = append(key, []byte(tokenID)...)
	}
	return key
}

// KeyNFT gets the key of nft stored by an denom and id
func keyNFT(denomID, tokenID string) []byte {
	key := append(PrefixNFT, delimiter...)
	if len(denomID) > 0 {
		key = append(key, []byte(denomID)...)
		key = append(key, delimiter...)
	}

	if len(denomID) > 0 && len(tokenID) > 0 {
		key = append(key, []byte(tokenID)...)
	}
	return key
}

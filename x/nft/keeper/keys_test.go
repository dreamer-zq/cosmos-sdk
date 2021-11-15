package keeper

import (
	"bytes"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestNftOfClassByOwnerStoreKey(t *testing.T) {
	address, _ := sdk.AccAddressFromBech32("cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47")
	denom1 := "denomid"
	denom2 := "denomid2"

	key1 := nftOfClassByOwnerStoreKey(address, denom1)
	fmt.Println(key1)
	key2 := nftOfClassByOwnerStoreKey(address, denom2)
	fmt.Println(key2)

	fmt.Println(bytes.Equal(key1, key2))
}

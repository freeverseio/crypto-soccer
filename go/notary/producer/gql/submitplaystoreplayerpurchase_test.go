package gql_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/awa/go-iap/playstore"
	"gotest.tools/assert"
)

func TestProva(t *testing.T) {
	t.Skip("Runme to check the purchase communication")
	packageName := "com.freeverse.phoenix"
	productID := "coinpack_45"

	credentials, err := ioutil.ReadFile("../../../../../key.json")
	assert.NilError(t, err)

	token := "hjjfppagdilpbmnmjaajpcgc.AO-J1Owne5VpLZzOtfFZkDY1k5T4kEXCgack0gmEYssqCgEYzlNgPtHdp72TPELzOl3T8XCYhc0k818EbCi7hiYcEgbCGNVyNCGd1I2wdz9pxGRHXs1-msWvAD9ztmd11v_hr9NqCSn1"
	client, err := playstore.New(credentials)
	assert.NilError(t, err)
	ctx := context.Background()
	result, err := client.VerifyProduct(ctx, packageName, productID, token)
	assert.NilError(t, err)
	t.Log(result)
}

package gql_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/awa/go-iap/playstore"
	"gotest.tools/assert"
)

// func TestSubmitPlayerPurchase(t *testing.T) {
// 	// t.Skip("Rective this test to test google api")
// 	ch := make(chan interface{}, 10)
// 	googleCredentials, err := ioutil.ReadFile("./key.json")
// 	assert.NilError(t, err)
// 	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials)

// 	in := input.SubmitPlayerPurchaseInput{}
// 	in.TeamId = "274877906944"
// 	in.PlayerId = "274877906944"
// 	in.PurchaseId = "GPA.3309-9448-9453-66100"
// 	in.Signature = "91366deb26195ac3b15b9e6fff99d425b2cc0d15d44dc8ee0377779400f92c4358a57754053facbe724e8a536e240b278cd651c756c46978eaebafc47767fd781b"

// 	id, err := r.SubmitPlayerPurchase(struct {
// 		Input input.SubmitPlayerPurchaseInput
// 	}{in})
// 	assert.NilError(t, err)
// 	assert.Equal(t, id, in.PlayerId)
// }

func TestProva(t *testing.T) {
	t.Skip("Runme to check the purchase communication")
	packageName := "com.freeverse.phoenix"
	productID := "coinpack_45"

	credentials, err := ioutil.ReadFile("./key.json")
	assert.NilError(t, err)

	token := "hjjfppagdilpbmnmjaajpcgc.AO-J1Owne5VpLZzOtfFZkDY1k5T4kEXCgack0gmEYssqCgEYzlNgPtHdp72TPELzOl3T8XCYhc0k818EbCi7hiYcEgbCGNVyNCGd1I2wdz9pxGRHXs1-msWvAD9ztmd11v_hr9NqCSn1"
	client, err := playstore.New(credentials)
	assert.NilError(t, err)
	ctx := context.Background()
	result, err := client.VerifyProduct(ctx, packageName, productID, token)
	assert.NilError(t, err)
	t.Log(result)
}

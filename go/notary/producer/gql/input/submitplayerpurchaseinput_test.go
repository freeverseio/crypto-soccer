package input_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"gotest.tools/assert"
)

// func TestInAppPurchase(t *testing.T) {
// 	// You need to prepare a public key for your Android app's in app billing
// 	// at https://console.developers.google.com.
// 	jsonKey, err := ioutil.ReadFile("jsonKey.json")
// 	assert.NilError(t, err)

// 	client, err := playstore.New(jsonKey)
// 	assert.NilError(t, err)
// 	ctx := context.Background()
// 	_, err = client.VerifyProduct(ctx, input.GooglePackage, input.GoogleProductID, "tocken")
// 	assert.NilError(t, err)
// }

func TestSubmitPlayerPurchaseInputHash(t *testing.T) {
	in := input.SubmitPlayerPurchaseInput{}

	hash, err := in.Hash()
	assert.Error(t, err, "Invalid TeamId")

	in.TeamId = "3"
	hash, err = in.Hash()
	assert.Error(t, err, "Invalid PlayerId")

	in.PlayerId = "5"
	hash, err = in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x65732122af23f04532e927990a9fe25d7ba5663403d037be02968b9e391f1446")
}

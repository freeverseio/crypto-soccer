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

	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19b")
}

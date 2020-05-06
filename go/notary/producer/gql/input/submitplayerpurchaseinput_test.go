package input_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/awa/go-iap/playstore"
	"gotest.tools/assert"
)

func TestInAppPurchase(t *testing.T) {
	// You need to prepare a public key for your Android app's in app billing
	// at https://console.developers.google.com.
	jsonKey, err := ioutil.ReadFile("jsonKey.json")
	assert.NilError(t, err)

	client, err := playstore.New(jsonKey)
	assert.NilError(t, err)
	ctx := context.Background()
	_, err = client.VerifyProduct(ctx, "package", "productID", "tocken")
	assert.NilError(t, err)
}

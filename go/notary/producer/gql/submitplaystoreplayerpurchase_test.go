package gql_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/awa/go-iap/playstore"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/androidpublisher/v3"
	"gotest.tools/assert"
)

const packageName = "com.freeverse.phoenix"
const productID = "coinpack_45"
const token = "hjjfppagdilpbmnmjaajpcgc.AO-J1Owne5VpLZzOtfFZkDY1k5T4kEXCgack0gmEYssqCgEYzlNgPtHdp72TPELzOl3T8XCYhc0k818EbCi7hiYcEgbCGNVyNCGd1I2wdz9pxGRHXs1-msWvAD9ztmd11v_hr9NqCSn1"

func TestAlex(t *testing.T) {
	data, err := ioutil.ReadFile("../../../../../key.json")
	assert.NilError(t, err)

	conf, err := google.JWTConfigFromJSON(data, androidpublisher.AndroidpublisherScope)
	assert.NilError(t, err)

	client := conf.Client(oauth2.NoContext)
	svc, err := androidpublisher.New(client)
	assert.NilError(t, err)

	ps := androidpublisher.NewPurchasesProductsService(svc)
	_, err = ps.Get(packageName, productID, token).Context(context.Background()).Do()
	assert.NilError(t, err)
}

func TestProva(t *testing.T) {
	credentials, err := ioutil.ReadFile("../../../../../key.json")
	assert.NilError(t, err)

	client, err := playstore.New(credentials)
	assert.NilError(t, err)

	_, err = client.VerifyProduct(context.Background(), packageName, productID, token)
	assert.NilError(t, err)
}

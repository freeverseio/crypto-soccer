package authproxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type ServerService interface {
	CountTeams(ctx context.Context, owner common.Address) (int, error)
	NewRequest(method string, body io.Reader) (*http.Request, error)
}

type GraphQLServerService struct {
	gqlurl string
}

func NewGraphQLServerService(url string) *GraphQLServerService {
	return &GraphQLServerService{
		gqlurl: url,
	}
}

func (b GraphQLServerService) NewRequest(method string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, b.gqlurl, body)
}

func (b GraphQLServerService) CountTeams(ctx context.Context, owner common.Address) (int, error) {
	gqlQuery := `{allTeams (condition: {owner: "` + owner.Hex() + `"}){totalCount}}`
	query, err := json.Marshal(map[string]string{"query": gqlQuery})
	if err != nil {
		return 0, errors.Wrap(err, "failed bulding auth query")
	}
	req, err := http.NewRequest(http.MethodPost, b.gqlurl, bytes.NewReader(query))
	if err != nil {
		return 0, errors.Wrap(err, "failed bulding auth request")
	}
	req.Header.Set("Content-Type", "application/json")

	// send request
	client := &http.Client{}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return 0, errors.Wrap(err, "failed sending auth request")
	}

	// check http response is ok
	if resp.StatusCode != 200 {
		errstr := "unknown"
		if resp.Body != nil {
			errbody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				errstr = string(errbody)
			}
		}
		return 0, fmt.Errorf("failed sending auth request, errcode=%v, err=%s", resp.StatusCode, errstr)
	}

	// parse qgl response, and return
	var response struct {
		Data struct {
			AllTeams struct {
				TotalCount int
			}
		}
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, errors.Wrap(err, "failed decoding auth response")
	}

	return response.Data.AllTeams.TotalCount, nil
}

type MockServerService struct {
	CountTeamFn  func() (int, error)
	NewRequestFn func() (*http.Request, error)
}

func (b MockServerService) CountTeams(ctx context.Context, owner common.Address) (int, error) {
	return b.CountTeamFn()
}

func (b MockServerService) NewRequest(method string, body io.Reader) (*http.Request, error) {
	return b.NewRequestFn()
}

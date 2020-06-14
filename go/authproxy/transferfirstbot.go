package authproxy

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func IsTransferFirstBotMutation(data string) (bool, error) {
	log.Debugf("Check matching on: %v", data)
	// return strings.Contains(data, "transferFirstBotToAddr"), nil
	ex := `(\s*)mutation(\s*)TransferTeamToPlayer\(\$timezoneIdx:(\s*)Int!,(\s*)\$countryIdx:(\s*)ID!,(\s*)\$address:(\s*)String!\)(\s*){(\s*)transferFirstBotToAddr\((\s*)timezone:(\s*)\$timezoneIdx(\s*)countryIdxInTimezone:(\s*)\$countryIdx(\s*)address:(\s*)\$address(\s*)\)(\s*)}(\s*)`
	// ex := `mutation(\s*)TransferTeamToPlayer\(\$timezoneIdx:(\s*)Int!,(\s*)\$countryIdx:(\s*)ID!,(\s*)\$address:(\s*)String!\)(\s*){(\s*)transferFirstBotToAddr\((\s*)timezone:(\s*)\$timezoneIdx(\s*)countryIdxInTimezone:(\s*)\$countryIdx(\s*)address:(\s*)\$address(\s*)\)(\s*)}\\",(\s*)`
	// ex := `mutation(\s*).*(\s*){(\s*)transferFirstBotToAddr(\s*)\((\s*)timezone(\s*):(\s*)\d{1,2}(\s*),(\s*)countryIdxInTimezone(\s*):(\s*)[0-9]+(\s*),(\s*)address(\s*):(\s*)"[a-zA-Z0-9]+"(\s*)\)(\s*)}`
	return regexp.MatchString(ex, data)
}

func MatchTransferFirstBotMutation(r *http.Request) (bool, error) {
	if r == nil {
		return false, errors.New("nil request")
	}
	if r.Body == http.NoBody {
		return false, nil
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false, errors.Wrap(err, "failed reading the body")
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	var query struct {
		Data string `json:"query"`
	}
	err = json.Unmarshal(body, &query)
	if err != nil {
		return false, err
	}
	return IsTransferFirstBotMutation(query.Data)
}

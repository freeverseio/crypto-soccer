package authproxy

import (
	"encoding/json"
	"net/http"
	"regexp"
)

func IsTransferFirstBotMutation(data string) (bool, error) {
	ex := `mutation(\s*).*(\s*){(\s*)transferFirstBotToAddr(\s*)\((\s*)timezone(\s*):(\s*)\d{1,2}(\s*),(\s*)countryIdxInTimezone(\s*):(\s*)[0-9]+(\s*),(\s*)address(\s*):(\s*)"[a-zA-Z0-9]+"(\s*)\)(\s*)}`
	return regexp.MatchString(ex, data)
}

func matchTransferFirstBotMutation(r *http.Request) (bool, error) {
	var query struct {
		Data string `json:"query"`
	}
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		return false, err
	}
	return IsTransferFirstBotMutation(query.Data)
}

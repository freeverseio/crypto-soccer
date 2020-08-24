package orgmap

import "github.com/freeverseio/crypto-soccer/go/storage"

type OrgMapDenyList struct {
	TeamID string `json:"team_id"` // team_id
}

func Map(vs []storage.Team, f func(storage.Team) OrgMapDenyList) []OrgMapDenyList {
	vsm := make([]OrgMapDenyList, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

package mock

type TeamStorageService struct {
	UpdateNameFn func(teamId string, name string) error
}

func (b TeamStorageService) UpdateName(teamId string, name string) error {
	return b.UpdateNameFn(teamId, name)
}

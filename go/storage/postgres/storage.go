package postgres

type StorageService struct {
	TeamStorageService  TeamStorageService
	MatchStorageService MatchStorageService
}

func NewStorageService() *StorageService {
	return &StorageService{
		*NewTeamStorageService(),
		*NewMatchStorageService(),
	}
}

func (b StorageService) TeamService() TeamStorageService {
	return b.TeamStorageService
}
func (b StorageService) MatchService() MatchStorageService {
	return b.MatchStorageService
}

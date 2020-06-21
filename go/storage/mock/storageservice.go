package mock

import "github.com/freeverseio/crypto-soccer/go/storage"

type StorageService struct {
	TeamStorageService  TeamStorageService
	MatchStorageService MatchStorageService
}

func NewStorageService() *StorageService {
	return &StorageService{
		TeamStorageService{},
		MatchStorageService{},
	}
}

func (b StorageService) TeamService() storage.TeamStorageService {
	return b.TeamStorageService
}
func (b StorageService) MatchService() storage.MatchStorageService {
	return b.MatchStorageService
}

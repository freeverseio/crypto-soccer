package memory

import (
	"github.com/freeverseio/crypto-soccer/go/storage"
)

func NewStorageService() *storage.StorageService {
	return &storage.StorageService{
		NewTeamStorageService(),
		NewMatchStorageService(),
	}
}

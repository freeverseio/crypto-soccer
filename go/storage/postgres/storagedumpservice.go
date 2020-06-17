package postgres

import log "github.com/sirupsen/logrus"

type StorageDumpService struct {
}

func (b StorageDumpService) Dump(fileName string) error {
	log.Warning("[postgres/storagedumpservice] not implemed")
	return nil
}

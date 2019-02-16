package storage

type GlobalsEntry struct {
	CurrentQuota uint
}

type SavePointEntry struct {
	LastBlock    uint64
	LastTxIndex  uint
	LastLogIndex uint
}

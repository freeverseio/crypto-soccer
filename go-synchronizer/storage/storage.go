package storage

type Storage interface {
	TeamAdd(ID uint64, name string) error
	TeamCount() (uint64,error)
}

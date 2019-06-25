package storage

type Team struct {
	Id   uint64
	Name string
}

type Storage interface {
	TeamAdd(ID uint64, name string) error
	TeamCount() (uint64, error)
}

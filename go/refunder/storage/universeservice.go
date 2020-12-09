package storage

type UniverseService interface {
	Begin() (UniverseTx, error)
}

type UniverseTx interface {
}

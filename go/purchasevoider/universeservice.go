package purchasevoider

type UniverseService interface {
	Begin() (UniverseTx, error)
}

type UniverseTx interface {
}

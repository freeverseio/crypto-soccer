package purchasevoider

type UniverseService interface {
	MarkForDeletion(id string) error
}

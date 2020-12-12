package mockup

type UniverseService struct {
	MarkForDeletionFn func(id string) error
}

func (b *UniverseService) MarkForDeletion(id string) error {
	return b.MarkForDeletionFn(id)
}

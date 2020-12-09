package market

type Service interface {
	Begin() (Tx, error)
}

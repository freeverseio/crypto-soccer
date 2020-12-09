package universe

type Service interface {
	Begin() (Tx, error)
}

package gql

type Resolver struct {
	c chan TransferFirstBotToAddrInput
}

func NewResolver(c chan TransferFirstBotToAddrInput) *Resolver {
	return &Resolver{c}
}

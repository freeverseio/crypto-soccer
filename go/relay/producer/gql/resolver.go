package gql

type Resolver struct {
	ch chan interface{}
}

func NewResolver(ch chan interface{}) *Resolver {
	return &Resolver{ch}
}

func (b *Resolver) Ping() bool {
	return true
}

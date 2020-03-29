package gql

type Resolver struct {
	c chan interface{}
}

func NewResolver(c chan interface{}) *Resolver {
	return &Resolver{c}
}

package gql

import (
	"errors"
	"strconv"
)

// UUID represents GraphQL's "UUID" scalar type. A custom type may be used instead.
type Sign string

func (Sign) ImplementsGraphQLType(name string) bool {
	return name == "Sign"
}

func (b *Sign) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		*b = Sign(input)
	default:
		err = errors.New("wrong type")
	}
	return err
}

func (b Sign) MarshalJSON() ([]byte, error) {
	return strconv.AppendQuote(nil, string(b)), nil
}

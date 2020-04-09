package gql

import (
	"errors"
	"strconv"
)

// UUID represents GraphQL's "UUID" scalar type. A custom type may be used instead.
type UUID string

func (UUID) ImplementsGraphQLType(name string) bool {
	return name == "UUID"
}

func (id *UUID) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		*id = UUID(input)
	default:
		err = errors.New("wrong type")
	}
	return err
}

func (id UUID) MarshalJSON() ([]byte, error) {
	return strconv.AppendQuote(nil, string(id)), nil
}

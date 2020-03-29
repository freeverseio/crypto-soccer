package gql

import (
	"net/http"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	log "github.com/sirupsen/logrus"
)

func NewServer(c chan interface{}) error {
	log.Info("New GraphQL server staring ...")

	resolver := NewResolver(c)
	schema := graphql.MustParseSchema(Schema, resolver)

	handler := relay.Handler{Schema: schema}

	http.Handle("/graphql/", &handler)
	http.Handle("/graphql", &handler) // Register without a trailing slash to avoid redirect.
	return http.ListenAndServe(":4000", nil)
}

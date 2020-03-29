package gql

import (
	"net/http"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	log "github.com/sirupsen/logrus"
)

func NewServer() error {
	log.Info("New GraphQL server staring ...")

	schema := graphql.MustParseSchema(Schema, NewResolver())

	handler := relay.Handler{Schema: schema}

	http.Handle("/graphql/", &handler)
	http.Handle("/graphql", &handler) // Register without a trailing slash to avoid redirect.
	return http.ListenAndServe(":4000", nil)
}

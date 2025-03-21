package gql

import (
	"database/sql"
	"net/http"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	log "github.com/sirupsen/logrus"
)

func NewServer(
	ch chan interface{},
	contracts contracts.Contracts,
	db *sql.DB,
) error {
	log.Info("New GraphQL server staring ...")

	resolver := NewResolver(ch, contracts, db)
	schema := graphql.MustParseSchema(Schema, resolver)

	handler := relay.Handler{Schema: schema}

	http.Handle("/graphql/", &handler)
	http.Handle("/graphql", &handler) // Register without a trailing slash to avoid redirect.
	return http.ListenAndServe(":4000", nil)
}

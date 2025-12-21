package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"gorm.io/gorm"
)

var (
	schema   *graphql.Schema
	resolver *Resolver
)

// SetupGraphQL 初始化 GraphQL schema 和 handler
func SetupGraphQL(db *gorm.DB) error {
	resolver = NewResolver(db)

	var err error
	schema, err = NewSchema(resolver)
	if err != nil {
		return err
	}

	return nil
}

// GetHandler 返回 GraphQL HTTP handler
func GetHandler() *handler.Handler {
	if schema == nil {
		return nil
	}

	return handler.New(&handler.Config{
		Schema:     schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	})
}

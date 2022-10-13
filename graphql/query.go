package graphql

import (
	"backend-2/api/graphql/field"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

func newQuery(db *gorm.DB) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"quotes": field.NewQuotes(db),
		},
	})
}

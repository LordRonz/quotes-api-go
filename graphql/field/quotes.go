package field

import (
	"backend-2/api/cmd/db/model"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

var quote = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Quote",
		Fields: graphql.Fields{
			"id":        &graphql.Field{Type: graphql.ID},
			"quote":     &graphql.Field{Type: graphql.String},
			"createdAt": &graphql.Field{Type: graphql.String},
			"updatedAt": &graphql.Field{Type: graphql.String},
			"deletedAt": &graphql.Field{Type: graphql.String},
		},
		Description: "Quotes data",
	},
)

func NewQuotes(db *gorm.DB) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(quote),
		Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
			var u []*model.Quote
			if err := db.Find(&u).Error; err != nil {
				// do something
			}

			return u, nil
		},
		Description: "quotes",
	}
}

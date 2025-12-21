package graphql

import (
	"github.com/graphql-go/graphql"
)

// NewSchema 創建並返回 GraphQL schema
func NewSchema(resolver *Resolver) (*graphql.Schema, error) {
	// 定義 User 類型
	userType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "User",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"name": &graphql.Field{
					Type: graphql.NewNonNull(graphql.String),
				},
				"email": &graphql.Field{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
		},
	)

	// 定義查詢類型
	queryType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"users": &graphql.Field{
					Type:        graphql.NewList(userType),
					Description: "獲取所有會員列表",
					Resolve:     resolver.GetUsers,
				},
				"user": &graphql.Field{
					Type:        userType,
					Description: "根據 ID 獲取單個會員",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.Int),
						},
					},
					Resolve: resolver.GetUserByID,
				},
			},
		},
	)

	// 定義 Mutation 類型（用於創建、更新、刪除）
	mutationType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				"createUser": &graphql.Field{
					Type:        userType,
					Description: "創建新會員",
					Args: graphql.FieldConfigArgument{
						"name": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"email": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: resolver.CreateUser,
				},
				"updateUser": &graphql.Field{
					Type:        userType,
					Description: "更新會員信息",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.Int),
						},
						"name": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"email": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: resolver.UpdateUser,
				},
				"deleteUser": &graphql.Field{
					Type:        graphql.Boolean,
					Description: "刪除會員",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.Int),
						},
					},
					Resolve: resolver.DeleteUser,
				},
			},
		},
	)

	// 創建 Schema
	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType,
			Mutation: mutationType,
		},
	)

	return &schema, err
}

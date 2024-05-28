package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// Contains information about a custom node on the Comfy Registry.

type Node struct {
	ent.Schema
}

func (Node) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Unique(),
		field.String("publisher_id").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.String("name").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.String("description").Optional().SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.String("author").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
		field.String("license").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.String("repository_url").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.String("icon_url").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
		field.Strings("tags").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Default([]string{}),
		field.Strings("test_field").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Default([]string{}),
	}
}

func (Node) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (Node) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("publisher", Publisher.Type).Field("publisher_id").Ref("nodes").Required().Unique(),
		edge.To("versions", NodeVersion.Type),
	}
}

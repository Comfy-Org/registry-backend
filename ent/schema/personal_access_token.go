package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type PersonalAccessToken struct {
	ent.Schema
}

func (PersonalAccessToken) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("name").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.String("description").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.String("publisher_id").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.String("token").Sensitive().Unique().SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
	}
}

func (PersonalAccessToken) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("token").Unique(),
	}
}

func (PersonalAccessToken) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("publisher", Publisher.Type).Ref("personal_access_tokens").
			Field("publisher_id").Unique().Required(),
	}
}

func (PersonalAccessToken) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

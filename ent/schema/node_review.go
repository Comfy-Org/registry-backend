package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type NodeReview struct {
	ent.Schema
}

func (NodeReview) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("node_id").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.String("user_id").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.Int("star").Default(0),
	}
}

func (NodeReview) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Field("user_id").Ref("reviews").Unique().Required(),
		edge.From("node", Node.Type).Field("node_id").Ref("reviews").Unique().Required(),
	}
}
func (NodeReview) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("node_id", "user_id").Unique(),
	}
}

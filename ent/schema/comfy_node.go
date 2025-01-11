package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type ComfyNode struct {
	ent.Schema
}

func (ComfyNode) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("name").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.UUID("node_version_id", uuid.UUID{}),
		field.String("category").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
		field.String("description").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
		field.String("input_types").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
		field.Bool("deprecated").Default(false).Optional(),
		field.Bool("experimental").Default(false).Optional(),
		field.JSON("output_is_list", []bool{}).Default([]bool{}).Optional(),
		field.String("return_names").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
		field.String("return_types").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
		field.String("function").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
	}
}

func (ComfyNode) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (ComfyNode) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("versions", NodeVersion.Type).Field("node_version_id").Ref("comfy_nodes").Required().Unique(),
	}
}

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

// Contains information about a specific version of a custom node on the Comfy Registry.

type NodeVersion struct {
	ent.Schema
}

func (NodeVersion) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("node_id").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.String("version").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Comment("Must be SemVer compliant"),
		field.String("changelog").Optional().SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.Strings("pip_dependencies").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.Bool("deprecated").Default(false),
	}
}

func (NodeVersion) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (NodeVersion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("node", Node.Type).Field("node_id").Ref("versions").Required().Unique(),
		edge.To("storage_file", StorageFile.Type).Unique(),
	}
}

func (NodeVersion) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("node_id", "version").Unique(),
	}
}

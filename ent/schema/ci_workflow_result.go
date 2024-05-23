package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// CIWorkflowResult holds the schema definition for the CIWorkflowResult entity.
type CIWorkflowResult struct {
	ent.Schema
}

// Fields of the CIWorkflowResult.
func (CIWorkflowResult) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.String("operating_system").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.String("gpu_type").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
		field.String("pytorch_version").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
		field.String("workflow_name").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
		field.String("run_id").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
		field.String("status").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
		field.Int64("start_time").Optional(),
		field.Int64("end_time").Optional(),
	}
}

// Edges of the CIWorkflowResult.
func (CIWorkflowResult) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("gitcommit", GitCommit.Type).Ref("results").Unique(),
		edge.To("storage_file", StorageFile.Type).Unique(),
	}
}

func (CIWorkflowResult) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

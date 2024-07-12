package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// Stores the artifacts and metadata about a single workflow execution for Comfy CI.
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
		field.String("workflow_name").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
		field.String("run_id").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
		field.String("job_id").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
		field.String("status").GoType(WorkflowRunStatusType("")).Default(string(WorkflowRunStatusTypeStarted)),
		field.Int64("start_time").Optional(),
		field.Int64("end_time").Optional(),
		field.String("python_version").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
		field.String("pytorch_version").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
		field.String("cuda_version").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
		field.String("comfy_run_flags").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
		field.Int("avg_vram").Optional().Comment("Average amount of VRAM used by the workflow in Megabytes"),
		field.Int("peak_vram").Optional().Comment("Peak amount of VRAM used by the workflow in Megabytes"),
		field.String("job_trigger_user").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional().Comment("User who triggered the job"),
		field.JSON("metadata", map[string]interface{}{}).SchemaType(map[string]string{
			dialect.Postgres: "jsonb",
		}).Optional().Comment("Stores miscellaneous metadata for each workflow run."),
	}
}

type WorkflowRunStatusType string

const (
	WorkflowRunStatusTypeCompleted WorkflowRunStatusType = "COMPLETED"
	WorkflowRunStatusTypeFailed    WorkflowRunStatusType = "FAILED"
	WorkflowRunStatusTypeStarted   WorkflowRunStatusType = "STARTED"
)

// Edges of the CIWorkflowResult.
func (CIWorkflowResult) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("gitcommit", GitCommit.Type).Ref("results").Unique(),
		edge.To("storage_file", StorageFile.Type),
	}
}

func (CIWorkflowResult) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

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

// GitCommit holds the schema definition for the GitCommit entity.
type GitCommit struct {
	ent.Schema
}

// Fields of the GitCommit.
func (GitCommit) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("commit_hash").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.String("branch_name").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.String("repo_name").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.String("commit_message").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.Time("commit_timestamp"),
		field.String("author").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
		field.Time("timestamp").Optional(),
		field.String("pr_number").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
	}
}

func (GitCommit) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("results", CIWorkflowResult.Type),
	}
}

func (GitCommit) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (GitCommit) Indexes() []ent.Index {
	// Hashes within the same repo should be unique.
	return []ent.Index{
		index.Fields("repo_name", "commit_hash").
			Unique(),
	}
}

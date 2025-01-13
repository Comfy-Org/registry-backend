package schema

import (
	"fmt"

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
		field.Enum("status").
			GoType(NodeVersionStatus("")).
			Default(string(NodeVersionStatusPending)),
		field.String("status_reason").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Default("").Comment("Give a reason for the status change. Eg. 'Banned due to security vulnerability'"),
		field.String("comfy_node_extract_status").
			GoType(ComfyNodeExtractStatus("")).
			Default(string(ComfyNodeExtractStatusPending)),
		field.JSON("comfy_node_cloud_build_info", ComfyNodeCloudBuildInfo{}).SchemaType(map[string]string{
			dialect.Postgres: "jsonb",
		}).Optional().Comment("Stores the Google Cloud Build information that extracts the comfy node information."),
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
		edge.To("comfy_nodes", ComfyNode.Type),
	}
}

func (NodeVersion) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("node_id", "version").Unique(),
	}
}

type NodeVersionStatus string

const (
	NodeVersionStatusActive  NodeVersionStatus = "active"
	NodeVersionStatusDeleted NodeVersionStatus = "deleted"
	NodeVersionStatusBanned  NodeVersionStatus = "banned"
	NodeVersionStatusPending NodeVersionStatus = "pending"
	NodeVersionStatusFlagged NodeVersionStatus = "flagged"
)

func (NodeVersionStatus) Values() (types []string) {
	return []string{
		string(NodeVersionStatusActive),
		string(NodeVersionStatusBanned),
		string(NodeVersionStatusDeleted),
		string(NodeVersionStatusPending),
		string(NodeVersionStatusFlagged),
	}
}

type ComfyNodeExtractStatus string

func (ComfyNodeExtractStatus) Values() (types []string) {
	return []string{
		string(ComfyNodeExtractStatusPending),
		string(ComfyNodeExtractStatusFailed),
		string(ComfyNodeExtractStatusSuccess),
	}
}

const (
	ComfyNodeExtractStatusPending = "pending"
	ComfyNodeExtractStatusFailed  = "failed"
	ComfyNodeExtractStatusSuccess = "success"
)

type ComfyNodeCloudBuildInfo struct {
	ProjectID     string `json:"project_id"`
	BuildID       string `json:"build_id"`
	ProjectNumber string `json:"project_number"`
	Location      string `json:"location"`
}

func (c *ComfyNodeCloudBuildInfo) String() string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds;region=%s/%s?project=%s", c.Location, c.BuildID, c.ProjectID)
}

package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// Contains information about a publisher on the Comfy Registry.
type Publisher struct {
	ent.Schema
}

func (Publisher) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Unique().Comment("The unique identifier of the publisher. Cannot be changed."),
		field.String("name").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Comment("The publicly visible name of the publisher."),
		field.String("description").Optional().SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.String("website").Optional().SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.String("support_email").Optional().SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.String("source_code_repo").Optional().SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.String("logo_url").Optional().SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional(),
		field.Enum("status").GoType(PublisherStatusType("")).Default(string(PublisherStatusTypeActive)),
	}
}

func (Publisher) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (Publisher) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("publisher_permissions", PublisherPermission.Type),
		edge.To("nodes", Node.Type),
		edge.To("personal_access_tokens", PersonalAccessToken.Type),
	}
}

type PublisherStatusType string

const (
	PublisherStatusTypeActive PublisherStatusType = "ACTIVE"
	PublisherStatusTypeBanned PublisherStatusType = "BANNED"
)

func (PublisherStatusType) Values() (types []string) {
	return []string{
		string(PublisherStatusTypeActive),
		string(PublisherStatusTypeBanned),
	}
}

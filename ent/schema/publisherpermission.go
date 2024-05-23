package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// This type describes which user has which permission on a publisher.
type PublisherPermission struct {
	ent.Schema
}

func (PublisherPermission) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("permission").GoType(PublisherPermissionType("")),
		field.String("user_id").StorageKey("user_id").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
		field.String("publisher_id").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}),
	}
}

func (PublisherPermission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("publisher_permissions").
			Field("user_id").Unique().Required(),
		edge.From("publisher", Publisher.Type).Ref("publisher_permissions").
			Field("publisher_id").Unique().Required(),
	}
}

type PublisherPermissionType string

const (
	PublisherPermissionTypeOwner  PublisherPermissionType = "owner"
	PublisherPermissionTypeMember PublisherPermissionType = "member"
)

func (PublisherPermissionType) Values() (types []string) {
	return []string{
		string(PublisherPermissionTypeOwner),
		string(PublisherPermissionTypeMember),
	}
}

package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Comment("The firebase UID of the user"),
		field.String("email").Optional(),
		field.String("name").Optional(),
		field.Bool("is_approved").Default(false).Comment("Whether the user is approved to use the platform"),
		field.Bool("is_admin").Default(false).Comment("Whether the user is approved to use the platform"),
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("publisher_permissions", PublisherPermission.Type),
	}
}

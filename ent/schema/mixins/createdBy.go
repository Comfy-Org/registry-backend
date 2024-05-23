package dripmixins

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// Write a ent Mixin that adds a "created_by" field to all entities.
// The field should be of type "string" and should be optional.
// The field should be named "created_by".
// The field should have a default value of "system".

type EditedByMixin struct {
	mixin.Schema
}

func (EditedByMixin) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("created_by", uuid.UUID{}).Optional(),
		field.UUID("updated_by", uuid.UUID{}).Optional(),
	}
}

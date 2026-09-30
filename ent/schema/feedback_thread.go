package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
	"time"
)

// FeedbackThread is private to Registry admins and current Publisher owners.
// Version and Publisher IDs are retained for audit when source resources are deleted.
type FeedbackThread struct{ ent.Schema }

func (FeedbackThread) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("version_id", uuid.UUID{}),
		field.String("publisher_id"),
		field.Enum("state").Values("awaiting_author", "awaiting_admin", "resolved").Default("awaiting_author"),
		field.String("created_by_user_id"),
		field.Int("last_message_seq").Default(0),
		field.Time("last_message_at").Default(time.Now),
		field.Int("revision").Default(0),
		field.String("resolved_by_user_id").Optional(),
		field.Time("resolved_at").Optional().Nillable(),
		field.Time("archived_at").Optional().Nillable(),
	}
}
func (FeedbackThread) Mixin() []ent.Mixin { return []ent.Mixin{mixin.Time{}} }
func (FeedbackThread) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("messages", FeedbackMessage.Type), edge.To("reads", FeedbackRead.Type), edge.To("events", FeedbackEvent.Type),
	}
}
func (FeedbackThread) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("version_id", "publisher_id").Unique(),
		index.Fields("publisher_id", "state", "last_message_at"),
		index.Fields("state", "last_message_at"),
	}
}

package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"time"
)

type FeedbackEvent struct{ ent.Schema }

func (FeedbackEvent) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New), field.UUID("thread_id", uuid.UUID{}), field.String("actor_user_id"),
		field.Enum("event_type").Values("resolved", "reopened", "archived", "superseded"), field.Time("created_at").Default(time.Now).Immutable(),
		field.UUID("replacement_version_id", uuid.UUID{}).Optional().Nillable(),
		field.String("replacement_version").Optional().Nillable(),
	}
}
func (FeedbackEvent) Edges() []ent.Edge {
	return []ent.Edge{edge.From("thread", FeedbackThread.Type).Ref("events").Field("thread_id").Unique().Required()}
}
func (FeedbackEvent) Indexes() []ent.Index {
	return []ent.Index{index.Fields("thread_id", "created_at")}
}

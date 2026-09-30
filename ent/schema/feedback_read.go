package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"time"
)

type FeedbackRead struct{ ent.Schema }

func (FeedbackRead) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New), field.UUID("thread_id", uuid.UUID{}), field.String("user_id"),
		field.Int("last_read_message_seq").Default(0).NonNegative(), field.Time("read_at").Default(time.Now),
	}
}
func (FeedbackRead) Edges() []ent.Edge {
	return []ent.Edge{edge.From("thread", FeedbackThread.Type).Ref("reads").Field("thread_id").Unique().Required()}
}
func (FeedbackRead) Indexes() []ent.Index {
	return []ent.Index{index.Fields("thread_id", "user_id").Unique(), index.Fields("user_id", "thread_id")}
}

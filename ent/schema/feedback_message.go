package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"time"
)

type FeedbackMessage struct{ ent.Schema }

func (FeedbackMessage) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("thread_id", uuid.UUID{}),
		field.Int("seq").Positive(),
		field.String("sender_user_id"),
		field.String("sender_name"),
		field.Enum("sender_role").Values("admin", "author"),
		field.String("body").NotEmpty(),
		field.UUID("client_message_id", uuid.UUID{}),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}
func (FeedbackMessage) Edges() []ent.Edge {
	return []ent.Edge{edge.From("thread", FeedbackThread.Type).Ref("messages").Field("thread_id").Unique().Required()}
}
func (FeedbackMessage) Indexes() []ent.Index {
	return []ent.Index{index.Fields("thread_id", "seq").Unique(), index.Fields("thread_id", "sender_user_id", "client_message_id").Unique()}
}

package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// Holds a generic GCP Storage Object
type StorageFile struct {
	ent.Schema
}

// Fields of the StorageFile.
func (StorageFile) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.String("bucket_name").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).NotEmpty(),
		field.String("object_name").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).NotEmpty().Optional(),
		field.String("file_path").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).NotEmpty(),
		field.String("file_type").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).NotEmpty().Comment("e.g., image, video"),
		field.String("file_url").SchemaType(map[string]string{
			dialect.Postgres: "text",
		}).Optional().Comment("Publicly accessible URL of the file, if available"),
	}
}

// Edges of the StorageFile.
func (StorageFile) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (StorageFile) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

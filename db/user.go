package db

import (
	"context"
	"registry-backend/ent"
	"registry-backend/ent/user"
)

// UpsertUser creates or updates a user in the database
func UpsertUser(ctx context.Context, client *ent.Client, firebaseUID string, email string, name string) error {
	return client.User.Create().
		SetID(firebaseUID).
		SetEmail(email).
		SetName(name).
		OnConflictColumns(user.FieldID).
		UpdateEmail().
		UpdateName().
		Exec(ctx)
}

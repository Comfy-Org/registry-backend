// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"registry-backend/ent/predicate"
	"registry-backend/ent/storagefile"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// StorageFileDelete is the builder for deleting a StorageFile entity.
type StorageFileDelete struct {
	config
	hooks    []Hook
	mutation *StorageFileMutation
}

// Where appends a list predicates to the StorageFileDelete builder.
func (sfd *StorageFileDelete) Where(ps ...predicate.StorageFile) *StorageFileDelete {
	sfd.mutation.Where(ps...)
	return sfd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (sfd *StorageFileDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, sfd.sqlExec, sfd.mutation, sfd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (sfd *StorageFileDelete) ExecX(ctx context.Context) int {
	n, err := sfd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (sfd *StorageFileDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(storagefile.Table, sqlgraph.NewFieldSpec(storagefile.FieldID, field.TypeUUID))
	if ps := sfd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, sfd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	sfd.mutation.done = true
	return affected, err
}

// StorageFileDeleteOne is the builder for deleting a single StorageFile entity.
type StorageFileDeleteOne struct {
	sfd *StorageFileDelete
}

// Where appends a list predicates to the StorageFileDelete builder.
func (sfdo *StorageFileDeleteOne) Where(ps ...predicate.StorageFile) *StorageFileDeleteOne {
	sfdo.sfd.mutation.Where(ps...)
	return sfdo
}

// Exec executes the deletion query.
func (sfdo *StorageFileDeleteOne) Exec(ctx context.Context) error {
	n, err := sfdo.sfd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{storagefile.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (sfdo *StorageFileDeleteOne) ExecX(ctx context.Context) {
	if err := sfdo.Exec(ctx); err != nil {
		panic(err)
	}
}

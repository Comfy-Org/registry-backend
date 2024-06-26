// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"registry-backend/ent/predicate"
	"registry-backend/ent/publisher"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// PublisherDelete is the builder for deleting a Publisher entity.
type PublisherDelete struct {
	config
	hooks    []Hook
	mutation *PublisherMutation
}

// Where appends a list predicates to the PublisherDelete builder.
func (pd *PublisherDelete) Where(ps ...predicate.Publisher) *PublisherDelete {
	pd.mutation.Where(ps...)
	return pd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (pd *PublisherDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, pd.sqlExec, pd.mutation, pd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (pd *PublisherDelete) ExecX(ctx context.Context) int {
	n, err := pd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (pd *PublisherDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(publisher.Table, sqlgraph.NewFieldSpec(publisher.FieldID, field.TypeString))
	if ps := pd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, pd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	pd.mutation.done = true
	return affected, err
}

// PublisherDeleteOne is the builder for deleting a single Publisher entity.
type PublisherDeleteOne struct {
	pd *PublisherDelete
}

// Where appends a list predicates to the PublisherDelete builder.
func (pdo *PublisherDeleteOne) Where(ps ...predicate.Publisher) *PublisherDeleteOne {
	pdo.pd.mutation.Where(ps...)
	return pdo
}

// Exec executes the deletion query.
func (pdo *PublisherDeleteOne) Exec(ctx context.Context) error {
	n, err := pdo.pd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{publisher.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (pdo *PublisherDeleteOne) ExecX(ctx context.Context) {
	if err := pdo.Exec(ctx); err != nil {
		panic(err)
	}
}

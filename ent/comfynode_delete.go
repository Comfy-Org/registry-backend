// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"registry-backend/ent/comfynode"
	"registry-backend/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// ComfyNodeDelete is the builder for deleting a ComfyNode entity.
type ComfyNodeDelete struct {
	config
	hooks    []Hook
	mutation *ComfyNodeMutation
}

// Where appends a list predicates to the ComfyNodeDelete builder.
func (cnd *ComfyNodeDelete) Where(ps ...predicate.ComfyNode) *ComfyNodeDelete {
	cnd.mutation.Where(ps...)
	return cnd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (cnd *ComfyNodeDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, cnd.sqlExec, cnd.mutation, cnd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (cnd *ComfyNodeDelete) ExecX(ctx context.Context) int {
	n, err := cnd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (cnd *ComfyNodeDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(comfynode.Table, sqlgraph.NewFieldSpec(comfynode.FieldID, field.TypeUUID))
	if ps := cnd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, cnd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	cnd.mutation.done = true
	return affected, err
}

// ComfyNodeDeleteOne is the builder for deleting a single ComfyNode entity.
type ComfyNodeDeleteOne struct {
	cnd *ComfyNodeDelete
}

// Where appends a list predicates to the ComfyNodeDelete builder.
func (cndo *ComfyNodeDeleteOne) Where(ps ...predicate.ComfyNode) *ComfyNodeDeleteOne {
	cndo.cnd.mutation.Where(ps...)
	return cndo
}

// Exec executes the deletion query.
func (cndo *ComfyNodeDeleteOne) Exec(ctx context.Context) error {
	n, err := cndo.cnd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{comfynode.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (cndo *ComfyNodeDeleteOne) ExecX(ctx context.Context) {
	if err := cndo.Exec(ctx); err != nil {
		panic(err)
	}
}

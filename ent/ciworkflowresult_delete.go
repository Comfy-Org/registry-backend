// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"registry-backend/ent/ciworkflowresult"
	"registry-backend/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// CIWorkflowResultDelete is the builder for deleting a CIWorkflowResult entity.
type CIWorkflowResultDelete struct {
	config
	hooks    []Hook
	mutation *CIWorkflowResultMutation
}

// Where appends a list predicates to the CIWorkflowResultDelete builder.
func (cwrd *CIWorkflowResultDelete) Where(ps ...predicate.CIWorkflowResult) *CIWorkflowResultDelete {
	cwrd.mutation.Where(ps...)
	return cwrd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (cwrd *CIWorkflowResultDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, cwrd.sqlExec, cwrd.mutation, cwrd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (cwrd *CIWorkflowResultDelete) ExecX(ctx context.Context) int {
	n, err := cwrd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (cwrd *CIWorkflowResultDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(ciworkflowresult.Table, sqlgraph.NewFieldSpec(ciworkflowresult.FieldID, field.TypeUUID))
	if ps := cwrd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, cwrd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	cwrd.mutation.done = true
	return affected, err
}

// CIWorkflowResultDeleteOne is the builder for deleting a single CIWorkflowResult entity.
type CIWorkflowResultDeleteOne struct {
	cwrd *CIWorkflowResultDelete
}

// Where appends a list predicates to the CIWorkflowResultDelete builder.
func (cwrdo *CIWorkflowResultDeleteOne) Where(ps ...predicate.CIWorkflowResult) *CIWorkflowResultDeleteOne {
	cwrdo.cwrd.mutation.Where(ps...)
	return cwrdo
}

// Exec executes the deletion query.
func (cwrdo *CIWorkflowResultDeleteOne) Exec(ctx context.Context) error {
	n, err := cwrdo.cwrd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{ciworkflowresult.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (cwrdo *CIWorkflowResultDeleteOne) ExecX(ctx context.Context) {
	if err := cwrdo.Exec(ctx); err != nil {
		panic(err)
	}
}
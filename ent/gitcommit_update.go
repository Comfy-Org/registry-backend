// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"registry-backend/ent/ciworkflowresult"
	"registry-backend/ent/gitcommit"
	"registry-backend/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// GitCommitUpdate is the builder for updating GitCommit entities.
type GitCommitUpdate struct {
	config
	hooks    []Hook
	mutation *GitCommitMutation
}

// Where appends a list predicates to the GitCommitUpdate builder.
func (gcu *GitCommitUpdate) Where(ps ...predicate.GitCommit) *GitCommitUpdate {
	gcu.mutation.Where(ps...)
	return gcu
}

// SetUpdateTime sets the "update_time" field.
func (gcu *GitCommitUpdate) SetUpdateTime(t time.Time) *GitCommitUpdate {
	gcu.mutation.SetUpdateTime(t)
	return gcu
}

// SetCommitHash sets the "commit_hash" field.
func (gcu *GitCommitUpdate) SetCommitHash(s string) *GitCommitUpdate {
	gcu.mutation.SetCommitHash(s)
	return gcu
}

// SetNillableCommitHash sets the "commit_hash" field if the given value is not nil.
func (gcu *GitCommitUpdate) SetNillableCommitHash(s *string) *GitCommitUpdate {
	if s != nil {
		gcu.SetCommitHash(*s)
	}
	return gcu
}

// SetBranchName sets the "branch_name" field.
func (gcu *GitCommitUpdate) SetBranchName(s string) *GitCommitUpdate {
	gcu.mutation.SetBranchName(s)
	return gcu
}

// SetNillableBranchName sets the "branch_name" field if the given value is not nil.
func (gcu *GitCommitUpdate) SetNillableBranchName(s *string) *GitCommitUpdate {
	if s != nil {
		gcu.SetBranchName(*s)
	}
	return gcu
}

// SetRepoName sets the "repo_name" field.
func (gcu *GitCommitUpdate) SetRepoName(s string) *GitCommitUpdate {
	gcu.mutation.SetRepoName(s)
	return gcu
}

// SetNillableRepoName sets the "repo_name" field if the given value is not nil.
func (gcu *GitCommitUpdate) SetNillableRepoName(s *string) *GitCommitUpdate {
	if s != nil {
		gcu.SetRepoName(*s)
	}
	return gcu
}

// SetCommitMessage sets the "commit_message" field.
func (gcu *GitCommitUpdate) SetCommitMessage(s string) *GitCommitUpdate {
	gcu.mutation.SetCommitMessage(s)
	return gcu
}

// SetNillableCommitMessage sets the "commit_message" field if the given value is not nil.
func (gcu *GitCommitUpdate) SetNillableCommitMessage(s *string) *GitCommitUpdate {
	if s != nil {
		gcu.SetCommitMessage(*s)
	}
	return gcu
}

// SetCommitTimestamp sets the "commit_timestamp" field.
func (gcu *GitCommitUpdate) SetCommitTimestamp(t time.Time) *GitCommitUpdate {
	gcu.mutation.SetCommitTimestamp(t)
	return gcu
}

// SetNillableCommitTimestamp sets the "commit_timestamp" field if the given value is not nil.
func (gcu *GitCommitUpdate) SetNillableCommitTimestamp(t *time.Time) *GitCommitUpdate {
	if t != nil {
		gcu.SetCommitTimestamp(*t)
	}
	return gcu
}

// SetAuthor sets the "author" field.
func (gcu *GitCommitUpdate) SetAuthor(s string) *GitCommitUpdate {
	gcu.mutation.SetAuthor(s)
	return gcu
}

// SetNillableAuthor sets the "author" field if the given value is not nil.
func (gcu *GitCommitUpdate) SetNillableAuthor(s *string) *GitCommitUpdate {
	if s != nil {
		gcu.SetAuthor(*s)
	}
	return gcu
}

// ClearAuthor clears the value of the "author" field.
func (gcu *GitCommitUpdate) ClearAuthor() *GitCommitUpdate {
	gcu.mutation.ClearAuthor()
	return gcu
}

// SetTimestamp sets the "timestamp" field.
func (gcu *GitCommitUpdate) SetTimestamp(t time.Time) *GitCommitUpdate {
	gcu.mutation.SetTimestamp(t)
	return gcu
}

// SetNillableTimestamp sets the "timestamp" field if the given value is not nil.
func (gcu *GitCommitUpdate) SetNillableTimestamp(t *time.Time) *GitCommitUpdate {
	if t != nil {
		gcu.SetTimestamp(*t)
	}
	return gcu
}

// ClearTimestamp clears the value of the "timestamp" field.
func (gcu *GitCommitUpdate) ClearTimestamp() *GitCommitUpdate {
	gcu.mutation.ClearTimestamp()
	return gcu
}

// AddResultIDs adds the "results" edge to the CIWorkflowResult entity by IDs.
func (gcu *GitCommitUpdate) AddResultIDs(ids ...uuid.UUID) *GitCommitUpdate {
	gcu.mutation.AddResultIDs(ids...)
	return gcu
}

// AddResults adds the "results" edges to the CIWorkflowResult entity.
func (gcu *GitCommitUpdate) AddResults(c ...*CIWorkflowResult) *GitCommitUpdate {
	ids := make([]uuid.UUID, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return gcu.AddResultIDs(ids...)
}

// Mutation returns the GitCommitMutation object of the builder.
func (gcu *GitCommitUpdate) Mutation() *GitCommitMutation {
	return gcu.mutation
}

// ClearResults clears all "results" edges to the CIWorkflowResult entity.
func (gcu *GitCommitUpdate) ClearResults() *GitCommitUpdate {
	gcu.mutation.ClearResults()
	return gcu
}

// RemoveResultIDs removes the "results" edge to CIWorkflowResult entities by IDs.
func (gcu *GitCommitUpdate) RemoveResultIDs(ids ...uuid.UUID) *GitCommitUpdate {
	gcu.mutation.RemoveResultIDs(ids...)
	return gcu
}

// RemoveResults removes "results" edges to CIWorkflowResult entities.
func (gcu *GitCommitUpdate) RemoveResults(c ...*CIWorkflowResult) *GitCommitUpdate {
	ids := make([]uuid.UUID, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return gcu.RemoveResultIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (gcu *GitCommitUpdate) Save(ctx context.Context) (int, error) {
	gcu.defaults()
	return withHooks(ctx, gcu.sqlSave, gcu.mutation, gcu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (gcu *GitCommitUpdate) SaveX(ctx context.Context) int {
	affected, err := gcu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (gcu *GitCommitUpdate) Exec(ctx context.Context) error {
	_, err := gcu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gcu *GitCommitUpdate) ExecX(ctx context.Context) {
	if err := gcu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (gcu *GitCommitUpdate) defaults() {
	if _, ok := gcu.mutation.UpdateTime(); !ok {
		v := gitcommit.UpdateDefaultUpdateTime()
		gcu.mutation.SetUpdateTime(v)
	}
}

func (gcu *GitCommitUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(gitcommit.Table, gitcommit.Columns, sqlgraph.NewFieldSpec(gitcommit.FieldID, field.TypeUUID))
	if ps := gcu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := gcu.mutation.UpdateTime(); ok {
		_spec.SetField(gitcommit.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := gcu.mutation.CommitHash(); ok {
		_spec.SetField(gitcommit.FieldCommitHash, field.TypeString, value)
	}
	if value, ok := gcu.mutation.BranchName(); ok {
		_spec.SetField(gitcommit.FieldBranchName, field.TypeString, value)
	}
	if value, ok := gcu.mutation.RepoName(); ok {
		_spec.SetField(gitcommit.FieldRepoName, field.TypeString, value)
	}
	if value, ok := gcu.mutation.CommitMessage(); ok {
		_spec.SetField(gitcommit.FieldCommitMessage, field.TypeString, value)
	}
	if value, ok := gcu.mutation.CommitTimestamp(); ok {
		_spec.SetField(gitcommit.FieldCommitTimestamp, field.TypeTime, value)
	}
	if value, ok := gcu.mutation.Author(); ok {
		_spec.SetField(gitcommit.FieldAuthor, field.TypeString, value)
	}
	if gcu.mutation.AuthorCleared() {
		_spec.ClearField(gitcommit.FieldAuthor, field.TypeString)
	}
	if value, ok := gcu.mutation.Timestamp(); ok {
		_spec.SetField(gitcommit.FieldTimestamp, field.TypeTime, value)
	}
	if gcu.mutation.TimestampCleared() {
		_spec.ClearField(gitcommit.FieldTimestamp, field.TypeTime)
	}
	if gcu.mutation.ResultsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   gitcommit.ResultsTable,
			Columns: []string{gitcommit.ResultsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(ciworkflowresult.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := gcu.mutation.RemovedResultsIDs(); len(nodes) > 0 && !gcu.mutation.ResultsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   gitcommit.ResultsTable,
			Columns: []string{gitcommit.ResultsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(ciworkflowresult.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := gcu.mutation.ResultsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   gitcommit.ResultsTable,
			Columns: []string{gitcommit.ResultsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(ciworkflowresult.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, gcu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{gitcommit.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	gcu.mutation.done = true
	return n, nil
}

// GitCommitUpdateOne is the builder for updating a single GitCommit entity.
type GitCommitUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *GitCommitMutation
}

// SetUpdateTime sets the "update_time" field.
func (gcuo *GitCommitUpdateOne) SetUpdateTime(t time.Time) *GitCommitUpdateOne {
	gcuo.mutation.SetUpdateTime(t)
	return gcuo
}

// SetCommitHash sets the "commit_hash" field.
func (gcuo *GitCommitUpdateOne) SetCommitHash(s string) *GitCommitUpdateOne {
	gcuo.mutation.SetCommitHash(s)
	return gcuo
}

// SetNillableCommitHash sets the "commit_hash" field if the given value is not nil.
func (gcuo *GitCommitUpdateOne) SetNillableCommitHash(s *string) *GitCommitUpdateOne {
	if s != nil {
		gcuo.SetCommitHash(*s)
	}
	return gcuo
}

// SetBranchName sets the "branch_name" field.
func (gcuo *GitCommitUpdateOne) SetBranchName(s string) *GitCommitUpdateOne {
	gcuo.mutation.SetBranchName(s)
	return gcuo
}

// SetNillableBranchName sets the "branch_name" field if the given value is not nil.
func (gcuo *GitCommitUpdateOne) SetNillableBranchName(s *string) *GitCommitUpdateOne {
	if s != nil {
		gcuo.SetBranchName(*s)
	}
	return gcuo
}

// SetRepoName sets the "repo_name" field.
func (gcuo *GitCommitUpdateOne) SetRepoName(s string) *GitCommitUpdateOne {
	gcuo.mutation.SetRepoName(s)
	return gcuo
}

// SetNillableRepoName sets the "repo_name" field if the given value is not nil.
func (gcuo *GitCommitUpdateOne) SetNillableRepoName(s *string) *GitCommitUpdateOne {
	if s != nil {
		gcuo.SetRepoName(*s)
	}
	return gcuo
}

// SetCommitMessage sets the "commit_message" field.
func (gcuo *GitCommitUpdateOne) SetCommitMessage(s string) *GitCommitUpdateOne {
	gcuo.mutation.SetCommitMessage(s)
	return gcuo
}

// SetNillableCommitMessage sets the "commit_message" field if the given value is not nil.
func (gcuo *GitCommitUpdateOne) SetNillableCommitMessage(s *string) *GitCommitUpdateOne {
	if s != nil {
		gcuo.SetCommitMessage(*s)
	}
	return gcuo
}

// SetCommitTimestamp sets the "commit_timestamp" field.
func (gcuo *GitCommitUpdateOne) SetCommitTimestamp(t time.Time) *GitCommitUpdateOne {
	gcuo.mutation.SetCommitTimestamp(t)
	return gcuo
}

// SetNillableCommitTimestamp sets the "commit_timestamp" field if the given value is not nil.
func (gcuo *GitCommitUpdateOne) SetNillableCommitTimestamp(t *time.Time) *GitCommitUpdateOne {
	if t != nil {
		gcuo.SetCommitTimestamp(*t)
	}
	return gcuo
}

// SetAuthor sets the "author" field.
func (gcuo *GitCommitUpdateOne) SetAuthor(s string) *GitCommitUpdateOne {
	gcuo.mutation.SetAuthor(s)
	return gcuo
}

// SetNillableAuthor sets the "author" field if the given value is not nil.
func (gcuo *GitCommitUpdateOne) SetNillableAuthor(s *string) *GitCommitUpdateOne {
	if s != nil {
		gcuo.SetAuthor(*s)
	}
	return gcuo
}

// ClearAuthor clears the value of the "author" field.
func (gcuo *GitCommitUpdateOne) ClearAuthor() *GitCommitUpdateOne {
	gcuo.mutation.ClearAuthor()
	return gcuo
}

// SetTimestamp sets the "timestamp" field.
func (gcuo *GitCommitUpdateOne) SetTimestamp(t time.Time) *GitCommitUpdateOne {
	gcuo.mutation.SetTimestamp(t)
	return gcuo
}

// SetNillableTimestamp sets the "timestamp" field if the given value is not nil.
func (gcuo *GitCommitUpdateOne) SetNillableTimestamp(t *time.Time) *GitCommitUpdateOne {
	if t != nil {
		gcuo.SetTimestamp(*t)
	}
	return gcuo
}

// ClearTimestamp clears the value of the "timestamp" field.
func (gcuo *GitCommitUpdateOne) ClearTimestamp() *GitCommitUpdateOne {
	gcuo.mutation.ClearTimestamp()
	return gcuo
}

// AddResultIDs adds the "results" edge to the CIWorkflowResult entity by IDs.
func (gcuo *GitCommitUpdateOne) AddResultIDs(ids ...uuid.UUID) *GitCommitUpdateOne {
	gcuo.mutation.AddResultIDs(ids...)
	return gcuo
}

// AddResults adds the "results" edges to the CIWorkflowResult entity.
func (gcuo *GitCommitUpdateOne) AddResults(c ...*CIWorkflowResult) *GitCommitUpdateOne {
	ids := make([]uuid.UUID, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return gcuo.AddResultIDs(ids...)
}

// Mutation returns the GitCommitMutation object of the builder.
func (gcuo *GitCommitUpdateOne) Mutation() *GitCommitMutation {
	return gcuo.mutation
}

// ClearResults clears all "results" edges to the CIWorkflowResult entity.
func (gcuo *GitCommitUpdateOne) ClearResults() *GitCommitUpdateOne {
	gcuo.mutation.ClearResults()
	return gcuo
}

// RemoveResultIDs removes the "results" edge to CIWorkflowResult entities by IDs.
func (gcuo *GitCommitUpdateOne) RemoveResultIDs(ids ...uuid.UUID) *GitCommitUpdateOne {
	gcuo.mutation.RemoveResultIDs(ids...)
	return gcuo
}

// RemoveResults removes "results" edges to CIWorkflowResult entities.
func (gcuo *GitCommitUpdateOne) RemoveResults(c ...*CIWorkflowResult) *GitCommitUpdateOne {
	ids := make([]uuid.UUID, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return gcuo.RemoveResultIDs(ids...)
}

// Where appends a list predicates to the GitCommitUpdate builder.
func (gcuo *GitCommitUpdateOne) Where(ps ...predicate.GitCommit) *GitCommitUpdateOne {
	gcuo.mutation.Where(ps...)
	return gcuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (gcuo *GitCommitUpdateOne) Select(field string, fields ...string) *GitCommitUpdateOne {
	gcuo.fields = append([]string{field}, fields...)
	return gcuo
}

// Save executes the query and returns the updated GitCommit entity.
func (gcuo *GitCommitUpdateOne) Save(ctx context.Context) (*GitCommit, error) {
	gcuo.defaults()
	return withHooks(ctx, gcuo.sqlSave, gcuo.mutation, gcuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (gcuo *GitCommitUpdateOne) SaveX(ctx context.Context) *GitCommit {
	node, err := gcuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (gcuo *GitCommitUpdateOne) Exec(ctx context.Context) error {
	_, err := gcuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gcuo *GitCommitUpdateOne) ExecX(ctx context.Context) {
	if err := gcuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (gcuo *GitCommitUpdateOne) defaults() {
	if _, ok := gcuo.mutation.UpdateTime(); !ok {
		v := gitcommit.UpdateDefaultUpdateTime()
		gcuo.mutation.SetUpdateTime(v)
	}
}

func (gcuo *GitCommitUpdateOne) sqlSave(ctx context.Context) (_node *GitCommit, err error) {
	_spec := sqlgraph.NewUpdateSpec(gitcommit.Table, gitcommit.Columns, sqlgraph.NewFieldSpec(gitcommit.FieldID, field.TypeUUID))
	id, ok := gcuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "GitCommit.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := gcuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, gitcommit.FieldID)
		for _, f := range fields {
			if !gitcommit.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != gitcommit.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := gcuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := gcuo.mutation.UpdateTime(); ok {
		_spec.SetField(gitcommit.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := gcuo.mutation.CommitHash(); ok {
		_spec.SetField(gitcommit.FieldCommitHash, field.TypeString, value)
	}
	if value, ok := gcuo.mutation.BranchName(); ok {
		_spec.SetField(gitcommit.FieldBranchName, field.TypeString, value)
	}
	if value, ok := gcuo.mutation.RepoName(); ok {
		_spec.SetField(gitcommit.FieldRepoName, field.TypeString, value)
	}
	if value, ok := gcuo.mutation.CommitMessage(); ok {
		_spec.SetField(gitcommit.FieldCommitMessage, field.TypeString, value)
	}
	if value, ok := gcuo.mutation.CommitTimestamp(); ok {
		_spec.SetField(gitcommit.FieldCommitTimestamp, field.TypeTime, value)
	}
	if value, ok := gcuo.mutation.Author(); ok {
		_spec.SetField(gitcommit.FieldAuthor, field.TypeString, value)
	}
	if gcuo.mutation.AuthorCleared() {
		_spec.ClearField(gitcommit.FieldAuthor, field.TypeString)
	}
	if value, ok := gcuo.mutation.Timestamp(); ok {
		_spec.SetField(gitcommit.FieldTimestamp, field.TypeTime, value)
	}
	if gcuo.mutation.TimestampCleared() {
		_spec.ClearField(gitcommit.FieldTimestamp, field.TypeTime)
	}
	if gcuo.mutation.ResultsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   gitcommit.ResultsTable,
			Columns: []string{gitcommit.ResultsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(ciworkflowresult.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := gcuo.mutation.RemovedResultsIDs(); len(nodes) > 0 && !gcuo.mutation.ResultsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   gitcommit.ResultsTable,
			Columns: []string{gitcommit.ResultsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(ciworkflowresult.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := gcuo.mutation.ResultsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   gitcommit.ResultsTable,
			Columns: []string{gitcommit.ResultsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(ciworkflowresult.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &GitCommit{config: gcuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, gcuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{gitcommit.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	gcuo.mutation.done = true
	return _node, nil
}

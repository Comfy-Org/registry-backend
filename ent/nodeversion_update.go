// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"registry-backend/ent/node"
	"registry-backend/ent/nodeversion"
	"registry-backend/ent/predicate"
	"registry-backend/ent/storagefile"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/dialect/sql/sqljson"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// NodeVersionUpdate is the builder for updating NodeVersion entities.
type NodeVersionUpdate struct {
	config
	hooks    []Hook
	mutation *NodeVersionMutation
}

// Where appends a list predicates to the NodeVersionUpdate builder.
func (nvu *NodeVersionUpdate) Where(ps ...predicate.NodeVersion) *NodeVersionUpdate {
	nvu.mutation.Where(ps...)
	return nvu
}

// SetUpdateTime sets the "update_time" field.
func (nvu *NodeVersionUpdate) SetUpdateTime(t time.Time) *NodeVersionUpdate {
	nvu.mutation.SetUpdateTime(t)
	return nvu
}

// SetNodeID sets the "node_id" field.
func (nvu *NodeVersionUpdate) SetNodeID(s string) *NodeVersionUpdate {
	nvu.mutation.SetNodeID(s)
	return nvu
}

// SetNillableNodeID sets the "node_id" field if the given value is not nil.
func (nvu *NodeVersionUpdate) SetNillableNodeID(s *string) *NodeVersionUpdate {
	if s != nil {
		nvu.SetNodeID(*s)
	}
	return nvu
}

// SetVersion sets the "version" field.
func (nvu *NodeVersionUpdate) SetVersion(s string) *NodeVersionUpdate {
	nvu.mutation.SetVersion(s)
	return nvu
}

// SetNillableVersion sets the "version" field if the given value is not nil.
func (nvu *NodeVersionUpdate) SetNillableVersion(s *string) *NodeVersionUpdate {
	if s != nil {
		nvu.SetVersion(*s)
	}
	return nvu
}

// SetChangelog sets the "changelog" field.
func (nvu *NodeVersionUpdate) SetChangelog(s string) *NodeVersionUpdate {
	nvu.mutation.SetChangelog(s)
	return nvu
}

// SetNillableChangelog sets the "changelog" field if the given value is not nil.
func (nvu *NodeVersionUpdate) SetNillableChangelog(s *string) *NodeVersionUpdate {
	if s != nil {
		nvu.SetChangelog(*s)
	}
	return nvu
}

// ClearChangelog clears the value of the "changelog" field.
func (nvu *NodeVersionUpdate) ClearChangelog() *NodeVersionUpdate {
	nvu.mutation.ClearChangelog()
	return nvu
}

// SetPipDependencies sets the "pip_dependencies" field.
func (nvu *NodeVersionUpdate) SetPipDependencies(s []string) *NodeVersionUpdate {
	nvu.mutation.SetPipDependencies(s)
	return nvu
}

// AppendPipDependencies appends s to the "pip_dependencies" field.
func (nvu *NodeVersionUpdate) AppendPipDependencies(s []string) *NodeVersionUpdate {
	nvu.mutation.AppendPipDependencies(s)
	return nvu
}

// SetDeprecated sets the "deprecated" field.
func (nvu *NodeVersionUpdate) SetDeprecated(b bool) *NodeVersionUpdate {
	nvu.mutation.SetDeprecated(b)
	return nvu
}

// SetNillableDeprecated sets the "deprecated" field if the given value is not nil.
func (nvu *NodeVersionUpdate) SetNillableDeprecated(b *bool) *NodeVersionUpdate {
	if b != nil {
		nvu.SetDeprecated(*b)
	}
	return nvu
}

// SetNode sets the "node" edge to the Node entity.
func (nvu *NodeVersionUpdate) SetNode(n *Node) *NodeVersionUpdate {
	return nvu.SetNodeID(n.ID)
}

// SetStorageFileID sets the "storage_file" edge to the StorageFile entity by ID.
func (nvu *NodeVersionUpdate) SetStorageFileID(id uuid.UUID) *NodeVersionUpdate {
	nvu.mutation.SetStorageFileID(id)
	return nvu
}

// SetNillableStorageFileID sets the "storage_file" edge to the StorageFile entity by ID if the given value is not nil.
func (nvu *NodeVersionUpdate) SetNillableStorageFileID(id *uuid.UUID) *NodeVersionUpdate {
	if id != nil {
		nvu = nvu.SetStorageFileID(*id)
	}
	return nvu
}

// SetStorageFile sets the "storage_file" edge to the StorageFile entity.
func (nvu *NodeVersionUpdate) SetStorageFile(s *StorageFile) *NodeVersionUpdate {
	return nvu.SetStorageFileID(s.ID)
}

// Mutation returns the NodeVersionMutation object of the builder.
func (nvu *NodeVersionUpdate) Mutation() *NodeVersionMutation {
	return nvu.mutation
}

// ClearNode clears the "node" edge to the Node entity.
func (nvu *NodeVersionUpdate) ClearNode() *NodeVersionUpdate {
	nvu.mutation.ClearNode()
	return nvu
}

// ClearStorageFile clears the "storage_file" edge to the StorageFile entity.
func (nvu *NodeVersionUpdate) ClearStorageFile() *NodeVersionUpdate {
	nvu.mutation.ClearStorageFile()
	return nvu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (nvu *NodeVersionUpdate) Save(ctx context.Context) (int, error) {
	nvu.defaults()
	return withHooks(ctx, nvu.sqlSave, nvu.mutation, nvu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (nvu *NodeVersionUpdate) SaveX(ctx context.Context) int {
	affected, err := nvu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (nvu *NodeVersionUpdate) Exec(ctx context.Context) error {
	_, err := nvu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (nvu *NodeVersionUpdate) ExecX(ctx context.Context) {
	if err := nvu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (nvu *NodeVersionUpdate) defaults() {
	if _, ok := nvu.mutation.UpdateTime(); !ok {
		v := nodeversion.UpdateDefaultUpdateTime()
		nvu.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (nvu *NodeVersionUpdate) check() error {
	if _, ok := nvu.mutation.NodeID(); nvu.mutation.NodeCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "NodeVersion.node"`)
	}
	return nil
}

func (nvu *NodeVersionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := nvu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(nodeversion.Table, nodeversion.Columns, sqlgraph.NewFieldSpec(nodeversion.FieldID, field.TypeUUID))
	if ps := nvu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := nvu.mutation.UpdateTime(); ok {
		_spec.SetField(nodeversion.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := nvu.mutation.Version(); ok {
		_spec.SetField(nodeversion.FieldVersion, field.TypeString, value)
	}
	if value, ok := nvu.mutation.Changelog(); ok {
		_spec.SetField(nodeversion.FieldChangelog, field.TypeString, value)
	}
	if nvu.mutation.ChangelogCleared() {
		_spec.ClearField(nodeversion.FieldChangelog, field.TypeString)
	}
	if value, ok := nvu.mutation.PipDependencies(); ok {
		_spec.SetField(nodeversion.FieldPipDependencies, field.TypeJSON, value)
	}
	if value, ok := nvu.mutation.AppendedPipDependencies(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, nodeversion.FieldPipDependencies, value)
		})
	}
	if value, ok := nvu.mutation.Deprecated(); ok {
		_spec.SetField(nodeversion.FieldDeprecated, field.TypeBool, value)
	}
	if nvu.mutation.NodeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   nodeversion.NodeTable,
			Columns: []string{nodeversion.NodeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(node.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := nvu.mutation.NodeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   nodeversion.NodeTable,
			Columns: []string{nodeversion.NodeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(node.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nvu.mutation.StorageFileCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   nodeversion.StorageFileTable,
			Columns: []string{nodeversion.StorageFileColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(storagefile.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := nvu.mutation.StorageFileIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   nodeversion.StorageFileTable,
			Columns: []string{nodeversion.StorageFileColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(storagefile.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, nvu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{nodeversion.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	nvu.mutation.done = true
	return n, nil
}

// NodeVersionUpdateOne is the builder for updating a single NodeVersion entity.
type NodeVersionUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *NodeVersionMutation
}

// SetUpdateTime sets the "update_time" field.
func (nvuo *NodeVersionUpdateOne) SetUpdateTime(t time.Time) *NodeVersionUpdateOne {
	nvuo.mutation.SetUpdateTime(t)
	return nvuo
}

// SetNodeID sets the "node_id" field.
func (nvuo *NodeVersionUpdateOne) SetNodeID(s string) *NodeVersionUpdateOne {
	nvuo.mutation.SetNodeID(s)
	return nvuo
}

// SetNillableNodeID sets the "node_id" field if the given value is not nil.
func (nvuo *NodeVersionUpdateOne) SetNillableNodeID(s *string) *NodeVersionUpdateOne {
	if s != nil {
		nvuo.SetNodeID(*s)
	}
	return nvuo
}

// SetVersion sets the "version" field.
func (nvuo *NodeVersionUpdateOne) SetVersion(s string) *NodeVersionUpdateOne {
	nvuo.mutation.SetVersion(s)
	return nvuo
}

// SetNillableVersion sets the "version" field if the given value is not nil.
func (nvuo *NodeVersionUpdateOne) SetNillableVersion(s *string) *NodeVersionUpdateOne {
	if s != nil {
		nvuo.SetVersion(*s)
	}
	return nvuo
}

// SetChangelog sets the "changelog" field.
func (nvuo *NodeVersionUpdateOne) SetChangelog(s string) *NodeVersionUpdateOne {
	nvuo.mutation.SetChangelog(s)
	return nvuo
}

// SetNillableChangelog sets the "changelog" field if the given value is not nil.
func (nvuo *NodeVersionUpdateOne) SetNillableChangelog(s *string) *NodeVersionUpdateOne {
	if s != nil {
		nvuo.SetChangelog(*s)
	}
	return nvuo
}

// ClearChangelog clears the value of the "changelog" field.
func (nvuo *NodeVersionUpdateOne) ClearChangelog() *NodeVersionUpdateOne {
	nvuo.mutation.ClearChangelog()
	return nvuo
}

// SetPipDependencies sets the "pip_dependencies" field.
func (nvuo *NodeVersionUpdateOne) SetPipDependencies(s []string) *NodeVersionUpdateOne {
	nvuo.mutation.SetPipDependencies(s)
	return nvuo
}

// AppendPipDependencies appends s to the "pip_dependencies" field.
func (nvuo *NodeVersionUpdateOne) AppendPipDependencies(s []string) *NodeVersionUpdateOne {
	nvuo.mutation.AppendPipDependencies(s)
	return nvuo
}

// SetDeprecated sets the "deprecated" field.
func (nvuo *NodeVersionUpdateOne) SetDeprecated(b bool) *NodeVersionUpdateOne {
	nvuo.mutation.SetDeprecated(b)
	return nvuo
}

// SetNillableDeprecated sets the "deprecated" field if the given value is not nil.
func (nvuo *NodeVersionUpdateOne) SetNillableDeprecated(b *bool) *NodeVersionUpdateOne {
	if b != nil {
		nvuo.SetDeprecated(*b)
	}
	return nvuo
}

// SetNode sets the "node" edge to the Node entity.
func (nvuo *NodeVersionUpdateOne) SetNode(n *Node) *NodeVersionUpdateOne {
	return nvuo.SetNodeID(n.ID)
}

// SetStorageFileID sets the "storage_file" edge to the StorageFile entity by ID.
func (nvuo *NodeVersionUpdateOne) SetStorageFileID(id uuid.UUID) *NodeVersionUpdateOne {
	nvuo.mutation.SetStorageFileID(id)
	return nvuo
}

// SetNillableStorageFileID sets the "storage_file" edge to the StorageFile entity by ID if the given value is not nil.
func (nvuo *NodeVersionUpdateOne) SetNillableStorageFileID(id *uuid.UUID) *NodeVersionUpdateOne {
	if id != nil {
		nvuo = nvuo.SetStorageFileID(*id)
	}
	return nvuo
}

// SetStorageFile sets the "storage_file" edge to the StorageFile entity.
func (nvuo *NodeVersionUpdateOne) SetStorageFile(s *StorageFile) *NodeVersionUpdateOne {
	return nvuo.SetStorageFileID(s.ID)
}

// Mutation returns the NodeVersionMutation object of the builder.
func (nvuo *NodeVersionUpdateOne) Mutation() *NodeVersionMutation {
	return nvuo.mutation
}

// ClearNode clears the "node" edge to the Node entity.
func (nvuo *NodeVersionUpdateOne) ClearNode() *NodeVersionUpdateOne {
	nvuo.mutation.ClearNode()
	return nvuo
}

// ClearStorageFile clears the "storage_file" edge to the StorageFile entity.
func (nvuo *NodeVersionUpdateOne) ClearStorageFile() *NodeVersionUpdateOne {
	nvuo.mutation.ClearStorageFile()
	return nvuo
}

// Where appends a list predicates to the NodeVersionUpdate builder.
func (nvuo *NodeVersionUpdateOne) Where(ps ...predicate.NodeVersion) *NodeVersionUpdateOne {
	nvuo.mutation.Where(ps...)
	return nvuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (nvuo *NodeVersionUpdateOne) Select(field string, fields ...string) *NodeVersionUpdateOne {
	nvuo.fields = append([]string{field}, fields...)
	return nvuo
}

// Save executes the query and returns the updated NodeVersion entity.
func (nvuo *NodeVersionUpdateOne) Save(ctx context.Context) (*NodeVersion, error) {
	nvuo.defaults()
	return withHooks(ctx, nvuo.sqlSave, nvuo.mutation, nvuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (nvuo *NodeVersionUpdateOne) SaveX(ctx context.Context) *NodeVersion {
	node, err := nvuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (nvuo *NodeVersionUpdateOne) Exec(ctx context.Context) error {
	_, err := nvuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (nvuo *NodeVersionUpdateOne) ExecX(ctx context.Context) {
	if err := nvuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (nvuo *NodeVersionUpdateOne) defaults() {
	if _, ok := nvuo.mutation.UpdateTime(); !ok {
		v := nodeversion.UpdateDefaultUpdateTime()
		nvuo.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (nvuo *NodeVersionUpdateOne) check() error {
	if _, ok := nvuo.mutation.NodeID(); nvuo.mutation.NodeCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "NodeVersion.node"`)
	}
	return nil
}

func (nvuo *NodeVersionUpdateOne) sqlSave(ctx context.Context) (_node *NodeVersion, err error) {
	if err := nvuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(nodeversion.Table, nodeversion.Columns, sqlgraph.NewFieldSpec(nodeversion.FieldID, field.TypeUUID))
	id, ok := nvuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "NodeVersion.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := nvuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, nodeversion.FieldID)
		for _, f := range fields {
			if !nodeversion.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != nodeversion.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := nvuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := nvuo.mutation.UpdateTime(); ok {
		_spec.SetField(nodeversion.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := nvuo.mutation.Version(); ok {
		_spec.SetField(nodeversion.FieldVersion, field.TypeString, value)
	}
	if value, ok := nvuo.mutation.Changelog(); ok {
		_spec.SetField(nodeversion.FieldChangelog, field.TypeString, value)
	}
	if nvuo.mutation.ChangelogCleared() {
		_spec.ClearField(nodeversion.FieldChangelog, field.TypeString)
	}
	if value, ok := nvuo.mutation.PipDependencies(); ok {
		_spec.SetField(nodeversion.FieldPipDependencies, field.TypeJSON, value)
	}
	if value, ok := nvuo.mutation.AppendedPipDependencies(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, nodeversion.FieldPipDependencies, value)
		})
	}
	if value, ok := nvuo.mutation.Deprecated(); ok {
		_spec.SetField(nodeversion.FieldDeprecated, field.TypeBool, value)
	}
	if nvuo.mutation.NodeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   nodeversion.NodeTable,
			Columns: []string{nodeversion.NodeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(node.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := nvuo.mutation.NodeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   nodeversion.NodeTable,
			Columns: []string{nodeversion.NodeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(node.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nvuo.mutation.StorageFileCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   nodeversion.StorageFileTable,
			Columns: []string{nodeversion.StorageFileColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(storagefile.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := nvuo.mutation.StorageFileIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   nodeversion.StorageFileTable,
			Columns: []string{nodeversion.StorageFileColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(storagefile.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &NodeVersion{config: nvuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, nvuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{nodeversion.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	nvuo.mutation.done = true
	return _node, nil
}

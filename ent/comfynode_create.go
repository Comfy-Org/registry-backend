// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"registry-backend/ent/comfynode"
	"registry-backend/ent/nodeversion"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ComfyNodeCreate is the builder for creating a ComfyNode entity.
type ComfyNodeCreate struct {
	config
	mutation *ComfyNodeMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreateTime sets the "create_time" field.
func (cnc *ComfyNodeCreate) SetCreateTime(t time.Time) *ComfyNodeCreate {
	cnc.mutation.SetCreateTime(t)
	return cnc
}

// SetNillableCreateTime sets the "create_time" field if the given value is not nil.
func (cnc *ComfyNodeCreate) SetNillableCreateTime(t *time.Time) *ComfyNodeCreate {
	if t != nil {
		cnc.SetCreateTime(*t)
	}
	return cnc
}

// SetUpdateTime sets the "update_time" field.
func (cnc *ComfyNodeCreate) SetUpdateTime(t time.Time) *ComfyNodeCreate {
	cnc.mutation.SetUpdateTime(t)
	return cnc
}

// SetNillableUpdateTime sets the "update_time" field if the given value is not nil.
func (cnc *ComfyNodeCreate) SetNillableUpdateTime(t *time.Time) *ComfyNodeCreate {
	if t != nil {
		cnc.SetUpdateTime(*t)
	}
	return cnc
}

// SetNodeVersionID sets the "node_version_id" field.
func (cnc *ComfyNodeCreate) SetNodeVersionID(u uuid.UUID) *ComfyNodeCreate {
	cnc.mutation.SetNodeVersionID(u)
	return cnc
}

// SetCategory sets the "category" field.
func (cnc *ComfyNodeCreate) SetCategory(s string) *ComfyNodeCreate {
	cnc.mutation.SetCategory(s)
	return cnc
}

// SetNillableCategory sets the "category" field if the given value is not nil.
func (cnc *ComfyNodeCreate) SetNillableCategory(s *string) *ComfyNodeCreate {
	if s != nil {
		cnc.SetCategory(*s)
	}
	return cnc
}

// SetDescription sets the "description" field.
func (cnc *ComfyNodeCreate) SetDescription(s string) *ComfyNodeCreate {
	cnc.mutation.SetDescription(s)
	return cnc
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (cnc *ComfyNodeCreate) SetNillableDescription(s *string) *ComfyNodeCreate {
	if s != nil {
		cnc.SetDescription(*s)
	}
	return cnc
}

// SetInputTypes sets the "input_types" field.
func (cnc *ComfyNodeCreate) SetInputTypes(s string) *ComfyNodeCreate {
	cnc.mutation.SetInputTypes(s)
	return cnc
}

// SetNillableInputTypes sets the "input_types" field if the given value is not nil.
func (cnc *ComfyNodeCreate) SetNillableInputTypes(s *string) *ComfyNodeCreate {
	if s != nil {
		cnc.SetInputTypes(*s)
	}
	return cnc
}

// SetDeprecated sets the "deprecated" field.
func (cnc *ComfyNodeCreate) SetDeprecated(b bool) *ComfyNodeCreate {
	cnc.mutation.SetDeprecated(b)
	return cnc
}

// SetNillableDeprecated sets the "deprecated" field if the given value is not nil.
func (cnc *ComfyNodeCreate) SetNillableDeprecated(b *bool) *ComfyNodeCreate {
	if b != nil {
		cnc.SetDeprecated(*b)
	}
	return cnc
}

// SetExperimental sets the "experimental" field.
func (cnc *ComfyNodeCreate) SetExperimental(b bool) *ComfyNodeCreate {
	cnc.mutation.SetExperimental(b)
	return cnc
}

// SetNillableExperimental sets the "experimental" field if the given value is not nil.
func (cnc *ComfyNodeCreate) SetNillableExperimental(b *bool) *ComfyNodeCreate {
	if b != nil {
		cnc.SetExperimental(*b)
	}
	return cnc
}

// SetOutputIsList sets the "output_is_list" field.
func (cnc *ComfyNodeCreate) SetOutputIsList(b []bool) *ComfyNodeCreate {
	cnc.mutation.SetOutputIsList(b)
	return cnc
}

// SetReturnNames sets the "return_names" field.
func (cnc *ComfyNodeCreate) SetReturnNames(s []string) *ComfyNodeCreate {
	cnc.mutation.SetReturnNames(s)
	return cnc
}

// SetReturnTypes sets the "return_types" field.
func (cnc *ComfyNodeCreate) SetReturnTypes(s string) *ComfyNodeCreate {
	cnc.mutation.SetReturnTypes(s)
	return cnc
}

// SetNillableReturnTypes sets the "return_types" field if the given value is not nil.
func (cnc *ComfyNodeCreate) SetNillableReturnTypes(s *string) *ComfyNodeCreate {
	if s != nil {
		cnc.SetReturnTypes(*s)
	}
	return cnc
}

// SetFunction sets the "function" field.
func (cnc *ComfyNodeCreate) SetFunction(s string) *ComfyNodeCreate {
	cnc.mutation.SetFunction(s)
	return cnc
}

// SetNillableFunction sets the "function" field if the given value is not nil.
func (cnc *ComfyNodeCreate) SetNillableFunction(s *string) *ComfyNodeCreate {
	if s != nil {
		cnc.SetFunction(*s)
	}
	return cnc
}

// SetID sets the "id" field.
func (cnc *ComfyNodeCreate) SetID(s string) *ComfyNodeCreate {
	cnc.mutation.SetID(s)
	return cnc
}

// SetVersionsID sets the "versions" edge to the NodeVersion entity by ID.
func (cnc *ComfyNodeCreate) SetVersionsID(id uuid.UUID) *ComfyNodeCreate {
	cnc.mutation.SetVersionsID(id)
	return cnc
}

// SetVersions sets the "versions" edge to the NodeVersion entity.
func (cnc *ComfyNodeCreate) SetVersions(n *NodeVersion) *ComfyNodeCreate {
	return cnc.SetVersionsID(n.ID)
}

// Mutation returns the ComfyNodeMutation object of the builder.
func (cnc *ComfyNodeCreate) Mutation() *ComfyNodeMutation {
	return cnc.mutation
}

// Save creates the ComfyNode in the database.
func (cnc *ComfyNodeCreate) Save(ctx context.Context) (*ComfyNode, error) {
	cnc.defaults()
	return withHooks(ctx, cnc.sqlSave, cnc.mutation, cnc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (cnc *ComfyNodeCreate) SaveX(ctx context.Context) *ComfyNode {
	v, err := cnc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cnc *ComfyNodeCreate) Exec(ctx context.Context) error {
	_, err := cnc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cnc *ComfyNodeCreate) ExecX(ctx context.Context) {
	if err := cnc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cnc *ComfyNodeCreate) defaults() {
	if _, ok := cnc.mutation.CreateTime(); !ok {
		v := comfynode.DefaultCreateTime()
		cnc.mutation.SetCreateTime(v)
	}
	if _, ok := cnc.mutation.UpdateTime(); !ok {
		v := comfynode.DefaultUpdateTime()
		cnc.mutation.SetUpdateTime(v)
	}
	if _, ok := cnc.mutation.Deprecated(); !ok {
		v := comfynode.DefaultDeprecated
		cnc.mutation.SetDeprecated(v)
	}
	if _, ok := cnc.mutation.Experimental(); !ok {
		v := comfynode.DefaultExperimental
		cnc.mutation.SetExperimental(v)
	}
	if _, ok := cnc.mutation.OutputIsList(); !ok {
		v := comfynode.DefaultOutputIsList
		cnc.mutation.SetOutputIsList(v)
	}
	if _, ok := cnc.mutation.ReturnNames(); !ok {
		v := comfynode.DefaultReturnNames
		cnc.mutation.SetReturnNames(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cnc *ComfyNodeCreate) check() error {
	if _, ok := cnc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "create_time", err: errors.New(`ent: missing required field "ComfyNode.create_time"`)}
	}
	if _, ok := cnc.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "update_time", err: errors.New(`ent: missing required field "ComfyNode.update_time"`)}
	}
	if _, ok := cnc.mutation.NodeVersionID(); !ok {
		return &ValidationError{Name: "node_version_id", err: errors.New(`ent: missing required field "ComfyNode.node_version_id"`)}
	}
	if _, ok := cnc.mutation.Deprecated(); !ok {
		return &ValidationError{Name: "deprecated", err: errors.New(`ent: missing required field "ComfyNode.deprecated"`)}
	}
	if _, ok := cnc.mutation.Experimental(); !ok {
		return &ValidationError{Name: "experimental", err: errors.New(`ent: missing required field "ComfyNode.experimental"`)}
	}
	if _, ok := cnc.mutation.OutputIsList(); !ok {
		return &ValidationError{Name: "output_is_list", err: errors.New(`ent: missing required field "ComfyNode.output_is_list"`)}
	}
	if _, ok := cnc.mutation.ReturnNames(); !ok {
		return &ValidationError{Name: "return_names", err: errors.New(`ent: missing required field "ComfyNode.return_names"`)}
	}
	if _, ok := cnc.mutation.VersionsID(); !ok {
		return &ValidationError{Name: "versions", err: errors.New(`ent: missing required edge "ComfyNode.versions"`)}
	}
	return nil
}

func (cnc *ComfyNodeCreate) sqlSave(ctx context.Context) (*ComfyNode, error) {
	if err := cnc.check(); err != nil {
		return nil, err
	}
	_node, _spec := cnc.createSpec()
	if err := sqlgraph.CreateNode(ctx, cnc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(string); ok {
			_node.ID = id
		} else {
			return nil, fmt.Errorf("unexpected ComfyNode.ID type: %T", _spec.ID.Value)
		}
	}
	cnc.mutation.id = &_node.ID
	cnc.mutation.done = true
	return _node, nil
}

func (cnc *ComfyNodeCreate) createSpec() (*ComfyNode, *sqlgraph.CreateSpec) {
	var (
		_node = &ComfyNode{config: cnc.config}
		_spec = sqlgraph.NewCreateSpec(comfynode.Table, sqlgraph.NewFieldSpec(comfynode.FieldID, field.TypeString))
	)
	_spec.OnConflict = cnc.conflict
	if id, ok := cnc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := cnc.mutation.CreateTime(); ok {
		_spec.SetField(comfynode.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = value
	}
	if value, ok := cnc.mutation.UpdateTime(); ok {
		_spec.SetField(comfynode.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = value
	}
	if value, ok := cnc.mutation.Category(); ok {
		_spec.SetField(comfynode.FieldCategory, field.TypeString, value)
		_node.Category = value
	}
	if value, ok := cnc.mutation.Description(); ok {
		_spec.SetField(comfynode.FieldDescription, field.TypeString, value)
		_node.Description = value
	}
	if value, ok := cnc.mutation.InputTypes(); ok {
		_spec.SetField(comfynode.FieldInputTypes, field.TypeString, value)
		_node.InputTypes = value
	}
	if value, ok := cnc.mutation.Deprecated(); ok {
		_spec.SetField(comfynode.FieldDeprecated, field.TypeBool, value)
		_node.Deprecated = value
	}
	if value, ok := cnc.mutation.Experimental(); ok {
		_spec.SetField(comfynode.FieldExperimental, field.TypeBool, value)
		_node.Experimental = value
	}
	if value, ok := cnc.mutation.OutputIsList(); ok {
		_spec.SetField(comfynode.FieldOutputIsList, field.TypeJSON, value)
		_node.OutputIsList = value
	}
	if value, ok := cnc.mutation.ReturnNames(); ok {
		_spec.SetField(comfynode.FieldReturnNames, field.TypeJSON, value)
		_node.ReturnNames = value
	}
	if value, ok := cnc.mutation.ReturnTypes(); ok {
		_spec.SetField(comfynode.FieldReturnTypes, field.TypeString, value)
		_node.ReturnTypes = value
	}
	if value, ok := cnc.mutation.Function(); ok {
		_spec.SetField(comfynode.FieldFunction, field.TypeString, value)
		_node.Function = value
	}
	if nodes := cnc.mutation.VersionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   comfynode.VersionsTable,
			Columns: []string{comfynode.VersionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(nodeversion.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.NodeVersionID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.ComfyNode.Create().
//		SetCreateTime(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ComfyNodeUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (cnc *ComfyNodeCreate) OnConflict(opts ...sql.ConflictOption) *ComfyNodeUpsertOne {
	cnc.conflict = opts
	return &ComfyNodeUpsertOne{
		create: cnc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ComfyNode.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (cnc *ComfyNodeCreate) OnConflictColumns(columns ...string) *ComfyNodeUpsertOne {
	cnc.conflict = append(cnc.conflict, sql.ConflictColumns(columns...))
	return &ComfyNodeUpsertOne{
		create: cnc,
	}
}

type (
	// ComfyNodeUpsertOne is the builder for "upsert"-ing
	//  one ComfyNode node.
	ComfyNodeUpsertOne struct {
		create *ComfyNodeCreate
	}

	// ComfyNodeUpsert is the "OnConflict" setter.
	ComfyNodeUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdateTime sets the "update_time" field.
func (u *ComfyNodeUpsert) SetUpdateTime(v time.Time) *ComfyNodeUpsert {
	u.Set(comfynode.FieldUpdateTime, v)
	return u
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *ComfyNodeUpsert) UpdateUpdateTime() *ComfyNodeUpsert {
	u.SetExcluded(comfynode.FieldUpdateTime)
	return u
}

// SetNodeVersionID sets the "node_version_id" field.
func (u *ComfyNodeUpsert) SetNodeVersionID(v uuid.UUID) *ComfyNodeUpsert {
	u.Set(comfynode.FieldNodeVersionID, v)
	return u
}

// UpdateNodeVersionID sets the "node_version_id" field to the value that was provided on create.
func (u *ComfyNodeUpsert) UpdateNodeVersionID() *ComfyNodeUpsert {
	u.SetExcluded(comfynode.FieldNodeVersionID)
	return u
}

// SetCategory sets the "category" field.
func (u *ComfyNodeUpsert) SetCategory(v string) *ComfyNodeUpsert {
	u.Set(comfynode.FieldCategory, v)
	return u
}

// UpdateCategory sets the "category" field to the value that was provided on create.
func (u *ComfyNodeUpsert) UpdateCategory() *ComfyNodeUpsert {
	u.SetExcluded(comfynode.FieldCategory)
	return u
}

// ClearCategory clears the value of the "category" field.
func (u *ComfyNodeUpsert) ClearCategory() *ComfyNodeUpsert {
	u.SetNull(comfynode.FieldCategory)
	return u
}

// SetDescription sets the "description" field.
func (u *ComfyNodeUpsert) SetDescription(v string) *ComfyNodeUpsert {
	u.Set(comfynode.FieldDescription, v)
	return u
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *ComfyNodeUpsert) UpdateDescription() *ComfyNodeUpsert {
	u.SetExcluded(comfynode.FieldDescription)
	return u
}

// ClearDescription clears the value of the "description" field.
func (u *ComfyNodeUpsert) ClearDescription() *ComfyNodeUpsert {
	u.SetNull(comfynode.FieldDescription)
	return u
}

// SetInputTypes sets the "input_types" field.
func (u *ComfyNodeUpsert) SetInputTypes(v string) *ComfyNodeUpsert {
	u.Set(comfynode.FieldInputTypes, v)
	return u
}

// UpdateInputTypes sets the "input_types" field to the value that was provided on create.
func (u *ComfyNodeUpsert) UpdateInputTypes() *ComfyNodeUpsert {
	u.SetExcluded(comfynode.FieldInputTypes)
	return u
}

// ClearInputTypes clears the value of the "input_types" field.
func (u *ComfyNodeUpsert) ClearInputTypes() *ComfyNodeUpsert {
	u.SetNull(comfynode.FieldInputTypes)
	return u
}

// SetDeprecated sets the "deprecated" field.
func (u *ComfyNodeUpsert) SetDeprecated(v bool) *ComfyNodeUpsert {
	u.Set(comfynode.FieldDeprecated, v)
	return u
}

// UpdateDeprecated sets the "deprecated" field to the value that was provided on create.
func (u *ComfyNodeUpsert) UpdateDeprecated() *ComfyNodeUpsert {
	u.SetExcluded(comfynode.FieldDeprecated)
	return u
}

// SetExperimental sets the "experimental" field.
func (u *ComfyNodeUpsert) SetExperimental(v bool) *ComfyNodeUpsert {
	u.Set(comfynode.FieldExperimental, v)
	return u
}

// UpdateExperimental sets the "experimental" field to the value that was provided on create.
func (u *ComfyNodeUpsert) UpdateExperimental() *ComfyNodeUpsert {
	u.SetExcluded(comfynode.FieldExperimental)
	return u
}

// SetOutputIsList sets the "output_is_list" field.
func (u *ComfyNodeUpsert) SetOutputIsList(v []bool) *ComfyNodeUpsert {
	u.Set(comfynode.FieldOutputIsList, v)
	return u
}

// UpdateOutputIsList sets the "output_is_list" field to the value that was provided on create.
func (u *ComfyNodeUpsert) UpdateOutputIsList() *ComfyNodeUpsert {
	u.SetExcluded(comfynode.FieldOutputIsList)
	return u
}

// SetReturnNames sets the "return_names" field.
func (u *ComfyNodeUpsert) SetReturnNames(v []string) *ComfyNodeUpsert {
	u.Set(comfynode.FieldReturnNames, v)
	return u
}

// UpdateReturnNames sets the "return_names" field to the value that was provided on create.
func (u *ComfyNodeUpsert) UpdateReturnNames() *ComfyNodeUpsert {
	u.SetExcluded(comfynode.FieldReturnNames)
	return u
}

// SetReturnTypes sets the "return_types" field.
func (u *ComfyNodeUpsert) SetReturnTypes(v string) *ComfyNodeUpsert {
	u.Set(comfynode.FieldReturnTypes, v)
	return u
}

// UpdateReturnTypes sets the "return_types" field to the value that was provided on create.
func (u *ComfyNodeUpsert) UpdateReturnTypes() *ComfyNodeUpsert {
	u.SetExcluded(comfynode.FieldReturnTypes)
	return u
}

// ClearReturnTypes clears the value of the "return_types" field.
func (u *ComfyNodeUpsert) ClearReturnTypes() *ComfyNodeUpsert {
	u.SetNull(comfynode.FieldReturnTypes)
	return u
}

// SetFunction sets the "function" field.
func (u *ComfyNodeUpsert) SetFunction(v string) *ComfyNodeUpsert {
	u.Set(comfynode.FieldFunction, v)
	return u
}

// UpdateFunction sets the "function" field to the value that was provided on create.
func (u *ComfyNodeUpsert) UpdateFunction() *ComfyNodeUpsert {
	u.SetExcluded(comfynode.FieldFunction)
	return u
}

// ClearFunction clears the value of the "function" field.
func (u *ComfyNodeUpsert) ClearFunction() *ComfyNodeUpsert {
	u.SetNull(comfynode.FieldFunction)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.ComfyNode.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(comfynode.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ComfyNodeUpsertOne) UpdateNewValues() *ComfyNodeUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(comfynode.FieldID)
		}
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(comfynode.FieldCreateTime)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ComfyNode.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *ComfyNodeUpsertOne) Ignore() *ComfyNodeUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ComfyNodeUpsertOne) DoNothing() *ComfyNodeUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ComfyNodeCreate.OnConflict
// documentation for more info.
func (u *ComfyNodeUpsertOne) Update(set func(*ComfyNodeUpsert)) *ComfyNodeUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ComfyNodeUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "update_time" field.
func (u *ComfyNodeUpsertOne) SetUpdateTime(v time.Time) *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *ComfyNodeUpsertOne) UpdateUpdateTime() *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetNodeVersionID sets the "node_version_id" field.
func (u *ComfyNodeUpsertOne) SetNodeVersionID(v uuid.UUID) *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetNodeVersionID(v)
	})
}

// UpdateNodeVersionID sets the "node_version_id" field to the value that was provided on create.
func (u *ComfyNodeUpsertOne) UpdateNodeVersionID() *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateNodeVersionID()
	})
}

// SetCategory sets the "category" field.
func (u *ComfyNodeUpsertOne) SetCategory(v string) *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetCategory(v)
	})
}

// UpdateCategory sets the "category" field to the value that was provided on create.
func (u *ComfyNodeUpsertOne) UpdateCategory() *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateCategory()
	})
}

// ClearCategory clears the value of the "category" field.
func (u *ComfyNodeUpsertOne) ClearCategory() *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.ClearCategory()
	})
}

// SetDescription sets the "description" field.
func (u *ComfyNodeUpsertOne) SetDescription(v string) *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetDescription(v)
	})
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *ComfyNodeUpsertOne) UpdateDescription() *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateDescription()
	})
}

// ClearDescription clears the value of the "description" field.
func (u *ComfyNodeUpsertOne) ClearDescription() *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.ClearDescription()
	})
}

// SetInputTypes sets the "input_types" field.
func (u *ComfyNodeUpsertOne) SetInputTypes(v string) *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetInputTypes(v)
	})
}

// UpdateInputTypes sets the "input_types" field to the value that was provided on create.
func (u *ComfyNodeUpsertOne) UpdateInputTypes() *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateInputTypes()
	})
}

// ClearInputTypes clears the value of the "input_types" field.
func (u *ComfyNodeUpsertOne) ClearInputTypes() *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.ClearInputTypes()
	})
}

// SetDeprecated sets the "deprecated" field.
func (u *ComfyNodeUpsertOne) SetDeprecated(v bool) *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetDeprecated(v)
	})
}

// UpdateDeprecated sets the "deprecated" field to the value that was provided on create.
func (u *ComfyNodeUpsertOne) UpdateDeprecated() *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateDeprecated()
	})
}

// SetExperimental sets the "experimental" field.
func (u *ComfyNodeUpsertOne) SetExperimental(v bool) *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetExperimental(v)
	})
}

// UpdateExperimental sets the "experimental" field to the value that was provided on create.
func (u *ComfyNodeUpsertOne) UpdateExperimental() *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateExperimental()
	})
}

// SetOutputIsList sets the "output_is_list" field.
func (u *ComfyNodeUpsertOne) SetOutputIsList(v []bool) *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetOutputIsList(v)
	})
}

// UpdateOutputIsList sets the "output_is_list" field to the value that was provided on create.
func (u *ComfyNodeUpsertOne) UpdateOutputIsList() *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateOutputIsList()
	})
}

// SetReturnNames sets the "return_names" field.
func (u *ComfyNodeUpsertOne) SetReturnNames(v []string) *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetReturnNames(v)
	})
}

// UpdateReturnNames sets the "return_names" field to the value that was provided on create.
func (u *ComfyNodeUpsertOne) UpdateReturnNames() *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateReturnNames()
	})
}

// SetReturnTypes sets the "return_types" field.
func (u *ComfyNodeUpsertOne) SetReturnTypes(v string) *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetReturnTypes(v)
	})
}

// UpdateReturnTypes sets the "return_types" field to the value that was provided on create.
func (u *ComfyNodeUpsertOne) UpdateReturnTypes() *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateReturnTypes()
	})
}

// ClearReturnTypes clears the value of the "return_types" field.
func (u *ComfyNodeUpsertOne) ClearReturnTypes() *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.ClearReturnTypes()
	})
}

// SetFunction sets the "function" field.
func (u *ComfyNodeUpsertOne) SetFunction(v string) *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetFunction(v)
	})
}

// UpdateFunction sets the "function" field to the value that was provided on create.
func (u *ComfyNodeUpsertOne) UpdateFunction() *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateFunction()
	})
}

// ClearFunction clears the value of the "function" field.
func (u *ComfyNodeUpsertOne) ClearFunction() *ComfyNodeUpsertOne {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.ClearFunction()
	})
}

// Exec executes the query.
func (u *ComfyNodeUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for ComfyNodeCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ComfyNodeUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *ComfyNodeUpsertOne) ID(ctx context.Context) (id string, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("ent: ComfyNodeUpsertOne.ID is not supported by MySQL driver. Use ComfyNodeUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *ComfyNodeUpsertOne) IDX(ctx context.Context) string {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// ComfyNodeCreateBulk is the builder for creating many ComfyNode entities in bulk.
type ComfyNodeCreateBulk struct {
	config
	err      error
	builders []*ComfyNodeCreate
	conflict []sql.ConflictOption
}

// Save creates the ComfyNode entities in the database.
func (cncb *ComfyNodeCreateBulk) Save(ctx context.Context) ([]*ComfyNode, error) {
	if cncb.err != nil {
		return nil, cncb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(cncb.builders))
	nodes := make([]*ComfyNode, len(cncb.builders))
	mutators := make([]Mutator, len(cncb.builders))
	for i := range cncb.builders {
		func(i int, root context.Context) {
			builder := cncb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ComfyNodeMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, cncb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = cncb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, cncb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, cncb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (cncb *ComfyNodeCreateBulk) SaveX(ctx context.Context) []*ComfyNode {
	v, err := cncb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cncb *ComfyNodeCreateBulk) Exec(ctx context.Context) error {
	_, err := cncb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cncb *ComfyNodeCreateBulk) ExecX(ctx context.Context) {
	if err := cncb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.ComfyNode.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ComfyNodeUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (cncb *ComfyNodeCreateBulk) OnConflict(opts ...sql.ConflictOption) *ComfyNodeUpsertBulk {
	cncb.conflict = opts
	return &ComfyNodeUpsertBulk{
		create: cncb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ComfyNode.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (cncb *ComfyNodeCreateBulk) OnConflictColumns(columns ...string) *ComfyNodeUpsertBulk {
	cncb.conflict = append(cncb.conflict, sql.ConflictColumns(columns...))
	return &ComfyNodeUpsertBulk{
		create: cncb,
	}
}

// ComfyNodeUpsertBulk is the builder for "upsert"-ing
// a bulk of ComfyNode nodes.
type ComfyNodeUpsertBulk struct {
	create *ComfyNodeCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.ComfyNode.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(comfynode.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ComfyNodeUpsertBulk) UpdateNewValues() *ComfyNodeUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(comfynode.FieldID)
			}
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(comfynode.FieldCreateTime)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ComfyNode.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *ComfyNodeUpsertBulk) Ignore() *ComfyNodeUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ComfyNodeUpsertBulk) DoNothing() *ComfyNodeUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ComfyNodeCreateBulk.OnConflict
// documentation for more info.
func (u *ComfyNodeUpsertBulk) Update(set func(*ComfyNodeUpsert)) *ComfyNodeUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ComfyNodeUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "update_time" field.
func (u *ComfyNodeUpsertBulk) SetUpdateTime(v time.Time) *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *ComfyNodeUpsertBulk) UpdateUpdateTime() *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetNodeVersionID sets the "node_version_id" field.
func (u *ComfyNodeUpsertBulk) SetNodeVersionID(v uuid.UUID) *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetNodeVersionID(v)
	})
}

// UpdateNodeVersionID sets the "node_version_id" field to the value that was provided on create.
func (u *ComfyNodeUpsertBulk) UpdateNodeVersionID() *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateNodeVersionID()
	})
}

// SetCategory sets the "category" field.
func (u *ComfyNodeUpsertBulk) SetCategory(v string) *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetCategory(v)
	})
}

// UpdateCategory sets the "category" field to the value that was provided on create.
func (u *ComfyNodeUpsertBulk) UpdateCategory() *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateCategory()
	})
}

// ClearCategory clears the value of the "category" field.
func (u *ComfyNodeUpsertBulk) ClearCategory() *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.ClearCategory()
	})
}

// SetDescription sets the "description" field.
func (u *ComfyNodeUpsertBulk) SetDescription(v string) *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetDescription(v)
	})
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *ComfyNodeUpsertBulk) UpdateDescription() *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateDescription()
	})
}

// ClearDescription clears the value of the "description" field.
func (u *ComfyNodeUpsertBulk) ClearDescription() *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.ClearDescription()
	})
}

// SetInputTypes sets the "input_types" field.
func (u *ComfyNodeUpsertBulk) SetInputTypes(v string) *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetInputTypes(v)
	})
}

// UpdateInputTypes sets the "input_types" field to the value that was provided on create.
func (u *ComfyNodeUpsertBulk) UpdateInputTypes() *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateInputTypes()
	})
}

// ClearInputTypes clears the value of the "input_types" field.
func (u *ComfyNodeUpsertBulk) ClearInputTypes() *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.ClearInputTypes()
	})
}

// SetDeprecated sets the "deprecated" field.
func (u *ComfyNodeUpsertBulk) SetDeprecated(v bool) *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetDeprecated(v)
	})
}

// UpdateDeprecated sets the "deprecated" field to the value that was provided on create.
func (u *ComfyNodeUpsertBulk) UpdateDeprecated() *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateDeprecated()
	})
}

// SetExperimental sets the "experimental" field.
func (u *ComfyNodeUpsertBulk) SetExperimental(v bool) *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetExperimental(v)
	})
}

// UpdateExperimental sets the "experimental" field to the value that was provided on create.
func (u *ComfyNodeUpsertBulk) UpdateExperimental() *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateExperimental()
	})
}

// SetOutputIsList sets the "output_is_list" field.
func (u *ComfyNodeUpsertBulk) SetOutputIsList(v []bool) *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetOutputIsList(v)
	})
}

// UpdateOutputIsList sets the "output_is_list" field to the value that was provided on create.
func (u *ComfyNodeUpsertBulk) UpdateOutputIsList() *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateOutputIsList()
	})
}

// SetReturnNames sets the "return_names" field.
func (u *ComfyNodeUpsertBulk) SetReturnNames(v []string) *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetReturnNames(v)
	})
}

// UpdateReturnNames sets the "return_names" field to the value that was provided on create.
func (u *ComfyNodeUpsertBulk) UpdateReturnNames() *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateReturnNames()
	})
}

// SetReturnTypes sets the "return_types" field.
func (u *ComfyNodeUpsertBulk) SetReturnTypes(v string) *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetReturnTypes(v)
	})
}

// UpdateReturnTypes sets the "return_types" field to the value that was provided on create.
func (u *ComfyNodeUpsertBulk) UpdateReturnTypes() *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateReturnTypes()
	})
}

// ClearReturnTypes clears the value of the "return_types" field.
func (u *ComfyNodeUpsertBulk) ClearReturnTypes() *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.ClearReturnTypes()
	})
}

// SetFunction sets the "function" field.
func (u *ComfyNodeUpsertBulk) SetFunction(v string) *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.SetFunction(v)
	})
}

// UpdateFunction sets the "function" field to the value that was provided on create.
func (u *ComfyNodeUpsertBulk) UpdateFunction() *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.UpdateFunction()
	})
}

// ClearFunction clears the value of the "function" field.
func (u *ComfyNodeUpsertBulk) ClearFunction() *ComfyNodeUpsertBulk {
	return u.Update(func(s *ComfyNodeUpsert) {
		s.ClearFunction()
	})
}

// Exec executes the query.
func (u *ComfyNodeUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the ComfyNodeCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for ComfyNodeCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ComfyNodeUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

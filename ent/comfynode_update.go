// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"registry-backend/ent/comfynode"
	"registry-backend/ent/nodeversion"
	"registry-backend/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/dialect/sql/sqljson"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ComfyNodeUpdate is the builder for updating ComfyNode entities.
type ComfyNodeUpdate struct {
	config
	hooks     []Hook
	mutation  *ComfyNodeMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the ComfyNodeUpdate builder.
func (cnu *ComfyNodeUpdate) Where(ps ...predicate.ComfyNode) *ComfyNodeUpdate {
	cnu.mutation.Where(ps...)
	return cnu
}

// SetUpdateTime sets the "update_time" field.
func (cnu *ComfyNodeUpdate) SetUpdateTime(t time.Time) *ComfyNodeUpdate {
	cnu.mutation.SetUpdateTime(t)
	return cnu
}

// SetName sets the "name" field.
func (cnu *ComfyNodeUpdate) SetName(s string) *ComfyNodeUpdate {
	cnu.mutation.SetName(s)
	return cnu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (cnu *ComfyNodeUpdate) SetNillableName(s *string) *ComfyNodeUpdate {
	if s != nil {
		cnu.SetName(*s)
	}
	return cnu
}

// SetNodeVersionID sets the "node_version_id" field.
func (cnu *ComfyNodeUpdate) SetNodeVersionID(u uuid.UUID) *ComfyNodeUpdate {
	cnu.mutation.SetNodeVersionID(u)
	return cnu
}

// SetNillableNodeVersionID sets the "node_version_id" field if the given value is not nil.
func (cnu *ComfyNodeUpdate) SetNillableNodeVersionID(u *uuid.UUID) *ComfyNodeUpdate {
	if u != nil {
		cnu.SetNodeVersionID(*u)
	}
	return cnu
}

// SetCategory sets the "category" field.
func (cnu *ComfyNodeUpdate) SetCategory(s string) *ComfyNodeUpdate {
	cnu.mutation.SetCategory(s)
	return cnu
}

// SetNillableCategory sets the "category" field if the given value is not nil.
func (cnu *ComfyNodeUpdate) SetNillableCategory(s *string) *ComfyNodeUpdate {
	if s != nil {
		cnu.SetCategory(*s)
	}
	return cnu
}

// ClearCategory clears the value of the "category" field.
func (cnu *ComfyNodeUpdate) ClearCategory() *ComfyNodeUpdate {
	cnu.mutation.ClearCategory()
	return cnu
}

// SetDescription sets the "description" field.
func (cnu *ComfyNodeUpdate) SetDescription(s string) *ComfyNodeUpdate {
	cnu.mutation.SetDescription(s)
	return cnu
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (cnu *ComfyNodeUpdate) SetNillableDescription(s *string) *ComfyNodeUpdate {
	if s != nil {
		cnu.SetDescription(*s)
	}
	return cnu
}

// ClearDescription clears the value of the "description" field.
func (cnu *ComfyNodeUpdate) ClearDescription() *ComfyNodeUpdate {
	cnu.mutation.ClearDescription()
	return cnu
}

// SetInputTypes sets the "input_types" field.
func (cnu *ComfyNodeUpdate) SetInputTypes(s string) *ComfyNodeUpdate {
	cnu.mutation.SetInputTypes(s)
	return cnu
}

// SetNillableInputTypes sets the "input_types" field if the given value is not nil.
func (cnu *ComfyNodeUpdate) SetNillableInputTypes(s *string) *ComfyNodeUpdate {
	if s != nil {
		cnu.SetInputTypes(*s)
	}
	return cnu
}

// ClearInputTypes clears the value of the "input_types" field.
func (cnu *ComfyNodeUpdate) ClearInputTypes() *ComfyNodeUpdate {
	cnu.mutation.ClearInputTypes()
	return cnu
}

// SetDeprecated sets the "deprecated" field.
func (cnu *ComfyNodeUpdate) SetDeprecated(b bool) *ComfyNodeUpdate {
	cnu.mutation.SetDeprecated(b)
	return cnu
}

// SetNillableDeprecated sets the "deprecated" field if the given value is not nil.
func (cnu *ComfyNodeUpdate) SetNillableDeprecated(b *bool) *ComfyNodeUpdate {
	if b != nil {
		cnu.SetDeprecated(*b)
	}
	return cnu
}

// ClearDeprecated clears the value of the "deprecated" field.
func (cnu *ComfyNodeUpdate) ClearDeprecated() *ComfyNodeUpdate {
	cnu.mutation.ClearDeprecated()
	return cnu
}

// SetExperimental sets the "experimental" field.
func (cnu *ComfyNodeUpdate) SetExperimental(b bool) *ComfyNodeUpdate {
	cnu.mutation.SetExperimental(b)
	return cnu
}

// SetNillableExperimental sets the "experimental" field if the given value is not nil.
func (cnu *ComfyNodeUpdate) SetNillableExperimental(b *bool) *ComfyNodeUpdate {
	if b != nil {
		cnu.SetExperimental(*b)
	}
	return cnu
}

// ClearExperimental clears the value of the "experimental" field.
func (cnu *ComfyNodeUpdate) ClearExperimental() *ComfyNodeUpdate {
	cnu.mutation.ClearExperimental()
	return cnu
}

// SetOutputIsList sets the "output_is_list" field.
func (cnu *ComfyNodeUpdate) SetOutputIsList(b []bool) *ComfyNodeUpdate {
	cnu.mutation.SetOutputIsList(b)
	return cnu
}

// AppendOutputIsList appends b to the "output_is_list" field.
func (cnu *ComfyNodeUpdate) AppendOutputIsList(b []bool) *ComfyNodeUpdate {
	cnu.mutation.AppendOutputIsList(b)
	return cnu
}

// ClearOutputIsList clears the value of the "output_is_list" field.
func (cnu *ComfyNodeUpdate) ClearOutputIsList() *ComfyNodeUpdate {
	cnu.mutation.ClearOutputIsList()
	return cnu
}

// SetReturnNames sets the "return_names" field.
func (cnu *ComfyNodeUpdate) SetReturnNames(s string) *ComfyNodeUpdate {
	cnu.mutation.SetReturnNames(s)
	return cnu
}

// SetNillableReturnNames sets the "return_names" field if the given value is not nil.
func (cnu *ComfyNodeUpdate) SetNillableReturnNames(s *string) *ComfyNodeUpdate {
	if s != nil {
		cnu.SetReturnNames(*s)
	}
	return cnu
}

// ClearReturnNames clears the value of the "return_names" field.
func (cnu *ComfyNodeUpdate) ClearReturnNames() *ComfyNodeUpdate {
	cnu.mutation.ClearReturnNames()
	return cnu
}

// SetReturnTypes sets the "return_types" field.
func (cnu *ComfyNodeUpdate) SetReturnTypes(s string) *ComfyNodeUpdate {
	cnu.mutation.SetReturnTypes(s)
	return cnu
}

// SetNillableReturnTypes sets the "return_types" field if the given value is not nil.
func (cnu *ComfyNodeUpdate) SetNillableReturnTypes(s *string) *ComfyNodeUpdate {
	if s != nil {
		cnu.SetReturnTypes(*s)
	}
	return cnu
}

// ClearReturnTypes clears the value of the "return_types" field.
func (cnu *ComfyNodeUpdate) ClearReturnTypes() *ComfyNodeUpdate {
	cnu.mutation.ClearReturnTypes()
	return cnu
}

// SetFunction sets the "function" field.
func (cnu *ComfyNodeUpdate) SetFunction(s string) *ComfyNodeUpdate {
	cnu.mutation.SetFunction(s)
	return cnu
}

// SetNillableFunction sets the "function" field if the given value is not nil.
func (cnu *ComfyNodeUpdate) SetNillableFunction(s *string) *ComfyNodeUpdate {
	if s != nil {
		cnu.SetFunction(*s)
	}
	return cnu
}

// ClearFunction clears the value of the "function" field.
func (cnu *ComfyNodeUpdate) ClearFunction() *ComfyNodeUpdate {
	cnu.mutation.ClearFunction()
	return cnu
}

// SetVersionsID sets the "versions" edge to the NodeVersion entity by ID.
func (cnu *ComfyNodeUpdate) SetVersionsID(id uuid.UUID) *ComfyNodeUpdate {
	cnu.mutation.SetVersionsID(id)
	return cnu
}

// SetVersions sets the "versions" edge to the NodeVersion entity.
func (cnu *ComfyNodeUpdate) SetVersions(n *NodeVersion) *ComfyNodeUpdate {
	return cnu.SetVersionsID(n.ID)
}

// Mutation returns the ComfyNodeMutation object of the builder.
func (cnu *ComfyNodeUpdate) Mutation() *ComfyNodeMutation {
	return cnu.mutation
}

// ClearVersions clears the "versions" edge to the NodeVersion entity.
func (cnu *ComfyNodeUpdate) ClearVersions() *ComfyNodeUpdate {
	cnu.mutation.ClearVersions()
	return cnu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (cnu *ComfyNodeUpdate) Save(ctx context.Context) (int, error) {
	cnu.defaults()
	return withHooks(ctx, cnu.sqlSave, cnu.mutation, cnu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cnu *ComfyNodeUpdate) SaveX(ctx context.Context) int {
	affected, err := cnu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cnu *ComfyNodeUpdate) Exec(ctx context.Context) error {
	_, err := cnu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cnu *ComfyNodeUpdate) ExecX(ctx context.Context) {
	if err := cnu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cnu *ComfyNodeUpdate) defaults() {
	if _, ok := cnu.mutation.UpdateTime(); !ok {
		v := comfynode.UpdateDefaultUpdateTime()
		cnu.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cnu *ComfyNodeUpdate) check() error {
	if cnu.mutation.VersionsCleared() && len(cnu.mutation.VersionsIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "ComfyNode.versions"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (cnu *ComfyNodeUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ComfyNodeUpdate {
	cnu.modifiers = append(cnu.modifiers, modifiers...)
	return cnu
}

func (cnu *ComfyNodeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := cnu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(comfynode.Table, comfynode.Columns, sqlgraph.NewFieldSpec(comfynode.FieldID, field.TypeUUID))
	if ps := cnu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cnu.mutation.UpdateTime(); ok {
		_spec.SetField(comfynode.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := cnu.mutation.Name(); ok {
		_spec.SetField(comfynode.FieldName, field.TypeString, value)
	}
	if value, ok := cnu.mutation.Category(); ok {
		_spec.SetField(comfynode.FieldCategory, field.TypeString, value)
	}
	if cnu.mutation.CategoryCleared() {
		_spec.ClearField(comfynode.FieldCategory, field.TypeString)
	}
	if value, ok := cnu.mutation.Description(); ok {
		_spec.SetField(comfynode.FieldDescription, field.TypeString, value)
	}
	if cnu.mutation.DescriptionCleared() {
		_spec.ClearField(comfynode.FieldDescription, field.TypeString)
	}
	if value, ok := cnu.mutation.InputTypes(); ok {
		_spec.SetField(comfynode.FieldInputTypes, field.TypeString, value)
	}
	if cnu.mutation.InputTypesCleared() {
		_spec.ClearField(comfynode.FieldInputTypes, field.TypeString)
	}
	if value, ok := cnu.mutation.Deprecated(); ok {
		_spec.SetField(comfynode.FieldDeprecated, field.TypeBool, value)
	}
	if cnu.mutation.DeprecatedCleared() {
		_spec.ClearField(comfynode.FieldDeprecated, field.TypeBool)
	}
	if value, ok := cnu.mutation.Experimental(); ok {
		_spec.SetField(comfynode.FieldExperimental, field.TypeBool, value)
	}
	if cnu.mutation.ExperimentalCleared() {
		_spec.ClearField(comfynode.FieldExperimental, field.TypeBool)
	}
	if value, ok := cnu.mutation.OutputIsList(); ok {
		_spec.SetField(comfynode.FieldOutputIsList, field.TypeJSON, value)
	}
	if value, ok := cnu.mutation.AppendedOutputIsList(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, comfynode.FieldOutputIsList, value)
		})
	}
	if cnu.mutation.OutputIsListCleared() {
		_spec.ClearField(comfynode.FieldOutputIsList, field.TypeJSON)
	}
	if value, ok := cnu.mutation.ReturnNames(); ok {
		_spec.SetField(comfynode.FieldReturnNames, field.TypeString, value)
	}
	if cnu.mutation.ReturnNamesCleared() {
		_spec.ClearField(comfynode.FieldReturnNames, field.TypeString)
	}
	if value, ok := cnu.mutation.ReturnTypes(); ok {
		_spec.SetField(comfynode.FieldReturnTypes, field.TypeString, value)
	}
	if cnu.mutation.ReturnTypesCleared() {
		_spec.ClearField(comfynode.FieldReturnTypes, field.TypeString)
	}
	if value, ok := cnu.mutation.Function(); ok {
		_spec.SetField(comfynode.FieldFunction, field.TypeString, value)
	}
	if cnu.mutation.FunctionCleared() {
		_spec.ClearField(comfynode.FieldFunction, field.TypeString)
	}
	if cnu.mutation.VersionsCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cnu.mutation.VersionsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(cnu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, cnu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{comfynode.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	cnu.mutation.done = true
	return n, nil
}

// ComfyNodeUpdateOne is the builder for updating a single ComfyNode entity.
type ComfyNodeUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *ComfyNodeMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdateTime sets the "update_time" field.
func (cnuo *ComfyNodeUpdateOne) SetUpdateTime(t time.Time) *ComfyNodeUpdateOne {
	cnuo.mutation.SetUpdateTime(t)
	return cnuo
}

// SetName sets the "name" field.
func (cnuo *ComfyNodeUpdateOne) SetName(s string) *ComfyNodeUpdateOne {
	cnuo.mutation.SetName(s)
	return cnuo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (cnuo *ComfyNodeUpdateOne) SetNillableName(s *string) *ComfyNodeUpdateOne {
	if s != nil {
		cnuo.SetName(*s)
	}
	return cnuo
}

// SetNodeVersionID sets the "node_version_id" field.
func (cnuo *ComfyNodeUpdateOne) SetNodeVersionID(u uuid.UUID) *ComfyNodeUpdateOne {
	cnuo.mutation.SetNodeVersionID(u)
	return cnuo
}

// SetNillableNodeVersionID sets the "node_version_id" field if the given value is not nil.
func (cnuo *ComfyNodeUpdateOne) SetNillableNodeVersionID(u *uuid.UUID) *ComfyNodeUpdateOne {
	if u != nil {
		cnuo.SetNodeVersionID(*u)
	}
	return cnuo
}

// SetCategory sets the "category" field.
func (cnuo *ComfyNodeUpdateOne) SetCategory(s string) *ComfyNodeUpdateOne {
	cnuo.mutation.SetCategory(s)
	return cnuo
}

// SetNillableCategory sets the "category" field if the given value is not nil.
func (cnuo *ComfyNodeUpdateOne) SetNillableCategory(s *string) *ComfyNodeUpdateOne {
	if s != nil {
		cnuo.SetCategory(*s)
	}
	return cnuo
}

// ClearCategory clears the value of the "category" field.
func (cnuo *ComfyNodeUpdateOne) ClearCategory() *ComfyNodeUpdateOne {
	cnuo.mutation.ClearCategory()
	return cnuo
}

// SetDescription sets the "description" field.
func (cnuo *ComfyNodeUpdateOne) SetDescription(s string) *ComfyNodeUpdateOne {
	cnuo.mutation.SetDescription(s)
	return cnuo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (cnuo *ComfyNodeUpdateOne) SetNillableDescription(s *string) *ComfyNodeUpdateOne {
	if s != nil {
		cnuo.SetDescription(*s)
	}
	return cnuo
}

// ClearDescription clears the value of the "description" field.
func (cnuo *ComfyNodeUpdateOne) ClearDescription() *ComfyNodeUpdateOne {
	cnuo.mutation.ClearDescription()
	return cnuo
}

// SetInputTypes sets the "input_types" field.
func (cnuo *ComfyNodeUpdateOne) SetInputTypes(s string) *ComfyNodeUpdateOne {
	cnuo.mutation.SetInputTypes(s)
	return cnuo
}

// SetNillableInputTypes sets the "input_types" field if the given value is not nil.
func (cnuo *ComfyNodeUpdateOne) SetNillableInputTypes(s *string) *ComfyNodeUpdateOne {
	if s != nil {
		cnuo.SetInputTypes(*s)
	}
	return cnuo
}

// ClearInputTypes clears the value of the "input_types" field.
func (cnuo *ComfyNodeUpdateOne) ClearInputTypes() *ComfyNodeUpdateOne {
	cnuo.mutation.ClearInputTypes()
	return cnuo
}

// SetDeprecated sets the "deprecated" field.
func (cnuo *ComfyNodeUpdateOne) SetDeprecated(b bool) *ComfyNodeUpdateOne {
	cnuo.mutation.SetDeprecated(b)
	return cnuo
}

// SetNillableDeprecated sets the "deprecated" field if the given value is not nil.
func (cnuo *ComfyNodeUpdateOne) SetNillableDeprecated(b *bool) *ComfyNodeUpdateOne {
	if b != nil {
		cnuo.SetDeprecated(*b)
	}
	return cnuo
}

// ClearDeprecated clears the value of the "deprecated" field.
func (cnuo *ComfyNodeUpdateOne) ClearDeprecated() *ComfyNodeUpdateOne {
	cnuo.mutation.ClearDeprecated()
	return cnuo
}

// SetExperimental sets the "experimental" field.
func (cnuo *ComfyNodeUpdateOne) SetExperimental(b bool) *ComfyNodeUpdateOne {
	cnuo.mutation.SetExperimental(b)
	return cnuo
}

// SetNillableExperimental sets the "experimental" field if the given value is not nil.
func (cnuo *ComfyNodeUpdateOne) SetNillableExperimental(b *bool) *ComfyNodeUpdateOne {
	if b != nil {
		cnuo.SetExperimental(*b)
	}
	return cnuo
}

// ClearExperimental clears the value of the "experimental" field.
func (cnuo *ComfyNodeUpdateOne) ClearExperimental() *ComfyNodeUpdateOne {
	cnuo.mutation.ClearExperimental()
	return cnuo
}

// SetOutputIsList sets the "output_is_list" field.
func (cnuo *ComfyNodeUpdateOne) SetOutputIsList(b []bool) *ComfyNodeUpdateOne {
	cnuo.mutation.SetOutputIsList(b)
	return cnuo
}

// AppendOutputIsList appends b to the "output_is_list" field.
func (cnuo *ComfyNodeUpdateOne) AppendOutputIsList(b []bool) *ComfyNodeUpdateOne {
	cnuo.mutation.AppendOutputIsList(b)
	return cnuo
}

// ClearOutputIsList clears the value of the "output_is_list" field.
func (cnuo *ComfyNodeUpdateOne) ClearOutputIsList() *ComfyNodeUpdateOne {
	cnuo.mutation.ClearOutputIsList()
	return cnuo
}

// SetReturnNames sets the "return_names" field.
func (cnuo *ComfyNodeUpdateOne) SetReturnNames(s string) *ComfyNodeUpdateOne {
	cnuo.mutation.SetReturnNames(s)
	return cnuo
}

// SetNillableReturnNames sets the "return_names" field if the given value is not nil.
func (cnuo *ComfyNodeUpdateOne) SetNillableReturnNames(s *string) *ComfyNodeUpdateOne {
	if s != nil {
		cnuo.SetReturnNames(*s)
	}
	return cnuo
}

// ClearReturnNames clears the value of the "return_names" field.
func (cnuo *ComfyNodeUpdateOne) ClearReturnNames() *ComfyNodeUpdateOne {
	cnuo.mutation.ClearReturnNames()
	return cnuo
}

// SetReturnTypes sets the "return_types" field.
func (cnuo *ComfyNodeUpdateOne) SetReturnTypes(s string) *ComfyNodeUpdateOne {
	cnuo.mutation.SetReturnTypes(s)
	return cnuo
}

// SetNillableReturnTypes sets the "return_types" field if the given value is not nil.
func (cnuo *ComfyNodeUpdateOne) SetNillableReturnTypes(s *string) *ComfyNodeUpdateOne {
	if s != nil {
		cnuo.SetReturnTypes(*s)
	}
	return cnuo
}

// ClearReturnTypes clears the value of the "return_types" field.
func (cnuo *ComfyNodeUpdateOne) ClearReturnTypes() *ComfyNodeUpdateOne {
	cnuo.mutation.ClearReturnTypes()
	return cnuo
}

// SetFunction sets the "function" field.
func (cnuo *ComfyNodeUpdateOne) SetFunction(s string) *ComfyNodeUpdateOne {
	cnuo.mutation.SetFunction(s)
	return cnuo
}

// SetNillableFunction sets the "function" field if the given value is not nil.
func (cnuo *ComfyNodeUpdateOne) SetNillableFunction(s *string) *ComfyNodeUpdateOne {
	if s != nil {
		cnuo.SetFunction(*s)
	}
	return cnuo
}

// ClearFunction clears the value of the "function" field.
func (cnuo *ComfyNodeUpdateOne) ClearFunction() *ComfyNodeUpdateOne {
	cnuo.mutation.ClearFunction()
	return cnuo
}

// SetVersionsID sets the "versions" edge to the NodeVersion entity by ID.
func (cnuo *ComfyNodeUpdateOne) SetVersionsID(id uuid.UUID) *ComfyNodeUpdateOne {
	cnuo.mutation.SetVersionsID(id)
	return cnuo
}

// SetVersions sets the "versions" edge to the NodeVersion entity.
func (cnuo *ComfyNodeUpdateOne) SetVersions(n *NodeVersion) *ComfyNodeUpdateOne {
	return cnuo.SetVersionsID(n.ID)
}

// Mutation returns the ComfyNodeMutation object of the builder.
func (cnuo *ComfyNodeUpdateOne) Mutation() *ComfyNodeMutation {
	return cnuo.mutation
}

// ClearVersions clears the "versions" edge to the NodeVersion entity.
func (cnuo *ComfyNodeUpdateOne) ClearVersions() *ComfyNodeUpdateOne {
	cnuo.mutation.ClearVersions()
	return cnuo
}

// Where appends a list predicates to the ComfyNodeUpdate builder.
func (cnuo *ComfyNodeUpdateOne) Where(ps ...predicate.ComfyNode) *ComfyNodeUpdateOne {
	cnuo.mutation.Where(ps...)
	return cnuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (cnuo *ComfyNodeUpdateOne) Select(field string, fields ...string) *ComfyNodeUpdateOne {
	cnuo.fields = append([]string{field}, fields...)
	return cnuo
}

// Save executes the query and returns the updated ComfyNode entity.
func (cnuo *ComfyNodeUpdateOne) Save(ctx context.Context) (*ComfyNode, error) {
	cnuo.defaults()
	return withHooks(ctx, cnuo.sqlSave, cnuo.mutation, cnuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cnuo *ComfyNodeUpdateOne) SaveX(ctx context.Context) *ComfyNode {
	node, err := cnuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (cnuo *ComfyNodeUpdateOne) Exec(ctx context.Context) error {
	_, err := cnuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cnuo *ComfyNodeUpdateOne) ExecX(ctx context.Context) {
	if err := cnuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cnuo *ComfyNodeUpdateOne) defaults() {
	if _, ok := cnuo.mutation.UpdateTime(); !ok {
		v := comfynode.UpdateDefaultUpdateTime()
		cnuo.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cnuo *ComfyNodeUpdateOne) check() error {
	if cnuo.mutation.VersionsCleared() && len(cnuo.mutation.VersionsIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "ComfyNode.versions"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (cnuo *ComfyNodeUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ComfyNodeUpdateOne {
	cnuo.modifiers = append(cnuo.modifiers, modifiers...)
	return cnuo
}

func (cnuo *ComfyNodeUpdateOne) sqlSave(ctx context.Context) (_node *ComfyNode, err error) {
	if err := cnuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(comfynode.Table, comfynode.Columns, sqlgraph.NewFieldSpec(comfynode.FieldID, field.TypeUUID))
	id, ok := cnuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "ComfyNode.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := cnuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, comfynode.FieldID)
		for _, f := range fields {
			if !comfynode.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != comfynode.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := cnuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cnuo.mutation.UpdateTime(); ok {
		_spec.SetField(comfynode.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := cnuo.mutation.Name(); ok {
		_spec.SetField(comfynode.FieldName, field.TypeString, value)
	}
	if value, ok := cnuo.mutation.Category(); ok {
		_spec.SetField(comfynode.FieldCategory, field.TypeString, value)
	}
	if cnuo.mutation.CategoryCleared() {
		_spec.ClearField(comfynode.FieldCategory, field.TypeString)
	}
	if value, ok := cnuo.mutation.Description(); ok {
		_spec.SetField(comfynode.FieldDescription, field.TypeString, value)
	}
	if cnuo.mutation.DescriptionCleared() {
		_spec.ClearField(comfynode.FieldDescription, field.TypeString)
	}
	if value, ok := cnuo.mutation.InputTypes(); ok {
		_spec.SetField(comfynode.FieldInputTypes, field.TypeString, value)
	}
	if cnuo.mutation.InputTypesCleared() {
		_spec.ClearField(comfynode.FieldInputTypes, field.TypeString)
	}
	if value, ok := cnuo.mutation.Deprecated(); ok {
		_spec.SetField(comfynode.FieldDeprecated, field.TypeBool, value)
	}
	if cnuo.mutation.DeprecatedCleared() {
		_spec.ClearField(comfynode.FieldDeprecated, field.TypeBool)
	}
	if value, ok := cnuo.mutation.Experimental(); ok {
		_spec.SetField(comfynode.FieldExperimental, field.TypeBool, value)
	}
	if cnuo.mutation.ExperimentalCleared() {
		_spec.ClearField(comfynode.FieldExperimental, field.TypeBool)
	}
	if value, ok := cnuo.mutation.OutputIsList(); ok {
		_spec.SetField(comfynode.FieldOutputIsList, field.TypeJSON, value)
	}
	if value, ok := cnuo.mutation.AppendedOutputIsList(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, comfynode.FieldOutputIsList, value)
		})
	}
	if cnuo.mutation.OutputIsListCleared() {
		_spec.ClearField(comfynode.FieldOutputIsList, field.TypeJSON)
	}
	if value, ok := cnuo.mutation.ReturnNames(); ok {
		_spec.SetField(comfynode.FieldReturnNames, field.TypeString, value)
	}
	if cnuo.mutation.ReturnNamesCleared() {
		_spec.ClearField(comfynode.FieldReturnNames, field.TypeString)
	}
	if value, ok := cnuo.mutation.ReturnTypes(); ok {
		_spec.SetField(comfynode.FieldReturnTypes, field.TypeString, value)
	}
	if cnuo.mutation.ReturnTypesCleared() {
		_spec.ClearField(comfynode.FieldReturnTypes, field.TypeString)
	}
	if value, ok := cnuo.mutation.Function(); ok {
		_spec.SetField(comfynode.FieldFunction, field.TypeString, value)
	}
	if cnuo.mutation.FunctionCleared() {
		_spec.ClearField(comfynode.FieldFunction, field.TypeString)
	}
	if cnuo.mutation.VersionsCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cnuo.mutation.VersionsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(cnuo.modifiers...)
	_node = &ComfyNode{config: cnuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, cnuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{comfynode.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	cnuo.mutation.done = true
	return _node, nil
}

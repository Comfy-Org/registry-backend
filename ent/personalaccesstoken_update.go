// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"registry-backend/ent/personalaccesstoken"
	"registry-backend/ent/predicate"
	"registry-backend/ent/publisher"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// PersonalAccessTokenUpdate is the builder for updating PersonalAccessToken entities.
type PersonalAccessTokenUpdate struct {
	config
	hooks     []Hook
	mutation  *PersonalAccessTokenMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the PersonalAccessTokenUpdate builder.
func (patu *PersonalAccessTokenUpdate) Where(ps ...predicate.PersonalAccessToken) *PersonalAccessTokenUpdate {
	patu.mutation.Where(ps...)
	return patu
}

// SetUpdateTime sets the "update_time" field.
func (patu *PersonalAccessTokenUpdate) SetUpdateTime(t time.Time) *PersonalAccessTokenUpdate {
	patu.mutation.SetUpdateTime(t)
	return patu
}

// SetName sets the "name" field.
func (patu *PersonalAccessTokenUpdate) SetName(s string) *PersonalAccessTokenUpdate {
	patu.mutation.SetName(s)
	return patu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (patu *PersonalAccessTokenUpdate) SetNillableName(s *string) *PersonalAccessTokenUpdate {
	if s != nil {
		patu.SetName(*s)
	}
	return patu
}

// SetDescription sets the "description" field.
func (patu *PersonalAccessTokenUpdate) SetDescription(s string) *PersonalAccessTokenUpdate {
	patu.mutation.SetDescription(s)
	return patu
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (patu *PersonalAccessTokenUpdate) SetNillableDescription(s *string) *PersonalAccessTokenUpdate {
	if s != nil {
		patu.SetDescription(*s)
	}
	return patu
}

// SetPublisherID sets the "publisher_id" field.
func (patu *PersonalAccessTokenUpdate) SetPublisherID(s string) *PersonalAccessTokenUpdate {
	patu.mutation.SetPublisherID(s)
	return patu
}

// SetNillablePublisherID sets the "publisher_id" field if the given value is not nil.
func (patu *PersonalAccessTokenUpdate) SetNillablePublisherID(s *string) *PersonalAccessTokenUpdate {
	if s != nil {
		patu.SetPublisherID(*s)
	}
	return patu
}

// SetToken sets the "token" field.
func (patu *PersonalAccessTokenUpdate) SetToken(s string) *PersonalAccessTokenUpdate {
	patu.mutation.SetToken(s)
	return patu
}

// SetNillableToken sets the "token" field if the given value is not nil.
func (patu *PersonalAccessTokenUpdate) SetNillableToken(s *string) *PersonalAccessTokenUpdate {
	if s != nil {
		patu.SetToken(*s)
	}
	return patu
}

// SetPublisher sets the "publisher" edge to the Publisher entity.
func (patu *PersonalAccessTokenUpdate) SetPublisher(p *Publisher) *PersonalAccessTokenUpdate {
	return patu.SetPublisherID(p.ID)
}

// Mutation returns the PersonalAccessTokenMutation object of the builder.
func (patu *PersonalAccessTokenUpdate) Mutation() *PersonalAccessTokenMutation {
	return patu.mutation
}

// ClearPublisher clears the "publisher" edge to the Publisher entity.
func (patu *PersonalAccessTokenUpdate) ClearPublisher() *PersonalAccessTokenUpdate {
	patu.mutation.ClearPublisher()
	return patu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (patu *PersonalAccessTokenUpdate) Save(ctx context.Context) (int, error) {
	patu.defaults()
	return withHooks(ctx, patu.sqlSave, patu.mutation, patu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (patu *PersonalAccessTokenUpdate) SaveX(ctx context.Context) int {
	affected, err := patu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (patu *PersonalAccessTokenUpdate) Exec(ctx context.Context) error {
	_, err := patu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (patu *PersonalAccessTokenUpdate) ExecX(ctx context.Context) {
	if err := patu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (patu *PersonalAccessTokenUpdate) defaults() {
	if _, ok := patu.mutation.UpdateTime(); !ok {
		v := personalaccesstoken.UpdateDefaultUpdateTime()
		patu.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (patu *PersonalAccessTokenUpdate) check() error {
	if patu.mutation.PublisherCleared() && len(patu.mutation.PublisherIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "PersonalAccessToken.publisher"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (patu *PersonalAccessTokenUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *PersonalAccessTokenUpdate {
	patu.modifiers = append(patu.modifiers, modifiers...)
	return patu
}

func (patu *PersonalAccessTokenUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := patu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(personalaccesstoken.Table, personalaccesstoken.Columns, sqlgraph.NewFieldSpec(personalaccesstoken.FieldID, field.TypeUUID))
	if ps := patu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := patu.mutation.UpdateTime(); ok {
		_spec.SetField(personalaccesstoken.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := patu.mutation.Name(); ok {
		_spec.SetField(personalaccesstoken.FieldName, field.TypeString, value)
	}
	if value, ok := patu.mutation.Description(); ok {
		_spec.SetField(personalaccesstoken.FieldDescription, field.TypeString, value)
	}
	if value, ok := patu.mutation.Token(); ok {
		_spec.SetField(personalaccesstoken.FieldToken, field.TypeString, value)
	}
	if patu.mutation.PublisherCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   personalaccesstoken.PublisherTable,
			Columns: []string{personalaccesstoken.PublisherColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(publisher.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := patu.mutation.PublisherIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   personalaccesstoken.PublisherTable,
			Columns: []string{personalaccesstoken.PublisherColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(publisher.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(patu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, patu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{personalaccesstoken.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	patu.mutation.done = true
	return n, nil
}

// PersonalAccessTokenUpdateOne is the builder for updating a single PersonalAccessToken entity.
type PersonalAccessTokenUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *PersonalAccessTokenMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdateTime sets the "update_time" field.
func (patuo *PersonalAccessTokenUpdateOne) SetUpdateTime(t time.Time) *PersonalAccessTokenUpdateOne {
	patuo.mutation.SetUpdateTime(t)
	return patuo
}

// SetName sets the "name" field.
func (patuo *PersonalAccessTokenUpdateOne) SetName(s string) *PersonalAccessTokenUpdateOne {
	patuo.mutation.SetName(s)
	return patuo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (patuo *PersonalAccessTokenUpdateOne) SetNillableName(s *string) *PersonalAccessTokenUpdateOne {
	if s != nil {
		patuo.SetName(*s)
	}
	return patuo
}

// SetDescription sets the "description" field.
func (patuo *PersonalAccessTokenUpdateOne) SetDescription(s string) *PersonalAccessTokenUpdateOne {
	patuo.mutation.SetDescription(s)
	return patuo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (patuo *PersonalAccessTokenUpdateOne) SetNillableDescription(s *string) *PersonalAccessTokenUpdateOne {
	if s != nil {
		patuo.SetDescription(*s)
	}
	return patuo
}

// SetPublisherID sets the "publisher_id" field.
func (patuo *PersonalAccessTokenUpdateOne) SetPublisherID(s string) *PersonalAccessTokenUpdateOne {
	patuo.mutation.SetPublisherID(s)
	return patuo
}

// SetNillablePublisherID sets the "publisher_id" field if the given value is not nil.
func (patuo *PersonalAccessTokenUpdateOne) SetNillablePublisherID(s *string) *PersonalAccessTokenUpdateOne {
	if s != nil {
		patuo.SetPublisherID(*s)
	}
	return patuo
}

// SetToken sets the "token" field.
func (patuo *PersonalAccessTokenUpdateOne) SetToken(s string) *PersonalAccessTokenUpdateOne {
	patuo.mutation.SetToken(s)
	return patuo
}

// SetNillableToken sets the "token" field if the given value is not nil.
func (patuo *PersonalAccessTokenUpdateOne) SetNillableToken(s *string) *PersonalAccessTokenUpdateOne {
	if s != nil {
		patuo.SetToken(*s)
	}
	return patuo
}

// SetPublisher sets the "publisher" edge to the Publisher entity.
func (patuo *PersonalAccessTokenUpdateOne) SetPublisher(p *Publisher) *PersonalAccessTokenUpdateOne {
	return patuo.SetPublisherID(p.ID)
}

// Mutation returns the PersonalAccessTokenMutation object of the builder.
func (patuo *PersonalAccessTokenUpdateOne) Mutation() *PersonalAccessTokenMutation {
	return patuo.mutation
}

// ClearPublisher clears the "publisher" edge to the Publisher entity.
func (patuo *PersonalAccessTokenUpdateOne) ClearPublisher() *PersonalAccessTokenUpdateOne {
	patuo.mutation.ClearPublisher()
	return patuo
}

// Where appends a list predicates to the PersonalAccessTokenUpdate builder.
func (patuo *PersonalAccessTokenUpdateOne) Where(ps ...predicate.PersonalAccessToken) *PersonalAccessTokenUpdateOne {
	patuo.mutation.Where(ps...)
	return patuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (patuo *PersonalAccessTokenUpdateOne) Select(field string, fields ...string) *PersonalAccessTokenUpdateOne {
	patuo.fields = append([]string{field}, fields...)
	return patuo
}

// Save executes the query and returns the updated PersonalAccessToken entity.
func (patuo *PersonalAccessTokenUpdateOne) Save(ctx context.Context) (*PersonalAccessToken, error) {
	patuo.defaults()
	return withHooks(ctx, patuo.sqlSave, patuo.mutation, patuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (patuo *PersonalAccessTokenUpdateOne) SaveX(ctx context.Context) *PersonalAccessToken {
	node, err := patuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (patuo *PersonalAccessTokenUpdateOne) Exec(ctx context.Context) error {
	_, err := patuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (patuo *PersonalAccessTokenUpdateOne) ExecX(ctx context.Context) {
	if err := patuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (patuo *PersonalAccessTokenUpdateOne) defaults() {
	if _, ok := patuo.mutation.UpdateTime(); !ok {
		v := personalaccesstoken.UpdateDefaultUpdateTime()
		patuo.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (patuo *PersonalAccessTokenUpdateOne) check() error {
	if patuo.mutation.PublisherCleared() && len(patuo.mutation.PublisherIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "PersonalAccessToken.publisher"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (patuo *PersonalAccessTokenUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *PersonalAccessTokenUpdateOne {
	patuo.modifiers = append(patuo.modifiers, modifiers...)
	return patuo
}

func (patuo *PersonalAccessTokenUpdateOne) sqlSave(ctx context.Context) (_node *PersonalAccessToken, err error) {
	if err := patuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(personalaccesstoken.Table, personalaccesstoken.Columns, sqlgraph.NewFieldSpec(personalaccesstoken.FieldID, field.TypeUUID))
	id, ok := patuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "PersonalAccessToken.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := patuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, personalaccesstoken.FieldID)
		for _, f := range fields {
			if !personalaccesstoken.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != personalaccesstoken.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := patuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := patuo.mutation.UpdateTime(); ok {
		_spec.SetField(personalaccesstoken.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := patuo.mutation.Name(); ok {
		_spec.SetField(personalaccesstoken.FieldName, field.TypeString, value)
	}
	if value, ok := patuo.mutation.Description(); ok {
		_spec.SetField(personalaccesstoken.FieldDescription, field.TypeString, value)
	}
	if value, ok := patuo.mutation.Token(); ok {
		_spec.SetField(personalaccesstoken.FieldToken, field.TypeString, value)
	}
	if patuo.mutation.PublisherCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   personalaccesstoken.PublisherTable,
			Columns: []string{personalaccesstoken.PublisherColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(publisher.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := patuo.mutation.PublisherIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   personalaccesstoken.PublisherTable,
			Columns: []string{personalaccesstoken.PublisherColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(publisher.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(patuo.modifiers...)
	_node = &PersonalAccessToken{config: patuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, patuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{personalaccesstoken.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	patuo.mutation.done = true
	return _node, nil
}

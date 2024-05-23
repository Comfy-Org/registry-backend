// Code generated by ent, DO NOT EDIT.

package ent

import (
	"registry-backend/ent/publisher"
	"registry-backend/ent/publisherpermission"
	"registry-backend/ent/schema"
	"registry-backend/ent/user"
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// PublisherPermissionCreate is the builder for creating a PublisherPermission entity.
type PublisherPermissionCreate struct {
	config
	mutation *PublisherPermissionMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetPermission sets the "permission" field.
func (ppc *PublisherPermissionCreate) SetPermission(spt schema.PublisherPermissionType) *PublisherPermissionCreate {
	ppc.mutation.SetPermission(spt)
	return ppc
}

// SetUserID sets the "user_id" field.
func (ppc *PublisherPermissionCreate) SetUserID(s string) *PublisherPermissionCreate {
	ppc.mutation.SetUserID(s)
	return ppc
}

// SetPublisherID sets the "publisher_id" field.
func (ppc *PublisherPermissionCreate) SetPublisherID(s string) *PublisherPermissionCreate {
	ppc.mutation.SetPublisherID(s)
	return ppc
}

// SetUser sets the "user" edge to the User entity.
func (ppc *PublisherPermissionCreate) SetUser(u *User) *PublisherPermissionCreate {
	return ppc.SetUserID(u.ID)
}

// SetPublisher sets the "publisher" edge to the Publisher entity.
func (ppc *PublisherPermissionCreate) SetPublisher(p *Publisher) *PublisherPermissionCreate {
	return ppc.SetPublisherID(p.ID)
}

// Mutation returns the PublisherPermissionMutation object of the builder.
func (ppc *PublisherPermissionCreate) Mutation() *PublisherPermissionMutation {
	return ppc.mutation
}

// Save creates the PublisherPermission in the database.
func (ppc *PublisherPermissionCreate) Save(ctx context.Context) (*PublisherPermission, error) {
	return withHooks(ctx, ppc.sqlSave, ppc.mutation, ppc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (ppc *PublisherPermissionCreate) SaveX(ctx context.Context) *PublisherPermission {
	v, err := ppc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ppc *PublisherPermissionCreate) Exec(ctx context.Context) error {
	_, err := ppc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ppc *PublisherPermissionCreate) ExecX(ctx context.Context) {
	if err := ppc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ppc *PublisherPermissionCreate) check() error {
	if _, ok := ppc.mutation.Permission(); !ok {
		return &ValidationError{Name: "permission", err: errors.New(`ent: missing required field "PublisherPermission.permission"`)}
	}
	if v, ok := ppc.mutation.Permission(); ok {
		if err := publisherpermission.PermissionValidator(v); err != nil {
			return &ValidationError{Name: "permission", err: fmt.Errorf(`ent: validator failed for field "PublisherPermission.permission": %w`, err)}
		}
	}
	if _, ok := ppc.mutation.UserID(); !ok {
		return &ValidationError{Name: "user_id", err: errors.New(`ent: missing required field "PublisherPermission.user_id"`)}
	}
	if _, ok := ppc.mutation.PublisherID(); !ok {
		return &ValidationError{Name: "publisher_id", err: errors.New(`ent: missing required field "PublisherPermission.publisher_id"`)}
	}
	if _, ok := ppc.mutation.UserID(); !ok {
		return &ValidationError{Name: "user", err: errors.New(`ent: missing required edge "PublisherPermission.user"`)}
	}
	if _, ok := ppc.mutation.PublisherID(); !ok {
		return &ValidationError{Name: "publisher", err: errors.New(`ent: missing required edge "PublisherPermission.publisher"`)}
	}
	return nil
}

func (ppc *PublisherPermissionCreate) sqlSave(ctx context.Context) (*PublisherPermission, error) {
	if err := ppc.check(); err != nil {
		return nil, err
	}
	_node, _spec := ppc.createSpec()
	if err := sqlgraph.CreateNode(ctx, ppc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	ppc.mutation.id = &_node.ID
	ppc.mutation.done = true
	return _node, nil
}

func (ppc *PublisherPermissionCreate) createSpec() (*PublisherPermission, *sqlgraph.CreateSpec) {
	var (
		_node = &PublisherPermission{config: ppc.config}
		_spec = sqlgraph.NewCreateSpec(publisherpermission.Table, sqlgraph.NewFieldSpec(publisherpermission.FieldID, field.TypeInt))
	)
	_spec.OnConflict = ppc.conflict
	if value, ok := ppc.mutation.Permission(); ok {
		_spec.SetField(publisherpermission.FieldPermission, field.TypeEnum, value)
		_node.Permission = value
	}
	if nodes := ppc.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   publisherpermission.UserTable,
			Columns: []string{publisherpermission.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.UserID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ppc.mutation.PublisherIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   publisherpermission.PublisherTable,
			Columns: []string{publisherpermission.PublisherColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(publisher.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.PublisherID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.PublisherPermission.Create().
//		SetPermission(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.PublisherPermissionUpsert) {
//			SetPermission(v+v).
//		}).
//		Exec(ctx)
func (ppc *PublisherPermissionCreate) OnConflict(opts ...sql.ConflictOption) *PublisherPermissionUpsertOne {
	ppc.conflict = opts
	return &PublisherPermissionUpsertOne{
		create: ppc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.PublisherPermission.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ppc *PublisherPermissionCreate) OnConflictColumns(columns ...string) *PublisherPermissionUpsertOne {
	ppc.conflict = append(ppc.conflict, sql.ConflictColumns(columns...))
	return &PublisherPermissionUpsertOne{
		create: ppc,
	}
}

type (
	// PublisherPermissionUpsertOne is the builder for "upsert"-ing
	//  one PublisherPermission node.
	PublisherPermissionUpsertOne struct {
		create *PublisherPermissionCreate
	}

	// PublisherPermissionUpsert is the "OnConflict" setter.
	PublisherPermissionUpsert struct {
		*sql.UpdateSet
	}
)

// SetPermission sets the "permission" field.
func (u *PublisherPermissionUpsert) SetPermission(v schema.PublisherPermissionType) *PublisherPermissionUpsert {
	u.Set(publisherpermission.FieldPermission, v)
	return u
}

// UpdatePermission sets the "permission" field to the value that was provided on create.
func (u *PublisherPermissionUpsert) UpdatePermission() *PublisherPermissionUpsert {
	u.SetExcluded(publisherpermission.FieldPermission)
	return u
}

// SetUserID sets the "user_id" field.
func (u *PublisherPermissionUpsert) SetUserID(v string) *PublisherPermissionUpsert {
	u.Set(publisherpermission.FieldUserID, v)
	return u
}

// UpdateUserID sets the "user_id" field to the value that was provided on create.
func (u *PublisherPermissionUpsert) UpdateUserID() *PublisherPermissionUpsert {
	u.SetExcluded(publisherpermission.FieldUserID)
	return u
}

// SetPublisherID sets the "publisher_id" field.
func (u *PublisherPermissionUpsert) SetPublisherID(v string) *PublisherPermissionUpsert {
	u.Set(publisherpermission.FieldPublisherID, v)
	return u
}

// UpdatePublisherID sets the "publisher_id" field to the value that was provided on create.
func (u *PublisherPermissionUpsert) UpdatePublisherID() *PublisherPermissionUpsert {
	u.SetExcluded(publisherpermission.FieldPublisherID)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.PublisherPermission.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *PublisherPermissionUpsertOne) UpdateNewValues() *PublisherPermissionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.PublisherPermission.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *PublisherPermissionUpsertOne) Ignore() *PublisherPermissionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *PublisherPermissionUpsertOne) DoNothing() *PublisherPermissionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the PublisherPermissionCreate.OnConflict
// documentation for more info.
func (u *PublisherPermissionUpsertOne) Update(set func(*PublisherPermissionUpsert)) *PublisherPermissionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&PublisherPermissionUpsert{UpdateSet: update})
	}))
	return u
}

// SetPermission sets the "permission" field.
func (u *PublisherPermissionUpsertOne) SetPermission(v schema.PublisherPermissionType) *PublisherPermissionUpsertOne {
	return u.Update(func(s *PublisherPermissionUpsert) {
		s.SetPermission(v)
	})
}

// UpdatePermission sets the "permission" field to the value that was provided on create.
func (u *PublisherPermissionUpsertOne) UpdatePermission() *PublisherPermissionUpsertOne {
	return u.Update(func(s *PublisherPermissionUpsert) {
		s.UpdatePermission()
	})
}

// SetUserID sets the "user_id" field.
func (u *PublisherPermissionUpsertOne) SetUserID(v string) *PublisherPermissionUpsertOne {
	return u.Update(func(s *PublisherPermissionUpsert) {
		s.SetUserID(v)
	})
}

// UpdateUserID sets the "user_id" field to the value that was provided on create.
func (u *PublisherPermissionUpsertOne) UpdateUserID() *PublisherPermissionUpsertOne {
	return u.Update(func(s *PublisherPermissionUpsert) {
		s.UpdateUserID()
	})
}

// SetPublisherID sets the "publisher_id" field.
func (u *PublisherPermissionUpsertOne) SetPublisherID(v string) *PublisherPermissionUpsertOne {
	return u.Update(func(s *PublisherPermissionUpsert) {
		s.SetPublisherID(v)
	})
}

// UpdatePublisherID sets the "publisher_id" field to the value that was provided on create.
func (u *PublisherPermissionUpsertOne) UpdatePublisherID() *PublisherPermissionUpsertOne {
	return u.Update(func(s *PublisherPermissionUpsert) {
		s.UpdatePublisherID()
	})
}

// Exec executes the query.
func (u *PublisherPermissionUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for PublisherPermissionCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *PublisherPermissionUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *PublisherPermissionUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *PublisherPermissionUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// PublisherPermissionCreateBulk is the builder for creating many PublisherPermission entities in bulk.
type PublisherPermissionCreateBulk struct {
	config
	err      error
	builders []*PublisherPermissionCreate
	conflict []sql.ConflictOption
}

// Save creates the PublisherPermission entities in the database.
func (ppcb *PublisherPermissionCreateBulk) Save(ctx context.Context) ([]*PublisherPermission, error) {
	if ppcb.err != nil {
		return nil, ppcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ppcb.builders))
	nodes := make([]*PublisherPermission, len(ppcb.builders))
	mutators := make([]Mutator, len(ppcb.builders))
	for i := range ppcb.builders {
		func(i int, root context.Context) {
			builder := ppcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*PublisherPermissionMutation)
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
					_, err = mutators[i+1].Mutate(root, ppcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = ppcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ppcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
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
		if _, err := mutators[0].Mutate(ctx, ppcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ppcb *PublisherPermissionCreateBulk) SaveX(ctx context.Context) []*PublisherPermission {
	v, err := ppcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ppcb *PublisherPermissionCreateBulk) Exec(ctx context.Context) error {
	_, err := ppcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ppcb *PublisherPermissionCreateBulk) ExecX(ctx context.Context) {
	if err := ppcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.PublisherPermission.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.PublisherPermissionUpsert) {
//			SetPermission(v+v).
//		}).
//		Exec(ctx)
func (ppcb *PublisherPermissionCreateBulk) OnConflict(opts ...sql.ConflictOption) *PublisherPermissionUpsertBulk {
	ppcb.conflict = opts
	return &PublisherPermissionUpsertBulk{
		create: ppcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.PublisherPermission.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ppcb *PublisherPermissionCreateBulk) OnConflictColumns(columns ...string) *PublisherPermissionUpsertBulk {
	ppcb.conflict = append(ppcb.conflict, sql.ConflictColumns(columns...))
	return &PublisherPermissionUpsertBulk{
		create: ppcb,
	}
}

// PublisherPermissionUpsertBulk is the builder for "upsert"-ing
// a bulk of PublisherPermission nodes.
type PublisherPermissionUpsertBulk struct {
	create *PublisherPermissionCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.PublisherPermission.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *PublisherPermissionUpsertBulk) UpdateNewValues() *PublisherPermissionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.PublisherPermission.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *PublisherPermissionUpsertBulk) Ignore() *PublisherPermissionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *PublisherPermissionUpsertBulk) DoNothing() *PublisherPermissionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the PublisherPermissionCreateBulk.OnConflict
// documentation for more info.
func (u *PublisherPermissionUpsertBulk) Update(set func(*PublisherPermissionUpsert)) *PublisherPermissionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&PublisherPermissionUpsert{UpdateSet: update})
	}))
	return u
}

// SetPermission sets the "permission" field.
func (u *PublisherPermissionUpsertBulk) SetPermission(v schema.PublisherPermissionType) *PublisherPermissionUpsertBulk {
	return u.Update(func(s *PublisherPermissionUpsert) {
		s.SetPermission(v)
	})
}

// UpdatePermission sets the "permission" field to the value that was provided on create.
func (u *PublisherPermissionUpsertBulk) UpdatePermission() *PublisherPermissionUpsertBulk {
	return u.Update(func(s *PublisherPermissionUpsert) {
		s.UpdatePermission()
	})
}

// SetUserID sets the "user_id" field.
func (u *PublisherPermissionUpsertBulk) SetUserID(v string) *PublisherPermissionUpsertBulk {
	return u.Update(func(s *PublisherPermissionUpsert) {
		s.SetUserID(v)
	})
}

// UpdateUserID sets the "user_id" field to the value that was provided on create.
func (u *PublisherPermissionUpsertBulk) UpdateUserID() *PublisherPermissionUpsertBulk {
	return u.Update(func(s *PublisherPermissionUpsert) {
		s.UpdateUserID()
	})
}

// SetPublisherID sets the "publisher_id" field.
func (u *PublisherPermissionUpsertBulk) SetPublisherID(v string) *PublisherPermissionUpsertBulk {
	return u.Update(func(s *PublisherPermissionUpsert) {
		s.SetPublisherID(v)
	})
}

// UpdatePublisherID sets the "publisher_id" field to the value that was provided on create.
func (u *PublisherPermissionUpsertBulk) UpdatePublisherID() *PublisherPermissionUpsertBulk {
	return u.Update(func(s *PublisherPermissionUpsert) {
		s.UpdatePublisherID()
	})
}

// Exec executes the query.
func (u *PublisherPermissionUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the PublisherPermissionCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for PublisherPermissionCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *PublisherPermissionUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

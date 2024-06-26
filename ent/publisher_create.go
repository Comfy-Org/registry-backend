// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"registry-backend/ent/node"
	"registry-backend/ent/personalaccesstoken"
	"registry-backend/ent/publisher"
	"registry-backend/ent/publisherpermission"
	"registry-backend/ent/schema"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// PublisherCreate is the builder for creating a Publisher entity.
type PublisherCreate struct {
	config
	mutation *PublisherMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreateTime sets the "create_time" field.
func (pc *PublisherCreate) SetCreateTime(t time.Time) *PublisherCreate {
	pc.mutation.SetCreateTime(t)
	return pc
}

// SetNillableCreateTime sets the "create_time" field if the given value is not nil.
func (pc *PublisherCreate) SetNillableCreateTime(t *time.Time) *PublisherCreate {
	if t != nil {
		pc.SetCreateTime(*t)
	}
	return pc
}

// SetUpdateTime sets the "update_time" field.
func (pc *PublisherCreate) SetUpdateTime(t time.Time) *PublisherCreate {
	pc.mutation.SetUpdateTime(t)
	return pc
}

// SetNillableUpdateTime sets the "update_time" field if the given value is not nil.
func (pc *PublisherCreate) SetNillableUpdateTime(t *time.Time) *PublisherCreate {
	if t != nil {
		pc.SetUpdateTime(*t)
	}
	return pc
}

// SetName sets the "name" field.
func (pc *PublisherCreate) SetName(s string) *PublisherCreate {
	pc.mutation.SetName(s)
	return pc
}

// SetDescription sets the "description" field.
func (pc *PublisherCreate) SetDescription(s string) *PublisherCreate {
	pc.mutation.SetDescription(s)
	return pc
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (pc *PublisherCreate) SetNillableDescription(s *string) *PublisherCreate {
	if s != nil {
		pc.SetDescription(*s)
	}
	return pc
}

// SetWebsite sets the "website" field.
func (pc *PublisherCreate) SetWebsite(s string) *PublisherCreate {
	pc.mutation.SetWebsite(s)
	return pc
}

// SetNillableWebsite sets the "website" field if the given value is not nil.
func (pc *PublisherCreate) SetNillableWebsite(s *string) *PublisherCreate {
	if s != nil {
		pc.SetWebsite(*s)
	}
	return pc
}

// SetSupportEmail sets the "support_email" field.
func (pc *PublisherCreate) SetSupportEmail(s string) *PublisherCreate {
	pc.mutation.SetSupportEmail(s)
	return pc
}

// SetNillableSupportEmail sets the "support_email" field if the given value is not nil.
func (pc *PublisherCreate) SetNillableSupportEmail(s *string) *PublisherCreate {
	if s != nil {
		pc.SetSupportEmail(*s)
	}
	return pc
}

// SetSourceCodeRepo sets the "source_code_repo" field.
func (pc *PublisherCreate) SetSourceCodeRepo(s string) *PublisherCreate {
	pc.mutation.SetSourceCodeRepo(s)
	return pc
}

// SetNillableSourceCodeRepo sets the "source_code_repo" field if the given value is not nil.
func (pc *PublisherCreate) SetNillableSourceCodeRepo(s *string) *PublisherCreate {
	if s != nil {
		pc.SetSourceCodeRepo(*s)
	}
	return pc
}

// SetLogoURL sets the "logo_url" field.
func (pc *PublisherCreate) SetLogoURL(s string) *PublisherCreate {
	pc.mutation.SetLogoURL(s)
	return pc
}

// SetNillableLogoURL sets the "logo_url" field if the given value is not nil.
func (pc *PublisherCreate) SetNillableLogoURL(s *string) *PublisherCreate {
	if s != nil {
		pc.SetLogoURL(*s)
	}
	return pc
}

// SetStatus sets the "status" field.
func (pc *PublisherCreate) SetStatus(sst schema.PublisherStatusType) *PublisherCreate {
	pc.mutation.SetStatus(sst)
	return pc
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (pc *PublisherCreate) SetNillableStatus(sst *schema.PublisherStatusType) *PublisherCreate {
	if sst != nil {
		pc.SetStatus(*sst)
	}
	return pc
}

// SetID sets the "id" field.
func (pc *PublisherCreate) SetID(s string) *PublisherCreate {
	pc.mutation.SetID(s)
	return pc
}

// AddPublisherPermissionIDs adds the "publisher_permissions" edge to the PublisherPermission entity by IDs.
func (pc *PublisherCreate) AddPublisherPermissionIDs(ids ...int) *PublisherCreate {
	pc.mutation.AddPublisherPermissionIDs(ids...)
	return pc
}

// AddPublisherPermissions adds the "publisher_permissions" edges to the PublisherPermission entity.
func (pc *PublisherCreate) AddPublisherPermissions(p ...*PublisherPermission) *PublisherCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pc.AddPublisherPermissionIDs(ids...)
}

// AddNodeIDs adds the "nodes" edge to the Node entity by IDs.
func (pc *PublisherCreate) AddNodeIDs(ids ...string) *PublisherCreate {
	pc.mutation.AddNodeIDs(ids...)
	return pc
}

// AddNodes adds the "nodes" edges to the Node entity.
func (pc *PublisherCreate) AddNodes(n ...*Node) *PublisherCreate {
	ids := make([]string, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return pc.AddNodeIDs(ids...)
}

// AddPersonalAccessTokenIDs adds the "personal_access_tokens" edge to the PersonalAccessToken entity by IDs.
func (pc *PublisherCreate) AddPersonalAccessTokenIDs(ids ...uuid.UUID) *PublisherCreate {
	pc.mutation.AddPersonalAccessTokenIDs(ids...)
	return pc
}

// AddPersonalAccessTokens adds the "personal_access_tokens" edges to the PersonalAccessToken entity.
func (pc *PublisherCreate) AddPersonalAccessTokens(p ...*PersonalAccessToken) *PublisherCreate {
	ids := make([]uuid.UUID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pc.AddPersonalAccessTokenIDs(ids...)
}

// Mutation returns the PublisherMutation object of the builder.
func (pc *PublisherCreate) Mutation() *PublisherMutation {
	return pc.mutation
}

// Save creates the Publisher in the database.
func (pc *PublisherCreate) Save(ctx context.Context) (*Publisher, error) {
	pc.defaults()
	return withHooks(ctx, pc.sqlSave, pc.mutation, pc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (pc *PublisherCreate) SaveX(ctx context.Context) *Publisher {
	v, err := pc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (pc *PublisherCreate) Exec(ctx context.Context) error {
	_, err := pc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pc *PublisherCreate) ExecX(ctx context.Context) {
	if err := pc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pc *PublisherCreate) defaults() {
	if _, ok := pc.mutation.CreateTime(); !ok {
		v := publisher.DefaultCreateTime()
		pc.mutation.SetCreateTime(v)
	}
	if _, ok := pc.mutation.UpdateTime(); !ok {
		v := publisher.DefaultUpdateTime()
		pc.mutation.SetUpdateTime(v)
	}
	if _, ok := pc.mutation.Status(); !ok {
		v := publisher.DefaultStatus
		pc.mutation.SetStatus(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (pc *PublisherCreate) check() error {
	if _, ok := pc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "create_time", err: errors.New(`ent: missing required field "Publisher.create_time"`)}
	}
	if _, ok := pc.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "update_time", err: errors.New(`ent: missing required field "Publisher.update_time"`)}
	}
	if _, ok := pc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Publisher.name"`)}
	}
	if _, ok := pc.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New(`ent: missing required field "Publisher.status"`)}
	}
	if v, ok := pc.mutation.Status(); ok {
		if err := publisher.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Publisher.status": %w`, err)}
		}
	}
	return nil
}

func (pc *PublisherCreate) sqlSave(ctx context.Context) (*Publisher, error) {
	if err := pc.check(); err != nil {
		return nil, err
	}
	_node, _spec := pc.createSpec()
	if err := sqlgraph.CreateNode(ctx, pc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(string); ok {
			_node.ID = id
		} else {
			return nil, fmt.Errorf("unexpected Publisher.ID type: %T", _spec.ID.Value)
		}
	}
	pc.mutation.id = &_node.ID
	pc.mutation.done = true
	return _node, nil
}

func (pc *PublisherCreate) createSpec() (*Publisher, *sqlgraph.CreateSpec) {
	var (
		_node = &Publisher{config: pc.config}
		_spec = sqlgraph.NewCreateSpec(publisher.Table, sqlgraph.NewFieldSpec(publisher.FieldID, field.TypeString))
	)
	_spec.OnConflict = pc.conflict
	if id, ok := pc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := pc.mutation.CreateTime(); ok {
		_spec.SetField(publisher.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = value
	}
	if value, ok := pc.mutation.UpdateTime(); ok {
		_spec.SetField(publisher.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = value
	}
	if value, ok := pc.mutation.Name(); ok {
		_spec.SetField(publisher.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := pc.mutation.Description(); ok {
		_spec.SetField(publisher.FieldDescription, field.TypeString, value)
		_node.Description = value
	}
	if value, ok := pc.mutation.Website(); ok {
		_spec.SetField(publisher.FieldWebsite, field.TypeString, value)
		_node.Website = value
	}
	if value, ok := pc.mutation.SupportEmail(); ok {
		_spec.SetField(publisher.FieldSupportEmail, field.TypeString, value)
		_node.SupportEmail = value
	}
	if value, ok := pc.mutation.SourceCodeRepo(); ok {
		_spec.SetField(publisher.FieldSourceCodeRepo, field.TypeString, value)
		_node.SourceCodeRepo = value
	}
	if value, ok := pc.mutation.LogoURL(); ok {
		_spec.SetField(publisher.FieldLogoURL, field.TypeString, value)
		_node.LogoURL = value
	}
	if value, ok := pc.mutation.Status(); ok {
		_spec.SetField(publisher.FieldStatus, field.TypeEnum, value)
		_node.Status = value
	}
	if nodes := pc.mutation.PublisherPermissionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   publisher.PublisherPermissionsTable,
			Columns: []string{publisher.PublisherPermissionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(publisherpermission.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.NodesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   publisher.NodesTable,
			Columns: []string{publisher.NodesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(node.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := pc.mutation.PersonalAccessTokensIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   publisher.PersonalAccessTokensTable,
			Columns: []string{publisher.PersonalAccessTokensColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(personalaccesstoken.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Publisher.Create().
//		SetCreateTime(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.PublisherUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (pc *PublisherCreate) OnConflict(opts ...sql.ConflictOption) *PublisherUpsertOne {
	pc.conflict = opts
	return &PublisherUpsertOne{
		create: pc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Publisher.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (pc *PublisherCreate) OnConflictColumns(columns ...string) *PublisherUpsertOne {
	pc.conflict = append(pc.conflict, sql.ConflictColumns(columns...))
	return &PublisherUpsertOne{
		create: pc,
	}
}

type (
	// PublisherUpsertOne is the builder for "upsert"-ing
	//  one Publisher node.
	PublisherUpsertOne struct {
		create *PublisherCreate
	}

	// PublisherUpsert is the "OnConflict" setter.
	PublisherUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdateTime sets the "update_time" field.
func (u *PublisherUpsert) SetUpdateTime(v time.Time) *PublisherUpsert {
	u.Set(publisher.FieldUpdateTime, v)
	return u
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *PublisherUpsert) UpdateUpdateTime() *PublisherUpsert {
	u.SetExcluded(publisher.FieldUpdateTime)
	return u
}

// SetName sets the "name" field.
func (u *PublisherUpsert) SetName(v string) *PublisherUpsert {
	u.Set(publisher.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *PublisherUpsert) UpdateName() *PublisherUpsert {
	u.SetExcluded(publisher.FieldName)
	return u
}

// SetDescription sets the "description" field.
func (u *PublisherUpsert) SetDescription(v string) *PublisherUpsert {
	u.Set(publisher.FieldDescription, v)
	return u
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *PublisherUpsert) UpdateDescription() *PublisherUpsert {
	u.SetExcluded(publisher.FieldDescription)
	return u
}

// ClearDescription clears the value of the "description" field.
func (u *PublisherUpsert) ClearDescription() *PublisherUpsert {
	u.SetNull(publisher.FieldDescription)
	return u
}

// SetWebsite sets the "website" field.
func (u *PublisherUpsert) SetWebsite(v string) *PublisherUpsert {
	u.Set(publisher.FieldWebsite, v)
	return u
}

// UpdateWebsite sets the "website" field to the value that was provided on create.
func (u *PublisherUpsert) UpdateWebsite() *PublisherUpsert {
	u.SetExcluded(publisher.FieldWebsite)
	return u
}

// ClearWebsite clears the value of the "website" field.
func (u *PublisherUpsert) ClearWebsite() *PublisherUpsert {
	u.SetNull(publisher.FieldWebsite)
	return u
}

// SetSupportEmail sets the "support_email" field.
func (u *PublisherUpsert) SetSupportEmail(v string) *PublisherUpsert {
	u.Set(publisher.FieldSupportEmail, v)
	return u
}

// UpdateSupportEmail sets the "support_email" field to the value that was provided on create.
func (u *PublisherUpsert) UpdateSupportEmail() *PublisherUpsert {
	u.SetExcluded(publisher.FieldSupportEmail)
	return u
}

// ClearSupportEmail clears the value of the "support_email" field.
func (u *PublisherUpsert) ClearSupportEmail() *PublisherUpsert {
	u.SetNull(publisher.FieldSupportEmail)
	return u
}

// SetSourceCodeRepo sets the "source_code_repo" field.
func (u *PublisherUpsert) SetSourceCodeRepo(v string) *PublisherUpsert {
	u.Set(publisher.FieldSourceCodeRepo, v)
	return u
}

// UpdateSourceCodeRepo sets the "source_code_repo" field to the value that was provided on create.
func (u *PublisherUpsert) UpdateSourceCodeRepo() *PublisherUpsert {
	u.SetExcluded(publisher.FieldSourceCodeRepo)
	return u
}

// ClearSourceCodeRepo clears the value of the "source_code_repo" field.
func (u *PublisherUpsert) ClearSourceCodeRepo() *PublisherUpsert {
	u.SetNull(publisher.FieldSourceCodeRepo)
	return u
}

// SetLogoURL sets the "logo_url" field.
func (u *PublisherUpsert) SetLogoURL(v string) *PublisherUpsert {
	u.Set(publisher.FieldLogoURL, v)
	return u
}

// UpdateLogoURL sets the "logo_url" field to the value that was provided on create.
func (u *PublisherUpsert) UpdateLogoURL() *PublisherUpsert {
	u.SetExcluded(publisher.FieldLogoURL)
	return u
}

// ClearLogoURL clears the value of the "logo_url" field.
func (u *PublisherUpsert) ClearLogoURL() *PublisherUpsert {
	u.SetNull(publisher.FieldLogoURL)
	return u
}

// SetStatus sets the "status" field.
func (u *PublisherUpsert) SetStatus(v schema.PublisherStatusType) *PublisherUpsert {
	u.Set(publisher.FieldStatus, v)
	return u
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *PublisherUpsert) UpdateStatus() *PublisherUpsert {
	u.SetExcluded(publisher.FieldStatus)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Publisher.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(publisher.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *PublisherUpsertOne) UpdateNewValues() *PublisherUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(publisher.FieldID)
		}
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(publisher.FieldCreateTime)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Publisher.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *PublisherUpsertOne) Ignore() *PublisherUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *PublisherUpsertOne) DoNothing() *PublisherUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the PublisherCreate.OnConflict
// documentation for more info.
func (u *PublisherUpsertOne) Update(set func(*PublisherUpsert)) *PublisherUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&PublisherUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "update_time" field.
func (u *PublisherUpsertOne) SetUpdateTime(v time.Time) *PublisherUpsertOne {
	return u.Update(func(s *PublisherUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *PublisherUpsertOne) UpdateUpdateTime() *PublisherUpsertOne {
	return u.Update(func(s *PublisherUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetName sets the "name" field.
func (u *PublisherUpsertOne) SetName(v string) *PublisherUpsertOne {
	return u.Update(func(s *PublisherUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *PublisherUpsertOne) UpdateName() *PublisherUpsertOne {
	return u.Update(func(s *PublisherUpsert) {
		s.UpdateName()
	})
}

// SetDescription sets the "description" field.
func (u *PublisherUpsertOne) SetDescription(v string) *PublisherUpsertOne {
	return u.Update(func(s *PublisherUpsert) {
		s.SetDescription(v)
	})
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *PublisherUpsertOne) UpdateDescription() *PublisherUpsertOne {
	return u.Update(func(s *PublisherUpsert) {
		s.UpdateDescription()
	})
}

// ClearDescription clears the value of the "description" field.
func (u *PublisherUpsertOne) ClearDescription() *PublisherUpsertOne {
	return u.Update(func(s *PublisherUpsert) {
		s.ClearDescription()
	})
}

// SetWebsite sets the "website" field.
func (u *PublisherUpsertOne) SetWebsite(v string) *PublisherUpsertOne {
	return u.Update(func(s *PublisherUpsert) {
		s.SetWebsite(v)
	})
}

// UpdateWebsite sets the "website" field to the value that was provided on create.
func (u *PublisherUpsertOne) UpdateWebsite() *PublisherUpsertOne {
	return u.Update(func(s *PublisherUpsert) {
		s.UpdateWebsite()
	})
}

// ClearWebsite clears the value of the "website" field.
func (u *PublisherUpsertOne) ClearWebsite() *PublisherUpsertOne {
	return u.Update(func(s *PublisherUpsert) {
		s.ClearWebsite()
	})
}

// SetSupportEmail sets the "support_email" field.
func (u *PublisherUpsertOne) SetSupportEmail(v string) *PublisherUpsertOne {
	return u.Update(func(s *PublisherUpsert) {
		s.SetSupportEmail(v)
	})
}

// UpdateSupportEmail sets the "support_email" field to the value that was provided on create.
func (u *PublisherUpsertOne) UpdateSupportEmail() *PublisherUpsertOne {
	return u.Update(func(s *PublisherUpsert) {
		s.UpdateSupportEmail()
	})
}

// ClearSupportEmail clears the value of the "support_email" field.
func (u *PublisherUpsertOne) ClearSupportEmail() *PublisherUpsertOne {
	return u.Update(func(s *PublisherUpsert) {
		s.ClearSupportEmail()
	})
}

// SetSourceCodeRepo sets the "source_code_repo" field.
func (u *PublisherUpsertOne) SetSourceCodeRepo(v string) *PublisherUpsertOne {
	return u.Update(func(s *PublisherUpsert) {
		s.SetSourceCodeRepo(v)
	})
}

// UpdateSourceCodeRepo sets the "source_code_repo" field to the value that was provided on create.
func (u *PublisherUpsertOne) UpdateSourceCodeRepo() *PublisherUpsertOne {
	return u.Update(func(s *PublisherUpsert) {
		s.UpdateSourceCodeRepo()
	})
}

// ClearSourceCodeRepo clears the value of the "source_code_repo" field.
func (u *PublisherUpsertOne) ClearSourceCodeRepo() *PublisherUpsertOne {
	return u.Update(func(s *PublisherUpsert) {
		s.ClearSourceCodeRepo()
	})
}

// SetLogoURL sets the "logo_url" field.
func (u *PublisherUpsertOne) SetLogoURL(v string) *PublisherUpsertOne {
	return u.Update(func(s *PublisherUpsert) {
		s.SetLogoURL(v)
	})
}

// UpdateLogoURL sets the "logo_url" field to the value that was provided on create.
func (u *PublisherUpsertOne) UpdateLogoURL() *PublisherUpsertOne {
	return u.Update(func(s *PublisherUpsert) {
		s.UpdateLogoURL()
	})
}

// ClearLogoURL clears the value of the "logo_url" field.
func (u *PublisherUpsertOne) ClearLogoURL() *PublisherUpsertOne {
	return u.Update(func(s *PublisherUpsert) {
		s.ClearLogoURL()
	})
}

// SetStatus sets the "status" field.
func (u *PublisherUpsertOne) SetStatus(v schema.PublisherStatusType) *PublisherUpsertOne {
	return u.Update(func(s *PublisherUpsert) {
		s.SetStatus(v)
	})
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *PublisherUpsertOne) UpdateStatus() *PublisherUpsertOne {
	return u.Update(func(s *PublisherUpsert) {
		s.UpdateStatus()
	})
}

// Exec executes the query.
func (u *PublisherUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for PublisherCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *PublisherUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *PublisherUpsertOne) ID(ctx context.Context) (id string, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("ent: PublisherUpsertOne.ID is not supported by MySQL driver. Use PublisherUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *PublisherUpsertOne) IDX(ctx context.Context) string {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// PublisherCreateBulk is the builder for creating many Publisher entities in bulk.
type PublisherCreateBulk struct {
	config
	err      error
	builders []*PublisherCreate
	conflict []sql.ConflictOption
}

// Save creates the Publisher entities in the database.
func (pcb *PublisherCreateBulk) Save(ctx context.Context) ([]*Publisher, error) {
	if pcb.err != nil {
		return nil, pcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(pcb.builders))
	nodes := make([]*Publisher, len(pcb.builders))
	mutators := make([]Mutator, len(pcb.builders))
	for i := range pcb.builders {
		func(i int, root context.Context) {
			builder := pcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*PublisherMutation)
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
					_, err = mutators[i+1].Mutate(root, pcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = pcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, pcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, pcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (pcb *PublisherCreateBulk) SaveX(ctx context.Context) []*Publisher {
	v, err := pcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (pcb *PublisherCreateBulk) Exec(ctx context.Context) error {
	_, err := pcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pcb *PublisherCreateBulk) ExecX(ctx context.Context) {
	if err := pcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Publisher.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.PublisherUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (pcb *PublisherCreateBulk) OnConflict(opts ...sql.ConflictOption) *PublisherUpsertBulk {
	pcb.conflict = opts
	return &PublisherUpsertBulk{
		create: pcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Publisher.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (pcb *PublisherCreateBulk) OnConflictColumns(columns ...string) *PublisherUpsertBulk {
	pcb.conflict = append(pcb.conflict, sql.ConflictColumns(columns...))
	return &PublisherUpsertBulk{
		create: pcb,
	}
}

// PublisherUpsertBulk is the builder for "upsert"-ing
// a bulk of Publisher nodes.
type PublisherUpsertBulk struct {
	create *PublisherCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Publisher.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(publisher.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *PublisherUpsertBulk) UpdateNewValues() *PublisherUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(publisher.FieldID)
			}
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(publisher.FieldCreateTime)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Publisher.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *PublisherUpsertBulk) Ignore() *PublisherUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *PublisherUpsertBulk) DoNothing() *PublisherUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the PublisherCreateBulk.OnConflict
// documentation for more info.
func (u *PublisherUpsertBulk) Update(set func(*PublisherUpsert)) *PublisherUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&PublisherUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "update_time" field.
func (u *PublisherUpsertBulk) SetUpdateTime(v time.Time) *PublisherUpsertBulk {
	return u.Update(func(s *PublisherUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *PublisherUpsertBulk) UpdateUpdateTime() *PublisherUpsertBulk {
	return u.Update(func(s *PublisherUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetName sets the "name" field.
func (u *PublisherUpsertBulk) SetName(v string) *PublisherUpsertBulk {
	return u.Update(func(s *PublisherUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *PublisherUpsertBulk) UpdateName() *PublisherUpsertBulk {
	return u.Update(func(s *PublisherUpsert) {
		s.UpdateName()
	})
}

// SetDescription sets the "description" field.
func (u *PublisherUpsertBulk) SetDescription(v string) *PublisherUpsertBulk {
	return u.Update(func(s *PublisherUpsert) {
		s.SetDescription(v)
	})
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *PublisherUpsertBulk) UpdateDescription() *PublisherUpsertBulk {
	return u.Update(func(s *PublisherUpsert) {
		s.UpdateDescription()
	})
}

// ClearDescription clears the value of the "description" field.
func (u *PublisherUpsertBulk) ClearDescription() *PublisherUpsertBulk {
	return u.Update(func(s *PublisherUpsert) {
		s.ClearDescription()
	})
}

// SetWebsite sets the "website" field.
func (u *PublisherUpsertBulk) SetWebsite(v string) *PublisherUpsertBulk {
	return u.Update(func(s *PublisherUpsert) {
		s.SetWebsite(v)
	})
}

// UpdateWebsite sets the "website" field to the value that was provided on create.
func (u *PublisherUpsertBulk) UpdateWebsite() *PublisherUpsertBulk {
	return u.Update(func(s *PublisherUpsert) {
		s.UpdateWebsite()
	})
}

// ClearWebsite clears the value of the "website" field.
func (u *PublisherUpsertBulk) ClearWebsite() *PublisherUpsertBulk {
	return u.Update(func(s *PublisherUpsert) {
		s.ClearWebsite()
	})
}

// SetSupportEmail sets the "support_email" field.
func (u *PublisherUpsertBulk) SetSupportEmail(v string) *PublisherUpsertBulk {
	return u.Update(func(s *PublisherUpsert) {
		s.SetSupportEmail(v)
	})
}

// UpdateSupportEmail sets the "support_email" field to the value that was provided on create.
func (u *PublisherUpsertBulk) UpdateSupportEmail() *PublisherUpsertBulk {
	return u.Update(func(s *PublisherUpsert) {
		s.UpdateSupportEmail()
	})
}

// ClearSupportEmail clears the value of the "support_email" field.
func (u *PublisherUpsertBulk) ClearSupportEmail() *PublisherUpsertBulk {
	return u.Update(func(s *PublisherUpsert) {
		s.ClearSupportEmail()
	})
}

// SetSourceCodeRepo sets the "source_code_repo" field.
func (u *PublisherUpsertBulk) SetSourceCodeRepo(v string) *PublisherUpsertBulk {
	return u.Update(func(s *PublisherUpsert) {
		s.SetSourceCodeRepo(v)
	})
}

// UpdateSourceCodeRepo sets the "source_code_repo" field to the value that was provided on create.
func (u *PublisherUpsertBulk) UpdateSourceCodeRepo() *PublisherUpsertBulk {
	return u.Update(func(s *PublisherUpsert) {
		s.UpdateSourceCodeRepo()
	})
}

// ClearSourceCodeRepo clears the value of the "source_code_repo" field.
func (u *PublisherUpsertBulk) ClearSourceCodeRepo() *PublisherUpsertBulk {
	return u.Update(func(s *PublisherUpsert) {
		s.ClearSourceCodeRepo()
	})
}

// SetLogoURL sets the "logo_url" field.
func (u *PublisherUpsertBulk) SetLogoURL(v string) *PublisherUpsertBulk {
	return u.Update(func(s *PublisherUpsert) {
		s.SetLogoURL(v)
	})
}

// UpdateLogoURL sets the "logo_url" field to the value that was provided on create.
func (u *PublisherUpsertBulk) UpdateLogoURL() *PublisherUpsertBulk {
	return u.Update(func(s *PublisherUpsert) {
		s.UpdateLogoURL()
	})
}

// ClearLogoURL clears the value of the "logo_url" field.
func (u *PublisherUpsertBulk) ClearLogoURL() *PublisherUpsertBulk {
	return u.Update(func(s *PublisherUpsert) {
		s.ClearLogoURL()
	})
}

// SetStatus sets the "status" field.
func (u *PublisherUpsertBulk) SetStatus(v schema.PublisherStatusType) *PublisherUpsertBulk {
	return u.Update(func(s *PublisherUpsert) {
		s.SetStatus(v)
	})
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *PublisherUpsertBulk) UpdateStatus() *PublisherUpsertBulk {
	return u.Update(func(s *PublisherUpsert) {
		s.UpdateStatus()
	})
}

// Exec executes the query.
func (u *PublisherUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the PublisherCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for PublisherCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *PublisherUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

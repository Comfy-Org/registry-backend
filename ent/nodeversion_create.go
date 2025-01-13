// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"registry-backend/ent/comfynode"
	"registry-backend/ent/node"
	"registry-backend/ent/nodeversion"
	"registry-backend/ent/schema"
	"registry-backend/ent/storagefile"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// NodeVersionCreate is the builder for creating a NodeVersion entity.
type NodeVersionCreate struct {
	config
	mutation *NodeVersionMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreateTime sets the "create_time" field.
func (nvc *NodeVersionCreate) SetCreateTime(t time.Time) *NodeVersionCreate {
	nvc.mutation.SetCreateTime(t)
	return nvc
}

// SetNillableCreateTime sets the "create_time" field if the given value is not nil.
func (nvc *NodeVersionCreate) SetNillableCreateTime(t *time.Time) *NodeVersionCreate {
	if t != nil {
		nvc.SetCreateTime(*t)
	}
	return nvc
}

// SetUpdateTime sets the "update_time" field.
func (nvc *NodeVersionCreate) SetUpdateTime(t time.Time) *NodeVersionCreate {
	nvc.mutation.SetUpdateTime(t)
	return nvc
}

// SetNillableUpdateTime sets the "update_time" field if the given value is not nil.
func (nvc *NodeVersionCreate) SetNillableUpdateTime(t *time.Time) *NodeVersionCreate {
	if t != nil {
		nvc.SetUpdateTime(*t)
	}
	return nvc
}

// SetNodeID sets the "node_id" field.
func (nvc *NodeVersionCreate) SetNodeID(s string) *NodeVersionCreate {
	nvc.mutation.SetNodeID(s)
	return nvc
}

// SetVersion sets the "version" field.
func (nvc *NodeVersionCreate) SetVersion(s string) *NodeVersionCreate {
	nvc.mutation.SetVersion(s)
	return nvc
}

// SetChangelog sets the "changelog" field.
func (nvc *NodeVersionCreate) SetChangelog(s string) *NodeVersionCreate {
	nvc.mutation.SetChangelog(s)
	return nvc
}

// SetNillableChangelog sets the "changelog" field if the given value is not nil.
func (nvc *NodeVersionCreate) SetNillableChangelog(s *string) *NodeVersionCreate {
	if s != nil {
		nvc.SetChangelog(*s)
	}
	return nvc
}

// SetPipDependencies sets the "pip_dependencies" field.
func (nvc *NodeVersionCreate) SetPipDependencies(s []string) *NodeVersionCreate {
	nvc.mutation.SetPipDependencies(s)
	return nvc
}

// SetDeprecated sets the "deprecated" field.
func (nvc *NodeVersionCreate) SetDeprecated(b bool) *NodeVersionCreate {
	nvc.mutation.SetDeprecated(b)
	return nvc
}

// SetNillableDeprecated sets the "deprecated" field if the given value is not nil.
func (nvc *NodeVersionCreate) SetNillableDeprecated(b *bool) *NodeVersionCreate {
	if b != nil {
		nvc.SetDeprecated(*b)
	}
	return nvc
}

// SetStatus sets the "status" field.
func (nvc *NodeVersionCreate) SetStatus(svs schema.NodeVersionStatus) *NodeVersionCreate {
	nvc.mutation.SetStatus(svs)
	return nvc
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (nvc *NodeVersionCreate) SetNillableStatus(svs *schema.NodeVersionStatus) *NodeVersionCreate {
	if svs != nil {
		nvc.SetStatus(*svs)
	}
	return nvc
}

// SetStatusReason sets the "status_reason" field.
func (nvc *NodeVersionCreate) SetStatusReason(s string) *NodeVersionCreate {
	nvc.mutation.SetStatusReason(s)
	return nvc
}

// SetNillableStatusReason sets the "status_reason" field if the given value is not nil.
func (nvc *NodeVersionCreate) SetNillableStatusReason(s *string) *NodeVersionCreate {
	if s != nil {
		nvc.SetStatusReason(*s)
	}
	return nvc
}

// SetComfyNodeExtractStatus sets the "comfy_node_extract_status" field.
func (nvc *NodeVersionCreate) SetComfyNodeExtractStatus(snes schema.ComfyNodeExtractStatus) *NodeVersionCreate {
	nvc.mutation.SetComfyNodeExtractStatus(snes)
	return nvc
}

// SetNillableComfyNodeExtractStatus sets the "comfy_node_extract_status" field if the given value is not nil.
func (nvc *NodeVersionCreate) SetNillableComfyNodeExtractStatus(snes *schema.ComfyNodeExtractStatus) *NodeVersionCreate {
	if snes != nil {
		nvc.SetComfyNodeExtractStatus(*snes)
	}
	return nvc
}

// SetComfyNodeCloudBuildInfo sets the "comfy_node_cloud_build_info" field.
func (nvc *NodeVersionCreate) SetComfyNodeCloudBuildInfo(sncbi schema.ComfyNodeCloudBuildInfo) *NodeVersionCreate {
	nvc.mutation.SetComfyNodeCloudBuildInfo(sncbi)
	return nvc
}

// SetNillableComfyNodeCloudBuildInfo sets the "comfy_node_cloud_build_info" field if the given value is not nil.
func (nvc *NodeVersionCreate) SetNillableComfyNodeCloudBuildInfo(sncbi *schema.ComfyNodeCloudBuildInfo) *NodeVersionCreate {
	if sncbi != nil {
		nvc.SetComfyNodeCloudBuildInfo(*sncbi)
	}
	return nvc
}

// SetID sets the "id" field.
func (nvc *NodeVersionCreate) SetID(u uuid.UUID) *NodeVersionCreate {
	nvc.mutation.SetID(u)
	return nvc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (nvc *NodeVersionCreate) SetNillableID(u *uuid.UUID) *NodeVersionCreate {
	if u != nil {
		nvc.SetID(*u)
	}
	return nvc
}

// SetNode sets the "node" edge to the Node entity.
func (nvc *NodeVersionCreate) SetNode(n *Node) *NodeVersionCreate {
	return nvc.SetNodeID(n.ID)
}

// SetStorageFileID sets the "storage_file" edge to the StorageFile entity by ID.
func (nvc *NodeVersionCreate) SetStorageFileID(id uuid.UUID) *NodeVersionCreate {
	nvc.mutation.SetStorageFileID(id)
	return nvc
}

// SetNillableStorageFileID sets the "storage_file" edge to the StorageFile entity by ID if the given value is not nil.
func (nvc *NodeVersionCreate) SetNillableStorageFileID(id *uuid.UUID) *NodeVersionCreate {
	if id != nil {
		nvc = nvc.SetStorageFileID(*id)
	}
	return nvc
}

// SetStorageFile sets the "storage_file" edge to the StorageFile entity.
func (nvc *NodeVersionCreate) SetStorageFile(s *StorageFile) *NodeVersionCreate {
	return nvc.SetStorageFileID(s.ID)
}

// AddComfyNodeIDs adds the "comfy_nodes" edge to the ComfyNode entity by IDs.
func (nvc *NodeVersionCreate) AddComfyNodeIDs(ids ...uuid.UUID) *NodeVersionCreate {
	nvc.mutation.AddComfyNodeIDs(ids...)
	return nvc
}

// AddComfyNodes adds the "comfy_nodes" edges to the ComfyNode entity.
func (nvc *NodeVersionCreate) AddComfyNodes(c ...*ComfyNode) *NodeVersionCreate {
	ids := make([]uuid.UUID, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return nvc.AddComfyNodeIDs(ids...)
}

// Mutation returns the NodeVersionMutation object of the builder.
func (nvc *NodeVersionCreate) Mutation() *NodeVersionMutation {
	return nvc.mutation
}

// Save creates the NodeVersion in the database.
func (nvc *NodeVersionCreate) Save(ctx context.Context) (*NodeVersion, error) {
	nvc.defaults()
	return withHooks(ctx, nvc.sqlSave, nvc.mutation, nvc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (nvc *NodeVersionCreate) SaveX(ctx context.Context) *NodeVersion {
	v, err := nvc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (nvc *NodeVersionCreate) Exec(ctx context.Context) error {
	_, err := nvc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (nvc *NodeVersionCreate) ExecX(ctx context.Context) {
	if err := nvc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (nvc *NodeVersionCreate) defaults() {
	if _, ok := nvc.mutation.CreateTime(); !ok {
		v := nodeversion.DefaultCreateTime()
		nvc.mutation.SetCreateTime(v)
	}
	if _, ok := nvc.mutation.UpdateTime(); !ok {
		v := nodeversion.DefaultUpdateTime()
		nvc.mutation.SetUpdateTime(v)
	}
	if _, ok := nvc.mutation.Deprecated(); !ok {
		v := nodeversion.DefaultDeprecated
		nvc.mutation.SetDeprecated(v)
	}
	if _, ok := nvc.mutation.Status(); !ok {
		v := nodeversion.DefaultStatus
		nvc.mutation.SetStatus(v)
	}
	if _, ok := nvc.mutation.StatusReason(); !ok {
		v := nodeversion.DefaultStatusReason
		nvc.mutation.SetStatusReason(v)
	}
	if _, ok := nvc.mutation.ComfyNodeExtractStatus(); !ok {
		v := nodeversion.DefaultComfyNodeExtractStatus
		nvc.mutation.SetComfyNodeExtractStatus(v)
	}
	if _, ok := nvc.mutation.ID(); !ok {
		v := nodeversion.DefaultID()
		nvc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (nvc *NodeVersionCreate) check() error {
	if _, ok := nvc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "create_time", err: errors.New(`ent: missing required field "NodeVersion.create_time"`)}
	}
	if _, ok := nvc.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "update_time", err: errors.New(`ent: missing required field "NodeVersion.update_time"`)}
	}
	if _, ok := nvc.mutation.NodeID(); !ok {
		return &ValidationError{Name: "node_id", err: errors.New(`ent: missing required field "NodeVersion.node_id"`)}
	}
	if _, ok := nvc.mutation.Version(); !ok {
		return &ValidationError{Name: "version", err: errors.New(`ent: missing required field "NodeVersion.version"`)}
	}
	if _, ok := nvc.mutation.PipDependencies(); !ok {
		return &ValidationError{Name: "pip_dependencies", err: errors.New(`ent: missing required field "NodeVersion.pip_dependencies"`)}
	}
	if _, ok := nvc.mutation.Deprecated(); !ok {
		return &ValidationError{Name: "deprecated", err: errors.New(`ent: missing required field "NodeVersion.deprecated"`)}
	}
	if _, ok := nvc.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New(`ent: missing required field "NodeVersion.status"`)}
	}
	if v, ok := nvc.mutation.Status(); ok {
		if err := nodeversion.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "NodeVersion.status": %w`, err)}
		}
	}
	if _, ok := nvc.mutation.StatusReason(); !ok {
		return &ValidationError{Name: "status_reason", err: errors.New(`ent: missing required field "NodeVersion.status_reason"`)}
	}
	if _, ok := nvc.mutation.ComfyNodeExtractStatus(); !ok {
		return &ValidationError{Name: "comfy_node_extract_status", err: errors.New(`ent: missing required field "NodeVersion.comfy_node_extract_status"`)}
	}
	if len(nvc.mutation.NodeIDs()) == 0 {
		return &ValidationError{Name: "node", err: errors.New(`ent: missing required edge "NodeVersion.node"`)}
	}
	return nil
}

func (nvc *NodeVersionCreate) sqlSave(ctx context.Context) (*NodeVersion, error) {
	if err := nvc.check(); err != nil {
		return nil, err
	}
	_node, _spec := nvc.createSpec()
	if err := sqlgraph.CreateNode(ctx, nvc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	nvc.mutation.id = &_node.ID
	nvc.mutation.done = true
	return _node, nil
}

func (nvc *NodeVersionCreate) createSpec() (*NodeVersion, *sqlgraph.CreateSpec) {
	var (
		_node = &NodeVersion{config: nvc.config}
		_spec = sqlgraph.NewCreateSpec(nodeversion.Table, sqlgraph.NewFieldSpec(nodeversion.FieldID, field.TypeUUID))
	)
	_spec.OnConflict = nvc.conflict
	if id, ok := nvc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := nvc.mutation.CreateTime(); ok {
		_spec.SetField(nodeversion.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = value
	}
	if value, ok := nvc.mutation.UpdateTime(); ok {
		_spec.SetField(nodeversion.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = value
	}
	if value, ok := nvc.mutation.Version(); ok {
		_spec.SetField(nodeversion.FieldVersion, field.TypeString, value)
		_node.Version = value
	}
	if value, ok := nvc.mutation.Changelog(); ok {
		_spec.SetField(nodeversion.FieldChangelog, field.TypeString, value)
		_node.Changelog = value
	}
	if value, ok := nvc.mutation.PipDependencies(); ok {
		_spec.SetField(nodeversion.FieldPipDependencies, field.TypeJSON, value)
		_node.PipDependencies = value
	}
	if value, ok := nvc.mutation.Deprecated(); ok {
		_spec.SetField(nodeversion.FieldDeprecated, field.TypeBool, value)
		_node.Deprecated = value
	}
	if value, ok := nvc.mutation.Status(); ok {
		_spec.SetField(nodeversion.FieldStatus, field.TypeEnum, value)
		_node.Status = value
	}
	if value, ok := nvc.mutation.StatusReason(); ok {
		_spec.SetField(nodeversion.FieldStatusReason, field.TypeString, value)
		_node.StatusReason = value
	}
	if value, ok := nvc.mutation.ComfyNodeExtractStatus(); ok {
		_spec.SetField(nodeversion.FieldComfyNodeExtractStatus, field.TypeString, value)
		_node.ComfyNodeExtractStatus = value
	}
	if value, ok := nvc.mutation.ComfyNodeCloudBuildInfo(); ok {
		_spec.SetField(nodeversion.FieldComfyNodeCloudBuildInfo, field.TypeJSON, value)
		_node.ComfyNodeCloudBuildInfo = value
	}
	if nodes := nvc.mutation.NodeIDs(); len(nodes) > 0 {
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
		_node.NodeID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := nvc.mutation.StorageFileIDs(); len(nodes) > 0 {
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
		_node.node_version_storage_file = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := nvc.mutation.ComfyNodesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   nodeversion.ComfyNodesTable,
			Columns: []string{nodeversion.ComfyNodesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(comfynode.FieldID, field.TypeUUID),
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
//	client.NodeVersion.Create().
//		SetCreateTime(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.NodeVersionUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (nvc *NodeVersionCreate) OnConflict(opts ...sql.ConflictOption) *NodeVersionUpsertOne {
	nvc.conflict = opts
	return &NodeVersionUpsertOne{
		create: nvc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.NodeVersion.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (nvc *NodeVersionCreate) OnConflictColumns(columns ...string) *NodeVersionUpsertOne {
	nvc.conflict = append(nvc.conflict, sql.ConflictColumns(columns...))
	return &NodeVersionUpsertOne{
		create: nvc,
	}
}

type (
	// NodeVersionUpsertOne is the builder for "upsert"-ing
	//  one NodeVersion node.
	NodeVersionUpsertOne struct {
		create *NodeVersionCreate
	}

	// NodeVersionUpsert is the "OnConflict" setter.
	NodeVersionUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdateTime sets the "update_time" field.
func (u *NodeVersionUpsert) SetUpdateTime(v time.Time) *NodeVersionUpsert {
	u.Set(nodeversion.FieldUpdateTime, v)
	return u
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *NodeVersionUpsert) UpdateUpdateTime() *NodeVersionUpsert {
	u.SetExcluded(nodeversion.FieldUpdateTime)
	return u
}

// SetNodeID sets the "node_id" field.
func (u *NodeVersionUpsert) SetNodeID(v string) *NodeVersionUpsert {
	u.Set(nodeversion.FieldNodeID, v)
	return u
}

// UpdateNodeID sets the "node_id" field to the value that was provided on create.
func (u *NodeVersionUpsert) UpdateNodeID() *NodeVersionUpsert {
	u.SetExcluded(nodeversion.FieldNodeID)
	return u
}

// SetVersion sets the "version" field.
func (u *NodeVersionUpsert) SetVersion(v string) *NodeVersionUpsert {
	u.Set(nodeversion.FieldVersion, v)
	return u
}

// UpdateVersion sets the "version" field to the value that was provided on create.
func (u *NodeVersionUpsert) UpdateVersion() *NodeVersionUpsert {
	u.SetExcluded(nodeversion.FieldVersion)
	return u
}

// SetChangelog sets the "changelog" field.
func (u *NodeVersionUpsert) SetChangelog(v string) *NodeVersionUpsert {
	u.Set(nodeversion.FieldChangelog, v)
	return u
}

// UpdateChangelog sets the "changelog" field to the value that was provided on create.
func (u *NodeVersionUpsert) UpdateChangelog() *NodeVersionUpsert {
	u.SetExcluded(nodeversion.FieldChangelog)
	return u
}

// ClearChangelog clears the value of the "changelog" field.
func (u *NodeVersionUpsert) ClearChangelog() *NodeVersionUpsert {
	u.SetNull(nodeversion.FieldChangelog)
	return u
}

// SetPipDependencies sets the "pip_dependencies" field.
func (u *NodeVersionUpsert) SetPipDependencies(v []string) *NodeVersionUpsert {
	u.Set(nodeversion.FieldPipDependencies, v)
	return u
}

// UpdatePipDependencies sets the "pip_dependencies" field to the value that was provided on create.
func (u *NodeVersionUpsert) UpdatePipDependencies() *NodeVersionUpsert {
	u.SetExcluded(nodeversion.FieldPipDependencies)
	return u
}

// SetDeprecated sets the "deprecated" field.
func (u *NodeVersionUpsert) SetDeprecated(v bool) *NodeVersionUpsert {
	u.Set(nodeversion.FieldDeprecated, v)
	return u
}

// UpdateDeprecated sets the "deprecated" field to the value that was provided on create.
func (u *NodeVersionUpsert) UpdateDeprecated() *NodeVersionUpsert {
	u.SetExcluded(nodeversion.FieldDeprecated)
	return u
}

// SetStatus sets the "status" field.
func (u *NodeVersionUpsert) SetStatus(v schema.NodeVersionStatus) *NodeVersionUpsert {
	u.Set(nodeversion.FieldStatus, v)
	return u
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *NodeVersionUpsert) UpdateStatus() *NodeVersionUpsert {
	u.SetExcluded(nodeversion.FieldStatus)
	return u
}

// SetStatusReason sets the "status_reason" field.
func (u *NodeVersionUpsert) SetStatusReason(v string) *NodeVersionUpsert {
	u.Set(nodeversion.FieldStatusReason, v)
	return u
}

// UpdateStatusReason sets the "status_reason" field to the value that was provided on create.
func (u *NodeVersionUpsert) UpdateStatusReason() *NodeVersionUpsert {
	u.SetExcluded(nodeversion.FieldStatusReason)
	return u
}

// SetComfyNodeExtractStatus sets the "comfy_node_extract_status" field.
func (u *NodeVersionUpsert) SetComfyNodeExtractStatus(v schema.ComfyNodeExtractStatus) *NodeVersionUpsert {
	u.Set(nodeversion.FieldComfyNodeExtractStatus, v)
	return u
}

// UpdateComfyNodeExtractStatus sets the "comfy_node_extract_status" field to the value that was provided on create.
func (u *NodeVersionUpsert) UpdateComfyNodeExtractStatus() *NodeVersionUpsert {
	u.SetExcluded(nodeversion.FieldComfyNodeExtractStatus)
	return u
}

// SetComfyNodeCloudBuildInfo sets the "comfy_node_cloud_build_info" field.
func (u *NodeVersionUpsert) SetComfyNodeCloudBuildInfo(v schema.ComfyNodeCloudBuildInfo) *NodeVersionUpsert {
	u.Set(nodeversion.FieldComfyNodeCloudBuildInfo, v)
	return u
}

// UpdateComfyNodeCloudBuildInfo sets the "comfy_node_cloud_build_info" field to the value that was provided on create.
func (u *NodeVersionUpsert) UpdateComfyNodeCloudBuildInfo() *NodeVersionUpsert {
	u.SetExcluded(nodeversion.FieldComfyNodeCloudBuildInfo)
	return u
}

// ClearComfyNodeCloudBuildInfo clears the value of the "comfy_node_cloud_build_info" field.
func (u *NodeVersionUpsert) ClearComfyNodeCloudBuildInfo() *NodeVersionUpsert {
	u.SetNull(nodeversion.FieldComfyNodeCloudBuildInfo)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.NodeVersion.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(nodeversion.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *NodeVersionUpsertOne) UpdateNewValues() *NodeVersionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(nodeversion.FieldID)
		}
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(nodeversion.FieldCreateTime)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.NodeVersion.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *NodeVersionUpsertOne) Ignore() *NodeVersionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *NodeVersionUpsertOne) DoNothing() *NodeVersionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the NodeVersionCreate.OnConflict
// documentation for more info.
func (u *NodeVersionUpsertOne) Update(set func(*NodeVersionUpsert)) *NodeVersionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&NodeVersionUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "update_time" field.
func (u *NodeVersionUpsertOne) SetUpdateTime(v time.Time) *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *NodeVersionUpsertOne) UpdateUpdateTime() *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetNodeID sets the "node_id" field.
func (u *NodeVersionUpsertOne) SetNodeID(v string) *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.SetNodeID(v)
	})
}

// UpdateNodeID sets the "node_id" field to the value that was provided on create.
func (u *NodeVersionUpsertOne) UpdateNodeID() *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.UpdateNodeID()
	})
}

// SetVersion sets the "version" field.
func (u *NodeVersionUpsertOne) SetVersion(v string) *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.SetVersion(v)
	})
}

// UpdateVersion sets the "version" field to the value that was provided on create.
func (u *NodeVersionUpsertOne) UpdateVersion() *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.UpdateVersion()
	})
}

// SetChangelog sets the "changelog" field.
func (u *NodeVersionUpsertOne) SetChangelog(v string) *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.SetChangelog(v)
	})
}

// UpdateChangelog sets the "changelog" field to the value that was provided on create.
func (u *NodeVersionUpsertOne) UpdateChangelog() *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.UpdateChangelog()
	})
}

// ClearChangelog clears the value of the "changelog" field.
func (u *NodeVersionUpsertOne) ClearChangelog() *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.ClearChangelog()
	})
}

// SetPipDependencies sets the "pip_dependencies" field.
func (u *NodeVersionUpsertOne) SetPipDependencies(v []string) *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.SetPipDependencies(v)
	})
}

// UpdatePipDependencies sets the "pip_dependencies" field to the value that was provided on create.
func (u *NodeVersionUpsertOne) UpdatePipDependencies() *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.UpdatePipDependencies()
	})
}

// SetDeprecated sets the "deprecated" field.
func (u *NodeVersionUpsertOne) SetDeprecated(v bool) *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.SetDeprecated(v)
	})
}

// UpdateDeprecated sets the "deprecated" field to the value that was provided on create.
func (u *NodeVersionUpsertOne) UpdateDeprecated() *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.UpdateDeprecated()
	})
}

// SetStatus sets the "status" field.
func (u *NodeVersionUpsertOne) SetStatus(v schema.NodeVersionStatus) *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.SetStatus(v)
	})
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *NodeVersionUpsertOne) UpdateStatus() *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.UpdateStatus()
	})
}

// SetStatusReason sets the "status_reason" field.
func (u *NodeVersionUpsertOne) SetStatusReason(v string) *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.SetStatusReason(v)
	})
}

// UpdateStatusReason sets the "status_reason" field to the value that was provided on create.
func (u *NodeVersionUpsertOne) UpdateStatusReason() *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.UpdateStatusReason()
	})
}

// SetComfyNodeExtractStatus sets the "comfy_node_extract_status" field.
func (u *NodeVersionUpsertOne) SetComfyNodeExtractStatus(v schema.ComfyNodeExtractStatus) *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.SetComfyNodeExtractStatus(v)
	})
}

// UpdateComfyNodeExtractStatus sets the "comfy_node_extract_status" field to the value that was provided on create.
func (u *NodeVersionUpsertOne) UpdateComfyNodeExtractStatus() *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.UpdateComfyNodeExtractStatus()
	})
}

// SetComfyNodeCloudBuildInfo sets the "comfy_node_cloud_build_info" field.
func (u *NodeVersionUpsertOne) SetComfyNodeCloudBuildInfo(v schema.ComfyNodeCloudBuildInfo) *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.SetComfyNodeCloudBuildInfo(v)
	})
}

// UpdateComfyNodeCloudBuildInfo sets the "comfy_node_cloud_build_info" field to the value that was provided on create.
func (u *NodeVersionUpsertOne) UpdateComfyNodeCloudBuildInfo() *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.UpdateComfyNodeCloudBuildInfo()
	})
}

// ClearComfyNodeCloudBuildInfo clears the value of the "comfy_node_cloud_build_info" field.
func (u *NodeVersionUpsertOne) ClearComfyNodeCloudBuildInfo() *NodeVersionUpsertOne {
	return u.Update(func(s *NodeVersionUpsert) {
		s.ClearComfyNodeCloudBuildInfo()
	})
}

// Exec executes the query.
func (u *NodeVersionUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for NodeVersionCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *NodeVersionUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *NodeVersionUpsertOne) ID(ctx context.Context) (id uuid.UUID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("ent: NodeVersionUpsertOne.ID is not supported by MySQL driver. Use NodeVersionUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *NodeVersionUpsertOne) IDX(ctx context.Context) uuid.UUID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// NodeVersionCreateBulk is the builder for creating many NodeVersion entities in bulk.
type NodeVersionCreateBulk struct {
	config
	err      error
	builders []*NodeVersionCreate
	conflict []sql.ConflictOption
}

// Save creates the NodeVersion entities in the database.
func (nvcb *NodeVersionCreateBulk) Save(ctx context.Context) ([]*NodeVersion, error) {
	if nvcb.err != nil {
		return nil, nvcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(nvcb.builders))
	nodes := make([]*NodeVersion, len(nvcb.builders))
	mutators := make([]Mutator, len(nvcb.builders))
	for i := range nvcb.builders {
		func(i int, root context.Context) {
			builder := nvcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*NodeVersionMutation)
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
					_, err = mutators[i+1].Mutate(root, nvcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = nvcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, nvcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, nvcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (nvcb *NodeVersionCreateBulk) SaveX(ctx context.Context) []*NodeVersion {
	v, err := nvcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (nvcb *NodeVersionCreateBulk) Exec(ctx context.Context) error {
	_, err := nvcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (nvcb *NodeVersionCreateBulk) ExecX(ctx context.Context) {
	if err := nvcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.NodeVersion.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.NodeVersionUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (nvcb *NodeVersionCreateBulk) OnConflict(opts ...sql.ConflictOption) *NodeVersionUpsertBulk {
	nvcb.conflict = opts
	return &NodeVersionUpsertBulk{
		create: nvcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.NodeVersion.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (nvcb *NodeVersionCreateBulk) OnConflictColumns(columns ...string) *NodeVersionUpsertBulk {
	nvcb.conflict = append(nvcb.conflict, sql.ConflictColumns(columns...))
	return &NodeVersionUpsertBulk{
		create: nvcb,
	}
}

// NodeVersionUpsertBulk is the builder for "upsert"-ing
// a bulk of NodeVersion nodes.
type NodeVersionUpsertBulk struct {
	create *NodeVersionCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.NodeVersion.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(nodeversion.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *NodeVersionUpsertBulk) UpdateNewValues() *NodeVersionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(nodeversion.FieldID)
			}
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(nodeversion.FieldCreateTime)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.NodeVersion.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *NodeVersionUpsertBulk) Ignore() *NodeVersionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *NodeVersionUpsertBulk) DoNothing() *NodeVersionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the NodeVersionCreateBulk.OnConflict
// documentation for more info.
func (u *NodeVersionUpsertBulk) Update(set func(*NodeVersionUpsert)) *NodeVersionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&NodeVersionUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "update_time" field.
func (u *NodeVersionUpsertBulk) SetUpdateTime(v time.Time) *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *NodeVersionUpsertBulk) UpdateUpdateTime() *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetNodeID sets the "node_id" field.
func (u *NodeVersionUpsertBulk) SetNodeID(v string) *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.SetNodeID(v)
	})
}

// UpdateNodeID sets the "node_id" field to the value that was provided on create.
func (u *NodeVersionUpsertBulk) UpdateNodeID() *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.UpdateNodeID()
	})
}

// SetVersion sets the "version" field.
func (u *NodeVersionUpsertBulk) SetVersion(v string) *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.SetVersion(v)
	})
}

// UpdateVersion sets the "version" field to the value that was provided on create.
func (u *NodeVersionUpsertBulk) UpdateVersion() *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.UpdateVersion()
	})
}

// SetChangelog sets the "changelog" field.
func (u *NodeVersionUpsertBulk) SetChangelog(v string) *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.SetChangelog(v)
	})
}

// UpdateChangelog sets the "changelog" field to the value that was provided on create.
func (u *NodeVersionUpsertBulk) UpdateChangelog() *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.UpdateChangelog()
	})
}

// ClearChangelog clears the value of the "changelog" field.
func (u *NodeVersionUpsertBulk) ClearChangelog() *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.ClearChangelog()
	})
}

// SetPipDependencies sets the "pip_dependencies" field.
func (u *NodeVersionUpsertBulk) SetPipDependencies(v []string) *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.SetPipDependencies(v)
	})
}

// UpdatePipDependencies sets the "pip_dependencies" field to the value that was provided on create.
func (u *NodeVersionUpsertBulk) UpdatePipDependencies() *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.UpdatePipDependencies()
	})
}

// SetDeprecated sets the "deprecated" field.
func (u *NodeVersionUpsertBulk) SetDeprecated(v bool) *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.SetDeprecated(v)
	})
}

// UpdateDeprecated sets the "deprecated" field to the value that was provided on create.
func (u *NodeVersionUpsertBulk) UpdateDeprecated() *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.UpdateDeprecated()
	})
}

// SetStatus sets the "status" field.
func (u *NodeVersionUpsertBulk) SetStatus(v schema.NodeVersionStatus) *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.SetStatus(v)
	})
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *NodeVersionUpsertBulk) UpdateStatus() *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.UpdateStatus()
	})
}

// SetStatusReason sets the "status_reason" field.
func (u *NodeVersionUpsertBulk) SetStatusReason(v string) *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.SetStatusReason(v)
	})
}

// UpdateStatusReason sets the "status_reason" field to the value that was provided on create.
func (u *NodeVersionUpsertBulk) UpdateStatusReason() *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.UpdateStatusReason()
	})
}

// SetComfyNodeExtractStatus sets the "comfy_node_extract_status" field.
func (u *NodeVersionUpsertBulk) SetComfyNodeExtractStatus(v schema.ComfyNodeExtractStatus) *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.SetComfyNodeExtractStatus(v)
	})
}

// UpdateComfyNodeExtractStatus sets the "comfy_node_extract_status" field to the value that was provided on create.
func (u *NodeVersionUpsertBulk) UpdateComfyNodeExtractStatus() *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.UpdateComfyNodeExtractStatus()
	})
}

// SetComfyNodeCloudBuildInfo sets the "comfy_node_cloud_build_info" field.
func (u *NodeVersionUpsertBulk) SetComfyNodeCloudBuildInfo(v schema.ComfyNodeCloudBuildInfo) *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.SetComfyNodeCloudBuildInfo(v)
	})
}

// UpdateComfyNodeCloudBuildInfo sets the "comfy_node_cloud_build_info" field to the value that was provided on create.
func (u *NodeVersionUpsertBulk) UpdateComfyNodeCloudBuildInfo() *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.UpdateComfyNodeCloudBuildInfo()
	})
}

// ClearComfyNodeCloudBuildInfo clears the value of the "comfy_node_cloud_build_info" field.
func (u *NodeVersionUpsertBulk) ClearComfyNodeCloudBuildInfo() *NodeVersionUpsertBulk {
	return u.Update(func(s *NodeVersionUpsert) {
		s.ClearComfyNodeCloudBuildInfo()
	})
}

// Exec executes the query.
func (u *NodeVersionUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the NodeVersionCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for NodeVersionCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *NodeVersionUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"registry-backend/ent/predicate"
	"registry-backend/ent/storagefile"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// StorageFileUpdate is the builder for updating StorageFile entities.
type StorageFileUpdate struct {
	config
	hooks    []Hook
	mutation *StorageFileMutation
}

// Where appends a list predicates to the StorageFileUpdate builder.
func (sfu *StorageFileUpdate) Where(ps ...predicate.StorageFile) *StorageFileUpdate {
	sfu.mutation.Where(ps...)
	return sfu
}

// SetUpdateTime sets the "update_time" field.
func (sfu *StorageFileUpdate) SetUpdateTime(t time.Time) *StorageFileUpdate {
	sfu.mutation.SetUpdateTime(t)
	return sfu
}

// SetBucketName sets the "bucket_name" field.
func (sfu *StorageFileUpdate) SetBucketName(s string) *StorageFileUpdate {
	sfu.mutation.SetBucketName(s)
	return sfu
}

// SetNillableBucketName sets the "bucket_name" field if the given value is not nil.
func (sfu *StorageFileUpdate) SetNillableBucketName(s *string) *StorageFileUpdate {
	if s != nil {
		sfu.SetBucketName(*s)
	}
	return sfu
}

// SetObjectName sets the "object_name" field.
func (sfu *StorageFileUpdate) SetObjectName(s string) *StorageFileUpdate {
	sfu.mutation.SetObjectName(s)
	return sfu
}

// SetNillableObjectName sets the "object_name" field if the given value is not nil.
func (sfu *StorageFileUpdate) SetNillableObjectName(s *string) *StorageFileUpdate {
	if s != nil {
		sfu.SetObjectName(*s)
	}
	return sfu
}

// ClearObjectName clears the value of the "object_name" field.
func (sfu *StorageFileUpdate) ClearObjectName() *StorageFileUpdate {
	sfu.mutation.ClearObjectName()
	return sfu
}

// SetFilePath sets the "file_path" field.
func (sfu *StorageFileUpdate) SetFilePath(s string) *StorageFileUpdate {
	sfu.mutation.SetFilePath(s)
	return sfu
}

// SetNillableFilePath sets the "file_path" field if the given value is not nil.
func (sfu *StorageFileUpdate) SetNillableFilePath(s *string) *StorageFileUpdate {
	if s != nil {
		sfu.SetFilePath(*s)
	}
	return sfu
}

// SetFileType sets the "file_type" field.
func (sfu *StorageFileUpdate) SetFileType(s string) *StorageFileUpdate {
	sfu.mutation.SetFileType(s)
	return sfu
}

// SetNillableFileType sets the "file_type" field if the given value is not nil.
func (sfu *StorageFileUpdate) SetNillableFileType(s *string) *StorageFileUpdate {
	if s != nil {
		sfu.SetFileType(*s)
	}
	return sfu
}

// SetFileURL sets the "file_url" field.
func (sfu *StorageFileUpdate) SetFileURL(s string) *StorageFileUpdate {
	sfu.mutation.SetFileURL(s)
	return sfu
}

// SetNillableFileURL sets the "file_url" field if the given value is not nil.
func (sfu *StorageFileUpdate) SetNillableFileURL(s *string) *StorageFileUpdate {
	if s != nil {
		sfu.SetFileURL(*s)
	}
	return sfu
}

// ClearFileURL clears the value of the "file_url" field.
func (sfu *StorageFileUpdate) ClearFileURL() *StorageFileUpdate {
	sfu.mutation.ClearFileURL()
	return sfu
}

// Mutation returns the StorageFileMutation object of the builder.
func (sfu *StorageFileUpdate) Mutation() *StorageFileMutation {
	return sfu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (sfu *StorageFileUpdate) Save(ctx context.Context) (int, error) {
	sfu.defaults()
	return withHooks(ctx, sfu.sqlSave, sfu.mutation, sfu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (sfu *StorageFileUpdate) SaveX(ctx context.Context) int {
	affected, err := sfu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (sfu *StorageFileUpdate) Exec(ctx context.Context) error {
	_, err := sfu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sfu *StorageFileUpdate) ExecX(ctx context.Context) {
	if err := sfu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (sfu *StorageFileUpdate) defaults() {
	if _, ok := sfu.mutation.UpdateTime(); !ok {
		v := storagefile.UpdateDefaultUpdateTime()
		sfu.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sfu *StorageFileUpdate) check() error {
	if v, ok := sfu.mutation.BucketName(); ok {
		if err := storagefile.BucketNameValidator(v); err != nil {
			return &ValidationError{Name: "bucket_name", err: fmt.Errorf(`ent: validator failed for field "StorageFile.bucket_name": %w`, err)}
		}
	}
	if v, ok := sfu.mutation.ObjectName(); ok {
		if err := storagefile.ObjectNameValidator(v); err != nil {
			return &ValidationError{Name: "object_name", err: fmt.Errorf(`ent: validator failed for field "StorageFile.object_name": %w`, err)}
		}
	}
	if v, ok := sfu.mutation.FilePath(); ok {
		if err := storagefile.FilePathValidator(v); err != nil {
			return &ValidationError{Name: "file_path", err: fmt.Errorf(`ent: validator failed for field "StorageFile.file_path": %w`, err)}
		}
	}
	if v, ok := sfu.mutation.FileType(); ok {
		if err := storagefile.FileTypeValidator(v); err != nil {
			return &ValidationError{Name: "file_type", err: fmt.Errorf(`ent: validator failed for field "StorageFile.file_type": %w`, err)}
		}
	}
	return nil
}

func (sfu *StorageFileUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := sfu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(storagefile.Table, storagefile.Columns, sqlgraph.NewFieldSpec(storagefile.FieldID, field.TypeUUID))
	if ps := sfu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := sfu.mutation.UpdateTime(); ok {
		_spec.SetField(storagefile.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := sfu.mutation.BucketName(); ok {
		_spec.SetField(storagefile.FieldBucketName, field.TypeString, value)
	}
	if value, ok := sfu.mutation.ObjectName(); ok {
		_spec.SetField(storagefile.FieldObjectName, field.TypeString, value)
	}
	if sfu.mutation.ObjectNameCleared() {
		_spec.ClearField(storagefile.FieldObjectName, field.TypeString)
	}
	if value, ok := sfu.mutation.FilePath(); ok {
		_spec.SetField(storagefile.FieldFilePath, field.TypeString, value)
	}
	if value, ok := sfu.mutation.FileType(); ok {
		_spec.SetField(storagefile.FieldFileType, field.TypeString, value)
	}
	if value, ok := sfu.mutation.FileURL(); ok {
		_spec.SetField(storagefile.FieldFileURL, field.TypeString, value)
	}
	if sfu.mutation.FileURLCleared() {
		_spec.ClearField(storagefile.FieldFileURL, field.TypeString)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, sfu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{storagefile.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	sfu.mutation.done = true
	return n, nil
}

// StorageFileUpdateOne is the builder for updating a single StorageFile entity.
type StorageFileUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *StorageFileMutation
}

// SetUpdateTime sets the "update_time" field.
func (sfuo *StorageFileUpdateOne) SetUpdateTime(t time.Time) *StorageFileUpdateOne {
	sfuo.mutation.SetUpdateTime(t)
	return sfuo
}

// SetBucketName sets the "bucket_name" field.
func (sfuo *StorageFileUpdateOne) SetBucketName(s string) *StorageFileUpdateOne {
	sfuo.mutation.SetBucketName(s)
	return sfuo
}

// SetNillableBucketName sets the "bucket_name" field if the given value is not nil.
func (sfuo *StorageFileUpdateOne) SetNillableBucketName(s *string) *StorageFileUpdateOne {
	if s != nil {
		sfuo.SetBucketName(*s)
	}
	return sfuo
}

// SetObjectName sets the "object_name" field.
func (sfuo *StorageFileUpdateOne) SetObjectName(s string) *StorageFileUpdateOne {
	sfuo.mutation.SetObjectName(s)
	return sfuo
}

// SetNillableObjectName sets the "object_name" field if the given value is not nil.
func (sfuo *StorageFileUpdateOne) SetNillableObjectName(s *string) *StorageFileUpdateOne {
	if s != nil {
		sfuo.SetObjectName(*s)
	}
	return sfuo
}

// ClearObjectName clears the value of the "object_name" field.
func (sfuo *StorageFileUpdateOne) ClearObjectName() *StorageFileUpdateOne {
	sfuo.mutation.ClearObjectName()
	return sfuo
}

// SetFilePath sets the "file_path" field.
func (sfuo *StorageFileUpdateOne) SetFilePath(s string) *StorageFileUpdateOne {
	sfuo.mutation.SetFilePath(s)
	return sfuo
}

// SetNillableFilePath sets the "file_path" field if the given value is not nil.
func (sfuo *StorageFileUpdateOne) SetNillableFilePath(s *string) *StorageFileUpdateOne {
	if s != nil {
		sfuo.SetFilePath(*s)
	}
	return sfuo
}

// SetFileType sets the "file_type" field.
func (sfuo *StorageFileUpdateOne) SetFileType(s string) *StorageFileUpdateOne {
	sfuo.mutation.SetFileType(s)
	return sfuo
}

// SetNillableFileType sets the "file_type" field if the given value is not nil.
func (sfuo *StorageFileUpdateOne) SetNillableFileType(s *string) *StorageFileUpdateOne {
	if s != nil {
		sfuo.SetFileType(*s)
	}
	return sfuo
}

// SetFileURL sets the "file_url" field.
func (sfuo *StorageFileUpdateOne) SetFileURL(s string) *StorageFileUpdateOne {
	sfuo.mutation.SetFileURL(s)
	return sfuo
}

// SetNillableFileURL sets the "file_url" field if the given value is not nil.
func (sfuo *StorageFileUpdateOne) SetNillableFileURL(s *string) *StorageFileUpdateOne {
	if s != nil {
		sfuo.SetFileURL(*s)
	}
	return sfuo
}

// ClearFileURL clears the value of the "file_url" field.
func (sfuo *StorageFileUpdateOne) ClearFileURL() *StorageFileUpdateOne {
	sfuo.mutation.ClearFileURL()
	return sfuo
}

// Mutation returns the StorageFileMutation object of the builder.
func (sfuo *StorageFileUpdateOne) Mutation() *StorageFileMutation {
	return sfuo.mutation
}

// Where appends a list predicates to the StorageFileUpdate builder.
func (sfuo *StorageFileUpdateOne) Where(ps ...predicate.StorageFile) *StorageFileUpdateOne {
	sfuo.mutation.Where(ps...)
	return sfuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (sfuo *StorageFileUpdateOne) Select(field string, fields ...string) *StorageFileUpdateOne {
	sfuo.fields = append([]string{field}, fields...)
	return sfuo
}

// Save executes the query and returns the updated StorageFile entity.
func (sfuo *StorageFileUpdateOne) Save(ctx context.Context) (*StorageFile, error) {
	sfuo.defaults()
	return withHooks(ctx, sfuo.sqlSave, sfuo.mutation, sfuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (sfuo *StorageFileUpdateOne) SaveX(ctx context.Context) *StorageFile {
	node, err := sfuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (sfuo *StorageFileUpdateOne) Exec(ctx context.Context) error {
	_, err := sfuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sfuo *StorageFileUpdateOne) ExecX(ctx context.Context) {
	if err := sfuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (sfuo *StorageFileUpdateOne) defaults() {
	if _, ok := sfuo.mutation.UpdateTime(); !ok {
		v := storagefile.UpdateDefaultUpdateTime()
		sfuo.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sfuo *StorageFileUpdateOne) check() error {
	if v, ok := sfuo.mutation.BucketName(); ok {
		if err := storagefile.BucketNameValidator(v); err != nil {
			return &ValidationError{Name: "bucket_name", err: fmt.Errorf(`ent: validator failed for field "StorageFile.bucket_name": %w`, err)}
		}
	}
	if v, ok := sfuo.mutation.ObjectName(); ok {
		if err := storagefile.ObjectNameValidator(v); err != nil {
			return &ValidationError{Name: "object_name", err: fmt.Errorf(`ent: validator failed for field "StorageFile.object_name": %w`, err)}
		}
	}
	if v, ok := sfuo.mutation.FilePath(); ok {
		if err := storagefile.FilePathValidator(v); err != nil {
			return &ValidationError{Name: "file_path", err: fmt.Errorf(`ent: validator failed for field "StorageFile.file_path": %w`, err)}
		}
	}
	if v, ok := sfuo.mutation.FileType(); ok {
		if err := storagefile.FileTypeValidator(v); err != nil {
			return &ValidationError{Name: "file_type", err: fmt.Errorf(`ent: validator failed for field "StorageFile.file_type": %w`, err)}
		}
	}
	return nil
}

func (sfuo *StorageFileUpdateOne) sqlSave(ctx context.Context) (_node *StorageFile, err error) {
	if err := sfuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(storagefile.Table, storagefile.Columns, sqlgraph.NewFieldSpec(storagefile.FieldID, field.TypeUUID))
	id, ok := sfuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "StorageFile.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := sfuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, storagefile.FieldID)
		for _, f := range fields {
			if !storagefile.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != storagefile.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := sfuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := sfuo.mutation.UpdateTime(); ok {
		_spec.SetField(storagefile.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := sfuo.mutation.BucketName(); ok {
		_spec.SetField(storagefile.FieldBucketName, field.TypeString, value)
	}
	if value, ok := sfuo.mutation.ObjectName(); ok {
		_spec.SetField(storagefile.FieldObjectName, field.TypeString, value)
	}
	if sfuo.mutation.ObjectNameCleared() {
		_spec.ClearField(storagefile.FieldObjectName, field.TypeString)
	}
	if value, ok := sfuo.mutation.FilePath(); ok {
		_spec.SetField(storagefile.FieldFilePath, field.TypeString, value)
	}
	if value, ok := sfuo.mutation.FileType(); ok {
		_spec.SetField(storagefile.FieldFileType, field.TypeString, value)
	}
	if value, ok := sfuo.mutation.FileURL(); ok {
		_spec.SetField(storagefile.FieldFileURL, field.TypeString, value)
	}
	if sfuo.mutation.FileURLCleared() {
		_spec.ClearField(storagefile.FieldFileURL, field.TypeString)
	}
	_node = &StorageFile{config: sfuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, sfuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{storagefile.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	sfuo.mutation.done = true
	return _node, nil
}

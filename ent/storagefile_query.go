// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"
	"registry-backend/ent/predicate"
	"registry-backend/ent/storagefile"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// StorageFileQuery is the builder for querying StorageFile entities.
type StorageFileQuery struct {
	config
	ctx        *QueryContext
	order      []storagefile.OrderOption
	inters     []Interceptor
	predicates []predicate.StorageFile
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the StorageFileQuery builder.
func (sfq *StorageFileQuery) Where(ps ...predicate.StorageFile) *StorageFileQuery {
	sfq.predicates = append(sfq.predicates, ps...)
	return sfq
}

// Limit the number of records to be returned by this query.
func (sfq *StorageFileQuery) Limit(limit int) *StorageFileQuery {
	sfq.ctx.Limit = &limit
	return sfq
}

// Offset to start from.
func (sfq *StorageFileQuery) Offset(offset int) *StorageFileQuery {
	sfq.ctx.Offset = &offset
	return sfq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (sfq *StorageFileQuery) Unique(unique bool) *StorageFileQuery {
	sfq.ctx.Unique = &unique
	return sfq
}

// Order specifies how the records should be ordered.
func (sfq *StorageFileQuery) Order(o ...storagefile.OrderOption) *StorageFileQuery {
	sfq.order = append(sfq.order, o...)
	return sfq
}

// First returns the first StorageFile entity from the query.
// Returns a *NotFoundError when no StorageFile was found.
func (sfq *StorageFileQuery) First(ctx context.Context) (*StorageFile, error) {
	nodes, err := sfq.Limit(1).All(setContextOp(ctx, sfq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{storagefile.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (sfq *StorageFileQuery) FirstX(ctx context.Context) *StorageFile {
	node, err := sfq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first StorageFile ID from the query.
// Returns a *NotFoundError when no StorageFile ID was found.
func (sfq *StorageFileQuery) FirstID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = sfq.Limit(1).IDs(setContextOp(ctx, sfq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{storagefile.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (sfq *StorageFileQuery) FirstIDX(ctx context.Context) uuid.UUID {
	id, err := sfq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single StorageFile entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one StorageFile entity is found.
// Returns a *NotFoundError when no StorageFile entities are found.
func (sfq *StorageFileQuery) Only(ctx context.Context) (*StorageFile, error) {
	nodes, err := sfq.Limit(2).All(setContextOp(ctx, sfq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{storagefile.Label}
	default:
		return nil, &NotSingularError{storagefile.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (sfq *StorageFileQuery) OnlyX(ctx context.Context) *StorageFile {
	node, err := sfq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only StorageFile ID in the query.
// Returns a *NotSingularError when more than one StorageFile ID is found.
// Returns a *NotFoundError when no entities are found.
func (sfq *StorageFileQuery) OnlyID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = sfq.Limit(2).IDs(setContextOp(ctx, sfq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{storagefile.Label}
	default:
		err = &NotSingularError{storagefile.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (sfq *StorageFileQuery) OnlyIDX(ctx context.Context) uuid.UUID {
	id, err := sfq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of StorageFiles.
func (sfq *StorageFileQuery) All(ctx context.Context) ([]*StorageFile, error) {
	ctx = setContextOp(ctx, sfq.ctx, "All")
	if err := sfq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*StorageFile, *StorageFileQuery]()
	return withInterceptors[[]*StorageFile](ctx, sfq, qr, sfq.inters)
}

// AllX is like All, but panics if an error occurs.
func (sfq *StorageFileQuery) AllX(ctx context.Context) []*StorageFile {
	nodes, err := sfq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of StorageFile IDs.
func (sfq *StorageFileQuery) IDs(ctx context.Context) (ids []uuid.UUID, err error) {
	if sfq.ctx.Unique == nil && sfq.path != nil {
		sfq.Unique(true)
	}
	ctx = setContextOp(ctx, sfq.ctx, "IDs")
	if err = sfq.Select(storagefile.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (sfq *StorageFileQuery) IDsX(ctx context.Context) []uuid.UUID {
	ids, err := sfq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (sfq *StorageFileQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, sfq.ctx, "Count")
	if err := sfq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, sfq, querierCount[*StorageFileQuery](), sfq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (sfq *StorageFileQuery) CountX(ctx context.Context) int {
	count, err := sfq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (sfq *StorageFileQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, sfq.ctx, "Exist")
	switch _, err := sfq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (sfq *StorageFileQuery) ExistX(ctx context.Context) bool {
	exist, err := sfq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the StorageFileQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (sfq *StorageFileQuery) Clone() *StorageFileQuery {
	if sfq == nil {
		return nil
	}
	return &StorageFileQuery{
		config:     sfq.config,
		ctx:        sfq.ctx.Clone(),
		order:      append([]storagefile.OrderOption{}, sfq.order...),
		inters:     append([]Interceptor{}, sfq.inters...),
		predicates: append([]predicate.StorageFile{}, sfq.predicates...),
		// clone intermediate query.
		sql:  sfq.sql.Clone(),
		path: sfq.path,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		CreateTime time.Time `json:"create_time,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.StorageFile.Query().
//		GroupBy(storagefile.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (sfq *StorageFileQuery) GroupBy(field string, fields ...string) *StorageFileGroupBy {
	sfq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &StorageFileGroupBy{build: sfq}
	grbuild.flds = &sfq.ctx.Fields
	grbuild.label = storagefile.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		CreateTime time.Time `json:"create_time,omitempty"`
//	}
//
//	client.StorageFile.Query().
//		Select(storagefile.FieldCreateTime).
//		Scan(ctx, &v)
func (sfq *StorageFileQuery) Select(fields ...string) *StorageFileSelect {
	sfq.ctx.Fields = append(sfq.ctx.Fields, fields...)
	sbuild := &StorageFileSelect{StorageFileQuery: sfq}
	sbuild.label = storagefile.Label
	sbuild.flds, sbuild.scan = &sfq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a StorageFileSelect configured with the given aggregations.
func (sfq *StorageFileQuery) Aggregate(fns ...AggregateFunc) *StorageFileSelect {
	return sfq.Select().Aggregate(fns...)
}

func (sfq *StorageFileQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range sfq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, sfq); err != nil {
				return err
			}
		}
	}
	for _, f := range sfq.ctx.Fields {
		if !storagefile.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if sfq.path != nil {
		prev, err := sfq.path(ctx)
		if err != nil {
			return err
		}
		sfq.sql = prev
	}
	return nil
}

func (sfq *StorageFileQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*StorageFile, error) {
	var (
		nodes = []*StorageFile{}
		_spec = sfq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*StorageFile).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &StorageFile{config: sfq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	if len(sfq.modifiers) > 0 {
		_spec.Modifiers = sfq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, sfq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (sfq *StorageFileQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := sfq.querySpec()
	if len(sfq.modifiers) > 0 {
		_spec.Modifiers = sfq.modifiers
	}
	_spec.Node.Columns = sfq.ctx.Fields
	if len(sfq.ctx.Fields) > 0 {
		_spec.Unique = sfq.ctx.Unique != nil && *sfq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, sfq.driver, _spec)
}

func (sfq *StorageFileQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(storagefile.Table, storagefile.Columns, sqlgraph.NewFieldSpec(storagefile.FieldID, field.TypeUUID))
	_spec.From = sfq.sql
	if unique := sfq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if sfq.path != nil {
		_spec.Unique = true
	}
	if fields := sfq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, storagefile.FieldID)
		for i := range fields {
			if fields[i] != storagefile.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := sfq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := sfq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := sfq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := sfq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (sfq *StorageFileQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(sfq.driver.Dialect())
	t1 := builder.Table(storagefile.Table)
	columns := sfq.ctx.Fields
	if len(columns) == 0 {
		columns = storagefile.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if sfq.sql != nil {
		selector = sfq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if sfq.ctx.Unique != nil && *sfq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range sfq.modifiers {
		m(selector)
	}
	for _, p := range sfq.predicates {
		p(selector)
	}
	for _, p := range sfq.order {
		p(selector)
	}
	if offset := sfq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := sfq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (sfq *StorageFileQuery) ForUpdate(opts ...sql.LockOption) *StorageFileQuery {
	if sfq.driver.Dialect() == dialect.Postgres {
		sfq.Unique(false)
	}
	sfq.modifiers = append(sfq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return sfq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (sfq *StorageFileQuery) ForShare(opts ...sql.LockOption) *StorageFileQuery {
	if sfq.driver.Dialect() == dialect.Postgres {
		sfq.Unique(false)
	}
	sfq.modifiers = append(sfq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return sfq
}

// StorageFileGroupBy is the group-by builder for StorageFile entities.
type StorageFileGroupBy struct {
	selector
	build *StorageFileQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (sfgb *StorageFileGroupBy) Aggregate(fns ...AggregateFunc) *StorageFileGroupBy {
	sfgb.fns = append(sfgb.fns, fns...)
	return sfgb
}

// Scan applies the selector query and scans the result into the given value.
func (sfgb *StorageFileGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, sfgb.build.ctx, "GroupBy")
	if err := sfgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*StorageFileQuery, *StorageFileGroupBy](ctx, sfgb.build, sfgb, sfgb.build.inters, v)
}

func (sfgb *StorageFileGroupBy) sqlScan(ctx context.Context, root *StorageFileQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(sfgb.fns))
	for _, fn := range sfgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*sfgb.flds)+len(sfgb.fns))
		for _, f := range *sfgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*sfgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := sfgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// StorageFileSelect is the builder for selecting fields of StorageFile entities.
type StorageFileSelect struct {
	*StorageFileQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (sfs *StorageFileSelect) Aggregate(fns ...AggregateFunc) *StorageFileSelect {
	sfs.fns = append(sfs.fns, fns...)
	return sfs
}

// Scan applies the selector query and scans the result into the given value.
func (sfs *StorageFileSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, sfs.ctx, "Select")
	if err := sfs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*StorageFileQuery, *StorageFileSelect](ctx, sfs.StorageFileQuery, sfs, sfs.inters, v)
}

func (sfs *StorageFileSelect) sqlScan(ctx context.Context, root *StorageFileQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(sfs.fns))
	for _, fn := range sfs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*sfs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := sfs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

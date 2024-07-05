// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"
	"registry-backend/ent/node"
	"registry-backend/ent/nodeversion"
	"registry-backend/ent/predicate"
	"registry-backend/ent/storagefile"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// NodeVersionQuery is the builder for querying NodeVersion entities.
type NodeVersionQuery struct {
	config
	ctx             *QueryContext
	order           []nodeversion.OrderOption
	inters          []Interceptor
	predicates      []predicate.NodeVersion
	withNode        *NodeQuery
	withStorageFile *StorageFileQuery
	withFKs         bool
	modifiers       []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the NodeVersionQuery builder.
func (nvq *NodeVersionQuery) Where(ps ...predicate.NodeVersion) *NodeVersionQuery {
	nvq.predicates = append(nvq.predicates, ps...)
	return nvq
}

// Limit the number of records to be returned by this query.
func (nvq *NodeVersionQuery) Limit(limit int) *NodeVersionQuery {
	nvq.ctx.Limit = &limit
	return nvq
}

// Offset to start from.
func (nvq *NodeVersionQuery) Offset(offset int) *NodeVersionQuery {
	nvq.ctx.Offset = &offset
	return nvq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (nvq *NodeVersionQuery) Unique(unique bool) *NodeVersionQuery {
	nvq.ctx.Unique = &unique
	return nvq
}

// Order specifies how the records should be ordered.
func (nvq *NodeVersionQuery) Order(o ...nodeversion.OrderOption) *NodeVersionQuery {
	nvq.order = append(nvq.order, o...)
	return nvq
}

// QueryNode chains the current query on the "node" edge.
func (nvq *NodeVersionQuery) QueryNode() *NodeQuery {
	query := (&NodeClient{config: nvq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := nvq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := nvq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(nodeversion.Table, nodeversion.FieldID, selector),
			sqlgraph.To(node.Table, node.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, nodeversion.NodeTable, nodeversion.NodeColumn),
		)
		fromU = sqlgraph.SetNeighbors(nvq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryStorageFile chains the current query on the "storage_file" edge.
func (nvq *NodeVersionQuery) QueryStorageFile() *StorageFileQuery {
	query := (&StorageFileClient{config: nvq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := nvq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := nvq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(nodeversion.Table, nodeversion.FieldID, selector),
			sqlgraph.To(storagefile.Table, storagefile.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, nodeversion.StorageFileTable, nodeversion.StorageFileColumn),
		)
		fromU = sqlgraph.SetNeighbors(nvq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first NodeVersion entity from the query.
// Returns a *NotFoundError when no NodeVersion was found.
func (nvq *NodeVersionQuery) First(ctx context.Context) (*NodeVersion, error) {
	nodes, err := nvq.Limit(1).All(setContextOp(ctx, nvq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{nodeversion.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (nvq *NodeVersionQuery) FirstX(ctx context.Context) *NodeVersion {
	node, err := nvq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first NodeVersion ID from the query.
// Returns a *NotFoundError when no NodeVersion ID was found.
func (nvq *NodeVersionQuery) FirstID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = nvq.Limit(1).IDs(setContextOp(ctx, nvq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{nodeversion.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (nvq *NodeVersionQuery) FirstIDX(ctx context.Context) uuid.UUID {
	id, err := nvq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single NodeVersion entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one NodeVersion entity is found.
// Returns a *NotFoundError when no NodeVersion entities are found.
func (nvq *NodeVersionQuery) Only(ctx context.Context) (*NodeVersion, error) {
	nodes, err := nvq.Limit(2).All(setContextOp(ctx, nvq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{nodeversion.Label}
	default:
		return nil, &NotSingularError{nodeversion.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (nvq *NodeVersionQuery) OnlyX(ctx context.Context) *NodeVersion {
	node, err := nvq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only NodeVersion ID in the query.
// Returns a *NotSingularError when more than one NodeVersion ID is found.
// Returns a *NotFoundError when no entities are found.
func (nvq *NodeVersionQuery) OnlyID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = nvq.Limit(2).IDs(setContextOp(ctx, nvq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{nodeversion.Label}
	default:
		err = &NotSingularError{nodeversion.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (nvq *NodeVersionQuery) OnlyIDX(ctx context.Context) uuid.UUID {
	id, err := nvq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of NodeVersions.
func (nvq *NodeVersionQuery) All(ctx context.Context) ([]*NodeVersion, error) {
	ctx = setContextOp(ctx, nvq.ctx, "All")
	if err := nvq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*NodeVersion, *NodeVersionQuery]()
	return withInterceptors[[]*NodeVersion](ctx, nvq, qr, nvq.inters)
}

// AllX is like All, but panics if an error occurs.
func (nvq *NodeVersionQuery) AllX(ctx context.Context) []*NodeVersion {
	nodes, err := nvq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of NodeVersion IDs.
func (nvq *NodeVersionQuery) IDs(ctx context.Context) (ids []uuid.UUID, err error) {
	if nvq.ctx.Unique == nil && nvq.path != nil {
		nvq.Unique(true)
	}
	ctx = setContextOp(ctx, nvq.ctx, "IDs")
	if err = nvq.Select(nodeversion.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (nvq *NodeVersionQuery) IDsX(ctx context.Context) []uuid.UUID {
	ids, err := nvq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (nvq *NodeVersionQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, nvq.ctx, "Count")
	if err := nvq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, nvq, querierCount[*NodeVersionQuery](), nvq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (nvq *NodeVersionQuery) CountX(ctx context.Context) int {
	count, err := nvq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (nvq *NodeVersionQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, nvq.ctx, "Exist")
	switch _, err := nvq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (nvq *NodeVersionQuery) ExistX(ctx context.Context) bool {
	exist, err := nvq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the NodeVersionQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (nvq *NodeVersionQuery) Clone() *NodeVersionQuery {
	if nvq == nil {
		return nil
	}
	return &NodeVersionQuery{
		config:          nvq.config,
		ctx:             nvq.ctx.Clone(),
		order:           append([]nodeversion.OrderOption{}, nvq.order...),
		inters:          append([]Interceptor{}, nvq.inters...),
		predicates:      append([]predicate.NodeVersion{}, nvq.predicates...),
		withNode:        nvq.withNode.Clone(),
		withStorageFile: nvq.withStorageFile.Clone(),
		// clone intermediate query.
		sql:  nvq.sql.Clone(),
		path: nvq.path,
	}
}

// WithNode tells the query-builder to eager-load the nodes that are connected to
// the "node" edge. The optional arguments are used to configure the query builder of the edge.
func (nvq *NodeVersionQuery) WithNode(opts ...func(*NodeQuery)) *NodeVersionQuery {
	query := (&NodeClient{config: nvq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	nvq.withNode = query
	return nvq
}

// WithStorageFile tells the query-builder to eager-load the nodes that are connected to
// the "storage_file" edge. The optional arguments are used to configure the query builder of the edge.
func (nvq *NodeVersionQuery) WithStorageFile(opts ...func(*StorageFileQuery)) *NodeVersionQuery {
	query := (&StorageFileClient{config: nvq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	nvq.withStorageFile = query
	return nvq
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
//	client.NodeVersion.Query().
//		GroupBy(nodeversion.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (nvq *NodeVersionQuery) GroupBy(field string, fields ...string) *NodeVersionGroupBy {
	nvq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &NodeVersionGroupBy{build: nvq}
	grbuild.flds = &nvq.ctx.Fields
	grbuild.label = nodeversion.Label
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
//	client.NodeVersion.Query().
//		Select(nodeversion.FieldCreateTime).
//		Scan(ctx, &v)
func (nvq *NodeVersionQuery) Select(fields ...string) *NodeVersionSelect {
	nvq.ctx.Fields = append(nvq.ctx.Fields, fields...)
	sbuild := &NodeVersionSelect{NodeVersionQuery: nvq}
	sbuild.label = nodeversion.Label
	sbuild.flds, sbuild.scan = &nvq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a NodeVersionSelect configured with the given aggregations.
func (nvq *NodeVersionQuery) Aggregate(fns ...AggregateFunc) *NodeVersionSelect {
	return nvq.Select().Aggregate(fns...)
}

func (nvq *NodeVersionQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range nvq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, nvq); err != nil {
				return err
			}
		}
	}
	for _, f := range nvq.ctx.Fields {
		if !nodeversion.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if nvq.path != nil {
		prev, err := nvq.path(ctx)
		if err != nil {
			return err
		}
		nvq.sql = prev
	}
	return nil
}

func (nvq *NodeVersionQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*NodeVersion, error) {
	var (
		nodes       = []*NodeVersion{}
		withFKs     = nvq.withFKs
		_spec       = nvq.querySpec()
		loadedTypes = [2]bool{
			nvq.withNode != nil,
			nvq.withStorageFile != nil,
		}
	)
	if nvq.withStorageFile != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, nodeversion.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*NodeVersion).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &NodeVersion{config: nvq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(nvq.modifiers) > 0 {
		_spec.Modifiers = nvq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, nvq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := nvq.withNode; query != nil {
		if err := nvq.loadNode(ctx, query, nodes, nil,
			func(n *NodeVersion, e *Node) { n.Edges.Node = e }); err != nil {
			return nil, err
		}
	}
	if query := nvq.withStorageFile; query != nil {
		if err := nvq.loadStorageFile(ctx, query, nodes, nil,
			func(n *NodeVersion, e *StorageFile) { n.Edges.StorageFile = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (nvq *NodeVersionQuery) loadNode(ctx context.Context, query *NodeQuery, nodes []*NodeVersion, init func(*NodeVersion), assign func(*NodeVersion, *Node)) error {
	ids := make([]string, 0, len(nodes))
	nodeids := make(map[string][]*NodeVersion)
	for i := range nodes {
		fk := nodes[i].NodeID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(node.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "node_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (nvq *NodeVersionQuery) loadStorageFile(ctx context.Context, query *StorageFileQuery, nodes []*NodeVersion, init func(*NodeVersion), assign func(*NodeVersion, *StorageFile)) error {
	ids := make([]uuid.UUID, 0, len(nodes))
	nodeids := make(map[uuid.UUID][]*NodeVersion)
	for i := range nodes {
		if nodes[i].node_version_storage_file == nil {
			continue
		}
		fk := *nodes[i].node_version_storage_file
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(storagefile.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "node_version_storage_file" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (nvq *NodeVersionQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := nvq.querySpec()
	if len(nvq.modifiers) > 0 {
		_spec.Modifiers = nvq.modifiers
	}
	_spec.Node.Columns = nvq.ctx.Fields
	if len(nvq.ctx.Fields) > 0 {
		_spec.Unique = nvq.ctx.Unique != nil && *nvq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, nvq.driver, _spec)
}

func (nvq *NodeVersionQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(nodeversion.Table, nodeversion.Columns, sqlgraph.NewFieldSpec(nodeversion.FieldID, field.TypeUUID))
	_spec.From = nvq.sql
	if unique := nvq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if nvq.path != nil {
		_spec.Unique = true
	}
	if fields := nvq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, nodeversion.FieldID)
		for i := range fields {
			if fields[i] != nodeversion.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
		if nvq.withNode != nil {
			_spec.Node.AddColumnOnce(nodeversion.FieldNodeID)
		}
	}
	if ps := nvq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := nvq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := nvq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := nvq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (nvq *NodeVersionQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(nvq.driver.Dialect())
	t1 := builder.Table(nodeversion.Table)
	columns := nvq.ctx.Fields
	if len(columns) == 0 {
		columns = nodeversion.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if nvq.sql != nil {
		selector = nvq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if nvq.ctx.Unique != nil && *nvq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range nvq.modifiers {
		m(selector)
	}
	for _, p := range nvq.predicates {
		p(selector)
	}
	for _, p := range nvq.order {
		p(selector)
	}
	if offset := nvq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := nvq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (nvq *NodeVersionQuery) ForUpdate(opts ...sql.LockOption) *NodeVersionQuery {
	if nvq.driver.Dialect() == dialect.Postgres {
		nvq.Unique(false)
	}
	nvq.modifiers = append(nvq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return nvq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (nvq *NodeVersionQuery) ForShare(opts ...sql.LockOption) *NodeVersionQuery {
	if nvq.driver.Dialect() == dialect.Postgres {
		nvq.Unique(false)
	}
	nvq.modifiers = append(nvq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return nvq
}

// Modify adds a query modifier for attaching custom logic to queries.
func (nvq *NodeVersionQuery) Modify(modifiers ...func(s *sql.Selector)) *NodeVersionSelect {
	nvq.modifiers = append(nvq.modifiers, modifiers...)
	return nvq.Select()
}

// NodeVersionGroupBy is the group-by builder for NodeVersion entities.
type NodeVersionGroupBy struct {
	selector
	build *NodeVersionQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (nvgb *NodeVersionGroupBy) Aggregate(fns ...AggregateFunc) *NodeVersionGroupBy {
	nvgb.fns = append(nvgb.fns, fns...)
	return nvgb
}

// Scan applies the selector query and scans the result into the given value.
func (nvgb *NodeVersionGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, nvgb.build.ctx, "GroupBy")
	if err := nvgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*NodeVersionQuery, *NodeVersionGroupBy](ctx, nvgb.build, nvgb, nvgb.build.inters, v)
}

func (nvgb *NodeVersionGroupBy) sqlScan(ctx context.Context, root *NodeVersionQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(nvgb.fns))
	for _, fn := range nvgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*nvgb.flds)+len(nvgb.fns))
		for _, f := range *nvgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*nvgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := nvgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// NodeVersionSelect is the builder for selecting fields of NodeVersion entities.
type NodeVersionSelect struct {
	*NodeVersionQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (nvs *NodeVersionSelect) Aggregate(fns ...AggregateFunc) *NodeVersionSelect {
	nvs.fns = append(nvs.fns, fns...)
	return nvs
}

// Scan applies the selector query and scans the result into the given value.
func (nvs *NodeVersionSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, nvs.ctx, "Select")
	if err := nvs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*NodeVersionQuery, *NodeVersionSelect](ctx, nvs.NodeVersionQuery, nvs, nvs.inters, v)
}

func (nvs *NodeVersionSelect) sqlScan(ctx context.Context, root *NodeVersionQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(nvs.fns))
	for _, fn := range nvs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*nvs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := nvs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (nvs *NodeVersionSelect) Modify(modifiers ...func(s *sql.Selector)) *NodeVersionSelect {
	nvs.modifiers = append(nvs.modifiers, modifiers...)
	return nvs
}

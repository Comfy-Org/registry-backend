// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"
	"registry-backend/ent/node"
	"registry-backend/ent/nodereview"
	"registry-backend/ent/predicate"
	"registry-backend/ent/user"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// NodeReviewQuery is the builder for querying NodeReview entities.
type NodeReviewQuery struct {
	config
	ctx        *QueryContext
	order      []nodereview.OrderOption
	inters     []Interceptor
	predicates []predicate.NodeReview
	withUser   *UserQuery
	withNode   *NodeQuery
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the NodeReviewQuery builder.
func (nrq *NodeReviewQuery) Where(ps ...predicate.NodeReview) *NodeReviewQuery {
	nrq.predicates = append(nrq.predicates, ps...)
	return nrq
}

// Limit the number of records to be returned by this query.
func (nrq *NodeReviewQuery) Limit(limit int) *NodeReviewQuery {
	nrq.ctx.Limit = &limit
	return nrq
}

// Offset to start from.
func (nrq *NodeReviewQuery) Offset(offset int) *NodeReviewQuery {
	nrq.ctx.Offset = &offset
	return nrq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (nrq *NodeReviewQuery) Unique(unique bool) *NodeReviewQuery {
	nrq.ctx.Unique = &unique
	return nrq
}

// Order specifies how the records should be ordered.
func (nrq *NodeReviewQuery) Order(o ...nodereview.OrderOption) *NodeReviewQuery {
	nrq.order = append(nrq.order, o...)
	return nrq
}

// QueryUser chains the current query on the "user" edge.
func (nrq *NodeReviewQuery) QueryUser() *UserQuery {
	query := (&UserClient{config: nrq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := nrq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := nrq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(nodereview.Table, nodereview.FieldID, selector),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, nodereview.UserTable, nodereview.UserColumn),
		)
		fromU = sqlgraph.SetNeighbors(nrq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryNode chains the current query on the "node" edge.
func (nrq *NodeReviewQuery) QueryNode() *NodeQuery {
	query := (&NodeClient{config: nrq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := nrq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := nrq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(nodereview.Table, nodereview.FieldID, selector),
			sqlgraph.To(node.Table, node.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, nodereview.NodeTable, nodereview.NodeColumn),
		)
		fromU = sqlgraph.SetNeighbors(nrq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first NodeReview entity from the query.
// Returns a *NotFoundError when no NodeReview was found.
func (nrq *NodeReviewQuery) First(ctx context.Context) (*NodeReview, error) {
	nodes, err := nrq.Limit(1).All(setContextOp(ctx, nrq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{nodereview.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (nrq *NodeReviewQuery) FirstX(ctx context.Context) *NodeReview {
	node, err := nrq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first NodeReview ID from the query.
// Returns a *NotFoundError when no NodeReview ID was found.
func (nrq *NodeReviewQuery) FirstID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = nrq.Limit(1).IDs(setContextOp(ctx, nrq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{nodereview.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (nrq *NodeReviewQuery) FirstIDX(ctx context.Context) uuid.UUID {
	id, err := nrq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single NodeReview entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one NodeReview entity is found.
// Returns a *NotFoundError when no NodeReview entities are found.
func (nrq *NodeReviewQuery) Only(ctx context.Context) (*NodeReview, error) {
	nodes, err := nrq.Limit(2).All(setContextOp(ctx, nrq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{nodereview.Label}
	default:
		return nil, &NotSingularError{nodereview.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (nrq *NodeReviewQuery) OnlyX(ctx context.Context) *NodeReview {
	node, err := nrq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only NodeReview ID in the query.
// Returns a *NotSingularError when more than one NodeReview ID is found.
// Returns a *NotFoundError when no entities are found.
func (nrq *NodeReviewQuery) OnlyID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = nrq.Limit(2).IDs(setContextOp(ctx, nrq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{nodereview.Label}
	default:
		err = &NotSingularError{nodereview.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (nrq *NodeReviewQuery) OnlyIDX(ctx context.Context) uuid.UUID {
	id, err := nrq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of NodeReviews.
func (nrq *NodeReviewQuery) All(ctx context.Context) ([]*NodeReview, error) {
	ctx = setContextOp(ctx, nrq.ctx, ent.OpQueryAll)
	if err := nrq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*NodeReview, *NodeReviewQuery]()
	return withInterceptors[[]*NodeReview](ctx, nrq, qr, nrq.inters)
}

// AllX is like All, but panics if an error occurs.
func (nrq *NodeReviewQuery) AllX(ctx context.Context) []*NodeReview {
	nodes, err := nrq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of NodeReview IDs.
func (nrq *NodeReviewQuery) IDs(ctx context.Context) (ids []uuid.UUID, err error) {
	if nrq.ctx.Unique == nil && nrq.path != nil {
		nrq.Unique(true)
	}
	ctx = setContextOp(ctx, nrq.ctx, ent.OpQueryIDs)
	if err = nrq.Select(nodereview.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (nrq *NodeReviewQuery) IDsX(ctx context.Context) []uuid.UUID {
	ids, err := nrq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (nrq *NodeReviewQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, nrq.ctx, ent.OpQueryCount)
	if err := nrq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, nrq, querierCount[*NodeReviewQuery](), nrq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (nrq *NodeReviewQuery) CountX(ctx context.Context) int {
	count, err := nrq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (nrq *NodeReviewQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, nrq.ctx, ent.OpQueryExist)
	switch _, err := nrq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (nrq *NodeReviewQuery) ExistX(ctx context.Context) bool {
	exist, err := nrq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the NodeReviewQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (nrq *NodeReviewQuery) Clone() *NodeReviewQuery {
	if nrq == nil {
		return nil
	}
	return &NodeReviewQuery{
		config:     nrq.config,
		ctx:        nrq.ctx.Clone(),
		order:      append([]nodereview.OrderOption{}, nrq.order...),
		inters:     append([]Interceptor{}, nrq.inters...),
		predicates: append([]predicate.NodeReview{}, nrq.predicates...),
		withUser:   nrq.withUser.Clone(),
		withNode:   nrq.withNode.Clone(),
		// clone intermediate query.
		sql:       nrq.sql.Clone(),
		path:      nrq.path,
		modifiers: append([]func(*sql.Selector){}, nrq.modifiers...),
	}
}

// WithUser tells the query-builder to eager-load the nodes that are connected to
// the "user" edge. The optional arguments are used to configure the query builder of the edge.
func (nrq *NodeReviewQuery) WithUser(opts ...func(*UserQuery)) *NodeReviewQuery {
	query := (&UserClient{config: nrq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	nrq.withUser = query
	return nrq
}

// WithNode tells the query-builder to eager-load the nodes that are connected to
// the "node" edge. The optional arguments are used to configure the query builder of the edge.
func (nrq *NodeReviewQuery) WithNode(opts ...func(*NodeQuery)) *NodeReviewQuery {
	query := (&NodeClient{config: nrq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	nrq.withNode = query
	return nrq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		NodeID string `json:"node_id,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.NodeReview.Query().
//		GroupBy(nodereview.FieldNodeID).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (nrq *NodeReviewQuery) GroupBy(field string, fields ...string) *NodeReviewGroupBy {
	nrq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &NodeReviewGroupBy{build: nrq}
	grbuild.flds = &nrq.ctx.Fields
	grbuild.label = nodereview.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		NodeID string `json:"node_id,omitempty"`
//	}
//
//	client.NodeReview.Query().
//		Select(nodereview.FieldNodeID).
//		Scan(ctx, &v)
func (nrq *NodeReviewQuery) Select(fields ...string) *NodeReviewSelect {
	nrq.ctx.Fields = append(nrq.ctx.Fields, fields...)
	sbuild := &NodeReviewSelect{NodeReviewQuery: nrq}
	sbuild.label = nodereview.Label
	sbuild.flds, sbuild.scan = &nrq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a NodeReviewSelect configured with the given aggregations.
func (nrq *NodeReviewQuery) Aggregate(fns ...AggregateFunc) *NodeReviewSelect {
	return nrq.Select().Aggregate(fns...)
}

func (nrq *NodeReviewQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range nrq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, nrq); err != nil {
				return err
			}
		}
	}
	for _, f := range nrq.ctx.Fields {
		if !nodereview.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if nrq.path != nil {
		prev, err := nrq.path(ctx)
		if err != nil {
			return err
		}
		nrq.sql = prev
	}
	return nil
}

func (nrq *NodeReviewQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*NodeReview, error) {
	var (
		nodes       = []*NodeReview{}
		_spec       = nrq.querySpec()
		loadedTypes = [2]bool{
			nrq.withUser != nil,
			nrq.withNode != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*NodeReview).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &NodeReview{config: nrq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(nrq.modifiers) > 0 {
		_spec.Modifiers = nrq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, nrq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := nrq.withUser; query != nil {
		if err := nrq.loadUser(ctx, query, nodes, nil,
			func(n *NodeReview, e *User) { n.Edges.User = e }); err != nil {
			return nil, err
		}
	}
	if query := nrq.withNode; query != nil {
		if err := nrq.loadNode(ctx, query, nodes, nil,
			func(n *NodeReview, e *Node) { n.Edges.Node = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (nrq *NodeReviewQuery) loadUser(ctx context.Context, query *UserQuery, nodes []*NodeReview, init func(*NodeReview), assign func(*NodeReview, *User)) error {
	ids := make([]string, 0, len(nodes))
	nodeids := make(map[string][]*NodeReview)
	for i := range nodes {
		fk := nodes[i].UserID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(user.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "user_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (nrq *NodeReviewQuery) loadNode(ctx context.Context, query *NodeQuery, nodes []*NodeReview, init func(*NodeReview), assign func(*NodeReview, *Node)) error {
	ids := make([]string, 0, len(nodes))
	nodeids := make(map[string][]*NodeReview)
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

func (nrq *NodeReviewQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := nrq.querySpec()
	if len(nrq.modifiers) > 0 {
		_spec.Modifiers = nrq.modifiers
	}
	_spec.Node.Columns = nrq.ctx.Fields
	if len(nrq.ctx.Fields) > 0 {
		_spec.Unique = nrq.ctx.Unique != nil && *nrq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, nrq.driver, _spec)
}

func (nrq *NodeReviewQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(nodereview.Table, nodereview.Columns, sqlgraph.NewFieldSpec(nodereview.FieldID, field.TypeUUID))
	_spec.From = nrq.sql
	if unique := nrq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if nrq.path != nil {
		_spec.Unique = true
	}
	if fields := nrq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, nodereview.FieldID)
		for i := range fields {
			if fields[i] != nodereview.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
		if nrq.withUser != nil {
			_spec.Node.AddColumnOnce(nodereview.FieldUserID)
		}
		if nrq.withNode != nil {
			_spec.Node.AddColumnOnce(nodereview.FieldNodeID)
		}
	}
	if ps := nrq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := nrq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := nrq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := nrq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (nrq *NodeReviewQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(nrq.driver.Dialect())
	t1 := builder.Table(nodereview.Table)
	columns := nrq.ctx.Fields
	if len(columns) == 0 {
		columns = nodereview.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if nrq.sql != nil {
		selector = nrq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if nrq.ctx.Unique != nil && *nrq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range nrq.modifiers {
		m(selector)
	}
	for _, p := range nrq.predicates {
		p(selector)
	}
	for _, p := range nrq.order {
		p(selector)
	}
	if offset := nrq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := nrq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (nrq *NodeReviewQuery) ForUpdate(opts ...sql.LockOption) *NodeReviewQuery {
	if nrq.driver.Dialect() == dialect.Postgres {
		nrq.Unique(false)
	}
	nrq.modifiers = append(nrq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return nrq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (nrq *NodeReviewQuery) ForShare(opts ...sql.LockOption) *NodeReviewQuery {
	if nrq.driver.Dialect() == dialect.Postgres {
		nrq.Unique(false)
	}
	nrq.modifiers = append(nrq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return nrq
}

// Modify adds a query modifier for attaching custom logic to queries.
func (nrq *NodeReviewQuery) Modify(modifiers ...func(s *sql.Selector)) *NodeReviewSelect {
	nrq.modifiers = append(nrq.modifiers, modifiers...)
	return nrq.Select()
}

// NodeReviewGroupBy is the group-by builder for NodeReview entities.
type NodeReviewGroupBy struct {
	selector
	build *NodeReviewQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (nrgb *NodeReviewGroupBy) Aggregate(fns ...AggregateFunc) *NodeReviewGroupBy {
	nrgb.fns = append(nrgb.fns, fns...)
	return nrgb
}

// Scan applies the selector query and scans the result into the given value.
func (nrgb *NodeReviewGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, nrgb.build.ctx, ent.OpQueryGroupBy)
	if err := nrgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*NodeReviewQuery, *NodeReviewGroupBy](ctx, nrgb.build, nrgb, nrgb.build.inters, v)
}

func (nrgb *NodeReviewGroupBy) sqlScan(ctx context.Context, root *NodeReviewQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(nrgb.fns))
	for _, fn := range nrgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*nrgb.flds)+len(nrgb.fns))
		for _, f := range *nrgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*nrgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := nrgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// NodeReviewSelect is the builder for selecting fields of NodeReview entities.
type NodeReviewSelect struct {
	*NodeReviewQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (nrs *NodeReviewSelect) Aggregate(fns ...AggregateFunc) *NodeReviewSelect {
	nrs.fns = append(nrs.fns, fns...)
	return nrs
}

// Scan applies the selector query and scans the result into the given value.
func (nrs *NodeReviewSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, nrs.ctx, ent.OpQuerySelect)
	if err := nrs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*NodeReviewQuery, *NodeReviewSelect](ctx, nrs.NodeReviewQuery, nrs, nrs.inters, v)
}

func (nrs *NodeReviewSelect) sqlScan(ctx context.Context, root *NodeReviewQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(nrs.fns))
	for _, fn := range nrs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*nrs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := nrs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (nrs *NodeReviewSelect) Modify(modifiers ...func(s *sql.Selector)) *NodeReviewSelect {
	nrs.modifiers = append(nrs.modifiers, modifiers...)
	return nrs
}

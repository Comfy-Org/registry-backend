// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"
	"registry-backend/ent/personalaccesstoken"
	"registry-backend/ent/predicate"
	"registry-backend/ent/publisher"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// PersonalAccessTokenQuery is the builder for querying PersonalAccessToken entities.
type PersonalAccessTokenQuery struct {
	config
	ctx           *QueryContext
	order         []personalaccesstoken.OrderOption
	inters        []Interceptor
	predicates    []predicate.PersonalAccessToken
	withPublisher *PublisherQuery
	modifiers     []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the PersonalAccessTokenQuery builder.
func (patq *PersonalAccessTokenQuery) Where(ps ...predicate.PersonalAccessToken) *PersonalAccessTokenQuery {
	patq.predicates = append(patq.predicates, ps...)
	return patq
}

// Limit the number of records to be returned by this query.
func (patq *PersonalAccessTokenQuery) Limit(limit int) *PersonalAccessTokenQuery {
	patq.ctx.Limit = &limit
	return patq
}

// Offset to start from.
func (patq *PersonalAccessTokenQuery) Offset(offset int) *PersonalAccessTokenQuery {
	patq.ctx.Offset = &offset
	return patq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (patq *PersonalAccessTokenQuery) Unique(unique bool) *PersonalAccessTokenQuery {
	patq.ctx.Unique = &unique
	return patq
}

// Order specifies how the records should be ordered.
func (patq *PersonalAccessTokenQuery) Order(o ...personalaccesstoken.OrderOption) *PersonalAccessTokenQuery {
	patq.order = append(patq.order, o...)
	return patq
}

// QueryPublisher chains the current query on the "publisher" edge.
func (patq *PersonalAccessTokenQuery) QueryPublisher() *PublisherQuery {
	query := (&PublisherClient{config: patq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := patq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := patq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(personalaccesstoken.Table, personalaccesstoken.FieldID, selector),
			sqlgraph.To(publisher.Table, publisher.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, personalaccesstoken.PublisherTable, personalaccesstoken.PublisherColumn),
		)
		fromU = sqlgraph.SetNeighbors(patq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first PersonalAccessToken entity from the query.
// Returns a *NotFoundError when no PersonalAccessToken was found.
func (patq *PersonalAccessTokenQuery) First(ctx context.Context) (*PersonalAccessToken, error) {
	nodes, err := patq.Limit(1).All(setContextOp(ctx, patq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{personalaccesstoken.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (patq *PersonalAccessTokenQuery) FirstX(ctx context.Context) *PersonalAccessToken {
	node, err := patq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first PersonalAccessToken ID from the query.
// Returns a *NotFoundError when no PersonalAccessToken ID was found.
func (patq *PersonalAccessTokenQuery) FirstID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = patq.Limit(1).IDs(setContextOp(ctx, patq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{personalaccesstoken.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (patq *PersonalAccessTokenQuery) FirstIDX(ctx context.Context) uuid.UUID {
	id, err := patq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single PersonalAccessToken entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one PersonalAccessToken entity is found.
// Returns a *NotFoundError when no PersonalAccessToken entities are found.
func (patq *PersonalAccessTokenQuery) Only(ctx context.Context) (*PersonalAccessToken, error) {
	nodes, err := patq.Limit(2).All(setContextOp(ctx, patq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{personalaccesstoken.Label}
	default:
		return nil, &NotSingularError{personalaccesstoken.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (patq *PersonalAccessTokenQuery) OnlyX(ctx context.Context) *PersonalAccessToken {
	node, err := patq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only PersonalAccessToken ID in the query.
// Returns a *NotSingularError when more than one PersonalAccessToken ID is found.
// Returns a *NotFoundError when no entities are found.
func (patq *PersonalAccessTokenQuery) OnlyID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = patq.Limit(2).IDs(setContextOp(ctx, patq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{personalaccesstoken.Label}
	default:
		err = &NotSingularError{personalaccesstoken.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (patq *PersonalAccessTokenQuery) OnlyIDX(ctx context.Context) uuid.UUID {
	id, err := patq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of PersonalAccessTokens.
func (patq *PersonalAccessTokenQuery) All(ctx context.Context) ([]*PersonalAccessToken, error) {
	ctx = setContextOp(ctx, patq.ctx, "All")
	if err := patq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*PersonalAccessToken, *PersonalAccessTokenQuery]()
	return withInterceptors[[]*PersonalAccessToken](ctx, patq, qr, patq.inters)
}

// AllX is like All, but panics if an error occurs.
func (patq *PersonalAccessTokenQuery) AllX(ctx context.Context) []*PersonalAccessToken {
	nodes, err := patq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of PersonalAccessToken IDs.
func (patq *PersonalAccessTokenQuery) IDs(ctx context.Context) (ids []uuid.UUID, err error) {
	if patq.ctx.Unique == nil && patq.path != nil {
		patq.Unique(true)
	}
	ctx = setContextOp(ctx, patq.ctx, "IDs")
	if err = patq.Select(personalaccesstoken.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (patq *PersonalAccessTokenQuery) IDsX(ctx context.Context) []uuid.UUID {
	ids, err := patq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (patq *PersonalAccessTokenQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, patq.ctx, "Count")
	if err := patq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, patq, querierCount[*PersonalAccessTokenQuery](), patq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (patq *PersonalAccessTokenQuery) CountX(ctx context.Context) int {
	count, err := patq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (patq *PersonalAccessTokenQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, patq.ctx, "Exist")
	switch _, err := patq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (patq *PersonalAccessTokenQuery) ExistX(ctx context.Context) bool {
	exist, err := patq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the PersonalAccessTokenQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (patq *PersonalAccessTokenQuery) Clone() *PersonalAccessTokenQuery {
	if patq == nil {
		return nil
	}
	return &PersonalAccessTokenQuery{
		config:        patq.config,
		ctx:           patq.ctx.Clone(),
		order:         append([]personalaccesstoken.OrderOption{}, patq.order...),
		inters:        append([]Interceptor{}, patq.inters...),
		predicates:    append([]predicate.PersonalAccessToken{}, patq.predicates...),
		withPublisher: patq.withPublisher.Clone(),
		// clone intermediate query.
		sql:  patq.sql.Clone(),
		path: patq.path,
	}
}

// WithPublisher tells the query-builder to eager-load the nodes that are connected to
// the "publisher" edge. The optional arguments are used to configure the query builder of the edge.
func (patq *PersonalAccessTokenQuery) WithPublisher(opts ...func(*PublisherQuery)) *PersonalAccessTokenQuery {
	query := (&PublisherClient{config: patq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	patq.withPublisher = query
	return patq
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
//	client.PersonalAccessToken.Query().
//		GroupBy(personalaccesstoken.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (patq *PersonalAccessTokenQuery) GroupBy(field string, fields ...string) *PersonalAccessTokenGroupBy {
	patq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &PersonalAccessTokenGroupBy{build: patq}
	grbuild.flds = &patq.ctx.Fields
	grbuild.label = personalaccesstoken.Label
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
//	client.PersonalAccessToken.Query().
//		Select(personalaccesstoken.FieldCreateTime).
//		Scan(ctx, &v)
func (patq *PersonalAccessTokenQuery) Select(fields ...string) *PersonalAccessTokenSelect {
	patq.ctx.Fields = append(patq.ctx.Fields, fields...)
	sbuild := &PersonalAccessTokenSelect{PersonalAccessTokenQuery: patq}
	sbuild.label = personalaccesstoken.Label
	sbuild.flds, sbuild.scan = &patq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a PersonalAccessTokenSelect configured with the given aggregations.
func (patq *PersonalAccessTokenQuery) Aggregate(fns ...AggregateFunc) *PersonalAccessTokenSelect {
	return patq.Select().Aggregate(fns...)
}

func (patq *PersonalAccessTokenQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range patq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, patq); err != nil {
				return err
			}
		}
	}
	for _, f := range patq.ctx.Fields {
		if !personalaccesstoken.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if patq.path != nil {
		prev, err := patq.path(ctx)
		if err != nil {
			return err
		}
		patq.sql = prev
	}
	return nil
}

func (patq *PersonalAccessTokenQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*PersonalAccessToken, error) {
	var (
		nodes       = []*PersonalAccessToken{}
		_spec       = patq.querySpec()
		loadedTypes = [1]bool{
			patq.withPublisher != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*PersonalAccessToken).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &PersonalAccessToken{config: patq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(patq.modifiers) > 0 {
		_spec.Modifiers = patq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, patq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := patq.withPublisher; query != nil {
		if err := patq.loadPublisher(ctx, query, nodes, nil,
			func(n *PersonalAccessToken, e *Publisher) { n.Edges.Publisher = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (patq *PersonalAccessTokenQuery) loadPublisher(ctx context.Context, query *PublisherQuery, nodes []*PersonalAccessToken, init func(*PersonalAccessToken), assign func(*PersonalAccessToken, *Publisher)) error {
	ids := make([]string, 0, len(nodes))
	nodeids := make(map[string][]*PersonalAccessToken)
	for i := range nodes {
		fk := nodes[i].PublisherID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(publisher.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "publisher_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (patq *PersonalAccessTokenQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := patq.querySpec()
	if len(patq.modifiers) > 0 {
		_spec.Modifiers = patq.modifiers
	}
	_spec.Node.Columns = patq.ctx.Fields
	if len(patq.ctx.Fields) > 0 {
		_spec.Unique = patq.ctx.Unique != nil && *patq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, patq.driver, _spec)
}

func (patq *PersonalAccessTokenQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(personalaccesstoken.Table, personalaccesstoken.Columns, sqlgraph.NewFieldSpec(personalaccesstoken.FieldID, field.TypeUUID))
	_spec.From = patq.sql
	if unique := patq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if patq.path != nil {
		_spec.Unique = true
	}
	if fields := patq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, personalaccesstoken.FieldID)
		for i := range fields {
			if fields[i] != personalaccesstoken.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
		if patq.withPublisher != nil {
			_spec.Node.AddColumnOnce(personalaccesstoken.FieldPublisherID)
		}
	}
	if ps := patq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := patq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := patq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := patq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (patq *PersonalAccessTokenQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(patq.driver.Dialect())
	t1 := builder.Table(personalaccesstoken.Table)
	columns := patq.ctx.Fields
	if len(columns) == 0 {
		columns = personalaccesstoken.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if patq.sql != nil {
		selector = patq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if patq.ctx.Unique != nil && *patq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range patq.modifiers {
		m(selector)
	}
	for _, p := range patq.predicates {
		p(selector)
	}
	for _, p := range patq.order {
		p(selector)
	}
	if offset := patq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := patq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (patq *PersonalAccessTokenQuery) ForUpdate(opts ...sql.LockOption) *PersonalAccessTokenQuery {
	if patq.driver.Dialect() == dialect.Postgres {
		patq.Unique(false)
	}
	patq.modifiers = append(patq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return patq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (patq *PersonalAccessTokenQuery) ForShare(opts ...sql.LockOption) *PersonalAccessTokenQuery {
	if patq.driver.Dialect() == dialect.Postgres {
		patq.Unique(false)
	}
	patq.modifiers = append(patq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return patq
}

// PersonalAccessTokenGroupBy is the group-by builder for PersonalAccessToken entities.
type PersonalAccessTokenGroupBy struct {
	selector
	build *PersonalAccessTokenQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (patgb *PersonalAccessTokenGroupBy) Aggregate(fns ...AggregateFunc) *PersonalAccessTokenGroupBy {
	patgb.fns = append(patgb.fns, fns...)
	return patgb
}

// Scan applies the selector query and scans the result into the given value.
func (patgb *PersonalAccessTokenGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, patgb.build.ctx, "GroupBy")
	if err := patgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*PersonalAccessTokenQuery, *PersonalAccessTokenGroupBy](ctx, patgb.build, patgb, patgb.build.inters, v)
}

func (patgb *PersonalAccessTokenGroupBy) sqlScan(ctx context.Context, root *PersonalAccessTokenQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(patgb.fns))
	for _, fn := range patgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*patgb.flds)+len(patgb.fns))
		for _, f := range *patgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*patgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := patgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// PersonalAccessTokenSelect is the builder for selecting fields of PersonalAccessToken entities.
type PersonalAccessTokenSelect struct {
	*PersonalAccessTokenQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (pats *PersonalAccessTokenSelect) Aggregate(fns ...AggregateFunc) *PersonalAccessTokenSelect {
	pats.fns = append(pats.fns, fns...)
	return pats
}

// Scan applies the selector query and scans the result into the given value.
func (pats *PersonalAccessTokenSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, pats.ctx, "Select")
	if err := pats.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*PersonalAccessTokenQuery, *PersonalAccessTokenSelect](ctx, pats.PersonalAccessTokenQuery, pats, pats.inters, v)
}

func (pats *PersonalAccessTokenSelect) sqlScan(ctx context.Context, root *PersonalAccessTokenQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(pats.fns))
	for _, fn := range pats.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*pats.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := pats.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
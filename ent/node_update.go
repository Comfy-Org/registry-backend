// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"registry-backend/ent/node"
	"registry-backend/ent/nodereview"
	"registry-backend/ent/nodeversion"
	"registry-backend/ent/predicate"
	"registry-backend/ent/publisher"
	"registry-backend/ent/schema"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/dialect/sql/sqljson"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// NodeUpdate is the builder for updating Node entities.
type NodeUpdate struct {
	config
	hooks    []Hook
	mutation *NodeMutation
}

// Where appends a list predicates to the NodeUpdate builder.
func (nu *NodeUpdate) Where(ps ...predicate.Node) *NodeUpdate {
	nu.mutation.Where(ps...)
	return nu
}

// SetUpdateTime sets the "update_time" field.
func (nu *NodeUpdate) SetUpdateTime(t time.Time) *NodeUpdate {
	nu.mutation.SetUpdateTime(t)
	return nu
}

// SetPublisherID sets the "publisher_id" field.
func (nu *NodeUpdate) SetPublisherID(s string) *NodeUpdate {
	nu.mutation.SetPublisherID(s)
	return nu
}

// SetNillablePublisherID sets the "publisher_id" field if the given value is not nil.
func (nu *NodeUpdate) SetNillablePublisherID(s *string) *NodeUpdate {
	if s != nil {
		nu.SetPublisherID(*s)
	}
	return nu
}

// SetName sets the "name" field.
func (nu *NodeUpdate) SetName(s string) *NodeUpdate {
	nu.mutation.SetName(s)
	return nu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (nu *NodeUpdate) SetNillableName(s *string) *NodeUpdate {
	if s != nil {
		nu.SetName(*s)
	}
	return nu
}

// SetDescription sets the "description" field.
func (nu *NodeUpdate) SetDescription(s string) *NodeUpdate {
	nu.mutation.SetDescription(s)
	return nu
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (nu *NodeUpdate) SetNillableDescription(s *string) *NodeUpdate {
	if s != nil {
		nu.SetDescription(*s)
	}
	return nu
}

// ClearDescription clears the value of the "description" field.
func (nu *NodeUpdate) ClearDescription() *NodeUpdate {
	nu.mutation.ClearDescription()
	return nu
}

// SetAuthor sets the "author" field.
func (nu *NodeUpdate) SetAuthor(s string) *NodeUpdate {
	nu.mutation.SetAuthor(s)
	return nu
}

// SetNillableAuthor sets the "author" field if the given value is not nil.
func (nu *NodeUpdate) SetNillableAuthor(s *string) *NodeUpdate {
	if s != nil {
		nu.SetAuthor(*s)
	}
	return nu
}

// ClearAuthor clears the value of the "author" field.
func (nu *NodeUpdate) ClearAuthor() *NodeUpdate {
	nu.mutation.ClearAuthor()
	return nu
}

// SetLicense sets the "license" field.
func (nu *NodeUpdate) SetLicense(s string) *NodeUpdate {
	nu.mutation.SetLicense(s)
	return nu
}

// SetNillableLicense sets the "license" field if the given value is not nil.
func (nu *NodeUpdate) SetNillableLicense(s *string) *NodeUpdate {
	if s != nil {
		nu.SetLicense(*s)
	}
	return nu
}

// SetRepositoryURL sets the "repository_url" field.
func (nu *NodeUpdate) SetRepositoryURL(s string) *NodeUpdate {
	nu.mutation.SetRepositoryURL(s)
	return nu
}

// SetNillableRepositoryURL sets the "repository_url" field if the given value is not nil.
func (nu *NodeUpdate) SetNillableRepositoryURL(s *string) *NodeUpdate {
	if s != nil {
		nu.SetRepositoryURL(*s)
	}
	return nu
}

// SetIconURL sets the "icon_url" field.
func (nu *NodeUpdate) SetIconURL(s string) *NodeUpdate {
	nu.mutation.SetIconURL(s)
	return nu
}

// SetNillableIconURL sets the "icon_url" field if the given value is not nil.
func (nu *NodeUpdate) SetNillableIconURL(s *string) *NodeUpdate {
	if s != nil {
		nu.SetIconURL(*s)
	}
	return nu
}

// ClearIconURL clears the value of the "icon_url" field.
func (nu *NodeUpdate) ClearIconURL() *NodeUpdate {
	nu.mutation.ClearIconURL()
	return nu
}

// SetTags sets the "tags" field.
func (nu *NodeUpdate) SetTags(s []string) *NodeUpdate {
	nu.mutation.SetTags(s)
	return nu
}

// AppendTags appends s to the "tags" field.
func (nu *NodeUpdate) AppendTags(s []string) *NodeUpdate {
	nu.mutation.AppendTags(s)
	return nu
}

// SetTotalInstall sets the "total_install" field.
func (nu *NodeUpdate) SetTotalInstall(i int64) *NodeUpdate {
	nu.mutation.ResetTotalInstall()
	nu.mutation.SetTotalInstall(i)
	return nu
}

// SetNillableTotalInstall sets the "total_install" field if the given value is not nil.
func (nu *NodeUpdate) SetNillableTotalInstall(i *int64) *NodeUpdate {
	if i != nil {
		nu.SetTotalInstall(*i)
	}
	return nu
}

// AddTotalInstall adds i to the "total_install" field.
func (nu *NodeUpdate) AddTotalInstall(i int64) *NodeUpdate {
	nu.mutation.AddTotalInstall(i)
	return nu
}

// SetTotalStar sets the "total_star" field.
func (nu *NodeUpdate) SetTotalStar(i int64) *NodeUpdate {
	nu.mutation.ResetTotalStar()
	nu.mutation.SetTotalStar(i)
	return nu
}

// SetNillableTotalStar sets the "total_star" field if the given value is not nil.
func (nu *NodeUpdate) SetNillableTotalStar(i *int64) *NodeUpdate {
	if i != nil {
		nu.SetTotalStar(*i)
	}
	return nu
}

// AddTotalStar adds i to the "total_star" field.
func (nu *NodeUpdate) AddTotalStar(i int64) *NodeUpdate {
	nu.mutation.AddTotalStar(i)
	return nu
}

// SetTotalReview sets the "total_review" field.
func (nu *NodeUpdate) SetTotalReview(i int64) *NodeUpdate {
	nu.mutation.ResetTotalReview()
	nu.mutation.SetTotalReview(i)
	return nu
}

// SetNillableTotalReview sets the "total_review" field if the given value is not nil.
func (nu *NodeUpdate) SetNillableTotalReview(i *int64) *NodeUpdate {
	if i != nil {
		nu.SetTotalReview(*i)
	}
	return nu
}

// AddTotalReview adds i to the "total_review" field.
func (nu *NodeUpdate) AddTotalReview(i int64) *NodeUpdate {
	nu.mutation.AddTotalReview(i)
	return nu
}

// SetStatus sets the "status" field.
func (nu *NodeUpdate) SetStatus(ss schema.NodeStatus) *NodeUpdate {
	nu.mutation.SetStatus(ss)
	return nu
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (nu *NodeUpdate) SetNillableStatus(ss *schema.NodeStatus) *NodeUpdate {
	if ss != nil {
		nu.SetStatus(*ss)
	}
	return nu
}

// SetPublisher sets the "publisher" edge to the Publisher entity.
func (nu *NodeUpdate) SetPublisher(p *Publisher) *NodeUpdate {
	return nu.SetPublisherID(p.ID)
}

// AddVersionIDs adds the "versions" edge to the NodeVersion entity by IDs.
func (nu *NodeUpdate) AddVersionIDs(ids ...uuid.UUID) *NodeUpdate {
	nu.mutation.AddVersionIDs(ids...)
	return nu
}

// AddVersions adds the "versions" edges to the NodeVersion entity.
func (nu *NodeUpdate) AddVersions(n ...*NodeVersion) *NodeUpdate {
	ids := make([]uuid.UUID, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return nu.AddVersionIDs(ids...)
}

// AddReviewIDs adds the "reviews" edge to the NodeReview entity by IDs.
func (nu *NodeUpdate) AddReviewIDs(ids ...uuid.UUID) *NodeUpdate {
	nu.mutation.AddReviewIDs(ids...)
	return nu
}

// AddReviews adds the "reviews" edges to the NodeReview entity.
func (nu *NodeUpdate) AddReviews(n ...*NodeReview) *NodeUpdate {
	ids := make([]uuid.UUID, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return nu.AddReviewIDs(ids...)
}

// Mutation returns the NodeMutation object of the builder.
func (nu *NodeUpdate) Mutation() *NodeMutation {
	return nu.mutation
}

// ClearPublisher clears the "publisher" edge to the Publisher entity.
func (nu *NodeUpdate) ClearPublisher() *NodeUpdate {
	nu.mutation.ClearPublisher()
	return nu
}

// ClearVersions clears all "versions" edges to the NodeVersion entity.
func (nu *NodeUpdate) ClearVersions() *NodeUpdate {
	nu.mutation.ClearVersions()
	return nu
}

// RemoveVersionIDs removes the "versions" edge to NodeVersion entities by IDs.
func (nu *NodeUpdate) RemoveVersionIDs(ids ...uuid.UUID) *NodeUpdate {
	nu.mutation.RemoveVersionIDs(ids...)
	return nu
}

// RemoveVersions removes "versions" edges to NodeVersion entities.
func (nu *NodeUpdate) RemoveVersions(n ...*NodeVersion) *NodeUpdate {
	ids := make([]uuid.UUID, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return nu.RemoveVersionIDs(ids...)
}

// ClearReviews clears all "reviews" edges to the NodeReview entity.
func (nu *NodeUpdate) ClearReviews() *NodeUpdate {
	nu.mutation.ClearReviews()
	return nu
}

// RemoveReviewIDs removes the "reviews" edge to NodeReview entities by IDs.
func (nu *NodeUpdate) RemoveReviewIDs(ids ...uuid.UUID) *NodeUpdate {
	nu.mutation.RemoveReviewIDs(ids...)
	return nu
}

// RemoveReviews removes "reviews" edges to NodeReview entities.
func (nu *NodeUpdate) RemoveReviews(n ...*NodeReview) *NodeUpdate {
	ids := make([]uuid.UUID, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return nu.RemoveReviewIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (nu *NodeUpdate) Save(ctx context.Context) (int, error) {
	nu.defaults()
	return withHooks(ctx, nu.sqlSave, nu.mutation, nu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (nu *NodeUpdate) SaveX(ctx context.Context) int {
	affected, err := nu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (nu *NodeUpdate) Exec(ctx context.Context) error {
	_, err := nu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (nu *NodeUpdate) ExecX(ctx context.Context) {
	if err := nu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (nu *NodeUpdate) defaults() {
	if _, ok := nu.mutation.UpdateTime(); !ok {
		v := node.UpdateDefaultUpdateTime()
		nu.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (nu *NodeUpdate) check() error {
	if v, ok := nu.mutation.Status(); ok {
		if err := node.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Node.status": %w`, err)}
		}
	}
	if _, ok := nu.mutation.PublisherID(); nu.mutation.PublisherCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Node.publisher"`)
	}
	return nil
}

func (nu *NodeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := nu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(node.Table, node.Columns, sqlgraph.NewFieldSpec(node.FieldID, field.TypeString))
	if ps := nu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := nu.mutation.UpdateTime(); ok {
		_spec.SetField(node.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := nu.mutation.Name(); ok {
		_spec.SetField(node.FieldName, field.TypeString, value)
	}
	if value, ok := nu.mutation.Description(); ok {
		_spec.SetField(node.FieldDescription, field.TypeString, value)
	}
	if nu.mutation.DescriptionCleared() {
		_spec.ClearField(node.FieldDescription, field.TypeString)
	}
	if value, ok := nu.mutation.Author(); ok {
		_spec.SetField(node.FieldAuthor, field.TypeString, value)
	}
	if nu.mutation.AuthorCleared() {
		_spec.ClearField(node.FieldAuthor, field.TypeString)
	}
	if value, ok := nu.mutation.License(); ok {
		_spec.SetField(node.FieldLicense, field.TypeString, value)
	}
	if value, ok := nu.mutation.RepositoryURL(); ok {
		_spec.SetField(node.FieldRepositoryURL, field.TypeString, value)
	}
	if value, ok := nu.mutation.IconURL(); ok {
		_spec.SetField(node.FieldIconURL, field.TypeString, value)
	}
	if nu.mutation.IconURLCleared() {
		_spec.ClearField(node.FieldIconURL, field.TypeString)
	}
	if value, ok := nu.mutation.Tags(); ok {
		_spec.SetField(node.FieldTags, field.TypeJSON, value)
	}
	if value, ok := nu.mutation.AppendedTags(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, node.FieldTags, value)
		})
	}
	if value, ok := nu.mutation.TotalInstall(); ok {
		_spec.SetField(node.FieldTotalInstall, field.TypeInt64, value)
	}
	if value, ok := nu.mutation.AddedTotalInstall(); ok {
		_spec.AddField(node.FieldTotalInstall, field.TypeInt64, value)
	}
	if value, ok := nu.mutation.TotalStar(); ok {
		_spec.SetField(node.FieldTotalStar, field.TypeInt64, value)
	}
	if value, ok := nu.mutation.AddedTotalStar(); ok {
		_spec.AddField(node.FieldTotalStar, field.TypeInt64, value)
	}
	if value, ok := nu.mutation.TotalReview(); ok {
		_spec.SetField(node.FieldTotalReview, field.TypeInt64, value)
	}
	if value, ok := nu.mutation.AddedTotalReview(); ok {
		_spec.AddField(node.FieldTotalReview, field.TypeInt64, value)
	}
	if value, ok := nu.mutation.Status(); ok {
		_spec.SetField(node.FieldStatus, field.TypeEnum, value)
	}
	if nu.mutation.PublisherCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   node.PublisherTable,
			Columns: []string{node.PublisherColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(publisher.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := nu.mutation.PublisherIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   node.PublisherTable,
			Columns: []string{node.PublisherColumn},
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
	if nu.mutation.VersionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   node.VersionsTable,
			Columns: []string{node.VersionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(nodeversion.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := nu.mutation.RemovedVersionsIDs(); len(nodes) > 0 && !nu.mutation.VersionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   node.VersionsTable,
			Columns: []string{node.VersionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(nodeversion.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := nu.mutation.VersionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   node.VersionsTable,
			Columns: []string{node.VersionsColumn},
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
	if nu.mutation.ReviewsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   node.ReviewsTable,
			Columns: []string{node.ReviewsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(nodereview.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := nu.mutation.RemovedReviewsIDs(); len(nodes) > 0 && !nu.mutation.ReviewsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   node.ReviewsTable,
			Columns: []string{node.ReviewsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(nodereview.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := nu.mutation.ReviewsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   node.ReviewsTable,
			Columns: []string{node.ReviewsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(nodereview.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, nu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{node.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	nu.mutation.done = true
	return n, nil
}

// NodeUpdateOne is the builder for updating a single Node entity.
type NodeUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *NodeMutation
}

// SetUpdateTime sets the "update_time" field.
func (nuo *NodeUpdateOne) SetUpdateTime(t time.Time) *NodeUpdateOne {
	nuo.mutation.SetUpdateTime(t)
	return nuo
}

// SetPublisherID sets the "publisher_id" field.
func (nuo *NodeUpdateOne) SetPublisherID(s string) *NodeUpdateOne {
	nuo.mutation.SetPublisherID(s)
	return nuo
}

// SetNillablePublisherID sets the "publisher_id" field if the given value is not nil.
func (nuo *NodeUpdateOne) SetNillablePublisherID(s *string) *NodeUpdateOne {
	if s != nil {
		nuo.SetPublisherID(*s)
	}
	return nuo
}

// SetName sets the "name" field.
func (nuo *NodeUpdateOne) SetName(s string) *NodeUpdateOne {
	nuo.mutation.SetName(s)
	return nuo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (nuo *NodeUpdateOne) SetNillableName(s *string) *NodeUpdateOne {
	if s != nil {
		nuo.SetName(*s)
	}
	return nuo
}

// SetDescription sets the "description" field.
func (nuo *NodeUpdateOne) SetDescription(s string) *NodeUpdateOne {
	nuo.mutation.SetDescription(s)
	return nuo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (nuo *NodeUpdateOne) SetNillableDescription(s *string) *NodeUpdateOne {
	if s != nil {
		nuo.SetDescription(*s)
	}
	return nuo
}

// ClearDescription clears the value of the "description" field.
func (nuo *NodeUpdateOne) ClearDescription() *NodeUpdateOne {
	nuo.mutation.ClearDescription()
	return nuo
}

// SetAuthor sets the "author" field.
func (nuo *NodeUpdateOne) SetAuthor(s string) *NodeUpdateOne {
	nuo.mutation.SetAuthor(s)
	return nuo
}

// SetNillableAuthor sets the "author" field if the given value is not nil.
func (nuo *NodeUpdateOne) SetNillableAuthor(s *string) *NodeUpdateOne {
	if s != nil {
		nuo.SetAuthor(*s)
	}
	return nuo
}

// ClearAuthor clears the value of the "author" field.
func (nuo *NodeUpdateOne) ClearAuthor() *NodeUpdateOne {
	nuo.mutation.ClearAuthor()
	return nuo
}

// SetLicense sets the "license" field.
func (nuo *NodeUpdateOne) SetLicense(s string) *NodeUpdateOne {
	nuo.mutation.SetLicense(s)
	return nuo
}

// SetNillableLicense sets the "license" field if the given value is not nil.
func (nuo *NodeUpdateOne) SetNillableLicense(s *string) *NodeUpdateOne {
	if s != nil {
		nuo.SetLicense(*s)
	}
	return nuo
}

// SetRepositoryURL sets the "repository_url" field.
func (nuo *NodeUpdateOne) SetRepositoryURL(s string) *NodeUpdateOne {
	nuo.mutation.SetRepositoryURL(s)
	return nuo
}

// SetNillableRepositoryURL sets the "repository_url" field if the given value is not nil.
func (nuo *NodeUpdateOne) SetNillableRepositoryURL(s *string) *NodeUpdateOne {
	if s != nil {
		nuo.SetRepositoryURL(*s)
	}
	return nuo
}

// SetIconURL sets the "icon_url" field.
func (nuo *NodeUpdateOne) SetIconURL(s string) *NodeUpdateOne {
	nuo.mutation.SetIconURL(s)
	return nuo
}

// SetNillableIconURL sets the "icon_url" field if the given value is not nil.
func (nuo *NodeUpdateOne) SetNillableIconURL(s *string) *NodeUpdateOne {
	if s != nil {
		nuo.SetIconURL(*s)
	}
	return nuo
}

// ClearIconURL clears the value of the "icon_url" field.
func (nuo *NodeUpdateOne) ClearIconURL() *NodeUpdateOne {
	nuo.mutation.ClearIconURL()
	return nuo
}

// SetTags sets the "tags" field.
func (nuo *NodeUpdateOne) SetTags(s []string) *NodeUpdateOne {
	nuo.mutation.SetTags(s)
	return nuo
}

// AppendTags appends s to the "tags" field.
func (nuo *NodeUpdateOne) AppendTags(s []string) *NodeUpdateOne {
	nuo.mutation.AppendTags(s)
	return nuo
}

// SetTotalInstall sets the "total_install" field.
func (nuo *NodeUpdateOne) SetTotalInstall(i int64) *NodeUpdateOne {
	nuo.mutation.ResetTotalInstall()
	nuo.mutation.SetTotalInstall(i)
	return nuo
}

// SetNillableTotalInstall sets the "total_install" field if the given value is not nil.
func (nuo *NodeUpdateOne) SetNillableTotalInstall(i *int64) *NodeUpdateOne {
	if i != nil {
		nuo.SetTotalInstall(*i)
	}
	return nuo
}

// AddTotalInstall adds i to the "total_install" field.
func (nuo *NodeUpdateOne) AddTotalInstall(i int64) *NodeUpdateOne {
	nuo.mutation.AddTotalInstall(i)
	return nuo
}

// SetTotalStar sets the "total_star" field.
func (nuo *NodeUpdateOne) SetTotalStar(i int64) *NodeUpdateOne {
	nuo.mutation.ResetTotalStar()
	nuo.mutation.SetTotalStar(i)
	return nuo
}

// SetNillableTotalStar sets the "total_star" field if the given value is not nil.
func (nuo *NodeUpdateOne) SetNillableTotalStar(i *int64) *NodeUpdateOne {
	if i != nil {
		nuo.SetTotalStar(*i)
	}
	return nuo
}

// AddTotalStar adds i to the "total_star" field.
func (nuo *NodeUpdateOne) AddTotalStar(i int64) *NodeUpdateOne {
	nuo.mutation.AddTotalStar(i)
	return nuo
}

// SetTotalReview sets the "total_review" field.
func (nuo *NodeUpdateOne) SetTotalReview(i int64) *NodeUpdateOne {
	nuo.mutation.ResetTotalReview()
	nuo.mutation.SetTotalReview(i)
	return nuo
}

// SetNillableTotalReview sets the "total_review" field if the given value is not nil.
func (nuo *NodeUpdateOne) SetNillableTotalReview(i *int64) *NodeUpdateOne {
	if i != nil {
		nuo.SetTotalReview(*i)
	}
	return nuo
}

// AddTotalReview adds i to the "total_review" field.
func (nuo *NodeUpdateOne) AddTotalReview(i int64) *NodeUpdateOne {
	nuo.mutation.AddTotalReview(i)
	return nuo
}

// SetStatus sets the "status" field.
func (nuo *NodeUpdateOne) SetStatus(ss schema.NodeStatus) *NodeUpdateOne {
	nuo.mutation.SetStatus(ss)
	return nuo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (nuo *NodeUpdateOne) SetNillableStatus(ss *schema.NodeStatus) *NodeUpdateOne {
	if ss != nil {
		nuo.SetStatus(*ss)
	}
	return nuo
}

// SetPublisher sets the "publisher" edge to the Publisher entity.
func (nuo *NodeUpdateOne) SetPublisher(p *Publisher) *NodeUpdateOne {
	return nuo.SetPublisherID(p.ID)
}

// AddVersionIDs adds the "versions" edge to the NodeVersion entity by IDs.
func (nuo *NodeUpdateOne) AddVersionIDs(ids ...uuid.UUID) *NodeUpdateOne {
	nuo.mutation.AddVersionIDs(ids...)
	return nuo
}

// AddVersions adds the "versions" edges to the NodeVersion entity.
func (nuo *NodeUpdateOne) AddVersions(n ...*NodeVersion) *NodeUpdateOne {
	ids := make([]uuid.UUID, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return nuo.AddVersionIDs(ids...)
}

// AddReviewIDs adds the "reviews" edge to the NodeReview entity by IDs.
func (nuo *NodeUpdateOne) AddReviewIDs(ids ...uuid.UUID) *NodeUpdateOne {
	nuo.mutation.AddReviewIDs(ids...)
	return nuo
}

// AddReviews adds the "reviews" edges to the NodeReview entity.
func (nuo *NodeUpdateOne) AddReviews(n ...*NodeReview) *NodeUpdateOne {
	ids := make([]uuid.UUID, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return nuo.AddReviewIDs(ids...)
}

// Mutation returns the NodeMutation object of the builder.
func (nuo *NodeUpdateOne) Mutation() *NodeMutation {
	return nuo.mutation
}

// ClearPublisher clears the "publisher" edge to the Publisher entity.
func (nuo *NodeUpdateOne) ClearPublisher() *NodeUpdateOne {
	nuo.mutation.ClearPublisher()
	return nuo
}

// ClearVersions clears all "versions" edges to the NodeVersion entity.
func (nuo *NodeUpdateOne) ClearVersions() *NodeUpdateOne {
	nuo.mutation.ClearVersions()
	return nuo
}

// RemoveVersionIDs removes the "versions" edge to NodeVersion entities by IDs.
func (nuo *NodeUpdateOne) RemoveVersionIDs(ids ...uuid.UUID) *NodeUpdateOne {
	nuo.mutation.RemoveVersionIDs(ids...)
	return nuo
}

// RemoveVersions removes "versions" edges to NodeVersion entities.
func (nuo *NodeUpdateOne) RemoveVersions(n ...*NodeVersion) *NodeUpdateOne {
	ids := make([]uuid.UUID, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return nuo.RemoveVersionIDs(ids...)
}

// ClearReviews clears all "reviews" edges to the NodeReview entity.
func (nuo *NodeUpdateOne) ClearReviews() *NodeUpdateOne {
	nuo.mutation.ClearReviews()
	return nuo
}

// RemoveReviewIDs removes the "reviews" edge to NodeReview entities by IDs.
func (nuo *NodeUpdateOne) RemoveReviewIDs(ids ...uuid.UUID) *NodeUpdateOne {
	nuo.mutation.RemoveReviewIDs(ids...)
	return nuo
}

// RemoveReviews removes "reviews" edges to NodeReview entities.
func (nuo *NodeUpdateOne) RemoveReviews(n ...*NodeReview) *NodeUpdateOne {
	ids := make([]uuid.UUID, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return nuo.RemoveReviewIDs(ids...)
}

// Where appends a list predicates to the NodeUpdate builder.
func (nuo *NodeUpdateOne) Where(ps ...predicate.Node) *NodeUpdateOne {
	nuo.mutation.Where(ps...)
	return nuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (nuo *NodeUpdateOne) Select(field string, fields ...string) *NodeUpdateOne {
	nuo.fields = append([]string{field}, fields...)
	return nuo
}

// Save executes the query and returns the updated Node entity.
func (nuo *NodeUpdateOne) Save(ctx context.Context) (*Node, error) {
	nuo.defaults()
	return withHooks(ctx, nuo.sqlSave, nuo.mutation, nuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (nuo *NodeUpdateOne) SaveX(ctx context.Context) *Node {
	node, err := nuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (nuo *NodeUpdateOne) Exec(ctx context.Context) error {
	_, err := nuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (nuo *NodeUpdateOne) ExecX(ctx context.Context) {
	if err := nuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (nuo *NodeUpdateOne) defaults() {
	if _, ok := nuo.mutation.UpdateTime(); !ok {
		v := node.UpdateDefaultUpdateTime()
		nuo.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (nuo *NodeUpdateOne) check() error {
	if v, ok := nuo.mutation.Status(); ok {
		if err := node.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Node.status": %w`, err)}
		}
	}
	if _, ok := nuo.mutation.PublisherID(); nuo.mutation.PublisherCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Node.publisher"`)
	}
	return nil
}

func (nuo *NodeUpdateOne) sqlSave(ctx context.Context) (_node *Node, err error) {
	if err := nuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(node.Table, node.Columns, sqlgraph.NewFieldSpec(node.FieldID, field.TypeString))
	id, ok := nuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Node.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := nuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, node.FieldID)
		for _, f := range fields {
			if !node.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != node.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := nuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := nuo.mutation.UpdateTime(); ok {
		_spec.SetField(node.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := nuo.mutation.Name(); ok {
		_spec.SetField(node.FieldName, field.TypeString, value)
	}
	if value, ok := nuo.mutation.Description(); ok {
		_spec.SetField(node.FieldDescription, field.TypeString, value)
	}
	if nuo.mutation.DescriptionCleared() {
		_spec.ClearField(node.FieldDescription, field.TypeString)
	}
	if value, ok := nuo.mutation.Author(); ok {
		_spec.SetField(node.FieldAuthor, field.TypeString, value)
	}
	if nuo.mutation.AuthorCleared() {
		_spec.ClearField(node.FieldAuthor, field.TypeString)
	}
	if value, ok := nuo.mutation.License(); ok {
		_spec.SetField(node.FieldLicense, field.TypeString, value)
	}
	if value, ok := nuo.mutation.RepositoryURL(); ok {
		_spec.SetField(node.FieldRepositoryURL, field.TypeString, value)
	}
	if value, ok := nuo.mutation.IconURL(); ok {
		_spec.SetField(node.FieldIconURL, field.TypeString, value)
	}
	if nuo.mutation.IconURLCleared() {
		_spec.ClearField(node.FieldIconURL, field.TypeString)
	}
	if value, ok := nuo.mutation.Tags(); ok {
		_spec.SetField(node.FieldTags, field.TypeJSON, value)
	}
	if value, ok := nuo.mutation.AppendedTags(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, node.FieldTags, value)
		})
	}
	if value, ok := nuo.mutation.TotalInstall(); ok {
		_spec.SetField(node.FieldTotalInstall, field.TypeInt64, value)
	}
	if value, ok := nuo.mutation.AddedTotalInstall(); ok {
		_spec.AddField(node.FieldTotalInstall, field.TypeInt64, value)
	}
	if value, ok := nuo.mutation.TotalStar(); ok {
		_spec.SetField(node.FieldTotalStar, field.TypeInt64, value)
	}
	if value, ok := nuo.mutation.AddedTotalStar(); ok {
		_spec.AddField(node.FieldTotalStar, field.TypeInt64, value)
	}
	if value, ok := nuo.mutation.TotalReview(); ok {
		_spec.SetField(node.FieldTotalReview, field.TypeInt64, value)
	}
	if value, ok := nuo.mutation.AddedTotalReview(); ok {
		_spec.AddField(node.FieldTotalReview, field.TypeInt64, value)
	}
	if value, ok := nuo.mutation.Status(); ok {
		_spec.SetField(node.FieldStatus, field.TypeEnum, value)
	}
	if nuo.mutation.PublisherCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   node.PublisherTable,
			Columns: []string{node.PublisherColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(publisher.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := nuo.mutation.PublisherIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   node.PublisherTable,
			Columns: []string{node.PublisherColumn},
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
	if nuo.mutation.VersionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   node.VersionsTable,
			Columns: []string{node.VersionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(nodeversion.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := nuo.mutation.RemovedVersionsIDs(); len(nodes) > 0 && !nuo.mutation.VersionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   node.VersionsTable,
			Columns: []string{node.VersionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(nodeversion.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := nuo.mutation.VersionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   node.VersionsTable,
			Columns: []string{node.VersionsColumn},
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
	if nuo.mutation.ReviewsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   node.ReviewsTable,
			Columns: []string{node.ReviewsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(nodereview.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := nuo.mutation.RemovedReviewsIDs(); len(nodes) > 0 && !nuo.mutation.ReviewsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   node.ReviewsTable,
			Columns: []string{node.ReviewsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(nodereview.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := nuo.mutation.ReviewsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   node.ReviewsTable,
			Columns: []string{node.ReviewsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(nodereview.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Node{config: nuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, nuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{node.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	nuo.mutation.done = true
	return _node, nil
}

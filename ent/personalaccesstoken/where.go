// Code generated by ent, DO NOT EDIT.

package personalaccesstoken

import (
	"registry-backend/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldLTE(FieldID, id))
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldEQ(FieldUpdateTime, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldEQ(FieldName, v))
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldEQ(FieldDescription, v))
}

// PublisherID applies equality check predicate on the "publisher_id" field. It's identical to PublisherIDEQ.
func PublisherID(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldEQ(FieldPublisherID, v))
}

// Token applies equality check predicate on the "token" field. It's identical to TokenEQ.
func Token(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldEQ(FieldToken, v))
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "create_time" field.
func CreateTimeNotIn(vs ...time.Time) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "create_time" field.
func CreateTimeGT(v time.Time) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldLTE(FieldCreateTime, v))
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldEQ(FieldUpdateTime, v))
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldNEQ(FieldUpdateTime, v))
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldIn(FieldUpdateTime, vs...))
}

// UpdateTimeNotIn applies the NotIn predicate on the "update_time" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldNotIn(FieldUpdateTime, vs...))
}

// UpdateTimeGT applies the GT predicate on the "update_time" field.
func UpdateTimeGT(v time.Time) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldGT(FieldUpdateTime, v))
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldGTE(FieldUpdateTime, v))
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldLT(FieldUpdateTime, v))
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldLTE(FieldUpdateTime, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldContainsFold(FieldName, v))
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldEQ(FieldDescription, v))
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldNEQ(FieldDescription, v))
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldIn(FieldDescription, vs...))
}

// DescriptionNotIn applies the NotIn predicate on the "description" field.
func DescriptionNotIn(vs ...string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldNotIn(FieldDescription, vs...))
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldGT(FieldDescription, v))
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldGTE(FieldDescription, v))
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldLT(FieldDescription, v))
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldLTE(FieldDescription, v))
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldContains(FieldDescription, v))
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldHasPrefix(FieldDescription, v))
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldHasSuffix(FieldDescription, v))
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldEqualFold(FieldDescription, v))
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldContainsFold(FieldDescription, v))
}

// PublisherIDEQ applies the EQ predicate on the "publisher_id" field.
func PublisherIDEQ(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldEQ(FieldPublisherID, v))
}

// PublisherIDNEQ applies the NEQ predicate on the "publisher_id" field.
func PublisherIDNEQ(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldNEQ(FieldPublisherID, v))
}

// PublisherIDIn applies the In predicate on the "publisher_id" field.
func PublisherIDIn(vs ...string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldIn(FieldPublisherID, vs...))
}

// PublisherIDNotIn applies the NotIn predicate on the "publisher_id" field.
func PublisherIDNotIn(vs ...string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldNotIn(FieldPublisherID, vs...))
}

// PublisherIDGT applies the GT predicate on the "publisher_id" field.
func PublisherIDGT(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldGT(FieldPublisherID, v))
}

// PublisherIDGTE applies the GTE predicate on the "publisher_id" field.
func PublisherIDGTE(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldGTE(FieldPublisherID, v))
}

// PublisherIDLT applies the LT predicate on the "publisher_id" field.
func PublisherIDLT(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldLT(FieldPublisherID, v))
}

// PublisherIDLTE applies the LTE predicate on the "publisher_id" field.
func PublisherIDLTE(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldLTE(FieldPublisherID, v))
}

// PublisherIDContains applies the Contains predicate on the "publisher_id" field.
func PublisherIDContains(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldContains(FieldPublisherID, v))
}

// PublisherIDHasPrefix applies the HasPrefix predicate on the "publisher_id" field.
func PublisherIDHasPrefix(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldHasPrefix(FieldPublisherID, v))
}

// PublisherIDHasSuffix applies the HasSuffix predicate on the "publisher_id" field.
func PublisherIDHasSuffix(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldHasSuffix(FieldPublisherID, v))
}

// PublisherIDEqualFold applies the EqualFold predicate on the "publisher_id" field.
func PublisherIDEqualFold(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldEqualFold(FieldPublisherID, v))
}

// PublisherIDContainsFold applies the ContainsFold predicate on the "publisher_id" field.
func PublisherIDContainsFold(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldContainsFold(FieldPublisherID, v))
}

// TokenEQ applies the EQ predicate on the "token" field.
func TokenEQ(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldEQ(FieldToken, v))
}

// TokenNEQ applies the NEQ predicate on the "token" field.
func TokenNEQ(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldNEQ(FieldToken, v))
}

// TokenIn applies the In predicate on the "token" field.
func TokenIn(vs ...string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldIn(FieldToken, vs...))
}

// TokenNotIn applies the NotIn predicate on the "token" field.
func TokenNotIn(vs ...string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldNotIn(FieldToken, vs...))
}

// TokenGT applies the GT predicate on the "token" field.
func TokenGT(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldGT(FieldToken, v))
}

// TokenGTE applies the GTE predicate on the "token" field.
func TokenGTE(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldGTE(FieldToken, v))
}

// TokenLT applies the LT predicate on the "token" field.
func TokenLT(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldLT(FieldToken, v))
}

// TokenLTE applies the LTE predicate on the "token" field.
func TokenLTE(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldLTE(FieldToken, v))
}

// TokenContains applies the Contains predicate on the "token" field.
func TokenContains(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldContains(FieldToken, v))
}

// TokenHasPrefix applies the HasPrefix predicate on the "token" field.
func TokenHasPrefix(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldHasPrefix(FieldToken, v))
}

// TokenHasSuffix applies the HasSuffix predicate on the "token" field.
func TokenHasSuffix(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldHasSuffix(FieldToken, v))
}

// TokenEqualFold applies the EqualFold predicate on the "token" field.
func TokenEqualFold(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldEqualFold(FieldToken, v))
}

// TokenContainsFold applies the ContainsFold predicate on the "token" field.
func TokenContainsFold(v string) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.FieldContainsFold(FieldToken, v))
}

// HasPublisher applies the HasEdge predicate on the "publisher" edge.
func HasPublisher() predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, PublisherTable, PublisherColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPublisherWith applies the HasEdge predicate on the "publisher" edge with a given conditions (other predicates).
func HasPublisherWith(preds ...predicate.Publisher) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(func(s *sql.Selector) {
		step := newPublisherStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.PersonalAccessToken) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.PersonalAccessToken) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.PersonalAccessToken) predicate.PersonalAccessToken {
	return predicate.PersonalAccessToken(sql.NotPredicates(p))
}

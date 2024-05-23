// Code generated by ent, DO NOT EDIT.

package node

import (
	"registry-backend/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id string) predicate.Node {
	return predicate.Node(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.Node {
	return predicate.Node(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.Node {
	return predicate.Node(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.Node {
	return predicate.Node(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...string) predicate.Node {
	return predicate.Node(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id string) predicate.Node {
	return predicate.Node(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.Node {
	return predicate.Node(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.Node {
	return predicate.Node(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.Node {
	return predicate.Node(sql.FieldLTE(FieldID, id))
}

// IDEqualFold applies the EqualFold predicate on the ID field.
func IDEqualFold(id string) predicate.Node {
	return predicate.Node(sql.FieldEqualFold(FieldID, id))
}

// IDContainsFold applies the ContainsFold predicate on the ID field.
func IDContainsFold(id string) predicate.Node {
	return predicate.Node(sql.FieldContainsFold(FieldID, id))
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.Node {
	return predicate.Node(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.Node {
	return predicate.Node(sql.FieldEQ(FieldUpdateTime, v))
}

// PublisherID applies equality check predicate on the "publisher_id" field. It's identical to PublisherIDEQ.
func PublisherID(v string) predicate.Node {
	return predicate.Node(sql.FieldEQ(FieldPublisherID, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Node {
	return predicate.Node(sql.FieldEQ(FieldName, v))
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.Node {
	return predicate.Node(sql.FieldEQ(FieldDescription, v))
}

// Author applies equality check predicate on the "author" field. It's identical to AuthorEQ.
func Author(v string) predicate.Node {
	return predicate.Node(sql.FieldEQ(FieldAuthor, v))
}

// License applies equality check predicate on the "license" field. It's identical to LicenseEQ.
func License(v string) predicate.Node {
	return predicate.Node(sql.FieldEQ(FieldLicense, v))
}

// RepositoryURL applies equality check predicate on the "repository_url" field. It's identical to RepositoryURLEQ.
func RepositoryURL(v string) predicate.Node {
	return predicate.Node(sql.FieldEQ(FieldRepositoryURL, v))
}

// IconURL applies equality check predicate on the "icon_url" field. It's identical to IconURLEQ.
func IconURL(v string) predicate.Node {
	return predicate.Node(sql.FieldEQ(FieldIconURL, v))
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.Node {
	return predicate.Node(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.Node {
	return predicate.Node(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.Node {
	return predicate.Node(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "create_time" field.
func CreateTimeNotIn(vs ...time.Time) predicate.Node {
	return predicate.Node(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "create_time" field.
func CreateTimeGT(v time.Time) predicate.Node {
	return predicate.Node(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.Node {
	return predicate.Node(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.Node {
	return predicate.Node(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.Node {
	return predicate.Node(sql.FieldLTE(FieldCreateTime, v))
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.Node {
	return predicate.Node(sql.FieldEQ(FieldUpdateTime, v))
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.Node {
	return predicate.Node(sql.FieldNEQ(FieldUpdateTime, v))
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.Node {
	return predicate.Node(sql.FieldIn(FieldUpdateTime, vs...))
}

// UpdateTimeNotIn applies the NotIn predicate on the "update_time" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.Node {
	return predicate.Node(sql.FieldNotIn(FieldUpdateTime, vs...))
}

// UpdateTimeGT applies the GT predicate on the "update_time" field.
func UpdateTimeGT(v time.Time) predicate.Node {
	return predicate.Node(sql.FieldGT(FieldUpdateTime, v))
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.Node {
	return predicate.Node(sql.FieldGTE(FieldUpdateTime, v))
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.Node {
	return predicate.Node(sql.FieldLT(FieldUpdateTime, v))
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.Node {
	return predicate.Node(sql.FieldLTE(FieldUpdateTime, v))
}

// PublisherIDEQ applies the EQ predicate on the "publisher_id" field.
func PublisherIDEQ(v string) predicate.Node {
	return predicate.Node(sql.FieldEQ(FieldPublisherID, v))
}

// PublisherIDNEQ applies the NEQ predicate on the "publisher_id" field.
func PublisherIDNEQ(v string) predicate.Node {
	return predicate.Node(sql.FieldNEQ(FieldPublisherID, v))
}

// PublisherIDIn applies the In predicate on the "publisher_id" field.
func PublisherIDIn(vs ...string) predicate.Node {
	return predicate.Node(sql.FieldIn(FieldPublisherID, vs...))
}

// PublisherIDNotIn applies the NotIn predicate on the "publisher_id" field.
func PublisherIDNotIn(vs ...string) predicate.Node {
	return predicate.Node(sql.FieldNotIn(FieldPublisherID, vs...))
}

// PublisherIDGT applies the GT predicate on the "publisher_id" field.
func PublisherIDGT(v string) predicate.Node {
	return predicate.Node(sql.FieldGT(FieldPublisherID, v))
}

// PublisherIDGTE applies the GTE predicate on the "publisher_id" field.
func PublisherIDGTE(v string) predicate.Node {
	return predicate.Node(sql.FieldGTE(FieldPublisherID, v))
}

// PublisherIDLT applies the LT predicate on the "publisher_id" field.
func PublisherIDLT(v string) predicate.Node {
	return predicate.Node(sql.FieldLT(FieldPublisherID, v))
}

// PublisherIDLTE applies the LTE predicate on the "publisher_id" field.
func PublisherIDLTE(v string) predicate.Node {
	return predicate.Node(sql.FieldLTE(FieldPublisherID, v))
}

// PublisherIDContains applies the Contains predicate on the "publisher_id" field.
func PublisherIDContains(v string) predicate.Node {
	return predicate.Node(sql.FieldContains(FieldPublisherID, v))
}

// PublisherIDHasPrefix applies the HasPrefix predicate on the "publisher_id" field.
func PublisherIDHasPrefix(v string) predicate.Node {
	return predicate.Node(sql.FieldHasPrefix(FieldPublisherID, v))
}

// PublisherIDHasSuffix applies the HasSuffix predicate on the "publisher_id" field.
func PublisherIDHasSuffix(v string) predicate.Node {
	return predicate.Node(sql.FieldHasSuffix(FieldPublisherID, v))
}

// PublisherIDEqualFold applies the EqualFold predicate on the "publisher_id" field.
func PublisherIDEqualFold(v string) predicate.Node {
	return predicate.Node(sql.FieldEqualFold(FieldPublisherID, v))
}

// PublisherIDContainsFold applies the ContainsFold predicate on the "publisher_id" field.
func PublisherIDContainsFold(v string) predicate.Node {
	return predicate.Node(sql.FieldContainsFold(FieldPublisherID, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Node {
	return predicate.Node(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Node {
	return predicate.Node(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Node {
	return predicate.Node(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Node {
	return predicate.Node(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Node {
	return predicate.Node(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Node {
	return predicate.Node(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Node {
	return predicate.Node(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Node {
	return predicate.Node(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Node {
	return predicate.Node(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Node {
	return predicate.Node(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Node {
	return predicate.Node(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Node {
	return predicate.Node(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Node {
	return predicate.Node(sql.FieldContainsFold(FieldName, v))
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.Node {
	return predicate.Node(sql.FieldEQ(FieldDescription, v))
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.Node {
	return predicate.Node(sql.FieldNEQ(FieldDescription, v))
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.Node {
	return predicate.Node(sql.FieldIn(FieldDescription, vs...))
}

// DescriptionNotIn applies the NotIn predicate on the "description" field.
func DescriptionNotIn(vs ...string) predicate.Node {
	return predicate.Node(sql.FieldNotIn(FieldDescription, vs...))
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.Node {
	return predicate.Node(sql.FieldGT(FieldDescription, v))
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.Node {
	return predicate.Node(sql.FieldGTE(FieldDescription, v))
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.Node {
	return predicate.Node(sql.FieldLT(FieldDescription, v))
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.Node {
	return predicate.Node(sql.FieldLTE(FieldDescription, v))
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.Node {
	return predicate.Node(sql.FieldContains(FieldDescription, v))
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.Node {
	return predicate.Node(sql.FieldHasPrefix(FieldDescription, v))
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.Node {
	return predicate.Node(sql.FieldHasSuffix(FieldDescription, v))
}

// DescriptionIsNil applies the IsNil predicate on the "description" field.
func DescriptionIsNil() predicate.Node {
	return predicate.Node(sql.FieldIsNull(FieldDescription))
}

// DescriptionNotNil applies the NotNil predicate on the "description" field.
func DescriptionNotNil() predicate.Node {
	return predicate.Node(sql.FieldNotNull(FieldDescription))
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.Node {
	return predicate.Node(sql.FieldEqualFold(FieldDescription, v))
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.Node {
	return predicate.Node(sql.FieldContainsFold(FieldDescription, v))
}

// AuthorEQ applies the EQ predicate on the "author" field.
func AuthorEQ(v string) predicate.Node {
	return predicate.Node(sql.FieldEQ(FieldAuthor, v))
}

// AuthorNEQ applies the NEQ predicate on the "author" field.
func AuthorNEQ(v string) predicate.Node {
	return predicate.Node(sql.FieldNEQ(FieldAuthor, v))
}

// AuthorIn applies the In predicate on the "author" field.
func AuthorIn(vs ...string) predicate.Node {
	return predicate.Node(sql.FieldIn(FieldAuthor, vs...))
}

// AuthorNotIn applies the NotIn predicate on the "author" field.
func AuthorNotIn(vs ...string) predicate.Node {
	return predicate.Node(sql.FieldNotIn(FieldAuthor, vs...))
}

// AuthorGT applies the GT predicate on the "author" field.
func AuthorGT(v string) predicate.Node {
	return predicate.Node(sql.FieldGT(FieldAuthor, v))
}

// AuthorGTE applies the GTE predicate on the "author" field.
func AuthorGTE(v string) predicate.Node {
	return predicate.Node(sql.FieldGTE(FieldAuthor, v))
}

// AuthorLT applies the LT predicate on the "author" field.
func AuthorLT(v string) predicate.Node {
	return predicate.Node(sql.FieldLT(FieldAuthor, v))
}

// AuthorLTE applies the LTE predicate on the "author" field.
func AuthorLTE(v string) predicate.Node {
	return predicate.Node(sql.FieldLTE(FieldAuthor, v))
}

// AuthorContains applies the Contains predicate on the "author" field.
func AuthorContains(v string) predicate.Node {
	return predicate.Node(sql.FieldContains(FieldAuthor, v))
}

// AuthorHasPrefix applies the HasPrefix predicate on the "author" field.
func AuthorHasPrefix(v string) predicate.Node {
	return predicate.Node(sql.FieldHasPrefix(FieldAuthor, v))
}

// AuthorHasSuffix applies the HasSuffix predicate on the "author" field.
func AuthorHasSuffix(v string) predicate.Node {
	return predicate.Node(sql.FieldHasSuffix(FieldAuthor, v))
}

// AuthorIsNil applies the IsNil predicate on the "author" field.
func AuthorIsNil() predicate.Node {
	return predicate.Node(sql.FieldIsNull(FieldAuthor))
}

// AuthorNotNil applies the NotNil predicate on the "author" field.
func AuthorNotNil() predicate.Node {
	return predicate.Node(sql.FieldNotNull(FieldAuthor))
}

// AuthorEqualFold applies the EqualFold predicate on the "author" field.
func AuthorEqualFold(v string) predicate.Node {
	return predicate.Node(sql.FieldEqualFold(FieldAuthor, v))
}

// AuthorContainsFold applies the ContainsFold predicate on the "author" field.
func AuthorContainsFold(v string) predicate.Node {
	return predicate.Node(sql.FieldContainsFold(FieldAuthor, v))
}

// LicenseEQ applies the EQ predicate on the "license" field.
func LicenseEQ(v string) predicate.Node {
	return predicate.Node(sql.FieldEQ(FieldLicense, v))
}

// LicenseNEQ applies the NEQ predicate on the "license" field.
func LicenseNEQ(v string) predicate.Node {
	return predicate.Node(sql.FieldNEQ(FieldLicense, v))
}

// LicenseIn applies the In predicate on the "license" field.
func LicenseIn(vs ...string) predicate.Node {
	return predicate.Node(sql.FieldIn(FieldLicense, vs...))
}

// LicenseNotIn applies the NotIn predicate on the "license" field.
func LicenseNotIn(vs ...string) predicate.Node {
	return predicate.Node(sql.FieldNotIn(FieldLicense, vs...))
}

// LicenseGT applies the GT predicate on the "license" field.
func LicenseGT(v string) predicate.Node {
	return predicate.Node(sql.FieldGT(FieldLicense, v))
}

// LicenseGTE applies the GTE predicate on the "license" field.
func LicenseGTE(v string) predicate.Node {
	return predicate.Node(sql.FieldGTE(FieldLicense, v))
}

// LicenseLT applies the LT predicate on the "license" field.
func LicenseLT(v string) predicate.Node {
	return predicate.Node(sql.FieldLT(FieldLicense, v))
}

// LicenseLTE applies the LTE predicate on the "license" field.
func LicenseLTE(v string) predicate.Node {
	return predicate.Node(sql.FieldLTE(FieldLicense, v))
}

// LicenseContains applies the Contains predicate on the "license" field.
func LicenseContains(v string) predicate.Node {
	return predicate.Node(sql.FieldContains(FieldLicense, v))
}

// LicenseHasPrefix applies the HasPrefix predicate on the "license" field.
func LicenseHasPrefix(v string) predicate.Node {
	return predicate.Node(sql.FieldHasPrefix(FieldLicense, v))
}

// LicenseHasSuffix applies the HasSuffix predicate on the "license" field.
func LicenseHasSuffix(v string) predicate.Node {
	return predicate.Node(sql.FieldHasSuffix(FieldLicense, v))
}

// LicenseEqualFold applies the EqualFold predicate on the "license" field.
func LicenseEqualFold(v string) predicate.Node {
	return predicate.Node(sql.FieldEqualFold(FieldLicense, v))
}

// LicenseContainsFold applies the ContainsFold predicate on the "license" field.
func LicenseContainsFold(v string) predicate.Node {
	return predicate.Node(sql.FieldContainsFold(FieldLicense, v))
}

// RepositoryURLEQ applies the EQ predicate on the "repository_url" field.
func RepositoryURLEQ(v string) predicate.Node {
	return predicate.Node(sql.FieldEQ(FieldRepositoryURL, v))
}

// RepositoryURLNEQ applies the NEQ predicate on the "repository_url" field.
func RepositoryURLNEQ(v string) predicate.Node {
	return predicate.Node(sql.FieldNEQ(FieldRepositoryURL, v))
}

// RepositoryURLIn applies the In predicate on the "repository_url" field.
func RepositoryURLIn(vs ...string) predicate.Node {
	return predicate.Node(sql.FieldIn(FieldRepositoryURL, vs...))
}

// RepositoryURLNotIn applies the NotIn predicate on the "repository_url" field.
func RepositoryURLNotIn(vs ...string) predicate.Node {
	return predicate.Node(sql.FieldNotIn(FieldRepositoryURL, vs...))
}

// RepositoryURLGT applies the GT predicate on the "repository_url" field.
func RepositoryURLGT(v string) predicate.Node {
	return predicate.Node(sql.FieldGT(FieldRepositoryURL, v))
}

// RepositoryURLGTE applies the GTE predicate on the "repository_url" field.
func RepositoryURLGTE(v string) predicate.Node {
	return predicate.Node(sql.FieldGTE(FieldRepositoryURL, v))
}

// RepositoryURLLT applies the LT predicate on the "repository_url" field.
func RepositoryURLLT(v string) predicate.Node {
	return predicate.Node(sql.FieldLT(FieldRepositoryURL, v))
}

// RepositoryURLLTE applies the LTE predicate on the "repository_url" field.
func RepositoryURLLTE(v string) predicate.Node {
	return predicate.Node(sql.FieldLTE(FieldRepositoryURL, v))
}

// RepositoryURLContains applies the Contains predicate on the "repository_url" field.
func RepositoryURLContains(v string) predicate.Node {
	return predicate.Node(sql.FieldContains(FieldRepositoryURL, v))
}

// RepositoryURLHasPrefix applies the HasPrefix predicate on the "repository_url" field.
func RepositoryURLHasPrefix(v string) predicate.Node {
	return predicate.Node(sql.FieldHasPrefix(FieldRepositoryURL, v))
}

// RepositoryURLHasSuffix applies the HasSuffix predicate on the "repository_url" field.
func RepositoryURLHasSuffix(v string) predicate.Node {
	return predicate.Node(sql.FieldHasSuffix(FieldRepositoryURL, v))
}

// RepositoryURLEqualFold applies the EqualFold predicate on the "repository_url" field.
func RepositoryURLEqualFold(v string) predicate.Node {
	return predicate.Node(sql.FieldEqualFold(FieldRepositoryURL, v))
}

// RepositoryURLContainsFold applies the ContainsFold predicate on the "repository_url" field.
func RepositoryURLContainsFold(v string) predicate.Node {
	return predicate.Node(sql.FieldContainsFold(FieldRepositoryURL, v))
}

// IconURLEQ applies the EQ predicate on the "icon_url" field.
func IconURLEQ(v string) predicate.Node {
	return predicate.Node(sql.FieldEQ(FieldIconURL, v))
}

// IconURLNEQ applies the NEQ predicate on the "icon_url" field.
func IconURLNEQ(v string) predicate.Node {
	return predicate.Node(sql.FieldNEQ(FieldIconURL, v))
}

// IconURLIn applies the In predicate on the "icon_url" field.
func IconURLIn(vs ...string) predicate.Node {
	return predicate.Node(sql.FieldIn(FieldIconURL, vs...))
}

// IconURLNotIn applies the NotIn predicate on the "icon_url" field.
func IconURLNotIn(vs ...string) predicate.Node {
	return predicate.Node(sql.FieldNotIn(FieldIconURL, vs...))
}

// IconURLGT applies the GT predicate on the "icon_url" field.
func IconURLGT(v string) predicate.Node {
	return predicate.Node(sql.FieldGT(FieldIconURL, v))
}

// IconURLGTE applies the GTE predicate on the "icon_url" field.
func IconURLGTE(v string) predicate.Node {
	return predicate.Node(sql.FieldGTE(FieldIconURL, v))
}

// IconURLLT applies the LT predicate on the "icon_url" field.
func IconURLLT(v string) predicate.Node {
	return predicate.Node(sql.FieldLT(FieldIconURL, v))
}

// IconURLLTE applies the LTE predicate on the "icon_url" field.
func IconURLLTE(v string) predicate.Node {
	return predicate.Node(sql.FieldLTE(FieldIconURL, v))
}

// IconURLContains applies the Contains predicate on the "icon_url" field.
func IconURLContains(v string) predicate.Node {
	return predicate.Node(sql.FieldContains(FieldIconURL, v))
}

// IconURLHasPrefix applies the HasPrefix predicate on the "icon_url" field.
func IconURLHasPrefix(v string) predicate.Node {
	return predicate.Node(sql.FieldHasPrefix(FieldIconURL, v))
}

// IconURLHasSuffix applies the HasSuffix predicate on the "icon_url" field.
func IconURLHasSuffix(v string) predicate.Node {
	return predicate.Node(sql.FieldHasSuffix(FieldIconURL, v))
}

// IconURLIsNil applies the IsNil predicate on the "icon_url" field.
func IconURLIsNil() predicate.Node {
	return predicate.Node(sql.FieldIsNull(FieldIconURL))
}

// IconURLNotNil applies the NotNil predicate on the "icon_url" field.
func IconURLNotNil() predicate.Node {
	return predicate.Node(sql.FieldNotNull(FieldIconURL))
}

// IconURLEqualFold applies the EqualFold predicate on the "icon_url" field.
func IconURLEqualFold(v string) predicate.Node {
	return predicate.Node(sql.FieldEqualFold(FieldIconURL, v))
}

// IconURLContainsFold applies the ContainsFold predicate on the "icon_url" field.
func IconURLContainsFold(v string) predicate.Node {
	return predicate.Node(sql.FieldContainsFold(FieldIconURL, v))
}

// HasPublisher applies the HasEdge predicate on the "publisher" edge.
func HasPublisher() predicate.Node {
	return predicate.Node(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, PublisherTable, PublisherColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPublisherWith applies the HasEdge predicate on the "publisher" edge with a given conditions (other predicates).
func HasPublisherWith(preds ...predicate.Publisher) predicate.Node {
	return predicate.Node(func(s *sql.Selector) {
		step := newPublisherStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasVersions applies the HasEdge predicate on the "versions" edge.
func HasVersions() predicate.Node {
	return predicate.Node(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, VersionsTable, VersionsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasVersionsWith applies the HasEdge predicate on the "versions" edge with a given conditions (other predicates).
func HasVersionsWith(preds ...predicate.NodeVersion) predicate.Node {
	return predicate.Node(func(s *sql.Selector) {
		step := newVersionsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Node) predicate.Node {
	return predicate.Node(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Node) predicate.Node {
	return predicate.Node(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Node) predicate.Node {
	return predicate.Node(sql.NotPredicates(p))
}
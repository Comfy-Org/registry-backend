// Code generated by ent, DO NOT EDIT.

package publisher

import (
	"registry-backend/ent/predicate"
	"registry-backend/ent/schema"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id string) predicate.Publisher {
	return predicate.Publisher(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.Publisher {
	return predicate.Publisher(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.Publisher {
	return predicate.Publisher(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.Publisher {
	return predicate.Publisher(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...string) predicate.Publisher {
	return predicate.Publisher(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id string) predicate.Publisher {
	return predicate.Publisher(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.Publisher {
	return predicate.Publisher(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.Publisher {
	return predicate.Publisher(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.Publisher {
	return predicate.Publisher(sql.FieldLTE(FieldID, id))
}

// IDEqualFold applies the EqualFold predicate on the ID field.
func IDEqualFold(id string) predicate.Publisher {
	return predicate.Publisher(sql.FieldEqualFold(FieldID, id))
}

// IDContainsFold applies the ContainsFold predicate on the ID field.
func IDContainsFold(id string) predicate.Publisher {
	return predicate.Publisher(sql.FieldContainsFold(FieldID, id))
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.Publisher {
	return predicate.Publisher(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.Publisher {
	return predicate.Publisher(sql.FieldEQ(FieldUpdateTime, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldEQ(FieldName, v))
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldEQ(FieldDescription, v))
}

// Website applies equality check predicate on the "website" field. It's identical to WebsiteEQ.
func Website(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldEQ(FieldWebsite, v))
}

// SupportEmail applies equality check predicate on the "support_email" field. It's identical to SupportEmailEQ.
func SupportEmail(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldEQ(FieldSupportEmail, v))
}

// SourceCodeRepo applies equality check predicate on the "source_code_repo" field. It's identical to SourceCodeRepoEQ.
func SourceCodeRepo(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldEQ(FieldSourceCodeRepo, v))
}

// LogoURL applies equality check predicate on the "logo_url" field. It's identical to LogoURLEQ.
func LogoURL(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldEQ(FieldLogoURL, v))
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.Publisher {
	return predicate.Publisher(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.Publisher {
	return predicate.Publisher(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.Publisher {
	return predicate.Publisher(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "create_time" field.
func CreateTimeNotIn(vs ...time.Time) predicate.Publisher {
	return predicate.Publisher(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "create_time" field.
func CreateTimeGT(v time.Time) predicate.Publisher {
	return predicate.Publisher(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.Publisher {
	return predicate.Publisher(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.Publisher {
	return predicate.Publisher(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.Publisher {
	return predicate.Publisher(sql.FieldLTE(FieldCreateTime, v))
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.Publisher {
	return predicate.Publisher(sql.FieldEQ(FieldUpdateTime, v))
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.Publisher {
	return predicate.Publisher(sql.FieldNEQ(FieldUpdateTime, v))
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.Publisher {
	return predicate.Publisher(sql.FieldIn(FieldUpdateTime, vs...))
}

// UpdateTimeNotIn applies the NotIn predicate on the "update_time" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.Publisher {
	return predicate.Publisher(sql.FieldNotIn(FieldUpdateTime, vs...))
}

// UpdateTimeGT applies the GT predicate on the "update_time" field.
func UpdateTimeGT(v time.Time) predicate.Publisher {
	return predicate.Publisher(sql.FieldGT(FieldUpdateTime, v))
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.Publisher {
	return predicate.Publisher(sql.FieldGTE(FieldUpdateTime, v))
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.Publisher {
	return predicate.Publisher(sql.FieldLT(FieldUpdateTime, v))
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.Publisher {
	return predicate.Publisher(sql.FieldLTE(FieldUpdateTime, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Publisher {
	return predicate.Publisher(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Publisher {
	return predicate.Publisher(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldContainsFold(FieldName, v))
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldEQ(FieldDescription, v))
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldNEQ(FieldDescription, v))
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.Publisher {
	return predicate.Publisher(sql.FieldIn(FieldDescription, vs...))
}

// DescriptionNotIn applies the NotIn predicate on the "description" field.
func DescriptionNotIn(vs ...string) predicate.Publisher {
	return predicate.Publisher(sql.FieldNotIn(FieldDescription, vs...))
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldGT(FieldDescription, v))
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldGTE(FieldDescription, v))
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldLT(FieldDescription, v))
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldLTE(FieldDescription, v))
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldContains(FieldDescription, v))
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldHasPrefix(FieldDescription, v))
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldHasSuffix(FieldDescription, v))
}

// DescriptionIsNil applies the IsNil predicate on the "description" field.
func DescriptionIsNil() predicate.Publisher {
	return predicate.Publisher(sql.FieldIsNull(FieldDescription))
}

// DescriptionNotNil applies the NotNil predicate on the "description" field.
func DescriptionNotNil() predicate.Publisher {
	return predicate.Publisher(sql.FieldNotNull(FieldDescription))
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldEqualFold(FieldDescription, v))
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldContainsFold(FieldDescription, v))
}

// WebsiteEQ applies the EQ predicate on the "website" field.
func WebsiteEQ(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldEQ(FieldWebsite, v))
}

// WebsiteNEQ applies the NEQ predicate on the "website" field.
func WebsiteNEQ(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldNEQ(FieldWebsite, v))
}

// WebsiteIn applies the In predicate on the "website" field.
func WebsiteIn(vs ...string) predicate.Publisher {
	return predicate.Publisher(sql.FieldIn(FieldWebsite, vs...))
}

// WebsiteNotIn applies the NotIn predicate on the "website" field.
func WebsiteNotIn(vs ...string) predicate.Publisher {
	return predicate.Publisher(sql.FieldNotIn(FieldWebsite, vs...))
}

// WebsiteGT applies the GT predicate on the "website" field.
func WebsiteGT(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldGT(FieldWebsite, v))
}

// WebsiteGTE applies the GTE predicate on the "website" field.
func WebsiteGTE(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldGTE(FieldWebsite, v))
}

// WebsiteLT applies the LT predicate on the "website" field.
func WebsiteLT(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldLT(FieldWebsite, v))
}

// WebsiteLTE applies the LTE predicate on the "website" field.
func WebsiteLTE(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldLTE(FieldWebsite, v))
}

// WebsiteContains applies the Contains predicate on the "website" field.
func WebsiteContains(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldContains(FieldWebsite, v))
}

// WebsiteHasPrefix applies the HasPrefix predicate on the "website" field.
func WebsiteHasPrefix(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldHasPrefix(FieldWebsite, v))
}

// WebsiteHasSuffix applies the HasSuffix predicate on the "website" field.
func WebsiteHasSuffix(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldHasSuffix(FieldWebsite, v))
}

// WebsiteIsNil applies the IsNil predicate on the "website" field.
func WebsiteIsNil() predicate.Publisher {
	return predicate.Publisher(sql.FieldIsNull(FieldWebsite))
}

// WebsiteNotNil applies the NotNil predicate on the "website" field.
func WebsiteNotNil() predicate.Publisher {
	return predicate.Publisher(sql.FieldNotNull(FieldWebsite))
}

// WebsiteEqualFold applies the EqualFold predicate on the "website" field.
func WebsiteEqualFold(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldEqualFold(FieldWebsite, v))
}

// WebsiteContainsFold applies the ContainsFold predicate on the "website" field.
func WebsiteContainsFold(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldContainsFold(FieldWebsite, v))
}

// SupportEmailEQ applies the EQ predicate on the "support_email" field.
func SupportEmailEQ(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldEQ(FieldSupportEmail, v))
}

// SupportEmailNEQ applies the NEQ predicate on the "support_email" field.
func SupportEmailNEQ(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldNEQ(FieldSupportEmail, v))
}

// SupportEmailIn applies the In predicate on the "support_email" field.
func SupportEmailIn(vs ...string) predicate.Publisher {
	return predicate.Publisher(sql.FieldIn(FieldSupportEmail, vs...))
}

// SupportEmailNotIn applies the NotIn predicate on the "support_email" field.
func SupportEmailNotIn(vs ...string) predicate.Publisher {
	return predicate.Publisher(sql.FieldNotIn(FieldSupportEmail, vs...))
}

// SupportEmailGT applies the GT predicate on the "support_email" field.
func SupportEmailGT(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldGT(FieldSupportEmail, v))
}

// SupportEmailGTE applies the GTE predicate on the "support_email" field.
func SupportEmailGTE(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldGTE(FieldSupportEmail, v))
}

// SupportEmailLT applies the LT predicate on the "support_email" field.
func SupportEmailLT(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldLT(FieldSupportEmail, v))
}

// SupportEmailLTE applies the LTE predicate on the "support_email" field.
func SupportEmailLTE(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldLTE(FieldSupportEmail, v))
}

// SupportEmailContains applies the Contains predicate on the "support_email" field.
func SupportEmailContains(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldContains(FieldSupportEmail, v))
}

// SupportEmailHasPrefix applies the HasPrefix predicate on the "support_email" field.
func SupportEmailHasPrefix(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldHasPrefix(FieldSupportEmail, v))
}

// SupportEmailHasSuffix applies the HasSuffix predicate on the "support_email" field.
func SupportEmailHasSuffix(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldHasSuffix(FieldSupportEmail, v))
}

// SupportEmailIsNil applies the IsNil predicate on the "support_email" field.
func SupportEmailIsNil() predicate.Publisher {
	return predicate.Publisher(sql.FieldIsNull(FieldSupportEmail))
}

// SupportEmailNotNil applies the NotNil predicate on the "support_email" field.
func SupportEmailNotNil() predicate.Publisher {
	return predicate.Publisher(sql.FieldNotNull(FieldSupportEmail))
}

// SupportEmailEqualFold applies the EqualFold predicate on the "support_email" field.
func SupportEmailEqualFold(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldEqualFold(FieldSupportEmail, v))
}

// SupportEmailContainsFold applies the ContainsFold predicate on the "support_email" field.
func SupportEmailContainsFold(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldContainsFold(FieldSupportEmail, v))
}

// SourceCodeRepoEQ applies the EQ predicate on the "source_code_repo" field.
func SourceCodeRepoEQ(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldEQ(FieldSourceCodeRepo, v))
}

// SourceCodeRepoNEQ applies the NEQ predicate on the "source_code_repo" field.
func SourceCodeRepoNEQ(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldNEQ(FieldSourceCodeRepo, v))
}

// SourceCodeRepoIn applies the In predicate on the "source_code_repo" field.
func SourceCodeRepoIn(vs ...string) predicate.Publisher {
	return predicate.Publisher(sql.FieldIn(FieldSourceCodeRepo, vs...))
}

// SourceCodeRepoNotIn applies the NotIn predicate on the "source_code_repo" field.
func SourceCodeRepoNotIn(vs ...string) predicate.Publisher {
	return predicate.Publisher(sql.FieldNotIn(FieldSourceCodeRepo, vs...))
}

// SourceCodeRepoGT applies the GT predicate on the "source_code_repo" field.
func SourceCodeRepoGT(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldGT(FieldSourceCodeRepo, v))
}

// SourceCodeRepoGTE applies the GTE predicate on the "source_code_repo" field.
func SourceCodeRepoGTE(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldGTE(FieldSourceCodeRepo, v))
}

// SourceCodeRepoLT applies the LT predicate on the "source_code_repo" field.
func SourceCodeRepoLT(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldLT(FieldSourceCodeRepo, v))
}

// SourceCodeRepoLTE applies the LTE predicate on the "source_code_repo" field.
func SourceCodeRepoLTE(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldLTE(FieldSourceCodeRepo, v))
}

// SourceCodeRepoContains applies the Contains predicate on the "source_code_repo" field.
func SourceCodeRepoContains(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldContains(FieldSourceCodeRepo, v))
}

// SourceCodeRepoHasPrefix applies the HasPrefix predicate on the "source_code_repo" field.
func SourceCodeRepoHasPrefix(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldHasPrefix(FieldSourceCodeRepo, v))
}

// SourceCodeRepoHasSuffix applies the HasSuffix predicate on the "source_code_repo" field.
func SourceCodeRepoHasSuffix(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldHasSuffix(FieldSourceCodeRepo, v))
}

// SourceCodeRepoIsNil applies the IsNil predicate on the "source_code_repo" field.
func SourceCodeRepoIsNil() predicate.Publisher {
	return predicate.Publisher(sql.FieldIsNull(FieldSourceCodeRepo))
}

// SourceCodeRepoNotNil applies the NotNil predicate on the "source_code_repo" field.
func SourceCodeRepoNotNil() predicate.Publisher {
	return predicate.Publisher(sql.FieldNotNull(FieldSourceCodeRepo))
}

// SourceCodeRepoEqualFold applies the EqualFold predicate on the "source_code_repo" field.
func SourceCodeRepoEqualFold(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldEqualFold(FieldSourceCodeRepo, v))
}

// SourceCodeRepoContainsFold applies the ContainsFold predicate on the "source_code_repo" field.
func SourceCodeRepoContainsFold(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldContainsFold(FieldSourceCodeRepo, v))
}

// LogoURLEQ applies the EQ predicate on the "logo_url" field.
func LogoURLEQ(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldEQ(FieldLogoURL, v))
}

// LogoURLNEQ applies the NEQ predicate on the "logo_url" field.
func LogoURLNEQ(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldNEQ(FieldLogoURL, v))
}

// LogoURLIn applies the In predicate on the "logo_url" field.
func LogoURLIn(vs ...string) predicate.Publisher {
	return predicate.Publisher(sql.FieldIn(FieldLogoURL, vs...))
}

// LogoURLNotIn applies the NotIn predicate on the "logo_url" field.
func LogoURLNotIn(vs ...string) predicate.Publisher {
	return predicate.Publisher(sql.FieldNotIn(FieldLogoURL, vs...))
}

// LogoURLGT applies the GT predicate on the "logo_url" field.
func LogoURLGT(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldGT(FieldLogoURL, v))
}

// LogoURLGTE applies the GTE predicate on the "logo_url" field.
func LogoURLGTE(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldGTE(FieldLogoURL, v))
}

// LogoURLLT applies the LT predicate on the "logo_url" field.
func LogoURLLT(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldLT(FieldLogoURL, v))
}

// LogoURLLTE applies the LTE predicate on the "logo_url" field.
func LogoURLLTE(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldLTE(FieldLogoURL, v))
}

// LogoURLContains applies the Contains predicate on the "logo_url" field.
func LogoURLContains(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldContains(FieldLogoURL, v))
}

// LogoURLHasPrefix applies the HasPrefix predicate on the "logo_url" field.
func LogoURLHasPrefix(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldHasPrefix(FieldLogoURL, v))
}

// LogoURLHasSuffix applies the HasSuffix predicate on the "logo_url" field.
func LogoURLHasSuffix(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldHasSuffix(FieldLogoURL, v))
}

// LogoURLIsNil applies the IsNil predicate on the "logo_url" field.
func LogoURLIsNil() predicate.Publisher {
	return predicate.Publisher(sql.FieldIsNull(FieldLogoURL))
}

// LogoURLNotNil applies the NotNil predicate on the "logo_url" field.
func LogoURLNotNil() predicate.Publisher {
	return predicate.Publisher(sql.FieldNotNull(FieldLogoURL))
}

// LogoURLEqualFold applies the EqualFold predicate on the "logo_url" field.
func LogoURLEqualFold(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldEqualFold(FieldLogoURL, v))
}

// LogoURLContainsFold applies the ContainsFold predicate on the "logo_url" field.
func LogoURLContainsFold(v string) predicate.Publisher {
	return predicate.Publisher(sql.FieldContainsFold(FieldLogoURL, v))
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v schema.PublisherStatus) predicate.Publisher {
	vc := v
	return predicate.Publisher(sql.FieldEQ(FieldStatus, vc))
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v schema.PublisherStatus) predicate.Publisher {
	vc := v
	return predicate.Publisher(sql.FieldNEQ(FieldStatus, vc))
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...schema.PublisherStatus) predicate.Publisher {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Publisher(sql.FieldIn(FieldStatus, v...))
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...schema.PublisherStatus) predicate.Publisher {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Publisher(sql.FieldNotIn(FieldStatus, v...))
}

// HasPublisherPermissions applies the HasEdge predicate on the "publisher_permissions" edge.
func HasPublisherPermissions() predicate.Publisher {
	return predicate.Publisher(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PublisherPermissionsTable, PublisherPermissionsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPublisherPermissionsWith applies the HasEdge predicate on the "publisher_permissions" edge with a given conditions (other predicates).
func HasPublisherPermissionsWith(preds ...predicate.PublisherPermission) predicate.Publisher {
	return predicate.Publisher(func(s *sql.Selector) {
		step := newPublisherPermissionsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasNodes applies the HasEdge predicate on the "nodes" edge.
func HasNodes() predicate.Publisher {
	return predicate.Publisher(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, NodesTable, NodesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasNodesWith applies the HasEdge predicate on the "nodes" edge with a given conditions (other predicates).
func HasNodesWith(preds ...predicate.Node) predicate.Publisher {
	return predicate.Publisher(func(s *sql.Selector) {
		step := newNodesStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasPersonalAccessTokens applies the HasEdge predicate on the "personal_access_tokens" edge.
func HasPersonalAccessTokens() predicate.Publisher {
	return predicate.Publisher(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PersonalAccessTokensTable, PersonalAccessTokensColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPersonalAccessTokensWith applies the HasEdge predicate on the "personal_access_tokens" edge with a given conditions (other predicates).
func HasPersonalAccessTokensWith(preds ...predicate.PersonalAccessToken) predicate.Publisher {
	return predicate.Publisher(func(s *sql.Selector) {
		step := newPersonalAccessTokensStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Publisher) predicate.Publisher {
	return predicate.Publisher(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Publisher) predicate.Publisher {
	return predicate.Publisher(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Publisher) predicate.Publisher {
	return predicate.Publisher(sql.NotPredicates(p))
}

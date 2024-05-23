// Code generated by ent, DO NOT EDIT.

package publisher

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the publisher type in the database.
	Label = "publisher"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time field in the database.
	FieldUpdateTime = "update_time"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldWebsite holds the string denoting the website field in the database.
	FieldWebsite = "website"
	// FieldSupportEmail holds the string denoting the support_email field in the database.
	FieldSupportEmail = "support_email"
	// FieldSourceCodeRepo holds the string denoting the source_code_repo field in the database.
	FieldSourceCodeRepo = "source_code_repo"
	// FieldLogoURL holds the string denoting the logo_url field in the database.
	FieldLogoURL = "logo_url"
	// EdgePublisherPermissions holds the string denoting the publisher_permissions edge name in mutations.
	EdgePublisherPermissions = "publisher_permissions"
	// EdgeNodes holds the string denoting the nodes edge name in mutations.
	EdgeNodes = "nodes"
	// EdgePersonalAccessTokens holds the string denoting the personal_access_tokens edge name in mutations.
	EdgePersonalAccessTokens = "personal_access_tokens"
	// Table holds the table name of the publisher in the database.
	Table = "publishers"
	// PublisherPermissionsTable is the table that holds the publisher_permissions relation/edge.
	PublisherPermissionsTable = "publisher_permissions"
	// PublisherPermissionsInverseTable is the table name for the PublisherPermission entity.
	// It exists in this package in order to avoid circular dependency with the "publisherpermission" package.
	PublisherPermissionsInverseTable = "publisher_permissions"
	// PublisherPermissionsColumn is the table column denoting the publisher_permissions relation/edge.
	PublisherPermissionsColumn = "publisher_id"
	// NodesTable is the table that holds the nodes relation/edge.
	NodesTable = "nodes"
	// NodesInverseTable is the table name for the Node entity.
	// It exists in this package in order to avoid circular dependency with the "node" package.
	NodesInverseTable = "nodes"
	// NodesColumn is the table column denoting the nodes relation/edge.
	NodesColumn = "publisher_id"
	// PersonalAccessTokensTable is the table that holds the personal_access_tokens relation/edge.
	PersonalAccessTokensTable = "personal_access_tokens"
	// PersonalAccessTokensInverseTable is the table name for the PersonalAccessToken entity.
	// It exists in this package in order to avoid circular dependency with the "personalaccesstoken" package.
	PersonalAccessTokensInverseTable = "personal_access_tokens"
	// PersonalAccessTokensColumn is the table column denoting the personal_access_tokens relation/edge.
	PersonalAccessTokensColumn = "publisher_id"
)

// Columns holds all SQL columns for publisher fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
	FieldDescription,
	FieldWebsite,
	FieldSupportEmail,
	FieldSourceCodeRepo,
	FieldLogoURL,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreateTime holds the default value on creation for the "create_time" field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the "update_time" field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the "update_time" field.
	UpdateDefaultUpdateTime func() time.Time
)

// OrderOption defines the ordering options for the Publisher queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCreateTime orders the results by the create_time field.
func ByCreateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreateTime, opts...).ToFunc()
}

// ByUpdateTime orders the results by the update_time field.
func ByUpdateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdateTime, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByDescription orders the results by the description field.
func ByDescription(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDescription, opts...).ToFunc()
}

// ByWebsite orders the results by the website field.
func ByWebsite(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldWebsite, opts...).ToFunc()
}

// BySupportEmail orders the results by the support_email field.
func BySupportEmail(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSupportEmail, opts...).ToFunc()
}

// BySourceCodeRepo orders the results by the source_code_repo field.
func BySourceCodeRepo(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSourceCodeRepo, opts...).ToFunc()
}

// ByLogoURL orders the results by the logo_url field.
func ByLogoURL(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLogoURL, opts...).ToFunc()
}

// ByPublisherPermissionsCount orders the results by publisher_permissions count.
func ByPublisherPermissionsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newPublisherPermissionsStep(), opts...)
	}
}

// ByPublisherPermissions orders the results by publisher_permissions terms.
func ByPublisherPermissions(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newPublisherPermissionsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByNodesCount orders the results by nodes count.
func ByNodesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newNodesStep(), opts...)
	}
}

// ByNodes orders the results by nodes terms.
func ByNodes(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newNodesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByPersonalAccessTokensCount orders the results by personal_access_tokens count.
func ByPersonalAccessTokensCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newPersonalAccessTokensStep(), opts...)
	}
}

// ByPersonalAccessTokens orders the results by personal_access_tokens terms.
func ByPersonalAccessTokens(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newPersonalAccessTokensStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newPublisherPermissionsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PublisherPermissionsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, PublisherPermissionsTable, PublisherPermissionsColumn),
	)
}
func newNodesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(NodesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, NodesTable, NodesColumn),
	)
}
func newPersonalAccessTokensStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PersonalAccessTokensInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, PersonalAccessTokensTable, PersonalAccessTokensColumn),
	)
}

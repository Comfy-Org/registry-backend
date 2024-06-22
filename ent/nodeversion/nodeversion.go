// Code generated by ent, DO NOT EDIT.

package nodeversion

import (
	"fmt"
	"registry-backend/ent/schema"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the nodeversion type in the database.
	Label = "node_version"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time field in the database.
	FieldUpdateTime = "update_time"
	// FieldNodeID holds the string denoting the node_id field in the database.
	FieldNodeID = "node_id"
	// FieldVersion holds the string denoting the version field in the database.
	FieldVersion = "version"
	// FieldChangelog holds the string denoting the changelog field in the database.
	FieldChangelog = "changelog"
	// FieldPipDependencies holds the string denoting the pip_dependencies field in the database.
	FieldPipDependencies = "pip_dependencies"
	// FieldDeprecated holds the string denoting the deprecated field in the database.
	FieldDeprecated = "deprecated"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldStatusReason holds the string denoting the status_reason field in the database.
	FieldStatusReason = "status_reason"
	// EdgeNode holds the string denoting the node edge name in mutations.
	EdgeNode = "node"
	// EdgeStorageFile holds the string denoting the storage_file edge name in mutations.
	EdgeStorageFile = "storage_file"
	// Table holds the table name of the nodeversion in the database.
	Table = "node_versions"
	// NodeTable is the table that holds the node relation/edge.
	NodeTable = "node_versions"
	// NodeInverseTable is the table name for the Node entity.
	// It exists in this package in order to avoid circular dependency with the "node" package.
	NodeInverseTable = "nodes"
	// NodeColumn is the table column denoting the node relation/edge.
	NodeColumn = "node_id"
	// StorageFileTable is the table that holds the storage_file relation/edge.
	StorageFileTable = "node_versions"
	// StorageFileInverseTable is the table name for the StorageFile entity.
	// It exists in this package in order to avoid circular dependency with the "storagefile" package.
	StorageFileInverseTable = "storage_files"
	// StorageFileColumn is the table column denoting the storage_file relation/edge.
	StorageFileColumn = "node_version_storage_file"
)

// Columns holds all SQL columns for nodeversion fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldNodeID,
	FieldVersion,
	FieldChangelog,
	FieldPipDependencies,
	FieldDeprecated,
	FieldStatus,
	FieldStatusReason,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "node_versions"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"node_version_storage_file",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
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
	// DefaultDeprecated holds the default value on creation for the "deprecated" field.
	DefaultDeprecated bool
	// DefaultStatusReason holds the default value on creation for the "status_reason" field.
	DefaultStatusReason string
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

const DefaultStatus schema.NodeVersionStatus = "pending"

// StatusValidator is a validator for the "status" field enum values. It is called by the builders before save.
func StatusValidator(s schema.NodeVersionStatus) error {
	switch s {
	case "active", "banned", "deleted", "pending":
		return nil
	default:
		return fmt.Errorf("nodeversion: invalid enum value for status field: %q", s)
	}
}

// OrderOption defines the ordering options for the NodeVersion queries.
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

// ByNodeID orders the results by the node_id field.
func ByNodeID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldNodeID, opts...).ToFunc()
}

// ByVersion orders the results by the version field.
func ByVersion(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldVersion, opts...).ToFunc()
}

// ByChangelog orders the results by the changelog field.
func ByChangelog(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldChangelog, opts...).ToFunc()
}

// ByDeprecated orders the results by the deprecated field.
func ByDeprecated(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDeprecated, opts...).ToFunc()
}

// ByStatus orders the results by the status field.
func ByStatus(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStatus, opts...).ToFunc()
}

// ByStatusReason orders the results by the status_reason field.
func ByStatusReason(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStatusReason, opts...).ToFunc()
}

// ByNodeField orders the results by node field.
func ByNodeField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newNodeStep(), sql.OrderByField(field, opts...))
	}
}

// ByStorageFileField orders the results by storage_file field.
func ByStorageFileField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newStorageFileStep(), sql.OrderByField(field, opts...))
	}
}
func newNodeStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(NodeInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, NodeTable, NodeColumn),
	)
}
func newStorageFileStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(StorageFileInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, StorageFileTable, StorageFileColumn),
	)
}

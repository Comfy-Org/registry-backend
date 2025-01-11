// Code generated by ent, DO NOT EDIT.

package comfynode

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the comfynode type in the database.
	Label = "comfy_node"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time field in the database.
	FieldUpdateTime = "update_time"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldNodeVersionID holds the string denoting the node_version_id field in the database.
	FieldNodeVersionID = "node_version_id"
	// FieldCategory holds the string denoting the category field in the database.
	FieldCategory = "category"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldInputTypes holds the string denoting the input_types field in the database.
	FieldInputTypes = "input_types"
	// FieldDeprecated holds the string denoting the deprecated field in the database.
	FieldDeprecated = "deprecated"
	// FieldExperimental holds the string denoting the experimental field in the database.
	FieldExperimental = "experimental"
	// FieldOutputIsList holds the string denoting the output_is_list field in the database.
	FieldOutputIsList = "output_is_list"
	// FieldReturnNames holds the string denoting the return_names field in the database.
	FieldReturnNames = "return_names"
	// FieldReturnTypes holds the string denoting the return_types field in the database.
	FieldReturnTypes = "return_types"
	// FieldFunction holds the string denoting the function field in the database.
	FieldFunction = "function"
	// EdgeVersions holds the string denoting the versions edge name in mutations.
	EdgeVersions = "versions"
	// Table holds the table name of the comfynode in the database.
	Table = "comfy_nodes"
	// VersionsTable is the table that holds the versions relation/edge.
	VersionsTable = "comfy_nodes"
	// VersionsInverseTable is the table name for the NodeVersion entity.
	// It exists in this package in order to avoid circular dependency with the "nodeversion" package.
	VersionsInverseTable = "node_versions"
	// VersionsColumn is the table column denoting the versions relation/edge.
	VersionsColumn = "node_version_id"
)

// Columns holds all SQL columns for comfynode fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
	FieldNodeVersionID,
	FieldCategory,
	FieldDescription,
	FieldInputTypes,
	FieldDeprecated,
	FieldExperimental,
	FieldOutputIsList,
	FieldReturnNames,
	FieldReturnTypes,
	FieldFunction,
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
	// DefaultDeprecated holds the default value on creation for the "deprecated" field.
	DefaultDeprecated bool
	// DefaultExperimental holds the default value on creation for the "experimental" field.
	DefaultExperimental bool
	// DefaultOutputIsList holds the default value on creation for the "output_is_list" field.
	DefaultOutputIsList []bool
	// DefaultReturnNames holds the default value on creation for the "return_names" field.
	DefaultReturnNames []string
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the ComfyNode queries.
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

// ByNodeVersionID orders the results by the node_version_id field.
func ByNodeVersionID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldNodeVersionID, opts...).ToFunc()
}

// ByCategory orders the results by the category field.
func ByCategory(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCategory, opts...).ToFunc()
}

// ByDescription orders the results by the description field.
func ByDescription(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDescription, opts...).ToFunc()
}

// ByInputTypes orders the results by the input_types field.
func ByInputTypes(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldInputTypes, opts...).ToFunc()
}

// ByDeprecated orders the results by the deprecated field.
func ByDeprecated(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDeprecated, opts...).ToFunc()
}

// ByExperimental orders the results by the experimental field.
func ByExperimental(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldExperimental, opts...).ToFunc()
}

// ByReturnTypes orders the results by the return_types field.
func ByReturnTypes(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldReturnTypes, opts...).ToFunc()
}

// ByFunction orders the results by the function field.
func ByFunction(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFunction, opts...).ToFunc()
}

// ByVersionsField orders the results by versions field.
func ByVersionsField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newVersionsStep(), sql.OrderByField(field, opts...))
	}
}
func newVersionsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(VersionsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, VersionsTable, VersionsColumn),
	)
}

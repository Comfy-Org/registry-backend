// Code generated by ent, DO NOT EDIT.

package ent

import (
	"registry-backend/ent/node"
	"registry-backend/ent/nodeversion"
	"registry-backend/ent/storagefile"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

// NodeVersion is the model entity for the NodeVersion schema.
type NodeVersion struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// NodeID holds the value of the "node_id" field.
	NodeID string `json:"node_id,omitempty"`
	// Must be SemVer compliant
	Version string `json:"version,omitempty"`
	// Changelog holds the value of the "changelog" field.
	Changelog string `json:"changelog,omitempty"`
	// PipDependencies holds the value of the "pip_dependencies" field.
	PipDependencies []string `json:"pip_dependencies,omitempty"`
	// Deprecated holds the value of the "deprecated" field.
	Deprecated bool `json:"deprecated,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the NodeVersionQuery when eager-loading is set.
	Edges                     NodeVersionEdges `json:"edges"`
	node_version_storage_file *uuid.UUID
	selectValues              sql.SelectValues
}

// NodeVersionEdges holds the relations/edges for other nodes in the graph.
type NodeVersionEdges struct {
	// Node holds the value of the node edge.
	Node *Node `json:"node,omitempty"`
	// StorageFile holds the value of the storage_file edge.
	StorageFile *StorageFile `json:"storage_file,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// NodeOrErr returns the Node value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e NodeVersionEdges) NodeOrErr() (*Node, error) {
	if e.Node != nil {
		return e.Node, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: node.Label}
	}
	return nil, &NotLoadedError{edge: "node"}
}

// StorageFileOrErr returns the StorageFile value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e NodeVersionEdges) StorageFileOrErr() (*StorageFile, error) {
	if e.StorageFile != nil {
		return e.StorageFile, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: storagefile.Label}
	}
	return nil, &NotLoadedError{edge: "storage_file"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*NodeVersion) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case nodeversion.FieldPipDependencies:
			values[i] = new([]byte)
		case nodeversion.FieldDeprecated:
			values[i] = new(sql.NullBool)
		case nodeversion.FieldNodeID, nodeversion.FieldVersion, nodeversion.FieldChangelog:
			values[i] = new(sql.NullString)
		case nodeversion.FieldCreateTime, nodeversion.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		case nodeversion.FieldID:
			values[i] = new(uuid.UUID)
		case nodeversion.ForeignKeys[0]: // node_version_storage_file
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the NodeVersion fields.
func (nv *NodeVersion) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case nodeversion.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				nv.ID = *value
			}
		case nodeversion.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field create_time", values[i])
			} else if value.Valid {
				nv.CreateTime = value.Time
			}
		case nodeversion.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field update_time", values[i])
			} else if value.Valid {
				nv.UpdateTime = value.Time
			}
		case nodeversion.FieldNodeID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field node_id", values[i])
			} else if value.Valid {
				nv.NodeID = value.String
			}
		case nodeversion.FieldVersion:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field version", values[i])
			} else if value.Valid {
				nv.Version = value.String
			}
		case nodeversion.FieldChangelog:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field changelog", values[i])
			} else if value.Valid {
				nv.Changelog = value.String
			}
		case nodeversion.FieldPipDependencies:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field pip_dependencies", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &nv.PipDependencies); err != nil {
					return fmt.Errorf("unmarshal field pip_dependencies: %w", err)
				}
			}
		case nodeversion.FieldDeprecated:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field deprecated", values[i])
			} else if value.Valid {
				nv.Deprecated = value.Bool
			}
		case nodeversion.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field node_version_storage_file", values[i])
			} else if value.Valid {
				nv.node_version_storage_file = new(uuid.UUID)
				*nv.node_version_storage_file = *value.S.(*uuid.UUID)
			}
		default:
			nv.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the NodeVersion.
// This includes values selected through modifiers, order, etc.
func (nv *NodeVersion) Value(name string) (ent.Value, error) {
	return nv.selectValues.Get(name)
}

// QueryNode queries the "node" edge of the NodeVersion entity.
func (nv *NodeVersion) QueryNode() *NodeQuery {
	return NewNodeVersionClient(nv.config).QueryNode(nv)
}

// QueryStorageFile queries the "storage_file" edge of the NodeVersion entity.
func (nv *NodeVersion) QueryStorageFile() *StorageFileQuery {
	return NewNodeVersionClient(nv.config).QueryStorageFile(nv)
}

// Update returns a builder for updating this NodeVersion.
// Note that you need to call NodeVersion.Unwrap() before calling this method if this NodeVersion
// was returned from a transaction, and the transaction was committed or rolled back.
func (nv *NodeVersion) Update() *NodeVersionUpdateOne {
	return NewNodeVersionClient(nv.config).UpdateOne(nv)
}

// Unwrap unwraps the NodeVersion entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (nv *NodeVersion) Unwrap() *NodeVersion {
	_tx, ok := nv.config.driver.(*txDriver)
	if !ok {
		panic("ent: NodeVersion is not a transactional entity")
	}
	nv.config.driver = _tx.drv
	return nv
}

// String implements the fmt.Stringer.
func (nv *NodeVersion) String() string {
	var builder strings.Builder
	builder.WriteString("NodeVersion(")
	builder.WriteString(fmt.Sprintf("id=%v, ", nv.ID))
	builder.WriteString("create_time=")
	builder.WriteString(nv.CreateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("update_time=")
	builder.WriteString(nv.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("node_id=")
	builder.WriteString(nv.NodeID)
	builder.WriteString(", ")
	builder.WriteString("version=")
	builder.WriteString(nv.Version)
	builder.WriteString(", ")
	builder.WriteString("changelog=")
	builder.WriteString(nv.Changelog)
	builder.WriteString(", ")
	builder.WriteString("pip_dependencies=")
	builder.WriteString(fmt.Sprintf("%v", nv.PipDependencies))
	builder.WriteString(", ")
	builder.WriteString("deprecated=")
	builder.WriteString(fmt.Sprintf("%v", nv.Deprecated))
	builder.WriteByte(')')
	return builder.String()
}

// NodeVersions is a parsable slice of NodeVersion.
type NodeVersions []*NodeVersion

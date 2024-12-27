// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"registry-backend/ent/comfynode"
	"registry-backend/ent/nodeversion"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

// ComfyNode is the model entity for the ComfyNode schema.
type ComfyNode struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// NodeVersionID holds the value of the "node_version_id" field.
	NodeVersionID uuid.UUID `json:"node_version_id,omitempty"`
	// Category holds the value of the "category" field.
	Category string `json:"category,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// InputTypes holds the value of the "input_types" field.
	InputTypes string `json:"input_types,omitempty"`
	// Deprecated holds the value of the "deprecated" field.
	Deprecated bool `json:"deprecated,omitempty"`
	// Experimental holds the value of the "experimental" field.
	Experimental bool `json:"experimental,omitempty"`
	// OutputIsList holds the value of the "output_is_list" field.
	OutputIsList []bool `json:"output_is_list,omitempty"`
	// ReturnNames holds the value of the "return_names" field.
	ReturnNames []string `json:"return_names,omitempty"`
	// ReturnTypes holds the value of the "return_types" field.
	ReturnTypes []string `json:"return_types,omitempty"`
	// Function holds the value of the "function" field.
	Function string `json:"function,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ComfyNodeQuery when eager-loading is set.
	Edges        ComfyNodeEdges `json:"edges"`
	selectValues sql.SelectValues
}

// ComfyNodeEdges holds the relations/edges for other nodes in the graph.
type ComfyNodeEdges struct {
	// Versions holds the value of the versions edge.
	Versions *NodeVersion `json:"versions,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// VersionsOrErr returns the Versions value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ComfyNodeEdges) VersionsOrErr() (*NodeVersion, error) {
	if e.Versions != nil {
		return e.Versions, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: nodeversion.Label}
	}
	return nil, &NotLoadedError{edge: "versions"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ComfyNode) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case comfynode.FieldOutputIsList, comfynode.FieldReturnNames, comfynode.FieldReturnTypes:
			values[i] = new([]byte)
		case comfynode.FieldDeprecated, comfynode.FieldExperimental:
			values[i] = new(sql.NullBool)
		case comfynode.FieldID, comfynode.FieldCategory, comfynode.FieldDescription, comfynode.FieldInputTypes, comfynode.FieldFunction:
			values[i] = new(sql.NullString)
		case comfynode.FieldCreateTime, comfynode.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		case comfynode.FieldNodeVersionID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ComfyNode fields.
func (cn *ComfyNode) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case comfynode.FieldID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value.Valid {
				cn.ID = value.String
			}
		case comfynode.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field create_time", values[i])
			} else if value.Valid {
				cn.CreateTime = value.Time
			}
		case comfynode.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field update_time", values[i])
			} else if value.Valid {
				cn.UpdateTime = value.Time
			}
		case comfynode.FieldNodeVersionID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field node_version_id", values[i])
			} else if value != nil {
				cn.NodeVersionID = *value
			}
		case comfynode.FieldCategory:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field category", values[i])
			} else if value.Valid {
				cn.Category = value.String
			}
		case comfynode.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				cn.Description = value.String
			}
		case comfynode.FieldInputTypes:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field input_types", values[i])
			} else if value.Valid {
				cn.InputTypes = value.String
			}
		case comfynode.FieldDeprecated:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field deprecated", values[i])
			} else if value.Valid {
				cn.Deprecated = value.Bool
			}
		case comfynode.FieldExperimental:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field experimental", values[i])
			} else if value.Valid {
				cn.Experimental = value.Bool
			}
		case comfynode.FieldOutputIsList:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field output_is_list", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &cn.OutputIsList); err != nil {
					return fmt.Errorf("unmarshal field output_is_list: %w", err)
				}
			}
		case comfynode.FieldReturnNames:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field return_names", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &cn.ReturnNames); err != nil {
					return fmt.Errorf("unmarshal field return_names: %w", err)
				}
			}
		case comfynode.FieldReturnTypes:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field return_types", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &cn.ReturnTypes); err != nil {
					return fmt.Errorf("unmarshal field return_types: %w", err)
				}
			}
		case comfynode.FieldFunction:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field function", values[i])
			} else if value.Valid {
				cn.Function = value.String
			}
		default:
			cn.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the ComfyNode.
// This includes values selected through modifiers, order, etc.
func (cn *ComfyNode) Value(name string) (ent.Value, error) {
	return cn.selectValues.Get(name)
}

// QueryVersions queries the "versions" edge of the ComfyNode entity.
func (cn *ComfyNode) QueryVersions() *NodeVersionQuery {
	return NewComfyNodeClient(cn.config).QueryVersions(cn)
}

// Update returns a builder for updating this ComfyNode.
// Note that you need to call ComfyNode.Unwrap() before calling this method if this ComfyNode
// was returned from a transaction, and the transaction was committed or rolled back.
func (cn *ComfyNode) Update() *ComfyNodeUpdateOne {
	return NewComfyNodeClient(cn.config).UpdateOne(cn)
}

// Unwrap unwraps the ComfyNode entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (cn *ComfyNode) Unwrap() *ComfyNode {
	_tx, ok := cn.config.driver.(*txDriver)
	if !ok {
		panic("ent: ComfyNode is not a transactional entity")
	}
	cn.config.driver = _tx.drv
	return cn
}

// String implements the fmt.Stringer.
func (cn *ComfyNode) String() string {
	var builder strings.Builder
	builder.WriteString("ComfyNode(")
	builder.WriteString(fmt.Sprintf("id=%v, ", cn.ID))
	builder.WriteString("create_time=")
	builder.WriteString(cn.CreateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("update_time=")
	builder.WriteString(cn.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("node_version_id=")
	builder.WriteString(fmt.Sprintf("%v", cn.NodeVersionID))
	builder.WriteString(", ")
	builder.WriteString("category=")
	builder.WriteString(cn.Category)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(cn.Description)
	builder.WriteString(", ")
	builder.WriteString("input_types=")
	builder.WriteString(cn.InputTypes)
	builder.WriteString(", ")
	builder.WriteString("deprecated=")
	builder.WriteString(fmt.Sprintf("%v", cn.Deprecated))
	builder.WriteString(", ")
	builder.WriteString("experimental=")
	builder.WriteString(fmt.Sprintf("%v", cn.Experimental))
	builder.WriteString(", ")
	builder.WriteString("output_is_list=")
	builder.WriteString(fmt.Sprintf("%v", cn.OutputIsList))
	builder.WriteString(", ")
	builder.WriteString("return_names=")
	builder.WriteString(fmt.Sprintf("%v", cn.ReturnNames))
	builder.WriteString(", ")
	builder.WriteString("return_types=")
	builder.WriteString(fmt.Sprintf("%v", cn.ReturnTypes))
	builder.WriteString(", ")
	builder.WriteString("function=")
	builder.WriteString(cn.Function)
	builder.WriteByte(')')
	return builder.String()
}

// ComfyNodes is a parsable slice of ComfyNode.
type ComfyNodes []*ComfyNode
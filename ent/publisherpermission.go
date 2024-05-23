// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"registry-backend/ent/publisher"
	"registry-backend/ent/publisherpermission"
	"registry-backend/ent/schema"
	"registry-backend/ent/user"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// PublisherPermission is the model entity for the PublisherPermission schema.
type PublisherPermission struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Permission holds the value of the "permission" field.
	Permission schema.PublisherPermissionType `json:"permission,omitempty"`
	// UserID holds the value of the "user_id" field.
	UserID string `json:"user_id,omitempty"`
	// PublisherID holds the value of the "publisher_id" field.
	PublisherID string `json:"publisher_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the PublisherPermissionQuery when eager-loading is set.
	Edges        PublisherPermissionEdges `json:"edges"`
	selectValues sql.SelectValues
}

// PublisherPermissionEdges holds the relations/edges for other nodes in the graph.
type PublisherPermissionEdges struct {
	// User holds the value of the user edge.
	User *User `json:"user,omitempty"`
	// Publisher holds the value of the publisher edge.
	Publisher *Publisher `json:"publisher,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// UserOrErr returns the User value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PublisherPermissionEdges) UserOrErr() (*User, error) {
	if e.User != nil {
		return e.User, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: user.Label}
	}
	return nil, &NotLoadedError{edge: "user"}
}

// PublisherOrErr returns the Publisher value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PublisherPermissionEdges) PublisherOrErr() (*Publisher, error) {
	if e.Publisher != nil {
		return e.Publisher, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: publisher.Label}
	}
	return nil, &NotLoadedError{edge: "publisher"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*PublisherPermission) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case publisherpermission.FieldID:
			values[i] = new(sql.NullInt64)
		case publisherpermission.FieldPermission, publisherpermission.FieldUserID, publisherpermission.FieldPublisherID:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the PublisherPermission fields.
func (pp *PublisherPermission) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case publisherpermission.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			pp.ID = int(value.Int64)
		case publisherpermission.FieldPermission:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field permission", values[i])
			} else if value.Valid {
				pp.Permission = schema.PublisherPermissionType(value.String)
			}
		case publisherpermission.FieldUserID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field user_id", values[i])
			} else if value.Valid {
				pp.UserID = value.String
			}
		case publisherpermission.FieldPublisherID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field publisher_id", values[i])
			} else if value.Valid {
				pp.PublisherID = value.String
			}
		default:
			pp.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the PublisherPermission.
// This includes values selected through modifiers, order, etc.
func (pp *PublisherPermission) Value(name string) (ent.Value, error) {
	return pp.selectValues.Get(name)
}

// QueryUser queries the "user" edge of the PublisherPermission entity.
func (pp *PublisherPermission) QueryUser() *UserQuery {
	return NewPublisherPermissionClient(pp.config).QueryUser(pp)
}

// QueryPublisher queries the "publisher" edge of the PublisherPermission entity.
func (pp *PublisherPermission) QueryPublisher() *PublisherQuery {
	return NewPublisherPermissionClient(pp.config).QueryPublisher(pp)
}

// Update returns a builder for updating this PublisherPermission.
// Note that you need to call PublisherPermission.Unwrap() before calling this method if this PublisherPermission
// was returned from a transaction, and the transaction was committed or rolled back.
func (pp *PublisherPermission) Update() *PublisherPermissionUpdateOne {
	return NewPublisherPermissionClient(pp.config).UpdateOne(pp)
}

// Unwrap unwraps the PublisherPermission entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pp *PublisherPermission) Unwrap() *PublisherPermission {
	_tx, ok := pp.config.driver.(*txDriver)
	if !ok {
		panic("ent: PublisherPermission is not a transactional entity")
	}
	pp.config.driver = _tx.drv
	return pp
}

// String implements the fmt.Stringer.
func (pp *PublisherPermission) String() string {
	var builder strings.Builder
	builder.WriteString("PublisherPermission(")
	builder.WriteString(fmt.Sprintf("id=%v, ", pp.ID))
	builder.WriteString("permission=")
	builder.WriteString(fmt.Sprintf("%v", pp.Permission))
	builder.WriteString(", ")
	builder.WriteString("user_id=")
	builder.WriteString(pp.UserID)
	builder.WriteString(", ")
	builder.WriteString("publisher_id=")
	builder.WriteString(pp.PublisherID)
	builder.WriteByte(')')
	return builder.String()
}

// PublisherPermissions is a parsable slice of PublisherPermission.
type PublisherPermissions []*PublisherPermission

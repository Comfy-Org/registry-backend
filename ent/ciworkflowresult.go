// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"registry-backend/ent/ciworkflowresult"
	"registry-backend/ent/gitcommit"
	"registry-backend/ent/schema"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

// CIWorkflowResult is the model entity for the CIWorkflowResult schema.
type CIWorkflowResult struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// OperatingSystem holds the value of the "operating_system" field.
	OperatingSystem string `json:"operating_system,omitempty"`
	// WorkflowName holds the value of the "workflow_name" field.
	WorkflowName string `json:"workflow_name,omitempty"`
	// RunID holds the value of the "run_id" field.
	RunID string `json:"run_id,omitempty"`
	// JobID holds the value of the "job_id" field.
	JobID string `json:"job_id,omitempty"`
	// Status holds the value of the "status" field.
	Status schema.WorkflowRunStatusType `json:"status,omitempty"`
	// StartTime holds the value of the "start_time" field.
	StartTime int64 `json:"start_time,omitempty"`
	// EndTime holds the value of the "end_time" field.
	EndTime int64 `json:"end_time,omitempty"`
	// PythonVersion holds the value of the "python_version" field.
	PythonVersion string `json:"python_version,omitempty"`
	// PytorchVersion holds the value of the "pytorch_version" field.
	PytorchVersion string `json:"pytorch_version,omitempty"`
	// CudaVersion holds the value of the "cuda_version" field.
	CudaVersion string `json:"cuda_version,omitempty"`
	// ComfyRunFlags holds the value of the "comfy_run_flags" field.
	ComfyRunFlags string `json:"comfy_run_flags,omitempty"`
	// Average amount of VRAM used by the workflow in Megabytes
	AvgVram int `json:"avg_vram,omitempty"`
	// Peak amount of VRAM used by the workflow in Megabytes
	PeakVram int `json:"peak_vram,omitempty"`
	// User who triggered the job
	JobTriggerUser string `json:"job_trigger_user,omitempty"`
	// Stores miscellaneous metadata for each workflow run.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the CIWorkflowResultQuery when eager-loading is set.
	Edges              CIWorkflowResultEdges `json:"edges"`
	git_commit_results *uuid.UUID
	selectValues       sql.SelectValues
}

// CIWorkflowResultEdges holds the relations/edges for other nodes in the graph.
type CIWorkflowResultEdges struct {
	// Gitcommit holds the value of the gitcommit edge.
	Gitcommit *GitCommit `json:"gitcommit,omitempty"`
	// StorageFile holds the value of the storage_file edge.
	StorageFile []*StorageFile `json:"storage_file,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// GitcommitOrErr returns the Gitcommit value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e CIWorkflowResultEdges) GitcommitOrErr() (*GitCommit, error) {
	if e.Gitcommit != nil {
		return e.Gitcommit, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: gitcommit.Label}
	}
	return nil, &NotLoadedError{edge: "gitcommit"}
}

// StorageFileOrErr returns the StorageFile value or an error if the edge
// was not loaded in eager-loading.
func (e CIWorkflowResultEdges) StorageFileOrErr() ([]*StorageFile, error) {
	if e.loadedTypes[1] {
		return e.StorageFile, nil
	}
	return nil, &NotLoadedError{edge: "storage_file"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*CIWorkflowResult) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case ciworkflowresult.FieldMetadata:
			values[i] = new([]byte)
		case ciworkflowresult.FieldStartTime, ciworkflowresult.FieldEndTime, ciworkflowresult.FieldAvgVram, ciworkflowresult.FieldPeakVram:
			values[i] = new(sql.NullInt64)
		case ciworkflowresult.FieldOperatingSystem, ciworkflowresult.FieldWorkflowName, ciworkflowresult.FieldRunID, ciworkflowresult.FieldJobID, ciworkflowresult.FieldStatus, ciworkflowresult.FieldPythonVersion, ciworkflowresult.FieldPytorchVersion, ciworkflowresult.FieldCudaVersion, ciworkflowresult.FieldComfyRunFlags, ciworkflowresult.FieldJobTriggerUser:
			values[i] = new(sql.NullString)
		case ciworkflowresult.FieldCreateTime, ciworkflowresult.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		case ciworkflowresult.FieldID:
			values[i] = new(uuid.UUID)
		case ciworkflowresult.ForeignKeys[0]: // git_commit_results
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the CIWorkflowResult fields.
func (cwr *CIWorkflowResult) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case ciworkflowresult.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				cwr.ID = *value
			}
		case ciworkflowresult.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field create_time", values[i])
			} else if value.Valid {
				cwr.CreateTime = value.Time
			}
		case ciworkflowresult.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field update_time", values[i])
			} else if value.Valid {
				cwr.UpdateTime = value.Time
			}
		case ciworkflowresult.FieldOperatingSystem:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field operating_system", values[i])
			} else if value.Valid {
				cwr.OperatingSystem = value.String
			}
		case ciworkflowresult.FieldWorkflowName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field workflow_name", values[i])
			} else if value.Valid {
				cwr.WorkflowName = value.String
			}
		case ciworkflowresult.FieldRunID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field run_id", values[i])
			} else if value.Valid {
				cwr.RunID = value.String
			}
		case ciworkflowresult.FieldJobID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field job_id", values[i])
			} else if value.Valid {
				cwr.JobID = value.String
			}
		case ciworkflowresult.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				cwr.Status = schema.WorkflowRunStatusType(value.String)
			}
		case ciworkflowresult.FieldStartTime:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field start_time", values[i])
			} else if value.Valid {
				cwr.StartTime = value.Int64
			}
		case ciworkflowresult.FieldEndTime:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field end_time", values[i])
			} else if value.Valid {
				cwr.EndTime = value.Int64
			}
		case ciworkflowresult.FieldPythonVersion:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field python_version", values[i])
			} else if value.Valid {
				cwr.PythonVersion = value.String
			}
		case ciworkflowresult.FieldPytorchVersion:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field pytorch_version", values[i])
			} else if value.Valid {
				cwr.PytorchVersion = value.String
			}
		case ciworkflowresult.FieldCudaVersion:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field cuda_version", values[i])
			} else if value.Valid {
				cwr.CudaVersion = value.String
			}
		case ciworkflowresult.FieldComfyRunFlags:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field comfy_run_flags", values[i])
			} else if value.Valid {
				cwr.ComfyRunFlags = value.String
			}
		case ciworkflowresult.FieldAvgVram:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field avg_vram", values[i])
			} else if value.Valid {
				cwr.AvgVram = int(value.Int64)
			}
		case ciworkflowresult.FieldPeakVram:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field peak_vram", values[i])
			} else if value.Valid {
				cwr.PeakVram = int(value.Int64)
			}
		case ciworkflowresult.FieldJobTriggerUser:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field job_trigger_user", values[i])
			} else if value.Valid {
				cwr.JobTriggerUser = value.String
			}
		case ciworkflowresult.FieldMetadata:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field metadata", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &cwr.Metadata); err != nil {
					return fmt.Errorf("unmarshal field metadata: %w", err)
				}
			}
		case ciworkflowresult.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field git_commit_results", values[i])
			} else if value.Valid {
				cwr.git_commit_results = new(uuid.UUID)
				*cwr.git_commit_results = *value.S.(*uuid.UUID)
			}
		default:
			cwr.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the CIWorkflowResult.
// This includes values selected through modifiers, order, etc.
func (cwr *CIWorkflowResult) Value(name string) (ent.Value, error) {
	return cwr.selectValues.Get(name)
}

// QueryGitcommit queries the "gitcommit" edge of the CIWorkflowResult entity.
func (cwr *CIWorkflowResult) QueryGitcommit() *GitCommitQuery {
	return NewCIWorkflowResultClient(cwr.config).QueryGitcommit(cwr)
}

// QueryStorageFile queries the "storage_file" edge of the CIWorkflowResult entity.
func (cwr *CIWorkflowResult) QueryStorageFile() *StorageFileQuery {
	return NewCIWorkflowResultClient(cwr.config).QueryStorageFile(cwr)
}

// Update returns a builder for updating this CIWorkflowResult.
// Note that you need to call CIWorkflowResult.Unwrap() before calling this method if this CIWorkflowResult
// was returned from a transaction, and the transaction was committed or rolled back.
func (cwr *CIWorkflowResult) Update() *CIWorkflowResultUpdateOne {
	return NewCIWorkflowResultClient(cwr.config).UpdateOne(cwr)
}

// Unwrap unwraps the CIWorkflowResult entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (cwr *CIWorkflowResult) Unwrap() *CIWorkflowResult {
	_tx, ok := cwr.config.driver.(*txDriver)
	if !ok {
		panic("ent: CIWorkflowResult is not a transactional entity")
	}
	cwr.config.driver = _tx.drv
	return cwr
}

// String implements the fmt.Stringer.
func (cwr *CIWorkflowResult) String() string {
	var builder strings.Builder
	builder.WriteString("CIWorkflowResult(")
	builder.WriteString(fmt.Sprintf("id=%v, ", cwr.ID))
	builder.WriteString("create_time=")
	builder.WriteString(cwr.CreateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("update_time=")
	builder.WriteString(cwr.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("operating_system=")
	builder.WriteString(cwr.OperatingSystem)
	builder.WriteString(", ")
	builder.WriteString("workflow_name=")
	builder.WriteString(cwr.WorkflowName)
	builder.WriteString(", ")
	builder.WriteString("run_id=")
	builder.WriteString(cwr.RunID)
	builder.WriteString(", ")
	builder.WriteString("job_id=")
	builder.WriteString(cwr.JobID)
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(fmt.Sprintf("%v", cwr.Status))
	builder.WriteString(", ")
	builder.WriteString("start_time=")
	builder.WriteString(fmt.Sprintf("%v", cwr.StartTime))
	builder.WriteString(", ")
	builder.WriteString("end_time=")
	builder.WriteString(fmt.Sprintf("%v", cwr.EndTime))
	builder.WriteString(", ")
	builder.WriteString("python_version=")
	builder.WriteString(cwr.PythonVersion)
	builder.WriteString(", ")
	builder.WriteString("pytorch_version=")
	builder.WriteString(cwr.PytorchVersion)
	builder.WriteString(", ")
	builder.WriteString("cuda_version=")
	builder.WriteString(cwr.CudaVersion)
	builder.WriteString(", ")
	builder.WriteString("comfy_run_flags=")
	builder.WriteString(cwr.ComfyRunFlags)
	builder.WriteString(", ")
	builder.WriteString("avg_vram=")
	builder.WriteString(fmt.Sprintf("%v", cwr.AvgVram))
	builder.WriteString(", ")
	builder.WriteString("peak_vram=")
	builder.WriteString(fmt.Sprintf("%v", cwr.PeakVram))
	builder.WriteString(", ")
	builder.WriteString("job_trigger_user=")
	builder.WriteString(cwr.JobTriggerUser)
	builder.WriteString(", ")
	builder.WriteString("metadata=")
	builder.WriteString(fmt.Sprintf("%v", cwr.Metadata))
	builder.WriteByte(')')
	return builder.String()
}

// CIWorkflowResults is a parsable slice of CIWorkflowResult.
type CIWorkflowResults []*CIWorkflowResult

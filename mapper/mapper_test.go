package mapper

import (
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/schema"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestMachineStatsToMap(t *testing.T) {

	vramTimeSeries := map[string]interface{}{
		"timeSeries": []float64{0.1, 0.2, 0.3},
	}
	// Test case 1: Valid machine stats object
	ms := &drip.MachineStats{
		CpuCapacity:    proto.String("4"),
		DiskCapacity:   proto.String("1024"),
		InitialCpu:     proto.String("2"),
		InitialDisk:    proto.String("512"),
		InitialRam:     proto.String("4096"),
		MemoryCapacity: proto.String("8192"),
		OsVersion:      proto.String("Ubuntu 20.04"),
		PipFreeze:      proto.String("package1==1.0.0\npackage2==2.0.0"),
		VramTimeSeries: &vramTimeSeries,
		MachineName:    proto.String("TestMachine"),
		GpuType:        proto.String("NVIDIA GeForce RTX 3080"),
	}

	expected := map[string]interface{}{
		"CpuCapacity":    proto.String("4"),
		"DiskCapacity":   proto.String("1024"),
		"InitialCpu":     proto.String("2"),
		"InitialDisk":    proto.String("512"),
		"InitialRam":     proto.String("4096"),
		"MemoryCapacity": proto.String("8192"),
		"OsVersion":      proto.String("Ubuntu 20.04"),
		"PipFreeze":      proto.String("package1==1.0.0\npackage2==2.0.0"),
		"VramTimeSeries": &vramTimeSeries,
		"MachineName":    proto.String("TestMachine"),
		"GpuType":        proto.String("NVIDIA GeForce RTX 3080"),
	}

	actual := MachineStatsToMap(ms)
	assert.Equal(t, expected, actual, "Expected and actual maps should be equal")

	// Test case 2: Nil machine stats object
	actual = MachineStatsToMap(nil)
	assert.Nil(t, actual, "Expected nil when machine stats is nil")
}

func TestCiWorkflowResultToActionJobResult(t *testing.T) {
	// Helper function to create a valid CIWorkflowResult for testing
	createValidCIWorkflowResult := func() *ent.CIWorkflowResult {
		id := uuid.New()
		commitID := uuid.New()
		now := time.Now()
		commitTimestamp := now

		return &ent.CIWorkflowResult{
			ID:              id,
			OperatingSystem: "Ubuntu 20.04",
			WorkflowName:    "Test Workflow",
			RunID:           "run-id",
			JobID:           "job-id",
			Status:          schema.WorkflowRunStatusTypeCompleted,
			StartTime:       now.Unix(),
			EndTime:         now.Add(1 * time.Hour).Unix(),
			PythonVersion:   "3.8.5",
			PytorchVersion:  "1.7.1",
			CudaVersion:     "11.0",
			ComfyRunFlags:   "test-flags",
			AvgVram:         2048,
			PeakVram:        4096,
			JobTriggerUser:  "test-user",
			Metadata: map[string]interface{}{
				"cpu": 4,
				"ram": 8192,
			},
			Edges: ent.CIWorkflowResultEdges{
				Gitcommit: &ent.GitCommit{
					ID:              commitID,
					CommitTimestamp: commitTimestamp,
					CommitHash:      "abc123",
					CommitMessage:   "Initial commit",
					RepoName:        "test-repo",
					PrNumber:        "123",
					Author:          "test-author",
				},
				StorageFile: []*ent.StorageFile{
					{
						FileURL: "http://example.com/file",
					},
				},
			},
		}
	}

	t.Run("Valid Workflow Result", func(t *testing.T) {
		result := createValidCIWorkflowResult()
		actionJobResult, err := CiWorkflowResultToActionJobResult(result)

		assert.NoError(t, err)
		assert.NotNil(t, actionJobResult)
		assert.Equal(t, result.ID, *actionJobResult.Id)
		assert.Equal(t, result.WorkflowName, *actionJobResult.WorkflowName)
		assert.Equal(t, result.OperatingSystem, *actionJobResult.OperatingSystem)
		assert.Equal(t, result.PythonVersion, *actionJobResult.PythonVersion)
		assert.Equal(t, result.PytorchVersion, *actionJobResult.PytorchVersion)
		assert.Equal(t, result.CudaVersion, *actionJobResult.CudaVersion)
		assert.Equal(t, result.Edges.Gitcommit.BranchName, *actionJobResult.BranchName)
		assert.Equal(t, "http://example.com/file", *actionJobResult.StorageFile.PublicUrl)
		assert.Equal(t, result.Edges.Gitcommit.CommitHash, *actionJobResult.CommitHash)
		assert.Equal(t, result.Edges.Gitcommit.ID.String(), *actionJobResult.CommitId)
		assert.Equal(t, result.Edges.Gitcommit.CommitTimestamp.Unix(), *actionJobResult.CommitTime)
		assert.Equal(t, result.Edges.Gitcommit.CommitMessage, *actionJobResult.CommitMessage)
		assert.Equal(t, result.Edges.Gitcommit.RepoName, *actionJobResult.GitRepo)
		assert.Equal(t, result.RunID, *actionJobResult.ActionRunId)
		assert.Equal(t, result.JobID, *actionJobResult.ActionJobId)
		assert.Equal(t, result.StartTime, *actionJobResult.StartTime)
		assert.Equal(t, result.EndTime, *actionJobResult.EndTime)
		assert.Equal(t, result.JobTriggerUser, *actionJobResult.JobTriggerUser)
		assert.Equal(t, result.AvgVram, *actionJobResult.AvgVram)
		assert.Equal(t, result.PeakVram, *actionJobResult.PeakVram)
		assert.Equal(t, result.ComfyRunFlags, *actionJobResult.ComfyRunFlags)
		assert.Equal(t, result.Edges.Gitcommit.PrNumber, *actionJobResult.PrNumber)
		assert.Equal(t, result.Edges.Gitcommit.Author, *actionJobResult.Author)
	})

	t.Run("Nil Workflow Result", func(t *testing.T) {
		actionJobResult, err := CiWorkflowResultToActionJobResult(nil)

		assert.Nil(t, actionJobResult)
		assert.NoError(t, err)
	})

	t.Run("Missing Storage File", func(t *testing.T) {
		result := createValidCIWorkflowResult()
		result.Edges.StorageFile = nil
		actionJobResult, err := CiWorkflowResultToActionJobResult(result)

		assert.NoError(t, err)
		assert.NotNil(t, actionJobResult)
		assert.Nil(t, actionJobResult.StorageFile)
	})

	t.Run("Missing Git Commit", func(t *testing.T) {
		result := createValidCIWorkflowResult()
		result.Edges.Gitcommit = nil
		_, err := CiWorkflowResultToActionJobResult(result)

		assert.Error(t, err)
	})

	t.Run("Invalid Metadata", func(t *testing.T) {
		result := createValidCIWorkflowResult()
		result.Metadata = nil // assuming MapToMachineStats handles nil correctly
		actionJobResult, err := CiWorkflowResultToActionJobResult(result)

		assert.NoError(t, err)
		assert.NotNil(t, actionJobResult)
		assert.Nil(t, actionJobResult.MachineStats)
	})

	t.Run("Status Conversion Error", func(t *testing.T) {
		result := createValidCIWorkflowResult()
		result.Status = "invalid_status"
		_, err := CiWorkflowResultToActionJobResult(result)

		assert.Error(t, err)
	})
}

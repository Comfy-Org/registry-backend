package mapper

import (
	"fmt"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/schema"
)

func CiWorkflowResultsToActionJobResults(results []*ent.CIWorkflowResult) ([]drip.ActionJobResult, error) {
	var jobResultsData []drip.ActionJobResult

	for _, result := range results {
		jobResultData, err := CiWorkflowResultToActionJobResult(result)
		if err != nil {
			return nil, err
		}
		jobResultsData = append(jobResultsData, *jobResultData)
	}
	return jobResultsData, nil
}

func CiWorkflowResultToActionJobResult(result *ent.CIWorkflowResult) (*drip.ActionJobResult, error) {
	var storageFileData *drip.StorageFile

	// Check if the StorageFile slice is not empty before accessing
	if len(result.Edges.StorageFile) > 0 {
		storageFileData = &drip.StorageFile{
			PublicUrl: &result.Edges.StorageFile[0].FileURL,
		}
	}
	commitId := result.Edges.Gitcommit.ID.String()
	commitUnixTime := result.Edges.Gitcommit.CommitTimestamp.Unix()
	apiStatus, err := DbWorkflowRunStatusToApi(result.Status)
	if err != nil {
		return nil, err
	}

	machineStats, err := MapToMachineStats(result.Metadata)
	if err != nil {
		return nil, err
	}

	return &drip.ActionJobResult{
		Id:              &result.ID,
		WorkflowName:    &result.WorkflowName,
		OperatingSystem: &result.OperatingSystem,
		PythonVersion:   &result.PythonVersion,
		PytorchVersion:  &result.PytorchVersion,
		CudaVersion:     &result.CudaVersion,
		StorageFile:     storageFileData,
		CommitHash:      &result.Edges.Gitcommit.CommitHash,
		CommitId:        &commitId,
		CommitTime:      &commitUnixTime,
		CommitMessage:   &result.Edges.Gitcommit.CommitMessage,
		GitRepo:         &result.Edges.Gitcommit.RepoName,
		ActionRunId:     &result.RunID,
		ActionJobId:     &result.JobID,
		StartTime:       &result.StartTime,
		EndTime:         &result.EndTime,
		JobTriggerUser:  &result.JobTriggerUser,
		AvgVram:         &result.AvgVram,
		PeakVram:        &result.PeakVram,
		ComfyRunFlags:   &result.ComfyRunFlags,

		Status:       &apiStatus,
		PrNumber:     &result.Edges.Gitcommit.PrNumber,
		Author:       &result.Edges.Gitcommit.Author,
		MachineStats: machineStats,
	}, nil
}

func ApiWorkflowRunStatusToDb(status drip.WorkflowRunStatus) (schema.WorkflowRunStatusType, error) {
	switch status {
	case drip.WorkflowRunStatusStarted:
		return schema.WorkflowRunStatusTypeStarted, nil
	case drip.WorkflowRunStatusCompleted:
		return schema.WorkflowRunStatusTypeCompleted, nil
	case drip.WorkflowRunStatusFailed:
		return schema.WorkflowRunStatusTypeFailed, nil
	default:
		// Throw an error
		return "", fmt.Errorf("unsupported workflow status: %v", status)

	}
}

func DbWorkflowRunStatusToApi(status schema.WorkflowRunStatusType) (drip.WorkflowRunStatus, error) {
	switch status {
	case schema.WorkflowRunStatusTypeStarted:
		return drip.WorkflowRunStatusStarted, nil
	case schema.WorkflowRunStatusTypeCompleted:
		return drip.WorkflowRunStatusCompleted, nil
	case schema.WorkflowRunStatusTypeFailed:
		return drip.WorkflowRunStatusFailed, nil
	default:
		// Throw an error
		return "", fmt.Errorf("unsupported workflow status: %v", status)
	}
}

func MachineStatsToMap(ms *drip.MachineStats) map[string]interface{} {
	return map[string]interface{}{
		"CpuCapacity":    ms.CpuCapacity,
		"DiskCapacity":   ms.DiskCapacity,
		"InitialCpu":     ms.InitialCpu,
		"InitialDisk":    ms.InitialDisk,
		"InitialRam":     ms.InitialRam,
		"MemoryCapacity": ms.MemoryCapacity,
		"OsVersion":      ms.OsVersion,
		"PipFreeze":      ms.PipFreeze,
		"VramTimeSeries": ms.VramTimeSeries,
		"MachineName":    ms.MachineName,
		"GpuType":        ms.GpuType,
	}
}

func MapToMachineStats(data map[string]interface{}) (*drip.MachineStats, error) {
	var ms drip.MachineStats

	if data == nil {
		return nil, nil
	}

	// Helper function to get string pointers from the map
	getStringPtr := func(key string) *string {
		if val, exists := data[key]; exists {
			if strVal, ok := val.(string); ok {
				return &strVal
			}
		}
		return nil // Return nil if the key does not exist or type assertion fails
	}

	ms.CpuCapacity = getStringPtr("CpuCapacity")
	ms.DiskCapacity = getStringPtr("DiskCapacity")
	ms.InitialCpu = getStringPtr("InitialCpu")
	ms.InitialDisk = getStringPtr("InitialDisk")
	ms.InitialRam = getStringPtr("InitialRam")
	ms.MemoryCapacity = getStringPtr("MemoryCapacity")
	ms.OsVersion = getStringPtr("OsVersion")
	ms.PipFreeze = getStringPtr("PipFreeze")
	ms.MachineName = getStringPtr("MachineName")
	ms.GpuType = getStringPtr("GpuType")

	if val, exists := data["VramTimeSeries"]; exists {
		if vram, ok := val.(map[string]interface{}); ok {
			ms.VramTimeSeries = &vram
		}
	}

	return &ms, nil
}

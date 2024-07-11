package mapper

import (
	"fmt"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/schema"
)

func CiWorkflowResultToActionJobResult(results []*ent.CIWorkflowResult) ([]drip.ActionJobResult, error) {
	var jobResultsData []drip.ActionJobResult

	for _, result := range results {
		storageFileData := drip.StorageFile{
			PublicUrl: &result.Edges.StorageFile.FileURL,
		}
		commitId := result.Edges.Gitcommit.ID.String()
		commitUnixTime := result.Edges.Gitcommit.CommitTimestamp.Unix()
		apiStatus, err := DbWorkflowRunStatusToApi(result.Status)
		if err != nil {
			return jobResultsData, err
		}
		jobResultData := drip.ActionJobResult{
			WorkflowName:    &result.WorkflowName,
			OperatingSystem: &result.OperatingSystem,
			GpuType:         &result.GpuType,
			PytorchVersion:  &result.PytorchVersion,
			StorageFile:     &storageFileData,
			CommitHash:      &result.Edges.Gitcommit.CommitHash,
			CommitId:        &commitId,
			CommitTime:      &commitUnixTime,
			CommitMessage:   &result.Edges.Gitcommit.CommitMessage,
			GitRepo:         &result.Edges.Gitcommit.RepoName,
			ActionRunId:     &result.RunID,
			StartTime:       &result.StartTime,
			EndTime:         &result.EndTime,
			JobTriggerUser:  &result.JobTriggerUser,
			AvgVram:         &result.AvgVram,
			PeakVram:        &result.PeakVram,
			PythonVersion:   &result.PythonVersion,
			Status:          &apiStatus,
			PrNumber:        &result.Edges.Gitcommit.PrNumber,
			Author:          &result.Edges.Gitcommit.Author,
		}
		jobResultsData = append(jobResultsData, jobResultData)
	}
	return jobResultsData, nil
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

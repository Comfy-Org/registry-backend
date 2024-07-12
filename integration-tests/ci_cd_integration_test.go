package integration

import (
	"context"
	"fmt"
	"registry-backend/config"
	"registry-backend/drip"
	"registry-backend/ent/gitcommit"
	"registry-backend/mock/gateways"
	"registry-backend/server/implementation"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestCICD(t *testing.T) {
	clientCtx := context.Background()
	client, cleanup := setupDB(t, clientCtx)
	defer cleanup()

	// Initialize the Service
	mockStorageService := new(gateways.MockStorageService)
	mockSlackService := new(gateways.MockSlackService)
	mockDiscordService := new(gateways.MockDiscordService)
	mockSlackService.
		On("SendRegistryMessageToSlack", mock.Anything).
		Return(nil) // Do nothing for all slack messsage calls.
	mockAlgolia := new(gateways.MockAlgoliaService)
	mockAlgolia.
		On("IndexNodes", mock.Anything, mock.Anything).
		Return(nil)
	impl := implementation.NewStrictServerImplementation(
		client, &config.Config{}, mockStorageService, mockSlackService, mockDiscordService, mockAlgolia)

	ctx := context.Background()
	now := time.Now()
	anHourAgo := now.Add(-1 * time.Hour)
	avgVram := 2132

	body := &drip.PostUploadArtifactJSONRequestBody{
		Repo:                "github.com/comfy/service",
		BranchName:          "develop",
		CommitHash:          "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		CommitMessage:       "new commit",
		CommitTime:          anHourAgo.Format(time.RFC3339),
		JobId:               "018fbe20-88a6-7d31-a194-eee8e2509da3",
		RunId:               "018fbe37-a7a8-74a3-a377-8d70d54f54d8",
		Os:                  "linux",
		WorkflowName:        "devops",
		CudaVersion:         proto.String("1.0.0"),
		BucketName:          proto.String("comfy-dev-bucket"),
		OutputFilesGcsPaths: proto.String("comfy-dev-file"),
		ComfyLogsGcsPath:    proto.String("comfy-dev-log"),
		StartTime:           anHourAgo.Unix(),
		EndTime:             now.Unix(),
		PrNumber:            "123",
		PythonVersion:       "3.8",
		PytorchVersion:      proto.String("1.0.0"),
		JobTriggerUser:      "comfy",
		Author:              "robin",
		AvgVram:             &avgVram,
		ComfyRunFlags:       proto.String("comfy"),
		Status:              drip.WorkflowRunStatusStarted,
		MachineStats: &drip.MachineStats{
			CpuCapacity:    proto.String("2.0"),
			InitialCpu:     proto.String("1.0"),
			InitialDisk:    proto.String("1.0"),
			DiskCapacity:   proto.String("2.0"),
			InitialRam:     proto.String("1.0"),
			MemoryCapacity: proto.String("2.0"),
			OsVersion:      proto.String("Ubuntu 24.10"),
			PipFreeze:      proto.String("requests==1.0.0"),
			MachineName:    proto.String("comfy-dev"),
			GpuType:        proto.String("NVIDIA Tesla V100"),
		},
	}

	t.Run("Post Upload Artifact", func(t *testing.T) {
		body := *body
		body.JobId = "018fbe4a-2844-7c2e-87f1-311605292452"
		body.RunId = "018fbe4a-5b1c-7a51-8e26-53e77961ee06"
		res, err := impl.PostUploadArtifact(ctx, drip.PostUploadArtifactRequestObject{Body: &body})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.PostUploadArtifact200JSONResponse{}, res, "should return 200")
	})

	t.Run("Re Post Upload Artifact", func(t *testing.T) {
		res, err := impl.PostUploadArtifact(ctx, drip.PostUploadArtifactRequestObject{Body: body})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.PostUploadArtifact200JSONResponse{}, res, "should return 200")
	})

	t.Run("Get Git Commit", func(t *testing.T) {
		expectedAvgVram := 2132
		expectedPeakVram := 0
		git, err := client.GitCommit.Query().Where(gitcommit.CommitHashEQ(body.CommitHash)).First(ctx)
		require.NoError(t, err)

		res, err := impl.GetGitcommit(ctx, drip.GetGitcommitRequestObject{Params: drip.GetGitcommitParams{
			CommitId:        proto.String(git.ID.String()),
			OperatingSystem: &body.Os,
			WorkflowName:    &body.WorkflowName,
			Branch:          &body.BranchName,
			RepoName:        &body.Repo,
		}})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.GetGitcommit200JSONResponse{}, res)
		res200 := res.(drip.GetGitcommit200JSONResponse)
		require.Len(t, *res200.JobResults, 1)
		assert.Equal(t, *res200.TotalNumberOfPages, 1)
		assert.Equal(t, drip.ActionJobResult{
			Id:              (*res200.JobResults)[0].Id,
			ActionRunId:     &body.RunId,
			ActionJobId:     &body.JobId,
			CommitHash:      &body.CommitHash,
			CommitId:        proto.String(git.ID.String()),
			CommitMessage:   &body.CommitMessage,
			CommitTime:      proto.Int64(anHourAgo.Unix()),
			EndTime:         proto.Int64(now.Unix()),
			GitRepo:         &body.Repo,
			OperatingSystem: &body.Os,
			StartTime:       proto.Int64(anHourAgo.Unix()),
			WorkflowName:    &body.WorkflowName,
			JobTriggerUser:  &body.JobTriggerUser,
			AvgVram:         &expectedAvgVram,
			PeakVram:        &expectedPeakVram,
			PythonVersion:   &body.PythonVersion,
			Status:          &body.Status,
			PrNumber:        &body.PrNumber,
			CudaVersion:     body.CudaVersion,
			PytorchVersion:  body.PytorchVersion,
			Author:          &body.Author,
			ComfyRunFlags:   body.ComfyRunFlags,
			StorageFile: &drip.StorageFile{
				PublicUrl: proto.String(fmt.Sprintf("https://storage.googleapis.com/%s/%s", *body.BucketName, *body.OutputFilesGcsPaths)),
			},
			MachineStats: body.MachineStats,
		}, (*res200.JobResults)[0])
	})

	t.Run("Get invalid Git Commit", func(t *testing.T) {
		fakeID, _ := uuid.NewV7()
		res, err := impl.GetGitcommit(ctx, drip.GetGitcommitRequestObject{Params: drip.GetGitcommitParams{
			CommitId: proto.String(fakeID.String())}})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.GetGitcommit200JSONResponse{}, res)
		assert.Len(t, *res.(drip.GetGitcommit200JSONResponse).JobResults, 0)
	})

	t.Run("Get Branch", func(t *testing.T) {
		res, err := impl.GetBranch(ctx, drip.GetBranchRequestObject{Params: drip.GetBranchParams{
			RepoName: body.Repo,
		}})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.GetBranch200JSONResponse{}, res)
		res200 := res.(drip.GetBranch200JSONResponse)
		require.Len(t, *res200.Branches, 1, "should return corrent number of branches")
		assert.Equal(t, body.BranchName, (*res200.Branches)[0], "should return correct branches")
	})

	t.Run("Get invalid branch", func(t *testing.T) {
		res, err := impl.GetBranch(ctx, drip.GetBranchRequestObject{Params: drip.GetBranchParams{
			RepoName: "notexist",
		}})
		require.NoError(t, err, "should return error")
		assert.IsType(t, drip.GetBranch200JSONResponse{}, res)
		assert.Len(t, *res.(drip.GetBranch200JSONResponse).Branches, 0, "should return empty branch")
	})
}

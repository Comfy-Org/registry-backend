package dripservices

import (
	"fmt"
	"registry-backend/config"
	"registry-backend/db"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/ciworkflowresult"
	"registry-backend/ent/gitcommit"
	"registry-backend/server/middleware/metric"
	"strings"
	"time"

	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// ComfyCIService provides methods to interact with CI-related data in the database.
type ComfyCIService struct {
	Config *config.Config
}

// NewComfyCIService creates a new instance of ComfyCIService.
func NewComfyCIService(config *config.Config) *ComfyCIService {
	return &ComfyCIService{
		Config: config,
	}
}

// ProcessCIRequest handles the incoming request and creates/updates the necessary entities.
func (s *ComfyCIService) ProcessCIRequest(ctx context.Context, client *ent.Client, req *drip.PostUploadArtifactRequestObject) error {
	// Check if git commit exists
	// If it does, remove all CiWorkflowRuns associated with it.
	existingCommit, err := client.GitCommit.Query().Where(gitcommit.CommitHashEQ(req.Body.CommitHash)).Where(gitcommit.RepoNameEQ(req.Body.Repo)).Only(ctx)
	if ent.IsNotSingular(err) {
		log.Ctx(ctx).Error().Err(err).Msgf("Failed to query git commit %s", req.Body.CommitHash)
		drip_metric.IncrementCustomCounterMetric(ctx, drip_metric.CustomCounterIncrement{
			Type:   "ci-git-commit-query-error",
			Val:    1,
			Labels: map[string]string{},
		})
	}
	if existingCommit != nil {
		_, err := client.CIWorkflowResult.Delete().Where(ciworkflowresult.HasGitcommitWith(gitcommit.IDEQ(existingCommit.ID))).Exec(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msgf("Failed to delete existing run results for git commit %s", req.Body.CommitHash)
			return err
		}
	}

	return db.WithTx(ctx, client, func(tx *ent.Tx) error {
		id, err := s.UpsertCommit(ctx, tx.Client(), req.Body.CommitHash, req.Body.BranchName, req.Body.Repo, req.Body.CommitTime, req.Body.CommitMessage)
		if err != nil {
			return err
		}
		gitcommit := tx.Client().GitCommit.GetX(ctx, id)
		if req.Body.OutputFilesGcsPaths != nil && req.Body.BucketName != nil {
			files, err := GetPublicUrlForOutputFiles(ctx, *req.Body.BucketName, *req.Body.OutputFilesGcsPaths)
			if err != nil {
				return err
			}

			for _, file := range files {
				// TODO(robinhuang): Get real filetype.
				file, err := s.UpsertStorageFile(ctx, tx.Client(), file.PublicURL, file.BucketName, file.FilePath, "image")

				if err != nil {
					log.Ctx(ctx).Error().Err(err).Msg("Failed to upsert storage file")
					drip_metric.IncrementCustomCounterMetric(ctx, drip_metric.CustomCounterIncrement{
						Type: "ci-upsert-storage-error",
						Val:  1,
						Labels: map[string]string{
							"bucket-name": file.BucketName,
						},
					})
					continue
				}

				cudaVersion := ""
				if req.Body.CudaVersion != nil {
					cudaVersion = *req.Body.CudaVersion
				}

				_, err = s.UpsertRunResult(ctx, tx.Client(), file, gitcommit, req.Body.Os, cudaVersion, req.Body.WorkflowName, req.Body.RunId, req.Body.StartTime, req.Body.EndTime)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// UpsertCommit creates or updates a GitCommit entity.
func (s *ComfyCIService) UpsertCommit(ctx context.Context, client *ent.Client, hash, branchName, repoName string, commitIsoTime string, commitMessage string) (uuid.UUID, error) {
	log.Ctx(ctx).Info().Msgf("Upserting commit %s", hash)
	commitTime, err := time.Parse(time.RFC3339, commitIsoTime)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := client.GitCommit.
		Create().
		SetCommitHash(hash).
		SetBranchName(branchName).
		SetRepoName(strings.ToLower(repoName)). // TODO(robinhuang): Write test for this.
		SetCommitTimestamp(commitTime).
		SetCommitMessage(commitMessage).
		OnConflict(
			// Careful, the order matters here.
			sql.ConflictColumns(gitcommit.FieldRepoName, gitcommit.FieldCommitHash),
		).
		UpdateNewValues().
		ID(ctx)

	if err != nil {
		return uuid.Nil, fmt.Errorf("GitCommit.Create: %w", err)
	}
	return id, nil
}

// UpsertRunResult creates or updates a ActionRunResult entity.
func (s *ComfyCIService) UpsertRunResult(ctx context.Context, client *ent.Client, file *ent.StorageFile, gitcommit *ent.GitCommit, os, gpuType, workflowName, runId string, startTime, endTime int64) (uuid.UUID, error) {
	log.Ctx(ctx).Info().Msgf("Upserting workflow result for commit %s", gitcommit.CommitHash)
	return client.CIWorkflowResult.
		Create().
		SetGitcommit(gitcommit).
		SetStorageFile(file).
		SetOperatingSystem(os).
		SetWorkflowName(workflowName).
		SetRunID(runId).
		SetStartTime(startTime).
		SetEndTime(endTime).
		OnConflict(
			sql.ConflictColumns(ciworkflowresult.FieldID),
		).
		UpdateNewValues().
		ID(ctx)
}

// UpsertStorageFile creates or updates a RunFile entity.
func (s *ComfyCIService) UpsertStorageFile(ctx context.Context, client *ent.Client, publicUrl, bucketName, filePath, fileType string) (*ent.StorageFile, error) {
	log.Ctx(ctx).Info().Msgf("Upserting storage file for URL %s", publicUrl)
	return client.StorageFile.
		Create().
		SetFileURL(publicUrl).
		SetFilePath(filePath).
		SetBucketName(bucketName).
		SetFileType(fileType).
		Save(ctx)
}

type ObjectInfo struct {
	BucketName string
	FilePath   string
	PublicURL  string
}

// GetPublicUrlForOutputFiles downloads the artifact, extracts it, and uploads each file to GCS
func GetPublicUrlForOutputFiles(ctx context.Context, bucketName, objects string) ([]ObjectInfo, error) {
	objectArr := strings.Split(objects, ",")
	var result []ObjectInfo
	for _, object := range objectArr {
		publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, object)
		log.Ctx(ctx).Info().Msgf("Public URL: %v", publicURL)
		result = append(result, ObjectInfo{
			BucketName: bucketName,
			FilePath:   object,
			PublicURL:  publicURL,
		})
	}
	return result, nil
}

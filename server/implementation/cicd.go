package implementation

import (
	"context"
	"registry-backend/drip"
	"registry-backend/ent/ciworkflowresult"
	"registry-backend/ent/gitcommit"
	"registry-backend/mapper"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (impl *DripStrictServerImplementation) GetGitcommit(ctx context.Context, request drip.GetGitcommitRequestObject) (drip.GetGitcommitResponseObject, error) {
	var commitId uuid.UUID = uuid.Nil
	if request.Params.CommitId != nil {
		log.Ctx(ctx).Info().Msgf("getting commit data for %s", *request.Params.CommitId)
		commitId = uuid.MustParse(*request.Params.CommitId)
	}

	if request.Params.OperatingSystem != nil {
		log.Ctx(ctx).Info().Msgf("getting commit data for %s", *request.Params.OperatingSystem)
	}

	repoName := "comfyanonymous/ComfyUI"
	if request.Params.RepoName != nil {
		repoName = *request.Params.RepoName
	}
	repoName = strings.ToLower(repoName)

	var operatingSystem string
	if request.Params.OperatingSystem != nil {
		operatingSystem = *request.Params.OperatingSystem
	} else {
		operatingSystem = "" // Assign a default value if nil
	}
	var branchName string
	if request.Params.Branch != nil {
		branchName = *request.Params.Branch
	} else {
		branchName = "" // Assign a default value if nil
	}

	var workflowName string
	if request.Params.WorkflowName != nil {
		workflowName = *request.Params.WorkflowName
	} else {
		workflowName = "" // Assign a default value if nil
	}
	log.Ctx(ctx).Info().Msgf("Querying database...")

	query := impl.Client.CIWorkflowResult.Query().
		WithGitcommit().
		WithStorageFile()

	query.Where(ciworkflowresult.HasGitcommitWith(gitcommit.RepoNameEQ(repoName)))
	query.Order(ciworkflowresult.ByGitcommitField(gitcommit.FieldCommitTimestamp, sql.OrderDesc()))
	log.Ctx(ctx).Info().Msgf("Filtering git commit by repo name %s", repoName)

	// Conditionally add the commitId filter
	if commitId != uuid.Nil {
		log.Ctx(ctx).Info().Msgf("Filtering git commit by db commit id %s", commitId)
		query.Where(ciworkflowresult.HasGitcommitWith(gitcommit.IDEQ(commitId)))
	}

	if branchName != "" {
		log.Ctx(ctx).Info().Msgf("Filtering git commit by branch %s", branchName)
		query.Where(ciworkflowresult.HasGitcommitWith(gitcommit.BranchNameEQ(branchName)))
	}

	// Continue building the query
	if operatingSystem != "" {
		log.Ctx(ctx).Info().Msgf("Filtering git commit by OS %s", operatingSystem)
		query.Where(ciworkflowresult.OperatingSystemEQ(operatingSystem))
	}
	if workflowName != "" {
		log.Ctx(ctx).Info().Msgf("Filtering git commit by workflow name %s", workflowName)
		query.Where(ciworkflowresult.WorkflowNameEQ(workflowName))
	}

	// Get total number of pages
	count, err := query.Count(ctx)
	log.Ctx(ctx).Info().Msgf("Got %d runs", count)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Error retrieving count of git commits w/ err: %v", err)
		return drip.GetGitcommit500Response{}, err
	}

	// Pagination
	page := 1
	pageSize := 10
	if request.Params.Page != nil {
		page = *request.Params.Page
	}
	if request.Params.PageSize != nil {
		pageSize = *request.Params.PageSize
	}
	query.Offset((page - 1) * pageSize).Limit(pageSize)

	numberOfPages := (count + pageSize - 1) / pageSize

	// Execute the query
	runs, err := query.All(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Error retrieving git commits w/ err: %v", err)
		return drip.GetGitcommit500Response{}, err
	}

	results, err := mapper.CiWorkflowResultsToActionJobResults(runs)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Error mapping git commits to action job results w/ err: %v", err)
		return drip.GetGitcommit500Response{}, err
	}
	log.Ctx(ctx).Info().Msgf("Git commits retrieved successfully")
	return drip.GetGitcommit200JSONResponse{
		JobResults:         &results,
		TotalNumberOfPages: &numberOfPages,
	}, nil
}

func (impl *DripStrictServerImplementation) GetWorkflowResult(ctx context.Context, request drip.GetWorkflowResultRequestObject) (drip.GetWorkflowResultResponseObject, error) {
	log.Ctx(ctx).Info().Msgf("Getting workflow result with ID %s", request.WorkflowResultId)
	workflowId := uuid.MustParse(request.WorkflowResultId)
	workflow, err := impl.Client.CIWorkflowResult.Query().WithGitcommit().WithStorageFile().Where(ciworkflowresult.IDEQ(workflowId)).First(ctx)

	if err != nil {
		log.Ctx(ctx).Error().Msgf("Error retrieving workflow result w/ err: %v", err)
		return drip.GetWorkflowResult500JSONResponse{
			Message: err.Error(),
		}, nil
	}

	result, err := mapper.CiWorkflowResultToActionJobResult(workflow)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Error mapping workflow result to action job result w/ err: %v", err)
		return drip.GetWorkflowResult500JSONResponse{
			Message: err.Error(),
		}, nil
	}

	log.Ctx(ctx).Info().Msgf("Workflow result retrieved successfully")
	return drip.GetWorkflowResult200JSONResponse(*result), nil
}

func (impl *DripStrictServerImplementation) GetBranch(ctx context.Context, request drip.GetBranchRequestObject) (drip.GetBranchResponseObject, error) {
	repoNameFilter := strings.ToLower(request.Params.RepoName)

	branches, err := impl.Client.GitCommit.
		Query().
		Where(gitcommit.RepoNameEQ(repoNameFilter)).
		GroupBy(gitcommit.FieldBranchName).
		Strings(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Error retrieving git's branchs w/ err: %v", err)
		return drip.GetBranch500Response{}, err
	}

	log.Ctx(ctx).Info().Msgf("Git branches from '%s' repo retrieved successfully", request.Params.RepoName)
	return drip.GetBranch200JSONResponse{Branches: &branches}, nil
}

func (impl *DripStrictServerImplementation) PostUploadArtifact(ctx context.Context, request drip.PostUploadArtifactRequestObject) (drip.PostUploadArtifactResponseObject, error) {
	err := impl.ComfyCIService.ProcessCIRequest(ctx, impl.Client, &request)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Error processiong CI request w/ err: %v", err)
		return drip.PostUploadArtifact500Response{}, err
	}

	log.Ctx(ctx).Info().Msgf("CI request with job id '%s' processed successfully", request.Body.JobId)
	return drip.PostUploadArtifact200JSONResponse{}, nil
}

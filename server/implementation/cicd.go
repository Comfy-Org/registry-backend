package implementation

import (
	"context"
	"fmt"
	"registry-backend/drip"
	"registry-backend/ent/ciworkflowresult"
	"registry-backend/ent/gitcommit"
	"registry-backend/ent/schema"
	"registry-backend/mapper"
	"registry-backend/tracing"
	"sort"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (impl *DripStrictServerImplementation) GetGitcommit(ctx context.Context, request drip.GetGitcommitRequestObject) (drip.GetGitcommitResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.GetGitcommit")()

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

func (impl *DripStrictServerImplementation) GetGitcommitsummary(ctx context.Context, request drip.GetGitcommitsummaryRequestObject) (drip.GetGitcommitsummaryResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.GetGitcommitsummary")()
	log.Ctx(ctx).Info().Msg("Getting git commit summary")

	// Prep relevant vars
	page := 1
	pageSize := 10
	if request.Params.Page != nil {
		page = *request.Params.Page
	}
	if request.Params.PageSize != nil {
		pageSize = *request.Params.PageSize
	}
	targetCount := (page + 2) * pageSize // Note: +1 for current page, +1 more to avoid partial result at end
	summaries := make(map[string]*drip.GitCommitSummary)
	commitsThusFar := 0
	readCount := 4096
	reachedEnd := false

	for len(summaries) < targetCount {

		// TODO: This is implementing as a dense processing at time of request, but this should really be preprocessed in advance
		// probably even stored in the database
		query := impl.Client.CIWorkflowResult.Query().
			WithGitcommit().
			WithStorageFile()

		// Apply filters
		repoName := "comfyanonymous/ComfyUI"
		if request.Params.RepoName != nil {
			repoName = *request.Params.RepoName
		}
		repoName = strings.ToLower(repoName)
		if request.Params.RepoName != nil {
			if commitsThusFar == 0 {
				log.Ctx(ctx).Info().Msgf("Filtering git commit by repo name %s", repoName)
			}
			query.Where(ciworkflowresult.HasGitcommitWith(gitcommit.RepoNameEQ(repoName)))
		}
		if request.Params.BranchName != nil {
			if commitsThusFar == 0 {
				log.Ctx(ctx).Info().Msgf("Filtering git commit by branch name %s", *request.Params.BranchName)
			}
			query.Where(ciworkflowresult.HasGitcommitWith(gitcommit.BranchNameEQ(*request.Params.BranchName)))
		}
		query.Order(ciworkflowresult.ByGitcommitField(gitcommit.FieldCommitTimestamp, sql.OrderDesc()))

		query.Offset(commitsThusFar).Limit(readCount)

		// Execute the query to get all commits
		commits, err := query.All(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Msgf("Error retrieving git commits w/ err: %v", err)
			message := fmt.Sprintf("Error retrieving git commits: %v", err)
			return drip.GetGitcommitsummary500JSONResponse{
				Message: &message,
			}, nil
		}
		log.Ctx(ctx).Info().Msgf("Retrieved %d commits to summarize", len(commits))
		if len(commits) == 0 {
			reachedEnd = true
			break
		}

		// Create summaries
		for _, commit := range commits {
			summary, exists := summaries[commit.Edges.Gitcommit.CommitHash]
			if !exists {
				summary = &drip.GitCommitSummary{
					CommitHash:    &commit.Edges.Gitcommit.CommitHash,
					Timestamp:     &commit.Edges.Gitcommit.CommitTimestamp,
					Author:        &commit.Edges.Gitcommit.Author,
					CommitName:    &commit.Edges.Gitcommit.CommitMessage,
					BranchName:    &commit.Edges.Gitcommit.BranchName,
					StatusSummary: &map[string]string{},
				}
				summaries[commit.Edges.Gitcommit.CommitHash] = summary
			}
			_, exists = (*summary.StatusSummary)[commit.OperatingSystem]
			if !exists {
				(*summary.StatusSummary)[commit.OperatingSystem] = string(commit.Status)
			} else if commit.Status == schema.WorkflowRunStatusTypeFailed {
				(*summary.StatusSummary)[commit.OperatingSystem] = string(commit.Status)
			}
		}
		// TODO: This heuristic hack is technically valid, but scaling up the database transfer hurts
		// And this is all just a placeholder hack anyway so meh
		/*
			if commitsThusFar == 0 {
				// Approximation of how many need to avoid too many loops
				heuristic := (targetCount * len(commits)) / len(summaries)
				readCount = max(1, min(10, heuristic+1)) * 4096
			}*/
		commitsThusFar += len(commits)
		log.Ctx(ctx).Info().Msgf("Retrieved %d commits to summarize", len(commits))
	}

	// Convert map to slice for pagination
	summarySlice := make([]drip.GitCommitSummary, 0, len(summaries))
	for _, summary := range summaries {
		summarySlice = append(summarySlice, *summary)
	}

	// Sort summaries by commit date (newest first)
	sort.Slice(summarySlice, func(i, j int) bool {
		dateI := summarySlice[i].Timestamp
		dateJ := summarySlice[j].Timestamp
		return dateI.After(*dateJ)
	})

	// Pagination
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > len(summarySlice) {
		end = len(summarySlice)
	}

	paginatedSummaries := summarySlice[start:end]
	totalPages := (len(summarySlice) + pageSize - 1) / pageSize
	displayedTotalPages := totalPages
	if !reachedEnd {
		displayedTotalPages = page + 2
	}

	log.Ctx(ctx).Info().Msgf("Git commit summaries retrieved successfully, %d summaries, %d of %d pages", len(paginatedSummaries), page, totalPages)
	return drip.GetGitcommitsummary200JSONResponse{
		CommitSummaries:    &paginatedSummaries,
		TotalNumberOfPages: &displayedTotalPages,
	}, nil
}

func (impl *DripStrictServerImplementation) GetWorkflowResult(ctx context.Context, request drip.GetWorkflowResultRequestObject) (drip.GetWorkflowResultResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.GetWorkflowResult")()

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
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.GetBranch")()

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
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.PostUploadArtifact")()

	err := impl.ComfyCIService.ProcessCIRequest(ctx, impl.Client, &request)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Error processing CI request w/ err: %v", err)
		return drip.PostUploadArtifact500Response{}, err
	}

	log.Ctx(ctx).Info().Msgf("CI request with job id '%s' processed successfully", request.Body.JobId)
	return drip.PostUploadArtifact200JSONResponse{}, nil
}

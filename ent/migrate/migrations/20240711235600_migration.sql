-- Modify "ci_workflow_results" table
ALTER TABLE "ci_workflow_results" DROP COLUMN "gpu_type", ADD COLUMN "job_id" text NULL, ADD COLUMN "cuda_version" text NULL;

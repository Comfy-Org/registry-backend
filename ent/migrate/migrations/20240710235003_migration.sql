UPDATE "ci_workflow_results"
SET "status" = 'STARTED'
WHERE "status" IS NULL;
-- Modify "ci_workflow_results" table
ALTER TABLE "ci_workflow_results" ALTER COLUMN "status" TYPE character varying, ALTER COLUMN "status" SET NOT NULL, ALTER COLUMN "status" SET DEFAULT 'STARTED', ADD COLUMN "python_version" text NULL, ADD COLUMN "vram" bigint NULL, ADD COLUMN "job_trigger_user" text NULL;
-- Modify "git_commits" table
ALTER TABLE "git_commits" ADD COLUMN "pr_number" text NULL;

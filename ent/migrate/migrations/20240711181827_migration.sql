-- Modify "ci_workflow_results" table
ALTER TABLE "ci_workflow_results" ADD COLUMN "metadata" jsonb NULL;

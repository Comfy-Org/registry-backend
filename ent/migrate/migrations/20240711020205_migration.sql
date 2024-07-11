-- Modify "ci_workflow_results" table
ALTER TABLE "ci_workflow_results" DROP COLUMN "ci_workflow_result_storage_file";
-- Modify "storage_files" table
ALTER TABLE "storage_files" ADD COLUMN "ci_workflow_result_storage_file" uuid NULL, ADD CONSTRAINT "storage_files_ci_workflow_results_storage_file" FOREIGN KEY ("ci_workflow_result_storage_file") REFERENCES "ci_workflow_results" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;

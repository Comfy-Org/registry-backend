-- Modify "ci_workflow_results" table
ALTER TABLE "ci_workflow_results" ADD COLUMN "peak_vram" bigint NULL;
-- Rename a column from "vram" to "avg_vram"
ALTER TABLE "ci_workflow_results" RENAME COLUMN "vram" TO "avg_vram";

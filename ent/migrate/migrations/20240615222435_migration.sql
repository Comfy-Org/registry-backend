-- Modify "node_versions" table
ALTER TABLE "node_versions" ADD COLUMN "status_reason" text NOT NULL DEFAULT '';

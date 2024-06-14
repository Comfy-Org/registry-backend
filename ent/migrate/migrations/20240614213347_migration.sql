-- Modify "node_versions" table
ALTER TABLE "node_versions" ADD COLUMN "status" character varying NOT NULL DEFAULT 'pending';

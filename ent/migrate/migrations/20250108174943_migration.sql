-- Modify "node_versions" table
ALTER TABLE "node_versions" ADD COLUMN "comfy_node_extract_status" character varying NOT NULL DEFAULT 'pending';

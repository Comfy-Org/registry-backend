-- Drop index "comfynode_id_node_version_id" from table: "comfy_nodes"
DROP INDEX "comfynode_id_node_version_id";
-- Modify "comfy_nodes" table
ALTER TABLE "comfy_nodes" ALTER COLUMN "id" TYPE uuid, ADD COLUMN "name" text NOT NULL;

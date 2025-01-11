-- Drop index "comfynode_id_node_version_id" from table: "comfy_nodes"
DROP INDEX "comfynode_id_node_version_id";
-- Modify "comfy_nodes" table
ALTER TABLE "comfy_nodes" ALTER COLUMN "deprecated" DROP NOT NULL, ALTER COLUMN "experimental" DROP NOT NULL, ALTER COLUMN "output_is_list" DROP NOT NULL, ALTER COLUMN "return_names" TYPE text, ALTER COLUMN "return_names" DROP NOT NULL;

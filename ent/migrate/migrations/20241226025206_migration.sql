-- Create "comfy_nodes" table
CREATE TABLE "comfy_nodes" ("id" text NOT NULL, "create_time" timestamptz NOT NULL, "update_time" timestamptz NOT NULL, "category" text NULL, "description" text NULL, "input_types" text NULL, "deprecated" boolean NOT NULL DEFAULT false, "experimental" boolean NOT NULL DEFAULT false, "output_is_list" jsonb NOT NULL, "return_names" jsonb NOT NULL, "return_types" jsonb NOT NULL, "function" text NOT NULL, "node_version_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "comfy_nodes_node_versions_comfy_nodes" FOREIGN KEY ("node_version_id") REFERENCES "node_versions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create index "comfynode_id_node_version_id" to table: "comfy_nodes"
CREATE UNIQUE INDEX "comfynode_id_node_version_id" ON "comfy_nodes" ("id", "node_version_id");
-- Modify "nodes" table
ALTER TABLE "nodes" ADD COLUMN "normalized_id" text NULL;
-- Create index "nodes_normalized_id_key" to table: "nodes"
CREATE UNIQUE INDEX "nodes_normalized_id_key" ON "nodes" ("normalized_id");

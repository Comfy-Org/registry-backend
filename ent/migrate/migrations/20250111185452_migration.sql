-- Drop the existing index on the "comfy_nodes" table if necessary
DROP INDEX IF EXISTS "comfynode_id_node_version_id";

-- Step 1: Add a temporary column with the `uuid` type if needed
ALTER TABLE "comfy_nodes"
    ADD COLUMN "id_temp" uuid;

-- Step 2: Update the new `id_temp` column with the existing `id` column data, casting it to `uuid`
UPDATE "comfy_nodes"
SET "id_temp" = "id"::uuid;

-- Step 3: Drop the old `id` column
ALTER TABLE "comfy_nodes" DROP COLUMN "id";

-- Step 4: Rename the temporary column to `id`
ALTER TABLE "comfy_nodes" RENAME COLUMN "id_temp" TO "id";

-- Step 5: Add the new `name` column
ALTER TABLE "comfy_nodes"
    ADD COLUMN "name" text NOT NULL;

-- Step 6: Add the primary key constraint to the `id` column
ALTER TABLE "comfy_nodes"
    ADD PRIMARY KEY ("id");

-- Step 7: Optionally, recreate the index if needed
-- For example, if you want to reindex with the new `id` column
CREATE INDEX "comfynode_id_node_version_id" ON "comfy_nodes" ("id");

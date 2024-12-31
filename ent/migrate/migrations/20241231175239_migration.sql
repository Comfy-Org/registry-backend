-- Modify "comfy_nodes" table
ALTER TABLE "comfy_nodes" ALTER COLUMN "return_types" TYPE text, ALTER COLUMN "return_types" DROP NOT NULL;

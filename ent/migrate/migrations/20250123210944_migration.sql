-- Modify "nodes" table
ALTER TABLE "nodes" ADD COLUMN "last_algolia_index_time" timestamptz NULL;

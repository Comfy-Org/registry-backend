-- Modify "nodes" table
ALTER TABLE "nodes" ADD COLUMN "status" character varying NOT NULL DEFAULT 'pending';

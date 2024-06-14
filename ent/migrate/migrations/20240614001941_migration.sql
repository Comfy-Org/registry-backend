-- Modify "publishers" table
ALTER TABLE "publishers" ADD COLUMN "status" character varying NOT NULL DEFAULT 'active';

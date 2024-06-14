-- Modify "publishers" table
ALTER TABLE "publishers" ADD COLUMN "status" character varying NOT NULL DEFAULT 'ACTIVE';
-- Modify "users" table
ALTER TABLE "users" ADD COLUMN "status" character varying NOT NULL DEFAULT 'ACTIVE';

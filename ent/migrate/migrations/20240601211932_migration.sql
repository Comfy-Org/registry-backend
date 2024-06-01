-- Modify "nodes" table
ALTER TABLE "nodes" ADD COLUMN "total_install" bigint NOT NULL DEFAULT 0, ADD COLUMN "total_star" bigint NOT NULL DEFAULT 0, ADD COLUMN "total_review" bigint NOT NULL DEFAULT 0;
-- Create "node_reviews" table
CREATE TABLE "node_reviews" ("id" uuid NOT NULL, "star" bigint NOT NULL DEFAULT 0, "node_id" text NOT NULL, "user_id" character varying NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "node_reviews_nodes_reviews" FOREIGN KEY ("node_id") REFERENCES "nodes" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "node_reviews_users_reviews" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create index "nodereview_node_id_user_id" to table: "node_reviews"
CREATE UNIQUE INDEX "nodereview_node_id_user_id" ON "node_reviews" ("node_id", "user_id");

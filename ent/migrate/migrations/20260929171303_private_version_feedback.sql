-- Modify "node_versions" table
-- Older installations do not have tags_admin; newer installations already do.
-- Preserve existing policy tags and initialize only previously missing values.
ALTER TABLE "node_versions" ADD COLUMN IF NOT EXISTS "tags_admin" text;
UPDATE "node_versions" SET "tags_admin" = '[]' WHERE "tags_admin" IS NULL;
ALTER TABLE "node_versions" ALTER COLUMN "tags_admin" SET NOT NULL;
-- Keep inserts from the previous backend release working during rollout.
ALTER TABLE "node_versions" ALTER COLUMN "tags_admin" SET DEFAULT '[]';
-- Create "feedback_threads" table
CREATE TABLE "feedback_threads" ("id" uuid NOT NULL, "create_time" timestamptz NOT NULL, "update_time" timestamptz NOT NULL, "version_id" uuid NOT NULL, "publisher_id" character varying NOT NULL, "state" character varying NOT NULL DEFAULT 'awaiting_author', "created_by_user_id" character varying NOT NULL, "last_message_seq" bigint NOT NULL DEFAULT 0, "last_message_at" timestamptz NOT NULL, "revision" bigint NOT NULL DEFAULT 0, "resolved_by_user_id" character varying NULL, "resolved_at" timestamptz NULL, "archived_at" timestamptz NULL, PRIMARY KEY ("id"));
-- Create index "feedbackthread_version_id_publisher_id" to table: "feedback_threads"
CREATE UNIQUE INDEX "feedbackthread_version_id_publisher_id" ON "feedback_threads" ("version_id", "publisher_id");
-- Create index "feedbackthread_publisher_id_state_last_message_at" to table: "feedback_threads"
CREATE INDEX "feedbackthread_publisher_id_state_last_message_at" ON "feedback_threads" ("publisher_id", "state", "last_message_at");
-- Create index "feedbackthread_state_last_message_at" to table: "feedback_threads"
CREATE INDEX "feedbackthread_state_last_message_at" ON "feedback_threads" ("state", "last_message_at");
-- Create "feedback_events" table
CREATE TABLE "feedback_events" ("id" uuid NOT NULL, "actor_user_id" character varying NOT NULL, "event_type" character varying NOT NULL, "created_at" timestamptz NOT NULL, "thread_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "feedback_events_feedback_threads_events" FOREIGN KEY ("thread_id") REFERENCES "feedback_threads" ("id") ON DELETE NO ACTION);
-- Create index "feedbackevent_thread_id_created_at" to table: "feedback_events"
CREATE INDEX "feedbackevent_thread_id_created_at" ON "feedback_events" ("thread_id", "created_at");
-- Create "feedback_messages" table
CREATE TABLE "feedback_messages" ("id" uuid NOT NULL, "seq" bigint NOT NULL, "sender_user_id" character varying NOT NULL, "sender_name" character varying NOT NULL, "sender_role" character varying NOT NULL, "body" character varying NOT NULL, "client_message_id" uuid NOT NULL, "created_at" timestamptz NOT NULL, "thread_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "feedback_messages_feedback_threads_messages" FOREIGN KEY ("thread_id") REFERENCES "feedback_threads" ("id") ON DELETE NO ACTION);
-- Create index "feedbackmessage_thread_id_seq" to table: "feedback_messages"
CREATE UNIQUE INDEX "feedbackmessage_thread_id_seq" ON "feedback_messages" ("thread_id", "seq");
-- Create index "feedbackmessage_thread_id_sender_user_id_client_message_id" to table: "feedback_messages"
CREATE UNIQUE INDEX "feedbackmessage_thread_id_sender_user_id_client_message_id" ON "feedback_messages" ("thread_id", "sender_user_id", "client_message_id");
-- Create "feedback_reads" table
CREATE TABLE "feedback_reads" ("id" uuid NOT NULL, "user_id" character varying NOT NULL, "last_read_message_seq" bigint NOT NULL DEFAULT 0, "read_at" timestamptz NOT NULL, "thread_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "feedback_reads_feedback_threads_reads" FOREIGN KEY ("thread_id") REFERENCES "feedback_threads" ("id") ON DELETE NO ACTION);
-- Create index "feedbackread_thread_id_user_id" to table: "feedback_reads"
CREATE UNIQUE INDEX "feedbackread_thread_id_user_id" ON "feedback_reads" ("thread_id", "user_id");
-- Create index "feedbackread_user_id_thread_id" to table: "feedback_reads"
CREATE INDEX "feedbackread_user_id_thread_id" ON "feedback_reads" ("user_id", "thread_id");

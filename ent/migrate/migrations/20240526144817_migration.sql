-- Create "git_commits" table
CREATE TABLE "git_commits" ("id" uuid NOT NULL, "create_time" timestamptz NOT NULL, "update_time" timestamptz NOT NULL, "commit_hash" text NOT NULL, "branch_name" text NOT NULL, "repo_name" text NOT NULL, "commit_message" text NOT NULL, "commit_timestamp" timestamptz NOT NULL, "author" text NULL, "timestamp" timestamptz NULL, PRIMARY KEY ("id"));
-- Create index "gitcommit_repo_name_commit_hash" to table: "git_commits"
CREATE UNIQUE INDEX "gitcommit_repo_name_commit_hash" ON "git_commits" ("repo_name", "commit_hash");
-- Create "storage_files" table
CREATE TABLE "storage_files" ("id" uuid NOT NULL, "create_time" timestamptz NOT NULL, "update_time" timestamptz NOT NULL, "bucket_name" text NOT NULL, "object_name" text NULL, "file_path" text NOT NULL, "file_type" text NOT NULL, "file_url" text NULL, PRIMARY KEY ("id"));
-- Create "ci_workflow_results" table
CREATE TABLE "ci_workflow_results" ("id" uuid NOT NULL, "create_time" timestamptz NOT NULL, "update_time" timestamptz NOT NULL, "operating_system" text NOT NULL, "gpu_type" text NULL, "pytorch_version" text NULL, "workflow_name" text NULL, "run_id" text NULL, "status" text NULL, "start_time" bigint NULL, "end_time" bigint NULL, "ci_workflow_result_storage_file" uuid NULL, "git_commit_results" uuid NULL, PRIMARY KEY ("id"), CONSTRAINT "ci_workflow_results_git_commits_results" FOREIGN KEY ("git_commit_results") REFERENCES "git_commits" ("id") ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT "ci_workflow_results_storage_files_storage_file" FOREIGN KEY ("ci_workflow_result_storage_file") REFERENCES "storage_files" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create "publishers" table
CREATE TABLE "publishers" ("id" text NOT NULL, "create_time" timestamptz NOT NULL, "update_time" timestamptz NOT NULL, "name" text NOT NULL, "description" text NULL, "website" text NULL, "support_email" text NULL, "source_code_repo" text NULL, "logo_url" text NULL, PRIMARY KEY ("id"));
-- Create "nodes" table
CREATE TABLE "nodes" ("id" text NOT NULL, "create_time" timestamptz NOT NULL, "update_time" timestamptz NOT NULL, "name" text NOT NULL, "description" text NULL, "author" text NULL, "license" text NOT NULL, "repository_url" text NOT NULL, "icon_url" text NULL, "tags" text NOT NULL, "publisher_id" text NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "nodes_publishers_nodes" FOREIGN KEY ("publisher_id") REFERENCES "publishers" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create "node_versions" table
CREATE TABLE "node_versions" ("id" uuid NOT NULL, "create_time" timestamptz NOT NULL, "update_time" timestamptz NOT NULL, "version" text NOT NULL, "changelog" text NULL, "pip_dependencies" text NOT NULL, "deprecated" boolean NOT NULL DEFAULT false, "node_id" text NOT NULL, "node_version_storage_file" uuid NULL, PRIMARY KEY ("id"), CONSTRAINT "node_versions_nodes_versions" FOREIGN KEY ("node_id") REFERENCES "nodes" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "node_versions_storage_files_storage_file" FOREIGN KEY ("node_version_storage_file") REFERENCES "storage_files" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create index "nodeversion_node_id_version" to table: "node_versions"
CREATE UNIQUE INDEX "nodeversion_node_id_version" ON "node_versions" ("node_id", "version");
-- Create "personal_access_tokens" table
CREATE TABLE "personal_access_tokens" ("id" uuid NOT NULL, "create_time" timestamptz NOT NULL, "update_time" timestamptz NOT NULL, "name" text NOT NULL, "description" text NOT NULL, "token" text NOT NULL, "publisher_id" text NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "personal_access_tokens_publishers_personal_access_tokens" FOREIGN KEY ("publisher_id") REFERENCES "publishers" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create index "personal_access_tokens_token_key" to table: "personal_access_tokens"
CREATE UNIQUE INDEX "personal_access_tokens_token_key" ON "personal_access_tokens" ("token");
-- Create index "personalaccesstoken_token" to table: "personal_access_tokens"
CREATE UNIQUE INDEX "personalaccesstoken_token" ON "personal_access_tokens" ("token");
-- Create "users" table
CREATE TABLE "users" ("id" character varying NOT NULL, "create_time" timestamptz NOT NULL, "update_time" timestamptz NOT NULL, "email" character varying NULL, "name" character varying NULL, "is_approved" boolean NOT NULL DEFAULT false, "is_admin" boolean NOT NULL DEFAULT false, PRIMARY KEY ("id"));
-- Create "publisher_permissions" table
CREATE TABLE "publisher_permissions" ("id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY, "permission" character varying NOT NULL, "publisher_id" text NOT NULL, "user_id" character varying NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "publisher_permissions_publishers_publisher_permissions" FOREIGN KEY ("publisher_id") REFERENCES "publishers" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "publisher_permissions_users_publisher_permissions" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);

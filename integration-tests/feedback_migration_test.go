package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"registry-backend/config"
	"registry-backend/db"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/schema"
)

func TestFeedbackMigrationSupportsExistingVersionsAndConversations(t *testing.T) {
	ctx := context.Background()
	migrations, err := filepath.Glob("../ent/migrate/migrations/*.sql")
	require.NoError(t, err)
	require.NotEmpty(t, migrations)
	const feedbackMigration = "20260929171303_private_version_feedback.sql"
	for _, existingTags := range []bool{false, true} {
		t.Run(map[bool]string{false: "legacy versions", true: "existing public policy tags"}[existingTags], func(t *testing.T) {
			dsn, cleanup := setupDatabaseURL(t, ctx)
			defer cleanup()
			database, err := sql.Open("postgres", dsn)
			require.NoError(t, err)
			defer database.Close()
			// Build the previous release from the committed SQL, without Schema.Create.
			for _, file := range migrations {
				if filepath.Base(file) >= feedbackMigration {
					break
				}
				source, err := os.ReadFile(file)
				require.NoError(t, err)
				_, err = database.ExecContext(ctx, string(source))
				require.NoError(t, err, file)
			}
			client, err := ent.Open("postgres", dsn)
			require.NoError(t, err)
			defer client.Close()
			client.User.Create().SetID("admin").SetIsAdmin(true).SaveX(ctx)
			client.User.Create().SetID("owner").SaveX(ctx)
			client.Publisher.Create().SetID("publisher").SetName("Publisher").SaveX(ctx)
			client.PublisherPermission.Create().SetPublisherID("publisher").SetUserID("owner").SetPermission(schema.PublisherPermissionTypeOwner).SaveX(ctx)
			client.Node.Create().SetID("example").SetNormalizedID("example").SetName("Example").SetPublisherID("publisher").SetLicense("MIT").SetRepositoryURL("https://example.invalid").SaveX(ctx)
			insertLegacyVersion := func(id uuid.UUID, number string) {
				t.Helper()
				_, err := database.ExecContext(ctx, `INSERT INTO node_versions(id, create_time, update_time, version, pip_dependencies, deprecated, node_id, status, status_reason) VALUES ($1, NOW(), NOW(), $2, '[]', false, 'example', 'flagged', 'private-scan-evidence')`, id, number)
				require.NoError(t, err)
			}
			versionID := uuid.New()
			insertLegacyVersion(versionID, "1.0.0")
			expectedTags := []string{}
			if existingTags {
				expectedTags = []string{"any-code-execute"}
				_, err = database.ExecContext(ctx, `ALTER TABLE node_versions ADD COLUMN tags_admin text; UPDATE node_versions SET tags_admin='["any-code-execute"]'`)
				require.NoError(t, err)
			}
			tx, err := database.BeginTx(ctx, nil)
			require.NoError(t, err)
			defer tx.Rollback()
			for _, file := range migrations {
				if filepath.Base(file) < feedbackMigration {
					continue
				}
				source, err := os.ReadFile(file)
				require.NoError(t, err)
				_, err = tx.ExecContext(ctx, string(source))
				require.NoError(t, err, file)
			}
			require.NoError(t, tx.Commit())
			db.RegisterFeedbackHooks(client)
			stored := client.NodeVersion.GetX(ctx, versionID)
			require.Equal(t, expectedTags, stored.TagsAdmin)
			require.Equal(t, "private-scan-evidence", stored.StatusReason)
			legacyID := uuid.New()
			insertLegacyVersion(legacyID, "1.0.1")
			require.Equal(t, []string{}, client.NodeVersion.GetX(ctx, legacyID).TagsAdmin, "previous-release inserts must retain the database default")

			impl := NewStrictServerImplementationWithMocks(client, &config.Config{})
			server := newRegistryHTTPTestServer(impl.DripStrictServerImplementation)
			call := func(actor, method, path string, input, output any) {
				t.Helper()
				body, err := json.Marshal(input)
				require.NoError(t, err)
				request := httptest.NewRequest(method, path, bytes.NewReader(body))
				request.Header.Set("Content-Type", "application/json")
				request.Header.Set("X-Test-User", actor)
				response := httptest.NewRecorder()
				server.ServeHTTP(response, request)
				require.Equal(t, 200, response.Code, response.Body.String())
				require.NoError(t, json.Unmarshal(response.Body.Bytes(), output))
			}
			adminPath := fmt.Sprintf("/admin/nodes/example/versions/%s/feedback", versionID)
			authorPath := fmt.Sprintf("/publishers/publisher/nodes/example/versions/%s/feedback", versionID)
			message := drip.FeedbackMessageInput{Target: drip.FeedbackWriteTarget{PublisherId: "publisher"}, Body: "Please revise installation.", ClientMessageId: uuid.New()}
			var conversation drip.FeedbackResponse
			call("admin", "POST", adminPath+"/messages", message, &conversation)
			require.Len(t, conversation.Messages, 1)
			message.Target, message.Body, message.ClientMessageId = conversation.Target, "Installation updated.", uuid.New()
			call("owner", "POST", authorPath+"/messages", message, &conversation)
			call("owner", "POST", authorPath+"/messages", message, &conversation)
			require.Len(t, conversation.Messages, 2, "a retry must not add another message")
			var read drip.FeedbackReadResult
			call("owner", "POST", authorPath+"/read", drip.FeedbackReadInput{Target: conversation.Target, LastReadMessageSeq: 2}, &read)
			require.Equal(t, 2, read.LastReadMessageSeq)
			call("admin", "PATCH", adminPath, drip.FeedbackStateInput{Target: conversation.Target, State: "resolved", ExpectedRevision: conversation.Thread.Revision}, &conversation)
			call("owner", "GET", authorPath, nil, &conversation)
			require.Equal(t, drip.FeedbackState("resolved"), conversation.Thread.State)
			require.Equal(t, 2, conversation.LastReadMessageSeq)
			require.Zero(t, conversation.UnreadCount)
			require.Len(t, conversation.Events, 1)
			require.Equal(t, drip.FeedbackEventEventType("resolved"), conversation.Events[0].EventType)
			require.Equal(t, "Installation updated.", conversation.Messages[1].Body)
			require.Equal(t, schema.NodeVersionStatusFlagged, client.NodeVersion.GetX(ctx, versionID).Status)
			require.Equal(t, 1, client.FeedbackThread.Query().CountX(ctx))
			require.Equal(t, 2, client.FeedbackMessage.Query().CountX(ctx))
			require.Equal(t, 1, client.FeedbackRead.Query().CountX(ctx))
			require.Equal(t, 1, client.FeedbackEvent.Query().CountX(ctx))
			requestBody, err := json.Marshal(drip.FeedbackSupersedeInput{Versions: []drip.FeedbackSupersedeVersion{{VersionId: versionID, Target: conversation.Target, ExpectedRevision: conversation.Thread.Revision}}})
			require.NoError(t, err)
			request := httptest.NewRequest("POST", fmt.Sprintf("/publishers/publisher/nodes/example/versions/%s/feedback/supersede", legacyID), bytes.NewReader(requestBody))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("X-Test-User", "owner")
			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)
			require.Equal(t, 204, response.Code, response.Body.String())
			call("admin", "GET", adminPath, nil, &conversation)
			require.Len(t, conversation.Events, 2)
			require.Equal(t, legacyID, *conversation.Events[1].ReplacementVersionId)
			require.Equal(t, "1.0.1", *conversation.Events[1].ReplacementVersion)
			require.True(t, client.NodeVersion.GetX(ctx, versionID).Deprecated)

		})
	}
}

package mapper_test

import (
	"encoding/json"
	"registry-backend/ent"
	"registry-backend/ent/schema"
	"registry-backend/mapper"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPublicVersionExcludesScanAndPreservesInstallationTags(t *testing.T) {
	// JSON fixtures also exercise compatibility with persisted version data.
	var version ent.NodeVersion
	require.NoError(t, json.Unmarshal([]byte(`{"version":"1.2.0","node_id":"example","status":"flagged","status_reason":"private-scanner-evidence","tags_admin":["any-code-execute"]}`), &version))
	encoded, err := json.Marshal(mapper.DbNodeVersionToApiNodeVersion(&version))
	require.NoError(t, err)
	var result map[string]any
	require.NoError(t, json.Unmarshal(encoded, &result))
	t.Run("raw scan never appears in public responses", func(t *testing.T) {
		require.NotContains(t, result, "status_reason")
		require.NotContains(t, string(encoded), "private-scanner-evidence")
	})
	t.Run("installation policy tags remain public", func(t *testing.T) {
		require.Equal(t, []any{"any-code-execute"}, result["tags_admin"])
	})
}

func TestSearchDocumentExcludesInternalFields(t *testing.T) {
	version := &ent.NodeVersion{
		NodeID: "example", Version: "1.0.0", StatusReason: "private-evidence", TagsAdmin: []string{"any-code-execute"},
		Status:                  schema.NodeVersionStatusFlagged,
		ComfyNodeCloudBuildInfo: schema.ComfyNodeCloudBuildInfo{ProjectID: "internal-build-project"},
		Edges:                   ent.NodeVersionEdges{StorageFile: &ent.StorageFile{FilePath: "internal-storage-path"}},
	}
	original, err := json.Marshal(version)
	require.NoError(t, err)
	document := mapper.PublicNodeVersionSearchDocument(version)
	encoded, err := json.Marshal(document)
	require.NoError(t, err)
	for _, field := range []string{"status_reason", "comfy_node_cloud_build_info", "comfy_node_extract_status", "edges"} {
		require.NotContains(t, document, field)
	}
	for _, private := range []string{"private-evidence", "internal-build-project", "internal-storage-path"} {
		require.NotContains(t, string(encoded), private)
	}
	require.Equal(t, []string{"any-code-execute"}, document["tags_admin"])
	after, err := json.Marshal(version)
	require.NoError(t, err)
	require.JSONEq(t, string(original), string(after), "indexing must preserve the source entity")
}

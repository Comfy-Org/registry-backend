# Private feedback API for admin agents

Agents can review versions, read conversations, send Markdown replies and resolve/reopen conversations over HTTP. The [OpenAPI document](../openapi.yml) is the contract; the generated server and web client use it. No browser UI is required for these operations.

Authenticate with `Authorization: Bearer <Firebase ID token>` for an active Registry admin. Every request checks the current database role. Publisher CLI tokens and legacy moderation JWTs are not feedback credentials. Keep token acquisition and refresh in the agent's existing identity integration.

## Select work

| Task | Request |
| --- | --- |
| Flagged versions with no feedback | `GET /admin/nodeversions?statuses=NodeVersionStatusFlagged&feedback_status=unprocessed` |
| Versions whose latest feedback is from an author | `GET /admin/nodeversions?feedback_status=needs_response` |
| Versions whose latest feedback is from an admin | `GET /admin/nodeversions?feedback_status=processed` |
| Open conversations waiting for an admin | `GET /admin/node-version-feedback?state=awaiting_admin` |
| Open conversations waiting for an author | `GET /admin/node-version-feedback?state=awaiting_author` |
| Resolved conversations | `GET /admin/node-version-feedback?state=resolved` |

The handling filter describes the last sender and includes resolved/archived conversations. Use the conversation `state` filter for a reply queue. Versions with no conversation appear in the version endpoint, not the conversation inbox. Version moderation remains independent of feedback state.

Version lists accept `nodeId`, exact `version`, `statuses`, `status_reason`, `page` and `pageSize` (1–100; default 10). Filters apply before totals and pagination. These responses include private scan evidence. Do not copy it wholesale into an author-visible message. Conversation lists accept `nodeId`, `publisherId`, `state` and an opaque `cursor`; they return at most 20 summaries and `next_cursor`. Omit a node/Publisher filter to leave that dimension unrestricted; an explicitly empty `nodeId=` or `publisherId=` matches no conversations. An explicitly empty or unknown `state` returns 422; omit `state` to include every state. Keep filters unchanged while following cursors. Archived threads can be listed but cannot be reopened or replied to.

Mutations can move records between work queues. Re-query page 1 after processing a numbered queue, or collect its IDs before writing. Do not advance numbered pages while removing items from the same filtered queue. Cursor queues are ordered by latest activity and are not a frozen snapshot; deduplicate by thread ID when processing a changing queue.

## Read, decide and write

1. `GET /admin/nodes/{nodeId}/versions/{versionId}/feedback`. `versionId` is the version UUID, not `1.0.0`. The response contains `target`, `thread`, `permissions`, messages and read information. An absent thread is `null`; `target` still identifies the intended Publisher for a first message.
2. Consume the relevant history. GET returns the latest 30 messages in ascending sequence order. Pass `next_before_seq` as `before_seq` until it is null to retrieve older pages. Deduplicate by message ID/sequence when polling; GET does not acknowledge messages.
3. Check `permissions.can_start` or `permissions.can_reply`. Compose a response for this specific `target`. Store a fresh UUID as `client_message_id` with the pending payload before sending.
4. `POST .../feedback/messages` with the unchanged target, Markdown source and that UUID. On an ambiguous failure, retry the **same payload and UUID**. A duplicate successful submission returns the conversation without inserting another message. The response supplies the updated state, revision and permissions.
5. To close or reopen, check `can_resolve` / `can_reopen`, then `PATCH .../feedback` with `state: "resolved"` or `"open"`, the observed `target`, and `expected_revision: thread.revision`. Reopen derives the awaiting side from the last sender. State changes neither approve the version nor alter the handling classification.
6. Optionally `POST .../feedback/read` with the current target and the highest contiguous `last_read_message_seq` actually consumed by this admin identity. Reading, replying and resolving are separate operations. Repeating a read acknowledgement cannot move it backward.

Message bodies are UTF-8 text with up to 5,000 Unicode code points after trimming. The web renders CommonMark safely; agents should send source text, not pre-rendered HTML. No attachments or message editing are supported. Messages retain the authenticated sender's user ID for attribution.

### Example: reply using the observed target

Set `REGISTRY_API_BASE`, `REGISTRY_ADMIN_TOKEN`, `REGISTRY_NODE_ID` and `REGISTRY_VERSION_ID` through the agent's existing configuration. Supply URL-encoded path segments for node/version IDs. This example requires curl, jq and uuidgen. Keep the generated request file and its UUID for retries; re-running payload creation generates a different logical message.

```sh
set -eu
umask 077
feedback_url="$REGISTRY_API_BASE/admin/nodes/$REGISTRY_NODE_ID/versions/$REGISTRY_VERSION_ID/feedback"
curl --fail-with-body --silent --show-error \
  -H "Authorization: Bearer $REGISTRY_ADMIN_TOKEN" \
  "$feedback_url" > conversation.json

jq -e '.permissions.can_start or .permissions.can_reply' conversation.json
client_message_id="$(uuidgen)"
jq --arg id "$client_message_id" \
  --arg body 'Please remove the **remote install script** and share the replacement nodepack version.' \
  '{target, client_message_id: $id, body: $body}' \
  conversation.json > reply.json

curl --fail-with-body --silent --show-error \
  -H "Authorization: Bearer $REGISTRY_ADMIN_TOKEN" \
  -H 'Content-Type: application/json' \
  --data-binary @reply.json \
  "$feedback_url/messages" > reply-response.json
```

For a close operation, fetch the current conversation, verify `can_resolve`, and build the body with `jq '{target, state: "resolved", expected_revision: .thread.revision}' conversation.json`. Send that JSON with `PATCH` to `feedback_url`. For reopen, verify `can_reopen` and use `state: "open"`.

## Handoff after a replacement upload

An owner can explicitly deprecate earlier feedback versions and close their conversations using `POST /publishers/{publisherId}/nodes/{nodeId}/versions/{replacementVersionId}/feedback/supersede`. This is an owner operation: an admin role alone does not authorize it. The URL UUID identifies the newer uploaded nodepack version; each item identifies an older version from the current node/Publisher summary inbox:

```json
{
  "versions": [
    {
      "version_id": "44444444-4444-4444-8444-444444444444",
      "target": {
        "publisher_id": "publisher",
        "thread_id": "55555555-5555-4555-8555-555555555555"
      },
      "expected_revision": 2
    }
  ]
}
```

The request supports 1–100 distinct versions within the 32 KB body limit. They must have current non-archived feedback and precede the replacement by upload time; neither selected versions nor the replacement can be deleted, and the replacement cannot be deprecated. `204` means every selection was deprecated and closed in one database transaction. A stale binding/revision or invalid replacement returns 409 and rolls back the batch. Inaccessible/foreign/no-feedback versions return 404. Duplicate IDs or an empty selection return 422.

Keep the exact request for ambiguous retries. A matching completed handoff is a no-op; a reopened or otherwise changed conversation conflicts. Search indexing uses the existing asynchronous gateway after commit and is retried with the same request.

For reviewers, `events[].event_type == "superseded"` means **the author will no longer follow this version's feedback**, not that the issue was verified or the version approved. The event includes `replacement_version_id` and the nodepack number `replacement_version`. Inspect it with `GET /admin/nodeversions?nodeId=example&version=1.0.1`, then use that replacement UUID for a separate feedback thread (a first message still requires a flagged version). Feedback, tags and moderation remain attached to their original version; uploading itself does not close or transfer a discussion. Ordinary admin resolve/reopen still works independently.

## Reference an earlier issue in new feedback

A feedback message on the replacement can cite the original issue with an ordinary Markdown link. Read the old thread or find it in the inbox; construct the URL from its `id`, `node_id`, `version_id` and `publisher_id` using the configured web origin and URL-encoded values:

```text
{webOrigin}/feedback/{thread.id}?nodeId={thread.node_id}&versionId={thread.version_id}&publisherId={thread.publisher_id}
```

For example, send this source as the new version's message `body` using its independently observed target and a fresh client message UUID:

```markdown
The installation issue remains in 1.0.1. See [example@1.0.0 feedback](https://registry.example/feedback/11111111-1111-4111-8111-111111111111?nodeId=example&versionId=44444444-4444-4444-8444-444444444444&publisherId=publisher).
```

Replace the illustrative hostname with the deployed web origin. Escape Markdown punctuation in a dynamically built label; the web's **Copy Markdown reference** action already produces this format. Store and retry the entire source with the same UUID as any other message. The API stores Markdown text and does not interpret references or fetch linked URLs. No additional endpoint, database relation or token is required. A reference does not move or reopen the original thread, and closed/deprecated history remains accessible only through existing authorization. The web verifies the linked thread and Publisher so transfer cannot redirect it to a new recipient's thread. Agents read the original history through its existing feedback GET, not through the HTML page.

## Handle results without parsing English messages

| HTTP result | Agent action |
| --- | --- |
| `204` | Owner handoff completed; no response body. |
| `200` | Consume the returned target, state, revision and permissions. |
| `401` | Refresh/reacquire the configured identity credential; do not retry with a Publisher token. |
| `403` / `404` | Stop work on the inaccessible resource. Private endpoints normally conceal revoked roles and inaccessible resources with 404. |
| `409` | GET again and reassess. The recipient, revision, state or retry payload may conflict. Never copy a pending message to a changed Publisher/thread or blindly replace its revision. If a lost close/reopen response already achieved the intended state, treat that state as observed success. |
| `413` / `422` | Correct the request/body before retrying. |
| `429` | Back off with jitter; respect `Retry-After` if present. |
| Network error / `5xx` | Use bounded retries. Message retries reuse the same target, body and UUID; state mutations require a fresh read if the outcome is unknown. |

Errors provide a JSON `message`; some authentication paths also provide `error`. The OpenAPI `FeedbackError` schema reflects this. Use status codes and the structured conversation state/permissions for decisions, not substring matching against `message`. Private responses use `Cache-Control: private, no-store`; keep credentials and conversation files out of public logs, repositories and shared caches.

## Verification

HTTP integration tests cover the three handling filters before pagination, each inbox state through reply/resolve/reopen, malformed/empty states, explicit empty node/Publisher filters, non-admin rejection, current-role revocation, target/revision conflicts, concurrent and repeated message submissions, history/inbox pagination and error-body conformance to OpenAPI. The HTTP test fixture supplies the principal at the Firebase boundary; live token issuance and refresh remain deployment integration checks.

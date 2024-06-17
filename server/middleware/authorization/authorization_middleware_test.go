package drip_authorization

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/schema"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var _ Assertor = mockAlwayErrorAssertor{}
var errMockAssertor = errors.New("assertion failed")

type mockAlwayErrorAssertor struct{}

// AssertAccessTokenBelongsToPublisher implements Assertor.
func (m mockAlwayErrorAssertor) AssertAccessTokenBelongsToPublisher(ctx context.Context, client *ent.Client, publisherID string, tokenId uuid.UUID) error {
	return errors.New("assertion failed")
}

// AssertNodeBanned implements Assertor.
func (m mockAlwayErrorAssertor) AssertNodeBanned(ctx context.Context, client *ent.Client, nodeID string) error {
	return errors.New("assertion failed")
}

// AssertNodeBelongsToPublisher implements Assertor.
func (m mockAlwayErrorAssertor) AssertNodeBelongsToPublisher(ctx context.Context, client *ent.Client, publisherID string, nodeID string) error {
	return errors.New("assertion failed")
}

// AssertPublisherBanned implements Assertor.
func (m mockAlwayErrorAssertor) AssertPublisherBanned(ctx context.Context, client *ent.Client, publisherID string) error {
	return errors.New("assertion failed")
}

// AssertPublisherPermissions implements Assertor.
func (m mockAlwayErrorAssertor) AssertPublisherPermissions(ctx context.Context, client *ent.Client, publisherID string, userID string, permissions []schema.PublisherPermissionType) (err error) {
	return errors.New("assertion failed")
}

// IsPersonalAccessTokenValidForPublisher implements Assertor.
func (m mockAlwayErrorAssertor) IsPersonalAccessTokenValidForPublisher(ctx context.Context, client *ent.Client, publisherID string, accessToken string) (bool, error) {
	return false, errors.New("assertion failed")
}

func TestNoAuthz(t *testing.T) {
	mw := NewAuthorizationManager(&ent.Client{}, mockAlwayErrorAssertor{}).AuthorizationMiddleware()
	req, res := httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)

	tests := []struct {
		op   string
		pass bool
		req  interface{}
	}{
		{op: "SomeOtherOperation", pass: true},

		{op: "CreatePublisher", pass: false, req: drip.CreatePublisherRequestObject{}},
		{op: "UpdatePublisher", pass: false, req: drip.UpdatePublisherRequestObject{}},
		{op: "CreateNode", pass: false, req: drip.CreateNodeRequestObject{}},
		{op: "DeleteNode", pass: false, req: drip.DeleteNodeRequestObject{}},
		{op: "UpdateNode", pass: false, req: drip.UpdateNodeRequestObject{}},
		{op: "GetNode", pass: false, req: drip.GetNodeRequestObject{}},
		{op: "PublishNodeVersion", pass: false, req: drip.PublishNodeVersionRequestObject{}},
		{op: "UpdateNodeVersion", pass: false, req: drip.UpdateNodeVersionRequestObject{}},
		{op: "DeleteNodeVersion", pass: false, req: drip.DeleteNodeVersionRequestObject{}},
		{op: "GetNodeVersion", pass: false, req: drip.GetNodeVersionRequestObject{}},
		{op: "ListNodeVersions", pass: false, req: drip.ListNodeVersionsRequestObject{}},
		{op: "InstallNode", pass: false, req: drip.InstallNodeRequestObject{}},
		{op: "CreatePersonalAccessToken", pass: false, req: drip.CreatePersonalAccessTokenRequestObject{}},
		{op: "DeletePersonalAccessToken", pass: false, req: drip.DeletePersonalAccessTokenRequestObject{}},
		{op: "GetPermissionOnPublisherNodes", pass: false, req: drip.GetPermissionOnPublisherNodesRequestObject{}},
		{op: "GetPermissionOnPublisher", pass: false, req: drip.GetPermissionOnPublisherRequestObject{}},
	}
	for _, test := range tests {
		handled := false
		h := func(ctx echo.Context, request interface{}) (interface{}, error) {
			handled = true
			return nil, nil
		}
		mw(h, test.op)(ctx, test.req)
		assert.Equal(t, test.pass, handled)
	}
}

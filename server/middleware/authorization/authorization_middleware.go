package drip_authorization

import (
	"registry-backend/drip"
	"registry-backend/ent/schema"
	"slices"

	strictecho "github.com/oapi-codegen/runtime/strictmiddleware/echo"
)

func (m *AuthorizationManager) AuthorizationMiddleware() drip.StrictMiddlewareFunc {
	subMiddlewares := map[string][]drip.StrictMiddlewareFunc{
		"CreatePublisher": {
			m.assertUserBanned(),
		},
		"UpdatePublisher": {
			m.assertUserBanned(),
			m.assertPublisherPermission(
				[]schema.PublisherPermissionType{schema.PublisherPermissionTypeOwner},
				func(req interface{}) (publisherID string) {
					return req.(drip.UpdatePublisherRequestObject).PublisherId
				},
			),
		},
		"CreateNode": {
			m.assertUserBanned(),
			m.assertPublisherPermission(
				[]schema.PublisherPermissionType{schema.PublisherPermissionTypeOwner},
				func(req interface{}) (publisherID string) {
					return req.(drip.CreateNodeRequestObject).PublisherId
				},
			),
		},
		"DeleteNode": {
			m.assertUserBanned(),
			m.assertPublisherPermission(
				[]schema.PublisherPermissionType{schema.PublisherPermissionTypeOwner},
				func(req interface{}) (publisherID string) {
					return req.(drip.DeleteNodeRequestObject).PublisherId
				},
			),
			m.assertNodeBelongsToPublisher(
				func(req interface{}) (publisherID string) {
					return req.(drip.DeleteNodeRequestObject).PublisherId
				},
				func(req interface{}) (nodeID string) {
					return req.(drip.DeleteNodeRequestObject).NodeId
				},
			),
		},
		"UpdateNode": {
			m.assertUserBanned(),
			m.assertNodeBanned(
				func(req interface{}) (nodeid string) {
					return req.(drip.UpdateNodeRequestObject).NodeId
				},
			),
			m.assertPublisherPermission(
				[]schema.PublisherPermissionType{schema.PublisherPermissionTypeOwner},
				func(req interface{}) (publisherID string) {
					return req.(drip.UpdateNodeRequestObject).PublisherId
				},
			),
			m.assertNodeBelongsToPublisher(
				func(req interface{}) (publisherID string) {
					return req.(drip.UpdateNodeRequestObject).PublisherId
				},
				func(req interface{}) (nodeID string) {
					return req.(drip.UpdateNodeRequestObject).NodeId
				},
			),
		},
		"GetNode": {
			m.assertNodeBanned(
				func(req interface{}) (nodeid string) {
					return req.(drip.GetNodeRequestObject).NodeId
				},
			),
		},
		"PublishNodeVersion": {
			m.assertPublisherBanned(
				func(req interface{}) (publisherID string) {
					return req.(drip.PublishNodeVersionRequestObject).PublisherId
				}),
			m.assertNodeBanned(
				func(req interface{}) (nodeid string) {
					return req.(drip.PublishNodeVersionRequestObject).NodeId
				},
			),
			m.assertNodeBelongsToPublisher(
				func(req interface{}) (publisherID string) {
					return req.(drip.PublishNodeVersionRequestObject).PublisherId
				},
				func(req interface{}) (NodeId string) {
					return req.(drip.PublishNodeVersionRequestObject).NodeId
				},
			),
			m.assertPersonalAccessTokenValid(
				func(req interface{}) (publisherID string) {
					return req.(drip.PublishNodeVersionRequestObject).PublisherId
				},
				func(req interface{}) (pat string) {
					return req.(drip.PublishNodeVersionRequestObject).Body.PersonalAccessToken
				},
			),
		},
		"UpdateNodeVersion": {
			m.assertUserBanned(),
			m.assertNodeBanned(
				func(req interface{}) (nodeid string) {
					return req.(drip.UpdateNodeVersionRequestObject).NodeId
				},
			),
			m.assertPublisherPermission(
				[]schema.PublisherPermissionType{schema.PublisherPermissionTypeOwner},
				func(req interface{}) (publisherID string) {
					return req.(drip.UpdateNodeVersionRequestObject).PublisherId
				},
			),
			m.assertNodeBelongsToPublisher(
				func(req interface{}) (publisherID string) {
					return req.(drip.UpdateNodeVersionRequestObject).PublisherId
				},
				func(req interface{}) (publisherID string) {
					return req.(drip.UpdateNodeVersionRequestObject).NodeId
				},
			),
		},
		"DeleteNodeVersion": {
			m.assertUserBanned(),
			m.assertNodeBanned(
				func(req interface{}) (nodeid string) {
					return req.(drip.DeleteNodeVersionRequestObject).NodeId
				},
			),
		},
		"GetNodeVersion": {
			m.assertNodeBanned(
				func(req interface{}) (nodeid string) {
					return req.(drip.GetNodeVersionRequestObject).NodeId
				},
			),
		},
		"ListNodeVersions": {
			m.assertNodeBanned(
				func(req interface{}) (nodeid string) {
					return req.(drip.ListNodeVersionsRequestObject).NodeId
				},
			),
		},
		"InstallNode": {
			m.assertNodeBanned(
				func(req interface{}) (nodeid string) {
					return req.(drip.InstallNodeRequestObject).NodeId
				},
			),
		},
		"CreatePersonalAccessToken": {
			m.assertUserBanned(),
			m.assertPublisherPermission(
				[]schema.PublisherPermissionType{schema.PublisherPermissionTypeOwner},
				func(req interface{}) (publisherID string) {
					return req.(drip.CreatePersonalAccessTokenRequestObject).PublisherId
				},
			),
		},
		"DeletePersonalAccessToken": {
			m.assertUserBanned(),
			m.assertPublisherPermission(
				[]schema.PublisherPermissionType{schema.PublisherPermissionTypeOwner},
				func(req interface{}) (publisherID string) {
					return req.(drip.DeletePersonalAccessTokenRequestObject).PublisherId
				},
			),
		},
		"ListPersonalAccessTokens": {
			m.assertPublisherPermission(
				[]schema.PublisherPermissionType{schema.PublisherPermissionTypeOwner},
				func(req interface{}) (publisherID string) {
					return req.(drip.ListPersonalAccessTokensRequestObject).PublisherId
				},
			),
		},
		"GetPermissionOnPublisherNodes": {
			m.assertPublisherPermission(
				[]schema.PublisherPermissionType{schema.PublisherPermissionTypeOwner},
				func(req interface{}) (publisherID string) {
					return req.(drip.GetPermissionOnPublisherNodesRequestObject).PublisherId
				},
			),
		},
		"GetPermissionOnPublisher": {
			m.assertPublisherPermission(
				[]schema.PublisherPermissionType{schema.PublisherPermissionTypeOwner},
				func(req interface{}) (publisherID string) {
					return req.(drip.GetPermissionOnPublisherRequestObject).PublisherId
				},
			),
		},
	}
	for _, v := range subMiddlewares {
		slices.Reverse(v)
	}

	return func(f strictecho.StrictEchoHandlerFunc, operationID string) strictecho.StrictEchoHandlerFunc {
		middlewares, ok := subMiddlewares[operationID]
		if !ok {
			return f
		}

		for _, mw := range middlewares {
			f = mw(f, operationID)
		}
		return f
	}
}

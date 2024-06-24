package mapper

import (
	"fmt"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/schema"

	"github.com/Masterminds/semver/v3"
	"github.com/google/uuid"
)

func ApiUpdateNodeVersionToUpdateFields(versionId string, updateRequest *drip.NodeVersionUpdateRequest, client *ent.Client) *ent.NodeVersionUpdateOne {
	update := client.NodeVersion.UpdateOneID(uuid.MustParse(versionId))
	if updateRequest.Changelog != nil {
		update.SetChangelog(*updateRequest.Changelog)
	}
	if updateRequest.Deprecated != nil {
		update.SetDeprecated(*updateRequest.Deprecated)
	}
	return update
}

func ValidatePublishNodeVersionRequest(request drip.PublishNodeVersionRequestObject) error {
	if request.NodeId != *request.Body.Node.Id {
		return fmt.Errorf("node ID in URL and body must be the same")
	}

	return nil
}

func ApiCreateNodeVersionToDb(nodeId string, nodeVersion *drip.NodeVersion, client *ent.Client) *ent.NodeVersionCreate {
	create := client.NodeVersion.Create()
	if nodeId != "" {
		create.SetNodeID(nodeId)
	}
	if nodeVersion.Version != nil {
		create.SetVersion(*nodeVersion.Version)
	}
	if nodeVersion.Changelog != nil {
		create.SetChangelog(*nodeVersion.Changelog)
	}
	if nodeVersion.Dependencies != nil {
		create.SetPipDependencies(*nodeVersion.Dependencies)
	}

	return create
}

func DbNodeVersionToApiNodeVersion(dbNodeVersion *ent.NodeVersion) *drip.NodeVersion {
	if dbNodeVersion == nil {
		return nil
	}

	id := dbNodeVersion.ID.String()
	var downloadUrl string
	status := DbNodeVersionStatusToApiNodeVersionStatus(dbNodeVersion.Status)
	if dbNodeVersion.Edges.StorageFile != nil {
		downloadUrl = dbNodeVersion.Edges.StorageFile.FileURL
	}

	apiVersion := &drip.NodeVersion{
		Id:           &id,
		Version:      &dbNodeVersion.Version,
		Changelog:    &dbNodeVersion.Changelog,
		Deprecated:   &dbNodeVersion.Deprecated,
		Dependencies: &dbNodeVersion.PipDependencies,
		CreatedAt:    &dbNodeVersion.CreateTime,
		Status:       status,
		StatusReason: &dbNodeVersion.StatusReason,
		DownloadUrl:  &downloadUrl,
	}
	return apiVersion
}

func CheckValidSemv(version string) bool {
	_, err := semver.NewVersion(version)
	return err == nil
}

func DbNodeVersionStatusToApiNodeVersionStatus(status schema.NodeVersionStatus) *drip.NodeVersionStatus {
	var nodeVersionStatus drip.NodeVersionStatus

	switch status {
	case schema.NodeVersionStatusActive:
		nodeVersionStatus = drip.NodeVersionStatusActive
	case schema.NodeVersionStatusBanned:
		nodeVersionStatus = drip.NodeVersionStatusBanned
	case schema.NodeVersionStatusDeleted:
		nodeVersionStatus = drip.NodeVersionStatusDeleted
	case schema.NodeVersionStatusPending:
		nodeVersionStatus = drip.NodeVersionStatusPending
	case schema.NodeVersionStatusFlagged:
		nodeVersionStatus = drip.NodeVersionStatusFlagged
	default:
		nodeVersionStatus = ""
	}

	return &nodeVersionStatus
}

func ApiNodeVersionStatusesToDbNodeVersionStatuses(status *[]drip.NodeVersionStatus) []schema.NodeVersionStatus {
	var nodeVersionStatus []schema.NodeVersionStatus

	if status == nil {
		return nodeVersionStatus
	}

	for _, s := range *status {
		dbNodeVersion := ApiNodeVersionStatusToDbNodeVersionStatus(s)
		nodeVersionStatus = append(nodeVersionStatus, dbNodeVersion)
	}

	return nodeVersionStatus
}

func ApiNodeVersionStatusToDbNodeVersionStatus(status drip.NodeVersionStatus) schema.NodeVersionStatus {
	var nodeVersionStatus schema.NodeVersionStatus

	switch status {
	case drip.NodeVersionStatusActive:
		nodeVersionStatus = schema.NodeVersionStatusActive
	case drip.NodeVersionStatusBanned:
		nodeVersionStatus = schema.NodeVersionStatusBanned
	case drip.NodeVersionStatusDeleted:
		nodeVersionStatus = schema.NodeVersionStatusDeleted
	case drip.NodeVersionStatusPending:
		nodeVersionStatus = schema.NodeVersionStatusPending
	case drip.NodeVersionStatusFlagged:
		nodeVersionStatus = schema.NodeVersionStatusFlagged
	default:
		nodeVersionStatus = ""
	}

	return nodeVersionStatus
}

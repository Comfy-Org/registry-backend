package mapper

import (
	"errors"
	"regexp"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/schema"
	"strings"
)

func ApiCreateNodeToDb(publisherId string, node *drip.Node, client *ent.Client) (*ent.NodeCreate, error) {
	newNode := client.Node.Create()
	newNode.SetPublisherID(publisherId)
	if node.Description != nil {
		newNode.SetDescription(*node.Description)
	}
	if node.Id != nil {
		newNode.SetID(*node.Id)
		newNode.SetNormalizedID(normalizeNodeID(*node.Id))
	}
	if node.Author != nil {
		newNode.SetAuthor(*node.Author)
	}
	if node.License != nil {
		newNode.SetLicense(*node.License)
	}
	if node.Name != nil {
		newNode.SetName(*node.Name)
	}
	if node.Category != nil {
		newNode.SetCategory(*node.Category)
	}
	if node.Tags != nil {
		newNode.SetTags(*node.Tags)
	}
	if node.Repository != nil {
		newNode.SetRepositoryURL(*node.Repository)
	}
	if node.Icon != nil {
		newNode.SetIconURL(*node.Icon)
	}

	return newNode, nil
}

func ApiUpdateNodeToUpdateFields(nodeID string, node *drip.Node, client *ent.Client) *ent.NodeUpdateOne {
	update := client.Node.UpdateOneID(nodeID)
	if node.Description != nil {
		update.SetDescription(*node.Description)
	}
	if node.Author != nil {
		update.SetAuthor(*node.Author)
	}
	if node.License != nil {
		update.SetLicense(*node.License)
	}
	if node.Name != nil {
		update.SetName(*node.Name)
	}
	if node.Tags != nil {
		update.SetTags(*node.Tags)
	}
	if node.Category != nil {
		update.SetCategory(*node.Category)
	}
	if node.Repository != nil {
		update.SetRepositoryURL(*node.Repository)
	}
	if node.Icon != nil {
		update.SetIconURL(*node.Icon)
	}

	return update
}

func ValidateNode(node *drip.Node) error {
	if node.Id == nil {
		return errors.New("node id is required")
	}

	IsValid, errMsg := IsValidNodeID(*node.Id)
	if !IsValid {
		return errors.New(errMsg)
	}

	if node.Description != nil && len(*node.Description) > 1000 {
		return errors.New("description is too long")
	}

	return nil
}

func IsValidNodeID(nodeID string) (bool, string) {
	// Ensure the ID is not empty and doesn't exceed 100 characters
	if len(nodeID) == 0 || len(nodeID) > 100 {
		return false, "node id must be between 1 and 100 characters"
	}

	// Ensure the name does not start or end with an underscore, hyphen, or period
	if strings.HasPrefix(nodeID, "_") || strings.HasPrefix(nodeID, "-") ||
		strings.HasPrefix(nodeID, ".") || strings.HasSuffix(nodeID, "_") ||
		strings.HasSuffix(nodeID, "-") || strings.HasSuffix(nodeID, ".") {
		return false, "node id must not start or end with an underscore, hyphen, or period"
	}

	// Regular expression pattern to validate the project name
	// ASCII letters, digits, underscores, hyphens, and periods are allowed
	pattern := `^[a-zA-Z0-9](?:[a-zA-Z0-9._-]*[a-zA-Z0-9])?$`
	regex := regexp.MustCompile(pattern)

	// Validate against the pattern
	if !regex.MatchString(nodeID) {
		return false, "node id can only contain ASCII letters, digits, " +
			"underscores, hyphens, and periods, and must not have invalid sequences"
	}

	// Additional validation for normalized equivalency can be added here if needed
	return true, ""
}

// NormalizeNodeID normalizes the node ID to lowercase
func normalizeNodeID(nodeID string) string {
	// TODO: consider normalizing the node ID to a specific format adhering to
	// 	https://packaging.python.org/en/latest/guides/writing-pyproject-toml/#name, replacing
	// 	underscores, hyphens, and periods with a single hyphen
	// regex := regexp.MustCompile(`[_\-.]+`)
	// return regex.ReplaceAllString(strings.ToLower(nodeID), "-")
	return strings.ToLower(nodeID)
}

func DbNodeToApiNode(node *ent.Node) *drip.Node {
	if node == nil {
		return nil
	}

	downloads := int(node.TotalInstall)
	rate := float32(0)
	if node.TotalReview > 0 {
		rate = float32(node.TotalStar) / float32(node.TotalReview)
	}

	return &drip.Node{
		Author:       &node.Author,
		Description:  &node.Description,
		Category:     &node.Category,
		Id:           &node.ID,
		License:      &node.License,
		Name:         &node.Name,
		Tags:         &node.Tags,
		Repository:   &node.RepositoryURL,
		Icon:         &node.IconURL,
		Downloads:    &downloads,
		Rating:       &rate,
		Status:       DbNodeStatusToApiNodeStatus(node.Status),
		StatusDetail: &node.StatusDetail,
	}
}

func DbNodeStatusToApiNodeStatus(status schema.NodeStatus) *drip.NodeStatus {
	var nodeStatus drip.NodeStatus

	switch status {
	case schema.NodeStatusActive:
		nodeStatus = drip.NodeStatusActive
	case schema.NodeStatusBanned:
		nodeStatus = drip.NodeStatusBanned
	case schema.NodeStatusDeleted:
		nodeStatus = drip.NodeStatusDeleted
	default:
		nodeStatus = ""
	}

	return &nodeStatus
}

package mapper

import (
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
		lowerCaseNodeID := strings.ToLower(*node.Id)
		newNode.SetID(lowerCaseNodeID)
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
	if node.Id != nil {
		if len(*node.Id) > 100 {
			return NewErrorBadRequest("node id is too long")
		}
		isValid, msg := IsValidNodeID(*node.Id)
		if !isValid {
			return NewErrorBadRequest(msg)
		}
	}
	if node.Description != nil {
		if len(*node.Description) > 1000 {
			return NewErrorBadRequest("description is too long")
		}
	}
	return nil
}

func IsValidNodeID(nodeID string) (bool, string) {
	if len(nodeID) == 0 || len(nodeID) > 50 {
		return false, "node id must be between 1 and 50 characters"
	}
	// Check if there are capital letters in the string
	if strings.ToLower(nodeID) != nodeID {
		return false, "Node ID can only contain lowercase letters"
	}
	// Regular expression pattern for Node ID validation (lowercase letters only)
	pattern := `^[a-z][a-z0-9-_]+(\.[a-z0-9-_]+)*$`
	// Compile the regular expression pattern
	regex := regexp.MustCompile(pattern)
	// Check if the string matches the pattern
	matches := regex.MatchString(nodeID)
	if !matches {
		return false, "Node ID can only contain lowercase letters, numbers, hyphens, underscores, and dots. Dots cannot be consecutive or be at the start or end of the id."
	}
	return true, ""
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

package entity

import (
	"registry-backend/ent"
	"registry-backend/ent/schema"
	"time"
)

type AlgoliaNode struct {
	ObjectID string `json:"objectID"`

	ID            string            `json:"id,omitempty"`
	CreateTime    time.Time         `json:"create_time,omitempty"`
	UpdateTime    time.Time         `json:"update_time,omitempty"`
	PublisherID   string            `json:"publisher_id,omitempty"`
	Name          string            `json:"name,omitempty"`
	Description   string            `json:"description,omitempty"`
	Category      string            `json:"category,omitempty"`
	Author        string            `json:"author,omitempty"`
	License       string            `json:"license,omitempty"`
	RepositoryURL string            `json:"repository_url,omitempty"`
	IconURL       string            `json:"icon_url,omitempty"`
	Tags          []string          `json:"tags,omitempty"`
	TotalInstall  int64             `json:"total_install,omitempty"`
	TotalStar     int64             `json:"total_star,omitempty"`
	TotalReview   int64             `json:"total_review,omitempty"`
	Status        schema.NodeStatus `json:"status,omitempty"`
	StatusDetail  string            `json:"status_detail,omitempty"`

	LatestVersion       string                   `json:"latest_version,omitempty"`
	LatestVersionStatus schema.NodeVersionStatus `json:"latest_version_status,omitempty"`
	ComfyNodeNames      []string                 `json:"comfy_nodes,omitempty"`
}

func (n *AlgoliaNode) ToEntNode() *ent.Node {
	node := &ent.Node{
		ID:            n.ID,
		CreateTime:    n.CreateTime,
		UpdateTime:    n.UpdateTime,
		PublisherID:   n.PublisherID,
		Name:          n.Name,
		Description:   n.Description,
		Category:      n.Category,
		Author:        n.Author,
		License:       n.License,
		RepositoryURL: n.RepositoryURL,
		IconURL:       n.IconURL,
		Tags:          n.Tags,
		TotalInstall:  n.TotalInstall,
		TotalStar:     n.TotalStar,
		TotalReview:   n.TotalReview,
		Status:        n.Status,
		StatusDetail:  n.StatusDetail,
	}
	if n.LatestVersion == "" {
		return node
	}

	node.Edges = ent.NodeEdges{
		Versions: []*ent.NodeVersion{{
			NodeID:  n.ID,
			Version: n.LatestVersion,
			Status:  n.LatestVersionStatus,
		}},
	}
	if len(n.ComfyNodeNames) == 0 {
		return node
	}

	node.Edges.Versions[0].Edges = ent.NodeVersionEdges{
		ComfyNodes: make([]*ent.ComfyNode, 0, len(n.ComfyNodeNames)),
	}
	for _, name := range n.ComfyNodeNames {
		node.Edges.Versions[0].Edges.ComfyNodes = append(node.Edges.Versions[0].Edges.ComfyNodes, &ent.ComfyNode{
			Name: name,
		})
	}
	return node
}

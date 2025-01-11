package mapper

import (
	"registry-backend/ent"
	"registry-backend/entity"
)

func AlgoliaNodeFromEntNode(node *ent.Node) entity.AlgoliaNode {
	n := entity.AlgoliaNode{
		ObjectID:      node.ID,
		ID:            node.ID,
		CreateTime:    node.CreateTime,
		UpdateTime:    node.UpdateTime,
		PublisherID:   node.PublisherID,
		Name:          node.Name,
		Description:   node.Description,
		Category:      node.Category,
		Author:        node.Author,
		License:       node.License,
		RepositoryURL: node.RepositoryURL,
		IconURL:       node.IconURL,
		Tags:          node.Tags,
		TotalInstall:  node.TotalInstall,
		TotalStar:     node.TotalStar,
		TotalReview:   node.TotalReview,
		Status:        node.Status,
		StatusDetail:  node.StatusDetail,

		LatestVersion:       "",
		LatestVersionStatus: "",
	}

	if node.Edges.Versions == nil {
		return n
	}

	var lv *ent.NodeVersion
	for _, v := range node.Edges.Versions {
		if lv == nil {
			lv = v
		} else if v.CreateTime.After(lv.CreateTime) {
			lv = v
		}
	}
	n.LatestVersion = lv.Version
	n.LatestVersionStatus = lv.Status
	n.ComfyNodeNames = make([]string, 0, len(lv.Edges.ComfyNodes))
	for _, v := range lv.Edges.ComfyNodes {
		n.ComfyNodeNames = append(n.ComfyNodeNames, v.Name)
	}

	return n
}

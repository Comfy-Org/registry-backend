package mapper

import (
	"registry-backend/ent"
	"registry-backend/entity"
)

func AlgoliaNodeFromEntNode(node *ent.Node) entity.AlgoliaNode {
	n := entity.AlgoliaNode{
		ObjectID: node.ID,
		Node:     new(ent.Node),
	}
	*n.Node = *node
	n.Edges = ent.NodeEdges{}
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

	n.LatestVersion = &struct {
		*ent.NodeVersion
		ComfyNodes map[string]*ent.ComfyNode `json:"comfy_nodes"`
	}{
		NodeVersion: new(ent.NodeVersion),
		ComfyNodes:  make(map[string]*ent.ComfyNode, len(lv.Edges.ComfyNodes)),
	}
	*n.LatestVersion.NodeVersion = *lv
	n.LatestVersion.NodeVersion.Edges = ent.NodeVersionEdges{}
	for _, v := range lv.Edges.ComfyNodes {
		n.LatestVersion.ComfyNodes[v.ID] = v
	}

	return n
}

package entity

import "registry-backend/ent"

type AlgoliaNode struct {
	ObjectID string `json:"objectID"`
	*ent.Node
	LatestVersion *struct {
		*ent.NodeVersion
		ComfyNodes map[string]*ent.ComfyNode `json:"comfy_nodes"`
	} `json:"latest_version"`
}

func (n *AlgoliaNode) ToEntNode() *ent.Node {
	node := n.Node
	if n.LatestVersion != nil {
		nv := n.LatestVersion.NodeVersion
		for _, v := range n.LatestVersion.ComfyNodes {
			nv.Edges.ComfyNodes = append(nv.Edges.ComfyNodes, v)
		}
		node.Edges.Versions = []*ent.NodeVersion{nv}
	}
	return node
}

package entity

import "registry-backend/ent"

// NodeFilter holds optional parameters for filtering node results
type NodeFilter struct {
	PublisherID   string
	Search        string
	IncludeBanned bool
}

// ListNodesResult is the structure that holds the paginated result of nodes
type ListNodesResult struct {
	Total      int         `json:"total"`
	Nodes      []*ent.Node `json:"nodes"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"totalPages"`
}

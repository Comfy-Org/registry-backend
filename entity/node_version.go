package entity

import (
	"registry-backend/ent"
	"registry-backend/ent/schema"
	"time"
)

type NodeVersionFilter struct {
	NodeId              string
	Status              []schema.NodeVersionStatus
	IncludeStatusReason bool
	MinAge              time.Duration
	PageSize            int
	Page                int
}

type ListNodeVersionsResult struct {
	Total        int                `json:"total"`
	NodeVersions []*ent.NodeVersion `json:"nodes"`
	Page         int                `json:"page"`
	Limit        int                `json:"limit"`
	TotalPages   int                `json:"totalPages"`
}

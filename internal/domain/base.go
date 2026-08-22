package domain

import "time"

type Base struct {
 ID int64 `json:"id"`
 CreatedAt time.Time `json:"created_at"`
 UpdatedAt time.Time `json:"updated_at"`
 DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type PageQuery struct { Page int; PageSize int; Keyword string; Status string }
type StatusHistory struct { Base; EntityID int64 `json:"entity_id"`; FromStatus string `json:"from_status"`; ToStatus string `json:"to_status"`; Reason string `json:"reason"`; ChangedBy int64 `json:"changed_by"` }
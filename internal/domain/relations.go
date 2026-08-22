package domain

import "time"

type UserRole struct { UserID int64 `json:"user_id"`; RoleID int64 `json:"role_id"` }
type RolePermission struct { RoleID int64 `json:"role_id"`; PermissionID int64 `json:"permission_id"` }
type RefreshToken struct { Base; UserID int64 `json:"user_id"`; TokenHash string `json:"-"`; ExpiresAt time.Time `json:"expires_at"`; RevokedAt *time.Time `json:"revoked_at,omitempty"` }
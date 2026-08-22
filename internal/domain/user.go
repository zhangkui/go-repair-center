package domain

import "time"

type User struct {
    Base
    Username string `json:"username,omitempty"`
    PasswordHash string `json:"password_hash,omitempty"`
    DisplayName string `json:"display_name,omitempty"`
    Phone string `json:"phone,omitempty"`
    Email string `json:"email,omitempty"`
    Status string `json:"status,omitempty"`
    LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

func (m User) EntityID() int64 { return m.ID }
func (m User) IsDeleted() bool { return m.DeletedAt != nil }
func (m User) CreatedUnix() int64 { return m.CreatedAt.Unix() }
func (m User) UpdatedUnix() int64 { return m.UpdatedAt.Unix() }
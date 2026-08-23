package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type SessionService struct {
	db    *sql.DB
	cache *CacheService
}

type SessionInfo struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type SessionSummary struct {
	UserID              int64         `json:"user_id"`
	ActiveSessionCount  int           `json:"active_session_count"`
	ExpiredSessionCount int           `json:"expired_session_count"`
	Sessions            []SessionInfo `json:"sessions"`
}

func NewSessionService(db *sql.DB, cache *CacheService) *SessionService {
	return &SessionService{db: db, cache: cache}
}

func (s *SessionService) ListByUser(ctx context.Context, userID int64) ([]SessionInfo, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT id, user_id, expires_at, revoked_at, created_at
FROM refresh_tokens
WHERE user_id = ? AND deleted_at IS NULL
ORDER BY created_at DESC
`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]SessionInfo, 0, 8)
	for rows.Next() {
		var item SessionInfo
		if err := rows.Scan(&item.ID, &item.UserID, &item.ExpiresAt, &item.RevokedAt, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *SessionService) ListActiveByUser(ctx context.Context, userID int64) ([]SessionInfo, error) {
	items, err := s.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	active := make([]SessionInfo, 0, len(items))
	now := time.Now()
	for _, item := range items {
		if item.RevokedAt == nil && item.ExpiresAt.After(now) {
			active = append(active, item)
		}
	}
	return active, nil
}

func (s *SessionService) RevokeByUser(ctx context.Context, userID int64) error {
	_, err := s.db.ExecContext(ctx, `
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE user_id = ? AND revoked_at IS NULL AND deleted_at IS NULL
`, userID)
	if err == nil && s.cache != nil {
		_ = s.cache.Delete(ctx, s.userSessionKey(userID))
	}
	return err
}

func (s *SessionService) RevokeByID(ctx context.Context, sessionID int64) error {
	_, err := s.db.ExecContext(ctx, `
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE id = ? AND revoked_at IS NULL AND deleted_at IS NULL
`, sessionID)
	if err == nil && s.cache != nil {
		_ = s.cache.Delete(ctx, SingleSessionSummaryKey(sessionID))
	}
	return err
}

func (s *SessionService) CleanupExpired(ctx context.Context) (int64, error) {
	result, err := s.db.ExecContext(ctx, `
UPDATE refresh_tokens
SET deleted_at = NOW(), updated_at = NOW()
WHERE expires_at < NOW() AND deleted_at IS NULL
`)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *SessionService) Summary(ctx context.Context, userID int64) (*SessionSummary, error) {
	if s.cache != nil {
		cacheKey := s.userSessionKey(userID)
		cached := &SessionSummary{}
		if ok, err := s.cache.GetJSON(ctx, cacheKey, cached); err == nil && ok {
			return cached, nil
		}
	}

	items, err := s.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	summary := &SessionSummary{
		UserID:   userID,
		Sessions: items,
	}
	for _, item := range items {
		if item.RevokedAt == nil && item.ExpiresAt.After(now) {
			summary.ActiveSessionCount++
		}
		if item.ExpiresAt.Before(now) {
			summary.ExpiredSessionCount++
		}
	}
	if s.cache != nil {
		_ = s.cache.SetJSON(ctx, s.userSessionKey(userID), summary, time.Minute)
	}
	return summary, nil
}

func (s *SessionService) CountActive(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := s.db.QueryRowContext(ctx, `
SELECT COUNT(*)
FROM refresh_tokens
WHERE user_id = ?
  AND revoked_at IS NULL
  AND expires_at > NOW()
  AND deleted_at IS NULL
`, userID).Scan(&count)
	return count, err
}

func (s *SessionService) TouchCache(ctx context.Context, userID int64) error {
	if s.cache == nil {
		return nil
	}
	summary, err := s.Summary(ctx, userID)
	if err != nil {
		return err
	}
	return s.cache.SetJSON(ctx, s.userSessionKey(userID), summary, time.Minute)
}

func (s *SessionService) userSessionKey(userID int64) string {
	return "sessions:user:" + strconvI64(userID)
}

func strconvI64(value int64) string {
	return fmt.Sprintf("%d", value)
}

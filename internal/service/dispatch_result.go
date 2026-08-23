package service

import "go-repair-center/internal/repository"

func newBatchDispatchResult(capacity int) *BatchDispatchResult {
	if capacity < 0 {
		capacity = 0
	}
	return &BatchDispatchResult{
		Dispatched: make([]repository.Record, 0, capacity),
		Skipped:    make([]repository.Record, 0),
	}
}

func (r *BatchDispatchResult) addConflicts(records []repository.Record) {
	if len(records) == 0 {
		return
	}
	r.Conflicts = append(r.Conflicts, records...)
}

func (r *BatchDispatchResult) addNotification(value string) {
	r.Notifications = append(r.Notifications, value)
}

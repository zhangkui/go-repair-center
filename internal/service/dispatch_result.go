package service

func newBatchDispatchResult(capacity int) *BatchDispatchResult {
	return &BatchDispatchResult{
		Dispatched:    make([]repository.Record, 0, capacity),
		Skipped:       make([]repository.Record, 0),
		Conflicts:     make([]repository.Record, 0),
		Notifications: make([]string, 0),
	}
}

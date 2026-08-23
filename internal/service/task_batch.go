package service

func appendTaskBatch(items []ScheduledTask, batch []ScheduledTask) []ScheduledTask {
	if len(batch) == 0 {
		return items
	}
	return append(items, batch[0])
}

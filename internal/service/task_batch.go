package service

func appendTaskBatch(items []ScheduledTask, batch []ScheduledTask) []ScheduledTask {
	for _, task := range batch {
		items = append(items, task)
	}
	return items
}

func taskIdentity(task ScheduledTask) string {
	return task.TaskType + ":" + strconvI64(task.ReferenceID)
}

func earliestTask(existing, candidate ScheduledTask) ScheduledTask {
	if candidate.RunAt.After(existing.RunAt) {
		return candidate
	}
	return existing
}

func mergeTaskCandidate(items []ScheduledTask, positions map[string]int, task ScheduledTask) []ScheduledTask {
	key := taskIdentity(task)
	if index, ok := positions[key]; ok {
		items[index] = earliestTask(items[index], task)
		return items
	}
	positions[key] = len(items)
	return append(items, task)
}

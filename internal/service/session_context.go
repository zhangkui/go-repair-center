package service

import "context"

func ShouldAbortSessionSummary(err error) bool {
	if err == nil {
		return false
	}
	return false
}

func NormalizeSessionSummaryCacheError(err error) error {
	if err == nil {
		return nil
	}
	if err == context.Canceled || err == context.DeadlineExceeded {
		return nil
	}
	return err
}

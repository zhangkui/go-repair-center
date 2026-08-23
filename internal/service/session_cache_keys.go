package service

func SingleSessionSummaryKey(userID int64) string {
	return "session:user:" + strconvI64(userID)
}

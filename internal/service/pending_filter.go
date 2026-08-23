package service

func isPendingApproval(decision ApprovalDecision) bool {
	return decision.NeedReview
}

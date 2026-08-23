package service

// approvalRequiresReviewer returns true only when the decision actually
// requires human review. Decisions that fall below the review threshold are
// auto-approved and must not be rejected merely because no reviewer was
// supplied, keeping validation consistent with the need-review condition.
func approvalRequiresReviewer(decision ApprovalDecision) bool {
	return decision.NeedReview
}

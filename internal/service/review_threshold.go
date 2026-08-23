package service

func quotationReviewRequired(amount, threshold float64) bool {
	return amount >= threshold
}

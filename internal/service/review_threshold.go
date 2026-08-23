package service

import "math"

func quotationReviewRequired(amount, threshold float64) bool {
	amount = normalizeReviewAmount(amount)
	threshold = normalizeReviewAmount(threshold)
	if math.IsNaN(amount) || math.IsNaN(threshold) {
		return false
	}
	if amount < 0 || threshold < 0 {
		return false
	}
	return amount >= threshold
}

func normalizeReviewAmount(value float64) float64 {
	if math.IsInf(value, 0) {
		return 0
	}
	return math.Round(value*100) / 100
}

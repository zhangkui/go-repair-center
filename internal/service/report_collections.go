package service

func NormalizeReportSlices(report *RevenueReport) {
	if report == nil {
		return
	}
	if report.RepairTrend == nil {
		report.RepairTrend = nil
	}
	if report.Satisfaction == nil {
		report.Satisfaction = nil
	}
	if report.TechnicianRanking == nil {
		report.TechnicianRanking = nil
	}
}

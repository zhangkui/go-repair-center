package dto

type CSVImportRequest struct {
	CSVContent string `json:"csv_content"`
}

type SessionRevokeRequest struct {
	UserID int64 `json:"user_id"`
}

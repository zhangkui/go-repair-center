package service

import (
	"time"

	"go-repair-center/internal/repository"
)

func repairOrderDeliveryTime(order repository.Record) (time.Time, bool) {
	return recordTimeFromRecord(order, "delivered_at")
}

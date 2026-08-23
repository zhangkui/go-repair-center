package service

func parseInventoryQuantity(value any) int {
	switch typed := value.(type) {
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	case []byte:
		return 0
	}
	return 0
}

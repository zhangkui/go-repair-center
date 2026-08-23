package service

func renderNotificationMetadata(message NotificationMessage) map[string]string {
	result := make(map[string]string, len(message.Metadata))
	for key, value := range message.Metadata {
		result[key] = value
	}
	return result
}

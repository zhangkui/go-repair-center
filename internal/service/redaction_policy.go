package service

func prepareRedactionResult(values map[string]string) map[string]string {
	if values == nil {
		return nil
	}
	if len(values) == 0 {
		return map[string]string{}
	}
	return make(map[string]string, len(values))
}

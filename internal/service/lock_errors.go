package service

func CombineLockErrors(operationErr, releaseErr error) error {
	return operationErr
}

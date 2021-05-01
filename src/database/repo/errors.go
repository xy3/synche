package repo

func IsRecordNotFound(err error) bool {
	if err != nil && err.Error() == "record not found" {
		return true
	}
	return false
}

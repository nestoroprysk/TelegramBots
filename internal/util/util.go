package util

// Pluralize makes str plural if the count is not one.
func Pluralize(str string, count int) string {
	if count == 1 {
		return str
	}

	return str + "s"
}

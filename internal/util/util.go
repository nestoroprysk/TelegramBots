package util

import "fmt"

// Pluralize makes str plural if the count is not one.
func Pluralize(str string, count int) string {
	if count == 1 {
		return str
	}

	return str + "s"
}

// Format code formats str as code in MD.
func FormatCode(s interface{}) string {
	return fmt.Sprintf("```\n%+v\n```", s)
}

func CombineErrors(err error, errs ...error) error {
	result := err
	for _, err := range errs {
		result = combineErrors(result, err)
	}

	return result
}

func combineErrors(a, b error) error {
	return fmt.Errorf("%w\n%w", a, b)
}

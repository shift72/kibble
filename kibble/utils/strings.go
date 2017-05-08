package utils

// Join - using the separator join the strings
func Join(separator string, values ...string) string {

	l := len(values)
	cleaned := make([]string, 0)
	for i := 0; i < l; i++ {
		if len(values[i]) > 0 {
			cleaned = append(cleaned, values[i])
		}
	}

	result := ""
	l = len(cleaned)
	for i := 0; i < l; i++ {
		result += cleaned[i]
		if i+1 != l {
			result += separator
		}
	}

	return result
}

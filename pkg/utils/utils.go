package utils

func MapGet(m map[string]string, key string) string {
	if m == nil {
		return ""
	}
	return m[key]
}

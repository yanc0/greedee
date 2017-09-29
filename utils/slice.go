package utils

// Contains return true if str is found in s
func ContainsString(s []string, str string) bool {
	for _, st := range s {
		if st == str {
			return true
		}
	}
	return false
}

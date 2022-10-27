package golang_kong_access

// stringArrayContains checks if the string s is present in the string array a
func stringArrayContains(a []string, s string) bool {
	for _, item := range a {
		if item == s {
			return true
		}
	}
	return false
}

package golang_kong_access

func stringArrayContains(a []string, s string) bool {
	for _, item := range a {
		if item == s {
			return true
		}
	}
	return false
}

package model

func convertToSet(arr []string) map[string]bool {
	set := make(map[string]bool)
	for _, s := range arr {
		set[s] = true
	}
	return set
}

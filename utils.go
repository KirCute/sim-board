package sim_board

func SliceContains[T comparable](slice []T, v T) bool {
	for _, v1 := range slice {
		if v == v1 {
			return true
		}
	}
	return false
}

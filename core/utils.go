package core

// FindOnArray search over string array, if value exists on array returns at position, else return -1
func FindOnArray(array []string, value string) int {
	for k, v := range array {
		if v == value {
			return k
		}
	}

	return -1
}
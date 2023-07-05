package helpers

func GetKeys[T any](mp map[string]T) []string {
	size := len(mp)
	ans := make([]string, size)

	i := 0
	for key := range mp {
		ans[i] = key
		i++
	}

	return ans
}

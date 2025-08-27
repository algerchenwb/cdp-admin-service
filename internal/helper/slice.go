package helper

func ArrayUniqueValue[T comparable](arr []T) []T {
	m := make(map[T]struct{})
	for _, v := range arr {
		m[v] = struct{}{}
	}

	result := make([]T, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

package slice

// SetDifference returns a slice containing the elements in slice a that are not in slice b.
func SetDifference[T comparable](a, b []T, unique bool) []T {
	result := make([]T, 0)

	if unique {
		a = Unique(a)
	}

	for _, aItem := range a {
		found := false

		for _, bItem := range b {
			if aItem == bItem {
				found = true

				break
			}
		}

		if !found {
			result = append(result, aItem)
		}
	}

	return result
}

// Unique returns a slice containing only the unique elements of the given slice.
func Unique[T comparable](s []T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0)

	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}

			result = append(result, v)
		}
	}

	return result
}

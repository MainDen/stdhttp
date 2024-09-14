package slicesx

func Select[S ~[]E, E any](s S, predicate func(E) bool) S {
	var result S
	for _, e := range s {
		if predicate(e) {
			result = append(result, e)
		}
	}
	return result
}

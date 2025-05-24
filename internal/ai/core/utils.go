package core

// AllMatch returns true if all elements satisfy the predicate.
// Returns false if any element fails the predicate.
// Returns error if predicate returns an error.
func AllMatch[T any](collection []T, predicate func(T) (bool, error)) (bool, error) {
	for _, item := range collection {
		result, err := predicate(item)
		if err != nil {
			return false, err
		}
		if !result {
			return false, nil
		}
	}

	return true, nil
}

// AnyMatch returns true if any element satisfies the predicate.
// Returns false if no element satisfies the predicate.
// Returns error if predicate returns an error.
func AnyMatch[T any](collection []T, predicate func(T) (bool, error)) (bool, error) {
	for _, item := range collection {
		result, err := predicate(item)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}

	return false, nil
}

// Count returns the number of elements that satisfy the predicate.
// Returns error if predicate returns an error.
func Count[T any](collection []T, predicate func(T) (bool, error)) (int, error) {
	count := 0
	for _, item := range collection {
		result, err := predicate(item)
		if err != nil {
			return 0, err
		}
		if result {
			count++
		}
	}

	return count, nil
}

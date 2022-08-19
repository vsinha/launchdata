package slices

// TODO Add this to the official slices library

// From: https://github.com/golang/go/wiki/SliceTricks#reversing
// Reverses an array in place
func Reverse[E any](a []E) {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
}

// Returned a newly allocated Reversed array
func Reversed[E any](a []E) []E {
	b := make([]E, len(a))
	for i := len(a) - 1; i >= 0; i-- {
		b[len(a)-1-i] = a[i]
	}

	return b
}

func Insert[S ~[]E, E any](slice S, index int, v ...E) S {
	tot := len(slice) + len(v)
	if tot <= cap(slice) {
		s2 := slice[:tot]
		copy(s2[index+len(v):], slice[index:])
		copy(s2[index:], v)
		return s2
	}
	s2 := make(S, tot)
	copy(s2, slice[:index])
	copy(s2[index:], v)
	copy(s2[index+len(v):], slice[index:])
	return s2
}

// Delete removes the elements s[i:j] from s, returning the modified slice.
// Delete panics if s[i:j] is not a valid slice of s.
// Delete modifies the contents of the slice s; it does not create a new slice.
// Delete is O(len(s)-(j-i)), so if many items must be deleted, it is better to
// make a single call deleting them all together than to delete one at a time.
func Delete[S ~[]E, E any](slice S, i, j int) S {
	return append(slice[:i], slice[j:]...)
}

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

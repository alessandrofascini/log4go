package utils

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func findIndexes(source string, finalLength int) (int, int) {
	if finalLength > 0 {
		return 0, min(len(source), finalLength)
	}
	n := len(source)
	// n+finalLength (finalLength is negative!)
	return max(n+finalLength, 0), n
}

func Truncate(s string, k int) string {
	start, end := findIndexes(s, k)
	return s[start:end]
}

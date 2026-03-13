package utils

// isValidFibonacci checks if a value is in the valid Fibonacci sequence for P-WVC
func IsValidFibonacci(value int) bool {
	validValues := []int{1, 2, 3, 5, 8, 13, 21, 34, 55, 89}
	for _, v := range validValues {
		if v == value {
			return true
		}
	}
	return false
}

// GetFibonacciSequence returns the valid Fibonacci sequence for P-WVC
func GetFibonacciSequence() []int {
	return []int{1, 2, 3, 5, 8, 13, 21, 34, 55, 89}
}

package util

//Clamp function
func Clamp(value int, min int, max int) int {
	if value > max {
		return max
	}
	if value < min {
		return min
	}
	return value
}

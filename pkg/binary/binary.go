package binary

import "math"

func PowerOfTwo(num int) []int {
	result := []int{}
	power := 0
	for num > 0 {
		if num%2 == 1 {
			result = append(result, int(math.Pow(2, float64(power))))
		}
		power++
		num = int(math.Floor(float64(num / 2)))
	}
	return result
}

package helpers

// Min number in the array
func Min(ints []int) int {
	var smallest = ints[0]
	for _, v := range ints {
		if v < smallest {
			smallest = v
		}
	}
	return smallest
}

// Max number in the array
func Max(ints []int) int {
	var biggest = ints[0]
	for _, v := range ints {
		if v > biggest {
			biggest = v
		}
	}
	return biggest
}

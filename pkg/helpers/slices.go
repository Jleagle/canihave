package helpers

func RemoveDuplicates(elements []string) []string {

	encountered := map[string]bool{}
	for v := range elements {
		encountered[elements[v]] = true
	}

	var result []string
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

func Reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func RemoveFromSlice(current []string, remove string) (ret []string) {
	for _, str := range current {
		if str != remove {
			ret = append(ret, str)
		}
	}
	return ret
}

func InSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

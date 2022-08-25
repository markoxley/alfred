package utils

func ListIndex(hs []string, n string) int {
	for i, v := range hs {
		if v == n {
			return i
		}
	}
	return -1
}

// listContains is a helper function to test if a list of strings [hs] contains a specified string [n]
//  @param hs The 'haystack' list of words
//  @param n The 'needle' to find
//  @return bool True is the word is found
func ListContains(hs []string, n string) bool {
	for _, i := range hs {
		if i == n {
			return true
		}
	}
	return false
}

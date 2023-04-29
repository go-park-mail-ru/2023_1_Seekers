package common

func GetFirstUtf(str string) string {
	for _, c := range str {
		return string(c)
	}
	return ""
}

func GetFirstRune(str string) rune {
	for _, c := range str {
		return c
	}
	return ' '
}

package pkg

func GetFirstUtf(str string) string {
	for _, c := range str {
		return string(c)
	}
	return ""
}

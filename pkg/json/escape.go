package json

import (
	"strings"
)

func Escape(i string) (string, error) {
	return strings.ReplaceAll(i, "\"", "\\\""), nil
	//b, err := json.Marshal(i)
	//if err != nil {
	//	return "", err
	//}
	//s := string(b)
	//return s[1 : len(s)-1], nil
}

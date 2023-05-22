package common

import (
	"golang.org/x/net/html"
	"strings"
)

func GetInnerText(htmlText string, maxLen int) string {
	tkn := html.NewTokenizer(strings.NewReader(htmlText))

	var data string

	for {
		tt := tkn.Next()
		if len(data) > maxLen {
			data = strings.TrimLeft(data, " ")
			return string([]rune(data)[:maxLen])
		}

		switch {
		case tt == html.ErrorToken:
			data = strings.TrimLeft(data, " ")
			if len(data) > maxLen {
				return string([]rune(data)[:maxLen])
			} else {
				return data
			}

		case tt == html.TextToken:
			t := tkn.Token()
			tokData := t.Data
			if strings.TrimSpace(tokData) != "" {
				data += " " + tokData
			}
		}
	}
}

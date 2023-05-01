package smtp

import (
	"fmt"
	"strings"
)

func ParseDomain(emailAddr string) (string, error) {
	result := strings.Split(emailAddr, "@")
	if len(result) != 2 {
		return "", fmt.Errorf("invalid address | %s", emailAddr)
	}
	domain := result[1]
	return domain, nil
}

func ParseLogin(emailAddr string) (string, error) {
	result := strings.Split(emailAddr, "@")
	if len(result) != 2 {
		return "", fmt.Errorf("invalid address | %s", emailAddr)
	}
	login := result[0]
	return login, nil
}

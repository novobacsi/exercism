package phonenumber

import (
	"errors"
	"fmt"
	"regexp"
)

var nonDigits = regexp.MustCompile(`[^0-9]`)

func Number(phoneNumber string) (string, error) {
	s := nonDigits.ReplaceAllString(phoneNumber, "")

	if len(s) == 11 {
		if s[0] != '1' {
			return "", errors.New("invalid country code prefix")
		}

		s = s[1:]
	}

	if len(s) != 10 {
		return "", errors.New("incorrect length")
	}

	if s[0] < '2' {
		return "", errors.New("invalid area code prefix")
	}

	if s[3] < '2' {
		return "", errors.New("invalid exchange code prefix")
	}

	return s, nil
}

func AreaCode(phoneNumber string) (string, error) {
	n, err := Number(phoneNumber)
	if err != nil {
		return "", err
	}

	return n[:3], nil
}

func Format(phoneNumber string) (string, error) {
	n, err := Number(phoneNumber)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(%s) %s-%s", n[:3], n[3:6], n[6:]), nil
}

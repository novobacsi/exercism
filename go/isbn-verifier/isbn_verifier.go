package isbnverifier

import "strings"

func IsValidISBN(isbn string) bool {
	s := strings.ReplaceAll(isbn, "-", "")

	if len(s) != 10 {
		return false
	}

	sum := 0
	for i := 0; i < len(s); i++ {
		v := 0
		switch {
		case s[i] >= '0' && s[i] <= '9':
			v = int(s[i] - '0')
		case s[i] == 'X' && i == 9:
			v = 10
		default:
			return false
		}
		sum += v * (10 - i)
	}

	return sum%11 == 0
}

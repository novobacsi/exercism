package luhn

import "strings"
import "fmt"

func Valid(id string) bool {
	id = strings.ReplaceAll(id, " ", "")

	if len(id) < 2 {
		return false
	}

	sum := 0
	for i := 0; i < len(id)-1; i++ {
		c := id[i]
		x := int(c) * 2
		if x > 9 {
			x = x - 9
		}
		sum += x
		fmt.Println(rune(sum))
	}

	return sum%10 == 0

}

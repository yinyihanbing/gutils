package sensitive

import "strings"

// 如果是大写字母, 统一转成小写
func CheckToLower(c rune) rune {
	str := []rune(strings.ToLower(string(c)))
	if len(str) > 0 {
		return str[0]
	}
	return c
}

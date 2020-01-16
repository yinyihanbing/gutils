package gutils

import (
	"strings"
	"strconv"
	"errors"
)

// 把字符串打散为切片
func SplitExtInt32(separator string, content string) ([]int32, error) {
	arr := strings.Split(content, separator)
	items := make([]int32, 0, len(arr))

	for _, v := range arr {
		if v != "" {
			item, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			items = append(items, int32(item))
		}
	}
	return items, nil
}

// 把字符串打散为切片
func SplitExtInt32Interface(separator string, content string) ([]interface{}, error) {
	arr := strings.Split(content, separator)
	items := make([]interface{}, 0, len(arr))

	for _, v := range arr {
		if v != "" {
			item, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			items = append(items, int32(item))
		}
	}
	return items, nil
}

// 把字符串打散为Map
func SplitExtMapInt32(separator1 string, separator2 string, content string) (map[int32]int32, error) {
	arr := strings.Split(content, separator1)
	items := map[int32]int32{}

	for _, v := range arr {
		if v != "" {
			arrStrItem := strings.Split(v, separator2)
			if len(arrStrItem) != 2 {
				return nil, errors.New("formatting error")
			}

			item1, err := strconv.Atoi(arrStrItem[0])
			if err != nil {
				return nil, err
			}

			item2, err := strconv.Atoi(arrStrItem[1])
			if err != nil {
				return nil, err
			}

			items[int32(item1)] = int32(item2)
		}
	}
	return items, nil
}

// 获取source的子串,如果start小于0或者end大于source长度则返回"", start:开始index，从0开始，包括0, end:结束index，以end结束，但不包括end
func SubString(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	return string(r[start:end])
}

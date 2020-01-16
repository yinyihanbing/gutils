package gutils

import "reflect"

// 获取切片中的大于等于指定值的数量
func GetSMinVCountI32(src []int32, minValue int32) int32 {
	if len(src) == 0 {
		return 0
	}
	var result int32 = 0
	for _, v := range src {
		if v >= minValue {
			result += 1
		}
	}
	return result
}

// 获取切片中的大于等于指定值的数量
func GetSMinVCountI64(src []int64, minValue int64) int64 {
	if len(src) == 0 {
		return 0
	}
	var result int64 = 0
	for _, v := range src {
		if v >= minValue {
			result += 1
		}
	}
	return result
}

// 获取切片中等于指定值的数量
func GetSVCountI32(src []int32, tar int32) int32 {
	if len(src) == 0 {
		return 0
	}
	var result int32 = 0
	for _, v := range src {
		if v == tar {
			result += 1
		}
	}
	return result
}

// 切片是否包含元素
func ContainSVI32(src []int32, tar int32) bool {
	for _, v := range src {
		if v == tar {
			return true
		}
	}
	return false
}

// 切片是否包含元素
func ContainSVI64(src []int64, tar int64) bool {
	for _, v := range src {
		if v == tar {
			return true
		}
	}
	return false
}

// 切片是否包含元素
func ContainSVStr(src []string, tar string) bool {
	for _, v := range src {
		if v == tar {
			return true
		}
	}
	return false
}

// 切片是否包含切片, src是否包含tar
func ContainSSI32(src []int32, tar []int32) bool {
	for _, v := range tar {
		if !ContainSVI32(src, v) {
			return false
		}
	}
	return true
}

// 获取切片中等于指定值的数量
func GetSVCountI64(src []int64, tar int64) int32 {
	if len(src) == 0 {
		return 0
	}
	var result int32 = 0
	for _, v := range src {
		if v == tar {
			result += 1
		}
	}
	return result
}

// 切片是否包含切片, src是否包含tar
func ContainSSI64(src []int64, tar []int64) bool {
	for _, v := range tar {
		if !ContainSVI64(src, v) {
			return false
		}
	}
	return true
}

// 切片是是否包含重复的元素, excludeElement=排除的元素
func IsSliceRepeatElementI64(src []int64, excludeElements ... int64) bool {
	var exists bool
	for i, v1 := range src {
		exists = false
		for _, v2 := range excludeElements {
			if v1 == v2 {
				exists = true
				break
			}
		}
		if exists {
			continue
		}

		for j, v2 := range src {
			if j != i && v1 == v2 {
				return true
			}
		}
	}
	return false
}

// 切片是是否包含重复的元素, excludeElement=排除的元素
func IsSliceRepeatElementI32(src []int32, excludeElements ... int32) bool {
	var exists bool
	for i, v1 := range src {
		exists = false
		for _, v2 := range excludeElements {
			if v1 == v2 {
				exists = true
				break
			}
		}
		if exists {
			continue
		}

		for j, v2 := range src {
			if j != i && v1 == v2 {
				return true
			}
		}
	}
	return false
}

// 冒泡排序
func BubbleSort(slice interface{}, compare func(a interface{}, b interface{}) bool) []interface{} {
	val := reflect.ValueOf(slice)

	sliceLen := val.Len()
	out := make([]interface{}, sliceLen)
	for i := 0; i < sliceLen; i++ {
		out[i] = val.Index(i).Interface()
	}

	for i := 0; i < len(out)-1; i ++ {
		for j := 0; j < (len(out) - 1 - i); j++ {
			if compare((out)[j], (out)[j+1]) {
				out[j], out[j+1] = out[j+1], out[j]
			}
		}
	}

	return out
}

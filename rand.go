package goutils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 随机下标
func RandGroup(p ...uint32) int {
	if p == nil {
		panic("args not found")
	}

	r := make([]uint32, len(p))
	for i := 0; i < len(p); i++ {
		if i == 0 {
			r[0] = p[0]
		} else {
			r[i] = r[i-1] + p[i]
		}
	}

	rl := r[len(r)-1]
	if rl == 0 {
		return 0
	}

	rn := uint32(rand.Int63n(int64(rl)))
	for i := 0; i < len(r); i++ {
		if rn < r[i] {
			return i
		}
	}

	panic("bug")
}

// 范围内随机(包含最大值和最小值)
func RandInterval(b1, b2 int32) int32 {
	if b1 == b2 {
		return b1
	}

	min, max := int64(b1), int64(b2)
	if min > max {
		min, max = max, min
	}
	return int32(rand.Int63n(max-min+1) + min)
}

// 范围内随机指定长度的切片(包含最大值和最小值)
func RandIntervalN(b1, b2 int32, n uint32) []int32 {
	if b1 == b2 {
		return []int32{b1}
	}

	min, max := int64(b1), int64(b2)
	if min > max {
		min, max = max, min
	}
	l := max - min + 1
	if int64(n) > l {
		n = uint32(l)
	}

	r := make([]int32, n)
	m := make(map[int32]int32)
	for i := uint32(0); i < n; i++ {
		v := int32(rand.Int63n(l) + min)

		if mv, ok := m[v]; ok {
			r[i] = mv
		} else {
			r[i] = v
		}

		lv := int32(l - 1 + min)
		if v != lv {
			if mv, ok := m[lv]; ok {
				m[v] = mv
			} else {
				m[v] = lv
			}
		}

		l--
	}

	return r
}

// 生成指定范围随机数,包含min,不包含max, 且min>0
func RandInt(min, max int) int{
	if min >= max {
		return max
	}
	if min == 0 || max == 0 {
		return 0
	}
	return rand.Intn(max-min) + min
}

// 生成指定范围随机数,包含min,不包含max, 且min>0
func RandInt32(min, max int32) int32 {
	if min >= max {
		return max
	}
	if min == 0 || max == 0 {
		return 0
	}
	return rand.Int31n(max-min) + min
}

// 生成指定范围随机数,包含min,不包含max, 且min>0
func RandInt64(min, max int64) int64 {
	if min >= max {
		return max
	}
	if min == 0 || max == 0 {
		return 0
	}
	return rand.Int63n(max-min) + min
}

// 随机打乱切片顺序
func SliceOutOfOrder(in []interface{}) []interface{} {
	l := len(in)
	for i := l - 1; i > 0; i-- {
		r := rand.Intn(i)
		in[r], in[i] = in[i], in[r]
	}
	return in
}

// 获取随机字符串(大小写字母), n=字符串长度
func GetRandomString(n int) string {
	const alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
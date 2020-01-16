package gutils

import (
	"fmt"
	"strings"
	"strconv"
	"math"
)

// 10进制转16进制
func DecHex(n int64) string {
	if n < 0 {
		return ""
	}
	if n == 0 {
		return "0"
	}
	hex := map[int64]int64{10: 65, 11: 66, 12: 67, 13: 68, 14: 69, 15: 70}
	s := ""
	for q := n; q > 0; q = q / 16 {
		m := q % 16
		if m > 9 && m < 16 {
			m = hex[m]
			s = fmt.Sprintf("%v%v", string(m), s)
			continue
		}
		s = fmt.Sprintf("%v%v", m, s)
	}
	return s
}

// 16进制转10进制
func HexDec(h string) (n int64) {
	s := strings.Split(strings.ToUpper(h), "")
	l := len(s)
	i := 0
	d := float64(0)
	hex := map[string]string{"A": "10", "B": "11", "C": "12", "D": "13", "E": "14", "F": "15"}
	for i = 0; i < l; i++ {
		c := s[i]
		if v, ok := hex[c]; ok {
			c = v
		}
		f, err := strconv.ParseFloat(c, 10)
		if err != nil {
			return -1
		}
		d += f * math.Pow(16, float64(l-i-1))
	}
	return int64(d)
}

// 取大值Int8
func MaxI8(x int8, y int8) int8 {
	if x > y {
		return x
	}
	return y
}

// 取大值Uint8
func MaxUi8(x uint8, y uint8) uint8 {
	if x > y {
		return x
	}
	return y
}

// 取大值Int
func MaxI(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

// 取大值Uint
func MaxUi(x uint, y uint) uint {
	if x > y {
		return x
	}
	return y
}

// 取大值Int32
func MaxI32(x int32, y int32) int32 {
	if x > y {
		return x
	}
	return y
}

// 取大值Uint32
func MaxUi32(x uint32, y uint32) uint32 {
	if x > y {
		return x
	}
	return y
}

// 取大值Int64
func MaxI64(x int64, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

// 取大值Uint64
func MaxUi64(x uint64, y uint64) uint64 {
	if x > y {
		return x
	}
	return y
}

// 取小值Int8
func MinI8(x int8, y int8) int8 {
	if x < y {
		return x
	}
	return y
}

// 取小值Uint8
func MinUi8(x uint8, y uint8) uint8 {
	if x < y {
		return x
	}
	return y
}

// 取小值Int
func MinI(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

// 取小值Uint
func MinUi(x uint, y uint) uint {
	if x < y {
		return x
	}
	return y
}

// 取小值Int32
func MinI32(x int32, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

// 取小值Uint32
func MinUi32(x uint32, y uint32) uint32 {
	if x < y {
		return x
	}
	return y
}

// 取小值Int64
func MinI64(x int64, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

// 取小值Uint64
func MinUi64(x uint64, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}

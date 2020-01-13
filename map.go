package goutils

// map添加map
func AppendMMI32(src map[int32]int32, appendV map[int32]int32) {
	for k, v := range appendV {
		if _, ok := src[k]; ok {
			src[k] += v
		} else {
			src[k] = v
		}
	}
}

// map添加值
func AppendMVI32(m map[int32]int32, k int32, v int32) {
	if k > 0 {
		if _, ok := m[k]; ok {
			m[k] += v
		} else {
			m[k] = v
		}
	}
}

// map添加值
func AppendMVI32I64(m map[int32]int64, k int32, v int64) {
	if k > 0 && v > 0 {
		if _, ok := m[k]; ok {
			m[k] += v
		} else {
			m[k] = v
		}
	}
}

// map添加map
func AppendMMF32(src map[int32]float32, appendV map[int32]float32) {
	for k, v := range appendV {
		if _, ok := src[k]; ok {
			src[k] += v
		} else {
			src[k] = v
		}
	}
}

// map添加值
func AppendMVI32F32(m map[int32]float32, k int32, v float32) {
	if k > 0 {
		if _, ok := m[k]; ok {
			m[k] += v
		} else {
			m[k] = v
		}
	}
}

// 属性Map获取值
func GetMVI32F32(m map[int32]float32, k int32) float32 {
	if k > 0 {
		if v, ok := m[k]; ok {
			return v
		}
	}
	return 0
}

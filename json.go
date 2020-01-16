package gutils

import (
	"io/ioutil"
	"encoding/json"
)

// 转换json对象
func ParseJsonWithFile(relativePath string, p interface{}) error {
	path := AbsPath(relativePath)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, p)
}
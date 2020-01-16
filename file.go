package gutils

import (
	"os"
	"bufio"
	"strings"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"math"
	"io"
	"errors"
	"path/filepath"
)

// 指定的文件或目录是否存在
func FileExists(p string) bool {
	p = ReplacePathSplit(p)
	if _, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// 创建目录
func Mkdir(p string) error {
	p = ReplacePathSplit(p)
	if !FileExists(p) {
		return os.MkdirAll(p, os.ModePerm)
	}
	return nil
}

// 创建文件
func CreateFile(filepath string) (*os.File, error) {
	filepath = ReplacePathSplit(filepath)
	var f *os.File
	var err error
	if _, err = os.Stat(filepath); os.IsNotExist(err) {
		f, err = os.Create(filepath)
		if err != nil {
			return nil, err
		}
	} else {
		f, err = os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0666)
	}
	return f, err
}

// 删除文件
func DeleteFile(p string) error {
	p = ReplacePathSplit(p)
	if FileExists(p) {
		return os.RemoveAll(p)
	}
	return nil
}

// 逐行读取文件
func ReadLine(fileName string, handler func(string) error) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		line = strings.TrimSpace(line)
		if err = handler(line); err != nil {
			return err
		}
	}
	return nil
}

// 从json文件中读取对象, path=文件路径, p=结构体引用
func ReadByJsonFile(path string, p interface{}) error {
	// 加载配置
	absPath := AbsPath(path)
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return err;
	}
	err = json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	return nil
}

// 获取文件大小后缀
func GetStrFileSize(size int64) string {
	v := float64(size)
	if v < math.Pow(float64(1024), float64(1)) {
		return fmt.Sprintf("%v Byte", v)
	} else if v < math.Pow(float64(1024), float64(2)) {
		return fmt.Sprintf("%.2f KB", v/math.Pow(float64(1024), float64(1)))
	} else if v < math.Pow(float64(1024), float64(3)) {
		return fmt.Sprintf("%.2f MB", v/math.Pow(float64(1024), float64(2)))
	} else if v < math.Pow(float64(1024), float64(4)) {
		return fmt.Sprintf("%.2f GB", v/math.Pow(float64(1024), float64(3)))
	} else if v < math.Pow(float64(1024), float64(5)) {
		return fmt.Sprintf("%.2f TB", v/math.Pow(float64(1024), float64(4)))
	} else if v < math.Pow(float64(1024), float64(6)) {
		return fmt.Sprintf("%.2f PB", v/math.Pow(float64(1024), float64(5)))
	} else if v < math.Pow(float64(1024), float64(7)) {
		return fmt.Sprintf("%.2f EB", v/math.Pow(float64(1024), float64(6)))
	} else if v < math.Pow(float64(1024), float64(8)) {
		return fmt.Sprintf("%.2f ZB", v/math.Pow(float64(1024), float64(7)))
	} else if v < math.Pow(float64(1024), float64(9)) {
		return fmt.Sprintf("%.2f YB", v/math.Pow(float64(1024), float64(8)))
	} else {
		return fmt.Sprintf("%v", v)
	}
}

// 拷贝文件夹
func CopyDir(srcPath string, destPath string) error {
	if srcInfo, err := os.Stat(srcPath); err != nil {
		return err
	} else {
		if !srcInfo.IsDir() {
			return errors.New("srcPath not dir")
		}
	}
	if destInfo, err := os.Stat(destPath); err != nil {
		return err
	} else {
		if !destInfo.IsDir() {
			return errors.New("destInfo not dir")
		}
	}

	err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if !f.IsDir() {
			path := strings.Replace(path, "\\", "/", -1)
			destNewPath := strings.Replace(path, srcPath, destPath, -1)
			if _, err := CopyFile(path, destNewPath); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// 生成目录并拷贝文件
func CopyFile(src, dest string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	destSplitPathDirs := strings.Split(dest, "/")

	destSplitPath := ""
	for index, dir := range destSplitPathDirs {
		if index < len(destSplitPathDirs)-1 {
			destSplitPath = destSplitPath + dir + "/"
			b := FileExists(destSplitPath)
			if b == false {
				err := os.Mkdir(destSplitPath, os.ModePerm)
				if err != nil {
					return 0, err
				}
			}
		}
	}
	dstFile, err := os.Create(dest)
	if err != nil {
		return 0, err
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}
package goutils

import (
	"path/filepath"
	"os"
	"strings"
	"path"
	"io"
	"fmt"
	"crypto/md5"
)

// 获取当前运行程序根目录绝对路径
func CurrentPath() string {
	if p, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		p = ReplacePathSplit(p)
		return p
	}
	return ""
}

// 获取绝对路径, p:路径
func AbsPath(p string) string {
	p = ReplacePathSplit(p)
	if !strings.Contains(p, ":") && !path.IsAbs(p) {
		p = path.Join(CurrentPath(), p)
		return p
	}
	return p
}

// 替换路径分隔符
func ReplacePathSplit(p string) string {
	return strings.Replace(p, "\\", "/", len(p))
}

// 获取不带后缀的文件名
func GetFilenameWithoutSuffix(p string) string {
	filenameWithSuffix := path.Base(p)
	fileSuffix := path.Ext(filenameWithSuffix)
	return strings.TrimSuffix(filenameWithSuffix, fileSuffix)
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

// 获取父目录
func GetParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

// 获取相对路径
func GetRelativePath(baseAbsPath, absPath string) string {
	baseAbsPath = strings.Replace(baseAbsPath, "\\", "/", len(baseAbsPath))
	absPath = strings.Replace(absPath, "\\", "/", len(absPath))
	relativePath := strings.Trim(strings.Replace(absPath, baseAbsPath, "", len(absPath)), "/")

	return relativePath
}

// 获取上层目录
func GetParentDir(dir string) string {
	return SubString(dir, 0, strings.LastIndex(dir, "/"))
}

// 获取文件Md5
func GetFileMd5(path string) (strMd5 string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		return "", err
	}
	strMd5 = fmt.Sprintf("%x", md5hash.Sum(nil))

	return strMd5, nil
}

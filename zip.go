package gutils

import (
	"bytes"
	"compress/zlib"
	"io"
	"archive/zip"
	"os"
	"path"
	"crypto/md5"
	"fmt"
	"path/filepath"
	"runtime"
	"errors"
)

// 进行zlib压缩
func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

// 进行zlib解压缩
func DoZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

//压缩文件 files 文件数组, 可以是不同dir下的文件或者文件夹,dest 压缩文件存放地址
func ZipCompress(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err := compress(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// 解压文件
func ZipDeCompress(zipFile, dest string, mode os.FileMode, delete bool) (err error) {
	f, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer func() {
		if e := f.Close(); e == nil {
			if delete {
				DeleteFile(zipFile)
			}
		}
	}()

	// 优先把所有文件夹创建出来
	for _, v := range f.File {
		info := v.FileInfo()
		if info.IsDir() {
			err = Mkdir(path.Join(dest, v.Name))
			if err != nil {
				return errors.New(fmt.Sprintf("%v, %v", path.Join(dest, v.Name), err))
			}
		}
	}

	for _, v := range f.File {
		// 跳过文件夹
		info := v.FileInfo()
		if info.IsDir() {
			continue
		}
		newFile, err := os.Create(path.Join(dest, v.Name))
		if err != nil {
			return errors.New(fmt.Sprintf("%v, %v", path.Join(dest, v.Name), err))
		}

		if runtime.GOOS == "linux" {
			if err := newFile.Chmod(mode); err != nil {
				newFile.Close()
				return errors.New(fmt.Sprintf("%v, chmod fail, %v", path.Join(dest, v.Name), err))
			}
		}

		srcFile, err := v.Open()
		if err != nil {
			newFile.Close()
			return errors.New(fmt.Sprintf("%v, %v", path.Join(dest, v.Name), err))
			return err
		}

		_, err = io.Copy(newFile, srcFile)
		if err != nil {
			newFile.Close()
			srcFile.Close()
			return errors.New(fmt.Sprintf("%v, copy fail, %v", path.Join(dest, v.Name), err))
		}

		newFile.Close()
		srcFile.Close()
	}

	return err
}

// 生成差异zip包
func CreateDiffZipPackage(zipFile1, zipFile2, outputDir string) (err error) {
	// 解压包1
	f1, err := zip.OpenReader(zipFile1)
	if err != nil {
		return err
	}
	defer f1.Close()

	// 解压包2
	f2, err := zip.OpenReader(zipFile2)
	if err != nil {
		return err
	}
	defer f2.Close()

	// 生成差异文件存储目录
	err = Mkdir(outputDir)
	if err != nil {
		return err
	}

	for _, v2 := range f2.File {
		if v2.FileInfo().IsDir() {
			continue
		}
		exists := false
		for _, v1 := range f1.File {
			if v1.Name == v2.Name {
				exists = true
				v1F, err := v1.Open()
				if err != nil {
					return err
				}
				defer v1F.Close()
				v1Md5 := md5.New()
				_, err = io.Copy(v1Md5, v1F)
				if err != nil {
					return err
				}
				strV1Md5 := fmt.Sprintf("%x", v1Md5.Sum([]byte("")))

				v2F, err := v2.Open()
				if err != nil {
					return err
				}
				defer v2F.Close()
				v2Md5 := md5.New()
				_, err = io.Copy(v2Md5, v2F)
				if err != nil {
					return err
				}
				strV2Md5 := fmt.Sprintf("%x", v2Md5.Sum([]byte("")))

				//Logs.Debug("file=%v,md5=%v,  file=%v,md5=%v", v1.Name, strV1Md5, v2.Name, strV2Md5)

				exists = strV1Md5 == strV2Md5
				break
			}
		}

		if !exists {
			v2F, err := v2.Open()
			if err != nil {
				return err
			}
			defer v2F.Close()
			err = Mkdir(path.Dir(path.Join(outputDir, v2.Name)))
			if err != nil {
				return err
			}
			nf, err := os.Create(path.Join(outputDir, v2.Name))
			if err != nil {
				return err
			}
			defer nf.Close()
			_, err = io.Copy(nf, v2F)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// 压缩文件夹(分包压缩), dirPath:文件夹路径, splitPackLimit:分包大小上限字节, zipCollection:是否将zip分包整合成1个zip包, delete被压缩的目录是否删除
func ZipSplitCompressByDir(dirPath string, splitPackLimit int64, zipCollection, delete bool) (err error) {
	needZipFiles := [][]string{}
	zipPackageSize := int64(0)
	zipPackageSizeLimit := splitPackLimit
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if len(needZipFiles) == 0 || zipPackageSize+info.Size() >= zipPackageSizeLimit {
			needZipFiles = append(needZipFiles, []string{path})
			zipPackageSize = info.Size()
		} else {
			zipPackageSize += info.Size()
			needZipFiles[len(needZipFiles)-1] = append(needZipFiles[len(needZipFiles)-1], path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	zipFiles := make([]string, len(needZipFiles))
	for i, v := range needZipFiles {
		zipFiles[i] = fmt.Sprintf("%v_%v.zip", dirPath, i+1)
		err := ZipCompressByFilesPath(dirPath, v, zipFiles[i])
		if err != nil {
			return err
		}
	}

	if zipCollection {
		outputPath := fmt.Sprintf("%v.zip", dirPath)
		err = ZipCompressByFilesPath(GetParentDirectory(dirPath), zipFiles, outputPath)
		if err != nil {
			return err
		}
		for _, v := range zipFiles {
			err := DeleteFile(v)
			if err != nil {
				return err
			}
		}
	}

	if delete {
		err := DeleteFile(dirPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func ZipCompressByFilesPath(zipBaseDir string, zipFilesPath []string, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	for _, v := range zipFilesPath {
		f, err := os.Open(v)
		if err != nil {
			return err
		}
		defer f.Close()

		info, err := f.Stat()
		if err != nil {
			return err
		}
		if info.IsDir() {
			continue
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Method = zip.Deflate
		header.SetModTime(info.ModTime())
		header.Name = GetRelativePath(zipBaseDir, v)
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, f)
		if err != nil {
			return err
		}
	}
	return nil
}

// 压缩字节组
func ZipCompressByBytes(data []byte) ([]byte, error) {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	if _, err := w.Write(data); err != nil {
		w.Close()
		return nil, err
	}
	w.Close()

	return in.Bytes(), nil
}

// 解压缩字节组
func ZipDeCompressByBytes(zipData []byte) ([]byte, error) {
	b := bytes.NewReader(zipData)
	var out bytes.Buffer
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(&out, r); err != nil {
		r.Close()
		return nil, err
	}
	r.Close()
	return out.Bytes(), nil
}

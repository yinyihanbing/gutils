package gutils

import (
	"os"
	"fmt"
	"reflect"
	"encoding/csv"
	"errors"
)

// CSV帮助类
var CsvHelper = new(csvHelper)

type csvHelper struct{}

// 导出数据到CSV
func (this *csvHelper) Export(savePath string, columnName []string, fieldsName []string, dataSlice interface{}) error {
	if FileExists(savePath) {
		return errors.New(fmt.Sprintf("existing file %v", savePath))
	}

	f, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// 写入UTF-8 BOM, 防止中文乱码
	f.WriteString("\xEF\xBB\xBF")

	// 写入列名
	w := csv.NewWriter(f)
	w.Write(columnName)

	// 写入数据
	v := reflect.ValueOf(dataSlice)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	n := v.Len()
	rowData := make([]string, 0, len(fieldsName))
	for i := 0; i < n; i++ {
		for _, fName := range fieldsName {
			e := v.Index(i)
			if e.Kind() == reflect.Ptr {
				e = e.Elem()
			}
			item := e.FieldByName(fName)
			rowData = append(rowData, fmt.Sprintf("%v", item))
		}
		w.Write(rowData)
		rowData = make([]string, 0, len(fieldsName))
	}

	w.Flush()

	return nil
}

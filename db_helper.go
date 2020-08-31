package gutils

import (
	"fmt"
	"strings"
	"reflect"
	"errors"
	"encoding/json"
	"database/sql"
)

// 数据库帮助类结构体
type DbHelper struct {
	DbDriverName string
	DbConn       string
	DbName       string
	DB           *sql.DB
}

// 实例化数据库帮助类
func NewDbHelper(dbConn string) *DbHelper {
	dbHelper := new(DbHelper)
	dbHelper.DbDriverName = "mysql"
	dbHelper.DbConn = dbConn
	dbHelper.setDbName()

	return dbHelper
}

// 获取连接串中的数据库名
func (this *DbHelper) setDbName() {
	switch this.DbDriverName {
	case "mysql":
		this.DbName = this.DbConn[strings.LastIndex(this.DbConn, "/")+1 : strings.LastIndex(this.DbConn, "?")]
	}
}

// 打开数据库连接
func (this *DbHelper) Open() error {
	db, err := sql.Open(this.DbDriverName, this.DbConn)
	if err != nil {
		return err
	}
	this.DB = db
	return nil
}

// 关闭数据库连接
func (this *DbHelper) Close() {
	this.DB.Close()
}

// 获取slice、map、ptr里结构体最终类型
func (this *DbHelper) getStructType(p interface{}) reflect.Type {
	reflectType := reflect.ValueOf(p).Type()
	for reflectType.Kind() == reflect.Slice || reflectType.Kind() == reflect.Map || reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

// 获取结构体的数据库字段值容器
func (this *DbHelper) getStructDbValueContainer(i interface{}) []interface{} {
	vType := this.getStructType(i)
	count := vType.NumField()

	container := make([]interface{}, count)

	for i := 0; i < count; i++ {
		f := vType.Field(i)
		switch f.Type.Kind() {
		case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct, reflect.Ptr:
			var strText string
			container[i] = &strText
		default:
			container[i] = reflect.New(f.Type).Interface()
			break
		}
	}

	return container
}

// 结构体的数据库字段值容器 转结构体, container=容器, p=结构体应用
func (this *DbHelper) transformDbRowDataToStruct(container []interface{}, p interface{}) (err error) {
	rv := reflect.ValueOf(p)
	if rv.Type().Kind() != reflect.Ptr {
		return errors.New("parameter `p` must be a ptr")
	}

	rv = rv.Elem()

	count := rv.NumField()
	for i := 0; i < count; i++ {
		val := reflect.ValueOf(container[i]).Elem()
		f := rv.Field(i)
		if f.Type().Kind() == reflect.Slice && f.Elem().Kind() == reflect.Uint8 {
			f.SetBytes([]byte(val.String()))
		} else {
			switch f.Type().Kind() {
			case reflect.Ptr, reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
				jsonStr := val.String()
				if jsonStr == "" {
					switch f.Type().Kind() {
					case reflect.Ptr, reflect.Map:
						jsonStr = "{}"
					case reflect.Struct, reflect.Slice, reflect.Array:
						jsonStr = "[]"
					default:
						jsonStr = ""
					}
				}
				m := reflect.New(f.Type()).Interface()
				err = json.Unmarshal([]byte(jsonStr), m)
				if err != nil {
					err = fmt.Errorf("json unmarshal error: %v", err)
					break
				}
				f.Set(reflect.ValueOf(m).Elem())
				break
			default:
				f.Set(val)
				break
			}
		}
	}

	return err
}

// 判断表是否存在
func (this *DbHelper) IsTableExists(tableName string) (exists bool, err error) {
	strSql := fmt.Sprintf("SELECT COUNT(1) FROM INFORMATION_SCHEMA.tables WHERE table_name = '%v' AND table_schema = '%v'", tableName, this.DbName)
	var count int
	err = this.DB.QueryRow(strSql).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// 执行无返回结果集的Sql
func (this *DbHelper) ExecuteNonQuery(strSql string) (err error) {
	_, err = this.DB.Exec(strSql)
	return err
}

// 查询多行数据, p=结构体切片引用, strSql=执行的查询Sql语句,  args=Sql语句参数
func (this *DbHelper) ExecuteDataSlice(p interface{}, strSql string, args ...interface{}) error {
	pType := reflect.TypeOf(p)
	if pType.Kind() != reflect.Ptr {
		return errors.New("parameter `p` must be a ptr")
	}
	if pType.Elem().Kind() != reflect.Slice {
		return errors.New("parameter `p` must be a slice ptr")
	}
	if pType.Elem().Elem().Kind() != reflect.Struct && (pType.Elem().Elem().Kind() == reflect.Ptr && pType.Elem().Elem().Elem().Kind() != reflect.Struct) {
		return errors.New("parameter `p` must be a struct slice ptr")
	}

	// 执行查询
	rows, err := this.DB.Query(strSql, args...)
	if err != nil {
		return err
	}

	// 没有数据直接返回
	if rows == nil {
		return nil
	}

	// 获取结构体的容器
	vContainer := this.getStructDbValueContainer(p)

	pStructType := this.getStructType(p)
	pStructKind := reflect.TypeOf(p).Elem().Elem().Kind()

	results := reflect.ValueOf(p)
	if results.Kind() == reflect.Ptr {
		results = results.Elem()
	}
	for rows.Next() {
		err = rows.Scan(vContainer...)
		if err != nil {
			rows.Close()
			break
		}
		newItem := reflect.New(pStructType).Interface()
		err = this.transformDbRowDataToStruct(vContainer, newItem)
		if err != nil {
			rows.Close()
			break
		}
		if pStructKind == reflect.Ptr {
			results.Set(reflect.Append(results, reflect.ValueOf(newItem)))
		} else {
			results.Set(reflect.Append(results, reflect.ValueOf(newItem).Elem()))
		}
	}

	return err
}

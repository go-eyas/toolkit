package resource

import (
	"errors"
	"github.com/jinzhu/gorm"
	"reflect"
)

type Field struct {
	ColumnName   string // 数据库列名
	StructKey    string // 结构体的名字
	isPrimaryKey bool   // 是否主键
	isIgnore     bool   // 该字段是否忽略
	JsonKey      string // json key
	Search       string // 查询方式
	Order        string // 排序方式
	Update       bool   // 是否可更新，默认不可以
	Create       bool   // 是否可在创建时指定
}

type Resource struct {
	db        *gorm.DB
	tableName string
	model     *gorm.DB
	scope     *gorm.Scope
	pk        string
	sample    interface{}
	Fields    []*Field
	fieldsStructKeyMap map[string]*Field
}

func NewGormResource(db *gorm.DB, v interface{}) *Resource {
	scope := db.NewScope(v)
	r := &Resource{
		sample:    v,
		db:        db,
		tableName: scope.TableName(),
		model:     db.Table(scope.TableName()),
		scope:     scope,
	}
	fields, pk := r.parseFields(scope)
	r.Fields = fields
	r.pk = pk
	keyMap := map[string]*Field{}
	for _, field := range fields {
		keyMap[field.StructKey] = field
	}
	r.fieldsStructKeyMap = keyMap
	return r
}

func (r *Resource) Row(pk interface{}) *gorm.DB {
	return r.model.Where(r.pk+" = ?", pk)
}

func (r *Resource) Create(v interface{}) error {
	var model interface{}
	rv := reflect.ValueOf(v)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Struct:
		model = v
	case reflect.Map:
		model = r.mapToCreateStruct(v)

	default:
		return errors.New("create unsupport type")
	}
	return r.model.Create(model).Error
}
func (r *Resource) Update(pk interface{}, v interface{}) error {
	return r.Row(pk).Updates(v).Error
}

func (r *Resource) Detail(pk interface{}, v interface{}) error {
	return r.Row(pk).First(v).Error
}
func (r *Resource) List(slice interface{}, args ...interface{}) (int64, error) {
	switch len(args) {
	case 0:
		return r.listQuery(slice, nil, nil)
	case 1:
		return r.listQuery(slice, args[0], nil)
	case 2:
		return r.listQuery(slice, args[0], args[1])
	}
	return 0, errors.New("list param error")
}

func (r *Resource) Delete(pk interface{}) error {
	return r.model.Delete(r.sample, r.pk+" = ?", pk).Error
}
// func (r *Resource) CreateHandle(res http.ResponseWriter, req http.Request) {
//
// }
// func (r *Resource) UpdateHandle(res http.ResponseWriter, req http.Request) {}
// func (r *Resource) QueryHandle(res http.ResponseWriter, req http.Request)  {}
// func (r *Resource) DeleteHandle(res http.ResponseWriter, req http.Request) {}

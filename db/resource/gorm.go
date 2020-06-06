package resource

import (
	"github.com/jinzhu/gorm"
	"net/http"
	"strings"
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
}

func parseFields(scope *gorm.Scope) ([]*Field, string) {
	gormFields := scope.Fields()
	fields := make([]*Field, len(gormFields))

	pk := scope.PrimaryKey()

	for i, f := range gormFields {
		field := &Field{
			// 默认值
			StructKey: f.Name,
			ColumnName:   f.DBName,
			JsonKey:      f.Name,
			Search:       "=",
			Order:        "DESC",
			isPrimaryKey: pk == f.DBName,
		}
		fields[i] = field
		tags := f.StructField.Tag
		// json tag
		jsonTag := tags.Get("json")
		if jsonTag != "" {
			field.JsonKey = jsonTag
		} else if jsonTag == "-" {
			field.JsonKey = ""
		}

		// parse resource tag
		resourceTag := tags.Get("resource")

		// 为空，全为默认值
		if resourceTag == "" {
			continue
		}
		// 忽略
		if resourceTag == "-" {
			field.isIgnore = true
		}
		tagList := strings.Split(resourceTag, ";")

		for _, value := range tagList {
			p := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(p[0]))
			v := ""
			if len(p) >= 2 {
				v = strings.Join(p[1:], ":")
			}
			switch k {
			case "PK":
				pk = field.ColumnName
				field.isPrimaryKey = true
			case "SEARCH":
				field.Search = strings.ToUpper(v)
			case "ORDER":
				field.Order = strings.ToUpper(v)
			case "CREATE":
				field.Create = true
			case "UPDATE":
				field.Update = true
			}
		}

	}
	return fields, pk
}

func NewGormResource(db *gorm.DB, v interface{}) *Resource {
	scope := db.NewScope(v)
	fields, pk := parseFields(scope)
	return &Resource{
		sample:    v,
		db:        db,
		tableName: scope.TableName(),
		model:     db.Table(scope.TableName()),
		scope:     scope,
		pk:        pk,
		Fields:    fields,
	}
}

func (r *Resource) Row(pk interface{}) *gorm.DB {
	return r.model.Where(r.pk+" = ?", pk)
}

func (r *Resource) Create(v interface{}) error {
	return r.model.Create(v).Error
}
func (r *Resource) Update(pk interface{}, v interface{}) error {
	return r.Row(pk).Updates(v).Error
}

func (r *Resource) Detail(pk interface{}, v interface{}) error {
	return r.Row(pk).First(v).Error
}
func (r *Resource) List(v interface{}, slice interface{}) (int64, error) {
	return 0, nil
}
func (r *Resource) Delete(pk interface{}) error {
	return r.model.Delete(r.sample, r.pk+" = ?", pk).Error
}
func (r *Resource) CreateHandle(res http.ResponseWriter, req http.Request) {

}
func (r *Resource) UpdateHandle(res http.ResponseWriter, req http.Request) {}
func (r *Resource) QueryHandle(res http.ResponseWriter, req http.Request)  {}
func (r *Resource) DeleteHandle(res http.ResponseWriter, req http.Request) {}

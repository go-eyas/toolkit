package resource

import (
  "encoding/json"
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
  db                 *gorm.DB
  tableName          string
  model              *gorm.DB
  modelTypeName      string
  scope              *gorm.Scope
  pk                 string
  sample             interface{}
  Fields             []*Field
  fieldsStructKeyMap map[string]*Field
  defaultOrder       []string
}

// NewGormResource 实例化资源
//
// example:
//
//   type Article struct {
//     ID      int64  `resource:"pk;search:=;order:desc" json:"id"`
//     Title   string `resource:"create;update;search:like" json:"title"`
//     Content string `resource:"create;update;search:like" json:"text"`
//     Status  byte   `resource:"search:=" json:"-"`
//   }
//
//   r := NewGormResource(db, &Article{})
func NewGormResource(db *gorm.DB, v interface{}) *Resource {
  scope := db.NewScope(v)
  rv := reflect.ValueOf(v)
  if rv.Kind() == reflect.Ptr {
    rv = rv.Elem()
  }
  rt := rv.Type()
  if rt.Kind() == reflect.Ptr {
    rt = rt.Elem()
  }

  r := &Resource{
    sample:    rv.Interface(),
    db:        db,
    tableName: scope.TableName(),
    model:     db.Table(scope.TableName()),
    scope:     scope,
  }
  r.modelTypeName = rt.PkgPath() + "." + rt.Name()
  fields, pk := r.parseFields(scope)
  r.Fields = fields
  r.pk = pk
  r.defaultOrder = []string{}
  keyMap := map[string]*Field{}
  for _, field := range fields {
    keyMap[field.StructKey] = field
    if field.Order == "DESC" || field.Order == "ASC" {
      r.defaultOrder = append(r.defaultOrder, field.ColumnName+" "+field.Order)
    }
  }
  r.fieldsStructKeyMap = keyMap
  return r
}

// Model 返回绑定了数据表的 gorm 实例
//   r.Model().Where("status = ?", status).Count()
func (r *Resource) Model() *gorm.DB {
  return r.model
}

// Row 返回绑定了主键值的 gorm 实例
//
//   r.Row().Update("status", 1)
func (r *Resource) Row(pk interface{}) *gorm.DB {
  return r.model.Where(r.pk+" = ?", pk)
}

// Create 创建资源，支持传入 struct、map，创建前会重置 resource tag 未设置 create 的字段为0值，使得创建记录时忽略值
//
//   err := r.Create(&Article{Title: "Hello", Status: 2}) 这里 status 会被重置为 0 ，因为 status 的 resource tag 未设置 create
func (r *Resource) Create(v interface{}) error {
  model, err := r.toCreateStruct(v, true)
  if err != nil {
    return err
  }
  return r.model.Create(model).Error
}

// CreateX 创建资源，支持传入 struct、map，传入的值均有效
//
//   err := r.CreateX(&Article{Title: "Hello", Status: 2}) 这里 status 成功设置值
func (r *Resource) CreateX(v interface{}) error {
  model, err := r.toCreateStruct(v, false)
  if err != nil {
    return err
  }
  return r.model.Create(model).Error
}

// Update 更新资源，支持传入 struct、map，只会更新 resource tag 设置了 update 的字段
//
//   err := r.Update(1, map[string]string{"title": "after title"})
func (r *Resource) Update(pk interface{}, v interface{}) error {
  updates, err := r.toUpdateMap(v, true)
  if err != nil {
    return err
  }
  if len(updates) > 0 {
    return r.Row(pk).Updates(updates).Error
  }
  return nil
}

// UpdateX 更新资源，支持传入 struct、map，更新传入的所有字段，如果传入的是 struct ，则会忽略 0 值
//
//   err := r.UpdateX(1, map[string]byte{"status": 1})
func (r *Resource) UpdateX(pk interface{}, v interface{}) error {
	updates, err := r.toUpdateMap(v, false)
	if err != nil {
		return err
	}
	if len(updates) > 0 {
		return r.Row(pk).Updates(updates).Error
	}
	return nil
}

// Detail 查询指定主键的记录
//
//   article := &Article{}
//   err := r.Detail(1, article)
func (r *Resource) Detail(pk interface{}, v interface{}) error {
  return r.Row(pk).First(v).Error
}

// List 查询资源列表，提供查询条件，排序规则，查询列表，查询规则会以resource tag 的search 值为准
//
//   list := []*Article{}
//   total, err := r.List(&list, map[string]byte{"status": 1})
//
func (r *Resource) List(slice interface{}, args ...interface{}) (int64, error) {
  switch len(args) {
  case 0:
    return r.listQuery(slice, nil, nil, nil)
  case 1:
    page := &Pagination{}
    query := args[0]
    raw, err := json.Marshal(query)
    if err != nil {
      return 0, errors.New("query parse error")
    }
    err = json.Unmarshal(raw, page)
    if err != nil {
      return 0, errors.New("query parse error")
    }
    order := r.getOrderArgs(query)
    return r.listQuery(slice, page, query, order)
  case 2:
    page := &Pagination{}
    query := args[0]
    raw, err := json.Marshal(query)
    if err != nil {
      return 0, errors.New("query parse error")
    }
    err = json.Unmarshal(raw, page)
    if err != nil {
      return 0, errors.New("query parse error")
    }
    return r.listQuery(slice, page, args[0], args[1])
  }
  return 0, errors.New("list param error")
}

func (r *Resource) ListPage(slice interface{}, page *Pagination, args ...interface{}) (int64, error) {
  switch len(args) {
  case 0:
    return r.listQuery(slice, nil, nil, nil)
  case 1:
    return r.listQuery(slice, page, args[0], nil)
  case 2:
    return r.listQuery(slice, page, args[0], args[1])
  }
  return 0, errors.New("list param error")
}

// Delete 删除指定主键的资源
//
//   err := r.Delete(1)
//
func (r *Resource) Delete(pk interface{}) error {
  return r.model.Delete(r.sample, r.pk+" = ?", pk).Error
}


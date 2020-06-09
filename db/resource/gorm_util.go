package resource

import (
  "fmt"
  "github.com/jinzhu/gorm"
  "reflect"
  "strings"
)

// 解析资源结构体
func (r *Resource) parseFields(scope *gorm.Scope) ([]*Field, string) {
  gormFields := scope.Fields()
  fields := make([]*Field, len(gormFields))

  pk := scope.PrimaryKey()

  for i, f := range gormFields {
    field := &Field{
      // 默认值
      StructKey:    f.Name,
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

// 新实例化资源
func (r *Resource) newModel() reflect.Value {
  dest := reflect.New(reflect.TypeOf(r.sample))
  if dest.Kind() == reflect.Ptr {
    return dest.Elem()
  }
  return dest
}

// 暂时用json先顶着用，以后再扩展 column 识别
func (r *Resource) mapToCreateStruct(mapValue interface{}) interface{} {
  dest := r.newModel()
  rv := reflect.ValueOf(mapValue)

  for _, field := range r.Fields {
    if !field.Create {
      continue
    }
    v := rv.MapIndex(reflect.ValueOf(field.StructKey))
    if !v.IsValid() {
      v = rv.MapIndex(reflect.ValueOf(field.ColumnName))
    }
    if !v.IsValid() {
      v = rv.MapIndex(reflect.ValueOf(field.JsonKey))
    }
    if !v.IsValid() || v.IsZero() {
      continue
    }
    v = reflect.ValueOf(v.Interface())
    destField := dest.FieldByName(field.StructKey)
    if destField.CanSet() && destField.Kind() == v.Kind() {
      destField.Set(v)
    }
  }

  return dest.Addr().Interface()
}

// 解析查询数据
func (r *Resource) parseQueryArgs(query interface{}) []*queryArg {
  if query == nil {
    return []*queryArg{}
  }
  qrv := reflect.ValueOf(query)
  if qrv.Kind() == reflect.Ptr {
    qrv = qrv.Elem()
  }
  qrt := qrv.Type()
  args := make([]*queryArg, 0)
  if qrv.Kind() == reflect.Struct {
    for i := 0; i < qrv.NumField(); i++ {
      fieldType := qrt.Field(i)
      fieldName := fieldType.Name
      resourceField, ok := r.fieldsStructKeyMap[fieldName]
      if !ok {
        continue
      }
      fvr := qrv.Field(i)
      if fvr.IsZero() {
        continue
      }
      fieldValue := fvr.Interface()

      if resourceField.Search == "LIKE" {
        fieldValue = fmt.Sprintf("%%%v%%", fieldValue)
        if fieldValue == "%%" {
          continue
        }
      }

      item := &queryArg{
        Column: resourceField.ColumnName,
        SearchType: resourceField.Search,
        SearchValue: fieldValue,
      }
      args = append(args, item)
    }
  } else if qrv.Kind() == reflect.Map {
    rv := qrv
    for _, field := range r.Fields {
      if field.Search == "NONE" || field.Search == "-"  {
        continue
      }

      v := rv.MapIndex(reflect.ValueOf(field.StructKey))
      if !v.IsValid() {
        v = rv.MapIndex(reflect.ValueOf(field.ColumnName))
      }
      if !v.IsValid() {
        v = rv.MapIndex(reflect.ValueOf(field.JsonKey))
      }
      if !v.IsValid() || v.IsZero() {
        continue
      }
      fieldValue := v.Interface()
      if field.Search == "LIKE" {
        fieldValue = fmt.Sprintf("%%%v%%", fieldValue)
      }
      item := &queryArg{
        Column:      field.ColumnName,
        SearchValue: fieldValue,
        SearchType:  field.Search,
      }
      args = append(args, item)
    }
  }
  return args
}

// 解析排序数据
func (r *Resource) parseOrderArgs(order interface{}) []string {
  if order == nil {
    return []string{}
  }
  if s, ok := order.([]string); ok {
    return s
  }

  args := []string{}
  ov := reflect.ValueOf(order)
  if ov.Kind() == reflect.Map {
    for _, key := range ov.MapKeys() {
      k := key.Interface()
      v := strings.ToUpper(fmt.Sprintf("%v", ov.MapIndex(key).Interface()))
      if v != "DESC" && v != "ASC" {
        continue
      }

      arg := fmt.Sprintf("%v %s", k, v)
      args = append(args, arg)
    }
  }
  return args
}

type Pagination struct {
  Offset int
  Limit  int
}

type queryArg struct {
  Column      string // 列名
  SearchValue interface{} // 查询值
  SearchType  string // 查询类型
}

func (r *Resource) listQuery(slice interface{}, page *Pagination, query interface{}, order interface{}) (int64, error) {
  // raw, err := json.Marshal(query)
  // if err != nil {
  //   return 0, errors.New("query parse error")
  // }
  // err = json.Unmarshal(raw, page)
  // if err != nil {
  //   return 0, errors.New("query parse error")
  // }

  var err error
  queryArgs := r.parseQueryArgs(query)
  q := r.model
  for _, arg := range queryArgs {
    if arg == nil {
      continue
    }

    q = q.Where(fmt.Sprintf("%s %s ?", arg.Column, arg.SearchType), arg.SearchValue)
  }

  orderArgs := r.parseOrderArgs(order)
  for _, arg := range orderArgs {
    q = q.Order(arg)
  }

  var count int64
  err = q.Count(&count).Error
  if err != nil {
    return 0, err
  }

  if page != nil {
    if page.Offset != 0 {
      q = q.Offset(page.Offset)
    }
    if page.Limit != 0 {
      q = q.Limit(page.Limit)
    }
  }

  err = q.Find(slice).Error

  return count, err




}

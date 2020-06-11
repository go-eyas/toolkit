package resource

import (
  "errors"
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
      Order:        "-",
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


func (r *Resource) toCreateStruct(v interface{}, protectField bool) (interface{}, error) {
  rv := reflect.ValueOf(v)

  if rv.Kind() == reflect.Ptr {
    rv = rv.Elem()
  }

  if rv.Kind() != reflect.Map && rv.Kind() != reflect.Struct {
    return nil, errors.New("create only support struct and map")
  }

  var dest reflect.Value
  rt := rv.Type()
  isOriginType := (rt.PkgPath() + "." + rt.Name()) == r.modelTypeName
  if isOriginType {
    dest = rv
  } else {
    dest = r.newModel()
  }


  for _, field := range r.Fields {
    if protectField && !field.Create {
      if isOriginType {
        destField := dest.FieldByName(field.StructKey)
        destField.Set(reflect.Zero(rv.FieldByName(field.StructKey).Type()))
      }
      continue
    }
    var v reflect.Value
    if rv.Kind() == reflect.Map {
      v = rv.MapIndex(reflect.ValueOf(field.StructKey))
      if !v.IsValid() {
        v = rv.MapIndex(reflect.ValueOf(field.ColumnName))
      }
      if !v.IsValid() {
        v = rv.MapIndex(reflect.ValueOf(field.JsonKey))
      }
      if !v.IsValid() || v.IsZero() {
        continue
      }
    } else if rv.Kind() == reflect.Struct {
      v = rv.FieldByName(field.StructKey)
      if !v.IsValid() || v.IsZero() {
        continue
      }
    }

    if v.IsValid() {
      v = reflect.ValueOf(v.Interface())
      destField := dest.FieldByName(field.StructKey)
      if destField.CanSet() && destField.Kind() == v.Kind() {
        destField.Set(v)
      }
    }
  }

  return dest.Addr().Interface(), nil
}

func (r *Resource) toUpdateMap(v interface{}, protectField bool) (map[string]interface{}, error) {
  result := map[string]interface{}{}
  rv := reflect.ValueOf(v)

  if rv.Kind() == reflect.Ptr {
    rv = rv.Elem()
  }

  if rv.Kind() != reflect.Map && rv.Kind() != reflect.Struct {
    return result, errors.New("update only support struct and map")
  }

  for _, field := range r.Fields {
    if protectField && !field.Update {
      continue
    }
    var v reflect.Value
    if rv.Kind() == reflect.Map {
      v = rv.MapIndex(reflect.ValueOf(field.StructKey))
      if !v.IsValid() {
        v = rv.MapIndex(reflect.ValueOf(field.ColumnName))
      }
      if !v.IsValid() {
        v = rv.MapIndex(reflect.ValueOf(field.JsonKey))
      }
      if !v.IsValid() || v.IsZero() {
        continue
      }
    } else if rv.Kind() == reflect.Struct {
      v = rv.FieldByName(field.StructKey)
      if !v.IsValid() || v.IsZero() {
        continue
      }
    }

    if v.IsValid() {
      v = reflect.ValueOf(v.Interface())
      result[field.ColumnName] = v.Interface()
    }
  }
  return result, nil
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
  baseOrder := append([]string{}, r.defaultOrder...)
  if order == nil {
    return baseOrder
  }

  args := baseOrder

  if s, ok := order.([]string); ok {
    args = append(args, s...)
  } else if s, ok := order.(string); ok {
    args = append(args, strings.Split(s, ",")...)
  } else {
    ov := reflect.ValueOf(order)
    if ov.Kind() == reflect.Map {
      for _, key := range ov.MapKeys() {
        k := key.Interface()
        v := strings.TrimSpace(strings.ToUpper(fmt.Sprintf("%v", ov.MapIndex(key).Interface())))
        if v != "DESC" && v != "ASC" {
          continue
        }

        arg := fmt.Sprintf("%v %s", k, v)
        args = append(args, arg)
      }
    }
  }


  // 去除重复
  orderMap := map[string]string{}
  for _, x := range args {
    for _, o := range strings.Split(x, ",") {
      p := strings.Split(strings.TrimSpace(o), " ")
      if len(p) < 2 {
        continue
      }
      orderMap[p[0]] = strings.TrimSpace(strings.Join(p[1:], " "))
    }
  }
  orders := []string{}
  for _, x := range args {
    for _, o := range strings.Split(x, ",") {
      p := strings.Split(strings.TrimSpace(o), " ")
      if len(p) < 2 {
        continue
      }
      k := p[0]
      val, ok := orderMap[k]
      if ok {
        orders = append(orders, k + " " + val)
      }
      delete(orderMap, k)
    }

  }

  return orders
}

func (r *Resource) getOrderArgs(v interface{}) []string {
  rv := reflect.ValueOf(v)
  var orderField reflect.Value
  if rv.Kind() == reflect.Ptr {
    rv = rv.Elem()
  }
  if rv.Kind() == reflect.Struct {
    orderField = rv.FieldByName("Order")

  } else if rv.Kind() == reflect.Map {
    orderField = rv.MapIndex(reflect.ValueOf("order"))
  }
  if orderField.IsValid() {
    return r.parseOrderArgs(orderField.Interface())
  }
  return []string{}
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

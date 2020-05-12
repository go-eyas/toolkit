# 事件监听器

最简洁的事件监听与分发

## API

默认使用全局监听器，在全局范围的监听均有效，如果需要局部的事件分发，可使用 New 重新实例化一个，实例化的事件分发与全局的完全隔离

支持链式操作

#### New()

重新实例化一个事件分发器

#### On(name string, handler func(interface{}))

增加监听

 * name 事件名称
 * handler 事件触发回调函数

#### Off(name string, handler func(interface{}))

取消事件监听，如果 handler 为空，则取消该事件的所有监听，如果不为空，则只取消指定的监听函数

 * name 事件名称
 * handler 事件触发回调函数

```go
emit.Off("evt") // 取消 evt 时间的所有函数监听
emit.Off("evt", fn1) // 只取消 fn1 函数的监听

```

#### Emit(name string, data interface{})

触发事件

 * name 事件名称
 * data 事件触发携带的数据

## 使用

```go
import "github.com/go-eyas/toolkit/emit"

fn1 := func(data interface{}) {
  fmt.Printf("fn1 receive data: %v", data)
}

fn2 := func(data interface{}) {
  fmt.Printf("fn2 receive data: %v", data)
}

fn3 := func(data interface{}) {
  fmt.Printf("fn3 receive data: %v", data)
}

emit.
  On("evt", fn1).
  On("evt", fn2, fn3).
  Emit("evt", "hello emitter")

emit.Off("evt", fn3)
emit.Emit("evt", "hello emitter again")

// or
e := emit.New().On(...).Off(...)

e.Emit(...)
```

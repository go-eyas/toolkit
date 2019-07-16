# GRPC

## 编译 proto

将会遍历执行目录与其子目录的所有后缀为 `.proto` 的文件编译为 go 源码，使用的编译命令是 

```sh
protoc -I=$GOPATH/src -I=. -I=$(pwd) --proto_path=$(pwd) --go_out=$(pwd) $(pwd)
```
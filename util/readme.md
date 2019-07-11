# 工具函数

```
package util // import "github.com/go-eyas/toolkit/util"


FUNCTIONS

func Assert(err error, msg interface{})
    Assert 断言 err != nil

func AssignMap(maps ...map[string]interface{}) map[string]interface{}
    AssignMap 合并多个map

func Base64Decoding(enc string) (string, error)
    Base64Decoding base64 解码

func Base64Encoding(str string) string
    Base64Encoding base64 编码

func FuncName(f interface{}) string
    FuncName 获取函数的名字

func HasFile(f string) bool
    HasFile 是否存在该文件

func RandomStr(length int) string
    RandomStr 生成随机字符串

func StructToMap(v interface{}) map[string]interface{}
    StructToMap 把结构体转成map，key使用json定义的key

func ToString(v interface{}) string
    ToString 把能转成字符串的都转成JSON字符串

func ToStruct(raw interface{}, v interface{})
    ToStruct 把一个结构体转成另一个结构体，以json key作为关联

```
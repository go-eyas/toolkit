package email

import (
	"os"
	"testing"

	"github.com/BurntSushi/toml"
)

var tomlConfig = `
host = "smtp.qq.com"
port = "465"
account = "893521870@qq.com"
password = "haha, wo cai bu gao su ni ne"
name = "unit test"
secure = true
[tpl.a]
bcc = ["Jeason <eyasliu@163.com>"] # 抄送
cc = [] # 抄送人
subject = "Welcome, {{.Name}}" # 主题
text = "Hello, I am {{.Name}}" # 文本
html = "<h1>Hello, I am {{.Name}}</h1>" # html 内容
`

func TestEmail(t *testing.T) {
	conf := &Config{}
	toml.Decode(tomlConfig, conf)
	conf.Password = os.Getenv("password")

	email := New(conf)
	err := email.SendByTpl("Yuesong Liu <liuyuesongde@163.com>", "a", struct{ Name string }{"Batman"})

	if err != nil {
		t.Errorf("mail send fail: %v", err)
	} else {
		t.Log("mail send success")
	}
}

func ExampleSample() {
	tomlConfig := `
host = "smtp.qq.com"
port = "465"
account = "893521870@qq.com"
password = "haha, wo cai bu gao su ni ne"
name = "unit test"
secure = true
[tpl.a]
bcc = ["Jeason <eyasliu@163.com>"] # 抄送
cc = [] # 抄送人
subject = "Welcome, {{.Name}}" # 主题
text = "Hello, I am {{.Name}}" # 文本
html = "<h1>Hello, I am {{.Name}}</h1>" # html 内容
`
	conf := &Config{}
	toml.Decode(tomlConfig, conf)
	email := New(conf)
	email.SendByTpl("Yuesong Liu <liuyuesongde@163.com>", "a", struct{ Name string }{"Batman"})
}

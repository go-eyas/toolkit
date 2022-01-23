# Fork From [lorca](https://github.com/zserge/lorca)

lorca 封装了非常实用的 chrome devtool protocol，但是暴露的 api 不足以使用更多的 cdp 功能，再次封装

# Lorca

构建桌面应用，纯 Go 语言实现，支持 chrome 最新浏览器

## 浏览器

默认会使用本地已安装的 chrome 内核浏览器，会自动判断是否已安装以下浏览器，按照以下判断顺序

- Google Chrome
- Microsoft Edge
- Opera
- 360 极速浏览器
- 360 安全浏览器
- QQ 浏览器
- 百度浏览器

如果已安装则直接使用，如果全都未安装，则自动下载一个绿色版 Chrome 启动程序

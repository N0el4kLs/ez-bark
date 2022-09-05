<h1 align="center">
    ezbark
</h1>



## Introduce

[`什么是BARK`](https://github.com/Finb/Bark)

`EzBark` 是一个简化 `bark api` 调用的工具,它支持命令行调用,消息文件调用以及第三方包调用.

`EzBark`的出现是为了更方便的进行消息推送.你可以将其融入到自己的脚本中,或是生成可执行文件配合`shell`脚本进行使用.

## Download

你可以使用以下命令安装二进制文件,也可以下载源码自己编译.

```
go install github.com/N0el4kLs/ez-bark/cmd/ezbark@latest
```

## Usage

1. 支持使用命令行输入,简单、快速完成消息推送
2. 文件输入.如果传入的内容固定，可以使用 `-initf` 选项生成消息缓存文件.将相应的内容写在文件后,使用 `-f`选项指定加载消息文件.
3. 如果你想推送的消息需要从网络上获取，可以将此仓库作为第三方包进行导入.



```
./ezbark -init http://127.0.0.1,abdefg
./ezbark -showconf=true
./ezbark -test=true
./ezbark -t title -m message -g group -s bell
./ezbark -initf filename.yml
./ezbark -f bark.yml
```

## Ezbark Go library
Usage example:

```go
package main

import ezbark "github.com/N0el4kLs/ez-bark"

func main() {
// use map[string]interface{} to store your own data
send := make(map[string]interface{})
options := ezbark.NewOptions()

// Use DefaultConf to load default config in ~/.config/ezbark/config.yml
// or you can use SetSend options to load bark server or devices
options.DefaultConf()

//fmt.Printf("%#v\n", options)

send["title"] = "title"
send["body"] = "title"
//send["server"] = "http://127.0.0.1:7070"
//send["key"] = "abcdefg"

options = options.SetSend(send)
options.Notice()
}
```

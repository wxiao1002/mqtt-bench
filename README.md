# mqtt-bench-go
MQTT Broker 测试工具

## 如何使用
1. 打开目录下测试文件，test.csv 编辑你的用户名密码与客户端id,
2. 修改 .vscode/launch.json 文件
```
             "-broker","tcp://192.168.1.100:1883",
             "-topic",  "/v1/12",
             "-payload", "{\"id\":1}",
             "-count", "1000",
             "-csv",  "./test.csv"
```
3. 使用 vscode 运行

////////////////////////////////
## 如何运行2
1. csv 文件，文件内容如下示范
```
ClientId,Username,Password
1,1,1
2,2,2
```
2. 运行 main.go
```
go run main.go -broker tcp://192.168.1.100:1883 -topic  /v1/12 -payload {"id":1}  -count 10 -csv  ./test.csv
```

## Installation
```
go install github.com/wxiao1002/mqtt-bench@main
```  

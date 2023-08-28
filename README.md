# MQTT Bench
 MQTT 发布消息测试工具

 ## 快速开始
 - 准备mqtt 账号密码csv 文件，参照 device_secret.csv
 - 运行
```
go run main.go
```
### 参数解析
broker ：mqtt broker 地址 <br/>    		
csvPath：读取用户密码地址，第一列用户，第二列密码 <br/> 
clients：创建多少个连接，csv 读取到的数目大于等于<br/> 
benchmarkTime: 压测时间，分钟开始<br/> 
messageIntervalInSec：生产消息的时间间隔，秒<br/> 

## 下一步计划
- 完成统计


		        
		

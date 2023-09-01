# MQTT Bench
 MQTT 发布消息测试工具

 ## 快速开始

 ### 编译好的程序包
-  通过 https://github.com/wxiao1002/mqtt-bench/releases 下载对应版本的程序包
- 运行 并传入参数,以下是linux 示范
```
./mqtt-bench-linux -csvPath /root/clients.csv -broker tcp://localhost:1883 -clients 1000 -benchmarkTime 20 -messageInterval 1
```
- 上述示范是 读取/root/clients.csv  文件解析用户名密码，生成1000个客户端并连接到 tcp://localhost:1883 ，以1秒一条的消息频率发布消息测试20 分钟

### 源码运行
- git clone https://github.com/wxiao1002/mqtt-bench.git
- 运行代码

```
go run main.go -csvPath /root/clients.csv -broker tcp://localhost:1883 -clients 1000 -benchmarkTime 20 -messageInterval 1
```
### 参数解析
- broker ：mqtt broker 地址    		
- csvPath：读取用户密码地址，第一列用户，第二列密码 
- clients：创建多少个连接，csv 读取到的数目大于等于该值
- benchmarkTime: 压测时间，分钟开始
- messageInterval：生产消息的时间间隔，秒
- topic：发布主题，为空的话 会是 api/{username}/attributes,username 是占位符到时候会自动替换成用户名 

# 运行结果

[![pPdN091.png](https://s1.ax1x.com/2023/08/29/pPdN091.png)](https://imgse.com/i/pPdN091)


		        
# 打包
- sh build/linux.sh	
- sh build/linux_arm.sh 
- sh build/windows.sh

# english doc
[./README-en.md]

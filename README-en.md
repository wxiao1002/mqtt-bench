# MQTT Bench
  MQTT publish message test tool

  ## Quick start

  ### Compiled package
- Download the corresponding version of the package through https://github.com/wxiao1002/mqtt-bench/releases
- Run and pass in parameters, the following is a linux demonstration
```
./mqtt-bench-linux -csvPath /root/clients.csv -broker tcp://localhost:1883 -clients 1000 -benchmarkTime 20 -messageInterval 1
```
- The above demonstration is to read the /root/clients.csv file to parse the username and password, generate 1000 clients and connect to tcp://localhost:1883, and publish a message at a frequency of 1 second for 20 minutes

### Source code running
- git clone https://github.com/wxiao1002/mqtt-bench.git
- run the code

```
go run main.go -csvPath /root/clients.csv -broker tcp://localhost:1883 -clients 1000 -benchmarkTime 20 -messageInterval 1
```
### Parameter analysis
- broker: mqtt broker address
- csvPath: read user password address, first column user, second column password
- clients: how many connections are created, the number read by csv is greater than or equal to this value
- benchmarkTime: pressure test time, starting from minute
- messageInterval: the time interval for producing messages, in seconds
- topic: publish topic, if it is empty, it will be api/{username}/attributes, username is a placeholder and will be automatically replaced with username

# operation result

[![pPdN091.png](https://s1.ax1x.com/2023/08/29/pPdN091.png)](https://imgse.com/i/pPdN091)
### 测试一

**模拟raft在跨数据中心的写操作**

目录在/MyRaft/mykvraft下面

启动多个3 5 7 个server


~~~go
go run myserver.go -address localhost:port members ip:port,ip:port  -delay 10
~~~

.......

-delay: 代表数据中心的延迟

调整client的mode 和 client nums

 ~~~go
go run myclient.go -servers localhost:port ... -mode write  -nums 10
 ~~~

-mode  read write 两种模式

-nums client数量
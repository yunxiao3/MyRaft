在三台server机器的mykvraft目录下分别运行下面命令

~~~shell
 go run myserver.go -address ip:5001 -members  membersIp1:5000,membersIp2:5002
~~~

在client机器的mykvraft目录下分别运行下面命令

~~~shell
go run myclient.go
~~~

~~~go
func request(num int)  {
	ck := Clerk{}
	ck.servers = make([]string, 3) 
    
    // 三个server ip需要配置
	ck.servers[0] = "ip:50011"
	ck.servers[1] = "ip:50001"
	ck.servers[1] = "ip:50021"


 	for i := 0; i < 1000 ; i++ {
		rand.Seed(time.Now().UnixNano())
		key := rand.Intn(100000)
		value := rand.Intn(100000)
		// 写操作
        ck.Put("key" + strconv.Itoa(key),"value"+ strconv.Itoa(value) )
		// 读操作
        fmt.Println(num, ck.Get("key1" ) )
		count++
	}
}


func main()  {
	
    clients := 10
	for i := 0; i < clients ; i++ {
		go request(i)		
	}
	time.Sleep(time.Second)
    // 一秒钟的读写数
	fmt.Println( count )
	time.Sleep(time.Second*1200)

}
~~~


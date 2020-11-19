package main

import (
	"log"
	"time"
	"fmt"
	//"strconv"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/reflection"
	crand "crypto/rand"
	"math/big"
	KV "../grpc/mykv"
	"strconv"
	"math/rand"
	//Per "../persister"

	//"math/rand"
)
type Clerk struct {
	servers []string
	// You will have to modify this struct.
	leaderId int
	id int64
	seq int64
}



func makeSeed() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := crand.Int(crand.Reader, max)
	x := bigx.Int64()
	return x
}


func MakeClerk(servers []string ) *Clerk {
	ck := new(Clerk)
	ck.servers = servers
	ck.id = makeSeed()
	ck.seq = 0
	// You'll have to add code here.
	return ck
}






func (ck *Clerk) Get(key string) string {
	args := &KV.GetArgs{Key:key}
	id := ck.leaderId
	for {
		reply, ok := ck.getValue(ck.servers[id], args)
		fmt.Println(reply.IsLeader)
		if (ok && reply.IsLeader){
			ck.leaderId = id;
			return reply.Value;
		}
		id = (id + 1) % len(ck.servers) 
	} 
}

func (ck *Clerk) getValue(address string , args  *KV.GetArgs) (*KV.GetReply, bool){
	// Initialize Client
	conn, err := grpc.Dial( address , grpc.WithInsecure(),grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect: %v", err)
	}
	defer conn.Close()
	client := KV.NewKVClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()
	reply, err := client.Get(ctx,args)
	if err != nil {

		log.Printf(" getValue could not greet: %v", err)
	}
	return reply, true
}



func (ck *Clerk) Put(key string, value string) bool {
	// You will have to modify this function.
	args := &KV.PutAppendArgs{Key:key,Value:value,Op:"Put", Id:ck.id, Seq:ck.seq }
	id := ck.leaderId
	for {
		reply, ok := ck.putAppendValue(ck.servers[id], args)
		if (ok && reply.IsLeader){
			ck.leaderId = id;
			return true
		}
		id = (id + 1) % len(ck.servers) 
	} 
}



func (ck *Clerk) Append(key string, value string) bool {
	// You will have to modify this function.
	args := &KV.PutAppendArgs{Key:key,Value:value,Op:"Append", Id:ck.id, Seq:ck.seq }
	id := ck.leaderId
	for {
		reply, ok := ck.putAppendValue(ck.servers[id], args)
		if (ok && reply.IsLeader){
			ck.leaderId = id;
			return true
		}
		id = (id + 1) % len(ck.servers) 
	} 
}


func (ck *Clerk) putAppendValue(address string , args  *KV.PutAppendArgs) (*KV.PutAppendReply, bool){
	// Initialize Client
	conn, err := grpc.Dial( address , grpc.WithInsecure(),grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect: %v", err)
	}
	defer conn.Close()
	client := KV.NewKVClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()
	reply, err := client.PutAppend(ctx,args)
	if err != nil {
		log.Printf("  putAppendValue could not greet: %v", err)
	}
	return reply, true
}




var count int

func request(num int)  {
	ck := Clerk{}
	ck.servers = make([]string, 3) 
	ck.servers[0] = "localhost:50011"
	ck.servers[1] = "localhost:50001"
	ck.servers[2] = "localhost:50021"


 	for i := 0; i < 10 ; i++ {
		rand.Seed(time.Now().UnixNano())
		key := "key" + strconv.Itoa(rand.Intn(100000))
		value := "value"+ strconv.Itoa(rand.Intn(100000))
		ok := ck.Put(key,value)
		fmt.Println(i ,"put   " ,value,ok )
		fmt.Println(i, "get   " ,ck.Get(key) )
		count++
	}
}


func main()  {
	fmt.Println( "count" )
	servers := 10
	//begin_time := time.Now().UnixNano()
	for i := 0; i < servers ; i++ {
		//go 
		request(i)		
	}
	request(1)
	//end_time := time.Now().UnixNano()

	//t := end_time - begin_time
	time.Sleep(time.Second)
	fmt.Println( count )

	time.Sleep(time.Second*1200)


}
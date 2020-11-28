package main

import (
	"log"
	"time"
	"fmt"
	"flag"
	"sync/atomic"

	//"strconv"
	"strings"
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
		fmt.Println(id)
		reply, ok := ck.getValue(ck.servers[id], args)
		fmt.Println(reply.IsLeader)
		if (ok && reply.IsLeader){
			ck.leaderId = id;
			return reply.Value;
		}else{
			fmt.Println("can not connect ", ck.servers[id], "or it's not leader")
		}
		id = (id + 1) % len(ck.servers) 
		
	} 
}

func (ck *Clerk) getValue(address string , args  *KV.GetArgs) (*KV.GetReply, bool){
	// Initialize Client
	conn, err := grpc.Dial( address , grpc.WithInsecure() ) //,grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect: %v", err)
		return nil, false
	}
	defer conn.Close()
	client := KV.NewKVClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 1)
	defer cancel()
	reply, err := client.Get(ctx,args)
	if err != nil {
		return  nil, false
		log.Printf(" getValue could not greet: %v", err)
	}
	return reply, true
}



func (ck *Clerk) Put(key string, value string) bool {
	// You will have to modify this function.
	args := &KV.PutAppendArgs{Key:key,Value:value,Op:"Put", Id:ck.id, Seq:ck.seq }
	id := ck.leaderId
	for {
		//fmt.Println(id)
		reply, ok := ck.putAppendValue(ck.servers[id], args)
		//fmt.Println(ok)

		if (ok && reply.IsLeader){
			ck.leaderId = id;
			return true;
		}else{
			fmt.Println(ok, "can not connect ", ck.servers[id], "or it's not leader")
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
	conn, err := grpc.Dial( address , grpc.WithInsecure() )//,grpc.WithBlock())
	if err != nil {
		return  nil, false
		log.Printf(" did not connect: %v", err)
	}
	defer conn.Close()
	client := KV.NewKVClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 1)
	defer cancel()
	reply, err := client.PutAppend(ctx,args)
	if err != nil {
		return  nil, false
		log.Printf("  putAppendValue could not greet: %v", err)
	}
	return reply, true
}




var count  int32  = 0

func Wirterequest(num int, servers []string)  {
	ck := Clerk{}
	ck.servers = make([]string, len(servers)) 

	for i:= 0; i <  len(servers); i++{
		ck.servers[i] = servers[i] + "1"
	}

 	for i := 0; i < num ; i++ {
		key := "key" + strconv.Itoa(rand.Intn(100000))
		value := "value"+ strconv.Itoa(rand.Intn(100000))
		ck.Put(key,value)
		atomic.AddInt32(&count,1)
	}
}


func Readequest(num int, servers []string)  {
	ck := Clerk{}
	ck.servers = make([]string, len(servers)) 

	for i:= 0; i <  len(servers); i++{
		ck.servers[i] = servers[i] + "1"
	}

 	for i := 0; i < num ; i++ {
		ck.Get("key") 
		atomic.AddInt32(&count,1)
	}
}




func main()  {

	var ser = flag.String("servers", "", "Input Your follower")
	flag.Parse()
	servers := strings.Split( *ser, ",")

	fmt.Println( "count" )
	serverNumm := 1
	for i := 0; i < serverNumm ; i++ {
		go  Wirterequest(10, servers)		
	} 
	time.Sleep(time.Second * 3)
	fmt.Println( count  / 3 )

	time.Sleep(time.Second*1200)


}
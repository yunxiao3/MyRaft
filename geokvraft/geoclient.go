package main

import (
	"log"
	"time"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/reflection"
	crand "crypto/rand"
	"math/big"
	KV "../grpc/mykv"
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


/* func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
} */


func MakeClerk(servers []string ) *Clerk {
	ck := new(Clerk)
	ck.servers = servers
	ck.id = makeSeed()
	ck.seq = 0
	// You'll have to add code here.
	return ck
}






func (ck *Clerk) Get(key string) string {
	
	//id := ck.leaderId
	args := &KV.GetArgs{Key:key}

	reply, ok := ck.getValue(ck.servers[0], args)
	if (ok){
		return reply.Value
	}
	return reply.Value


/* 
	//reply  := &KV.GetReply{}
	for {	
		//ok := ck.servers[id].Call("KVServer.Get",&args, &reply )	
		
		
		if (ok && reply.IsLeader){
			ck.leaderId = id;
			//fmt.Println("get value ",  reply.Value)
			return reply.Value;
		}
//		fmt.Println("Wrong leader !")

		id = (id + 1) % len(ck.servers) 
	} */
	
	//return reply.Value
}


func (ck *Clerk) getValue(address string , args  *KV.GetArgs) (*KV.GetReply, bool){

	fmt.Println("getValue")

	// Initialize Client
	conn, err := grpc.Dial( address , grpc.WithInsecure(),grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := KV.NewKVClient(conn)


	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//args := &RPC.AppendEntriesArgs{}
	reply, err := client.Get(ctx,args)
	if err != nil {
		log.Printf("could not greet: %v", err)
	}


	return reply, true
	//log.Printf("Append reply: %s", r)
	//fmt.Println("Append name is ABC")
}


/* 

func (ck *Clerk) PutAppend(key string, value string, op string) {
	// You will have to modify this function.
	args := KV.PutAppendArgs{key,value,op, ck.id,ck.seq}
	ck.seq++
	id := ck.leaderId
	for {
		reply := KV.PutAppendReply{}
		ok := ck.servers[id].Call("KVServer.PutAppend", &args, &reply)
		if (ok && reply.IsLeader){
			ck.leaderId = id;
		//	fmt.Println("putappend", key , value)
			return 
		}
		fmt.Print()
		id = (id + 1) % len(ck.servers)
	} 
}

func (ck *Clerk) Put(key string, value string) {
	ck.PutAppend(key, value, "Put")
}
func (ck *Clerk) Append(key string, value string) {
	ck.PutAppend(key, value, "Append")
} */


func main()  {
	ck := Clerk{}
	ck.servers = make([]string, 5) 
	ck.servers[0] = "localhost:4000"
	for {
		fmt.Println( ck.Get("get") )
	}

}
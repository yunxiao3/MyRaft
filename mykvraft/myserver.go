package main

import (
	
	"log"
	"net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"../myraft"
	"sync"
	"flag"
	"time"
	"strings"
	"fmt"
	Per "../persister"
	config "../config"
	KV "../grpc/mykv"
)




type KVServer struct {
	mu      sync.Mutex
	me      int
	rf      *myraft.Raft
	applyCh chan int
	//applyCh chan raft.ApplyMsg
	dead    int32 // set by Kill()
	maxraftstate int // snapshot if log grows this big

	// Your definitions here.
	Seq map[int64]int64
	db map[string]string
	chMap map[int]chan config.Op
	persist  *Per.Persister

}


func (kv *KVServer) PutAppend(ctx context.Context,args *KV.PutAppendArgs) ( *KV.PutAppendReply, error){
	// Your code here.


	//time.Sleep(time.Second)
	reply := &KV.PutAppendReply{};
	_ , reply.IsLeader = kv.rf.GetState()
	//reply.IsLeader = false;
	if !reply.IsLeader{
		return reply, nil
	}
	oringalOp := config.Op{args.Op, args.Key,args.Value, args.Id, args.Seq}
	index, _, isLeader := kv.rf.Start(oringalOp)
	if !isLeader {
		fmt.Println("Leader Changed !")
		reply.IsLeader = false
		return reply, nil
	}

	apply := <- kv.applyCh
	fmt.Println(apply)

	fmt.Println(index)

	return reply, nil

}



func (kv *KVServer) Get(ctx context.Context, args *KV.GetArgs) ( *KV.GetReply, error){
	
	reply := &KV.GetReply{}
	_ , reply.IsLeader = kv.rf.GetState()
	//reply.IsLeader = false;
	if !reply.IsLeader{
		return reply, nil
	}

	//fmt.Println()


	oringalOp := config.Op{"Get", args.Key,"" , 0, 0}
	_, _, isLeader  := kv.rf.Start(oringalOp)
	if !isLeader {
		return reply, nil
	}
	//fmt.Println(index)

	reply.IsLeader = true
	//kv.mu.Lock()
	//fmt.Println("Asdsada")
	reply.Value = string( kv.persist.Get(args.Key) )
	//fmt.Println("Asdsada")

	//kv.mu.Unlock()
	return reply, nil
}


func (kv *KVServer) equal(a config.Op, b config.Op) bool  {
	return (a.Option == b.Option && a.Key == b.Key &&a.Value == b.Value) 
}



func (kv *KVServer) RegisterServer(address string)  {
	// Register Server 
	for{

		lis, err := net.Listen("tcp", address)
		fmt.Println(address)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		KV.RegisterKVServer(s, kv )
		// Register reflection service on gRPC server.
		reflection.Register(s)
		if err := s.Serve(lis); err != nil {
			fmt.Println("failed to serve: %v", err)
		}
		
	}
	
}



func main()  {

	var add = flag.String("address", "", "Input Your address")
	var mems = flag.String("members", "", "Input Your follower")
	flag.Parse()

	server := KVServer{}

	// Local address	
	address := *add

	persist := &Per.Persister{}
	persist.Init("../db/"+address)
	
	//for i := 0; i <= int (address[ len(address) - 1] - '0'); i++{
	server.applyCh = make(chan int, 1)
	//}
	
	//server.applyCh = make(chan int, 1)
	fmt.Println(server.applyCh)

	server.persist  = persist

	// Members's address
	members := strings.Split( *mems, ",")


	fmt.Println(address, members)



	go server.RegisterServer(address+"1")
	server.rf = myraft.MakeRaft(address , members ,persist, &server.mu,server.applyCh)
	

	time.Sleep(time.Second*1200)

}

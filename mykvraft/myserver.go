package main

import (
	
	"log"
	"net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"os"
	"../myraft"
	"sync"
	"time"
	"fmt"
	KV "../grpc/mykv"
)

type Op struct {
	// Your definitions here.
	// Field names must start with capital letters,
	// otherwise RPC will break.
	Option string
	Key string 
	Value string
	Id int64
	Seq int64
}



type KVServer struct {
	mu      sync.Mutex
	me      int
	rf      *myraft.Raft
	//applyCh chan raft.ApplyMsg
	dead    int32 // set by Kill()
	maxraftstate int // snapshot if log grows this big

	// Your definitions here.
	Seq map[int64]int64
	db map[string]string
	chMap map[int]chan Op
}


func (kv *KVServer) PutAppend(ctx context.Context,args *KV.PutAppendArgs) ( *KV.PutAppendReply, error){
	// Your code here.
	reply := &KV.PutAppendReply{};
	_ , isLeader := kv.rf.GetState()
	reply.IsLeader = false;
	if !isLeader{
		return reply, nil
	}
	oringalOp := Op{args.Op, args.Key,args.Value, args.Id, args.Seq}
	index, _, isLeader := kv.rf.Start(oringalOp)
	if !isLeader {
		return reply, nil
	}
	fmt.Println(index)

	return reply, nil

	/* ch := kv.getChan(index)
	op := kv.replyNotify(ch)
	
	if kv.equal(oringalOp , op){
		reply.IsLeader = true
		return
	}  */
}



func (kv *KVServer) Get(ctx context.Context, args *KV.GetArgs) ( *KV.GetReply, error){
	reply := &KV.GetReply{}
	_ , isLeader := kv.rf.GetState()
	reply.IsLeader = false;
	if !isLeader{
		return reply, nil
	} 
	oringalOp := Op{"Get", args.Key,"" , 0, 0}
	index, _, isLeader  := kv.rf.Start(oringalOp)
	if !isLeader {
		return reply, nil
	}
	fmt.Println(index)
	return reply, nil
	/*  
	ch := kv.getChan(index)
	op := kv.replyNotify(ch)
	if args.Key == op.Key{
		reply.IsLeader = true
		kv.mu.Lock()
		reply.Value = kv.db[args.Key]
		kv.mu.Unlock()
		return
	} 

	*/
}


func (kv *KVServer) equal(a Op, b Op) bool  {
	return (a.Option == b.Option && a.Key == b.Key &&a.Value == b.Value) 
}



func (kv *KVServer) RegisterServer(address string)  {
	// Register Server 
	for{

		lis, err := net.Listen("tcp", address)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		KV.RegisterKVServer(s, kv /* &Raft{} */)
		// Register reflection service on gRPC server.
		reflection.Register(s)
		if err := s.Serve(lis); err != nil {
			fmt.Println("failed to serve: %v", err)
		}
		
	}
	
}






func main()  {
	server := KVServer{}
	go server.RegisterServer("localhost:40000")
	fmt.Println(os.Args)
	server.rf = myraft.MakeRaft(os.Args)

	time.Sleep(time.Second*1200)

}

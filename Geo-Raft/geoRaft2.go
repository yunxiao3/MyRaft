package main

import (
	"log"
	//"net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"sync"
	//"google.golang.org/grpc/reflection"
	RPC "./grpc/raft"
	"fmt"
	"time"
)

type State int
const (
    Follower State = iota // value --> 0
    Candidate             // value --> 1
    Leader                // value --> 2
)

type Log struct {
    Term    int         "term when entry was received by leader"
    Command interface{} "command for state machine,"
}

type ApplyMsg struct {
    CommandValid bool
    Command      interface{}
    CommandIndex int
}

type Raft struct {
    mu        sync.Mutex          // Lock to protect shared access to this peer's state
   // peers     []*labrpc.ClientEnd // RPC end points of all peers
   // persister *Persister          // Object to hold this peer's persisted state
    me        int32                 // this peer's index into peers[]
	port string
	ip string
    // state a Raft server must maintain.
    state     State

    //Persistent state on all servers:(Updated on stable storage before responding to RPCs)
    currentTerm int32    "latest term server has seen (initialized to 0 increases monotonically)"
    votedFor    int32    "candidateId that received vote in current term (or null if none)"
    log         []Log  "log entries;(first index is 1)"

    //Volatile state on all servers:
    commitIndex int32    "index of highest log entry known to be committed (initialized to 0, increases monotonically)"
    lastApplied int32    "index of highest log entry applied to state machine (initialized to 0, increases monotonically)"

    //Volatile state on leaders：(Reinitialized after election)
    nextIndex   []int32  "for each server,index of the next log entry to send to that server"
    matchIndex  []int32  "for each server,index of highest log entry known to be replicated on server(initialized to 0, im)"

    //channel
    applyCh     chan ApplyMsg // from Make()
    killCh      chan bool //for Kill()
    //handle rpc
    voteCh      chan bool
	appendLogCh chan bool
	

	// New 

	client RPC.RAFTClient


}

func send(ch chan bool) {
    select {
    case <-ch: //if already set, consume it then resent to avoid block
    default:
    }
    ch <- true
}

// SayHello implements helloworld.GreeterServer
func (raft *Raft) RequestVote(ctx context.Context, in *RPC.RequestVoteArgs) (*RPC.RequestVoteReply, error) {
	/* raft.msg = raft.msg + in.Name

	if (in.Name == "ABC"){
		fmt.Println("RequestVoteArgs name is ABC")
		
		return &pb.RequestVoteReply{Message: "Hello ABC " + in.Name}, nil
	}else{
		fmt.Println("RequestVoteArgs name is raft ")
		return &pb.RequestVoteReply{Message: "Hello RAFT" + raft.msg}, nil
	} */
	return &RPC.RequestVoteReply{Term:1}, nil
}
 
func (raft *Raft)AppendEntries(ctx context.Context, in *RPC.AppendEntriesArgs) (*RPC.AppendEntriesReply, error) {
	//raft.msg = raft.msg + in.Term
	reply := &RPC.AppendEntriesReply{}
	reply.Term = 1
    reply.Success = false
    reply.ConflictTerm = -1
    reply.ConflictIndex = 0
	if (in.Term == 1){
		fmt.Println("Append name is ABC")
		return reply, nil
	}else{
		return reply, nil
	}
}


func (rf *Raft)init(port string, ip string){


	//rf := &Raft{}

    rf.state = Follower
    rf.currentTerm = 0
    rf.votedFor = -1
    rf.log = make([]Log,1) //(first index is 1)

    rf.commitIndex = 0
    rf.lastApplied = 0

    //rf.applyCh = applyCh
    //because gorountne only send the chan to below goroutine,to avoid block, need 1 buffer
    rf.voteCh = make(chan bool,1)
    rf.appendLogCh = make(chan bool,1)
    rf.killCh = make(chan bool,1)

	rf.port = port
	rf.ip = ip
	//rf.msg = ""
	



	// Register Server 
/* 	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RPC.RegisterRAFTServer(s, &Raft{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	} */


	// Initialize Client
	conn, err := grpc.Dial( "localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
//	defer conn.Close()
	rf.client = RPC.NewRAFTClient(conn)
}

func (raft *Raft) StartRequestVote(){

	fmt.Println("StartRequestVote")


	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	args := &RPC.RequestVoteArgs{Term:	1,CandidateId: 1, LastLogIndex: 1,LastLogTerm :1};


	r, err := raft.client.RequestVote(ctx,args)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r)
	fmt.Println("Append name is ABC")
}



func (raft *Raft) StartAppendEntries(){

	fmt.Println("StartAppendEntries")


	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	args := &RPC.AppendEntriesArgs{}


	r, err := raft.client.AppendEntries(ctx,args)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r)
	fmt.Println("Append name is ABC")
}



func main() {
	raft := Raft{}
	raft.init(":50052","localhost")

	raft.StartAppendEntries()
}
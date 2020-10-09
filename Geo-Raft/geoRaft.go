package main

import (
	"log"
	"net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"sync"
	"google.golang.org/grpc/reflection"
	RPC "./grpc/raft"
	"fmt"
	"time"
	"./labgob"
	"bytes"
)

type State int
const (
    Follower State = iota // value --> 0
    Candidate             // value --> 1
    Leader                // value --> 2
)

type Log struct {
    Term    int32         "term when entry was received by leader"
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

const NULL int32 = -1



//Helper function
func send(ch chan bool) {
    select {
    case <-ch: //if already set, consume it then resent to avoid block
    default:
    }
    ch <- true
}

func (rf *Raft) getPrevLogIdx(i int) int32 {
    return rf.nextIndex[i] - 1
}

func (rf *Raft) getPrevLogTerm(i int) int32 {
    prevLogIdx := rf.getPrevLogIdx(i)
    if prevLogIdx < 0 {
        return -1
    }
    return rf.log[prevLogIdx].Term
}

func (rf *Raft) getLastLogIdx() int32 {
    return int32(len(rf.log) - 1)
}

func (rf *Raft) getLastLogTerm() int32 {
    idx := rf.getLastLogIdx()
    if idx < 0 {
        return -1
    }
    return rf.log[idx].Term
}

// SayHello implements helloworld.GreeterServer
func (rf *Raft) RequestVote(ctx context.Context, in *RPC.RequestVoteArgs) (*RPC.RequestVoteReply, error) {
	
	fmt.Println("RequestVote CALL")
	reply := &RPC.RequestVoteReply{}
	rf.currentTerm++
	reply.Term = rf.currentTerm
    reply.VoteGranted = false
	
	return reply, nil
}
 
func (rf *Raft)AppendEntries(ctx context.Context, in *RPC.AppendEntriesArgs) (*RPC.AppendEntriesReply, error) {
	//raft.msg = raft.msg + in.Term
	
	fmt.Println("AppendEntries CALL")

	reply := &RPC.AppendEntriesReply{}
	rf.currentTerm ++
	reply.Term = rf.currentTerm 
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


func  RegisterServer(port string)  {
	// Register Server 
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RPC.RegisterRAFTServer(s, &Raft{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}


//Follower Section:
func (rf *Raft) beFollower(term int32) {
    rf.state = Follower
    rf.votedFor = NULL
    rf.currentTerm = term
    //rf.persist()
}


func (rf *Raft) beLeader() {
    if rf.state != Candidate {
        return
	}
	//fmt.Println(rf.me, " Become Leader !", rf.currentTerm)
    rf.state = Leader
    //initialize leader data
    rf.nextIndex = make([]int32,/* len(rf.peers) */10)
    rf.matchIndex = make([]int32,/* len(rf.peers) */10)//initialized to 0
    for i := 0; i < len(rf.nextIndex); i++ {//(initialized to leader last log index + 1)
        rf.nextIndex[i] = rf.getLastLogIdx() + 1
    }
}


//Candidate Section:
// If AppendEntries RPC received from new leader: convert to follower implemented in AppendEntries RPC Handler
func (rf *Raft) beCandidate() { //Reset election timer are finished in caller
	//fmt.Println(rf.me,"become Candidate", rf.currentTerm)
	rf.state = Candidate
    rf.currentTerm++ //Increment currentTerm
    rf.votedFor = rf.me //vote myself first
    //rf.persist()
    //ask for other's vote
    go rf.startElection() //Send RequestVote RPCs to all other servers
}



func (rf *Raft) startAppendLog() {
	idx := 0

	

	appendLog := append(make([]Log,0),rf.log[rf.nextIndex[idx]:]...)

	w := new(bytes.Buffer)
    e := labgob.NewEncoder(w)
    e.Encode(appendLog)

    data := w.Bytes()

	args := RPC.AppendEntriesArgs{
		Term: rf.currentTerm,
		LeaderId: rf.me,
		PrevLogIndex: rf.getPrevLogIdx(idx),
		PrevLogTerm: rf.getPrevLogTerm(idx),
		//If last log index ≥ nextIndex for a follower:send AppendEntries RPC with log entries starting at nextIndex
		//nextIndex > last log index, rf.log[rf.nextIndex[idx]:] will be empty then like a heartbeat
		Log: data,
		LeaderCommit: rf.commitIndex,
	}
	rf.StartAppendEntries(&args)
}


//If election timeout elapses: start new election handled in caller
func (rf *Raft) startElection() {

	args := RPC.RequestVoteArgs{
        Term: rf.currentTerm,
        CandidateId: rf.me,
        LastLogIndex: rf.getLastLogIdx(),
        LastLogTerm: rf.getLastLogTerm(),

	};

	// TODO.....
	rf.StartRequestVote(&args)
   
}

func (rf *Raft)init(port string, ip string){


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


	// Add New 
	rf.port = port
	rf.ip = ip
	

	go RegisterServer(rf.port)

	

	// Initialize Client
	conn, err := grpc.Dial( "localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//	defer conn.Close()
	rf.client = RPC.NewRAFTClient(conn)
}


func (raft *Raft) StartAppendEntries(args  *RPC.AppendEntriesArgs){

	fmt.Println("StartAppendEntries")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//args := &RPC.AppendEntriesArgs{}
	r, err := raft.client.AppendEntries(ctx,args)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r)
	fmt.Println("Append name is ABC")
}


func (raft *Raft) StartRequestVote(args *RPC.RequestVoteArgs){

	fmt.Println("StartRequestVote")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := raft.client.RequestVote(ctx, args)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Term)

}



func main() {
	raft := Raft{}
	raft.init(":50051","localhost")
	for i:= 0; i < 10; i++{
		raft.startElection()
		//raft.StartAppendEntries()
	}
}
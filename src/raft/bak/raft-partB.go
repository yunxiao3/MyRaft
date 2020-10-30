package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//


import (
	"sync"
	"sync/atomic"	
	"../labrpc"
	"math/rand"
	"time"
	"fmt"
	"sort"
)




const (
	Follower int = 1;
	Candidate int = 2;
	Leader int = 3;
)

type LogEntry struct{
	Term int
	Command interface{}
}

//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]
	dead      int32               // set by Kill()


	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

	/* Follower == 1 || Candidate  == 2 || Leader == 3  */
	state int

	// Volatile state on all servers:
	currentTerm int
	commitIndex int 
	lastAplied int
	isVoted bool 
	timeout int
	log[] LogEntry

	// Volatile state on leaders:
	nextIndex []int
	matchIndex []int

	// Debug information to give every sever a number 
	number int

	// Using chan to implement the communication between goroutine 
	resetChan chan bool
	applyCh chan ApplyMsg
}

// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in Lab 3 you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh; at that point you can add fields to
// ApplyMsg, but set CommandValid to false for these other uses.
//
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
}

//
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {

	// Your data here (2A, 2B).
	Term int
	CandidateId int
	LastLogIndex int 
	LastLogTerm int 
	
}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (2A).
	Term int
	VoteGranted bool 
	Message string
	
}

type AppendEntry struct{
	Term int
	LeaderId int
	PrevLogIndex int
	PrevLogTerm int 
	Entries[] LogEntry
	LeaderCommit int
}

type AppendEntryReply struct{
	Term int
	Success bool 
}


// Debug information
var N = 0



/* =======  Conmunication ======= */



//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (2A, 2B).
	rf.resetChan <- true
	rf.mu.Lock()
	defer rf.mu.Unlock() 
	// Reply message 
	reply.Term = rf.currentTerm
	reply.VoteGranted = false
	// Receiving a bigger term, server should immediately call becomeFollower 
	// to uptdate the value of trem and isVoted 
	if args.Term > rf.currentTerm {
		rf.becomeFollower(args.Term)
	} 

	//fmt.Println(args.CandidateId," request Vote ", args)


	// If candidate's term is less than server's term server or has already voted,
	// the server should immediately refuse candidate's voting request
	if args.Term < rf.currentTerm || rf.isVoted == true {
		//fmt.Println("args.Term < rf.currentTe")
		return
	}
	// If candidate's LastLogTerm is less than server's
	// the server should immediately refuse candidate's voting request
	if args.LastLogTerm < rf.getLastLogTerm(){
		//fmt.Println("rgs.LastLogTerm < rf.currentT")
		return 
	}
	// If candidate's LastLogTerm equal to server's and its LastLogIndex is less than server's
	// the server should immediately refuse candidate's voting request
	if args.LastLogTerm == rf.getLastLogTerm() && args.LastLogIndex < rf.getLastLogIndex(){
		//fmt.Println("astLogTerm == rf.getL && args.LastLogIndex < rf.getLastLogIndex(")
		return
	}
	reply.VoteGranted = true
}


func sendChan(ch chan bool){
	select{
		case <- ch:
		default:
	}
	ch <- true
}
//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}


func (rf *Raft) startElection(){

	// Counting the number of acquired votes 
	var votes int32 = 1
	rf.isVoted = true
	
	for i:=0; i < len(rf.peers); i++{

		if i == rf.me {
			continue
		}
		
		go func (i int)  {
			args := &RequestVoteArgs{
				rf.currentTerm,
				rf.me,
				rf.getLastLogIndex(),
				rf.getLastLogTerm(),
			}
			reply := &RequestVoteReply{} 
		
			ret := rf.sendRequestVote(i,args,reply)
			if ret != true {
				return
			}
			
			if reply.Term > rf.currentTerm {
				rf.becomeFollower(reply.Term)
			}

			if reply.VoteGranted == true {
				atomic.AddInt32(&votes,1)
			}

			if atomic.LoadInt32(&votes) > int32( len(rf.peers) / 2) {
				rf.becomeLeader()
			}
		}(i)
	}
}

func (rf *Raft) getLastLogIndex() int{
	return len(rf.log) - 1
}

func (rf *Raft) getLastLogTerm() int{
	index := rf.getLastLogIndex()
	if (index < 0){
		return -1
	}
	return rf.log[index].Term
}


// Send  Append Entries 
func (rf *Raft) sendAppendEntries (server int, args *AppendEntry, reply *AppendEntryReply) bool {
	ok := rf.peers[server].Call("Raft.AppendEntries", args, reply)
	return ok
}

//
// example AppendEntries RPC handler.
//
func (rf *Raft) AppendEntries(args *AppendEntry, reply *AppendEntryReply) {
	// Your code here (2A, 2B).

	// Notify state manage rountine to update relection time 
	rf.resetChan <- true
	// Reply Message
	reply.Success = false
	reply.Term = rf.currentTerm

	// If Append Entry's term is bigger than server's update local term or change state 
	if args.Term > rf.currentTerm {
		rf.becomeFollower(args.Term)
	}

	// If  Append Entry's term is less than server's, server should directly refuse append request 
	if args.Term < rf.currentTerm {
		return
	}	


/* 	if args.PrevLogIndex >= len(rf.log){
		return
	} */
	prevLogTerm := -1
	if args.PrevLogIndex >= 0 && args.PrevLogIndex < len(rf.log){
		prevLogTerm = rf.log[args.PrevLogIndex].Term
	}
	if prevLogTerm != args.PrevLogTerm {
		return
	}

	reply.Success = true

	// Append entry to Follower's log 
	index := args.PrevLogIndex 
	for i := 0; i < len(args.Entries); i++ {
		// find the right place to begin to append log 
		index++
		// if append log's index bigger than local log's max size, append logs should immedately 
		// be accepted 
		if(index >= len(rf.log)){
			rf.log = append(rf.log,args.Entries[i:]...)
			//fmt.Println(rf.me , " Append Log ",args.Entries)
			break
		}
	
		// If an existing entry conflicts with a new one 
		// (same index but different terms), delete the existing entry and all that follow
		if(rf.log[index].Term != args.Entries[i].Term){
			rf.log = rf.log[:index]
			rf.log = append(rf.log,args.Entries[i:]...)
			//fmt.Println(rf.me , " Append Log", args.Entries[i:])
			break
		}	

	}

	if args.LeaderCommit > rf.commitIndex{
		rf.commitIndex = Min(args.LeaderCommit,rf.getLastLogIndex())
		//fmt.Println("Commit Log")
		rf.updateLastApplied()
	}



}

func  Min(a int, b int) int {
	if a > b{
		return b
	}
	return a
}
func (rf *Raft) updateLastApplied(){
 	for rf.lastAplied < rf.commitIndex{
		rf.lastAplied++
		applyMsg := ApplyMsg{
			true,
			rf.log[rf.lastAplied].Command,
			rf.lastAplied,
		}
		rf.applyCh <- applyMsg
	}
}

func (rf *Raft) getPrevLogIndex(index int) int {
	preIndex := rf.nextIndex[index] - 1
	if preIndex >= 0{
		return preIndex
	}
	return 0
}
func (rf *Raft) getPrevLogTerm(index int) int {
	preIndex := rf.getPrevLogIndex(index)
	if (preIndex < 0 || preIndex >= len(rf.log ) ) {
		fmt.Println("out of range ", preIndex, " ", len(rf.log)  )
		return 0
	}


	return rf.log[preIndex].Term 
}


func (rf *Raft) leaderAppendEntries(){

	for i:=0; i < len(rf.peers); i++ {
		if i == rf.me{
			continue
		}
		//Concurrent execution can improve system's performance
		go func(index int){
			// Other goroutines may change server's state
			if rf.state != Leader{
				return
			}
			preIndex := rf.getPrevLogIndex(index)
			preTerm := rf.getPrevLogTerm(index)
			appendLog := append(make([]LogEntry,0), rf.log[rf.nextIndex[index]:]...)
			
 			args := &AppendEntry{
				rf.currentTerm,
				rf.me,
				preIndex,
				preTerm,
				appendLog,        
				rf.commitIndex,
			}
			reply := &AppendEntryReply{}

			ret := rf.sendAppendEntries(index,args,reply)

			if (ret != true){
				return
			}
			if (reply.Term > rf.currentTerm){
				rf.becomeFollower(reply.Term)
				return
			}

			if reply.Success {
				rf.matchIndex[index] = args.PrevLogIndex + len(args.Entries) 
				nextIndex := rf.nextIndex[index] + len(args.Entries)
				rf.nextIndex[index] = nextIndex
				rf.updateCommitIndex()
			}else{
				if (rf.nextIndex[index] > 1){
					rf.nextIndex[index] -= 1
				}
			}
 
		}(i)	
	}
}

func (rf *Raft) updateCommitIndex(){
	rf.matchIndex[rf.me] = len(rf.log) - 1
	copyMatchIndex := make([]int,len(rf.matchIndex))
	copy(copyMatchIndex, rf.matchIndex)
	sort.Ints(copyMatchIndex)
	N := copyMatchIndex[len(copyMatchIndex) / 2]
	if N > rf.commitIndex && rf.log[N].Term == rf.currentTerm{
		rf.commitIndex = N	
		rf.updateLastApplied()
	}
}
/* ============ Raft Eclection ========== */

func (rf *Raft) becomeCandidate() {
	rf.mu.Lock()
	rf.resetChan <- true
	rf.state = Candidate
	rf.currentTerm++    
	rf.isVoted = true
	rf.mu.Unlock() 

	go rf.startElection()
}

func (rf *Raft) becomeLeader(){
/* 	rf.mu.Lock()
	defer rf.mu.Unlock() */

	rf.resetChan <- true
	if (rf.state != Candidate){
		return
	}
	// Debug Information
	//fmt.Println("A leader produced by election ! ", rf.number, " term: " , rf.currentTerm)
	rf.state = Leader
	rf.nextIndex = make([]int, len(rf.peers))
	rf.matchIndex = make([]int, len(rf.peers))
	
	// Init nextIndex 
	for i:=0; i < len(rf.nextIndex); i++{
		rf.nextIndex[i] =  rf.getLastLogIndex() + 1
	}
}


func (rf *Raft) becomeFollower(term int){
/* 	rf.mu.Lock()
	defer rf.mu.Unlock() */ 
	rf.state = Follower
	rf.currentTerm = term
	rf.isVoted = false
	rf.resetChan <- true
}


/* ========== Storage ========= */

func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	// w := new(bytes.Buffer)
	// e := labgob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}


//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }
}


/* ==============  Management ============ */
//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me

	// Your initialization code here (2A, 2B, 2C).
	
	// Init server's state
	rf.state = Follower
	rf.currentTerm = 0
	rf.commitIndex = 0 
	rf.lastAplied = 0
	rf.isVoted = false
	rf.nextIndex = make([]int, len(peers))
	rf.matchIndex = make([]int, len(peers))
	rf.log = make([]LogEntry,1)
	rf.resetChan = make(chan bool,1)
	rf.applyCh = applyCh
	// Debug information to give every server a number 
	N++
	rf.number = N 

	//Create a background goroutine to manage state's change 
	go func(){
		// Loop and change raft's state
		for {
			recelectTime := time.Duration(rand.Intn(200) + 150) * time.Millisecond
			switch rf.state {

				case Candidate, Follower : {
					select{
						case <- rf.resetChan :
						case <- time.After(recelectTime):
						rf.becomeCandidate()
					}
				}
				case Leader:{
					rf.leaderAppendEntries()
					time.Sleep(time.Duration(50) * time.Millisecond)
				}
			}
		}
	}()
	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	return rf
}




// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {
	// Your code here (2A).

	rf.mu.Lock()
	defer rf.mu.Unlock() 

	var term int
	var isleader bool
	term = rf.currentTerm
	isleader = rf.state == Leader

	return term, isleader
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//






//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//


func (rf *Raft) Start(command interface{}) (int, int, bool) {
	// Your code here (2B).
	//fmt.Print(rf.me, "  Start ")
	rf.mu.Lock()
	defer rf.mu.Unlock()
	
	index := -1
	term := rf.currentTerm
	isLeader := rf.state == Leader


	if isLeader {
		index = rf.getLastLogIndex() + 1
		newLog := LogEntry{
				term,
				command,
		}
		rf.log = append(rf.log,newLog) 
		//fmt.Println("Insert log ", command)
	}
	return index, term, isLeader
}

//
// the tester doesn't halt goroutines created by Raft after each test,
// but it does call the Kill() method. your code can use killed() to
// check whether Kill() has been called. the use of atomic avoids the
// need for a lock.
//
// the issue is that long-running goroutines use memory and may chew
// up CPU time, perhaps causing later tests to fail and generating
// confusing debug output. any goroutine with a long-running loop
// should call killed() to check whether it should stop.
//
func (rf *Raft) Kill() {
	atomic.StoreInt32(&rf.dead, 1)
	// Your code here, if desired.
}

func (rf *Raft) killed() bool {
	z := atomic.LoadInt32(&rf.dead)
	return z == 1
}
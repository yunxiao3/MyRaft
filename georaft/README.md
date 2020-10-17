## Geo-Raft 算法





~~~protobuf
service RAFT {
    rpc RequestVote (RequestVoteArgs) returns (RequestVoteReply) {};
    rpc AppendEntries (AppendEntriesArgs) returns (AppendEntriesReply){};
}
message AppendEntriesArgs {
    int32 Term = 1;      //  "leader’s term"
    int32 LeaderId = 2;  //  "so follower can redirect clients"
    int32 PrevLogIndex = 3; // "index of log entry immediately preceding new ones"
    int32 PrevLogTerm = 4;      //"term of prevLogIndex entry"
    bytes Log = 5;  //"log entries to store "
    int32 LeaderCommit = 6;    // "leader’s commitIndex"
}
message AppendEntriesReply{
    int32 Term = 1;             
    bool Success = 2;       
    int32 ConflictIndex  = 3;
    int32 ConflictTerm  = 4;
}


service Secretary {
    rpc L2SAppendEntries (L2SAppendEntriesArgs) returns (RequestVoteReply) {}
}

message L2SAppendEntriesArgs {
    int32 Term = 1;         //  "leader’s term"
    int32 PrevLogIndex = 3; // "index of log entry immediately preceding new ones"
    int32 PrevLogTerm = 4;  //  "term of prevLogIndex entry"
    bytes Log = 5;          //  "log entries to store 
    int32 LeaderCommit = 6; //  "leader’s commitIndex"
}
message L2SAppendEntriesReply{
    int32 Term = 1;             
    bool  Success = 2;       
    bytes Matchindex = 3;
}

~~~


#include<stdio.h>


struct raft_state{

    /* latest term server has seen (initialized to 0
    on first boot, increases monotonically) */
    int currentTerm;
    /* candidateId that received vote in current
    term (or null if none) */
    int votedFor;
    /* og entries; each entry contains command
    for state machine, and term when entry
    was received by leader (first index is 1) */
    void* logs;
    /* index of highest log entry known to be
    committed (initialized to 0, increases monotonically)*/
    int commitIndex;
    /* 
    index of highest log entry applied to state
    machine (initialized to 0, increases monotonically) */
    int lastApplied;

    /* data */
};

struct RequestVote{
    /* candidate’s term */
    int term;   
    /* candidate requesting vote */
    int candidateId;
    /* index of candidate’s last log entry */
    int lastLogIndex;
    /* term of candidate’s last log entry */
    int lastLogTerm;
    /* true means candidate received vote */
    int voteGranted;
};


struct  AppendEntries{
    /* leader’s term */
    int term;
    /* so follower can redirect clients */
    int leaderId;
    /* index of log entry immediately preceding new ones */
    int prevLogIndex;
    /* term of prevLogIndex entry */
    int prevLogTerm;
    /* log entries to store (empty for heartbeat;
    may send more than one for efficiency) */
    void* entries;
    /* leader’s commitIndex */
    int leaderCommit;
    /* true if follower contained entry matching
    prevLogIndex and prevLogTerm */
    int success;


};

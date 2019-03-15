#include<stdio.h>

typedef enum {
	FOLLOWER,
	CANDIDATE,
	LEADER
} raft_role_t;

typedef struct {
	int heartbeat_time;
	int election_min_time;
	int election_max_time;
	char *ip;
	int port;
} raft_config_t;

typedef struct {

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

    /* index of highest log entry applied to state
    machine (initialized to 0, increases monotonically) */
    int lastApplied;

} raft_state;

typedef struct {
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
} RequestVote;


typedef struct{
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
} AppendEntries ;

typedef struct
{
    /* size of array */
    int size;

    /* the amount of elements in the array */
    int count;

    /* position of the queue */
    int front, back;

    /* we compact the log, and thus need to increment the Base Log Index */
    int base;

    // raft_entry_t* entries;

    /* callbacks */
    // raft_cbs_t *cb;
    // void* raft;
} log;


typedef {
    /* raft basic setting */
    raft_config_t config;

    raft_role_t role;
    int term;
    void *logs;



}
#include<stdio.h>
#include<stdlib.h>

#define RAFT_MSG_UPDATE 0 // append entry
#define RAFT_MSG_DONE   1 // entry appended
#define RAFT_MSG_CLAIM  2 // vote for me
#define RAFT_MSG_VOTE   3 // my vote

#define bool char
#define true 1
#define false 0

typedef enum {
	FOLLOWER,
	CANDIDATE,
	LEADER
} raft_role_t;

typedef struct {
    /*the unit is ms */
	int heartbeat_time;
	int election_min_time;
	int election_max_time;
    int timeout;

	char *ip;
	int port;
} *raft_config_t;

typedef struct{
    /* raft basic setting */
    raft_config_t config;

    /* if the role is leader time means heartbeat time otherwise means election time or timeout's time */
    int time;
    /* Election of the total votes */
    int votes;
    /* the node role */
    raft_role_t role;
    /* lastest term */
    int term;
    /* the leader*/
    raft_node_t leader;

    /* the total node numbers */
    int peernum;

    void *logs;
} *raft_node_t; 


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








/*
 * Raft message "class" hierarchy:
 *
 *   raft_msg_data_t <-- raft_msg_update_t
 *                   <-- raft_msg_done_t
 *                   <-- raft_msg_claim_t
 *                   <-- raft_msg_vote_t
 *
 * 'update' is sent by a leader to all other peers
 *   'done' is sent in reply to 'update'
 *  'claim' is sent by a candidate to all other peers
 *   'vote' is sent in reply to 'claim'
 */

typedef struct raft_msg_data_t {
	int msgtype;
	int curterm;
	int from;
	int seqno;
} raft_msg_data_t;

typedef struct raft_msg_update_t {
	raft_msg_data_t msg;
	bool snapshot; // true if this message contains a snapshot
	int previndex; // the index of the preceding log entry
	int prevterm;  // the term of the preceding log entry

	bool empty;    // the message is just a heartbeat if empty

	int entryterm;
	int totallen;  // the length of the whole update

	int acked;     // the leader's acked number

	int offset;    // the offset of this chunk inside the whole update
	int len;       // the length of the chunk
	char data[1];
} raft_msg_update_t;

typedef struct raft_msg_done_t {
	raft_msg_data_t msg;
	int entryterm;  // the term of the appended entry
	// raft_progress_t progress; // the progress after appending
	int applied;
	bool success;
	// the message is considered acked when the last chunk appends successfully
} raft_msg_done_t;

typedef struct raft_msg_claim_t {
	raft_msg_data_t msg;
	int index; // the index of my last completely received entry
	int lastterm;  // the term of my last entry
} raft_msg_claim_t;

typedef struct raft_msg_vote_t {
	raft_msg_data_t msg;
	bool granted;
} raft_msg_vote_t;

typedef union {
	raft_msg_update_t u;
	raft_msg_done_t d;
	raft_msg_claim_t c;
	raft_msg_vote_t v;
} raft_msg_any_t;

// typedef raft_node *raft_node_t;



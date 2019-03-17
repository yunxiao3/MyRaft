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

typedef struct raft_peer_t {
	bool up;

	int seqno;  // the rpc sequence number
	// raft_progress_t acked; // the number of entries:bytes acked by this peer
	int applied; // the number of entries applied by this peer

	char *host;
	int port;
	struct sockaddr_in addr;
	int silent_ms; // how long was this peer silent
} raft_peer_t;

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

    /*the sock handle  */
    int sock;

    /* the peers */
    raft_peer_t *peers;

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


/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */

typedef struct
{
    void *buf;

    unsigned int len;
} raft_entry_data_t;

typedef struct
{
    /** the entry's term at the point it was created */
    int term;

    /** the entry's unique ID */
    int id;

    /** type of entry */
    int type;

    raft_entry_data_t data;
} msg_entry_t;

/** requestVote 请求投票
   * 竞选者Candidate去竞选Leader时发送给其它node的投票请求。
   * 其它Leader或者Candidate收到term比自己大的投票请求时，会自动变成Follower*/
typedef struct
{
    /** 当前任期号，通过任期号的大小与其它Candidate竞争Leader */
    int term;

    /** 竞选者的id */
    int candidate_id;

    /** 竞选者本地保存的最新一条日志的index */
    int last_log_idx;

    /** 竞选者本地保存的最新一条日志的任期号*/
    int last_log_term;
} msg_requestvote_t;


/** 投票请求的回复response.
  * 该response主要是给返回某个node是否接收了Candidate的投票请求. */
typedef struct
{
    /** node的任期号，Candidate根据投票结果和node的任期号来更新自己的任期号 */
    int term;

    /** 投票结果，如果node给Candidate投票则为true */
    int vote_granted;
} msg_requestvote_response_t;

/**  添加日志请求.
  * Follower可以从该消息中知道哪些日志可以安全地提交到状态机FSM中去。
  * Leader可以将该消息作为心跳消息定期发送。
  * 旧的Leader和Candidate收到该消息后可能会自动变成Follower */
typedef struct
{
    /** Leader当前的任期号 */
    int term;

    /** 最新日志的前一条日志的index，用于Follower确认与Leader的日志完全一致 */
    int prev_log_idx;

    /** 最新日志的前一条日志的任期号term */
    int prev_log_term;

    /** leader当前已经确认提交到状态机FSM的日志索引index，这意味着Follower也可以安全地将该索引index以前的日志提交 */
    int leader_commit;

    /** 这条添加日志消息携带的日志条数，该实现中最多只有一条 */
    int n_entries;

    /** 这条添加日志消息中携带的日志数组 */
    msg_entry_t* entries;
} msg_appendentries_t;

/** 添加日志回复.
 * 旧的Leader或Candidate收到该消息会变成Follower *／
 * */
typedef struct
{
    /** 当前任期号 */
    int term;

    /** node成功添加日志时返回ture，即prev_log_index和prev_log_term都比对成功。否则返回false */
    int success;

    /* 下面两个字段不是Raft论文中规定的字段:
    /* 用来优化日志追加过程，以加速日志的追加。Raft原文中的追加过程是一次只能追加一条日志*/

    /** 处理添加日志请求后本地的最大日志索引 */
    int current_idx;

    /** 从添加日志请求中接受的第一条日志索引 */
    int first_idx;
} msg_appendentries_response_t;



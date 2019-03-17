#include<stdio.h>

typedef void* raft_node_t;

typedef enum {
	FOLLOWER,
	CANDIDATE,
	LEADER
} raft_role_t;


typedef struct{
    void *buf;
    unsigned int len;
} raft_entry_data_t;

typedef struct {
    /*the unit is ms */
	int heartbeat_time;
	int election_min_time;
	int election_max_time;
    int timeout;
	char *ip;
	int port;
} raft_config_t;


typedef struct{
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
typedef struct{
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


typedef struct {
    /* 所有服务器比较固定的状态: */

    /* 服务器最后一次知道的任期号（初始化为 0，持续递增） */
    int current_term;

    /* 记录在当前分期内给哪个Candidate投过票， */
    int voted_for;

    /* 日志条目集；每一个条目包含一个用户状态机执行的指令，和收到时的任期号 */
    void* log;

    /* 变动比较频繁的变量: */

    /* 已知的最大的已经被提交的日志条目的索引值 */
    int commit_idx;

    /* 最后被应用到状态机的日志条目索引值（初始化为 0，持续递增） */
    int last_applied_idx;

    /* 三种状态：follower/leader/candidate */
    raft_role_t role;

    /* 计时器，周期函数每次执行时会递增改值 */
    // int timeout_elapsed;

    raft_node_t* nodes;
    int num_nodes;

    /* 每个noed的配置信息 */
    raft_config_t* config;
    
    /* leader表示心跳时间，follower表示超时时间，candidate表示竞选使劲 */
    int time;

    // int election_timeout;
    // int request_timeout;

    /* 保存Leader的信息，没有Leader时为NULL */
    raft_node_t* current_leader;

    /* callbacks，由调用该raft实现的调用者来实现，网络IO和持久存储
     * 都由调用者在callback中实现  暂时不实现这一部分功能*/
    // raft_cbs_t cb;


    void* udata;

    /* 自己的信息 */
    raft_node_t* node;

    /* 该raft实现每次只进行一个服务器的配置更改，该变量记录raft server
     * 是否正在进行配置更改*/
    // int voting_cfg_change_log_idx;
} raft_server_private_t;

typedef struct
{
    void* udata;  /*一般保存与其它机器的连接信息，由使用者决定怎么实现连接*/

    int next_idx; /*对于每一个服务器，需要发送给他的下一个日志条目的索引值（初始化为领导人最后索引值加一）*/
    int match_idx; /*对于每一个服务器，已经复制给他的日志的最高索引值*/

    int flags; /*有三种取值，是相或的关系 1:该机器有给我投票 2:该机器有投票权  3: 该机器有最新的日志*/

    int id; /*机器对应的id值，这个每台机器在全局都是唯一的*/
} raft_node_private_t;
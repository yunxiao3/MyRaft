#include<stdio.h>
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

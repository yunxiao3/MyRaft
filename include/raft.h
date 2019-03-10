typedef struct {

    raft_term_t current_term;

    raft_node_id_t voted_for;

    void* log;

    raft_index_t commit_idx;


    raft_index_t last_applied_idx;


    int state;

    int timeout_elapsed;

    raft_node_t* nodes;
    int num_nodes;

    int election_timeout;
    int election_timeout_rand;
    int request_timeout;


    raft_node_t* current_leader;

    raft_cbs_t cb;
    void* udata;

    raft_node_t* node;


    raft_index_t voting_cfg_change_log_idx;


    int connected;

    int snapshot_in_progress;
    int snapshot_flags;


    raft_index_t snapshot_last_idx;
    raft_term_t snapshot_last_term;

    raft_index_t saved_snapshot_last_idx;
    raft_term_t saved_snapshot_last_term;
} raft_server_private_t;

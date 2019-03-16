#include"../include/raft.h"

void raft_peer_init(raft_peer_t *p) {
	p->up = false;
	p->seqno = 0;
	reset_progress(&p->acked);
	p->applied = 0;

	p->host = DEFAULT_LISTENHOST;
	p->port = DEFAULT_LISTENPORT;
	p->silent_ms = 0;
}
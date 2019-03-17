#include"../include/raft.h"

#define DEFAULT_LISTENHOST "0.0.0.0"
#define DEFAULT_LISTENPORT 6543

void raft_peer_init(raft_peer_t *p) {
	p->up = false;
	p->seqno = 0;
	// reset_progress(&p->acked);
	p->applied = 0;

	p->host = DEFAULT_LISTENHOST;
	p->port = DEFAULT_LISTENPORT;
	p->silent_ms = 0;
}

/* bool raft_peers_init(raft_node_t raft) {
	int i;
	raft->peers = malloc(raft->config.peernum_max * sizeof(raft_peer_t));
	if (!raft->peers) {
		shout("failed to allocate memory for raft peers\n");
		return false;
	}
	for (i = 0; i < raft->config.peernum_max; i++) {
		raft_peer_init(raft->peers + i);
	}
	return true;
} */

static void raft_send(raft_node_t r, int dst, void *m, int mlen) {
	/* assert(r->peers[dst].up);
	assert(mlen <= r->config.msg_len_max);
	assert(msg_size_is((raft_msg_t)m, mlen));
	assert(((raft_msg_t)m)->msgtype >= 0);
	assert(((raft_msg_t)m)->msgtype < 4);
	assert(dst >= 0);
	assert(dst < r->config.peernum_max);
	assert(dst != r->me);
	assert(((raft_msg_t)m)->from == r->me);
 */
	raft_peer_t *peer = r->peers + dst;

	int sent = sendto(
		r->sock, m, mlen, 0,
		(struct sockaddr*)&peer->addr, sizeof(peer->addr)
	);
	/* if (sent == -1) {
		shout(
			"failed to send a msg to [%d]: %s\n",
			dst, strerror(errno)
		);
	} */
}


int raft_create_udp_socket(raft_node_t r) {
	// assert(r->me != NOBODY);
	raft_peer_t *me = r->peers + r->me;
	struct addrinfo *hint;
	struct addrinfo *addrs = NULL;
	struct addrinfo *a;
	char portstr[6];
	int rc;
	memset(&hint, 0, sizeof(hint));
	hint->ai_socktype = SOCK_DGRAM;
	hint->ai_family = AF_INET;
	hint->ai_protocol = getprotobyname("udp")->p_proto;

	snprintf(portstr, 6, "%d", me->port);

	if ((rc = getaddrinfo(me->host, portstr, &hint, &addrs)) != 0)
	{
		shout(
			"cannot convert the host string '%s'"
			" to a valid address: %s\n", me->host, gai_strerror(rc));
		return -1;
	}

	for (a = addrs; a != NULL; a = a->ai_next)
	{
		int sock = socket(AF_INET, SOCK_DGRAM, IPPROTO_UDP);
		if (sock < 0) {
			shout("cannot create socket: %s\n", strerror(errno));
			continue;
		}
		socket_set_reuseaddr(sock);
		socket_set_recv_timeout(sock, r->config.heartbeat_ms);

		debug("binding udp %s:%d\n", me->host, me->port);
		if (bind(sock, a->ai_addr, a->ai_addrlen) < 0) {			
			shout("cannot bind the socket: %s\n", strerror(errno));
			close(sock);
			continue;
		}
		r->sock = sock;
		assert(a->ai_addrlen <= sizeof(me->addr));
		memcpy(&me->addr, a->ai_addr, a->ai_addrlen);
		return sock;
	}

	shout("cannot resolve the host string '%s' to a valid address\n",
		me->host
	);
	return -1;
}




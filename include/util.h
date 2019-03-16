#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <netdb.h>
#include <unistd.h>
#include <assert.h>
#include <errno.h>
#include <limits.h>
#include <string.h>
#include<stdio.h>
#include<stdlib.h>
#include<time.h>
#include"raft.h"

int getRandomTime(int min, int max){
    srand((unsigned)time(NULL));
    int a = min - 10;
    while(a < min){
        a = rand() % max;
    }
    return a;
}




void raft_reset_time(raft_node_t r) {
	if (r->role == LEADER) {
		r->time = r->config->heartbeat_time;
	} 
    if(r->role == CANDIDATE){
		r->time = getRandomTime(
			r->config->election_min_time,
			r->config->election_max_time
		);
	}
    if(r->role == FOLLOWER){
        r -> time = r->config->timeout;
    }
}







void initConfig(raft_config_t c, char *ip, int port){
    c->heartbeat_time = 150;
	c->election_min_time = 100;
	c->election_max_time = 200;
    c->timeout = 50;
	// memcpy(&raft->config, config, sizeof(raft_config_t));
	c->port = port;
}

void initRaftNode(raft_node_t me){
    
    me->role = FOLLOWER;
    raft_reset_time(me);
    me->term = 0;

}


bool raft_become_leader(raft_node_t r) {
	if (r->votes * 2 > r->peernum) {
		// got the support of a majority
		r->role = LEADER;
		// r->leader = r->me;
/* 		raft_reset_bytes_acked(r);
		raft_reset_silent_time(r, NOBODY);
		raft_reset_timer(r);
		shout("became the leader\n"); */
		return true;
	}
	return false;
}

/* int raft_create_udp_socket(raft_node_t r) {

	struct addrinfo hint;
	struct addrinfo *addrs = NULL;
	struct addrinfo *a;
	char portstr[6];
	int rc;
	memset(&hint, 0, sizeof(hint));
	hint.ai_socktype = SOCK_DGRAM;
	hint.ai_family = AF_INET;
	hint.ai_protocol = getprotobyname("udp")->p_proto;

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
} */


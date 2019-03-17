#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<errno.h>
#include<sys/types.h>
#include<sys/socket.h>
#include<netinet/in.h>
#include<termios.h>
#include<sys/types.h>   
#include<sys/stat.h>    
#include<fcntl.h>
#include<unistd.h>
#include<sys/ioctl.h>
#include<signal.h>
#include<errno.h>
 
#define MAXLINE 256
#define PORT	6666
int listenfd;
int connfd;
pthread_t read_id, write_id;
 
/*
linux ctrl + C 会产生 SIGINT信号
接收到SIGINT 信号进入该函数
*/
void stop(int signo)
{
	printf("stop\n");
	close(connfd);	
	close(listenfd);
	 _exit(0);
}
 
/*
当客户端断开连接的时候，
在服务端socket send进程可以收到收到信号SIGPIPE，
收到SIGPIPE信号进入该函数结束创建的线程。
*/
void signal_pipe(int signo)
{
	pthread_kill(read_id,SIGQUIT);//向read线程发送SIGQUIT
	pthread_join(read_id,NULL);	  //阻塞线程运行，直到read 线程退出。
	
	close(connfd);                //关闭连接
	printf("read pthread out \n");
	
	pthread_exit(0);              //结束write 线程
}
 
/*
read 线程接收到SIGQUIT信号，
执行线程退出操作
*/
void pthread_out(int signo)
{
	pthread_exit(0);
}
 
/*
read 线程执行函数
*/
void* read_func(void* arg)
{
	char readbuff[MAXLINE];
	int n = 0;
	int fd;
 
	fd = *(int*)arg;    /*main 主进程传递过来的连接文件描述符*/
	memset(&readbuff,0,sizeof(readbuff));
 
	signal(SIGQUIT,pthread_out); /* 注册SIGQUIT 信号*/ 
	while(1)
	{
    	n = recv(fd, readbuff, MAXLINE, 0);  /*recv 在这里是阻塞运行*/
		if(n > 0)
		{
			printf("server recv data: %s \n",readbuff);
		}
	};
}
/*
write 线程执行函数
*/
void* write_func(void* arg)
{
	char writebuff[MAXLINE];
	char* write = "I am server";
	unsigned char i = 0;
	int num = 0;
	int fd;
 
	fd = *(int*)arg;
	memset(&writebuff,0,sizeof(writebuff));
	
	signal(SIGPIPE,signal_pipe); /* 注册 SIGPIPE信号 */
	while(1)
	{
		sleep(1);
		send(fd,write,strlen(write)+1,0);/*向客户端发送数据*/
	}
}
 
int main(int argc, char** argv)
{
	char buff[MAXLINE];
	int num;
	int addrlen;
	struct sockaddr_in server_addr;  /*服务器地址结构*/
	struct sockaddr_in client_addr;  /*客户端地址结构*/
 
	if((listenfd = socket(AF_INET,SOCK_STREAM,0)) == -1)/*建立一个流式套接字*/
	{
		printf("create socket error: %s(errno: %d)\n",strerror(errno),errno);
		exit(0);
	}
	
	/*设置服务端地址*/
	addrlen = sizeof(struct sockaddr_in);
	memset(&server_addr, 0, addrlen);
	server_addr.sin_family = AF_INET;    /*AF_INET表示 IPv4 Intern 协议*/
	server_addr.sin_addr.s_addr = htonl(INADDR_ANY); /*INADDR_ANY 可以监听任意IP */
	server_addr.sin_port = htons(PORT); /*设置端口*/
    
	/*绑定地址结构到套接字描述符*/
	if( bind(listenfd, (struct sockaddr*)&server_addr, sizeof(server_addr)) == -1)
	{
		printf("bind socket error: %s(errno: %d)\n",strerror(errno),errno);
		exit(0);
	}
 
	/*设置监听队列，这里设置为1，表示只能同时处理一个客户端的连接*/
	if( listen(listenfd, 1) == -1)
	{
		printf("listen socket error: %s(errno: %d)\n",strerror(errno),errno);
		exit(0);
	}
	
	signal(SIGINT,stop); /*注册SIGINT 信号*/
	while(1)
	{	
		printf("wait client accpt \n");
		if( (connfd = accept(listenfd, (struct sockaddr*)&client_addr, &addrlen)) == -1)/*接收客户端的连接，这里会阻塞，直到有客户端连接*/
		{
			printf("accept socket error: %s(errno: %d)",strerror(errno),errno);
			continue;
		}
	
		if(pthread_create(&read_id,NULL,read_func,&connfd))/*创建 read 线程*/
		{
			printf("pthread_create read_func err\n");
		}
 
		if(pthread_create(&write_id,NULL,write_func,&connfd))/*创建 write 线程*/
		{
			printf("pthread_create write_func err\n");
		}
		
		pthread_join(write_id,NULL); /*阻塞，直到write进程退出后才进行新的客户端连接*/
		printf("write pthread out \n");
	}
}
 
 
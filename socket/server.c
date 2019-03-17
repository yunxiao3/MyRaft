#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<errno.h>
#include<unistd.h>
#include<sys/wait.h>
#include<sys/types.h>
#include<sys/socket.h>
#include<netinet/in.h>
#define SERV_PORT 9998
#define MAXLINE 4096

int main(int argc, char** argv){
    
    int    socket_fd, connect_fd;
    struct sockaddr_in     servaddr;
    char    buf[MAXLINE],sendbuf[MAXLINE];
    int     len;
    pid_t   pid;
    if( (socket_fd = socket(AF_INET, SOCK_STREAM, 0)) == -1 ){
        printf("create socket error: %s(errno: %d)\n",strerror(errno),errno);
        exit(0);
    }
    //初始化
    memset(&servaddr, 0, sizeof(servaddr));
    servaddr.sin_family = AF_INET;  
    servaddr.sin_addr.s_addr = htonl(INADDR_ANY);//IP地址设置成INADDR_ANY,让系统自动获取本机的IP地址。
    servaddr.sin_port = htons(SERV_PORT);
    //调用bind函数，绑定
    if( bind(socket_fd, (struct sockaddr*)&servaddr, sizeof(servaddr)) <0)
    {
         printf("bind socket error: %s(errno: %d)\n",strerror(errno),errno);
         exit(0);
    }
    //调用listen函数监听端口
    if( listen(socket_fd, 10) <0)
    {
         printf("listen socket error: %s(errno: %d)\n",strerror(errno),errno);
         exit(0);
    }
    printf("waiting for client's connection......\n");

       //调用fork函数实现多进程，注意accept函数在fork调用之后
        if(pid=fork()<0)
        {
            perror("fork ");
            exit(-1);
         }
         if(pid==0)
        {
          connect_fd = accept(socket_fd,(struct sockaddr*)NULL,NULL);   
          while((len= recv(connect_fd, buf, MAXLINE, 0))>0)
        {
       printf("receive message from client: %s\n", buf);
       //向客户端发送回应数据
        printf("send message to client: \n");
        fgets(sendbuf, 4096, stdin);
       if( send(connect_fd, sendbuf, strlen(sendbuf), 0) < 0){
            printf("send messaeg error: %s(errno: %d)\n", strerror(errno), errno);
            exit(0);
        }
        }
         }
  close(connect_fd);
  close(socket_fd);
}

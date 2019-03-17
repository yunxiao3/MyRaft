#include <stdio.h>
#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <netinet/in.h>
#include <errno.h>
#include <stdlib.h>

#define MAXLINE    1024
#define SERV_PORT  9998

int main(int argc, char *argv[])
        {
            char sendbuf[MAXLINE],receivebuf[MAXLINE];
            struct sockaddr_in servaddr;      //定义sockaddr类型结构体servaddr
            int client_sockfd;
            int rec_len;
            // 判断命令端输入的参数是否正确 
             if( argc != 2)
             {
                 printf("usage: ./client <ipaddress>\n");
                 exit(0);
              }
             // 调用socket函数
             if((client_sockfd=socket(AF_INET,SOCK_STREAM,0))<0)
              {
                  perror("socket");//调用perror函数，该函数会自行打印出错信息
                  exit(0);
              }
             //初始化结构体
            memset(&servaddr,0,sizeof(servaddr));      //数据初始化-清零 
            servaddr.sin_family = AF_INET;             //设置IPv4通信 
            servaddr.sin_port = htons(SERV_PORT);      //设置服务器端口号 
            // IP地址转换函数inet_pton，将点分十进制转换为二进制 
             if( inet_pton(AF_INET, argv[1], &servaddr.sin_addr) <= 0)
              {
                   printf("inet_pton error for %s\n",argv[1]);
                   exit(0);
              }
             //调用connect函数
            if( connect(client_sockfd, (struct sockaddr *)&servaddr, sizeof(servaddr))<0)
              {
                  perror("connected failed");
                  exit(0);
              }
            //循环发送接收数据，send发送数据，recv接收数据 

       while(1)
           {
             printf("send msg to server: \n");
             fgets(sendbuf, 1024, stdin);//从标准输入读数据
             // 向服务器端发送数据 
             if( send(client_sockfd, sendbuf, strlen(sendbuf), 0) < 0)
              {
                  printf("send msg error: %s(errno: %d)\n", strerror(errno), errno);
                  exit(0);
              }
             // 接受服务器端传过来的数据 
             if((rec_len = recv(client_sockfd,receivebuf, MAXLINE,0)) == -1)
              {
                   perror("recv error");
                   exit(1);
              }
             printf("Response from server: %s\n",receivebuf);
           }
       //关闭套接字 
    close(client_sockfd);
    return 0;
        }



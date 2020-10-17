
package main

import (
	 SE "../georaft"
	 "fmt"
	 "log"
	 "golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
	RPC "../grpc/georaft"
	
)







func  L2SsendAppendEntries(address string , args  *RPC.L2SAppendEntriesArgs){

	fmt.Println("StartAppendEntries")

	// Initialize Client
	conn, err := grpc.Dial( address , grpc.WithInsecure(),grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := RPC.NewSecretaryClient(conn)


	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//args := &RPC.AppendEntriesArgs{}
	r, err := client.L2SAppendEntries(ctx,args)
	if err != nil {
		log.Printf("could not greet: %v", err)
	}
	log.Printf("Append reply: %s", r)
	//fmt.Println("Append name is ABC")
}



func main()  {

	members := make([]string, 2)

	members[0] = "localhost:5000"
	members[1] = "localhost:5001"

	se := SE.MakeSecretary(members)


	if (se != nil){
		fmt.Println("SUCCESS")
		 for{
			args := &RPC.L2SAppendEntriesArgs{}
			L2SsendAppendEntries("localhost:5000", args)
			time.Sleep(time.Second )

			fmt.Println("L2SsendAppendEntries")
		} 
	}else{
		fmt.Println("FAIL")
	}



}
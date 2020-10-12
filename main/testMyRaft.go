package main


import (
	"../myraft"
	 "os"
	/*"fmt"*/
	"time"
)






 func main() {

	raft := myraft.MakeRaft(os.Args)

	raft.GetState()

	time.Sleep(time.Second*12)



	/* if (len(os.Args) > 1){
		raft.members = make([]string, len(os.Args) - 1)
		for i:= 0; i < len(os.Args) - 1; i++{
			//fmt.Printf("args[%v]=[%v]\n",k,v)
			raft.members[i] = os.Args[i+1]
			fmt.Printf(raft.members[i])
		}
	}

	
	raft.init(raft.members[0])
	fmt.Println("init raft address " , raft.members[0])	

	time.Sleep(time.Second*12)


	time.Sleep(time.Second*200) */


}  
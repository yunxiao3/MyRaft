 package main
 
 
 import (
	 "log"
	// "os"
	 "time"
	 "golang.org/x/net/context"
	 "google.golang.org/grpc"
	 pb "./raft"
 )
  
 const (
	 address     = "localhost:50051"
	 defaultName = "world"
 )


 type Raft struct{
	client pb.RAFTClient
}

func (raft *Raft) init(){
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
//	defer conn.Close()
	raft.client = pb.NewRAFTClient(conn)
}

func (raft *Raft)RequestVote(){

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := raft.client.RequestVote(ctx, &pb.RequestVoteArgs{Name: "name"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}


 func main() {

	raft := Raft{}
	raft.init()
	raft.RequestVote()
	 // Set up a connection to the server.
	/*  conn, err := grpc.Dial(address, grpc.WithInsecure())
	 if err != nil {
		 log.Fatalf("did not connect: %v", err)
	 }
	 defer conn.Close()
	 c := pb.NewRAFTClient(conn)
  
	 // Contact the server and print out its response.
	 name := defaultName
	 if len(os.Args) > 1 {
		 name = os.Args[1]
	 }
	 ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	 defer cancel()
	 r, err := c.RequestVote(ctx, &pb.RequestVoteArgs{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)


	 for  i := 0;  i < 1500; i++ {
		r, err = c.RequestVote(ctx, &pb.RequestVoteArgs{Name: name})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r.Message)
	 } */
	 
 }
 
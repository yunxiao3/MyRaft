package main

//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto
 
import (
	"log"
	"net"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	RPC "./grpc/raft"
 
)
 
const (
	port = ":50051"
)
 
// server is used to implement helloworld.GreeterServer.
type Raft struct{

	ip string 
	port string
	msg string
}
 
// SayHello implements helloworld.GreeterServer
 func (raft *Raft) RequestVote(ctx context.Context, in *RPC.RequestVoteArgs) (*RPC.RequestVoteReply, error) {
	/* raft.msg = raft.msg + in.Name

	if (in.Name == "ABC"){
		fmt.Println("RequestVoteArgs name is ABC")
		
		return &pb.RequestVoteReply{Message: "Hello ABC " + in.Name}, nil
	}else{
		fmt.Println("RequestVoteArgs name is raft ")
		return &pb.RequestVoteReply{Message: "Hello RAFT" + raft.msg}, nil
	} */
	return &RPC.RequestVoteReply{}, nil
} 
 
func (raft *Raft)AppendEntries(ctx context.Context, in *RPC.AppendEntriesArgs) (*RPC.AppendEntriesReply, error) {
	//raft.msg = raft.msg + in.Term
	reply := &RPC.AppendEntriesReply{}
	reply.Term = 1
    reply.Success = false
    reply.ConflictTerm = -1
    reply.ConflictIndex = 0
	if (in.Term == 1){
		fmt.Println("Append name is ABC")
		return reply, nil
	}else{
		return reply, nil
	}
}

func (raft *Raft)init(port string, ip string){
	raft.port = port
	raft.ip = ip
	raft.msg = ""


	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RPC.RegisterRAFTServer(s, &Raft{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}



func main() {
	raft := Raft{}
	raft.init(":50051","localhost")
}

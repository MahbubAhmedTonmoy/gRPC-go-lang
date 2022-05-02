package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	pb "github.com/MahbubAhmedTonmoy/GoGrpc/student/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var address string = "0.0.0.0:50051"

type Server struct {
	pb.StudentServiceServer
}

func main() {
	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("Failed to listen on: %v\n", err)
	}
	log.Printf("Listening on %s\n", address)
            
	s := grpc.NewServer()
	pb.RegisterStudentServiceServer(s, &Server{})
	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v\n", err)
	}
}
func (s *Server) Student(ctx context.Context, in *pb.StudentRequest) (*pb.StudentResult, error) {
	log.Printf("Student function was invoked with %v\n", in)

	return &pb.StudentResult{
		Result: "hello" + in.Name,
	}, nil
}

func (s *Server) LongStudent(stream pb.StudentService_LongStudentServer) error {
	log.Println("long student funce was invocked")

	res := ""

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.StudentResult{
				Result: res,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream %v\n", err)
		}
		res += fmt.Sprintf("hello %s!\n", req.Name)
	}
}
func (s *Server) StudentManyTime(in *pb.StudentRequest, stream pb.StudentService_StudentManyTimeServer) error {
	log.Printf("Student many times function was invoked with %v\n", in)

	for i := 0; i < 10; i++ {
		res := fmt.Sprintf("Hello %s, number %d", in.Name, i)
		stream.Send(&pb.StudentResult{
			Result: res,
		})
	}
	return nil
}

func (s *Server) StudentEveryWhere(stream pb.StudentService_StudentEveryWhereServer) error {
	log.Printf("Everywhere was invocked")

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v\n", err)
		}
		res := "Hello " + req.Name + "!"

		err = stream.Send(&pb.StudentResult{
			Result: res,
		})
		if err != nil {
			log.Fatalf("Error while sending data to client: %v\n", err)
		}
	}
}

func (s *Server) StudentDeadLine(ctx context.Context, in *pb.StudentRequest) (*pb.StudentResult, error) {
	log.Printf("StudentDeadLine was invocked")
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.DeadlineExceeded {
			log.Panicf("the client cancled the request")
			return nil, status.Error(codes.Canceled, "client cancled the request")
		}
		time.Sleep(1 * time.Second)
	}
	return &pb.StudentResult{
		Result: "Hello " + in.Name,
	}, nil
}

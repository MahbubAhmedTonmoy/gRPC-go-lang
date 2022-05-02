package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "github.com/MahbubAhmedTonmoy/GoGrpc/student/proto"
)

var add string = "localhost:50051"

func main() {
	conn, err := grpc.Dial(add, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("failed to connecting %v\n", err)
	}
	defer conn.Close()

	c := pb.NewStudentServiceClient(conn)

	// doStudent(c)
	// doStudentManyTimes(c)
	// doLongStudent(c)
	//doEverYWhere(c)
	doDeadLine(c, 5*time.Second)
}
func doStudent(c pb.StudentServiceClient) {
	log.Print("doStudent was invoked")
	res, err := c.Student(context.Background(), &pb.StudentRequest{
		Name: "mahbub ahmed",
	})

	if err != nil {
		log.Fatalf("Could not student %v\n", err)
	}
	log.Printf("Student %s\n", res.Result)
}
func doStudentManyTimes(c pb.StudentServiceClient) {
	log.Print("doStudent many times was invoked")
	stream, err := c.StudentManyTime(context.Background(), &pb.StudentRequest{
		Name: "mahbub ahmed",
	})

	if err != nil {
		log.Fatalf("Could not student %v\n", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {

			log.Fatalf("Error while reading stream %v\n", err)
		}

		log.Printf("Student %s\n", msg.Result)
	}
}

func doLongStudent(c pb.StudentServiceClient) {
	log.Print("do Long Student was invoked")
	reqs := []*pb.StudentRequest{
		{Name: "mahbub"},
		{Name: "ahmed"},
		{Name: "tonmoy"},
		{Name: "boltu"},
	}
	stream, err := c.LongStudent(context.Background())

	if err != nil {
		log.Fatalf("error while calling long student %v\n", err)
	}
	for _, req := range reqs {
		log.Printf("Sending req %v\n", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from long student %v\n", err)
	}

	log.Printf("Long student %s\n", res.Result)
}

func doEverYWhere(c pb.StudentServiceClient) {
	log.Printf("do wvery where invoked")
	stream, err := c.StudentEveryWhere(context.Background())
	if err != nil {
		log.Fatalf("Error while creatiing stream %v\n", err)
	}
	reqs := []*pb.StudentRequest{
		{Name: "mahbub"},
		{Name: "ahmed"},
		{Name: "tonmoy"},
		{Name: "boltu"},
	}
	waitc := make(chan struct{})

	go func() {
		for _, req := range reqs {
			log.Printf("Sending req %v\n", req)
			stream.Send(req)
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v\n", err)
			}
			log.Printf("Received: %v\n", res.Result)
		}
		close(waitc)
	}()
	<-waitc
}

func doDeadLine(c pb.StudentServiceClient, timeout time.Duration) {
	log.Printf("doDeadLine was invoked")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()
	req := &pb.StudentRequest{
		Name: "mahbub",
	}

	res, err := c.StudentDeadLine(ctx, req)

	if err != nil {
		e, ok := status.FromError(err)

		if ok {
			if e.Code() == codes.DeadlineExceeded {
				log.Printf("Deadline Exceeded!")
				return
			} else {
				log.Fatalf("Unexpected grpc error %v\n", err)
			}
		} else {
			log.Fatalf("A non grpc error %v\n", err)
		}
	}

	log.Printf("result %s\n", res.Result)
}

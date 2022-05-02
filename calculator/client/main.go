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

	pb "github.com/MahbubAhmedTonmoy/GoGrpc/calculator/proto"
)

var add string = "localhost:50051"

func main() {
	conn, err := grpc.Dial(add, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal("failed to connecting %v\n", err)
	}
	defer conn.Close()

	c := pb.NewCalculatorServiceClient(conn)

	// doSum(c)
	//doPrime(c)
	//doAvg(c)
	//doMax(c)
	doSqrt(c, -25)
}

func doSum(c pb.CalculatorServiceClient) {
	log.Print("doSum was invoked")
	res, err := c.Sum(context.Background(), &pb.SumRequest{
		FirstNumber:  10,
		SecondNumber: 20,
	})

	if err != nil {
		log.Fatal("Could not sum %v\n", err)
	}
	log.Printf("sum %s\n", res.Result)
}
func doPrime(c pb.CalculatorServiceClient) {
	log.Print("doPrime was invoked")
	stream, err := c.Primes(context.Background(), &pb.PrimeRequest{
		Number: 120,
	})

	if err != nil {
		log.Fatal("error while call prime %v\n", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error while reading stream %v\n", err)
		}

		log.Printf("Primes %d\n", res.Result)
	}
}

func doAvg(c pb.CalculatorServiceClient) {
	log.Printf("do Avg func called")
	stream, err := c.Avg(context.Background())

	if err != nil {
		log.Fatalf("error while opening the stream %v\n", err)
	}
	numbers := []int64{3, 4, 9, 54, 23}

	for _, number := range numbers {
		log.Printf("sending number %d\n", number)

		stream.Send(&pb.AvgRequest{Number: number})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response %v\n", err)
	}
	log.Printf("Avg %f\n", res.Result)
}

func doMax(c pb.CalculatorServiceClient) {
	log.Printf("do Avg func called")
	stream, err := c.Max(context.Background())

	if err != nil {
		log.Fatalf("error while opening the stream %v\n", err)
	}

	waitc := make(chan struct{})

	go func() {
		numbers := []int32{4, 7, 2, 19, 4, 6, 32}
		for _, number := range numbers {
			log.Printf("sending number %d\n", number)
			stream.Send(&pb.MaxRequest{
				Number: number,
			})
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
				log.Fatalf("Error while reading stream %v\n", err)
				break
			}

			log.Printf("max number %d\n", res.Result)
		}
		close(waitc)
	}()
	<-waitc
}
func doSqrt(c pb.CalculatorServiceClient, n int32) {
	log.Print("doSqrt was invoked")
	res, err := c.Sqrt(context.Background(), &pb.SqrtRequest{
		Number: n,
	})

	if err != nil {
		e, ok := status.FromError(err)
		if ok {
			log.Printf("Error message from server %s\n", e.Message())
			log.Printf("Error code from server %s\n", e.Code())

			if e.Code() == codes.InvalidArgument {
				log.Printf("send negative number")
				return
			}
		} else {
			log.Fatalf("a non grpc error %v\n", err)
		}
	}
	log.Printf("sqrt  %f\n", res.Result)
}

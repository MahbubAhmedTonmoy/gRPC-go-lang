package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	pb "github.com/MahbubAhmedTonmoy/GoGrpc/calculator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var address string = "0.0.0.0:50051"

type Server struct {
	pb.CalculatorServiceServer
}

func main() {
	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("Failed to listen on: %v\n", err)
	}
	log.Printf("Listening on %s\n", address)

	s := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(s, &Server{})
	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v\n", err)
	}
}

func (s *Server) Sum(ctx context.Context, in *pb.SumRequest) (*pb.SumResponse, error) {
	log.Printf("Sum function was invoked with %v\n", in)

	result := in.FirstNumber + in.SecondNumber
	return &pb.SumResponse{
		Result: result,
	}, nil
}

func (s *Server) Primes(in *pb.PrimeRequest, stream pb.CalculatorService_PrimesServer) error {
	log.Printf("Primes function was invoked with %v\n", in)
	number := in.Number
	divisor := int64(2)

	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&pb.PrimeResponse{
				Result: divisor,
			})
			number /= divisor
		} else {
			divisor++
		}
	}
	return nil
}

func (s *Server) Avg(stream pb.CalculatorService_AvgServer) error {
	log.Printf("Avg function was invoked")

	var sum int64 = 0
	count := 0

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&pb.AvgResponse{
				Result: float64(sum) / float64(count),
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream %v\n", err)
		}
		log.Printf("receiving number %v\n", req.Number)
		sum += req.Number
		count++
	}
}

func (s *Server) Max(stream pb.CalculatorService_MaxServer) error {
	log.Printf("Max function was invoked")
	var maximum int32 = 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream %v\n", err)
		}
		if number := req.Number; number > maximum {
			maximum = number
			err := stream.Send(&pb.MaxResponse{Result: maximum})
			if err != nil {
				log.Fatalf("Error while sending data to client %v\n", err)
			}
		}
	}
}

func (s *Server) Sqrt(ctx context.Context, in *pb.SqrtRequest) (*pb.SqrtResponse, error) {
	log.Printf("Sqrt function was invoked with %v\n", in)

	number := in.Number
	if number < 0 {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("received a negative number %d\n", number),
		)
	}
	return &pb.SqrtResponse{
		Result: math.Sqrt(float64(number)),
	}, nil
}

package main

import (
	"context"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/MahbubAhmedTonmoy/GoGrpc/movie/proto"
)

var add string = "localhost:50051"

func main() {
	conn, err := grpc.Dial(add, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("failed to connecting %v\n", err)
	}
	defer conn.Close()

	c := pb.NewMovieServiceClient(conn)
	//createMovie(c)
	//getMovieById(c, "626f5af19cbc3e55dca64619")
	//updateMovide(c, "626f5af19cbc3e55dca64619")
	//getMovides(c)
	deleteMovide(c, "626f5af19cbc3e55dca64619")
}
func createMovie(c pb.MovieServiceClient) string {
	log.Printf("-------create movie was invoked--------")
	movie := pb.Movie{
		ProducerId: "Devgun",
		Title:      "test",
		Content:    "Test terst tesrtt esda",
	}
	res, err := c.CreateMovie(context.Background(), &movie)
	if err != nil {
		log.Fatalf("Unexpected error: %v\n", err)
	}
	log.Printf("movie created %s\n", res.Id)
	return res.Id
}

func getMovieById(c pb.MovieServiceClient, id string) *pb.Movie {
	log.Printf("-------getMovieById was invoked--------")

	req := &pb.MovieId{Id: id}

	res, err := c.GetMovie(context.Background(), req)
	log.Printf("response", res)
	if err != nil {
		log.Printf("Error happened while reading: %v\n", err)
	}
	return res
}
func updateMovide(c pb.MovieServiceClient, id string) {
	log.Printf("-------updateMovide was invoked--------")
	movie := &pb.Movie{
		Id:         id,
		ProducerId: "update Devgun",
		Title:      "update test",
		Content:    "update Test terst tesrtt esda",
	}

	_, err := c.UpdateMovie(context.Background(), movie)
	if err != nil {
		log.Printf("Error happened while reading: %v\n", err)
	}
}
func getMovides(c pb.MovieServiceClient) {
	log.Printf("-------getMovides was invoked--------")

	stream, err := c.ListMovie(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Printf("Error happened while reading: %v\n", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v\n", err)
		}
		log.Println(res)
	}
}

func deleteMovide(c pb.MovieServiceClient, id string) {
	log.Printf("-------deleteMovide was invoked--------")
	movie := &pb.MovieId{
		Id: id,
	}

	_, err := c.DeleteMovie(context.Background(), movie)
	if err != nil {
		log.Printf("Error happened while reading: %v\n", err)
	}
	log.Printf("movie deleted")
}

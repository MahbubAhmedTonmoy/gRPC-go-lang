package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/MahbubAhmedTonmoy/GoGrpc/movie/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var address string = "0.0.0.0:50051"
var collection *mongo.Collection

type Server struct {
	pb.MovieServiceServer
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("MovieDb").Collection("movie")

	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("Failed to listen on: %v\n", err)
	}
	log.Printf("Listening on %s\n", address)

	s := grpc.NewServer()
	pb.RegisterMovieServiceServer(s, &Server{})
	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v\n", err)
	}
}

type MovieItem struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	ProducerId string             `bson:"producer_id"`
	Title      string             `bson:"title"`
	Content    string             `bson:"content"`
}

func documentToMovie(data *MovieItem) *pb.Movie {
	return &pb.Movie{
		Id:         data.Id.Hex(),
		ProducerId: data.ProducerId,
		Title:      data.Title,
		Content:    data.Content,
	}
}

func (s *Server) CreateMovie(ctx context.Context, in *pb.Movie) (*pb.MovieId, error) {
	log.Printf("Create movie was invoked with %v\n", in)
	data := MovieItem{
		ProducerId: in.ProducerId,
		Title:      in.Title,
		Content:    in.Content,
	}

	response, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v\n", err))
	}
	oid, ok := response.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("can not convert oid: %v\n", err))
	}
	return &pb.MovieId{
		Id: oid.Hex(),
	}, nil
}

func (s *Server) GetMovie(ctx context.Context, in *pb.MovieId) (*pb.Movie, error) {
	log.Printf("GetMovie movie was invoked with %v\n", in)
	oid, err := primitive.ObjectIDFromHex(in.Id)

	if err != nil {
		return nil, status.Error(
			codes.Internal,
			"Can not parse ID",
		)
	}
	data := &MovieItem{}
	filter := bson.M{"_id": oid}
	res := collection.FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Error(
			codes.NotFound,
			"Can not find movie with id provided",
		)
	}
	return documentToMovie(data), nil
}

func (s *Server) UpdateMovie(ctx context.Context, in *pb.Movie) (*emptypb.Empty, error) {
	log.Printf("Updte movie was invoked with %v\n", in)

	oid, err := primitive.ObjectIDFromHex(in.Id)

	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Cannot parse ID",
		)
	}
	data := &MovieItem{
		ProducerId: in.ProducerId,
		Title:      in.Title,
		Content:    in.Content,
	}
	res, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$set": data},
	)
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			"Cannot update",
		)
	}

	if res.MatchedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			"Cannot found ID",
		)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) ListMovie(in *emptypb.Empty, stream pb.MovieService_ListMovieServer) error {
	log.Printf("Get all invocked")
	cur, err := collection.Find(context.Background(), primitive.D{{}})
	if err != nil {
		return status.Errorf(
			codes.Internal,
			"Unknown internal error",
		)
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		data := &MovieItem{}
		err := cur.Decode(data)

		if err != nil {
			return status.Errorf(
				codes.Internal,
				"error while decode data",
			)
		}
		stream.Send(documentToMovie(data))
	}
	if err = cur.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("unknown error %v\n", err),
		)
	}
	return nil
}

func (s *Server) DeleteMovie(ctx context.Context, in *pb.MovieId) (*emptypb.Empty, error) {
	log.Printf("delete movie was invoked with %v\n", in)

	oid, err := primitive.ObjectIDFromHex(in.Id)

	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Cannot parse ID",
		)
	}
	res, err := collection.DeleteOne(
		ctx,
		bson.M{"_id": oid},
	)
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			"Cannot update",
		)
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			"Cannot found ID",
		)
	}
	return &emptypb.Empty{}, nil
}

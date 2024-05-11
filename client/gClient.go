package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"protoapi"
	"time"

	"google.golang.org/grpc"
)

var port = ":8010"

func AskingDateTime(ctx context.Context, m protoapi.RandomClient) (*protoapi.DateTime, error) {
	req := &protoapi.RequestDateTime{
		Value: "Please send me the date and time",
	}

	return m.GetDate(ctx, req)
}

func AskPass(ctx context.Context, m protoapi.RandomClient, seed int64, length int64) (*protoapi.RandomPass, error) {
	req := &protoapi.RequestPass{
		Seed:   seed,
		Length: length,
	}

	return m.GetRandomPass(ctx, req)
}

func AskRandom(ctx context.Context, m protoapi.RandomClient, seed int64, place int64) (*protoapi.RandomInt, error) {
	req := &protoapi.RandomParams{
		Seed:  seed,
		Place: place,
	}

	return m.GetRandom(ctx, req)
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Using default port:", port)
	} else {
		port = os.Args[1]
	}

	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Dial:", err)
		return
	}

	rand.Seed(time.Now().Unix())
	seed := int64(rand.Intn(100))

	//return date
	client := protoapi.NewRandomClient(conn)
	r, err := AskingDateTime(context.Background(), client)
	assertError(err)
	fmt.Println("Server Date and Time:", r.Value)

	//random pass
	length := int64(rand.Intn(20))
	p, err := AskPass(context.Background(), client, 100, length+1)
	assertError(err)
	fmt.Println("Random password:", p.Password)

	//random integer 1
	place := int64(rand.Intn(100))
	i, err := AskRandom(context.Background(), client, seed, place)
	assertError(err)
	fmt.Println("Random integer 1:", i.Value)

	//random integer 2
	k, err := AskRandom(context.Background(), client, seed, place-1)
	assertError(err)
	fmt.Println("Random integer 2:", k.Value)

}

func assertError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

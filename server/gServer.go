package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"protoapi"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var min = 0
var max = 100

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func getString(len int64) string {
	temp := ""
	startChar := "!"
	var i int64 = 1

	for {
		myRand := random(0, 94)
		newChar := string(startChar[0] + byte(myRand))
		temp = temp + newChar

		if i == len {
			break
		}
		i++
	}
	return temp
}

type RandomServer struct {
	protoapi.UnimplementedRandomServer
}

func (RandomServer) GetDate(ctx context.Context, r *protoapi.RequestDateTime) (*protoapi.DateTime, error) {
	currentTime := time.Now()
	res := &protoapi.DateTime{
		Value: currentTime.String(),
	}

	return res, nil
}

func (RandomServer) GetRandom(ctx context.Context, r *protoapi.RandomParams) (*protoapi.RandomInt, error) {
	rand.Seed(r.GetSeed())
	place := r.GetPlace()
	temp := random(min, max)

	for {
		place--
		if place <= 0 {
			break
		}
		temp = random(min, max)
	}

	res := &protoapi.RandomInt{
		Value: int64(temp),
	}

	return res, nil
}

func (RandomServer) GetRandomPass(ctx context.Context, r *protoapi.RequestPass) (*protoapi.RandomPass, error) {
	rand.Seed(r.GetSeed())
	temp := getString(r.GetLength())

	res := &protoapi.RandomPass{
		Password: temp,
	}

	return res, nil
}

var port = ":8010"

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Using default port:", port)
	} else {
		port = os.Args[1]
	}

	server := grpc.NewServer()

	var randomServer RandomServer

	protoapi.RegisterRandomServer(server, randomServer)

	reflection.Register(server)

	listener, err := net.Listen("tcp", port)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Serving request...")
	server.Serve(listener)
}

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	skip "github.com/zackarysantana/goskip"
	"github.com/zackarysantana/goskip/examples"
)

type UsersValue struct {
	Name    string `json:"name"`
	Active  bool   `json:"active"`
	Friends []int  `json:"friends"`
}

type GroupsValue struct {
	Name    string `json:"name"`
	Members []int  `json:"members"`
}

func CreateClients(ctx context.Context) (skip.ControlClient, skip.StreamClient, func(), error) {
	shutdown, err := examples.StartSkipContainer(ctx, "examples/groups/skip.ts")
	if err != nil {
		return nil, nil, nil, err
	}

	controlClient := skip.NewControlClient(os.Getenv("SKIP_CONTROL_URL"))
	streamClient := skip.NewStreamClient(os.Getenv("SKIP_STREAM_URL"))

	return controlClient, streamClient, shutdown, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	controlClient, streamClient, shutdown, err := CreateClients(ctx)
	if err != nil {
		panic(err)
	}
	defer shutdown()

	go func() {
		uuid, err := controlClient.CreateResourceInstance(ctx, "active_friends", 0)
		if err != nil {
			panic(err)
		}

		err = streamClient.Stream(ctx, string(uuid), skip.ReadStream(func(event skip.StreamType, data []skip.CollectionValue[float64, float64]) error {
			fmt.Printf("Received Event: %s, Data: %v\n", event, data)
			return nil
		}))
		if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			panic(err)
		}
	}()
	time.Sleep(1 * time.Second)

	snapshot, err := skip.ReadResourceSnapshot[float64, float64](controlClient.GetResourceSnapshot(ctx, "active_friends", 0))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Resource snapshot: '%v'\n", snapshot)

	key, err := skip.ReadResourceKey[float64](controlClient.GetResourceKey(ctx, "active_friends", 1001, 0))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Resource key: %v\n", key)

	fmt.Println("Setting Carol to active")
	err = controlClient.UpdateInputCollection(ctx, "users", []skip.CollectionData{
		{
			Key: 2,
			Values: skip.Values(
				UsersValue{
					Name:    "Carol",
					Active:  true,
					Friends: []int{0, 1},
				},
			),
		},
	})
	if err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)

	fmt.Println("Setting Alice to inactive")
	err = controlClient.UpdateInputCollection(ctx, "users", []skip.CollectionData{
		{
			Key: 1,
			Values: skip.Values(
				UsersValue{
					Name:    "Alice",
					Active:  false,
					Friends: []int{0, 2},
				},
			),
		},
	})
	if err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)

	fmt.Println("Setting Eve as Bob's friend")
	err = controlClient.UpdateInputCollection(ctx, "users", []skip.CollectionData{
		{
			Key: 0,
			Values: skip.Values(
				UsersValue{
					Name:    "Bob",
					Active:  true,
					Friends: []int{1, 2, 3},
				},
			),
		},
	})
	if err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)

	fmt.Println("Removing Carol and adding Eve to group 2")
	err = controlClient.UpdateInputCollection(ctx, "groups", []skip.CollectionData{
		{
			Key: 1002,
			Values: skip.Values(
				GroupsValue{
					Name:    "Group 2",
					Members: []int{0, 3},
				},
			),
		},
	})
	if err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Second)
}

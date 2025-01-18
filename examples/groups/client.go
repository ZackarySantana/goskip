package main

import (
	"context"
	"fmt"
	"time"

	skip "github.com/zackarysantana/goskip"
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

func main() {
	ctx := context.Background()
	controlClient := skip.NewControlClient("http://localhost:8081/v1")
	streamClient := skip.NewStreamingClient("http://localhost:8080/v1")

	go func() {
		uuid, err := controlClient.CreateResourceInstance(ctx, "active_friends", 0)
		if err != nil {
			panic(err)
		}

		err = streamClient.StreamData(ctx, string(uuid), func(event skip.StreamType, data []skip.CollectionUpdate) {
			fmt.Printf("Received Event: %s, Data: %v\n", event, data)
		})
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(1 * time.Second)

	data, err := controlClient.GetResourceSnapshot(ctx, "active_friends", 0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Resource snapshot: %s\n", string(data))

	data, err = controlClient.GetResourceKey(ctx, "active_friends", 1001, 0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Resource key: %s\n", string(data))

	fmt.Println("Setting Carol to active")
	err = controlClient.UpdateInputCollection(ctx, "users", []skip.CollectionUpdate{
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
	err = controlClient.UpdateInputCollection(ctx, "users", []skip.CollectionUpdate{
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
	err = controlClient.UpdateInputCollection(ctx, "users", []skip.CollectionUpdate{
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
	err = controlClient.UpdateInputCollection(ctx, "groups", []skip.CollectionUpdate{
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

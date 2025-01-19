# Groups Example

This example copies Skip's groups and adapts it to Go.

## Skip's Example

-   [groups.ts](https://github.com/SkipLabs/skip/blob/main/skipruntime-ts/examples/groups.ts) defines the `skip` service.
-   [groups-server.ts](https://github.com/SkipLabs/skip/blob/main/skipruntime-ts/examples/groups-server.ts) is a layer over the `skip` service.
-   [groups-client.ts](https://github.com/SkipLabs/skip/blob/main/skipruntime-ts/examples/groups-client.ts) queries the server to change and read the state.

## goskip's Example

-   [skip](./skip.ts) defines the `skip` service.
-   [client.go](./client.go) queries the skip service to change and read the state.

### Running the Example

The example uses test containeres, so the only tools you need installed are `go` and `docker`.

To run the example, run the following command:

```bash
go run examples/groups/client.go
```

You should see the output like:

```bash
~/goskip$ go run examples/groups/client.go
Received Event: init, Data: [{Key:1001 Values:[1]}]
Setting Carol to active
Received Event: update, Data: [{Key:1001 Values:[1 2]} {Key:1002 Values:[2]}]
Setting Alice to inactive
Received Event: update, Data: [{Key:1001 Values:[2]}]
Setting Eve as Bob's friend
Received Event: update, Data: [{Key:1001 Values:[2 3]} {Key:1002 Values:[2]}]
Removing Carol and adding Eve to group 2
Received Event: update, Data: [{Key:1002 Values:[3]}]
```

Running the `client.go` file again will result in no updates, as the objects are already in that state.

```bash
~/goskip$ go run examples/groups/client.go
Received Event: init, Data: [{Key:1001 Values:[2 3]} {Key:1002 Values:[3]} {Key:1002 Values:[3]}]
Setting Carol to active
Setting Alice to inactive
Setting Eve as Bob's friend
Removing Carol and adding Eve to group 2
```

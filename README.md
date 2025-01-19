# goskip &middot; [![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/ZackarySantana/goskip/blob/main/LICENSE)

goskip is an unoffical open-source client for [Skip](https://github.com/SkipLabs/skip).

## Usage

If you're unfamiliar with Skip, you can find the documentation on the [Skip website](https://skiplabs.io/docs/introduction).

This repository is a wrapper for the skip API. It exposes two clients, `ControlClient` and `StreamClient`. These clients are used across collections and resources and then read functions are used to parse the underlying data.

To get started, create the clients:

```go
controlClient := skip.NewControlClient("<control_url>")
streamClient := skip.NewStreamClient("<stream_url>")
```

Then you can use the exported methods on the clients to call the API.

### Streaming Data

To create a resource instance:

```go
uuid, err := controlClient.CreateResourceInstance(ctx, "<resource_name>", <params>)
```

You can then use the `uuid` to stream data:

```go
err := streamClient.Stream(ctx, uuid, func(event skip.StreamType, data []byte) error {
    // handle untyped data
    return nil
})
```

The above example handles untyped data. If you want to handle typed data, you can use the `skip.ReadStream` function:

```go
err := streamClient.Stream(ctx, uuid, skip.ReadStream(func(event skip.StreamType, data []skip.CollectionValue[<key_type>, <value_type>]) error {
    // handle typed data
    return nil
}))
```

### Updating Data

To insert data, create a data type and use the `UpdateInputCollection` method:

```go
type DataValue struct {
    Name       string `json:"name"`
    DrankWater bool `json:"drank_water"`
}

err = controlClient.UpdateInputCollection(ctx, "<collection_name>", []skip.CollectionData{
    {
        Key: <key_value>,
        Values: skip.Values(
            DataValue{
                Name:    "Some Name",
                DrankWater: true,
            },
            DataValue{
                Name:    "Other Name",
                DrankWater: false,
            },
        ),
    },
})
```

### Snapshoting Data

Using Skip's API, you can snapshot a resource collection or an individual key in a collection:

```go
// Collection:
snapshot, err := skip.ReadResourceSnapshot[<key_type>, <value_type>](controlClient.GetResourceSnapshot(ctx, "<resource_name>", <params>))

// Individual Key:
key, err := skip.ReadResourceKey[<value_type>](controlClient.GetResourceKey(ctx, "<resource_name>", <resource_key>, <params>))
```

## Examples

The [examples](./examples) directory contains examples that have a `client.go` file and a `skip` directory. To run an example, run:

```bash
go run examples/<example>/client.go
```

These are the examples available:

-   [Groups](./examples/groups)

## Contributing

Contributions and pull requests are welcome! Feel free to drop an issue if you have any ideas or suggestions.

## License

goskip is [MIT licensed](./LICENSE).

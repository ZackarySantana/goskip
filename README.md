# goskip &middot; [![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/ZackarySantana/goskip/blob/main/LICENSE)

goskip is an unoffical open-source client for [skip](https://github.com/SkipLabs/skip).

## Installation

To get started, install the NPM packages for the Skip runtime API, server, and helpers:

`npm install @skipruntime/api @skipruntime/server @skipruntime/helpers`

From there, you're ready to start building a reactive service!
See the [getting started guide](https://skiplabs.io/docs/getting_started) to walk through some of Skip's core concepts by example and get up to speed.

## Documentation

[Godocs](https://pkg.go.dev/github.com/zackarysantana/goskip)

## Examples

The [examples](./examples) directory contains examples that have a `main.go` file and a `skip` directory. To run an example, open two terminals and run:

```bash
bun run examples/<example>/skip.ts
```

```bash
go run examples/<example>/main.go
```

These are the examples available:

-   [Groups](./examples/groups)

## Contributing

Contributions and pull requests are welcome! Feel free to drop an issue if you have any ideas or suggestions.

## License

goskip is [MIT licensed](./LICENSE).

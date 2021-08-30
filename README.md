# udpbench

`udpbench` is a UDP server benchmark tool.

Send a UUID to the server, then wait same UUID from server. (A server must return received UUID)

And then, `udpbench` prints total request count and duration.

## How to use

```shell
$ udpbench --help
Usage of ./udpbench:
      --address string    Server IP address (default "127.0.0.1")
      --count int         Number of request from a worker (default 10)
      --parallelism int   Worker parallelism number (default 10)
      --port int          Server port (default 8080)
pflag: help requested
```

## Example server

See [example](example)

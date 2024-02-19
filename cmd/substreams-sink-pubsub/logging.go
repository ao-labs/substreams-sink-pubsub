package main

import (
	"github.com/streamingfast/cli"
	"github.com/streamingfast/logging"
)

var zlog, tracer = logging.RootLogger("sink-kv", "github.com/streamingfast/substreams-sink-pubsub/cmd/substreams-sink-pubsub")

func init() {
	cli.SetLogger(zlog, tracer)
}

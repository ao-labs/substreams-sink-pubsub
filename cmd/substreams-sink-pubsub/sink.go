package main

import (
	pubsub "cloud.google.com/go/pubsub"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	. "github.com/streamingfast/cli"
	"github.com/streamingfast/cli/sflags"
	"github.com/streamingfast/shutter"
	sink "github.com/streamingfast/substreams-sink"
	spubsub "github.com/streamingfast/substreams-sink-pubsub"
)

var sinkCmd = Command(sinkRunE,
	"sink <manifest-path> <module-name> <topic-name> [<block-range>]",
	"Substreams Pubsub sinking",
	Flags(func(flags *pflag.FlagSet) {
		sink.AddFlagsToSet(flags)

		flags.String("cursor_path", "/tmp/sink-state/", "Sink cursor's path")
		flags.String("project", "", "Google Cloud Project ID")
		flags.StringP("endpoint", "e", "", "Substreams gRPC endpoint (e.g. 'mainnet.eth.streamingfast.io:443')")
	}),
	Description(`
		Publishs block data on a google PubSub from a Substreams output. 

		The required arguments are:
		- <manifest-path>: URL or local path to a '.yaml' file (e.g. './examples/simple/substreams.yaml').
		- <moduleName>: The module name returning publish instructions in the substreams.
		- <topicName>: The PubSub topic name to publish the messages to.
		
		The optional arguments are:
		- <start>:<stop>: The range of block to sync, if not provided, will sync from the module's initial block and then forever.

		If <start>:<stop> is not provided, assumes the whole chain.
	`),
	ExamplePrefixed("substreams-sink-pubsub sink", `
		# Publish block data messages produced by map_clocks for the whole chain
		-e mainnet.eth.streamingfast.io:443 ./examples/simple/substreams.yaml map_clocks "topic" --project "1"
		# Publish block data messages produced by map_clocks for a specific range of blocks
		-e mainnet.eth.streamingfast.io:443 ./examples/simple/substreams.yaml map_clocks "topic" 0:1000 --project "1"
	`),
)

func sinkRunE(cmd *cobra.Command, args []string) error {
	app := shutter.New()
	ctx := cmd.Context()

	manifestPath, module, topicName, blockRange := extractInjectArgs(cmd, args)
	endpoint := sflags.MustGetString(cmd, "endpoint")
	cursorPath := sflags.MustGetString(cmd, "cursor_path")
	projectID := sflags.MustGetString(cmd, "project")

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("creating pubsub client: %w", err)
	}

	topic := client.Topic(topicName)
	topic.EnableMessageOrdering = true

	sinker, err := sink.NewFromViper(
		cmd,
		"sf.substreams.sink.pubsub.v1.Publish",
		endpoint, manifestPath, module, blockRange,
		zlog, tracer,
	)
	if err != nil {
		return fmt.Errorf("unable to setup sinker: %w", err)
	}

	s := spubsub.NewSink(sinker, zlog, cursorPath, client, topic)

	s.OnTerminating(func(err error) {
		if err != nil {
			app.Shutdown(err)
			return
		}
	})

	app.OnTerminating(func(err error) {
		s.Shutdown(err)
	})

	s.Run(ctx)

	return nil
}

func extractInjectArgs(_ *cobra.Command, args []string) (manifestPath, moduleName, topicName, blockRange string) {
	manifestPath = args[0]
	moduleName = args[1]

	if len(args) >= 3 {
		topicName = args[2]
	}

	if len(args) == 4 {
		blockRange = args[3]
	}
	return
}

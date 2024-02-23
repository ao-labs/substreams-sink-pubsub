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
	"sink <endpoint> <projectId> <topicName> <moduleName> [<manifest-path>] [<block-range>]",
	"Substreams Pubsub sinking",
	Flags(func(flags *pflag.FlagSet) {
		sink.AddFlagsToSet(flags)

		flags.String("cursor_path", "/tmp/sink-state/", "Sink cursor's path")

	}),
	Description(`
		Publishs block data on a google PubSub from a Substreams output. 

		The required arguments are:
		- <endpoint>: The Substreams endpoint to reach (e.g. 'mainnet.eth.streamingfast.io:443').
		- <projectId>: The Google Cloud project ID.
		- <topicName>: The PubSub topic name on which data will be published.
		- <moduleName>: The module name returning publish instructions in the substreams.
		
		The optional arguments are:
		- <manifest>: URL or local path to a '.yaml' file (e.g. './examples/pubsub_substream/substreams.yaml').
		- <start>:<stop>: The range of block to sync, if not provided, will sync from the module's initial block and then forever.

		If the <manifest> is not provided, assume '.' contains a Substreams project to run. If
		<start>:<stop> is not provided, assumes the whole chain.
	`),
	ExamplePrefixed("substreams-sink-pubsub sink", `
		# Publish block data messages produced by map_clocks for the whole chain
		mainnet.eth.streamingfast.io:443 1 topic map_clocks ./examples/pubsub_substream/substreams.yaml

		# Publish block data messages produced by map_clocks for a specific range of blocks
		mainnet.eth.streamingfast.io:443 1 topic map_clocks ./examples/pubsub_substream/substreams.yaml 0:100000
	`),
)

func sinkRunE(cmd *cobra.Command, args []string) error {
	app := shutter.New()
	ctx := cmd.Context()

	endpoint, projectID, topicName, module, manifestPath, blockRange := extractInjectArgs(cmd, args)
	cursorPath := sflags.MustGetString(cmd, "cursor_path")

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

func extractInjectArgs(_ *cobra.Command, args []string) (endpoint, projectId, topicName, moduleName, manifestPath, blockRange string) {
	endpoint = args[0]
	projectId = args[1]
	topicName = args[2]
	moduleName = args[3]

	if len(args) >= 5 {
		manifestPath = args[4]
	}

	if len(args) == 6 {
		blockRange = args[5]
	}
	return
}

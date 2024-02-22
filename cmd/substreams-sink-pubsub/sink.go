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
	"sink",
	"Substreams Pubsub sink",
	Flags(func(flags *pflag.FlagSet) {
		sink.AddFlagsToSet(flags)

		flags.String("module", "", "An explicit module to sink, if not provided, expecting the Substreams manifest to defined 'sink' configuration")
		flags.String("cursor_path", "/tmp/sink-state/", "Sink cursor's path")
		flags.String("topic_name", "topic", "Pubsub topic name")
		flags.String("project_id", "1", "Pubsub project id")

	}),
)

func sinkRunE(cmd *cobra.Command, args []string) error {
	app := shutter.New()
	ctx := cmd.Context()

	endpoint, manifestPath, blockRange := extractInjectArgs(cmd, args)
	module := sflags.MustGetString(cmd, "module")
	cursorPath := sflags.MustGetString(cmd, "cursor_path")
	topicName := sflags.MustGetString(cmd, "topic_name")
	projectID := sflags.MustGetString(cmd, "project_id")

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("creating pubsub client: %w", err)
	}

	topic := client.Topic(topicName)
	topic.EnableMessageOrdering = true

	mapTopics := make(map[string]*pubsub.Topic)
	mapTopics[topicName] = topic

	sinker, err := sink.NewFromViper(
		cmd,
		"pubsub.v1.PublishOperations",
		endpoint, manifestPath, module, blockRange,
		zlog, tracer,
	)
	if err != nil {
		return fmt.Errorf("unable to setup sinker: %w", err)
	}

	s := spubsub.NewSink(sinker, zlog, tracer, cursorPath, client, mapTopics)

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

func extractInjectArgs(_ *cobra.Command, args []string) (endpoint, manifestPath, blockRange string) {
	endpoint = args[0]
	if len(args) >= 2 {
		manifestPath = args[1]
	}

	if len(args) == 3 {
		blockRange = args[2]
	}
	return
}

package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/streamingfast/cli"
	"github.com/streamingfast/cli/sflags"
	"github.com/streamingfast/shutter"
	sink "github.com/streamingfast/substreams-sink"
	pubsub "github.com/streamingfast/substreams-sink-pubsub"
)

var sinkCmd = cli.Command(sinkRunE, "sink", "Substreams Pubsub sink")

func sinkRunE(cmd *cobra.Command, args []string) error {
	app := shutter.New()
	ctx := cmd.Context()

	endpoint, manifestPath, blockRange := extractInjectArgs(cmd, args)
	module := sflags.MustGetString(cmd, "module")

	sinker, err := sink.NewFromViper(
		cmd,
		"need to set this",
		endpoint, manifestPath, module, blockRange,
		zlog, tracer,
	)
	if err != nil {
		return fmt.Errorf("unable to setup sinker: %w", err)
	}

	s := pubsub.NewSink(sinker, zlog, tracer)

	s.OnTerminating(func(err error) {
		if err != nil {
			app.Shutdown(err)
			return
		}
	})

	app.OnTerminating(func(err error) {
		s.Shutdown(err)
	})

	go func() {
		s.Run(ctx)
	}()

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

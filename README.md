# Substreams sink to Google PubSub

Substreams sink for Google PubSub helps quickly and easily sink blockchain data using Substreams modules to a PubSub topic.

## Prerequisites

Before sinking any data to a PubSub, make sure to have the following prerequisites:

- **PubSub Creation**:
  Create a PubSub with a Google cloud projectID associated and a topic on which to publish the data.

- **Substreams creation**:
   A Substreams that has  a `map` module with an output type of `sf.substreams.sink.pubsub.v1.Publish` message https://github.com/streamingfast/substreams-sink-pubsub/blob/develop/proto/sf/substreams/sink/pubsub/v1/pubsub.proto.

## Installation

Install the `substreams-sink-pubsub` binary from source, by running the following command:

```bash
go install ./cmd/substreams-sink-pubsub
```

## Running

The `substreams-sink-pubsub` binary offers a sink tool. This sink tool sinks the data from the Substreams to the PubSub associated with your Google cloud `project_id`.
This is publishing all the block relative data depending on the Substreams module you are using, on a specified `topic_name`.

> [!NOTE]
> You can do `docker compose up` in the root of the repository to spin up a GCP PubSub local emulator to test out the sink easily, available on post 8888 by default. If you use the emulator, ensure to also do `export PUBSUB_EMULATOR_HOST=localhost:8888` in your terminal so the sink can correctly reach it.

Run the sink using:

```bash
substreams-sink-pubsub sink --project=acme https://raw.githubusercontent.com/streamingfast/substreams-sink-pubsub/refs/heads/develop/examples/simple/simple-v0.1.0.spkg map_clocks dev-topic 100000:+1000
```

> [!NOTE]
> Check `substreams-sink-pubsub sink --help` for full command description and options

### Examples

We provide two pre-built Substreams to use as example(s):
- [./examples/simple](./examples/simple/) a simple mapper that maps `sf.substreams.v1.Clock` so it works on any network
- [./examples/ethERC20Transfers](./examples/ethERC20Transfers/) an ERC20 transfers Substreams

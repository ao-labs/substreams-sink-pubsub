# Substreams sink to Google PubSub

Substreams sink for Google PubSub helps quickly and easily sync blockchain data using Substreams modules to a PubSub topic.

## Prerequisites

Before sinking any data to a PubSub, make sure to have the following prerequisites:

- **PubSub Creation**:
  Create a PubSub with a Google cloud projectID associated and a topic on which to publish the data.

- **Substreams creation**:
  - Use the `simple` provided in the [examples](./examples) directory or create your own substreams.
  - Compile the `simple` project (or your own substreams):

      ```bash
      cd examples/simple
      cargo build --target wasm32-unknown-unknown --release
      ```
    **Note:** *If you are creating your own substreams, make sure to create a `map` module with an output type of `sf.substreams.sink.pubsub.v1.Publish` message https://github.com/streamingfast/substreams-sink-pubsub/blob/develop/proto/sf/substreams/sink/pubsub/v1/pubsub.proto*

## Installation

Install the `substreams-sink-pubsub` binary from source, by running the following command:

```bash
go install ./cmd/substreams-sink-pubsub
```

## Running 

The `substreams-sink-pubsub` binary offers a sink tool. This sink tool sinks the data from the substreams to the PubSub associated with your Google cloud `project_id`. 
This is publishing all the block relative data depending on the substreams module you are using, on a specified `topic_name`. 

Run the sink providing the `substreams manifest` and the substreams `module name` (the one having the `map` module with an output type of sf.substreams.sink.pubsub.v1.Publish),
using the following command:

```bash 
substreams-sink-pubsub sink -e <endpoint> --project <project_id> <substreams_manifest> <substreams_module_name> <topic_name> 
```

**Note:** *--help flag can be used to get more information on the flags used in the sink command.*

## Example

As an example, let's sink the ethereum blockchain data from the `simple` module's named `map_clocks`, provided in the [examples](./examples) directory.

Run the following command, to publish the data on the PubSub topic `myTopic` associated with the Google cloud project-id `myProjectId`:

```bash
substreams-sink-pubsub sink -e mainnet.eth.streamingfast.io:443 --project "1" ./examples/simple/substreams.yaml map_clocks "topic"
```




